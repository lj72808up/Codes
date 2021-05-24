package models

import (
	"AStar/models/airflowDag"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

//--------------------------------------------

type Airflow struct{}

type AirflowDag struct {
	DagId           string `orm:"pk;size(250)" `
	LastScheduleRun string `orm:"size(255)" field:"Last run"`
	Schedule        string `orm:"size(255)" field:"调度时间"`
	User            string `orm:"size(2000)" field:"所有者"`
}

type AirflowDagStats struct {
	Dagid        string `orm:"pk;size(250)" `
	StateSuccess int    `orm:"size(11)"`
	StateFailed  int    `orm:"size(11)"`
	StateRunning int    `orm:"size(11)"`
}

type AirflowTaskStats struct {
	Dagid                   string `orm:"pk;size(250)" `
	TaskStateSuc            int    `orm:"size(11)"`
	TaskStateRunning        int    `orm:"size(11)"`
	TaskStateFailed         int    `orm:"size(11)"`
	TaskStateUpFailed       int    `orm:"size(11)"`
	TaskStateSkipped        int    `orm:"size(11)"`
	TaskStateUpForRetry     int    `orm:"size(11)"`
	TaskStateUpForReshedule int    `orm:"size(11)"`
	TaskStateQueued         int    `orm:"size(11)"`
	TaskStateNull           int    `orm:"size(11)"`
	TaskStateSheduled       int    `orm:"size(11)"`
}

type AirflowDagShow struct {
	Name        string
	ExhibitName string
	IsOpen      int
	//LastScheduleRun       		 	string
	//LastScheduleStatus				string
	Schedule string
	User     string
	Status   int
	/*
		StateSuccess	      			int
		StateFailed        				int
		StateRunning     				int
		TaskStateSuc	        		int
		TaskStateRunning        		int
		TaskStateFailed     			int
		TaskStateUpFailed     			int
		TaskStateSkipped     			int
		TaskStateUpForRetry     		int
		TaskStateUpForReshedule     	int
		TaskStateQueued     			int
		TaskStateNull     				int
		TaskStateSheduled     			int

	*/
}

/*
// 获取所有的Dag信息
func   GetAllDagInfo() []AirflowDagShow {
	o := orm.NewOrm()
	var views []AirflowDagShow
	sql := `
			SELECT
				ad.dag_id AS dag_id,
				ad.is_open AS is_open,
				ad.last_schedule_run AS last_schedule_run,
				ad.schedule AS schedule,
				ad.owners AS owners,
				IFNULL(ads.state_success, 0) AS state_success,
				IFNULL(ads.state_failed, 0) AS state_failed,
				IFNULL(ads.state_running, 0) AS state_running,
				IFNULL(ats.task_state_suc, 0) AS task_state_suc,
				IFNULL(ats.task_state_running, 0) AS task_state_running,
				IFNULL(ats.task_state_failed, 0) AS task_state_failed,
				IFNULL(ats.task_state_up_failed, 0) AS task_state_up_failed,
				IFNULL(ats.task_state_skipped, 0) AS task_state_skipped,
				IFNULL(ats.task_state_up_for_retry,0) AS task_state_up_for_retry,
				IFNULL(ats.task_state_up_for_reshedule,0) AS task_state_up_for_reshedule,
				IFNULL(ats.task_state_queued, 0) AS task_state_queued,
				IFNULL(ats.task_state_null, 0) AS task_state_null,
				IFNULL(ats.task_state_sheduled, 0) AS task_state_sheduled
			FROM
				airflow_dag ad
			LEFT JOIN airflow_dag_stats ads ON ad.dag_id = ads.dag_id
			LEFT JOIN airflow_task_stats ats ON ad.dag_id = ats.dag_id
	`
	_, err := o.Raw(sql).QueryRows(&views)
	if err != nil {
		panic(err.Error())
	}
	return views
}
*/

// 获取所有的Dag信息
func GetAllDagInfo(isAll bool, user string) []AirflowDagShow {
	o := orm.NewOrm()
	var views []AirflowDagShow
	//增加是否有写的权限, 有则返回全部dag, 无则只返回自己创建的dag
	condSql := " where status!=? "
	if !isAll {
		condSql = fmt.Sprintf("%s AND user='%s' ", condSql, user)
	}
	sql := fmt.Sprintf(`
			SELECT
				ad.name AS name,
				ad.exhibit_name AS exhibit_name,
				ad.is_open AS is_open,
				ad.cron_expression AS schedule,
				ad.user AS user,
				ad.status AS status
			FROM
				dag_air_flow ad
			%s
			order by modify_time desc
	`, condSql)
	_, err := o.Raw(sql, airflowDag.REMOVED).QueryRows(&views)
	if err != nil {
		panic(err.Error())
	}
	return views
}

// 改变Dag开关标志位
func ChangeDagSwitchAPI(dagName string, isOpen string) bool {
	url := AirflowInfo.AirflowApiUrl + "/" + "experimental/dag_paused"

	//API 接口参数写的不好， is_paused 为 true 时是打开 dag
	is_paused := "false"
	if isOpen == "1" {
		is_paused = "true"
	}

	/*
		param, _ := json.Marshal(map[string]string{
			"dag_id": dagId,
			"is_paused":is_paused,
		})
	*/

	post_url := fmt.Sprintf("%s?dag_id=%s&is_paused=%s", url, dagName, is_paused)
	//req := httplib.Post(post_url).SetTimeout(time.Second * 5, time.Second * 5).Body(param)
	req := httplib.Post(post_url).SetTimeout(time.Second*5, time.Second*5)
	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
	}
	var res map[string]interface{}
	_ = json.Unmarshal(bytes, &res)

	res_code, exist := res["code"]

	if exist && int(res_code.(float64)) == 0 {
		return true
	}
	return false

}

