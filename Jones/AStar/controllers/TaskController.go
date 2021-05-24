package controllers

import (
	"AStar/models"
	pb "AStar/protobuf"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"github.com/saintfish/chardet"
	"strconv"
	"strings"
	"time"
)

type TaskController struct {
	BaseController
}

// @Title 查看所有SubmitLog
// @Summary 查看所有SubmitLog
// @Description 查看所有SubmitLog
// @Success 200 {object} models.TableTitle
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/log [get]
func (this *TaskController) GetAll() {
	var QueryLogs []models.QuerySubmitLog
	orm.NewOrm().QueryTable("query_submit_log").OrderBy("-create_time").All(&QueryLogs)
	tsb := GetSubmitLog(QueryLogs)
	data, _ := json.Marshal(tsb)
	this.Ctx.WriteString(string(data))
}

// @Title 按用户查看所有SubmitLog
// @Param	username path string true "用户名"
// @Summary 按用户查看所有SubmitLog
// @Description 按用户查看所有SubmitLog
// @Success 200 {object} models.TableTitle
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/log/user/:username [get]
func (this *TaskController) GetByUser() {
	username := this.Ctx.Input.Params()[":username"]
	var QueryLogs []models.QuerySubmitLog
	orm.NewOrm().QueryTable("query_submit_log").Filter("session", username).OrderBy("-create_time").All(&QueryLogs)
	tsb := GetSubmitLog(QueryLogs)
	data, _ := json.Marshal(tsb)
	this.Ctx.WriteString(string(data))
}

// @Title 按日期查看所有SubmitLog
// @Param	date path string true "日期(yyyy-MM-dd)"
// @Summary 按日期查看所有SubmitLog
// @Description 按日期查看所有SubmitLog
// @Success 200 {object} models.TableTitle
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/log/date/:date [get]
func (this *TaskController) GetByDate() {
	date := this.Ctx.Input.Params()[":date"]
	var QueryLogs []models.QuerySubmitLog
	orm.NewOrm().QueryTable("query_submit_log").Filter("date", date).OrderBy("-create_time").All(&QueryLogs)
	tsb := GetSubmitLog(QueryLogs)
	data, _ := json.Marshal(tsb)
	this.Ctx.WriteString(string(data))
}

// @Title 获取Rsync地址的根目录
// @Summary 获取Rsync地址的根目录，用于拼接成完整的rsync地址
// @Description 获取Rsync地址的根目录，用于拼接成完整的rsync地址
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/log/rsync [get]
func (this *TaskController) GetRsync() {
	rsyncPath := []string{models.ServerInfo.RsyncPath, models.ServerInfo.PaganiRsyncPath}
	jsonData, _ := json.Marshal(rsyncPath)
	this.Ctx.WriteString(string(jsonData))
}

// @Title 按配置条件查看所有SubmitLog
// @Param	session query string false "用户名"
// @Param	date query string false "日期(yyyy-MM-dd)"
// @Param	status query int false "任务状态"
// @Summary 按配置条件查看所有SubmitLog
// @Description 按配置条件查看所有SubmitLog
// @Success 200 {object} models.TableTitle
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/log/config [get]
func (this *TaskController) GetByConfig() {
	startDate := this.GetString("startDate")
	date := this.GetString("date")
	session := this.GetString("session")
	status := this.GetString("status")

	querySeter := orm.NewOrm().QueryTable("query_submit_log")
	var queryLogs []models.QuerySubmitLog
	querySeter = querySeter.Filter("session", session).
		Filter("date__gte", startDate).
		Filter("date__lte", date)
	if status != "" {
		querySeter = querySeter.Filter("status", status)
	}
	_, _ = querySeter.OrderBy("-create_time").All(&queryLogs)
	tsb := GetSubmitLog(queryLogs)
	data, _ := json.Marshal(tsb)
	this.Ctx.WriteString(string(data))
}

