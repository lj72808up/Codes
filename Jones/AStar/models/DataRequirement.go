package models

import (
	"AStar/utils"
	b64 "encoding/base64"
	"github.com/astaxie/beego/orm"
	"time"
)

type DataRequirement struct {
	Rid         int64         `orm:"pk;auto" field:"id,hidden"`
	ProjectName string        `orm:"size(200)" field:"需求名称"`
	CreateTime  string        `orm:"size(20)" field:"需求提出时间"`
	StartTime   string        `orm:"size(20)" field:"开始时间"`
	EndTime     string        `orm:"size(20)" field:"结束时间"`
	Apartment   string        `orm:"size(20)" field:"来源部门"`
	ExpectTime  string        `orm:"size(20)" field:"预期完成时间"`
	Proposer    string        `orm:"size(50)" field:"提出人"`
	Developer   string        `orm:"size(50)" field:"开发者"`
	DevelopTime string        `orm:"size(20)" field:"开发时间,hidden"`
	Status      DevelopStatus `orm:"size(3)" field:"状态"`
	DemandText  string        `orm:"type(text)" field:"文档描述,drop"` // 需求文本
	Priority    string        `orm:"size(3)" field:"优先级"`          // 优先级
}

type DataRequirementHistory struct {
	Hid         int64  `orm:"pk;auto" field:"历史ID"` // 历史ID
	Operator    string `orm:"size(50)" field:"修改人"`
	OperateTime string `orm:"size(50)" field:"修改时间"`
	// 需求表
	Rid         int64  `field:"原始id,hidden"` // rid 一旦提交就不会改变.
	ProjectName string `orm:"size(200)" field:"需求名称"`
	//CreateTime  string        `orm:"size(20)" field:"需求提出时间"`
	StartTime  string `orm:"size(20)" field:"开始时间"`
	EndTime    string `orm:"size(20)" field:"结束时间"`
	Apartment  string `orm:"size(20)" field:"来源部门"`
	ExpectTime string `orm:"size(20)" field:"预期完成时间"`
	//Proposer    string        `orm:"size(50)" field:"需求提出人"`
	Developer string `orm:"size(50)" field:"开发者"`
	//DevelopTime string        `orm:"size(20)" field:"开发时间,hidden"`
	Status     DevelopStatus `orm:"size(3)" field:"状态"`
	DemandText string        `orm:"type(text)" field:"文档描述,hidden"` // 需求文本
	Priority   string        `orm:"size(3)" field:"优先级"`            // 优先级
}

func (this *DataRequirementHistory) InceptDemand(demand DataRequirement) {
	this.Rid = demand.Rid
	this.ProjectName = demand.ProjectName
	//this.CreateTime = demand.CreateTime
	this.StartTime = demand.StartTime
	this.EndTime = demand.EndTime
	this.Apartment = demand.Apartment
	this.ExpectTime = demand.ExpectTime
	//this.Proposer = demand.Proposer
	this.Developer = demand.Developer
	//this.DevelopTime = demand.DevelopTime
	this.Status = demand.Status
	this.DemandText = demand.DemandText
	this.Priority = demand.Priority
}

func GetAllDemandHistories(rid int64) utils.TableData {
	o := orm.NewOrm()
	res := []DataRequirementHistory{}
	_, _ = o.Raw("SELECT * FROM data_requirement_history WHERE rid = ? ORDER BY hid DESC", rid).QueryRows(&res)
	data := transView(res)
	return data
}

func transView(views []DataRequirementHistory) utils.TableData {
	var rows []map[string]interface{}
	for _, view := range views {
		tagRow := utils.GetSnakeStrMapOfStruct(view)
		rows = append(rows, tagRow)
	}
	table := utils.TableData{Data: rows}
	table.GetSnakeTitles(DataRequirementHistory{})
	return table
}

type DevelopStatus int

const (
	Pending       DevelopStatus = 1  //  "待处理"
	Planning      DevelopStatus = 2  // 规划中（排期中)
	DemandReview  DevelopStatus = 3  // 需求评审中
	TechReview    DevelopStatus = 4  // 技术评审中
	Developing    DevelopStatus = 5  //  "开发中"
	Finish        DevelopStatus = 6  //  "开发完成"
	Examine       DevelopStatus = 7  // 验收中
	FinishProduct DevelopStatus = 8  //  "完成（产品）"
	FinishTech    DevelopStatus = 9  // 完成（技术）
	HangUp        DevelopStatus = 10 // 挂起
	//ExamineFinish DevelopStatus = 11 // 验收完成
)

func GetAllDevelopStatus() map[DevelopStatus]string {
	statusCollection := map[DevelopStatus]string{
		Pending: "待处理", Planning: "规划中（排期中)", DemandReview: "需求评审中", TechReview: "技术评审中",
		Developing: "开发中", Finish: "开发完成", Examine: "验收完成", FinishProduct: "完成（产品)", FinishTech: "完成（技术）", HangUp: "挂起",
	}
	return statusCollection
}

func GetAllDataRequests(canEdit bool, proposer string) utils.TableData {

	o := orm.NewOrm()
	var demands []DataRequirement
	if canEdit {
		_, _ = o.Raw("SELECT * FROM data_requirement ORDER BY rid DESC").QueryRows(&demands)
	} else {
		_, _ = o.Raw("SELECT * FROM data_requirement WHERE proposer = ? ORDER BY rid DESC", proposer).QueryRows(&demands)
	}

	var rows []map[string]interface{}
	for _, view := range demands {
		tagRow := utils.GetSnakeStrMapOfStruct(view)
		rows = append(rows, tagRow)
	}
	table := utils.TableData{Data: rows}
	table.GetSnakeTitles(DataRequirement{})
	return table
}

