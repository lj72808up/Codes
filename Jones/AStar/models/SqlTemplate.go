package models

import (
	pb "AStar/protobuf"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"github.com/golang/protobuf/proto"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SqlTemplate struct {
	TemplateId    int64  `orm:"pk;auto"`
	Title         string `orm:"size(50)"`
	Creator       string `orm:"size(50)"`
	Template      string `orm:"type(text)"`
	Param         string
	Dsid          int64
	Pid           int64  `orm:"default(0)"`
	Description   string `orm:"size(255)"`
	Createredash  bool   `orm:"-"`
}

type SqlTemplateSeconLevel struct {
	TemplateId    int64  `orm:"pk;auto"`
	Title         string `orm:"size(50)"`
	Creator       string `orm:"size(50)"`
	Template      string `orm:"type(text)"`
	Param         string
	Dsid          int64
	Pid           int64  `orm:"default(0)"`
	Description   string `orm:"size(255)"`
	Createredash  bool   `orm:"-"`
	QueryId       int64
}

type TemplateTreeNode struct {
	Title         string             `json:"label"`
	TemplateId    string             `json:"templateId"`
	Children      []TemplateTreeNode `json:"children,omitempty"`
	Pid           string             `json:"pid"`
	RedashQueryId string             `json:"redashid"`
}

type SqlTemplateHeat struct {
	Hid        int64  `orm:"pk;auto"`
	TemplateId int64  `orm:"size(100)"`
	QueryId    string `orm:"size(32)"`
	Date       string `orm:"size(11)"`
	User       string `orm:"size(100)"`
}

type LimitedAccessElement struct {
	Eid     int64  `orm:"pk;auto"`
	Subject string `orm:"size(50)"`
	Url     string `orm:"size(50)"`
	Element string `orm:"size(50)"`
}

type GroupWorksheetMapping struct {
	Gwid        int64  `orm:"pk;auto"`
	GroupName   string `orm:"size(255)"`
	WorksheetId string `orm:"size(50)"`
	CanModify   bool   `orm:"default(false)"`
	Valid       bool   `orm:"default(true)"`
}

// 多字段索引
func (u *GroupWorksheetMapping) TableIndex() [][]string {
	return [][]string{
		{"Valid", "GroupName"}, {"CanModify"},
	}
}


func (u *GroupWorksheetMapping) TableUnique() [][]string {
	return [][]string{
		[]string{"GroupName", "WorksheetId"},
	}
}

func (this *LimitedAccessElement) GetLimitedElement() int64 {
	o := orm.NewOrm()
	err := o.Raw("SELECT eid FROM limited_access_element WHERE Subject=? and url=?", this.Subject, this.Url).QueryRow(this)
	if err != nil {
		fmt.Println(err.Error())
	}
	return this.Eid
}

//返回成功插入的 queryid，失败返回0
// -1 连接未创建
// -2
func (this SqlTemplate) InsertRedashQuery() (int, string) {
	err_msg := ""
	redash_queryid := 0

	//正则提取参数
	template := this.Template
	reg1 := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	submatchall := reg1.FindAllString(template, -1)

	var paramList []string
	var paramRangeList map[string]string
	paramRangeList = make(map[string]string)

	for _, element := range submatchall {
		element = strings.Trim(element, "{{")
		element = strings.Trim(element, "}}")
		fmt.Println(element)

		if strings.Contains(element, ".start") || strings.Contains(element, ".end") {
			keySplits := strings.Split(element, ".")
			paramRangeList[keySplits[0]] = keySplits[0]
		} else {
			paramList = append(paramList, element)
		}
	}

	//获取Redash datasoource id
	o := orm.NewOrm()
	sql := fmt.Sprintf("SELECT redashid from connection_manager where dsid=%d ;", this.Dsid)
	redash_dsid := 0
	o.Raw(sql).QueryRow(&redash_dsid)
	if redash_dsid == 0 {
		err_msg = "数据库连接Redash端未创建"
		return redash_queryid, err_msg
	}

	var request_map map[string]interface{}
	request_map = make(map[string]interface{})

	request_map["query"] = this.Template
	request_map["data_source_id"] = redash_dsid
	request_map["name"] = this.Title
	request_map["tags"] = []string{this.Creator}

	if len(paramList) != 0 || len(paramRangeList) != 0 {
		thisMap := make(map[string][]map[string]string)
		for _, param := range paramList {
			cur_param := make(map[string]string)
			cur_param["name"] = param
			cur_param["title"] = param
			cur_param["value"] = " "
			cur_param["type"] = "number"
			cur_param["locals"] = "[]"
			thisMap["parameters"] = append(thisMap["parameters"], cur_param)
		}

		//添加date range 参数
		for _, param := range paramRangeList {
			cur_param := make(map[string]string)
			cur_param["name"] = param
			cur_param["title"] = param
			cur_param["value"] = " "
			cur_param["type"] = "date-range"
			cur_param["locals"] = "[]"
			thisMap["parameters"] = append(thisMap["parameters"], cur_param)
		}

		request_map["options"] = thisMap
	}
	byte_param, _ := json.Marshal(request_map)
	param := string(byte_param)

	//通过redash 接口插入
	url := "http://inner-redash.adtech.sogou/api/queries"
	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body([]byte(param))

	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
		return redash_queryid, err_msg
	}
	redash_queryid, err := strconv.Atoi(string(bytes))
	if err != nil {
		Logger.Error(err.Error())
		redash_queryid = 0
		err_msg = "RedashQueryAPI创建失败"
		return redash_queryid, err_msg
	}

	//插入本地RedashQuery表中
	var rdb FrontRedashQuery
	rdb.Name = this.Title
	rdb.Creator = this.Creator
	rdb.QueryId = redash_queryid
	rdb.Index = fmt.Sprintf("queryshow/%d", redash_queryid)
	rdb.Url = fmt.Sprintf("http://inner-redash.adtech.sogou/queries/%d/source", redash_queryid)
	rdb.TemplateId = this.TemplateId
	o.Insert(&rdb)

	err_msg = ""
	return redash_queryid, err_msg
}