// @Title 查询所有用户的SubmitLog
// @Summary 按配置条件查看所有SubmitLog
// @Description 按配置条件查看所有SubmitLog
// @Success 200 {object} models.TableTitle
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task_god/log/config [get]
func (this *TaskController) GetByConfigGod() {
	configList := []string{"date", "session", "status"}
	querySeter := orm.NewOrm().QueryTable("query_submit_log")
	var QueryLogs []models.QuerySubmitLog
	for _, config := range configList {
		if this.GetString(config) != "" {
			querySeter = querySeter.Filter(config, this.GetString(config))
		}
	}
	querySeter.OrderBy("-create_time").All(&QueryLogs)
	tsb := GetSubmitLog(QueryLogs)
	data, _ := json.Marshal(tsb)
	this.Ctx.WriteString(string(data))
}

// @Title 根据ID查看发任务详情
// @Summary 根据ID查看发任务详情
// @Description 根据ID查看发任务详情
// @Success 200 {object} string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task_god/query/:queryId [get]
func (this *TaskController) GetTaskById() {
	queryId := this.Ctx.Input.Params()[":queryId"]
	o := orm.NewOrm()
	submitLog := []models.QuerySubmitLog{}
	_, _ = o.Raw("SELECT * FROM query_submit_log WHERE query_id=?", queryId).QueryRows(&submitLog)
	data := GetSubmitLog(submitLog).Data[0]
	res, _ := json.Marshal(data)
	this.Ctx.WriteString(string(res))
}

func GetSubmitLog(QueryLogs []models.QuerySubmitLog) utils.TableData {

	// 创建表内容
	var bodys []map[string]interface{}

	// 对保存数据库中值的QueryLogs进行遍历
	for _, queryLog := range QueryLogs {

		// 将queryLog映射成map，保存成tmpData
		tmpData := utils.GetSnakeStrMapOfStruct(queryLog)

		subTime := fmt.Sprintf("%.2f", float64(tmpData["last_update"].(int64)-tmpData["create_time"].(int64))/1000)
		tmpData["sub_time"] = subTime

		for key, value := range tmpData {

			// 将最后更新时间更改为可读时间
			if key == "last_update" || key == "create_time" {
				timeUnix := value.(int64)
				t := utils.TimeFormat(time.Unix(timeUnix/1000, 0))
				tmpData[key] = t
			}
		}

		statusCode := tmpData["status"]
		switch tmpData["status"] {
		case pb.QueryStatus_Waiting:
			tmpData["status"] = "任务提交中"
		case pb.QueryStatus_Running:
			tmpData["status"] = "任务执行中"
		case pb.QueryStatus_Killing:
			tmpData["status"] = "正在取消"
		case pb.QueryStatus_Finished:
			tmpData["status"] = "任务完成"
		case pb.QueryStatus_Killed:
			tmpData["status"] = "任务取消"
		case pb.QueryStatus_SparkReRunning:
			tmpData["status"] = "任务执行中"
		case pb.QueryStatus_SparkReRunFinished:
			tmpData["status"] = "任务完成"
		case pb.QueryStatus_SyncQueryRunning:
			tmpData["status"] = "同步任务执行中"
		case pb.QueryStatus_SyncQueryFinished:
			tmpData["status"] = "同步任务完成"
		case pb.QueryStatus_SyncQueryFailed:
			tmpData["status"] = "同步任务失败"
		case pb.QueryStatus_CustomQueryWaiting:
			tmpData["status"] = "自定义任务提交中"
		case pb.QueryStatus_CustomQueryRunning:
			tmpData["status"] = "自定义任务执行中"
		case pb.QueryStatus_CustomQueryFinished:
			tmpData["status"] = "自定义任务完成"
		case pb.QueryStatus_CustomQueryFailed:
			tmpData["status"] = "自定义任务失败"
		}

		tmpData["status_code"] = statusCode

		bodys = append(bodys, tmpData)
	}

	// 返回数据

	tsb := utils.TableData{Data: bodys}
	tsb.GetSnakeTitles(models.QuerySubmitLog{})
	tsb.Cols = append(tsb.Cols, utils.TableTitle{Prop: "sub_time", Label: "执行时间（s）", Hidden: true})
	return tsb
}

