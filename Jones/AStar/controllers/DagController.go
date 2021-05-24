package controllers

import (
	"AStar/models"
	"AStar/models/airflowDag/hive"
	"AStar/utils"
	"AStar/utils/Sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type DagController struct {
	BaseController
}

var DagNameSpace = "airflow"
var timeOut = 5 * time.Minute

// @Title 修改DAG任务
// @Summary 修改DAG任务
// @Description 修改DAG任务
// @Success 200
// @Accept json
// @router /updateDag [post]
func (this *DagController) UpdateDag() {
	dagName := this.GetString("dagName", "")
	var view models.DagView
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &view)
	defer func() {
		if r := recover(); r != nil {
			models.Logger.Error(fmt.Sprintf("%v", r))
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()
	view.UpdateDag(string(this.Ctx.Input.RequestBody), dagName)
	this.Ctx.WriteString(utils.VueResponse{
		Res:  true,
		Info: "创建成功",
	}.String())
}

// @Title 指派dag
// @Summary 指派dag
// @Description 指派dag
// @Success 200
// @Accept json
// @router /assignDag [post]
func (this *DagController) AssignDag() {
	var param map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &param)
	var assignUser = param["user"]
	var dagId = param["dagId"]
	o := orm.NewOrm()
	_, _ = o.Raw("UPDATE dag_air_flow SET user = ? WHERE name = ?", assignUser, dagId).Exec()
}

