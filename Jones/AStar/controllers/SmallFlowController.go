package controllers

import (
	"AStar/models"
	pb "AStar/protobuf"
	"AStar/utils"
	"AStar/utils/SmallFlow"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type SmallFlowController struct {
	BaseController
}

// @Title 提交小流量任务pc
// @Summary 提交小流量任务pc
// @Param para body models.FlowTaskPcDef true "任务信息"
// @Description 提交小流量任务pc
// @Success 200
// @Failure 400
// @Failure 403
// @Failure 404
// @Accept json
// @router /small/submit2pc [post]
func (this *SmallFlowController) Submit2Pc() {

	var queryCollections = []string{}

	var taskID = utils.DateFormatWithoutLine(time.Now()) + "_" + utils.GetMd5ForQuery("small", this.userName, time.Now().String())[9:30]

	var tmpFullTbName = taskID + "_FULL_{dt}"
	//var tmpSameTbName = queryID + "_SAME"

	var taskJSON = models.FlowTaskPcDef{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &taskJSON)
	if err != nil {
		models.Logger.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		return
	}

	var startDate = taskJSON.DateRange[0]
	var endDate = taskJSON.DateRange[1]

	var taskType = taskJSON.Type
	var session = this.userName
	var queryName = taskJSON.Name

	var ipRange = taskJSON.IpRestrict
	var locRange = taskJSON.AdPosRestrict
	var ipRst = SmallFlow.GetIPRestriction(strings.Join(ipRange, ","))

	// 试验路对照路判断+实验样式判断的where语句
	var abTest = SmallFlow.GetPcABTestCase(taskJSON)
	var abTestWhere = SmallFlow.GetPcABTestWhere(taskJSON)
	var locRst = SmallFlow.GetLocRestriction(strings.Join(locRange, ","))

	// 全量结果SQL语句
	var fullDataSQL = SmallFlow.GetFullDataSQL(taskID, taskType, abTest, abTestWhere, ipRst, locRst)
	queryCollections = append(queryCollections, fullDataSQL)

	// 同质结果SQL语句
	var sameDataSQL = SmallFlow.GetSameQuerySQL(tmpFullTbName)
	queryCollections = append(queryCollections, sameDataSQL)

	var smallFlow = pb.SmallFlow{
		TaskId:    taskID,
		TaskName:  queryName,
		TaskType:  taskType,
		StartDate: startDate,
		EndDate:   endDate,
		Session:   session,
		Sqls:      queryCollections,
		Json:      string(this.Ctx.Input.RequestBody),
	}

	models.SmallFlowTask{}.InsertSmallFlowTask(smallFlow)

	this.Ctx.ResponseWriter.Write([]byte(smallFlow.TaskId))

}

// @Title 提交小流量任务
// @Summary 提交小流量任务
// @Param para body models.FlowTaskDef true "任务信息"
// @Description 提交小流量任务
// @Success 200
// @Failure 400
// @Failure 403
// @Failure 404
// @Accept json
// @router /small/submit2ws [post]
func (this *SmallFlowController) Submit2Ws() {

	var queryCollections = []string{}

	var taskID = utils.DateFormatWithoutLine(time.Now()) + "_" + utils.GetMd5ForQuery("small", this.userName, time.Now().String())[9:30]

	var tmpFullTbName = taskID + "_FULL_{dt}"
	//var tmpSameTbName = queryID + "_SAME"

	var taskJSON = models.FlowTaskDef{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &taskJSON)
	if err != nil {
		models.Logger.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		return
	}

	var startDate = taskJSON.DateRange[0]
	var endDate = taskJSON.DateRange[1]

	var taskType = taskJSON.Type
	var session = this.userName
	var queryName = taskJSON.Name

	var ipRange = taskJSON.IpRestrict
	var locRange = taskJSON.AdPosRestrict
	var ipRst = SmallFlow.GetIPRestriction(strings.Join(ipRange, ","))

	//var abTest = SmallFlow.GetABTest(strconv.Itoa(exprFlowPos), exprFlowVal, exprFlowType, strconv.Itoa(ctrlFlowPos), ctrlFlowVal, ctrlFlowType)

	// 试验路对照路判断+实验样式判断的where语句
	var abTest = SmallFlow.GetABTestCase(taskJSON)

	var abTestWhere = SmallFlow.GetABTestWhere(taskJSON)
	var locRst = SmallFlow.GetLocRestriction(strings.Join(locRange, ","))

	// 全量结果SQL语句
	var fullDataSQL = SmallFlow.GetFullDataSQL(taskID, taskType, abTest, abTestWhere, ipRst, locRst)
	queryCollections = append(queryCollections, fullDataSQL)

	// 同质结果SQL语句
	var sameDataSQL = SmallFlow.GetSameQuerySQL(tmpFullTbName)
	queryCollections = append(queryCollections, sameDataSQL)

	var smallFlow = pb.SmallFlow{
		TaskId:    taskID,
		TaskName:  queryName,
		TaskType:  taskType,
		StartDate: startDate,
		EndDate:   endDate,
		Session:   session,
		Sqls:      queryCollections,
		Json:      string(this.Ctx.Input.RequestBody),
	}

	//fmt.Println(smallFlow.Sqls)

	models.SmallFlowTask{}.InsertSmallFlowTask(smallFlow)

	this.Ctx.ResponseWriter.Write([]byte(smallFlow.TaskId))

}

// @Title 查询小流量任务
// @Summary 查询小流量任务
// @Param para body models.SmallFlowQueryJson true "任务信息"
// @Description 查询小流量任务
// @Success 200  {object} models.TableTitle
// @Failure 400
// @Failure 403
// @Failure 404
// @Accept json
// @router /small/query [post]
func (this *SmallFlowController) Query() {

	// 解析前端json保存在smallFlowQuery变量中
	var sfQuery models.SmallFlowQueryJson
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &sfQuery)

	// 解析异常报错
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte("para error"))
		return
	}

	// 拼接sql语句
	var sql = SmallFlow.GetSmallFlowResSQL(sfQuery.TaskId, sfQuery.DateRange, sfQuery.IsMobilQq, sfQuery.IsSame, !sfQuery.IsAvg, sfQuery.CustomId, sfQuery.Query)

	// 获取要查看的日期数组和要查看的排位数组
	// 其中：
	// 日期的数组从当前方法的dfQuery JSON中取出
	// 排位的数组需要从数据库记录的提交任务的JSON中取出

	// 取出日期数组
	var dateArr = utils.GetDateRange(sfQuery.DateRange[0], sfQuery.DateRange[1])
	// 当选择了isAvt，也就是求日均结果的时候，日期数组应该仅包含一个空字符串(注意不是一个空数组)
	if !sfQuery.IsAvg {
		dateArr = []string{""}
	}
	// 取出日期数组 --完成


	// 取出排位数组，需要从数据库当中取，相对应的ORM结构体是SmallFlowTask
	var tasks []models.SmallFlowTask
	orm.NewOrm().QueryTable("small_flow_task").Filter("task_id", sfQuery.TaskId).All(&tasks)
	var taskJSON = models.FlowTaskDef{}
	err2 := json.Unmarshal([]byte(tasks[0].Json), &taskJSON)
	if err2 != nil {
		models.Logger.Error(err2.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	var locArr = taskJSON.AdPosRestrict
	locArr = append(locArr, "全部")
	// 取出排位数组 -- 完成

	// 取出用户数组
	var cusArr []string
	if len(sfQuery.CustomId) != 0 {
		cusArr = sfQuery.CustomId
	} else {
		cusArr = append(cusArr, "")
	}

	// showArgs数组用于辅助GetTitles方法来判断哪些title需要被设置为hidden
	// 需要结合SmallFlowRe中的自定义字段来判断
	var showArgs = make(map[string]bool)
	showArgs["dt"] = sfQuery.IsAvg
	var hasCustomerId = false
	if len(sfQuery.CustomId) != 0 {
		hasCustomerId = true
	}
	showArgs["custom_id"] = hasCustomerId
	// 至此showArgs填充完毕

	// 用于存储sql的查询结果
	var queryRes []models.SmallFlowRes

	//执行查询
	orm.NewOrm().Raw(sql).QueryRows(&queryRes)

	//fmt.Println(sql)

	// 定义填充数据和计算A/B-1结果后的数据结构
	var filledQueryResMapExpr = make(map[string]models.SmallFlowRes)
	var filledQueryResMapCtrl = make(map[string]models.SmallFlowRes)
	// 创建空结果用于填充不存在的pw和日期等
	var emptyRes = models.SmallFlowRes{
		Id:         0,
		Dt:         "-",
		Tag:        "-",
		Pw:         "-",
		SumPvPage:  "-",
		SumPvAd:    "-",
		SumCharge:  "-",
		SumClick:   "-",
		AvgCtrPage: "-",
		AvgCtrAd:   "-",
		AvgAcp:     "-",
		AvgRpmPage: "-",
		CustomId:   "-",
	}

	// 对取出的数据库中的数据queryRes进行操作
	// 生成 dt_pw 即形如 20180810_2 的key，并生成空的结果集
	for _, dt := range dateArr {
		for _, loc := range locArr {
			for _, cus := range cusArr {
				var key = dt + "_" + loc + "_" + cus
				emptyRes.Dt = dt
				emptyRes.Pw = loc
				emptyRes.CustomId = cus
				filledQueryResMapExpr[key] = emptyRes
				filledQueryResMapCtrl[key] = emptyRes
			}
		}
	}

	// 将实际有的值填充进结果集

	for _, val := range queryRes {
		var key = val.Dt + "_" + val.Pw + "_" + val.CustomId
		if val.Tag == "实验组" {
			filledQueryResMapExpr[key] = val
		} else { // if val.Tag == "对照组"
			filledQueryResMapCtrl[key] = val
		}

	}

	var fieldIdx = 0
	switch sfQuery.Query {
	case "整体页面":
		fieldIdx = 0
	case "实验样式-样式":
		fieldIdx = 1
	case "实验样式-页面":
		fieldIdx = 2
	}

	// 生成resBody用于返回给前端，是表格中的数据
	var resBody []map[string]interface{}

	// 将struct结果按循环顺序转换成map，添加进用于返回给前端的resBody当中
	for _, dt := range dateArr {
		for _, loc := range locArr {
			for _, cus := range cusArr {
				resBody = append(resBody, utils.GetSnakeStrMapOfStruct(filledQueryResMapExpr[dt+"_"+loc+"_"+cus]))
			}
		}
		for _, loc := range locArr {
			for _, cus := range cusArr {
				resBody = append(resBody, utils.GetSnakeStrMapOfStruct(filledQueryResMapCtrl[dt+"_"+loc+"_"+cus]))
			}
		}

		for _, loc := range locArr {
			for _, cus := range cusArr {
				var resExpr = filledQueryResMapExpr[dt+"_"+loc+"_"+cus]
				var resCtrl = filledQueryResMapCtrl[dt+"_"+loc+"_"+cus]
				var calcedRes = SmallFlow.CalcABTestResult(resExpr, resCtrl)
				var resBodyTable = utils.GetSnakeStrMapOfStruct(calcedRes)
				resBody = append(resBody, resBodyTable)
			}
		}
	}
	//for _, dt := range dateArr {
	//	for _, loc := range locArr {
	//		for _, cus := range cusArr {
	//			resBody = append(resBody, utils.GetSnakeStrMapOfStruct(filledQueryResMapCtrl[dt+"_"+loc+"_"+cus]))
	//		}
	//	}
	//}

	// 在上述结果的基础上还需要继续添加A/B-1的值
	//for _, dt := range dateArr {
	//	for _, loc := range locArr {
	//		for _, cus := range cusArr {
	//			var resExpr = filledQueryResMapExpr[dt+"_"+loc+"_"+cus]
	//			var resCtrl = filledQueryResMapCtrl[dt+"_"+loc+"_"+cus]
	//			var calcedRes = SmallFlow.CalcABTestResult(resExpr, resCtrl)
	//			var resBodyTable = utils.GetSnakeStrMapOfStruct(calcedRes)
	//			resBody = append(resBody, resBodyTable)
	//		}
	//	}
	//}

	// 生成json并返回
	tsb := utils.TableData{Data: resBody}
	tsb.GetSnakeTitlesWithHiddenForMuti(models.SmallFlowRes{}, showArgs, fieldIdx)
	resJSON, _ := json.Marshal(tsb)
	this.Ctx.WriteString(string(resJSON))
}

// @Title 下载小流量任务
// @Summary 下载小流量任务
// @Param para body models.SmallFlowQueryJson true "任务信息"
// @Description 下载小流量任务
// @Success 200  {object} models.TableTitle
// @Failure 400
// @Failure 403
// @Failure 404
// @Accept json
// @router /small/download [post]
func (this *SmallFlowController) Download() {

	// 解析前端json保存在smallFlowQuery变量中
	var sfQuery models.SmallFlowQueryJson
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &sfQuery)

	// 解析异常报错
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte("para error"))
		return
	}

	// 拼接sql语句
	sql, _ := SmallFlow.GetSmallFlowDownSQL(sfQuery.TaskId, sfQuery.DateRange, sfQuery.IsMobilQq, sfQuery.IsSame, !sfQuery.IsAvg, sfQuery.CustomId, sfQuery.Download)
	//fmt.Println(sql)

	// 获取table
	var tables []string
	tables = append(tables, "adtl_john.small_flow_" + sfQuery.TaskId)

	// 指定查询引擎
	var queryEngine = pb.QueryEngine_hive

	// 这些值是用于给QueryInfo存储的 wait to fill
	var dimListForQueryInfo []string
	var outputListForQueryInfo []string
	var filterListForQueryInfo []string

	// 同质有查询词维度
	if sfQuery.IsSame {
		dimListForQueryInfo = []string{ "日期", "客户ID", "一级行业", "二级行业", "三级行业", "查询词",
		"实验组", "ADPV",  "ADCLICK", "ADCOST", "ADCTR",  "ADACP", "ADRPM",
		"对照组", "ADPV",  "ADCLICK", "ADCOST", "ADCTR",  "ADACP", "ADRPM",}
	} else {
		dimListForQueryInfo = []string{ "日期", "客户ID", "一级行业", "二级行业", "三级行业",
			"实验组", "ADPV",  "ADCLICK", "ADCOST", "ADCTR",  "ADACP", "ADRPM",
			"对照组", "ADPV",  "ADCLICK", "ADCOST", "ADCTR",  "ADACP", "ADRPM",}
	}


	// 不同于query 放入MQ,返回queryId
	nowTime := time.Now()
	// 生成Md5结果，作为QueryInfo的QueryId
	Md5Res := utils.GetMd5ForQuery(sql, this.userName, nowTime.String())
	// 生成QueryInfo
	queryInfo := pb.QueryInfo{
		QueryId:    Md5Res,                 // query id
		Sql:        []string{sql},          // sql statement
		Session:    this.userName,          // session
		Engine:     queryEngine,            // 引擎
		Dimensions: dimListForQueryInfo,    // 维度表头
		Outputs:    outputListForQueryInfo, // 输出表头
		Filters:    filterListForQueryInfo, // 过滤字段
		Status:     pb.QueryStatus_Waiting, // 状态
		LastUpdate: nowTime.UnixNano() / 1e6,// 最后更新时间戳
		CreateTime: nowTime.UnixNano() / 1e6,// 创建时间戳
		Tables:     tables,                 // 查询所用表
		Type:       "小流量下载",                   // 查询类别"小流量下载"
	}
	var queryTaskName = sfQuery.TaskId + sfQuery.Download
	var queryJSONStr = string(this.Ctx.Input.RequestBody)

	// 将QueryInfo入库MySql
	models.QuerySubmitLog{}.InsertQuerySubmitLog(queryInfo, queryTaskName, queryJSONStr)
	// 编码QueryInfo
	queryBytes, _ := proto.Marshal(&queryInfo)
	// 得到MQ
	models.RabbitMQModel{}.Publish(queryBytes)

	// 返回queryId
	this.Ctx.WriteString(string(queryInfo.QueryId))
}

