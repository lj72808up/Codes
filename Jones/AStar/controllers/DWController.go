package controllers

import (
	"AStar/models"
	"AStar/protobuf"
	pb "AStar/protobuf"
	"AStar/utils"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"github.com/golang/protobuf/proto"
	"github.com/pingcap/log"
	"github.com/saintfish/chardet"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DWController struct {
	BaseController
}

// @Title 获取Meta信息
// @Summary 获取Meta信息
// @Description 获取Meta信息
// @Success 200 {object} models.MetaTreeNode
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/metaAll [get]
func (this *DWController) MetaAll() {
	//metaUrl := utils.SparkRouteMap(this.controllerName, this.actionName)
	//metaUrl:="http://10.134.101.120:4803/get_meta"
	metaURL := utils.SparkRoute(models.ServerInfo.SparkHost, models.ServerInfo.SparkPort, models.RouteMap[this.controllerName+"_"+this.actionName])
	req := httplib.Get(metaURL).SetTimeout(time.Hour, time.Hour)
	bytes, error := req.Bytes()
	if error != nil {
		models.Logger.Error(error.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte("数仓信息获取失败"))
	}
	hivemeta := &adtl_platform.HiveMeta{}
	proto.Unmarshal(bytes, hivemeta)

	metaTree := models.MetaTreeNode{Label: "meta"}
	for _, database := range hivemeta.Databases {
		var tmpdatabase []models.MetaTreeNode
		for _, table := range database.Talbes {
			var tmptable []models.MetaTreeNode
			for _, field := range table.Fileds {
				tmptable = append(tmptable, models.MetaTreeNode{Label: field})
			}
			tmpdatabase = append(tmpdatabase, models.MetaTreeNode{table.NameTable, "", 0, tmptable})
		}
		metaTree.Children = append(metaTree.Children, models.MetaTreeNode{database.NameDatabase, "", 0, tmpdatabase})
	}
	jsonByte, _ := json.Marshal(metaTree)
	this.Ctx.WriteString(string(jsonByte))
}

// @Title 获取Meta信息
// @Summary 获取Meta信息
// @Description 获取Meta信息
// @Success 200 {object} models.MetaTreeNode
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/meta/get [post]
func (this *DWController) Meta() {
	query := models.SqlQueryModel{}
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	dsType := params["dsType"]
	dsId, _ := strconv.ParseInt(params["dsId"], 10, 64)
	model := query.GetMetaTreeByPrivilege(dsType, dsId, this.userName)
	res, _ := json.Marshal(model)
	this.Ctx.WriteString(string(res))
}

// @Title 获取sql的queryId
// @Param	sql formData string true "sql"
// @Summary 获取sql的queryId
// @Description 获取sql的queryId
// @Success 200 {string} queryId
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/sql/id [post]
func (this *DWController) QueryID() {
	sqlstr := this.GetString("sql")
	Md5Res := utils.GetMd5ForQuery(sqlstr, this.userName, time.Now().String())
	this.Ctx.WriteString(string(Md5Res))
}

// @Title 验证hive sql
// @Param	sql formData string true "sql"
// @Summary 获取sql的queryId
// @Description 获取sql的queryId
// @Success 200 {string} queryId
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/sqlVerify [post]
func (this *DWController) SqlVerify() {
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)

	reqURL := utils.SparkRoute(models.ServerInfo.SparkHost, models.ServerInfo.SparkPort, "sql_verify")
	param, _ := json.Marshal(map[string]string{
		"sql":     params["sql"],
		"sqlType": params["sqlType"],
		"dsId":    params["dsId"],
	})
	paramJson, _ := json.Marshal(param)
	models.Logger.Info("sqlVerify参数: %s", string(paramJson))
	req := httplib.Post(reqURL).SetTimeout(time.Hour, time.Hour).Body(param)

	var res utils.VueResponse
	bytes, err := req.Bytes()
	if err != nil {
		log.Error(err.Error())
		res = utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}
	} else {
		var sparkServerRes utils.SparkServerResponse
		_ = json.Unmarshal(bytes, &sparkServerRes)
		res = utils.VueResponse{
			Res:  sparkServerRes.Flag,
			Info: sparkServerRes.Msg,
		}
	}
	resBytes, _ := json.Marshal(res)
	this.Ctx.WriteString(string(resBytes))
}

