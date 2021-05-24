package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type DataRequirementController struct {
	BaseController
}

var DataRequirementSpace = "requirement"

// @Title 获取所有需求列表
// @Summary 获取所有需求列表
// @Description 获取所有需求列表
// @Success 200
// @Accept json
// @router /getAllDemand [get]
func (this *DataRequirementController) GetAllDemand() {
	checkUrl := fmt.Sprintf("%s/%s_%s", RootNameSpace, DataRequirementSpace, "authority")
	canEdit := models.Enforcer.Enforce(this.userName, checkUrl,  "POST")
	utils.Logger.Info("当前用户%s对需求列表是否有编辑权限: %t",this.userName, canEdit)
	table := models.GetAllDataRequests(canEdit, this.userName)
	bytes, _ := json.Marshal(table)
	this.Ctx.WriteString(string(bytes))
}

// @Title 根据Id获取需求
// @Summary 根据Id获取需求
// @Description 根据Id获取需求
// @Success 200
// @Accept json
// @router /getDemandById/:rid [get]
func (this *DataRequirementController) GetDemandById() {
	rid, _ := strconv.ParseInt(this.Ctx.Input.Param(":rid"), 10, 64)
	o := orm.NewOrm()
	var demand models.DataRequirement
	_ = o.Raw("SELECT * FROM data_requirement WHERE rid=?", rid).QueryRow(&demand)
	bytes, _ := json.Marshal(demand)
	this.Ctx.WriteString(string(bytes))
}

// @Title 新增数据需求
// @Summary 新增数据需求
// @Description 新增数据需求
// @Success 200
// @Accept json
// @router /post [post]
func (this *DataRequirementController) Post() {
	var requirement models.DataRequirement
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &requirement)
	rid, err := requirement.Add(this.userName)

	newHistory := models.DataRequirementHistory{}
	newHistory.InceptDemand(requirement)
	newHistory.Operator = this.userName
	newHistory.OperateTime = utils.TimeFormat(time.Now())
	o := orm.NewOrm()
	_,_ = o.Insert(&newHistory)

	var info string
	if err == nil {
		info = fmt.Sprintf("%d", rid)
	} else {
		info = err.Error()
	}
	this.Ctx.WriteString(utils.VueResponse{
		Res:  err == nil,
		Info: info,
	}.String())
}

// @Title 更新数据需求
// @Summary 更新数据需求
// @Description 更新数据需求
// @Success 200
// @Accept json
// @router /update/:rid [post]
func (this *DataRequirementController) Update() {
	rid, _ := strconv.ParseInt(this.Ctx.Input.Param(":rid"), 10, 64)
	var newDemand models.DataRequirement
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &newDemand)
	newDemand.Rid = rid
	o := orm.NewOrm()
	oldDemand := models.DataRequirement{}
	_ = o.Raw("SELECT * FROM data_requirement WHERE rid = ?", rid).QueryRow(&oldDemand)
	// 保存旧的历史
/*	history := models.DataRequirementHistory{}
	history.InceptDemand(oldDemand)
	history.Operator = this.userName
	history.OperateTime = utils.TimeFormat(time.Now())
	_,_ = o.Insert(&history)*/
	// Todo 历史表插入新的需求
	newHistory := models.DataRequirementHistory{}
	newHistory.InceptDemand(newDemand)
	newHistory.Operator = this.userName
	newHistory.OperateTime = utils.TimeFormat(time.Now())
	_,_ = o.Insert(&newHistory)
	// 更新新的东西
	_, err := o.Update(&newDemand, "ProjectName", "ExpectTime", "Priority", "Status", "DemandText", "Developer","Apartment","StartTime","EndTime")

	var info string
	if err == nil {
		info = "更新成功"
	} else {
		info = err.Error()
	}
	this.Ctx.WriteString(utils.VueResponse{
		Res:  err == nil,
		Info: info,
	}.String())
}

// @Title 获取所有需求历史
// @Summary 获取所有需求历史
// @Description 获取所有需求历史
// @Success 200
// @Accept json
// @router /getHistory/:rid [get]
func (this *DataRequirementController) GetHistory() {
	rid, _ := strconv.ParseInt(this.Ctx.Input.Param(":rid"), 10, 64)
	histories := models.GetAllDemandHistories(rid)
	bytes,_ := json.Marshal(histories)
	this.Ctx.WriteString(string(bytes))
}