// 创建绿皮 Redash 页面
//返回成功插入的 queryid，失败返回0
// -1 连接未创建
// -2
func (this SqlTemplate) InsertLvpiRedashQuery(encode string, date_pianyi bool) (int, string) {
	err_msg := ""
	redash_queryid := 0

	//正则提取参数
	template := this.Template
	reg1 := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	submatchall := reg1.FindAllString(template, -1)

	var paramList []string

	for _, element := range submatchall {
		element = strings.Trim(element, "{{")
		element = strings.Trim(element, "}}")
		fmt.Println(element)
		paramList = append(paramList, element)
	}

	//获取Redash datasoource id
	o := orm.NewOrm()
	sql := fmt.Sprintf("SELECT redashid from connection_manager where dsid=%d ;", this.Dsid)
	redash_dsid := 0
	o.Raw(sql).QueryRow(&redash_dsid)
	if redash_dsid == 0 {
		err_msg = "数据库连接Redash端未创建"
		return redash_queryid, err_msg
	}

	var request_map map[string]interface{}
	request_map = make(map[string]interface{})

	request_map["query"] = this.Template
	request_map["data_source_id"] = redash_dsid
	request_map["name"] = this.Title
	request_map["tags"] = []string{this.Creator}
	request_map["lvpi_flag"] = true

	if len(paramList) != 0 {
		thisMap := make(map[string][]map[string]string)
		for _, param := range paramList {
			cur_param := make(map[string]string)
			cur_param["name"] = param
			cur_param["title"] = param
			cur_param["value"] = "d_now"
			cur_param["type"] = "date"
			cur_param["locals"] = "[]"


			if date_pianyi {
				if param == "基准日期"{
					cur_param["value"] = "d_bias_yesterday"
				} else {
					cur_param["value"] = "d_bias_before_yesterday"
				}
				cur_param["type"] = "date-lvpi"
			} else {
				if param == "基准日期"{
					cur_param["value"] = "d_yesterday"
				} else {
					cur_param["value"] = "d_day_before_yesterday"
				}
				cur_param["type"] = "date"
			}
			thisMap["parameters"] = append(thisMap["parameters"], cur_param)
		}
		lvpi_params := make(map[string]string)
		lvpi_params["encode"] = encode
		if date_pianyi {
			lvpi_params["bias"] = "true"
		} else {
			lvpi_params["bias"] = "false"
		}
		thisMap["lvpi"] = append(thisMap["lvpi"], lvpi_params)
		request_map["options"] = thisMap
	}
	byte_param, _ := json.Marshal(request_map)
	param := string(byte_param)

	//通过redash 接口插入
	url := "http://inner-redash.adtech.sogou/api/queries"
	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body([]byte(param))

	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
		return redash_queryid, err_msg
	}
	redash_queryid, err := strconv.Atoi(string(bytes))
	if err != nil {
		Logger.Error(err.Error())
		redash_queryid = 0
		err_msg = "RedashQueryAPI创建失败"
		return redash_queryid, err_msg
	}

	//插入本地RedashLvpiQuery表中
	var rdb FrontRedashQuery
	rdb.Name = this.Title
	rdb.Creator = this.Creator
	rdb.QueryId = redash_queryid
	rdb.Index = fmt.Sprintf("queryshow/%d", redash_queryid)
	rdb.Url = fmt.Sprintf("http://inner-redash.adtech.sogou/queries/%d/source", redash_queryid)
	rdb.TemplateId = this.TemplateId
	rdb.IsLvpi = 1
	o.Insert(&rdb)

	err_msg = ""
	return redash_queryid, err_msg
}


