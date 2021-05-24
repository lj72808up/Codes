package models

import (
	"AStar/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"net/http"
	"sync"
	"time"
)

type Task struct {
	JonesDagName  string // jone 数据库中的 dag名称
	JonesModuleId string
	AirflowId     string // 提交后生成的 airflow 任务 id
	Img           string // 任务镜像名
	Cmd           string // 任务命令
	Env           string // 环境变量
	ScheduleTime  string // 执行的是哪个时间的任务
	Username      string // 提交任务的同用户
}

var wg sync.WaitGroup

func workFunction(triggerDag DagTriggerTask) {
	runChannel(triggerDag)
}

func worker(wg *sync.WaitGroup) {
	for dagTriggerTask := range ManualTriggerDagChannel {
		workFunction(dagTriggerTask)
	}
}

func CreateWorkerPool(noOfWorkers int) {
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Done()
}

func runChannel(triggerDag DagTriggerTask) {
	Logger.Info("触发手动Dag: %s, 时间: %s", triggerDag.DagName, triggerDag.ScheduleTime)
	hid := triggerDag.HistoryId
	nodes := triggerDag.NodeList // 一个 dag 一次调度的模块数量
	/////  start-  insert db
	err := UpdateTriggerHistory(triggerDag.HistoryId, "running", "")
	if err != nil {
		utils.Logger.Error("%s", err.Error())
	}
	/////  end-  insert db
	tasksInDagPerTime := dagNodeToTask(nodes, triggerDag.DagName, triggerDag.ScheduleTime, triggerDag.UserName)
	var runtimeErr error
	remoteIds := []map[string]string{}
	for _, t := range tasksInDagPerTime {
		// (1) 提交任务
		airflowId := CreateRemoteTask(t, t.ScheduleTime)
		remoteIds = append(remoteIds, map[string]string{
			"remoteId": airflowId+"",
			"modelName": t.JonesModuleId,
		})
		// (2) 轮训接口, 看是否执行完毕
		isFinish := false
		for ; (!isFinish) && (err == nil); {
			isFinish, runtimeErr = IsFinish(airflowId, triggerDag.UserName)
			time.Sleep(30 * time.Second)
		}
		if runtimeErr != nil {
			break
		}
	}

	/// start update db
	remoteIdStr,_ := json.Marshal(remoteIds)
	if runtimeErr != nil {
		if err := UpdateTriggerHistory(hid, "fail", string(remoteIdStr)); err != nil {
			utils.Logger.Error("%s", err.Error())
		}
	} else {
		if err := UpdateTriggerHistory(hid, "success", string(remoteIdStr)); err != nil {
			utils.Logger.Error("%s", err.Error())
		}
	}
	/// end update db
}

func dagNodeToTask(nodes []DagNode, dagName string, scheduleTime string, userName string) []Task {
	var tasks []Task
	dbUser, passwd, hostAndPort, dataBase := utils.GetServerDBInfo()
	for _, node := range nodes {
		// 将每个节点转边成任务加入channel
		modelName := node.Name
		env := fmt.Sprintf(`{
				"JONES_DAG_NAME":"%s",
				"JONES_JDBC":"jdbc:mysql://%s/%s",
				"JONES_JDBC_PASSWORD":"%s",
				"JONES_JDBC_USER":"%s",
				"JONES_TASK_NAME":"%s",
				"password":"adtd_platform",
				"user":"adtd_platform"
			}`, dagName, hostAndPort, dataBase, passwd, dbUser, modelName)

		// 放入每个任务
		tasks = append(tasks, Task{
			JonesDagName:  dagName,
			JonesModuleId: modelName,
			AirflowId:     "",
			Img:           "harbor.adtech.sogou/hadoop/hadoop:jones_platform",
			Cmd:           fmt.Sprintf("cd jones_platform;sh run.sh %s", scheduleTime),
			Env:           env,
			ScheduleTime:  scheduleTime,
			Username:      userName,
		})
	}
	return tasks
}