func (this *DataRequirement) Add(proposer string) (int64, error) {
	o := orm.NewOrm()
	this.CreateTime = utils.TimeFormat(time.Now())
	this.Proposer = proposer
	_, err := o.Insert(this)
	if err != nil {
		Logger.Error("%v", err)
	} else {
		Logger.Info("成功创建新需求, id为%d", this.Rid)
	}
	return this.Rid, err
}

func (this *DataRequirement) DeleteById(id int) {
	o := orm.NewOrm()
	_, _ = o.Raw("DELETE FROM data_requirement WHERE rid = ?", id).Exec()
}

func (this *DataRequirement) Update() error {
	// 更新需求描述与优先级 (description, priority)
	o := orm.NewOrm()
	_, err := o.Update(this, "demand_text", "priority")
	return err
}

type DemandUserDefineIndicator struct {
	Udid        int64  `orm:"pk;auto"`
	Rid         int64  // 需求Id
	AssetName   string `orm:"size(255)"`
	Description string `orm:"size(255)"`
}

func (this *DataRequirement) AddUserDefineIndicator(requestId int64, userDefineIndicators *[]DemandUserDefineIndicator) error {
	o := orm.NewOrm()
	_, _ = o.Raw("DELETE FROM demand_user_define_indicator WHERE rid = ?", requestId).Exec()
	_, err := o.InsertMulti(100, userDefineIndicators)
	if err != nil {
		Logger.Error("%v", err)
	}
	return err
}

func (this *DataRequirement) GetUserDefineIndicators(requestId int64) []DemandUserDefineIndicator {
	o := orm.NewOrm()
	userDefineIndicators := []DemandUserDefineIndicator{}
	_, _ = o.Raw("SELECT * FROM demand_user_define_indicator WHERE rid = ?", requestId).QueryRows(&userDefineIndicators)
	return userDefineIndicators
}

type DemandAttachment struct {
	Aid       int64  `orm:"pk;auto"`
	Filename  string `orm:"size(200)"`
	Base64Doc string `orm:"type(text)"` // 文件内容, base64转码
	Rid       int64  // 关联的需求ID
}

// 上传二进制文档
func (this *DemandAttachment) Add(fileName string, binDoc *[]byte, requestId int64) error {
	base64Doc := b64.StdEncoding.EncodeToString(*binDoc)
	this.Filename = fileName
	this.Base64Doc = base64Doc
	this.Rid = requestId
	o := orm.NewOrm()
	_, err := o.Insert(this)
	return err
}

func (this *DemandAttachment) DeleteById(id int64) error {
	o := orm.NewOrm()
	_, err := o.Raw("DELETE FROM demand_attachment WHERE aid=?", id).Exec()
	if err != nil {
		utils.Logger.Error("%v", err)
	}
	return err
}

type AttachmentView struct {
	Filename  string
	Base64Doc string
}

func (this *DemandAttachment) GetById(id int64) (string, *[]byte) {
	attchment := AttachmentView{}
	o := orm.NewOrm()
	_ = o.Raw("SELECT filename, base64_doc FROM demand_attachment WHERE Aid = ?", id).QueryRow(&attchment)
	binaryDoc, _ := b64.StdEncoding.DecodeString(attchment.Base64Doc)
	return attchment.Filename, &binaryDoc
}

type DemandAssetRelation struct {
	Daid    int64 `orm:"pk;auto" field:"id,hidden"`
	Rid     int64
	Assetid int64
	Atype   string `orm:"size(20)"`
}

// 更新需求和数据资产的关系
func (this *DemandAssetRelation) Update(relations *[]DemandAssetRelation) error {
	delSql := "DELETE FROM demand_asset_relation WHERE Rid = ?"
	o := orm.NewOrm()
	requestId := (*relations)[0].Rid
	_, _ = o.Raw(delSql, requestId).Exec()
	_, err := o.InsertMulti(100, relations)
	return err
}

type DataRequirementComment struct {
	Cid         int64  `orm:"pk;auto"`
	DemandId    int64  // 需求ID
	Comment     string `orm:"size(255)"` // 评论
	Commentator string `orm:"size(50)"`  // 评论人
	Time        string `orm:"size(50)"`  // 评论时间
}

func (this *DataRequirementComment) Add() error {
	o := orm.NewOrm()
	this.Time = utils.TimeFormat(time.Now())
	_, err := o.Insert(this)
	return err
}

type UserNameMapping struct {
	Id                int64  `orm:"pk;auto"`
	Username          string `orm:"size(50)"`
	Uid               string `orm:"size(50);unique"`
	LastTrigLoginTime string `orm:"size(50)"`
}

func (this *UserNameMapping) AddOrUpdate() {
	this.LastTrigLoginTime = utils.TimeFormat(time.Now())
	o := orm.NewOrm()
	_, _ = o.Raw("REPLACE INTO user_name_mapping(username, uid, last_trig_login_time) VALUES (?,?,?)",
		this.Username, this.Uid, this.LastTrigLoginTime).Exec()
}

func GetUsernameFromMapping(uid string) string {
	o := orm.NewOrm()
	username := ""
	_ = o.Raw("SELECT username FROM user_name_mapping WHERE uid = ?", uid).QueryRow(&username)
	return username
}
