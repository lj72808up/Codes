package models

import (
	"AStar/models/airflowDag"
	adtl_platform "AStar/protobuf"
	"AStar/utils"
	"AStar/utils/Sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"github.com/flosch/pongo2"
	"github.com/golang/protobuf/proto"
	"regexp"
	"strings"
	"time"
)

type DagAirFlow struct {
	Id             int64                `orm:"pk;auto"`
	Name           string               `orm:"size(50);unique"`
	ExhibitName    string               `orm:"size(50)"`
	FileName       string               `orm:"size(255)"`
	StartDate      string               `orm:"size(20)"`
	MailList       string               `orm:"size(255)"`
	CronExpression string               `orm:"size(20)"`
	Delay          int64                `orm:"size(4)"`
	Description    string               `orm:"type(text)"`
	User           string               `orm:"size(50)"`
	ModifyTime     string               `orm:"size(20)"`
	Status         airflowDag.DagStatus `orm:"type(integer)"`
	IsOpen         int                  `orm:"size(1)" field:"Dag开关"`
	//LastScheduleRun string 	 `orm:"size(255)" field:"Last run"`
	//LastScheduleStatus string `orm:"default(notrun)"  field:"Last run status"`
	//Owners         string 	 `orm:"default(datacenter)" field:"所有者"`
}

func (this *DagAirFlow) GetByName(dagName string) {
	o := orm.NewOrm()
	_ = o.Raw("select * from dag_air_flow where name=?", dagName).QueryRow(this)
}

func (this *DagAirFlow) Add() {
	o := orm.NewOrm()
	if _, err := o.Insert(this); err != nil {
		panic(err)
	}
}

func (this *DagAirFlow) UpdateStatus() {
	o := orm.NewOrm()
	var creatingDags []DagAirFlow
	_, _ = o.Raw("SELECT * FROM dag_air_flow WHERE status = ?", airflowDag.NEW).QueryRows(&creatingDags)
	var removingDags []DagAirFlow
	_, _ = o.Raw("SELECT * FROM dag_air_flow WHERE status = ?", airflowDag.REMOVING).QueryRows(&removingDags)

	if (len(creatingDags) + len(removingDags)) == 0 {
		Logger.Info("airflow暂无新增或删除中的任务")
		return
	}
	// 1. 访问删除接口
	for _, dag := range removingDags {
		DeleteDagByApi(dag.Name)
		//fmt.Println(dag)
	}

	// 2. 访问最终列表
	dagListUrl := fmt.Sprintf("http://api.airflow.adtech.sogou/api/experimental/dag_list?owner=%s", "datacenter")
	req := httplib.Get(dagListUrl).SetTimeout(time.Minute, time.Minute)
	bytes, err := req.Bytes()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dagList := airflowDag.DagListView{}
	_ = json.Unmarshal(bytes, &dagList)

	dagMap := make(map[string]bool) // 查找状态
	for _, dag := range dagList.Data {
		dagMap[dag.DagId] = true
	}

	// 3. 数据库获取所有创建中的
	for _, dag := range creatingDags {
		if dagMap[dag.Name] == true { //
			dag.Status = airflowDag.ESTABLISH
			_, _ = o.Update(&dag, "status")
		}
	}
	// 4. 数据库获取所有删除中的
	for _, dag := range removingDags {
		if dagMap[dag.Name] == false {
			dag.Status = airflowDag.REMOVED
			_, _ = o.Update(&dag, "status")
		}
	}
}

type DagAirFlowImg struct {
	Id      int64  `orm:"pk;auto"`
	Img     string `orm:"size(255)"`
	Command string `orm:"size(255)"`
}

func GetAllDagImgs() []DagAirFlowImg {
	o := orm.NewOrm()
	o.Begin()
	var imgs []DagAirFlowImg
	_, _ = o.Raw("SELECT * FROM dag_air_flow_img").QueryRows(&imgs)
	return imgs
}