func CreateStruct(fields []reflect.StructField) reflect.Value {
	var structType reflect.Type
	structType = reflect.StructOf(fields)
	return reflect.Zero(structType)
}

// @Title 查看可用的小流量任务
// @Summary 查看可用的小流量任务
// @Description 查看可用的小流量任务
// @Success 200
// @Failure 400
// @Failure 403
// @Failure 404
// @Accept json
// @router /small/log/enable [get]
func (this *SmallFlowController) GetAllEnable() {
	var sfTasks []models.SmallFlowTask
	orm.NewOrm().QueryTable("small_flow_task").Filter("is_enabled", 1).All(&sfTasks)
	tsb := GetSmallFlowTasks(sfTasks)
	data, _ := json.Marshal(tsb)
	this.Ctx.WriteString(string(data))
}

// @Title 查看下线的小流量任务
// @Summary 查看下线的小流量任务
// @Description 查看下线的小流量任务
// @Success 200
// @Failure 400
// @Failure 403
// @Failure 404
// @Accept json
// @router /small/log/disable [get]
func (this *SmallFlowController) GetAllDisable() {
	var sfTasks []models.SmallFlowTask
	orm.NewOrm().QueryTable("small_flow_task").Filter("is_enabled", 2).All(&sfTasks)
	tsb := GetSmallFlowTasks(sfTasks)
	data, _ := json.Marshal(tsb)
	this.Ctx.WriteString(string(data))
}