// @Title 创建DAG任务
// @Summary 创建DAG任务
// @Description 创建DAG任务
// @Success 200
// @Accept json
// @router /addDag [post]
func (this *DagController) PostDag() {
	var view models.DagView
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &view)
	defer func() {
		if r := recover(); r != nil {
			models.Logger.Error(fmt.Sprintf("%v", r))
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()
	view.PostDag(view, this.userName, string(this.Ctx.Input.RequestBody))
	this.Ctx.WriteString(utils.VueResponse{
		Res:  true,
		Info: "创建成功",
	}.String())
}

// @Title 创建DAG失败时进行保存
// @Summary 创建DAG失败时进行保存
// @Description 创建DAG失败时进行保存
// @Success 200
// @Accept json
// @router /saveDagWhenPostFail [post]
func (this *DagController) SaveDagWhenPostFail() {
	var view models.DagView
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &view)
	defer func() {
		if r := recover(); r != nil {
			models.Logger.Error(fmt.Sprintf("%v", r))
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()
	view.SaveDagWhenFail(this.userName, string(this.Ctx.Input.RequestBody))
	this.Ctx.WriteString(utils.VueResponse{
		Res:  true,
		Info: "保存成功",
	}.String())
}

// @Title 创建DAG任务前, 解析输出表, 自动创建表
// @Summary 创建DAG任务前, 解析输出表, 自动创建表
// @Description 创建DAG任务前, 解析输出表, 自动创建表
// @Success 200
// @Accept json
// @router /preAddDag [post]
func (this *DagController) PrePostDag() {
	var view models.DagView
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &view)
	defer func() {
		if r := recover(); r != nil {
			models.Logger.Error(fmt.Sprintf("%v", r))
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()
	roles := models.Enforcer.GetRolesForUser(this.userName)
	view.PrePostDag(roles, this.userName)
	this.Ctx.WriteString(utils.VueResponse{
		Res:  true,
		Info: "远程DAG输出表创建成功, 元数据已导入",
	}.String())
}

type TaskView struct {
	TaskId        string `json:"task_id" field:"任务"`
	SchedulerTime string `json:"scheduler_time" field:"调度批次"`
	ExecutionDate string `json:"execution_date" field:"执行批次"`
	StartDate     string `json:"start_date" field:"开始时间"`
	EndDate       string `json:"end_date" field:"结束时间"`
	TryNum        string `json:"try_num" field:"重试次数,hidden"`
	State         string `json:"state" field:"状态,hidden"`
	Epoch         int    `json:"epoch" field:"轮数,hidden"`
}

func (this *TaskView) MakeFormatTime() {
	toFormat := "2006-01-02 15:04:05"
	this.StartDate = utils.TransDateByStandard(this.StartDate, time.RFC1123, toFormat)
	this.EndDate = utils.TransDateByStandard(this.EndDate, time.RFC1123, toFormat)
	this.SchedulerTime = utils.TransDateByStandard(this.SchedulerTime, time.RFC1123, toFormat)
	//this.ExecutionDate = utils.TransDateByStandard(this.ExecutionDate, time.RFC3339, toFormat)
}

type JobsInstance struct {
	Data   []TaskView `json:"data"`
	Msg    string     `json:"msg"`
	Status int64      `json:"status"`
}

// @Title 获取dag任务中job的详情
// @Summary 获取dag任务中job的详情
// @Description 获取dag任务中job的详情
// @Success 200
// @Accept json
// @router /getDetailByDagId [post]
func (this *DagController) GetTaskDetail() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	dagId := attr["dagId"]
	baseDate := attr["baseDate"]
	cnt, _ := strconv.Atoi(attr["cnt"])

	taskId := attr["taskId"]
	url := fmt.Sprintf("http://api.airflow.adtech.sogou/api/experimental/dagrun_list?dag_id=%s&nums=%d&base_date=%s", dagId, cnt, baseDate)
	if taskId != "" {
		url = fmt.Sprintf("%s&task_id=%s", url, taskId)
	}
	req := httplib.Get(url).SetTimeout(timeOut, timeOut)
	bytes, _ := req.Bytes()
	var jobs JobsInstance
	_ = json.Unmarshal(bytes, &jobs)
	tasks := jobs.Data
	// dag中的task个数
	url2 := fmt.Sprintf("http://api.airflow.adtech.sogou/api/experimental/dagrun_list?dag_id=%s&nums=1&base_date=%s", dagId, baseDate)
	req2 := httplib.Get(url2).SetTimeout(timeOut, timeOut)
	bytes2, _ := req2.Bytes()
	var jobs2 JobsInstance
	_ = json.Unmarshal(bytes2, &jobs2)
	taskCnt := len(jobs2.Data) // dag中任务的个数
	for i, _ := range tasks {
		tasks[i].MakeFormatTime()
		tasks[i].Epoch = i / taskCnt
	}
	/////// 按照调度轮次从最近到最远排序
	i := 0
	j := len(tasks) - taskCnt
	epoch := 1 // 轮数

	for i < j {
		for increase := 0; increase < taskCnt; increase++ {
			tmp := tasks[i]
			tasks[i] = tasks[j]
			tasks[j] = tmp
			i = i + 1
			j = j + 1
		}
		epoch = epoch + 1
		j = len(tasks) - taskCnt*epoch
	}
	///////

	// transView
	var bodys []map[string]interface{}
	for _, data := range tasks {
		// 转成map
		tmpData := utils.GetSnakeStrMapOfStruct(data)
		bodys = append(bodys, tmpData)
	}

	// 返回数据
	tsb := utils.TableData{Data: bodys}
	tsb.GetSnakeTitles(TaskView{})
	res, _ := json.Marshal(tsb)

	this.Ctx.WriteString(string(res))
}

// @Title 获取1个dag中的task个数
// @Summary 获取1个dag中的task个数
// @Description 获取1个dag中的task个数
// @Success 200
// @Accept json
// @router /getTaskCntInOneDag/:dagId [get]
func (this *DagController) GetTaskCntInOneDag() {
	dagId := this.Ctx.Input.Params()[":dagId"]
	batchCnt := 1
	url := fmt.Sprintf("http://api.airflow.adtech.sogou/api/experimental/dagrun_list?dag_id=%s&nums=%d", dagId, batchCnt)
	req := httplib.Get(url).SetTimeout(timeOut, timeOut)
	bytes, _ := req.Bytes()
	var jobs JobsInstance
	_ = json.Unmarshal(bytes, &jobs)
	taskCount := len(jobs.Data)
	res, _ := json.Marshal(map[string]int{
		"taskCount": taskCount,
	})
	this.Ctx.WriteString(string(res))
}

// @Title 获取所有的镜像
// @Summary 获取所有的镜像
// @Description 获取所有的镜像
// @Success 200
// @Accept json
// @router /getAllDagImgs [get]
func (this *DagController) GetAllDagImgs() {
	imgs := models.GetAllDagImgs()
	res, _ := json.Marshal(imgs)
	this.Ctx.WriteString(string(res))
}

// @Title 检测dag名称是否重名
// @Summary 检测dag名称是否重名
// @Description 检测dag名称是否重名
// @Success 200
// @Accept json
// @router /checkDagName/:dagId [get]
func (this *DagController) CheckDagName() {
	dagId := this.Ctx.Input.Params()[":dagId"]
	cnt := 1
	url := fmt.Sprintf("http://api.airflow.adtech.sogou/api/experimental/dagrun_list?dag_id=%s&nums=%d", dagId, cnt)

	req := httplib.Get(url).SetTimeout(timeOut, timeOut)
	bytes, _ := req.Bytes()
	var jobs JobsInstance
	_ = json.Unmarshal(bytes, &jobs)
	tasks := jobs.Data
	if len(tasks) > 0 {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "重名了",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: "没重名",
		}.String())
	}
}