type DagView struct {
	Edges     []DagEdge `json:"edges"`
	Nodes     []DagNode `json:"nodes"`
	DagConfig DagConfig `json:"dagConfig"`
}

type CreateOutputView struct {
	OutTableType     string            `json:"outTableType"`
	TableName        string            `json:"tableName"`
	NewOrOldOutTable string            `json:"newOrOldOutTable"`
	TableDesc        string            `json:"tableDesc"`
	StorageEngine    string            `json:"storageEngine"`
	Sql              string            `json:"sql"`
	StaticFields     map[string]string `json:"sqlFields,omitempty"`
	TaskType         string            `json:"taskType"` // 小时级还是天级
	InputType        string            `json:"inputType"`
	TaskTypeInfo     TaskTypeInfo      `json:"taskTypeInfo"`
	Params           map[string]string `json:"params,omitempty"`
	Lvpi             LvpiNode          `json:"lvpi,omitempty"`
	ClickHouse       CkNode            `json:"clickhouse,omitempty"`
}

type TaskTypeInfo struct {
	Dsid             int64               `json:"dsid"`
	HdfsLocation     string              `json:"hdfsLocation"`
	Splitter         string              `json:"splitter"`
	ValidFieldCount  int                 `json:"validFieldCount"`
	Schema           []map[string]string `json:"schema"`
	ExtendSchema     []map[string]string `json:"extendSchema"`
	HasExtend        bool                `json:"hasExtend"`
	HdfsDoneLocation string              `json:"hdfsDoneLocation"`
	IsCheckDone      bool                `json:"isCheckDone"`
}

type LvpiNode struct {
	DsId      int64  `json:"dsid"`
	JdbcStr   string `json:"jdbcStr"`
	TableName string `json:"tableName"`
}

type CkNode struct {
	TableName     string   `json:"tableName"`
	Indices       string   `json:"indexStr,omitempty"`
	Partitions    string   `json:"partitionStr,omitempty"`
	IndicesArr    []string `json:"indices,omitempty"`
	PartitionsArr []string `json:"partitions,omitempty"`
}

func (this *DagView) PrePostDag(roles []string, user string) {
	outputs := make([]CreateOutputView, 0)

	sqlTemplate := SqlTemplate{}
	for _, node := range this.Nodes {
		// 替换参数和内置变量, 形成可执行SQL
		sqlParams := make([]map[string]string, 0)
		_ = json.Unmarshal([]byte(node.SqlParam), &sqlParams)
		runSql := sqlTemplate.GenerateSql(node.OriginalSql, sqlParams, "dag")
		// 字段及其注解
		staticFieldsDescription := make(map[string]string)
		for _, f := range node.StaticFields {
			staticFieldsDescription[f.FieldName] = f.Description
		}

		var req CreateOutputView
		dummyTaskTypeInfo := TaskTypeInfo{
			Dsid:             0,
			HdfsLocation:     "",
			Splitter:         "",
			ValidFieldCount:  0,
			Schema:           []map[string]string{{"id": "", "name": ""}},
			ExtendSchema:     []map[string]string{{"id": "", "name": ""}},
			HasExtend:        false,
			HdfsDoneLocation: "",
			IsCheckDone:      false,
		}
		if node.OutTableType == "hive" {
			req = CreateOutputView{
				OutTableType:     node.OutTableType,
				TableName:        node.OutTableName,
				NewOrOldOutTable: node.NewOrOldOutTable,
				TableDesc:        node.OutTableDesc,
				StorageEngine:    node.StorageEngine,
				Sql:              runSql,
				StaticFields:     staticFieldsDescription,
				TaskType:         this.DagConfig.DagType, // 小时级/天级
				InputType:        node.TaskType,          // hdfs还是hive
				TaskTypeInfo:     node.TaskTypeInfo,
				Params:           make(map[string]string, 0), // 占位符, 没有会导致json中存在null
			}
		} else if node.OutTableType == "mysql" {
			lvpi := node.Lvpi
			lvpi.TableName = node.OutTableName
			lvpi.GetJdbcString() // 补充 jdbcString 字段

			req = CreateOutputView{
				OutTableType:     node.OutTableType,
				TableName:        node.OutTableName,
				NewOrOldOutTable: node.NewOrOldOutTable,
				TaskTypeInfo:     dummyTaskTypeInfo,
				TableDesc:        node.OutTableDesc,
				Sql:              runSql,
				StaticFields:     staticFieldsDescription,
				Lvpi:             lvpi,
			}
		} else if node.OutTableType == "clickhouse" {

			ck := node.ClickHouse
			ck.TableName = node.OutTableName
			ck.PartitionsArr = strings.Split(ck.Partitions, ",")
			ck.IndicesArr = strings.Split(ck.Indices, ",")

			req = CreateOutputView{
				OutTableType:     node.OutTableType,
				TableName:        node.OutTableName,
				NewOrOldOutTable: node.NewOrOldOutTable,
				TaskTypeInfo:     dummyTaskTypeInfo,
				TableDesc:        node.OutTableDesc,
				Sql:              runSql,
				StaticFields:     staticFieldsDescription,
				ClickHouse:       ck,
			}
		} else {
			continue
		}

		outputs = append(outputs, req)
	}
	// 访问hive接口
	if len(outputs) > 0 {
		this.createOutPut(outputs, roles, user)
	}
}

