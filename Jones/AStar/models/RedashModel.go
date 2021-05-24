package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type FrontRedashAllTableRes struct {
	Lid   int64  `json:"lid"`
	Pid   int64  `json:"pid"`
	Title string `json:"title"    orm:"column(title)"`
	Index string `json:"index" 	orm:"column(index)"`
	Info  string `json:"info" 	orm:"column(info)"`
}
type FrontRedashDashboard struct {
	Id          int    `orm:"pk;auto"`
	Creator     string `orm:"size(200)"`
	Name        string `orm:"size(200)"`
	Url         string `json:"slug" 		orm:"size(200)"`
	DashboradId int    `json:"id"  		orm:"size(50)"`
	Index       string `orm:"size(200)`
}

type FrontRedashQuery struct {
	Id         int    `orm:"pk;auto"`
	Creator    string `orm:"size(200)"`
	Name       string `orm:"size(200)"`
	Url        string `json:"slug" 		orm:"size(200)"`
	QueryId    int    `json:"id"  		orm:"size(20)"`
	Index      string `orm:"size(200)`
	TemplateId int64  `orm:"default(0)"`
	IsLvpi     int    `orm:"default(0)"`
}

type FrontRedashLvpiQuery struct {
	Id         int    `orm:"pk;auto"`
	Creator    string `orm:"size(200)"`
	Name       string `orm:"size(200)"`
	Url        string `json:"slug" 		orm:"size(200)"`
	QueryId    int    `json:"id"  		orm:"size(20)"`
	Index      string `orm:"size(200)`
	TemplateId int64  `orm:"default(0)"`
}

type FrontViewInfo struct {
	Name    string `json:"name"    orm:"column(name)"`
	Url     string `json:"url" 		orm:"column(url)"`
	Index   string `json:"index" 	orm:"column(index)"`
	Creator string `json:"creator" 	orm:"column(creator)"`
}

func InsertRedash(attr map[string]string) int64 {
	var newLid int64
	newLid = -1

	superiorDomain, _ := strconv.ParseInt(strings.TrimSpace(attr["superiorDomain"]), 10, 64)
	redashname := strings.TrimSpace(attr["redashname"])
	redashURL := strings.TrimSpace(attr["redashURL"])
	redashIndex := strings.TrimSpace(attr["redashIndex"])
	type_str := strings.TrimSpace(attr["type"])

	info := make(map[string]string)
	info[redashIndex] = redashURL
	info_str, _ := json.Marshal(info)

	pid_init := 20
	urlSql := fmt.Sprintf(" SELECT lid from top_label where title='%s'", type_str)
	orm.NewOrm().Raw(urlSql).QueryRow(&pid_init)
	//if Pid_init == nil {
	//	Pid_init = 20
	//}

	var rd TopLabel
	rd.Title = redashname
	rd.Index = "redashshow/" + redashIndex
	rd.Info = string(info_str)
	rd.Pid = int64(pid_init)
	rd.NameSpace = "redash"

	newLid, err := orm.NewOrm().Insert(&rd)

	(&DomainLvpi{
		Name:          rd.Title,
		Administrator: "", // 二级利率没有管理员
		Pid:           superiorDomain,
		Index:         rd.Index,
		Info:          rd.Info,
	}).Add()

	if err != nil {
		Logger.Error(err.Error())
		newLid = -1
	}
	return newLid
}

func CheckRedash(attr map[string]string) map[string]string {
	o := orm.NewOrm()

	var resTitle []string
	var resNum int64
	info := make(map[string]string)

	redashURL := strings.TrimSpace(attr["redashURL"])
	checkURL := strings.Split(redashURL, "&refresh")[0]

	urlSql := fmt.Sprintf(" SELECT title from top_label where info like '%%%s%%'", checkURL)

	resNum, _ = o.Raw(urlSql).QueryRows(&resTitle)
	if resNum != 0 {
		info["code"] = "0001"
		info["title"] = resTitle[0]
	} else {
		info["code"] = "0000"
	}

	return info
}

func GetAllRedashInfo(redash_type string) []FrontRedashAllTableRes {
	o := orm.NewOrm()
	var resData []FrontRedashAllTableRes

	urlSql := fmt.Sprintf(" SELECT lid,pid, title,`index` ,info from top_label where name_space='redash' and `index` like '%s%%' and info!='' order by lid ", redash_type)
	_, err := o.Raw(urlSql).QueryRows(&resData)
	if err != nil {
		fmt.Printf("SQL select  :%s\n", err)
	}

	for idx, frontRedashRow := range resData {
		url_str := frontRedashRow.Info
		idx_str := strings.Split(frontRedashRow.Index, "/")[1]

		mapInfo := make(map[string]string)
		err := json.Unmarshal([]byte(url_str), &mapInfo)
		if err != nil {
			fmt.Printf("Unmarshal Info Error :%s\n", idx_str)
			continue
		}
		url := mapInfo[idx_str]
		resData[idx].Info = url
	}
	return resData
}