// @Title 上传需求文档
// @Summary 上传需求文档
// @Description 上传需求文档
// @Success 200
// @Accept json
// @router /postAttachment/:requestId [post]
func (this *DataRequirementController) PostAttachment() {
	requestId, _ := strconv.ParseInt(this.Ctx.Input.Params()[":requestId"], 10, 64)
	f, h, err := this.GetFile("file")
	if err != nil {
		models.Logger.Error("%v", err)
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	} else {
		fileName := h.Filename
		binaryDoc, _ := utils.ReadBytes(f, h.Size)
		attachment := models.DemandAttachment{}
		if err := attachment.Add(fileName, &binaryDoc, requestId); err != nil {
			this.Ctx.WriteString(utils.VueResponse{
				Res:  false,
				Info: err.Error(),
			}.String())
		} else {
			this.Ctx.WriteString(utils.VueResponse{
				Res:  true,
				Info: "上传成功",
			}.String())
		}
	}
}

func transAssetType(chineseWord string) string {
	transType := ""
	if chineseWord == "原子指标" {
		transType = "predefine"
	} else if chineseWord == "计算指标" {
		transType = "combine"
	} else if chineseWord == "维度" {
		transType = "dimension"
	}
	return transType
}

// @Title 更新需求与数据资产关系
// @Summary 更新需求与数据资产关系
// @Description 更新需求与数据资产关系
// @Success 200
// @Accept json
// @router /postRequireAssetRelation/:requestId [post]
func (this *DataRequirementController) PostRequireAssetRelation() {
	requestId, _ := strconv.ParseInt(this.Ctx.Input.Params()[":requestId"], 10, 64)
	println(requestId)
	relations := []models.DemandAssetRelation{}
	var assets []map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &assets)
	// 转化参数为DemandAssetRelation列表
	for _, asset := range assets {
		assetId, _ := strconv.ParseInt(asset["origin_id"], 10, 64)
		assetType := transAssetType(asset["type_name"])
		r := models.DemandAssetRelation{
			Rid:     requestId, // 需求Id
			Assetid: assetId,
			Atype:   assetType,
		}
		relations = append(relations, r)
	}
	rModel := models.DemandAssetRelation{}

	if err := rModel.Update(&relations); err != nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "更新成功",
		}.String())
	}
}

// @Title 获取开发状态映射
// @Summary 获取开发状态映射
// @Description 获取开发状态映射
// @Success 200
// @Accept json
// @router /GetDevStatusMapping [get]
func (this *DataRequirementController) GetDevStatusMapping() {
	devStatusMapping := models.GetAllDevelopStatus()
	bytes, _ := json.Marshal(devStatusMapping)
	this.Ctx.WriteString(string(bytes))
}

// @Title 获取需求关联的数据资产
// @Summary 获取需求关联的数据资产
// @Description 获取需求关联的数据资产
// @Success 200
// @Accept json
// @router /getAssetRelationById/:rid [get]
func (this *DataRequirementController) GetAssetRelationById() {
	rid, _ := strconv.ParseInt(this.Ctx.Input.Param(":rid"), 10, 64)
	o := orm.NewOrm()
	relations := []models.DemandAssetRelation{}
	_, _ = o.Raw("SELECT * FROM demand_asset_relation WHERE rid=?", rid).QueryRows(&relations)
	bytes, _ := json.Marshal(relations)
	this.Ctx.WriteString(string(bytes))
}

// @Title 上传需求附件
// @Summary 上传需求附件
// @Description 上传需求附件
// @Success 200
// @Accept json
// @router /attachment/upload/:rid [post]
func (this *DataRequirementController) UploadAttachment() {
	rid, _ := strconv.ParseInt(this.Ctx.Input.Param(":rid"), 10, 64)
	f, h, _ := this.GetFile("file")
	defer f.Close()
	bytes, _ := utils.ReadBytes(f, h.Size)
	// 附件存储数据库
	attachment := models.DemandAttachment{}
	if err := attachment.Add(h.Filename, &bytes, rid); err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "附件上传成功",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	}
}