// @Title 提交SQL任务
// @Param	sql formData string true "sql"
// @Param	queryId formData string true "queryId"
// @Summary 获取sql的queryId
// @Description 获取sql的queryId
// @Success 200 {object} models.TableTitle
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/sql/task [post]
func (this *DWController) QuerySQL() {
	NameStr := this.GetString("name")
	SQLStr := this.GetString("sql")
	Md5Res := this.GetString("queryId")
	queryWordsPartition := this.GetString("queryWordsPartition")
	sqlType := this.GetString("sqlType")
	dsIdStr := this.GetString("dsId")
	dsid, _ := strconv.ParseInt(dsIdStr, 10, 64)

	jsonMap := make(map[string]interface{}) //map[string]string{}
	filterObj := make(map[string]interface{})
	filterObj["name"] = "queryWordPartition"
	filterObj["value"] = queryWordsPartition
	filterObj["type"] = pb.FilterOperation_name[10]
	var filterList []map[string]interface{}
	filterList = append(filterList, filterObj)
	jsonMap["filters"] = filterList
	jsonMap["dsid"] = dsid
	jsonMap["sqlType"] = sqlType

	queryJSONStr, _ := json.Marshal(jsonMap)
	filterObjJson, _ := json.Marshal(filterObj)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()

	var queryEngine adtl_platform.QueryEngine // 查询引擎
	var queryType string

	switch sqlType {
	case "hive":
		queryEngine = pb.QueryEngine_spark
		queryType = "sql查询"
	case "mysql":
		queryEngine = pb.QueryEngine_jdbc
		queryType = "sql查询"
		//dsid, _ = strconv.ParseInt(dsIdStr, 10, 64)
	case "clickhouse":
		queryEngine = pb.QueryEngine_jdbc_clickhouse
		queryType = "sql查询"
	case "kylin":
		queryEngine = pb.QueryEngine_kylin
		queryType = "sql查询"
	}
	// build queryInfo
	session := this.userName
	if strings.Contains(NameStr, "_outputcsv") {
		session = "csv"
	}

	queryInfo := pb.QueryInfo{
		QueryId:    Md5Res,                            // query id
		Sql:        []string{SQLStr},                  // sql statement
		Session:    this.userName,                     // session
		Engine:     queryEngine,                       // 引擎
		Status:     pb.QueryStatus_CustomQueryWaiting, // 状态
		LastUpdate: time.Now().UnixNano() / 1e6,       // 最后更新时间戳
		CreateTime: time.Now().UnixNano() / 1e6,       // 创建时间戳
		Tables:     []string{},                        // 查询所用表
		Type:       queryType,                         // 查询类别"无线搜索查询"
		Filters:    []string{string(filterObjJson)},   // 过滤字段
		Dsid:       dsid,
	}

	strLen := len(SQLStr)
	if strLen > 100 {
		strLen = 100
	}
	// 编码QueryInfo
	models.QuerySubmitLog{}.InsertQuerySubmitLog(queryInfo, NameStr, string(queryJSONStr))
	// 得到MQ
	queryInfo.Session = session
	queryBytes, _ := proto.Marshal(&queryInfo)
	models.RabbitMQModel{}.Publish(queryBytes)

	this.Ctx.WriteString(utils.VueResponse{
		Res:  true,
		Info: "",
	}.String())
}