func (this *LvpiNode) GetJdbcString() string {
	conn := ConnectionManager{}
	conn.GetById(this.DsId)
	// jdbc:mysql://10.135.73.50:5029/?useSSL=false&useUnicode=true&characterEncoding=UTF-8
	jdbcString := fmt.Sprintf("jdbc:mysql://%s:%s/%s?user=%s&password=%s&useSSL=false&useUnicode=true",
		conn.Host, conn.Port, conn.Database, conn.User, conn.Passwd)
	Logger.Info(jdbcString)
	this.JdbcStr = jdbcString
	return jdbcString
}

func (this *DagView) createOutPut(outputs []CreateOutputView, roles []string, user string) {
	importMeta := ImportMetaData{}
	param, _ := json.Marshal(outputs)
	// 访问sparkServer的接口, 先创建出表, 再导入元数据
	url := utils.SparkRoute(ServerInfo.SparkHost, ServerInfo.SparkPort, "dag_createTable")
	req := httplib.Post(url).SetTimeout(time.Minute*1, time.Minute*5).Body(param)
	bytes, err := req.Bytes()
	Logger.Info("sparkServer:%s", string(bytes))
	if err != nil {
		panic(err.Error())
	}
	hiveTables := &adtl_platform.HiveTableSet{}
	_ = proto.Unmarshal(bytes, hiveTables)

	if len(bytes) > 0 && ((hiveTables.Tables == nil) || (len(hiveTables.Tables) == 0)) {
		Logger.Error(string(bytes))
		var res utils.SparkServerResponse
		_ = json.Unmarshal(bytes, &res)
		panic(res.Msg)
	}

	// hive 元数据导入
	if (hiveTables.Tables != nil) && (len(hiveTables.Tables) > 0) {
		// (1) 获取hive类型对应的connId
		connId := GetAllConn("hive")[0].Dsid

		// (2) 导入元数据
		o := orm.NewOrm()
		now := utils.TimeFormat(time.Now())
		partitionName := getPartitionNameByDagType(this.DagConfig.DagType)
		tableType := "hive"

		for _, hTable := range hiveTables.Tables {
			hName := hTable.NameTable
			dbName := "datacenter"
			tbName := hName
			expandTableName := strings.Split(hName, ".")
			if len(expandTableName) > 1 {
				dbName = expandTableName[0]
				tbName = expandTableName[1]
			}
			hTable.NameTable = tbName
			// 检测该表是否已经导入过
			if !CheckDuplicateTable(tableType, dbName, tbName) {
				tid := importMeta.InsertTableFromMeta(dbName, hTable, o, roles, fmt.Sprint(connId))
				newVersion := importMeta.InsertVersionFromMeta(tid, now, o)
				// 没有数据的分区, show partitions显示不出分区名
				if partitionName != "" {
					hTable.ExpendInfo["partitions"] = partitionName
				}
				importMeta.InsertFieldFromMeta(tid, newVersion, hTable, o)
				AddTablePrivilegeMapping(tid, roles)
			}
		}
	}

	// 对于mysql和ck, 接口不会返回元数据, 用astar导
	this.importMetaDataFromType(user, outputs)
}