func InsertDashBoard(dashname string, usrid string) (int, string) {
	res_dashid := -1
	res_dashurl := ""

	dashname = strings.TrimSpace(dashname)
	usrid = strings.TrimSpace(usrid)

	var request_map map[string]string
	request_map = make(map[string]string)

	request_map["name"] = dashname
	request_map["usrid"] = usrid

	byte_param, _ := json.Marshal(request_map)
	param := string(byte_param)

	//通过redash 接口插入
	url := fmt.Sprintf("http://inner-redash.adtech.sogou/api/dashboards")
	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body([]byte(param))

	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
		return -1, ""
	}

	var rdb FrontRedashDashboard
	err := json.Unmarshal(bytes, &rdb)
	if err != nil {
		Logger.Error(error.Error())
		return 0, ""
	}
	res_dashid = rdb.DashboradId
	slug := rdb.Url
	res_dashurl = fmt.Sprintf("http://inner-redash.adtech.sogou/dashboard/%s", rdb.Url)

	rdb.Creator = usrid
	rdb.Name = dashname
	rdb.Url = res_dashurl
	rdb.DashboradId = res_dashid
	rdb.Index = fmt.Sprintf("dashboardshow/%s", slug)

	_, err = orm.NewOrm().Insert(&rdb)
	if err != nil {
		Logger.Error(err.Error())
		return -1, ""
	}
	/*
		//插入TopLabel
		pid_init := 38
		urlSql :=fmt.Sprintf(" SELECT lid from top_label where title='%s'","DashBoard")
		orm.NewOrm().Raw(urlSql).QueryRow(&pid_init)

		info := make(map[string]string)
		info[slug] = res_dashurl
		info_str,_ := json.Marshal(info)

		var rd TopLabel
		rd.Title =  dashname
		rd.Index =  fmt.Sprintf("dashboardshow/%s",slug)
		rd.Info  =  string(info_str)
		rd.Pid   =   int64(pid_init)
		rd.NameSpace = "redash"

		//newLid , err := orm.NewOrm().Insert(&rd)
		_ , err = orm.NewOrm().Insert(&rd)
		if err != nil {
			Logger.Error(err.Error())
		}*/
	return res_dashid, res_dashurl

}

func DeleteDashBoard(index string) bool {
	suc_flag := false
	index = strings.TrimSpace(index)

	o := orm.NewOrm()
	urlSql := fmt.Sprintf(" delete from front_redash_dashboard where `index`='%s' ", index)
	res, err := o.Raw(urlSql).Exec()
	if err != nil {
		fmt.Printf("Delete get err:%s | res:%s", err, res)
		return suc_flag
	}

	array := strings.Split(index, "/")
	if len(array) != 2 {
		return suc_flag
	}
	slug := array[1]
	//通过redash 接口插入
	url := fmt.Sprintf("http://inner-redash.adtech.sogou/api/dashboards/%s", slug)
	req := httplib.Delete(url).SetTimeout(time.Second*5, time.Second*5)

	_, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
		return suc_flag
	}

	suc_flag = true
	return suc_flag
}

func UpdateDashBoard(index string, name string, user string) bool {
	suc_flag := false
	index = strings.TrimSpace(index)

	o := orm.NewOrm()

	if len(name) == 0 {
		urlSql := fmt.Sprintf(" select name from  front_redash_dashboard where `index`='%s' ", index)
		_ = orm.NewOrm().Raw(urlSql).QueryRow(&name)
	}
	if len(user) == 0 {
		urlSql := fmt.Sprintf(" select creator from  front_redash_dashboard where `index`='%s' ", index)
		_ = orm.NewOrm().Raw(urlSql).QueryRow(&user)
	}

	dsid := 0
	urlSql := fmt.Sprintf(" select dashborad_id from front_redash_dashboard where `index`='%s' ", index)
	err := o.Raw(urlSql).QueryRow(&dsid)
	if err != nil || dsid == 0 {
		fmt.Printf("Update get err:%s | sql:%s", err, urlSql)
		return suc_flag
	}

	var request_map map[string]string
	request_map = make(map[string]string)
	request_map["name"] = name
	request_map["usrid"] = user
	byte_param, _ := json.Marshal(request_map)
	param := string(byte_param)

	//通过redash 接口修改
	url := fmt.Sprintf("http://inner-redash.adtech.sogou/api/dashboards/%d", dsid)
	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body([]byte(param))

	_, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
		return suc_flag
	}

	urlSql = fmt.Sprintf(" update front_redash_dashboard set name='%s',creator='%s' where `index`='%s' ", name, user, index)
	res, err := o.Raw(urlSql).Exec()
	if err != nil || dsid == 0 {
		fmt.Printf("Update get err:%s | sql:%s", err, res)
		return suc_flag
	}

	suc_flag = true
	return suc_flag
}