// @Title 提交自定义SQL任务
// @Param	sql formData string true "sql"
// @Param	queryId formData string true "queryId"
// @Summary 获取sql的queryId
// @Description 获取sql的queryId
// @Success 200 {object} models.TableTitle
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/sql/custom [post]
func (this *DWController) QuerySQLTask() {
	SQLStr := strings.TrimSpace(this.GetString("sql"))
	Md5Res := this.GetString("queryId")
	// build queryInfo
	queryInfo := adtl_platform.QueryInfo{
		QueryId:    Md5Res,
		Session:    this.userName,
		Sql:        []string{SQLStr},
		Status:     adtl_platform.QueryStatus_CustomQueryWaiting,
		LastUpdate: time.Now().UnixNano() / 1e6,
		Engine:     1,
		CreateTime: time.Now().UnixNano() / 1e6,
		Dimensions: []string{"imei", "加密安卓id"},
	}

	// marshal queryInfo
	queryBytes, _ := proto.Marshal(&queryInfo)

	// create request
	reqURL := utils.SparkRoute(models.ServerInfo.PaganiHost, models.ServerInfo.SparkPort, models.RouteMap[this.controllerName+"_"+this.actionName])
	req := httplib.Post(reqURL).SetTimeout(time.Hour, time.Hour).Body(queryBytes)

	// InsertQuerySubmitLog into mysql db to log this query
	models.QuerySubmitLog{}.InsertQuerySubmitLog(queryInfo, "", "")

	// it calls Response inner. Bytes returns the body []byte in response.
	bytes, error := req.Bytes()

	// error occurs so that query failed
	if error != nil {
		models.Logger.Error(error.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte("SQL提交失败,请联系管理员查看!"))
	}

	// prepare transform data
	tran := &adtl_platform.TranStream{}
	proto.Unmarshal(bytes, tran)

	var yks []map[string]string

	for _, line := range tran.Line {
		var tmap = make(map[string]string)
		for idxline, field := range line.Field {
			tmap[tran.Schema[idxline]] = field
		}
		yks = append(yks, tmap)
	}

	var resTitle []utils.TableTitle
	var resBody []map[string]interface{}

	for _, line := range tran.Line {
		var tmap = make(map[string]interface{})
		for idx, field := range line.Field {
			tmap[tran.Schema[idx]] = field
		}
		resBody = append(resBody, tmap)
	}

	for _, schema := range tran.Schema {
		resTitle = append(resTitle, utils.TableTitle{Prop: schema, Label: schema, Sortable: true, Hidden: false})
	}
	tsb := utils.TableData{Cols: resTitle, Data: resBody}
	data, _ := json.Marshal(tsb)

	this.Ctx.WriteString(string(data))
}

// @Title 获取前端各部分的具体选项
// @Summary 获取前端各部分的具体选项
// @Description 获取前端各部分的具体选项
// @Success 200 {object} models.FrontOptions
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/sql/options [get]
func (this *DWController) GetOptions() {

	var ff []models.FrontFilterModel
	orm.NewOrm().QueryTable("front_filter_model").All(&ff)
	tempMap := make(map[string]map[string]int)
	for _, obj := range ff {
		var frontMap map[string][]string
		json.Unmarshal([]byte(obj.FiltersJson), &frontMap)
		for key, strList := range frontMap {
			tmpSetMap := make(map[string]int)
			tempMap[key] = tmpSetMap
			for _, str := range strList {
				if str != "" {
					tempMap[key][str] = 1
				}
			}
		}
	}

	finalMap := make(map[string][]string)

	for key, strList := range tempMap {
		var tmpStrList []string
		for strKey, _ := range strList {
			tmpStrList = append(tmpStrList, strKey)
		}
		sort.Strings(tmpStrList)
		finalMap[key] = tmpStrList
	}

	data, _ := json.Marshal(finalMap)
	this.Ctx.WriteString(string(data))

	//data, _ := json.Marshal(models.FrontOptions)
	//this.Ctx.WriteString(string(data))

	//tmp := make(map[string][]string)
	//tmp["dqOptions"] = []string{"大区1","大区2"}
	//tmp["wzOptions"] = []string{"位置1","位置2"}
	//
	//data, _ := json.Marshal(tmp)
	//this.Ctx.WriteString(string(data))

}

// @Title 获取前端各部分的具体选项
// @Param	date body []string true "时间"
// @Summary 获取前端各部分的具体选项
// @Description 获取前端各部分的具体选项
// @Success 200 {object} models.FrontOptions
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/sql/options_wap/date [post]
func (this *DWController) GetOptionsDate() {
	var dateRange []string
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &dateRange)
	if err != nil {
		models.Logger.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		return
	}

	var ff []models.FrontFilterModel
	cond := orm.NewCondition()
	cond1 := cond.And("date__gte", dateRange[0]).And("date__lte", dateRange[1]).Or("date__eq", "20180901")
	orm.NewOrm().QueryTable("front_filter_model").
		SetCond(cond1).
		All(&ff)
	tempMap := make(map[string]map[string]int)
	for _, obj := range ff {
		var frontMap map[string][]string
		json.Unmarshal([]byte(obj.FiltersJson), &frontMap)
		for key, strList := range frontMap {
			tmpSetMap := make(map[string]int)
			_, ok := tempMap[key]
			if !ok {
				tempMap[key] = tmpSetMap
			}
			for _, str := range strList {
				if str != "" {
					tempMap[key][str] = 1
				}
			}
		}
	}

	finalMap := make(map[string][]string)

	for key, strList := range tempMap {
		var tmpStrList []string
		for strKey, _ := range strList {
			tmpStrList = append(tmpStrList, strKey)
		}
		sort.Strings(tmpStrList)
		finalMap[key] = tmpStrList
	}

	data, _ := json.Marshal(finalMap)
	this.Ctx.WriteString(string(data))

}