func ChangeDagSwitchDB(dagName string, newStatus string) bool {
	flag := true
	o := orm.NewOrm()

	urlSql := fmt.Sprintf(" Update dag_air_flow set is_open='%s' where name='%s'", newStatus, dagName)
	res, err := o.Raw(urlSql).Exec()
	if err != nil {
		fmt.Printf("update dag_air_flow get err:%s | res:%s", err, res)
		flag = false
	}
	return flag
}

func DeleteDagFile(dagName string) bool {
	flag := true
	o := orm.NewOrm()

	dagFile := "/tmp/tmp.file"
	user := ""
	urlSql := fmt.Sprintf(" select user,file_name from dag_air_flow  where name='%s'", dagName)
	err := o.Raw(urlSql).QueryRow(&user, &dagFile)
	if err != nil {
		fmt.Printf("update dag_air_flow get err:%s | res:%s", err)
		flag = false
	}

	airbase := beego.AppConfig.String("airflowFileRoot")
	//dag_file_full_path := fmt.Sprintf("%s/%s/%s", airbase, utils.DateFormatWithoutLine(time.Now()), user)

	//删除本地文件
	rmlocalCmd := fmt.Sprintf("rm -f %s", dagFile)
	splits := strings.Split(rmlocalCmd, " ")
	utils.ExecShell(splits[0], splits[1:])

	//同步Dag目录到远端
	rsyncCmd1 := fmt.Sprintf("rsync -aP --delete %s %s", airbase, beego.AppConfig.String("airflowRemoteFileRoot1"))
	rsyncCmd2 := fmt.Sprintf("rsync -aP --delete %s %s", airbase, beego.AppConfig.String("airflowRemoteFileRoot2"))
	cmds := []string{rsyncCmd1, rsyncCmd2}
	for _, cmd := range cmds {
		fmt.Println("rsync 执行:", cmd)
		splits := strings.Split(cmd, " ")
		utils.ExecShell(splits[0], splits[1:])
	}
	return flag
}

// 创建中、删除中、已删除
func UpdateDagStatus(dagName string, newStatus int) bool {
	flag := true
	o := orm.NewOrm()

	res, err := o.Raw("Update dag_air_flow set status=? where name=?", newStatus, dagName).Exec()
	if err != nil {
		fmt.Printf("update dag_air_flow status get err:%s | res:%s\n", err.Error(), res)
		flag = false
	}
	return flag
}

// 删除Dag by API
func DeleteDagApi(dagName string) {
	url := AirflowInfo.AirflowApiUrl + "/" + "experimental/dag_delete"

	param, _ := json.Marshal(map[string]string{
		"dag_id": dagName,
	})

	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body(param)
	bytes, err := req.Bytes()
	if err != nil {
		Logger.Error(err.Error())
	}
	fmt.Println(string(bytes))
}

// 重新标志dag 调度状态 、清除、标志成功、标记失败
func MarkDagApi(dagName string, runTime string, newStatus string) bool {
	url := AirflowInfo.AirflowApiUrl + "/" + "experimental/"
	if newStatus == "notrun" {
		url += "dag_clear"
	} else if newStatus == "success" {
		url += "dag_mark_success"
	} else if newStatus == "failed" {
		url += "dag_mark_failed"
	}

	param, _ := json.Marshal(map[string]string{
		"dag_id":         dagName,
		"execution_date": runTime,
	})

	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body(param)
	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
	}
	var res map[string]interface{}
	_ = json.Unmarshal(bytes, &res)

	res_code, exist := res["code"]
	res_info, _ := res["msg"]
	if exist && int(res_code.(float64)) == 0 {
		return true
	}

	fmt.Printf("Delete Dag by Api get err:%s", res_info)
	return false
}

// 改变调度执行状态
func UpdateDagRunStatus(dagName string, runTime string, newStatus string) bool {
	flag := true
	o := orm.NewOrm()

	urlSql := fmt.Sprintf(" Update dag_air_flow set last_schedule_status='%s' where name='%s'", newStatus, dagName)
	res, err := o.Raw(urlSql).Exec()
	if err != nil {
		fmt.Printf("update dag_air_flow last_schedule_status get err:%s | res:%s", err, res)
		flag = false
	}
	return flag
}

// 重新标志Task 调度状态 、清除、标志成功、标记失败
func MarkTaskApi(dagName string, taskId string, runTime string, newStatus string) bool {
	url := AirflowInfo.AirflowApiUrl + "/" + "experimental/"
	if newStatus == "notrun" {
		url += "task_clear"
	} else if newStatus == "success" {
		url += "task_mark_success"
	} else if newStatus == "failed" {
		url += "task_mark_failed"
	}

	param, _ := json.Marshal(map[string]interface{}{
		"dag_id":         dagName,
		"task_id":        taskId,
		"execution_date": runTime,
		"downstream":     false,
		"upstream":       false,
	})

	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body(param)
	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
	}
	var res map[string]interface{}
	_ = json.Unmarshal(bytes, &res)

	res_code, exist := res["code"]
	res_info, _ := res["msg"]
	if exist && int(res_code.(float64)) == 0 {
		return true
	}

	fmt.Printf("MarkTaskApi get err:%s", res_info)
	return false
}
