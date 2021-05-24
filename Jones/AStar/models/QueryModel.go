package models

import (
	pb "AStar/protobuf"
	"AStar/utils"
	"time"
	"github.com/astaxie/beego/orm"
)

type QuerySubmitLog struct {
	Id           int            `field:"Id,drop"`
	Date         string         `field:"任务日期,drop"`
	QueryId      string         `orm:"size(32)" field:"任务ID"`
	QueryName    string         `field:"任务名称"`
	Session      string         `field:"提交用户"`
	SqlStr       string         `orm:"type(text)" field:"SQL,drop"`
	QueryJson    string         `orm:"type(text)" field:"QueryJson,drop"`
	Status       pb.QueryStatus `orm:"type(int)" field:"任务状态"`
	CreateTime   int64          `field:"任务创建时间"`
	StartTime    int64          `field:"任务创建时间,drop"`
	LastUpdate   int64          `field:"最近更改时间"`
	DataPath     string         `field:"RSYNC地址,hidden"`
	ExceptionMsg string         `orm:"type(text)" field:"异常信息,drop"`
	Engine       pb.QueryEngine `orm:"type(int)" field:"计算引擎,hidden"`
	Type         string         `field:"任务种类"`
}

func (querySubmitLog QuerySubmitLog) InsertQuerySubmitLog(queryInfo pb.QueryInfo, queryName string, queryJson string) {
	querySubmitLog.Date = utils.DateFormat(time.Now())
	querySubmitLog.QueryId = queryInfo.QueryId
	querySubmitLog.QueryName = queryName
	querySubmitLog.Session = queryInfo.Session
	querySubmitLog.SqlStr = queryInfo.Sql[0]
	querySubmitLog.QueryJson = queryJson
	querySubmitLog.Status = queryInfo.Status
	querySubmitLog.Engine = queryInfo.Engine
	querySubmitLog.LastUpdate = queryInfo.LastUpdate
	querySubmitLog.CreateTime = queryInfo.CreateTime
	querySubmitLog.Type = queryInfo.Type
	_, err := orm.NewOrm().Insert(&querySubmitLog)
	if err != nil {
		Logger.Error(err.Error())
	}
}

func (querySubmitLog QuerySubmitLog) Update(queryInfo pb.QueryInfo, cols ...string) {
	querySubmitLog.QueryId = queryInfo.QueryId
	querySubmitLog.Status = queryInfo.Status
	querySubmitLog.LastUpdate = time.Now().Unix()
	_, err := orm.NewOrm().Update(querySubmitLog, cols...)
	if err != nil {
		Logger.Error(err.Error())
	}
}

//var QuerySubmitLogHeader = map[string]string{
//	"Date":       "任务日期",
//	"QueryId":    "任务ID",
//	"Session":    "提交用户",
//	"QueryName":  "任务名称",
//	"Status":     "任务状态",
//	"CreateTime": "任务创建时间",
//	"LastUpdate": "最近更改时间",
//	"DataPath":   "RSYNC路径",
//}
//
//func GetQuerySubmitLogHeaderList() []string {
//	QuerySubmitLogHeader := []string{
//		"Date",
//		"QueryId",
//		"Session",
//		"QueryName",
//		"Status",
//		"CreateTime",
//		"LastUpdate",
//		"DataPath",
//		"StatusCode",
//	}
//	return QuerySubmitLogHeader
//}

func (querySubmitLog QuerySubmitLog) GetQueryJson(queryId string) (string, error) {
	o := orm.NewOrm()
	submitLog := QuerySubmitLog{QueryId: queryId}
	err := o.Read(&submitLog, "query_id")
	if err == nil {
		return submitLog.QueryJson, nil
	}
	return "", err
}


func (querySubmitLog *QuerySubmitLog) GetSubmitLog() {
	o := orm.NewOrm()
	err := o.Read(querySubmitLog, "query_id")
	if err != nil {
		panic(err.Error())
	}
}