// @Title 获取前端各部分的具体选项
// @Param	date body []string true "时间"
// @Summary 获取前端各部分的具体选项
// @Description 获取前端各部分的具体选项
// @Success 200 {object} models.FrontOptions
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/sql/options_pc/date [post]
func (this *DWController) GetOptionsPCDate() {
	var dateRange []string
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &dateRange)
	if err != nil {
		models.Logger.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		return
	}

	var ff []models.FrontFilterModel
	cond := orm.NewCondition()
	cond1 := cond.And("date__gte", dateRange[0]).And("date__lte", dateRange[1]).Or("date__eq", "20180901")
	orm.NewOrm().QueryTable("front_filter_model").
		SetCond(cond1).
		All(&ff)

	tempMap := make(map[string]map[string]int)
	for _, obj := range ff {
		var frontMap map[string][]string
		json.Unmarshal([]byte(obj.PcJson), &frontMap)
		for key, strList := range frontMap {
			tmpSetMap := make(map[string]int)
			_, ok := tempMap[key]
			if !ok {
				tempMap[key] = tmpSetMap
			}
			for _, str := range strList {
				if str != "" {
					tempMap[key][str] = 1
				}
			}
		}
	}

	finalMap := make(map[string][]string)

	for key, strList := range tempMap {
		var tmpStrList []string
		for strKey, _ := range strList {
			tmpStrList = append(tmpStrList, strKey)
		}
		sort.Strings(tmpStrList)
		finalMap[key] = tmpStrList
	}

	data, _ := json.Marshal(finalMap)
	this.Ctx.WriteString(string(data))

}

// @Title 获取前端各部分的具体选项
// @Param	date body []string true "时间"
// @Summary 获取前端各部分的具体选项
// @Description 获取前端各部分的具体选项
// @Success 200 {object} models.FrontOptions
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/sql/options_yh/date [post]
func (this *DWController) GetOptionsYHDate() {
	var dateRange []string
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &dateRange)
	if err != nil {
		models.Logger.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		return
	}

	var ff []models.FrontFilterModel
	cond := orm.NewCondition()
	cond1 := cond.And("date__gte", dateRange[0]).And("date__lte", dateRange[1]).Or("date__eq", "20180901")
	orm.NewOrm().QueryTable("front_filter_model").
		SetCond(cond1).
		All(&ff)

	tempMap := make(map[string]map[string]int)
	for _, obj := range ff {
		var frontMap map[string][]string
		json.Unmarshal([]byte(obj.GalaxyJson), &frontMap)
		for key, strList := range frontMap {
			tmpSetMap := make(map[string]int)
			tempMap[key] = tmpSetMap
			for _, str := range strList {
				if str != "" {
					tempMap[key][str] = 1
				}
			}
		}
	}

	finalMap := make(map[string][]string)

	for key, strList := range tempMap {
		var tmpStrList []string
		for strKey, _ := range strList {
			tmpStrList = append(tmpStrList, strKey)
		}
		sort.Strings(tmpStrList)
		finalMap[key] = tmpStrList
	}

	data, _ := json.Marshal(finalMap)
	this.Ctx.WriteString(string(data))

}

// @Title 上传csv文件, 并解析成表 (多列)
// @Summary 上传csv文件, 并解析成表
// @Description 上传csv文件, 并解析成表
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/upload/file [post]
func (this *DWController) Upload() {
	// 接收文件
	uploadRoot := fmt.Sprintf("%s/%s/%s", models.ServerInfo.UploadRoot, this.userName,
		utils.DateFormatWithoutLine(time.Now()))
	utils.MkDir(uploadRoot)
	f, h, err := this.GetFile("file")
	fmt.Println(this.GetString("username", "")) // 表单中的其他参数

	now := time.Now()
	md5Ts := md5.Sum([]byte(fmt.Sprint(now)))
	toFileName := fmt.Sprintf("%s_%s_%x", this.userName, utils.DateFormatWithoutLine(now), md5Ts)

	toFile := fmt.Sprintf("%s/%s", uploadRoot, toFileName)
	fmt.Printf("upload %s to %s\n", h.Filename, toFile)
	defer f.Close()

	bytes, _ := utils.ReadBytes(f, h.Size)
	detector := chardet.NewTextDetector()
	result, _ := detector.DetectBest(bytes)

	var resStr string
	if result.Charset != "UTF-8" {
		resStr = utils.ConvertToString(string(bytes), result.Charset, "UTF-8")
	} else {
		resStr = string(bytes)
	}

	fileData := strings.Replace(resStr, "\r", "", -1)
	err = ioutil.WriteFile(toFile, []byte(fileData), 0777)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("write file successful")
	// 解析表头
	//fi, _ := os.Open(toFile)
	//br := bufio.NewReader(fi)
	//firstLine, _, _ := br.ReadLine()
	hdfsDir := fmt.Sprintf("/user/adtd_platform/queryword/part=%s", toFileName)
	utils.LoadQueryWords(hdfsDir, toFile)

	os.RemoveAll(toFile)
	fmt.Println("已删除文件:", toFile)
	this.Ctx.WriteString(fmt.Sprintf("part='%s'", toFileName))
}