//更新Redash里的  queryid；成功返回redashid，失败返回0
func (this SqlTemplate) UpdateRedash(title string, template string) int {
	//获取Redash datasoource id
	o := orm.NewOrm()
	sql1 := fmt.Sprintf("SELECT dsid from sql_template where template_id=%d ;", this.TemplateId)
	o.Raw(sql1).QueryRow(&this.Dsid)

	sql := fmt.Sprintf("SELECT redashid from connection_manager where dsid=%d ;", this.Dsid)
	redash_dsid := 0
	o.Raw(sql).QueryRow(&redash_dsid)

	var request_map map[string]string
	request_map = make(map[string]string)

	request_map["query"] = template
	request_map["data_source_id"] = strconv.Itoa(redash_dsid)
	request_map["name"] = title

	byte_param, _ := json.Marshal(request_map)
	param := string(byte_param)

	//通过redash 接口插入
	//临时 停用 函数
	//////////////////////////////url := fmt.Sprintf("http://inner-redash.adtech.sogou/api/queries/%d", this.RedashQueryId)
	url := ""
	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body([]byte(param))

	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
	}
	redash_queryid, err := strconv.Atoi(string(bytes))
	if err != nil {
		Logger.Error(error.Error())
		redash_queryid = 0
	}
	return redash_queryid
}

func (this SqlTemplate) InsertQueryInfoTopLabel(redash_query_id int) {
	//插入TopLabel
	pid_init := 40
	urlSql := fmt.Sprintf(" SELECT lid from top_label where title='%s'", "RedashQuery")
	orm.NewOrm().Raw(urlSql).QueryRow(&pid_init)

	info := make(map[int]string)
	info[redash_query_id] = fmt.Sprintf("http://inner-redash.adtech.sogou/queries/%d/source", redash_query_id)
	info_str, _ := json.Marshal(info)

	var rd TopLabel
	rd.Title = this.Title
	rd.Index = fmt.Sprintf("queryshow/%d", redash_query_id)
	rd.Info = string(info_str)
	rd.Pid = int64(pid_init)
	rd.NameSpace = "redash"

	//newLid , err := orm.NewOrm().Insert(&rd)
	_, err := orm.NewOrm().Insert(&rd)
	if err != nil {
		Logger.Error(err.Error())
	}
}

func (this *SqlTemplate) ChangeAffiliation (templateId int64, folderId int64) error {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE sql_template SET pid=? WHERE template_id = ?", folderId, templateId).Exec()
	if err != nil {
		utils.Logger.Error("%s",err.Error())
	}
	return err
}

func (this SqlTemplate) NewTemplate() (int64, string) {
	templateId, err := orm.NewOrm().Insert(&this)
	if err != nil {
		Logger.Error(err.Error())
	}
	this.TemplateId = templateId

	//插入redash
	if this.Createredash == true {
		res_msg := this.CreateQueryDatasource()
		if len(res_msg) != 0 {
			return 0, res_msg
		}

		_, err := this.InsertRedashQuery()
		if err != "" {
			return 0, err
		}
		//if redash_query_id != 0 {
		//	this.InsertQueryInfoTopLabel(redash_query_id)
		//}
	}

	return templateId, ""
}