func (this *DagView) importMetaDataFromType(user string, outputs []CreateOutputView) {
	ckConn := GetAllConn("clickhouse")[0]
	for _, output := range outputs {
		db, tb := utils.TransTableIdentifier(output.TableName)
		if output.OutTableType == "mysql" {
			param := map[string]string{
				"tableType":    "mysql",
				"importId":     fmt.Sprintf("%d", output.Lvpi.DsId),
				"databaseName": db,
				"tableName":    tb,
			}
			GenerateFieldByDefSingle(param, user)
		} else if output.OutTableType == "clickhouse" {
			param := map[string]string{
				"tableType":    "clickhouse",
				"importId":     fmt.Sprintf("%d", ckConn.Dsid),
				"databaseName": db,
				"tableName":    tb,
			}
			GenerateFieldByDefSingle(param, user)
		}
	}
}

// dag输出表的分区信息, 小时级是hour, 天级是dt
func getPartitionNameByDagType(dagType string) string {
	partitionName := ""
	switch dagType {
	case "day":
		partitionName = "dt"
	case "hour":
		partitionName = "hour"
	}
	return partitionName
}

func (this *DagView) UpdateDag(msg string, dagName string) {
	now := utils.TimeFormat(time.Now())
	o := orm.NewOrm()
	dag := DagAirFlow{}
	_ = o.Raw("SELECT * FROM dag_air_flow WHERE name=?", dagName).QueryRow(&dag)

	if dag.Id == 0 {
		panic(fmt.Sprintf("未找到名字为%s的DAG", dagName))
	}
	// dag文件
	dagFile := dag.FileName
	content := this.MakeAirflowFile()
	_, _ = utils.WriteFileFullPath(dagFile, content)
	// 推送
	airbase := beego.AppConfig.String("airflowFileRoot")
	rsyncCmd1 := fmt.Sprintf("rsync -aP --delete %s %s", airbase, beego.AppConfig.String("airflowRemoteFileRoot1"))
	rsyncCmd2 := fmt.Sprintf("rsync -aP --delete %s %s", airbase, beego.AppConfig.String("airflowRemoteFileRoot2"))
	cmds := []string{rsyncCmd1, rsyncCmd2}
	for _, cmd := range cmds {
		fmt.Println("rsync 执行:", cmd)
		splits := strings.Split(cmd, " ")
		utils.ExecShell(splits[0], splits[1:])
	}
	// 修改"description, modify_time, status=new"字段
	dag.Description = msg
	dag.ModifyTime = now
	dag.Status = airflowDag.NEW
	if _, err := o.Update(&dag, "Description", "ModifyTime", "Status"); err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
}

func (this *DagView) SaveDagWhenFail(user string, body string) {
	now := utils.TimeFormat(time.Now())
	toFile := makeFilePath(this.DagConfig.Name, user)
	// 不存保存, 存在则更新
	dagAirFlow := DagAirFlow{
		Name:           this.DagConfig.Name,
		ExhibitName:    this.DagConfig.ExhibitName,
		FileName:       toFile,
		StartDate:      this.DagConfig.StartDate,
		MailList:       this.DagConfig.MailList,
		CronExpression: this.DagConfig.CronExpression,
		Delay:          this.DagConfig.Delay,
		Description:    body, //string(this.Ctx.Input.RequestBody),
		User:           user,
		ModifyTime:     now,
		Status:         airflowDag.BUILDFAIL, // 创建失败
	}
	dagAirFlow.Add()
}

