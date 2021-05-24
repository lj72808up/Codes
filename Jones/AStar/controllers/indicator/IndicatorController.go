package indicator

import (
	"AStar/controllers"
	"AStar/models"
	"AStar/models/indicator"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type IndicatorController struct {
	controllers.BaseController
}

var IndicatorNameSpace = "indicator"

// @Title GetServiceTag
// @Summary 获取业务标签
// @Description 获取业务标签
// @router /serviceTag/getAll [get]
func (this *IndicatorController) GetServiceTag() {
	tag := indicator.ServiceTag{}
	vxeData := tag.GetAll()
	bytes, _ := json.Marshal(vxeData)
	this.Ctx.WriteString(string(bytes))
}

// @Title PutServiceTag
// @Summary 新增业务标签
// @Description 新增业务标签
// @router /serviceTag/putServiceTag [post]
func (this *IndicatorController) PutServiceTag() {
	attr := map[string]string{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	tag := indicator.ServiceTag{}
	tagName := attr["name"]
	tagDesc := attr["desc"]
	if err := tag.AddOne(tagName, tagDesc); err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "新增成功",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: fmt.Sprintf("%v", err),
		}.String())
	}
}

// @Title EditServiceTag
// @Summary 修改业务标签
// @Description 修改业务标签
// @router /serviceTag/editServiceTag [post]
func (this *IndicatorController) EditServiceTag() {
	tag := indicator.ServiceTag{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &tag)
	if err := tag.UpdateOne(tag.Tid, tag.TagName, tag.Description); err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "新增成功",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: fmt.Sprintf("%v", err),
		}.String())
	}
}

// @Title GetIndicator
// @Summary 新增业务标签
// @Description 新增业务标签
// @router /getAllIndicator [get]
func (this *IndicatorController) GetIndicator() {
	attr := map[string]string{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	index := indicator.Indicator{}
	vxeData, _ := index.GetPreDefineIndicator()
	bytes, _ := json.Marshal(vxeData)
	this.Ctx.WriteString(string(bytes))
}

// @Title getAllCombineIndicator
// @Summary 获取所有计算指标
// @Description 获取所有计算指标
// @router /getAllCombineIndicator [get]
func (this *IndicatorController) GetAllCombineIndicator() {
	attr := map[string]string{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	index := indicator.Indicator{}
	vxeData, _ := index.GetCombineIndicator()
	bytes, _ := json.Marshal(vxeData)
	this.Ctx.WriteString(string(bytes))
}

// @Title PostPreIndicator
// @Summary 新增原子指标
// @Description 新增原子指标
// @router /postPreIndicator [post]
func (this *IndicatorController) PostPreIndicator() {
	view := indicator.View{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &view)
	index := indicator.Indicator{}
	// 设置指标类型为预定义指标
	view.IndexType = indicator.PREDEFINE
	if id, err := index.AddOne(view); err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: fmt.Sprintf("%d", id),
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: fmt.Sprintf("%v", err),
		}.String())
	}
}

// @Title postCombineIndicator
// @Summary 新增计算指标
// @Description 新增计算指标
// @router /postCombineIndicator [post]
func (this *IndicatorController) PostCombineIndicator() {
	view := indicator.View{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &view)
	index := indicator.Indicator{}
	// 设置指标类型为预定义指标
	view.IndexType = indicator.COMBINE
	var err error
	id := view.Id
	if view.Id == 0 {
		id, err = index.AddOne(view)
	} else {
		err = index.UpdateOne(view)
	}
	if err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: fmt.Sprintf("%d", id),
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: fmt.Sprintf("%v", err),
		}.String())
	}
}

// @Title EditPreIndicator
// @Summary 修改预定义指标
// @Description 修改预定义指标
// @router /editPreIndicator [post]
func (this *IndicatorController) EditPreIndicator() {
	view := indicator.View{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &view)
	index := indicator.Indicator{}
	// 设置指标类型为预定义指标
	view.IndexType = indicator.PREDEFINE
	if err := index.UpdateOne(view); err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "修改成功",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: fmt.Sprintf("%v", err),
		}.String())
	}
}