func (this *SqlTemplate) AddTemplateHeat(queryId string, user string) {
	o := orm.NewOrm()
	now := utils.DateFormat(time.Now())
	templateHeat := SqlTemplateHeat{
		TemplateId: this.TemplateId,
		QueryId:    queryId,
		Date:       now,
		User:       user,
	}
	_, _ = o.Insert(&templateHeat)
}

func (this *SqlTemplate) GetSqlTemplate(templateId int64) {
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM sql_template WHERE template_id=?", templateId).QueryRow(this)
	if err == orm.ErrNoRows {
		Logger.Error("查询不到")
	} else if err == orm.ErrMissPK {
		Logger.Error("找不到主键")
	} else {
		Logger.Info("title:%s \nTemplate:%s", this.Title, this.Template)
	}
}

func (this *SqlTemplate) DeleteById(id int64) error {
	o := orm.NewOrm()
	_, err := o.Raw("DELETE FROM sql_template WHERE template_id = ?", id).Exec()
	return err
}

func (this *SqlTemplate) CanModifyById(templateId int64, roles []string) bool {
	o := orm.NewOrm()
	modifyRoles := make([]string, 0)
	rolesWildCard := utils.MakeRolePlaceHolder(roles)
	sql := fmt.Sprintf("SELECT group_name FROM group_worksheet_mapping "+
		"WHERE valid=true and worksheet_id=? and group_name in %s AND can_modify=TRUE group by group_name", rolesWildCard)
	_, _ = o.Raw(sql, templateId, roles).QueryRows(&modifyRoles)
	Logger.Info("可以编辑%d的角色数为:%v", templateId, modifyRoles)
	canSaveModify := len(modifyRoles) > 0
	return canSaveModify
}

// 限制只有2个层级的树节点
func (this *SqlTemplate) GetAllTemplate(roles []string, connType string,user string) []TemplateTreeNode {
	o := orm.NewOrm()

	var firstLevel []SqlTemplate
	var secondLevel []SqlTemplateSeconLevel

	rolesWildCard := utils.MakeRolePlaceHolder(roles)

	firstSql := fmt.Sprintf(`select t1.* from 
					(SELECT * FROM sql_template where pid=0 ) t1 
				left join
					(SELECT worksheet_id FROM group_worksheet_mapping WHERE valid = true and group_name in %s 
					 group by worksheet_id) t2
				on t1.template_ID = t2.worksheet_id
				where t2.worksheet_id is not null`, rolesWildCard)

	if connType != "all" {
		firstSql = fmt.Sprintf(`
				SELECT
					t1.*
				FROM
					(
						SELECT * FROM sql_template WHERE pid = 0
					) t1
				JOIN connection_manager t2 ON t1.dsid = t2.dsid
				LEFT JOIN (
					SELECT worksheet_id FROM group_worksheet_mapping
					WHERE valid = TRUE AND group_name in %s group by worksheet_id
				) t3 ON t1.template_ID = t3.worksheet_id
				WHERE
					t3.worksheet_id IS NOT NULL
					and t2.conn_type= ?`, rolesWildCard)
	}

	if connType != "all" {
		_, _ = o.Raw(firstSql, roles, connType).QueryRows(&firstLevel)
	} else {
		// 取mysql hive clickHouse 所有的sql模板
		_, _ = o.Raw(firstSql, roles).QueryRows(&firstLevel)
	}

	secondSql := fmt.Sprintf(`SELECT
										st.*,IFNULL(frq.query_id,0) AS query_id
									FROM
										sql_template st
									LEFT JOIN 
										( SELECT * FROM front_redash_query WHERE creator = '%s') frq 
									ON st.template_id = frq.template_id
									WHERE
										st.pid != 0`, user)
	_, _ = o.Raw(secondSql).QueryRows(&secondLevel)

	// 转成view
	nodes := make([]TemplateTreeNode, len(firstLevel))
	for index, x := range firstLevel {
		nodes[index] = TemplateTreeNode{
			Title:      x.Title,
			TemplateId: fmt.Sprint(x.TemplateId),
			Children:   nil,
			Pid:        fmt.Sprint(x.Pid),
		}
	}

	if len(secondLevel) != 0 {
		for index, _ := range nodes {
			var children []TemplateTreeNode
			for _, y := range secondLevel {
				if fmt.Sprint(y.Pid) == nodes[index].TemplateId {
					children = append(children, TemplateTreeNode{
						Title:         y.Title,
						TemplateId:    fmt.Sprint(y.TemplateId),
						Children:      nil,
						Pid:           fmt.Sprint(y.Pid),
						RedashQueryId: fmt.Sprint(y.QueryId),
					})
				}
			}
			nodes[index].Children = children
		}
	}
	//fmt.Println(nodes)
	return nodes
}