// @Title 获取Dag任务详情
// @Summary 获取Dag任务详情
// @Description 获取Dag任务详情
// @Success 200
// @Accept json
// @router /getDag [get]
func (this *DagController) GetDagByName() {
	dagName := this.GetString("dagName")
	var dag models.DagAirFlow
	dag.GetByName(dagName)
	this.Ctx.WriteString(dag.Description)
}

type GetSql struct {
	Sql    string              `json:"sql"`
	Params []map[string]string `json:"params"`
}

// @Title 获取sql的字段
// @Summary 获取sql的字段
// @Description 获取sql的字段
// @Success 200
// @Accept json
// @router /getSqlFields [post]
func (this *DagController) GetSqlFields() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	var sqlParam []map[string]string
	_ = json.Unmarshal([]byte(attr["params"]), &sqlParam)

	sqlTemplate := models.SqlTemplate{}
	sql := sqlTemplate.GenerateSql(attr["sql"], sqlParam, "dag")

	parser := Sql.SemanticParser{}
	var fields []string
	if attr["taskType"] == "hive" {
		metaURL := utils.SparkRoute(models.ServerInfo.SparkHost, models.ServerInfo.SparkPort, "dag_getFieldName")
		param, _ := json.Marshal(map[string]string{"sql": sql, "taskType": attr["taskType"]})
		req := httplib.Post(metaURL).SetTimeout(5*time.Minute, 5*time.Minute).Body(param)
		bytes, err := req.Bytes()
		if err != nil {
			this.Ctx.WriteString(utils.VueResponse{
				Res:  false,
				Info: err.Error(),
			}.String())
		} else {
			sparkserverRes := utils.SparkServerResponse{}
			_ = json.Unmarshal(bytes, &sparkserverRes)
			this.Ctx.WriteString(utils.VueResponse{
				Res:  sparkserverRes.Flag,
				Info: sparkserverRes.Msg,
			}.String())
		}
	} else {
		fields = parser.GetFields(sql)
		res, _ := json.Marshal(fields)
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: string(res),
		}.String())
	}
}

// @Title 获取hive存储引擎
// @Summary 获取hive存储引擎
// @Description 获取hive存储引擎
// @Success 200
// @Accept json
// @router /getHiveStorage [get]
func (this *DagController) GetHiveStorage() {
	types := []hive.StorageType{
		hive.OrcSnappy, hive.TSV,
	}
	bytes, _ := json.Marshal(types)
	this.Ctx.WriteString(string(bytes))

}

// @Title dag输出表从旧表选择时, 获取输出表的选项
// @Summary dag输出表从旧表选择时, 获取输出表的选项
// @Description dag输出表从旧表选择时, 获取输出表的选项
// @Success 200
// @Accept json
// @router /getOutTables [post]
func (this *DagController) GetOutTables() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	role := models.Enforcer.GetRolesForUser(this.userName)[0]
	dsId, _ := strconv.ParseInt(attr["dsId"], 10, 64)
	tables := models.GetOutTables(role, attr["tableType"], dsId)
	res, _ := json.Marshal(tables)
	this.Ctx.WriteString(string(res))
}