// @Title 获取执行的数据
// @Param	data_path formData string true "data_path"
// @Param	status formData int true "status"
// @Summary 获取执行的数据
// @Description 获取执行的数据
// @Success 200 {object} models.TableTitle
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/data/path [post]
func (this *TaskController) Limit() {

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

	// 约定的key
	key := "data_path"

	// 取出key对应的data_path值
	dataPath := this.GetString(key)
	status, _ := this.GetInt("status")

/*	start := time.Now().Unix()
	//查看是否又缓存数据，有则使用缓存
	//查找缓存出现错误则继续请求SparkServer
	fmt.Println("访问redis..")
	taskid := strings.Split(dataPath, "/")[1]
	existsId, existErr := existsTaskIdCache(taskid)
	if existErr == nil && existsId == true {
		stringJSON, err := getTaskIdCache(taskid)
		if err == nil {
			this.Ctx.WriteString(stringJSON)

			end := time.Now().Unix()
			fmt.Printf("redis handle cost is :%d \n", (end - start))
			return
		}
	}*/

	// 组装成queryInfo并封装成二进制流
	infoBytes, _ := proto.Marshal(&pb.QueryInfo{DataPath: dataPath})

	var reqURL string
	// 请求spark-server
	if status == 22 {
		reqURL = utils.SparkRoute(models.ServerInfo.PaganiHost, models.ServerInfo.SparkPort, models.RouteMap[this.controllerName+"_"+this.actionName])
	} else {
		reqURL = utils.SparkRoute(models.ServerInfo.SparkHost, models.ServerInfo.SparkPort, models.RouteMap[this.controllerName+"_"+this.actionName])
	}
	//fmt.Println("redis未找到, 访问get_limit")
	pbRequest := httplib.Post(reqURL).SetTimeout(time.Hour, time.Hour).Body(infoBytes)

	// 获取response流
	bytesReceived, err := pbRequest.Bytes()
	if err != nil {
		panic(err.Error())
	}

	// 取回的TranStream的protobuffer
	pbResponse := &pb.TranStream{}

	proto.Unmarshal(bytesReceived, pbResponse)

	fmt.Println("get_limit访问完毕")
	// 生成表格
	var resTitle []utils.TableTitle
	var resBody []map[string]interface{}

	resTitle = append(resTitle, utils.TableTitle{Prop: "idx", Label: "行号", Sortable: true, Hidden: false})

	for _, schema := range pbResponse.Schema {
		commentLabel := getCommentOfHeader(schema)
		resTitle = append(resTitle, utils.TableTitle{Prop: schema, Label: commentLabel, Sortable: false, Hidden: false})
	}

	for idx, line := range pbResponse.Line {
		var tmap = make(map[string]interface{})
		tmap["idx"] = idx + 1
		for idxline, field := range line.Field {
			tmap[pbResponse.Schema[idxline]] = field
		}
		resBody = append(resBody, tmap)
	}

	// 生成TableData的struct
	tsb := utils.TableData{Cols: resTitle, Data: resBody}

	// 生成json返回给前端
	resJSON, _ := json.Marshal(tsb)

	this.Ctx.WriteString(utils.VueResponse{
		Res:  true,
		Info: string(resJSON),
	}.String())
}

func existsTaskIdCache(taskid string) (bool, error) {
	c, err := redis.Dial("tcp", models.RedisConfig.RedisHost+":"+models.RedisConfig.RedisPort, redis.DialPassword(models.RedisConfig.RedisPwd))
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return false, err
	}
	defer c.Close()

	is_key_exit, err := redis.Bool(c.Do("EXISTS", "_spark:"+taskid+":schema"))
	return is_key_exit, err
}