func (this *SqlTemplate) GetChildren(pid int64) []TemplateTreeNode {
	o := orm.NewOrm()
	sql := fmt.Sprint("SELECT template_id, title FROM sql_template where pid=", pid)
	var children []TemplateTreeNode
	_, err := o.Raw(sql).QueryRows(&children)
	if err != nil {
		panic(err.Error())
	}
	return children
}

func (this SqlTemplate) GenerateSql(sql string, param []map[string]string, paddingType string) string {
	runSql := sql
	// 替换模板中的参数为{{|safe}} 模式, 避免 单引号 转义
	sql = strings.Replace(sql, "}}", "|safe}}", -1)
	param, chineseKeyMapping := handleChineseParam(param)
	for k, v := range chineseKeyMapping {
		sql = strings.Replace(sql, k, v, -1)
	}
	tmpl := utils.SqlUtil{}
	p := make(map[string]interface{})

	if param != nil && len(param) != 0 {
		p = this.TransParam(param)
	}

	if paddingType == "dag" {
		p = tmpl.PaddingDagInternalVarFromUserParam(p)
	}

	// 最终替换
	if p != nil && len(p) > 0 {
		runSql = tmpl.SubstituteJinja(sql, p)
	}

	Logger.Info("runSql:%s\n", runSql)
	return runSql
}

func handleChineseParam(param []map[string]string) ([]map[string]string, map[string]string) {
	cnt := 0 // 中文个数
	prefix := "chineseVar"
	transParam := make([]map[string]string, 0)
	chineseKeyMapping := make(map[string]string)
	for _, m := range param {
		paramName := m["name"]
		if utils.IsChinese(paramName) {
			newParamName := fmt.Sprintf("%s%d", prefix, cnt)
			cnt = cnt + 1
			chineseKeyMapping[paramName] = newParamName
			item := make(map[string]string) // 复制mapping中的项
			for k, v := range m {
				if k == "name" {
					item["name"] = newParamName
				} else {
					item[k] = v
				}
			}
			transParam = append(transParam, item)
		} else {
			item := make(map[string]string) // 复制mapping中的项
			for k, v := range m {
				item[k] = v
			}
			transParam = append(transParam, item)
		}
	}
	return transParam, chineseKeyMapping
}

type TimeInterval struct {
	start string `json:"start"`
	end   string `json:"end"`
}

func (this SqlTemplate) TransParam(param []map[string]string) map[string]interface{} {
	p := make(map[string]interface{})
	intervalparams := make(map[string]*TimeInterval)

	for _, m := range param {
		key := strings.Trim(m["name"], " ")
		if !strings.Contains(key, ".") {
			value := m["value"]
			p[key] = value
		}
		// 对于间隔型变量
		if m["type"] == "timeInterval" {
			keySplits := strings.Split(key, ".")
			keyPrefix := strings.Trim(keySplits[0], " ") // 前缀
			keySuffix := strings.Trim(keySplits[1], " ") // 后缀
			// 判断 prefix 是否已经加入到 intervalParams
			isExist := false
			for kName, v := range intervalparams {
				// 已经存在, 则修改 prefix 和 suffix
				if kName == keyPrefix {
					isExist = true
					if keySuffix == "start" {
						v.start = m["value"]
					}
					if keySuffix == "end" {
						v.end = m["value"]
					}
				}
			}
			if !isExist {
				// 不存在, 建立一个变量
				if keySuffix == "start" {
					newItem := TimeInterval{
						start: m["value"],
						end:   "",
					}
					intervalparams[keyPrefix] = &newItem
				}
				if keySuffix == "end" {
					newItem := TimeInterval{
						start: "",
						end:   m["value"],
					}
					intervalparams[keyPrefix] = &newItem
				}
			}
		}
	}
	// 将 intervalParam 加入到 p
	for k, v := range intervalparams {
		p[k] = v
	}
	return p
}

func (this SqlTemplate) GetAllDataSource(connType string) []ConnectionManager {
	return GetAllConn(connType)
}