// @Title dag输出表从旧表选择时, 校验输出表和sql产生的schema是否一致
// @Summary dag输出表从旧表选择时, 校验输出表和sql产生的schema是否一致
// @Description dag输出表从旧表选择时, 校验输出表和sql产生的schema是否一致
// @Success 200
// @Accept json
// @router /checkOutTable [post]
func (this *DagController) CheckOutTable() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	outTable := attr["outTableName"]
	outType := attr["outTableType"]
	taskType := attr["taskType"]
	var sqlParam []map[string]string
	_ = json.Unmarshal([]byte(attr["params"]), &sqlParam)

	defer func() {
		if r := recover(); r != nil {
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()

	sqlTemplate := models.SqlTemplate{}
	sql := sqlTemplate.GenerateSql(attr["sql"], sqlParam, "dag")
	sql = strings.Replace(sql, ";", "", -1)

	flag, msg := models.CheckMeta(sql, outTable, outType, taskType)
	info := ""
	if !flag {
		info = fmt.Sprintf("sql与输出表的schema不匹配, 或未配置sql参数. %s", msg)
	}
	this.Ctx.WriteString(utils.VueResponse{
		Res:  flag,
		Info: info,
	}.String())
}

// @Title 推完dag文件后, 自动导入表之间的关系
// @Summary 推完dag文件后, 自动导入表之间的关系
// @Description 推完dag文件后, 自动导入表之间的关系
// @Success 200
// @Accept json
// @router /afterPostDag [post]
func (this *DagController) AfterPostDag() {
	var view models.DagView
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &view)
	rolesStr := strings.Join(models.Enforcer.GetRolesForUser(this.userName),",")
	defer func() {
		if r := recover(); r != nil {
			models.Logger.Error(fmt.Sprintf("%v", r))
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()
	view.AfterPostDag(rolesStr)
	this.Ctx.WriteString(utils.VueResponse{
		Res:  true,
		Info: "表关系添加成功",
	}.String())
}


// @Title 手动触发单日期跑
// @Summary 手动触发单日期跑
// @Description 手动触发单日期跑
// @Success 200
// @Accept json
// @router /manualTrigger [post]
func (this *DagController) ManualTrigger() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	startTime := attr["startTime"]
	endTime := attr["endTime"]
	dagName := attr["dagName"]

	var dag models.DagAirFlow
	dag.GetByName(dagName)

	var dagView models.DagView
	_ = json.Unmarshal([]byte(dag.Description), &dagView)

	dagView.PostDagToQueue(dagView, startTime, endTime, this.userName) // 后台开启
}

// @Title 获取手动重跑的记录
// @Summary 获取手动重跑的记录
// @Description 获取手动重跑的记录
// @Success 200
// @Accept json
// @router /getManualTriggerHistory [post]
func (this *DagController) GetManualTriggerHistory() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	dagName := attr["dagName"]

	histories := models.GetAllTriggerHistory(dagName)

	var bodys []map[string]interface{}
	for _, data := range histories {
		// 转成map
		tmpData := utils.GetSnakeStrMapOfStruct(data)
		bodys = append(bodys, tmpData)
	}

	// 返回数据
	tsb := utils.TableData{Data: bodys}
	tsb.GetSnakeTitles(models.DagTriggerHistory{})
	res, _ := json.Marshal(tsb)

	this.Ctx.WriteString(string(res))
}

// @Title 获取模块日志
// @Summary 获取模块日志
// @Description 获取模块日志
// @Success 200
// @Accept json
// @router /getTriggerLog [post]
func (this *DagController) GetTriggerLog() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	remoteId := attr["remoteId"]
	userName := attr["userName"]
	log := models.GetTaskLog(remoteId, userName)
	bytes,_ := json.Marshal(log)
	this.Ctx.WriteString(string(bytes))
}