// @Title EditCombineIndicator
// @Summary 修改计算指标
// @Description 修改计算指标
// @router /editCombineIndicator [post]
func (this *IndicatorController) EditCombineIndicator() {
	view := indicator.View{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &view)
	index := indicator.Indicator{}
	// 设置指标类型为预定义指标
	view.IndexType = indicator.COMBINE
	if err := index.UpdateOne(view); err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "修改成功",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: fmt.Sprintf("%v", err),
		}.String())
	}
}

// @Title GetIndicatorDependenciesById
// @Summary 根据指标id获取依赖
// @Description 修改计算指标
// @router /getIndicatorDependenciesById/:id [get]
func (this *IndicatorController) GetIndicatorDependenciesById() {
	idStr := this.Ctx.Input.Params()[":id"]
	id, _ := strconv.Atoi(idStr)
	dependencyIds := indicator.GetDependencyIds(id)
	bytes, _ := json.Marshal(dependencyIds)
	this.Ctx.WriteString(string(bytes))
}

// @Title GetIndicatorFactRelation
// @Summary 获取原子指标关联的事实表
// @Description 获取原子指标关联的事实表
// @router /getIndicatorFactRelation/:id [get]
func (this *IndicatorController) GetIndicatorFactRelation() {
	idStr := this.Ctx.Input.Params()[":id"]
	id, _ := strconv.Atoi(idStr)
	indexFactRelation := indicator.IndicatorFactRelation{}
	table := indexFactRelation.GetFactRelations(id)
	bytes, _ := json.Marshal(table)
	this.Ctx.WriteString(string(bytes))
}

// @Title AddIndicatorFactRelation
// @Summary 新增原子指标关联的事实表
// @Description 新增原子指标关联的事实表
// @router /predefine/addIndicatorFactRelation [post]
func (this *IndicatorController) AddIndicatorFactRelation() {
	factRalation := indicator.IndicatorFactRelation{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &factRalation)
	if err := factRalation.Add(); err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "修改成功",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: fmt.Sprintf("%v", err),
		}.String())
	}
}

// @Title DelIndicatorFactRelation
// @Summary 删除原子指标关联的事实表
// @Description 删除原子指标关联的事实表
// @router /predefine/delIndicatorFactRelation/:id [post]
func (this *IndicatorController) DelIndicatorFactRelation() {
	idStr := this.Ctx.Input.Params()[":id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	o := orm.NewOrm()
	if _, err := o.Raw("DELETE FROM indicator_fact_relation WHERE ifid = ?", id).Exec(); err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "修改成功",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: fmt.Sprintf("%v", err),
		}.String())
	}
}

// @Title DeleteIndicator
// @Summary 删除指标
// @Description 删除指标
// @router /delIndicator/:id [post]
func (this *IndicatorController) DelIndicator() {
	idStr := this.Ctx.Input.Params()[":id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	index := indicator.Indicator{}
	err := index.DelteById(id)
	if err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "删除成功",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: fmt.Sprintf("%s", err.Error()),
		}.String())
	}
}

// @Title GetEdit
// @Summary 查看是否可以编辑指标
// @Description 查看是否可以编辑指标
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getEdit [get]
func (this *IndicatorController) GetEdit() {
	userName, _, _ := this.Ctx.Request.BasicAuth()
	method := "POST"
	checkUrl := fmt.Sprintf("%s/%s_%s", controllers.RootNameSpace, IndicatorNameSpace, "authority")
	models.Logger.Info(checkUrl)
	flag := models.Enforcer.Enforce(userName, checkUrl, method)
	res, _ := json.Marshal(map[string]bool{
		"readOnly": !flag,
	})
	this.Ctx.WriteString(string(res))
}

// @Title GetDictionary
// @Summary 获取数据字典
// @Description 获取数据字典
// @router /dictionary/getAll [get]
func (this *IndicatorController) GetDictionary() {
	index := indicator.Indicator{}
	table := index.GetDictionaryItems("", "", "")
	bytes, _ := json.Marshal(table)
	this.Ctx.WriteString(string(bytes))
}