func getTaskIdCache(taskid string) (string, error) {
	//要产出的数据
	resJSON := ""
	var resTitle []utils.TableTitle
	var resBody []map[string]interface{}
	resTitle = append(resTitle, utils.TableTitle{Prop: "idx", Label: "行号", Sortable: true, Hidden: false})

	//建立Redis链接
	conn, err := redis.Dial("tcp", models.RedisConfig.RedisHost+":"+models.RedisConfig.RedisPort, redis.DialPassword(models.RedisConfig.RedisPwd))
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return resJSON, err
	}
	defer conn.Close()

	//获取该taskid 的所有key
	iter := 0
	keys := []string{}
	pattern := taskid + "*"

	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern, "COUNT", "5000"))
		if err != nil {
			fmt.Errorf("error retrieving '%s' keys", pattern)
			continue
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	//如果key 为空， 说明结果集为空
	if len(keys) == 0 {
		tsb := utils.TableData{Cols: resTitle, Data: resBody}
		byteJSON, _ := json.Marshal(tsb)
		return string(byteJSON), nil
	}

	//获取该taskid的所有值
	conn.Send("MULTI")
	for _, key := range keys {
		conn.Send("HGETALL", key)
	}
	res, err := redis.Values(conn.Do("EXEC"))

	//生成表格内容
	for idx, data := range res {
		mapData, _ := redis.StringMap(data, nil)
		mapData["idx"] = strconv.Itoa(idx + 1)

		dest := make(map[string]interface{})
		for k, v := range mapData {
			dest[k] = interface{}(v)
		}
		resBody = append(resBody, dest)
	}

	//生成表格头
	if len(res) > 0 {
		mapData, _ := redis.StringMap(res[0], nil)
		for k, _ := range mapData {
			commentLabel := getCommentOfHeader(k)
			resTitle = append(resTitle, utils.TableTitle{Prop: k, Label: commentLabel, Sortable: false, Hidden: false})
		}
	}

	//返回拼接得结果
	tsb := utils.TableData{Cols: resTitle, Data: resBody}
	byteJSON, _ := json.Marshal(tsb)
	return string(byteJSON), nil
}

func getCommentOfHeader(key string) string {
	field, ok := models.WsSQLModel.Fields[key]
	if ok {
		return field.Comment
	} else {
		return key
	}
}

// @Title 根据query_id取得对应的queryJSON
// @Param	query_id path string true "query_id"
// @Summary 根据query_id取得对应的queryJSON
// @Description 根据query_id取得对应的queryJSON
// @Success 200 {string} json
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/log/query/:query_id [get]
func (this *TaskController) GetJsonById() {
	query_id := this.Ctx.Input.Params()[":query_id"]
	queryJson, err := models.QuerySubmitLog{}.GetQueryJson(query_id)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte(err.Error()))
	} else {
		this.Ctx.WriteString(queryJson)
	}
}

// @Title 上传文件获取文件内容
// @Summary 提交Jones任务
// @Param file  formData file true ""
// @Description 提交Jones任务
// @Success 200 []string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/file [post]
func (this *TaskController) File() {
	file, head, _ := this.GetFile("file")
	bytes, err := utils.ReadBytes(file, head.Size)

	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(bytes)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	var resStr string
	if result.Charset != "UTF-8" {
		resStr = utils.ConvertToString(string(bytes), result.Charset, "UTF-8")
	} else {
		resStr = string(bytes)
	}

	fileData := strings.Replace(resStr, "\r", "", -1)
	valueData := strings.Split(fileData, "\n")
	//if(len(valueData) > 5000) {
	//	this.Ctx.ResponseWriter.WriteHeader(400)
	//	this.Ctx.ResponseWriter.Write([]byte("文件行数需小于500"))
	//	return
	//}
	var strArr []string
	for _, str := range valueData {
		if str != "" {
			strArr = append(strArr, strings.Replace(strings.TrimSpace(str), "'", "", -1))
		}
	}
	JSONData, _ := json.Marshal(strArr)
	this.Ctx.WriteString(string(JSONData))
}