func makeFilePath(dagName string, user string) string {
	fileName := fmt.Sprintf("%s.py", utils.GetMd5(dagName))
	airbase := beego.AppConfig.String("airflowFileRoot")
	dateUserbase := fmt.Sprintf("%s/%s/%s", airbase, utils.DateFormatWithoutLine(time.Now()), user)
	utils.MkDir(dateUserbase)
	toFile := fmt.Sprintf("%s/%s", dateUserbase, fileName)
	return toFile
}

type RelationReq struct {
	Sql      string `json:"sql"`
	OutTable string `json:"out"`
	OutType  string `json:"outType"`
	TaskType string `json:"taskType"`
	Roles    string `json:"roles"`
}

func (this *DagView) AfterPostDag(rolesStr string) {
	upstreams := this.makeHiveRelationReqs(rolesStr)
	for _, stream := range upstreams {
		outType := stream["outType"]
		outDB, outTB := transTableName(stream["outTable"]) // 输出表
		outTid := GetTableId(outType, outDB, outTB)
		if outTid == 0 {
			Logger.Error("%s is not exist in table_meta_data", stream["outTable"])
			panic(fmt.Sprintf("%s is not exist in table_meta_data", stream["outTable"]))
		}
		upTables := strings.Split(stream["upstream"], ",") // 逗号分隔的上游表皮
		for _, upTable := range upTables {
			upDB, upTB := transTableName(upTable)
			upStreamTableType := stream["taskType"]
			upId := GetTableId(upStreamTableType, upDB, upTB)
			relation := TableRelation{
				TableAId: outTid,
				TableBId: upId,
			}
			relation.AddOne()
		}
	}
}

func transTableName(tableName string) (db string, tb string) {
	strs := strings.Split(tableName, ".")
	if len(strs) > 1 {
		return strs[0], strs[1]
	} else {
		return "datacenter", tableName
	}
}

func (this *DagView) makeHiveRelationReqs(rolesStr string) []map[string]string {
	var hiveReqs []RelationReq
	sqlTemplate := SqlTemplate{}
	nodes := this.Nodes
	for _, node := range nodes {
		if node.TaskType == "hive" && (node.OutTableType == "hive" ||
			node.OutTableType == "mysql" || node.OutTableType == "clickhouse") {
			// 生成代yyyyMMdd的参数, 替换参数后的sql
			var sqlParams []map[string]string
			_ = json.Unmarshal([]byte(node.SqlParam), &sqlParams)
			runSql := sqlTemplate.GenerateSql(node.OriginalSql, sqlParams, "dag")
			req := RelationReq{
				Sql:      runSql,
				OutTable: node.OutTableName,
				OutType:  node.OutTableType,
				TaskType: node.TaskType,
				Roles:    rolesStr,
			}
			hiveReqs = append(hiveReqs, req)
		}
	}
	// visit
	if len(hiveReqs) > 0 {
		param, _ := json.Marshal(hiveReqs)
		url := utils.SparkRoute(ServerInfo.SparkHost, ServerInfo.SparkPort, "dag_getUpstream")
		resp := httplib.Post(url).SetTimeout(time.Minute*1, time.Minute*5).Body(param)
		bytes, err := resp.Bytes()
		// response
		if err != nil {
			Logger.Error(err.Error())
		}
		var upstreams []map[string]string // {outTable:***  upstream:***}
		_ = json.Unmarshal(bytes, &upstreams)
		return upstreams
	}
	return []map[string]string{}
}