func (this SqlTemplate) UpdateTemplate(title string, template string, desc string,opt_user string) {
	o := orm.NewOrm()

	////通过template id 获取 dsid 和 查看是有已经有 存在的 redash_query_id
	//sql := fmt.Sprintf("SELECT redash_query_id,dsid,creator from sql_template where template_id=%d ;", this.TemplateId)
	//o.Raw(sql).QueryRow(&this.RedashQueryId, &this.Dsid, &this.Creator)
	//
	////if this.RedashQueryId == 0 {
	////	//更新redash 里的query信息
	////	this.Title = title
	////	this.Template = template
	////	this.Description = desc
	////
	////	redash_query_id := this.InsertRedash()
	////	this.RedashQueryId = redash_query_id
	////} else {
	////	this.UpdateRedash(title, template)
	////}
	//if this.RedashQueryId != 0 {
	//	this.UpdateRedash(title, template)
	//}

	if o.Read(&this) == nil {
		this.Title = title
		this.Template = template
		this.Description = desc
		num, err := o.Update(&this)
		if err == nil {
			Logger.Info("update success, update rows: ", num)
		} else {
			Logger.Error("update failure")
		}
	}
}

func (this SqlTemplate) CreateQueryDatasource() string {
	o := orm.NewOrm()
	res_msg := ""
	//如果依赖的数据连接不存在， 则创建连接
	redash_dsid := 0
	ds_type := ""
	sql := fmt.Sprintf("SELECT redashid,conn_type from connection_manager where dsid=%d ;", this.Dsid)
	o.Raw(sql).QueryRow(&redash_dsid, &ds_type)

	if len(ds_type) == 0 {
		res_msg = "当前Query所依赖DataSource数据库不存在"
		return res_msg
	}

	if ds_type != "mysql" && ds_type != "clickhouse" {
		res_msg = "当前引擎不支持创建可视化"
		return res_msg
	}

	if redash_dsid == 0 {
		//创建所依赖的 Redash Data Source
		var conn ConnectionManager
		sql = fmt.Sprintf("SELECT * from connection_manager where dsid=%d ;", this.Dsid)
		o.Raw(sql).QueryRow(&conn)

		suc_flag, redash_dsid := conn.AddRedashConn()
		if !suc_flag {
			res_msg = "Create Redash DataSource Failed"
			return res_msg
		}
		conn.Redashid = redash_dsid
		o.Update(&conn, "redashid")
	}
	return res_msg
}

func (this SqlTemplate) CreateQueryViewModel(newCreator string) (string, int) {
	o := orm.NewOrm()
	res_msg := ""

	//通过template id 获取 dsid 和 查看是有已经有 存在的 redash_query_id
	sql := fmt.Sprintf("SELECT * from sql_template where template_id=%d ;", this.TemplateId)
	o.Raw(sql).QueryRow(&this)
	this.Creator = newCreator

	//if this.RedashQueryId != 0 {
	//	res_msg = "可视化Query已存在!"
	//	return res_msg, 0
	//}
	if this.Dsid == 0 {
		res_msg = "当前Query所依赖DataSourceId为0!"
		return res_msg, 0
	}

	res_msg = this.CreateQueryDatasource()
	if len(res_msg) != 0 {
		return res_msg, 0
	}

	redash_query_id, err := this.InsertRedashQuery()
	//if redash_query_id != 0 && len(err) == 0 {
	//	this.RedashQueryId = redash_query_id
	//	o.Update(&this, "redash_query_id")
	//	//this.InsertQueryInfoTopLabel(redash_query_id)
	//}
	return err, redash_query_id
}


func (this SqlTemplate) SQLCreateQueryViewModel() (string, int) {
	res_msg := ""

	if this.Dsid == 0 {
		res_msg = "当前Query所依赖DataSourceId为0!"
		return res_msg, 0
	}

	res_msg = this.CreateQueryDatasource()
	if len(res_msg) != 0 {
		return res_msg, 0
	}

	redash_query_id, err := this.InsertRedashQuery()
	return err, redash_query_id
}