// @Title 重跑某个任务
// @Summary 重跑某个任务
// @Description 重跑某个任务
// @Success 200 []string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/reRun [post]
func (this *TaskController) ReRun() {
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var attr map[string]string
	json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	fmt.Println(attr)
	queryId := attr["queryId"]
	submitLog := models.QuerySubmitLog{QueryId: queryId}
	submitLog.GetSubmitLog()
	nowTime := time.Now()

	var params map[string]string
	json.Unmarshal([]byte(submitLog.QueryJson), &params)
	//解析dsid字段
	dsid, err := strconv.ParseInt(params["dsid"], 10, 64)
	if err != nil {
		panic(err.Error())
	}

	// 生成QueryInfo
	queryInfo := pb.QueryInfo{
		QueryId:    utils.GetMd5ForQuery(fmt.Sprint(queryId), this.userName, nowTime.String()), // query id
		Sql:        []string{submitLog.SqlStr},                                                 // sql statement
		Session:    this.userName,                                                              // session :this.userName
		Engine:     submitLog.Engine,                                                           // 引擎
		Dimensions: []string{""},                                                               // 维度表头
		Outputs:    []string{""},                                                               // 输出表头
		Filters:    []string{""},                                                               // 过滤字段
		Status:     pb.QueryStatus_Waiting,                                                     // 状态
		LastUpdate: nowTime.UnixNano() / 1e6,                                                   // 最后更新时间戳
		CreateTime: nowTime.UnixNano() / 1e6,                                                   // 创建时间戳
		Tables:     []string{""},                                                               // 查询所用表
		Type:       submitLog.Type,                                                             // 查询类别
		Dsid:       dsid,
	}

	// 将QueryInfo入库MySql
	models.QuerySubmitLog{}.InsertQuerySubmitLog(queryInfo, submitLog.QueryName, submitLog.QueryJson)
	// 编码QueryInfo
	queryBytes, _ := proto.Marshal(&queryInfo)
	// 得到MQ
	models.RabbitMQModel{}.Publish(queryBytes)

	response, _ := json.Marshal(map[string]string{"queryId": queryInfo.QueryId})
	this.Ctx.WriteString(string(response))
}

// @Title 查看任务的异常
// @Summary 查看任务的异常
// @Description 查看任务的异常
// @Success 200 []string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/checkException [post]
func (this *TaskController) CheckException() {
	var attr map[string]string
	json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	queryId := attr["queryId"]
	submitLog := models.QuerySubmitLog{QueryId: queryId}
	submitLog.GetSubmitLog()

	response, _ := json.Marshal(map[string]string{"res": submitLog.ExceptionMsg})
	this.Ctx.WriteString(string(response))

}

// @Title 查看queryLog
// @Summary 查看queryLog
// @Description 查看queryLog
// @Success 200 []string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/queryLog/:queryId [get]
func (this *TaskController) GetSubmitQueryLog() {
	queryId := this.Ctx.Input.Params()[":queryId"]
	submitLog := models.QuerySubmitLog{QueryId: queryId}
	submitLog.GetSubmitLog()

	response, _ := json.Marshal(submitLog)
	this.Ctx.WriteString(string(response))

}

// @Title 更改任务所属人
// @Summary 更改任务所属人
// @Description 更改任务所属人
// @Success 200 []string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/changeTaskUser [post]
func (this *TaskController) ChangeTaskUser() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	taskId := attr["assignTaskId"]
	userName := attr["assignUserName"]
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE query_submit_log SET session=? WHERE query_id=?", userName, taskId).Exec()
	if err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	}
}

// @Title 删除任务记录
// @Summary 删除任务记录
// @Description 删除任务记录
// @Success 200 []string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/deleteTaskRecoder/:taskId [delete]
func (this *TaskController) DeleteTaskRecoder() {
	taskId := this.Ctx.Input.Params()[":taskId"]
	o := orm.NewOrm()
	_, err := o.Raw("DELETE FROM query_submit_log WHERE query_id = ?", taskId).Exec()
	if err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	}
}

// @Title 杀死spark任务
// @Summary 杀死spark任务
// @Description 杀死spark任务
// @Success 200 []string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /task/killByJobId [post]
func (this *TaskController) KillByJobId() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	queryId := attr["queryId"]
	reqURL := utils.SparkRoute(models.ServerInfo.SparkHost, models.ServerInfo.SparkPort, "killByJobId")
	req := httplib.Post(reqURL).SetTimeout(time.Second*30, time.Second*30).Body(queryId)
	bytes, err := req.Bytes()
	if err != nil {
		models.Logger.Error("%s", err.Error())
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	} else {
		var res utils.SparkServerResponse
		_ = json.Unmarshal(bytes, &res)
		if res.Flag {
			this.Ctx.WriteString(utils.VueResponse{
				Res:  true,
				Info: "删除成功",
			}.String())
		} else {
			this.Ctx.WriteString(utils.VueResponse{
				Res:  false,
				Info: res.Msg,
			}.String())
		}
	}
}