func (this *DagView) PostDag(dag DagView, user string, body string) {
	now := utils.TimeFormat(time.Now())

	// 替换文件
	content := this.MakeAirflowFile()
	airbase := beego.AppConfig.String("airflowFileRoot")
	//last write file method
	/*
		fileName := fmt.Sprintf("%s.py", utils.GetMd5(this.DagConfig.Name))
		dateUserbase := fmt.Sprintf("%s/%s/%s", airbase, utils.DateFormatWithoutLine(time.Now()), user)
		toFile, _ := utils.WriteFile(dateUserbase, fileName, content)*/

	toFile := makeFilePath(this.DagConfig.Name, user)
	_, _ = utils.WriteFileFullPath(toFile, content)

	// 插入数据库
	nodes := dag.Nodes
	// sql中的yyyyMMdd和HH参数由镜像替换
	for i, _ := range nodes {
		var param []map[string]string
		var sqlParams []map[string]string
		_ = json.Unmarshal([]byte(nodes[i].SqlParam), &sqlParams)
		for _, p := range sqlParams {
			if (p["name"] != "yyyyMMdd") && (p["name"] != "HH") {
				param = append(param, p)
			}
		}
	}

	dagAirFlow := DagAirFlow{
		Name:           this.DagConfig.Name,
		ExhibitName:    this.DagConfig.ExhibitName,
		FileName:       toFile,
		StartDate:      this.DagConfig.StartDate,
		MailList:       this.DagConfig.MailList,
		CronExpression: this.DagConfig.CronExpression,
		Delay:          this.DagConfig.Delay,
		Description:    body, //string(this.Ctx.Input.RequestBody),
		User:           user,
		ModifyTime:     now,
		Status:         airflowDag.NEW,
	}
	dagAirFlow.Add()

	// 推送
	rsyncCmd1 := fmt.Sprintf("rsync -aP --delete %s %s", airbase, beego.AppConfig.String("airflowRemoteFileRoot1"))
	rsyncCmd2 := fmt.Sprintf("rsync -aP --delete %s %s", airbase, beego.AppConfig.String("airflowRemoteFileRoot2"))
	cmds := []string{rsyncCmd1, rsyncCmd2}
	for _, cmd := range cmds {
		fmt.Println("rsync 执行:", cmd)
		splits := strings.Split(cmd, " ")
		utils.ExecShell(splits[0], splits[1:])
	}
}

type DagTriggerTask struct {
	DagName      string
	NodeList     []DagNode
	ScheduleTime string
	UserName     string
	HistoryId    int64
}

func (this *DagView) PostDagToQueue(dag DagView, start string, end string, userName string) {
	nodes := dag.Nodes
	edges := dag.Edges
	nodeSize := len(nodes)
	inDegree := make([]int, nodeSize)
	ans := []DagNode{} //  按照图的依赖关系解析出的列表
	idMapping := make(map[int64]int)
	for idx, node := range nodes {
		idMapping[node.Id] = idx
	}

	for _, edge := range edges {
		inDegree[idMapping[edge.DstId]]++
	}
	num := 0
	for num < nodeSize {
		for i := 0; i < nodeSize; i++ {
			if inDegree[i] == 0 {
				num++
				ans = append(ans, nodes[i])
				inDegree[i] = -1
				for _, edge := range edges {
					if edge.SrcId == nodes[i].Id && inDegree[idMapping[edge.DstId]] > 0 {
						inDegree[idMapping[edge.DstId]]--
					}
				}
			}
		}
	}

	dateList := utils.GetDateRange(start, end)
	hours := []string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09",
		"10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23",}
	if dag.DagConfig.DagType == "hour" {
		dummyDateList := []string{}
		for _, date := range dateList {
			for _, hour := range hours {
				dummyDateList = append(dummyDateList, fmt.Sprintf("%s%s", date, hour))
			}
		}
		dateList = dummyDateList
	}

	for _, dateStr := range dateList {
		println("放入dag:  " + dag.DagConfig.Name + ";;" + dateStr)
		hid, err := (&DagTriggerHistory{
			DagName:      dag.DagConfig.Name,
			Username:     userName,
			ScheduleTime: dateStr,
			Status:       "pending",
		}).Insert()
		if err != nil {
			utils.Logger.Error("%s", err.Error())
		}
		ManualTriggerDagChannel <- DagTriggerTask{
			DagName:      dag.DagConfig.Name,
			NodeList:     ans,
			ScheduleTime: dateStr,
			UserName:     userName,
			HistoryId:    hid,
		}
	}

}