// time: 跑的哪个时间段的任务
func CreateRemoteTask(t Task, scheduleTime string) string {
	// http 请求体
	paramMap := map[string]string{
		"docker_image": t.Img,
		"run_cmd":      t.Cmd,
		"env_params":   t.Env,
	}
	params, _ := json.Marshal(paramMap)

	var cookie = http.Cookie{
		Name:  "_adtech_user",
		Value: t.Username,
	}

	req := httplib.Post("http://op.adtech.sogou/airflow/api/task/").
		SetTimeout(5*time.Minute, 5*time.Minute).SetCookie(&cookie).
		Body(params).Header("Content-Type", "application/json")

	bytes, err := req.Bytes()

	println(string(bytes))
	if err != nil {
		utils.Logger.Error("%s", err)
		return ""
	} else {
		resp := map[string]int{}
		_ = json.Unmarshal(bytes, &resp)
		//utils.Logger.Info("%d 开始执行",resp["task_id"])
		return fmt.Sprintf("%d", resp["task_id"])
	}
}

// 获取任务执行进度
func IsFinish(airflowId string, username string) (bool, error) {
	url := fmt.Sprintf("http://op.adtech.sogou/airflow/api/task/%s/", airflowId)
	var cookie = http.Cookie{
		Name:  "_adtech_user",
		Value: username,
	}
	req := httplib.Get(url).SetTimeout(5*time.Minute, 5*time.Minute).SetCookie(&cookie)
	bytes, err := req.Bytes()

	println(string(bytes))
	if err != nil {
		utils.Logger.Error("%s", err)
		return true, nil
	} else {
		status := arbitrateFinish(bytes)
		if status == "success" {
			return true, nil
		} else if status == "failure" {
			utils.Logger.Error("%s", "还任务执行失败, 后面不应该继续执行")
			return true, errors.New("模块执行失败, 后续不再执行")
		} else {
			return false, nil
		}
	}
}

func arbitrateFinish(bytes []byte) string {
	resp := map[string]string{}
	_ = json.Unmarshal(bytes, &resp)
	return resp["job_status"]
}

// 获取任务日志
func GetTaskLog(airflowId string, username string) string {
	url := "http://op.adtech.sogou/airflow/api/task_log?task_id=" + airflowId
	var cookie = http.Cookie{
		Name:  "_adtech_user",
		Value: username,
	}
	req := httplib.Get(url).SetTimeout(5*time.Minute, 5*time.Minute).
		SetCookie(&cookie)
	bytes, err := req.Bytes()

	if err != nil {
		utils.Logger.Error("%s", err.Error())
		return err.Error()
	} else {
		return string(bytes)
	}
}

type DagTriggerHistory struct {
	Hid      int64  `orm:"pk;auto" field:"历史Id"`
	DagName  string `orm:"size(50);index" field:"名称"`
	Username string `orm:"size(50);index" field:"提交用户"`
	//RemoteId     int    `orm:"index" field:"远程Id"` // 远端接口生成的id
	ScheduleTime string `orm:"size(50)" field:"调度时间"`
	Status       string `field:"状态"` // 手动触发的dag的状态: running, fail, success
	RemoteIds    string `orm:"size(255)" field:"remoteIds,hidden"`
}

func GetAllTriggerHistory(dagName string) []DagTriggerHistory {
	o := orm.NewOrm()
	histories := []DagTriggerHistory{}
	_, _ = o.Raw("SELECT * FROM dag_trigger_history WHERE dag_name = ? ORDER BY hid DESC", dagName).
		QueryRows(&histories)
	return histories
}

func (this *DagTriggerHistory) Insert() (int64, error) {
	o := orm.NewOrm()
	_, err := o.Insert(this)
	return this.Hid, err
}

func UpdateTriggerHistory(hid int64, status string, remoteIdStr string) error {
	o := orm.NewOrm()
	//ids := strings.Join(remoteIds, ",")
	_, err := o.Raw("UPDATE dag_trigger_history SET status = ? ,remote_ids = ? WHERE hid = ?", status, remoteIdStr, hid).Exec()
	return err
}
