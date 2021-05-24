package models

import (
	pb "AStar/protobuf"
	"github.com/astaxie/beego/orm"
)

// 小流量任务的struct
type
SmallFlowTask struct {
	Id           int    ` field:"主键,hidden"`
	TaskId       string `orm:"size(32)"  field:"任务ID"`
	TaskName     string ` field:"任务名称"`
	TaskType     string ` field:"任务类型"`
	StartDate    string ` field:"开始日期,hidden"`
	EndDate      string ` field:"结束日期,hidden"`
	Session      string ` field:"用户名"`
	IsEnabled    int   ` field:"是否启用,hidden"`
	FullSql      string `orm:"type(text)" field:"全量SQL,hidden"`
	SameQuerySql string `orm:"type(text)" field:"同质SQL,hidden"`
	Json         string `orm:"type(text)" field:"JSON串,hidden"`
}

// SmallFlowTask的insert ORM 方法
func (smallFlow SmallFlowTask) InsertSmallFlowTask(task pb.SmallFlow) {
	smallFlow.TaskId = task.TaskId
	smallFlow.TaskName = task.TaskName
	smallFlow.TaskType = task.TaskType
	smallFlow.StartDate = task.StartDate
	smallFlow.EndDate = task.EndDate
	smallFlow.Session = task.Session
	smallFlow.IsEnabled = 0
	smallFlow.FullSql = task.Sqls[0]
	smallFlow.SameQuerySql = task.Sqls[1]
	smallFlow.Json = task.Json
	_, err := orm.NewOrm().Insert(&smallFlow)
	if err != nil {
		Logger.Error(err.Error())
	}
}

func (smallFlow* SmallFlowTask) UpdateEnable() error {
	_, err := orm.NewOrm().Update(smallFlow, "IsEnabled")
	if err != nil {
		Logger.Error(err.Error())
	}
	return err
}

func (smallFlow SmallFlowTask) GetQueryJson(taskId string) (string, error) {
	o := orm.NewOrm()
	submitLog := SmallFlowTask{TaskId: taskId}
	err := o.Read(&submitLog, "task_id")
	if err == nil {
		return submitLog.Json, nil
	}
	return "", err
}

type FlowTaskPcDef struct {
	Name           string         `json:"name"`
	DateRange      []string       `json:"date_range"`
	Type           string         `json:"type"`
	AdPosRestrict  []string       `json:"ad_pos_restrict"`
	IpRestrict     []string       `json:"ip_restrict"`
	ExprExtendReserver StyleReserver `json:"expr_extend_reserve"`
	CtrlExtendReserver StyleReserver `json:"ctrl_extend_reserve"`
}

type FlowTaskDef struct {
	Name           string         `json:"name"`
	DateRange      []string       `json:"date_range"`
	Type           string         `json:"type"`
	AdPosRestrict  []string       `json:"ad_pos_restrict"`
	IpRestrict     []string       `json:"ip_restrict"`
	ExprFlow       FlowDef        `json:"expr_flow"`
	CtrlFlow       FlowDef        `json:"ctrl_flow"`
	ExtendReserver ExtendReserver `json:"extend_reserve,omitempty"`
	StyleReserver  StyleReserver  `json:"style_reserve,omitempty"`
	ClickFlag      ClickFlag      `json:"click_flag,omitempty"`
}


type SmallFlowRes struct {
	Id         int64  `orm:"null" field:"主键,hidden"`
	Dt         string `orm:"null" field:"日期,hidden"`
	Tag        string `orm:"null" field:"组别;组别;组别"`
	Pw         string `orm:"null" field:"位次;位次;位次"`
	SumPvPage  string `orm:"null" field:"页面ADPV;实验样式ADPV;页面ADPV"`
	SumPvAd    string `orm:"null" field:"页面广告条数;实验样式广告条数;页面广告条数"`
	SumCharge  string `orm:"null" field:"页面消耗(元);实验样式消耗(元);页面消耗(元)"`
	SumClick   string `orm:"null" field:"页面点击;实验样式点击;页面点击"`
	AvgCtrPage string `orm:"null" field:"页面CTR3(%);实验样式CTR3(%);页面CTR3(%)"`
	AvgCtrAd   string `orm:"null" field:"页面CTR2(%);实验样式CTR2(%);页面CTR2(%)"`
	AvgAcp     string `orm:"null" field:"页面ACP;实验样式ACP;页面ACP"`
	AvgRpmPage string `orm:"null" field:"页面RPM3;实验样式RPM3;页面RPM3"`
	CustomId   string `orm:"null" field:"客户ID,hidden"`
}