// @Title 下载需求附件
// @Summary 下载需求附件
// @Description 下载需求附件
// @Success 200
// @Accept json
// @router /attachment/download/:aid [post]
func (this *DataRequirementController) DownloadAttachment() {
	aid, _ := strconv.ParseInt(this.Ctx.Input.Param(":aid"), 10, 64)
	attachment := models.DemandAttachment{}
	fileName, contentBytes := attachment.GetById(aid)

	dir := "download/demand/doc"
	toFile := fmt.Sprintf("%s/%s", dir, fileName)
	utils.MkDir(dir)
	if err := ioutil.WriteFile(toFile, *contentBytes, 0777); err != nil {
		utils.Logger.Error(err.Error())
		panic(err)
	}
	utils.Logger.Info("写入附件到本地:%s 成功", toFile)
	this.Ctx.Output.Download(toFile)
	_ = os.Remove(toFile)
}

type AttachmentView struct {
	Aid      int64
	Filename string
}

// @Title 获取上传需求附件列表
// @Summary 获取上传需求附件列表
// @Description 获取上传需求附件列表
// @Success 200
// @Accept json
// @router /attachment/get/:rid [get]
func (this *DataRequirementController) GetAttachmentList() {
	rid, _ := strconv.ParseInt(this.Ctx.Input.Param(":rid"), 10, 64)
	o := orm.NewOrm()
	views := []AttachmentView{}
	if _, err := o.Raw("SELECT aid, filename FROM demand_attachment WHERE rid = ?", rid).QueryRows(&views); err != nil {
		utils.Logger.Error("%v", err)
	}
	bytes, _ := json.Marshal(views)
	this.Ctx.WriteString(string(bytes))
}

// @Title 删除上传的附件
// @Summary 删除上传的附件
// @Description 删除上传的附件
// @Success 200
// @Accept json
// @router /attachment/delete/:aid [get]
func (this *DataRequirementController) DelAttachmentById() {
	aid, _ := strconv.ParseInt(this.Ctx.Input.Param(":aid"), 10, 64)
	attachment := models.DemandAttachment{}
	_ = attachment.DeleteById(aid)
}

// @Title 增加自定义指标
// @Summary 增加自定义指标
// @Description 增加自定义指标
// @Success 200
// @Accept json
// @router /indicator/userDefine/post/:rid [post]
func (this *DataRequirementController) PostUserDefineIndicator() {
	rid, _ := strconv.ParseInt(this.Ctx.Input.Param(":rid"), 10, 64)
	var udfs []models.DemandUserDefineIndicator
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &udfs)
	demand := models.DataRequirement{}
	err := demand.AddUserDefineIndicator(rid, &udfs)
	if err != nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "添加成功",
		}.String())
	}
}

// @Title 获取自定义指标
// @Summary 获取自定义指标
// @Description 获取自定义指标
// @Success 200
// @Accept json
// @router /indicator/userDefine/get/:rid [get]
func (this *DataRequirementController) GetUserDefineIndicator() {
	rid, _ := strconv.ParseInt(this.Ctx.Input.Param(":rid"), 10, 64)
	demand := models.DataRequirement{}
	udIndicators := demand.GetUserDefineIndicators(rid)
	bytes, _ := json.Marshal(udIndicators)
	this.Ctx.WriteString(string(bytes))
}

// @Title 增加需求评论
// @Summary 增加需求评论
// @Description 增加需求评论
// @Success 200
// @Accept json
// @router /indicator/addDataRequirement [post]
func (this *DataRequirementController) AddDataRequirement() {
	dataComment := models.DataRequirementComment{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &dataComment)

	userName := models.GetUsernameFromMapping(this.userName)
	if userName == "" {
		userName = this.userName
	}
	dataComment.Commentator = userName

	err := dataComment.Add()
	if err != nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "评论成功",
		}.String())
	}
}


// @Title 展示评论
// @Summary 展示评论
// @Description 展示评论
// @Success 200
// @Accept json
// @router /indicator/getDataRequirement [get]
func (this *DataRequirementController) GetDataRequirement() {
	id := this.GetString("demandId") // 需求ID
	o := orm.NewOrm()
	commentArr := []models.DataRequirementComment{}
	_, _ = o.Raw("SELECT * FROM data_requirement_comment WHERE demand_id=?", id).QueryRows(&commentArr)
	bytes, _ := json.Marshal(commentArr)
	this.Ctx.WriteString(string(bytes))
}