func (this *DagView) MakeAirflowFile() *string {
	tpl, err := pongo2.FromFile("models/airflowDag/airflow.tmpl")
	if err != nil {
		panic(err)
	}
	// 解析年月日
	date, _ := time.Parse("2006-01-02", this.DagConfig.StartDate)
	param := make(map[string]interface{})
	// Dag
	param["start_year"] = date.Year()
	param["start_month"] = date.Month()
	param["start_day"] = date.Day()

	param["dagType"] = this.DagConfig.DagType
	param["delay"] = this.DagConfig.Delay
	param["retries"] = this.DagConfig.Retries
	param["retry_delay"] = this.DagConfig.RetryDelay

	emails := strings.Split(this.DagConfig.MailList, ",")
	emailUsers := fmt.Sprintf("\"%s\"", strings.Join(emails, "\",\""))
	param["emailUsers"] = emailUsers

	param["cron_expression"] = this.DagConfig.CronExpression
	param["dag_name"] = this.DagConfig.Name

	jonesJdbc, jdbcUser, jdbcPasswd := getJdbcUrl()

	// Tasks
	var tasks []map[string]string
	for _, node := range this.Nodes {
		taskParam := make(map[string]string)
		taskParam["cmd"] = node.Command

		envMap := make(map[string]string)
		for _, m := range node.Envs {
			envMap[m["env_key"]] = m["env_val"]
		}
		envMap["JONES_JDBC"] = jonesJdbc
		envMap["JONES_JDBC_USER"] = jdbcUser
		envMap["JONES_JDBC_PASSWORD"] = jdbcPasswd
		envMap["JONES_DAG_NAME"] = this.DagConfig.Name
		envMap["JONES_TASK_NAME"] = node.Name
		envMap["user"] = "adtd_platform"
		envMap["password"] = "adtd_platform"
		envRes, _ := json.Marshal(envMap)

		taskParam["env_json"] = string(envRes)
		taskParam["task_var"] = fmt.Sprintf("node_%d", node.Id)
		taskParam["task_name"] = node.Name
		taskParam["img"] = node.Image
		tasks = append(tasks, taskParam)
	}
	param["tasks"] = tasks

	// Edges
	param["edges"] = this.Edges

	out, err := tpl.Execute(param)
	if err != nil {
		panic(err)
	}
	return &out
}

func getJdbcUrl() (string, string, string) {
	//url := "yuhancheng:yuhancheng@tcp(10.139.36.81:3306)/adtl_test?charset=utf8&loc=Asia%2FShanghai"
	url := beego.AppConfig.String("datasourcename")
	r, _ := regexp.Compile("^(.*):(.*)@tcp\\((.*)\\)/(.*)\\?(.*)$")
	params := r.FindStringSubmatch(url)

	user := params[1]
	passwd := params[2]
	hostPort := params[3]
	dataBase := params[4]

	jdbcUrl := fmt.Sprintf("jdbc:mysql://%s/%s",
		hostPort, dataBase)
	fmt.Println(jdbcUrl)
	return jdbcUrl, user, passwd
}

type DagEdge struct {
	Id    int64 `json:"id"`
	SrcId int64 `json:"src_node_id"`
	DstId int64 `json:"dst_node_id"`
}

type DagNode struct {
	Id               int64               `json:"id"`
	Name             string              `json:"name"`
	TaskDesc         string              `json:"taskDesc"`
	TaskType         string              `json:"taskType"`
	TaskTypeInfo     TaskTypeInfo        `json:"taskTypeInfo"`
	Image            string              `json:"image"`
	Command          string              `json:"command"`
	OriginalSql      string              `json:"originalSql"`
	SqlParam         string              `json:"sqlParam"`
	StaticFields     []StaticField       `json:"sqlFields"`
	OutTableName     string              `json:"outTableName"`
	OutTableType     string              `json:"outTableType"`
	OutTableDesc     string              `json:"outTableDesc"`
	NewOrOldOutTable string              `json:"newOrOldOutTable"`
	StorageEngine    string              `json:"storageEngine"`
	Envs             []map[string]string `json:"envs"`
	Lvpi             LvpiNode            `json:"lvpi"`
	ClickHouse       CkNode              `json:"clickhouse"`
}