func GetIndexIsWrite(index string, roles []string) bool {
	roles_str := strings.Join(roles, "','")
	roles_str = fmt.Sprintf(("'%s'"), roles_str)

	res_lid := -1
	sql := fmt.Sprintf("SELECT lid FROM name_space_writable WHERE role_name IN (%s) AND lid IN ( SELECT lid FROM top_label WHERE `index` = '%s')", roles_str, index)
	orm.NewOrm().Raw(sql).QueryRow(&res_lid)

	return res_lid != -1
}

func GetAllDashboardInfo(user string, canWrite bool) []FrontViewInfo {
	query_limit := ""
	if canWrite == false {
		query_limit = fmt.Sprintf(" where creator='%s' ", user)
	}

	var resInfo []FrontViewInfo
	sql := fmt.Sprintf("select `name`,url,`index`,creator from front_redash_dashboard %s", query_limit)
	orm.NewOrm().Raw(sql).QueryRows(&resInfo)
	return resInfo
}

func GetAllRedashQueryInfo(user string, canWrite bool) []FrontViewInfo {
	query_limit := ""
	if canWrite == false {
		query_limit = fmt.Sprintf(" where creator='%s' ", user)
	}

	var resInfo []FrontViewInfo
	sql := fmt.Sprintf("select `name`,url,`index`,creator from front_redash_query %s", query_limit)
	orm.NewOrm().Raw(sql).QueryRows(&resInfo)
	return resInfo
}

func DeleteRedashQuery(index string) bool {
	suc_flag := false
	index = strings.TrimSpace(index)

	o := orm.NewOrm()
	urlSql := fmt.Sprintf(" delete from front_redash_query where `index`='%s' ", index)
	res, err := o.Raw(urlSql).Exec()
	if err != nil {
		fmt.Printf("Delete get err:%s | res:%s", err, res)
		return suc_flag
	}

	array := strings.Split(index, "/")
	if len(array) != 2 {
		return suc_flag
	}
	slug := array[1]
	//通过redash 接口插入
	url := fmt.Sprintf("http://inner-redash.adtech.sogou/api/queries/%s", slug)
	req := httplib.Delete(url).SetTimeout(time.Second*5, time.Second*5)

	_, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
		return suc_flag
	}

	suc_flag = true
	return suc_flag
}

//更新Redash里的  queryid；成功返回redashid，失败返回0
func UpdateRedashQueryCreator(queryid int, newUser string) bool {
	suc_flag := false

	//更新redash
	var request_map map[string]interface{}
	request_map = make(map[string]interface{})
	request_map["tags"] = []string{newUser}

	byte_param, _ := json.Marshal(request_map)
	param := string(byte_param)

	url := fmt.Sprintf("http://inner-redash.adtech.sogou/api/queries/%d", queryid)
	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body([]byte(param))

	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
	}
	_, err := strconv.Atoi(string(bytes))
	if err != nil {
		Logger.Error(error.Error())
		suc_flag = false
		return suc_flag
	}

	//更新本地
	urlSql := fmt.Sprintf(" update front_redash_query set creator='%s' where query_id='%d' ", newUser, queryid)
	_, err = orm.NewOrm().Raw(urlSql).Exec()
	if err != nil {
		fmt.Printf("UpdateRedashQueryCreator SQL get err:%s | sql:%s", err, urlSql)
		suc_flag = false
	} else {
		suc_flag = true
	}
	return suc_flag
}

func ChangeAssignUser(ass_type string, index string, newUser string) bool {
	suc_flag := false
	index = strings.TrimSpace(index)
	urlSql := ""

	if ass_type == "dashboard" {
		suc_flag = UpdateDashBoard(index, "", newUser)
		//urlSql = fmt.Sprintf(" update front_redash_dashboard set creator='%s' where `index`='%s' ", newUser, index)
	} else {
		o := orm.NewOrm()
		query_id := 0
		urlSql = fmt.Sprintf(" select query_id from front_redash_query  where `index`='%s' ", index)
		err := o.Raw(urlSql).QueryRow(&query_id)
		if err != nil || query_id == 0 {
			return suc_flag
		}
		suc_flag = UpdateRedashQueryCreator(query_id, newUser)
	}
	return suc_flag
}