func (this SqlTemplate) LvpiCreateQueryViewModel(encode string, date_pianyi bool) (string, int) {
	res_msg := ""

	if this.Dsid == 0 {
		res_msg = "当前Query所依赖DataSourceId为0!"
		return res_msg, 0
	}

	res_msg = this.CreateQueryDatasource()
	if len(res_msg) != 0 {
		return res_msg, 0
	}

	redash_query_id, err := this.InsertLvpiRedashQuery(encode, date_pianyi)
	return err, redash_query_id
}

func ParseHeader(sql string) []string {
	sql = strings.ToLower(sql)
	start := strings.Index(sql, "select")
	if start != -1 {
		start = start + 6
	} else {
		start = 0
	}
	end := strings.Index(sql, "from")
	if end == -1 {
		end = len(sql)
	}
	headerStr := strings.Trim(sql[start:end], " ")
	split := strings.Split(headerStr, ",")
	var res []string
	for _, v := range split {
		res = append(res, strings.Trim(v, " "))
	}
	return res
}

func (this SqlTemplate) PublishMq(sql string, dsid int64, username string, runSql string, param string,
	queryWordsPartition string, title string) string {
	session := username
	if strings.Contains(title, "_outputcsv") {
		session = "csv"
	}

	// 不同于query 放入MQ,返回queryId
	nowTime := time.Now()
	// 生成Md5结果，作为QueryInfo的QueryId
	Md5Res := utils.GetMd5ForQuery(fmt.Sprint(this.TemplateId), this.Creator, nowTime.String())
	//dataSource := this.GetDataSource()
	connManager := ConnectionManager{}
	connManager.GetById(this.Dsid)
	// 指定查询引擎
	var queryEngine pb.QueryEngine
	var queryType string

	switch connManager.ConnType {
	case "mysql":
		queryEngine = pb.QueryEngine_jdbc
		queryType = "工单"
	case "clickhouse":
		queryEngine = pb.QueryEngine_jdbc_clickhouse
		queryType = "工单"
	case "hive":
		queryEngine = pb.QueryEngine_jdbc_spark
		queryType = "工单"
	case "kylin":
		queryEngine = pb.QueryEngine_kylin
		queryType = "工单"
	}

	filterObj := make(map[string]interface{})
	filterObj["name"] = "queryWordPartition"
	filterObj["value"] = queryWordsPartition
	filterObj["type"] = pb.FilterOperation_name[10]
	var filterList []map[string]interface{}
	filterList = append(filterList, filterObj)
	filterObjJson, _ := json.Marshal(filterObj)

	// mysql中的QueryInfo
	queryInfo := pb.QueryInfo{
		QueryId:    Md5Res,                          // query id
		Sql:        []string{sql},                   // sql statement
		Session:    username,                        // session :this.userName
		Engine:     queryEngine,                     // 引擎
		Filters:    []string{string(filterObjJson)}, // 过滤字段
		Status:     pb.QueryStatus_Waiting,          // 状态
		LastUpdate: nowTime.UnixNano() / 1e6,        // 最后更新时间戳
		CreateTime: nowTime.UnixNano() / 1e6,        // 创建时间戳
		Tables:     []string{""},                    // 查询所用表
		Type:       queryType,                       // 查询类别"小流量下载"
		Dsid:       dsid,
	}
	var queryTaskName = this.Title
	jsonMap := make(map[string]interface{}) //map[string]string{}
	jsonMap["dsid"] = fmt.Sprint(dsid)
	jsonMap["param"] = param
	jsonMap["filters"] = filterList

	queryJSONStr, _ := json.Marshal(jsonMap)

	//mq中的sql 拷贝一个queryinfo, 只是sql是truesql
	mqQueryInfo := pb.QueryInfo{
		QueryId:    Md5Res,                          // query id
		Sql:        []string{runSql},                // sql statement
		Session:    session,                         // session :this.userName
		Engine:     queryEngine,                     // 引擎
		Filters:    []string{string(filterObjJson)}, // 过滤字段
		Status:     pb.QueryStatus_Waiting,          // 状态
		LastUpdate: nowTime.UnixNano() / 1e6,        // 最后更新时间戳
		CreateTime: nowTime.UnixNano() / 1e6,        // 创建时间戳
		Tables:     []string{""},                    // 查询所用表
		Type:       queryType,                       // 查询类别"小流量下载"
		Dsid:       dsid,
	}
	// 编码QueryInfo
	queryBytes, _ := proto.Marshal(&mqQueryInfo)
	// 先推到MQ
	RabbitMQModel{}.Publish(queryBytes)
	// 再将QueryInfo入库MySql
	QuerySubmitLog{}.InsertQuerySubmitLog(queryInfo, queryTaskName, string(queryJSONStr))

	// 返回queryId
	return queryInfo.QueryId
}