// @Title 查看待上线的小流量任务
// @Summary 查看待上线的小流量任务
// @Description 查看待上线的小流量任务
// @Success 200
// @Failure 400
// @Failure 403
// @Failure 404
// @Accept json
// @router /small/log/ready [get]
func (this *SmallFlowController) GetAllReady() {
	var sfTasks []models.SmallFlowTask
	orm.NewOrm().QueryTable("small_flow_task").Filter("is_enabled", 0).All(&sfTasks)
	tsb := GetSmallFlowTasks(sfTasks)
	data, _ := json.Marshal(tsb)
	this.Ctx.WriteString(string(data))
}


// @Title 上线的小流量任务
// @Param	task_id path string true "task_id"
// @Summary 上线的小流量任务
// @Description 上线的小流量任务
// @Success 200
// @Failure 400
// @Failure 403
// @Failure 404
// @Accept json
// @router /small/task/online/:task_id [get]
func (this *SmallFlowController) TaskOnline() {
	taskId:=this.Ctx.Input.Params()[":task_id"]
	var sfTasks []models.SmallFlowTask
	orm.NewOrm().QueryTable("small_flow_task").Filter("task_id", taskId).All(&sfTasks)
	if len(sfTasks) <= 0 {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte("未找到正确的task_id"))
		return
	}
	sfTasks[0].IsEnabled = 1
	err := sfTasks[0].UpdateEnable()
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte("上线失败请联系管理员"))
		return
	}
	this.Ctx.ResponseWriter.WriteHeader(200)
	this.Ctx.ResponseWriter.Write([]byte("ok"))
}