type StaticField struct {
	FieldName   string `json:"fieldName"`
	Description string `json:"description"`
}

type DagConfig struct {
	Name           string `json:"name"`
	ExhibitName    string `json:"exhibitName"`
	StartDate      string `json:"startDate"`
	MailList       string `json:"mailList"`
	CronExpression string `json:"cronExpression"`
	DagType        string `json:"dagType"`
	Delay          int64  `json:"delay"`
	Retries        int64  `json:"retries"`
	RetryDelay     int64  `json:"retryDelay"`
}

func DeleteDagByApi(dagName string) {
	url := "http://api.airflow.adtech.sogou/api/experimental/dag_delete"
	param, _ := json.Marshal(map[string]string{
		"dag_id": dagName,
	})
	fmt.Println("restful删除Dag", dagName)

	req := httplib.Post(url).SetTimeout(time.Minute, time.Minute).Body(param)
	_, err := req.Bytes()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetOutTables(role string, tableType string, dsId int64) []string {
	o := orm.NewOrm()
	tables := make([]string, 0)
	if tableType == "mysql" {
		sql := `SELECT CONCAT(t1.database_name,'.',t1.table_name) FROM table_meta_data t1
LEFT JOIN 
(SELECT table_id FROM group_metatable_mapping WHERE group_name=? AND (read_flag=1 OR exec_flag=1)) t2
ON t1.tid = t2.table_id
WHERE t2.table_id IS NOT NULL
AND t1.table_type=? AND t1.dsid=?`
		_, _ = o.Raw(sql, role, tableType, dsId).QueryRows(&tables)
	} else {
		sql := `SELECT CONCAT(t1.database_name,'.',t1.table_name) FROM table_meta_data t1
LEFT JOIN 
(SELECT table_id FROM group_metatable_mapping WHERE group_name=? AND (read_flag=1 OR exec_flag=1)) t2
ON t1.tid = t2.table_id
WHERE t2.table_id IS NOT NULL
AND t1.table_type=?`
		_, _ = o.Raw(sql, role, tableType).QueryRows(&tables)
	}
	return tables
}

func CheckMeta(sql string, outTable string, outType string, taskType string) (bool, string) {
	if taskType == "hive" {
		url := utils.SparkRoute(ServerInfo.SparkHost, ServerInfo.SparkPort, "dag_checkOutTable")
		param, _ := json.Marshal(map[string]string{
			"sql":      sql,
			"outTable": outTable,
			"outType":  outType,
		})
		req := httplib.Post(url).SetTimeout(time.Minute*1, time.Minute*5).Body(param)
		bytes, err := req.Bytes()
		if err != nil {
			panic(err.Error())
		}
		var checkRes utils.SparkServerResponse
		_ = json.Unmarshal(bytes, &checkRes)
		return checkRes.Flag, checkRes.Msg
	} else {
		// 只校验字段个数
		db, tb := utils.TransTableIdentifier(outTable)
		var fieldCnt int
		o := orm.NewOrm()
		targetSql := `SELECT count(t2.fid)
FROM (
	SELECT version,
	    table_id
	FROM version_meta_data
	WHERE table_id = 
	    (SELECT tid
	    FROM table_meta_data
	    WHERE table_name=? AND database_name=? AND table_type=? LIMIT 1)
	    ORDER BY  create_time DESC LIMIT 1
) t1
LEFT JOIN field_meta_data t2
ON t1.version = t2.version
    AND t1.table_id = t2.table_id
`
		_ = o.Raw(targetSql, tb, db, outType).QueryRow(&fieldCnt)

		parser := Sql.SemanticParser{}
		fields := parser.GetFields(sql)
		return len(fields) == fieldCnt, ""
	}
}