// 创建工单时增加权限
func (this GroupWorksheetMapping) AddOneMapping() {
	o := orm.NewOrm()
	_, err := o.Insert(&this)
	if err != nil {
		panic(err)
	}
}

//更新角色和工单列表
func (this SqlTemplate) AddworkSheetMapping(groupName string, worksheets []SqlTemplateView) {
	o := orm.NewOrm()
	res, err := o.Raw("DELETE FROM group_worksheet_mapping WHERE group_name = ?", groupName).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		Logger.Info("group worksheet mapping affected nums: ", num)
	} else {
		panic(err)
	}

	mappings := make([]GroupWorksheetMapping, 0)
	for _, sheet := range worksheets {
		mapping := GroupWorksheetMapping{
			GroupName:   groupName,
			WorksheetId: fmt.Sprint(sheet.TemplateId),
			Valid:       true,
			CanModify:   sheet.CanModify,
		}
		mappings = append(mappings, mapping)
	}
	// bulk: 一次性并列插入的条数
	if len(mappings) > 0 {
		_, err2 := o.InsertMulti(100, mappings)
		if err2 != nil {
			panic(err2)
		}
	}
}

type SqlTemplateView struct {
	TemplateId   int64  `json:"worksheetId"`
	Title        string `json:"worksheetName"`
	HasPrivilege bool   `json:"hasPrivilege"`
	CanModify    bool   `json:"canModify"`
}

// 获取所有工单
func (this SqlTemplate) GetAllSqlTemplateCategory() []SqlTemplate {
	o := orm.NewOrm()
	sql := fmt.Sprint("SELECT * FROM sql_template WHERE pid = 0")
	var templates []SqlTemplate
	_, err := o.Raw(sql).QueryRows(&templates)
	if err != nil {
		panic(err.Error())
	}
	return templates
}

// 查询用户下有权限的工单
func (this SqlTemplate) GetCheckedCategory(groupName string) []SqlTemplateView {
	o := orm.NewOrm()
	var views []SqlTemplateView
	sql := `select 
    t1.template_id, t1.title, 
    case 
        when t2.worksheet_id is null THEN false
        else true 
    END as has_privilege,
    CASE 
    	WHEN t2.can_modify IS NOT NULL THEN t2.can_modify
    	ELSE FALSE
    END AS can_modify
FROM
(SELECT template_id, title FROM sql_template WHERE pid = 0) t1 
left join
(SELECT worksheet_id,can_modify FROM group_worksheet_mapping WHERE valid = true and group_name = ?) t2
on
 t1.template_id = t2.worksheet_id
order by t1.template_id asc`
	_, err := o.Raw(sql, groupName).QueryRows(&views)
	if err != nil {
		panic(err.Error())
	}
	return views
}

func (this SqlTemplate) GetUserViewAccess(roles []string) bool {
	o := orm.NewOrm()

	queryOptionLableId := 0
	o.Raw("select lid from top_label where `index`='redashqueryoption' ").QueryRow(&queryOptionLableId)
	if queryOptionLableId == 0 {
		return false
	}

	viewRoles := make([]string, 0)
	rolesWildCard := fmt.Sprintf("('%s')", strings.Join(roles, "','"))
	sql := fmt.Sprintf("SELECT group_name FROM group_label_mapping "+
		"WHERE valid=true and label_id=? and group_name in %s group by group_name", rolesWildCard)
	_, _ = o.Raw(sql, queryOptionLableId).QueryRows(&viewRoles)
	Logger.Info("当前用户有权限查看可视化的用户组为:%v", viewRoles)
	access := len(viewRoles) > 0
	return access
}


func (this SqlTemplate) GetUserJumpAccess(user string) bool {
	o := orm.NewOrm()
	redash_creator := ""

	sql := fmt.Sprintf("select rq.creator from sql_template tp join front_redash_query rq " +
		"where tp.template_id=rq.template_id and tp.template_id='%d' and rq.creator='%s' ", this.TemplateId,user)
	_= o.Raw(sql).QueryRow(&redash_creator)

	access := (redash_creator == user)
	return access
}