// @Title 下线的小流量任务
// @Param	task_id path string true "task_id"
// @Summary 下线的小流量任务
// @Description 下线的小流量任务
// @Success 200
// @Failure 400
// @Failure 403
// @Failure 404
// @Accept json
// @router /small/task/offline/:task_id [get]
func (this *SmallFlowController) TaskOffline() {
	taskId:=this.Ctx.Input.Params()[":task_id"]

	var sfTasks []models.SmallFlowTask
	orm.NewOrm().QueryTable("small_flow_task").Filter("task_id", taskId).All(&sfTasks)
	if len(sfTasks) <= 0 {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte("未找到正确的task_id"))
		return
	}
	sfTasks[0].IsEnabled = 2
	err := sfTasks[0].UpdateEnable()
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte("上线失败请联系管理员"))
		return
	}
	this.Ctx.ResponseWriter.WriteHeader(200)
	this.Ctx.ResponseWriter.Write([]byte("ok"))
}


func GetSmallFlowTasks(sfTasks [] models.SmallFlowTask) utils.TableData {
	var resBody []map[string]interface{}
		for _, task := range sfTasks {
		var resBodyTable = utils.GetSnakeStrMapOfStruct(task)
		resBody = append(resBody, resBodyTable)
	}

	// 返回数据
	tsb := utils.TableData{Data: resBody}
	tsb.GetSnakeTitles(models.SmallFlowTask{})
	return tsb
}