// @Title 需要数据
// @Summary 需要数据
// @Description 需要数据
// @Success 200 {object} models.MetaTreeNode
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/explore [post]
func (this *DWController) GetExploreData() {
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)

	tableName := params["tableName"]
	dsId, _ := strconv.ParseInt(params["dsId"], 10, 64)
	dsType := params["dsType"]

	query := models.SqlQueryModel{}
	vxeData := query.GetExploreData(tableName, 10, dsType, dsId)
	res, _ := json.Marshal(vxeData)
	this.Ctx.WriteString(string(res))
}

// @Summary 索引字段
// @Description 索引字段
// @Success 200 {string}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dw/searchField [post]
func (this *DWController) SearchField() {
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)

	var tables []SearchFieldTable
	fieldName := params["fieldName"]
	roles := models.Enforcer.GetRolesForUser(this.userName)
	roleWildCard := utils.MakeRolePlaceHolder(roles)

	o := orm.NewOrm()
	sql := fmt.Sprintf(`
SELECT t2.tid,t2.table_name, t2.table_type, t2.dsid, t2.database_name
FROM 
(SELECT p1.table_id FROM
  (SELECT table_id FROM field_meta_data WHERE field_name LIKE ? GROUP BY table_id) p1
  LEFT JOIN 
  (select table_id from group_metatable_mapping where group_name in %s AND (read_flag =1 or exec_flag=1) AND valid=1 group by table_id) p2
  ON p1.table_id = p2.table_id 
  where p2.table_id is not null
) t1
LEFT JOIN
table_meta_data t2
ON t1.table_id=t2.tid
WHERE t2.tid IS NOT null
`, roleWildCard)
	_, _ = o.Raw(sql, fmt.Sprintf("%%%s%%", fieldName), roles).QueryRows(&tables)

	var hiveTables []SearchFieldTableTreeNode
	var ckTables []SearchFieldTableTreeNode
	var mysqlTables []SearchFieldTableTreeNode

	for _, t := range tables {
		node := SearchFieldTableTreeNode{
			Label:     t.TableName,
			TableName: fmt.Sprintf("%s.%s", t.DatabaseName, t.TableName),
			TableId:   t.Tid,
			TableType: t.TableType,
			DsId:      t.Dsid,
		}
		if t.TableType == "hive" {
			hiveTables = append(hiveTables, node)
		} else if t.TableType == "mysql" {
			mysqlTables = append(mysqlTables, node)
		} else if t.TableType == "clickhouse" {
			ckTables = append(ckTables, node)
		}
	}

	var children []SearchFieldTableTreeNode
	if len(hiveTables) > 0 {
		children = append(children, SearchFieldTableTreeNode{
			Label:    "hive",
			Children: hiveTables,
		})
	}
	if len(ckTables) > 0 {
		children = append(children, SearchFieldTableTreeNode{
			Label:    "clickhouse",
			Children: ckTables,
		})
	}
	if len(mysqlTables) > 0 {
		children = append(children, SearchFieldTableTreeNode{
			Label:    "mysql",
			Children: mysqlTables,
		})
	}

	metaTree := SearchFieldTableTreeNode{
		Label:    "meta",
		Children: children,
	}

	res, _ := json.Marshal(metaTree)
	this.Ctx.WriteString(string(res))
}

type SearchFieldTable struct {
	Tid          int64
	TableName    string
	DatabaseName string
	TableType    string
	Dsid         int64
}

type SearchFieldTableTreeNode struct {
	Label     string                     `json:"label"`
	TableName string                     `json:"tableName"`
	TableType string                     `json:"tableType"`
	TableId   int64                      `json:"tid"`
	DsId      int64                      `json:"dsid"`
	Children  []SearchFieldTableTreeNode `json:"children,omitempty"`
}