// @Title SearchDictionary
// @Summary 根据关键词和业务标签搜索数据字典
// @Description 获取数据字典
// @router /dictionary/search [post]
func (this *IndicatorController) SearchDictionary() {
	attr := map[string]string{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	searchWord := attr["searchWord"]
	businessLabel := attr["businessLabel"]
	assetType := attr["assetType"]
	index := indicator.Indicator{}
	table := index.GetDictionaryItems(searchWord, businessLabel, assetType)
	bytes, _ := json.Marshal(table)
	this.Ctx.WriteString(string(bytes))
}

// @Title GetDictionaryById
// @Summary 获取某个资产的详细信息
// @Description 获取某个资产的详细信息
// @router /dictionary/searchById [post]
func (this *IndicatorController) GetDictionaryById() {
	attr := map[string]string{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	originId, _ := strconv.ParseInt(attr["originId"], 10, 64)
	assetType := attr["assetType"]
	o := orm.NewOrm()
	if (assetType == "predefine") || (assetType == "combine") {
		index := indicator.Indicator{}
		_ = o.Raw("SELECT * FROM indicator WHERE id = ?", originId).QueryRow(&index)
		bytes, _ := json.Marshal(index)
		this.Ctx.WriteString(string(bytes))
	} else if assetType == "dimension" {
		dimension := indicator.DimensionView{}
		sql := `
SELECT t1.*, 
	t2.table_name AS primary_table_name,
	t2.database_name AS primary_database_name,
	t3.datasource_name AS primary_table_data_source,
	t3.conn_type AS data_source_type,
	t3.dsid AS data_source
FROM dimension t1 
LEFT JOIN 
	table_meta_data t2
ON t1.primary_table_id = t2.tid
LEFT JOIN 
	connection_manager t3
ON t2.dsid = t3.dsid
WHERE t1.id = ?
`
		_ = o.Raw(sql, originId).QueryRow(&dimension)
		bytes, _ := json.Marshal(dimension)
		this.Ctx.WriteString(string(bytes))
	}
}

// @Title CheckIndicatorDuplicate
// @Summary 判断指标是否重复了
// @Description 判断指标是否重复了
// @router /dictionary/checkIndicatorDuplicate [post]
func (this *IndicatorController) CheckIndicatorDuplicate() {
	attr := map[string]string{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	name := attr["name"]
	fieldName := attr["fieldName"]
	id, _ := strconv.ParseInt(attr["id"], 10, 64)
	o := orm.NewOrm()
	indicators := []indicator.Indicator{}
	if id == 0 {
		_, _ = o.Raw("SELECT * FROM indicator WHERE name = ? or field_name = ?", name, fieldName).QueryRows(&indicators)
	} else {
		_, _ = o.Raw("SELECT * FROM indicator WHERE (name = ? or field_name = ?) and id != ?", name, fieldName, id).QueryRows(&indicators)
	}

	bytes, _ := json.Marshal(indicators)
	this.Ctx.WriteString(string(bytes))
}

// @Title 增加修饰词
// @Summary 增加修饰词
// @Description 增加修饰词
// @router /qualifier/add [post]
func (this *IndicatorController) AddQualifier() {
	qualifier := indicator.Qualifier{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &qualifier)
	if err := qualifier.Add(); err != nil {
		models.Logger.Error("%v", err)
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: fmt.Sprintf("%s", err.Error()),
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "添加成功",
		}.String())
	}
}

// @Title 获取时间修饰词
// @Summary 获取时间修饰词
// @Description 获取时间修饰词
// @router /qualifier/time/get [get]
func (this *IndicatorController) GetAllTimeQualifier() {
	timeQualifiers := []indicator.Qualifier{}
	o := orm.NewOrm()
	_, _ = o.Raw("SELECT * FROM qualifier WHERE qualifier_type = '时间周期'").QueryRows(&timeQualifiers)
	bytes,_ := json.Marshal(timeQualifiers)
	this.Ctx.WriteString(string(bytes))
}