// @Title 根据task_id取得对应的queryJSON
// @Param	task_id path string true "task_id"
// @Summary 根据task_id取得对应的queryJSON
// @Description 根据task_id取得对应的queryJSON
// @Success 200 {string}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /small/log/query/:task_id [get]
func (this *SmallFlowController) GetJsonByTaskId() {
	taskId:=this.Ctx.Input.Params()[":task_id"]
	queryJson, err := models.SmallFlowTask{}.GetQueryJson(taskId)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte(err.Error()))
	} else {
		this.Ctx.WriteString(queryJson)
	}
}



// @Title GetIpPc
// @Description 根据类型获取IP列表
// @Success 200 ok
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /small/ip/pc [get]
func (this *SmallFlowController) GetIpPc() {
	resp, err := http.Get("http://stat.union.sogou-inc.com/server/GetIpByPath.php?path=PC_Search/main/BiddingServer&type=json")
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
	}
	defer resp.Body.Close()
	cnt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
	}
	this.Ctx.WriteString(string(cnt))
}


// @Title GetIpWap
// @Description 根据类型获取IP列表
// @Success 200 ok
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /small/ip/wap [get]
func (this *SmallFlowController) GetIpWap() {
	resp, err := http.Get("http://stat.union.sogou-inc.com/server/GetIpByPath.php?path=Mobile_Search/main/BiddingServer&type=json")
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
	}
	defer resp.Body.Close()
	cnt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
	}
	this.Ctx.WriteString(string(cnt))
}

