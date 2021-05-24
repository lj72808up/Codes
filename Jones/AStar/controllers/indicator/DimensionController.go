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

type DimensionController struct {
	controllers.BaseController
}

// @Title GetDimensions
// @Summary 获取维度
// @Description 获取维度
// @router /dimension/getAll [get]
func (this *IndicatorController) GetDimensions() {
	dimension := indicator.Dimension{}
	isDerivative,_ := strconv.ParseBool(this.Ctx.Input.Query("isDerivative"))
	table := dimension.GetAllDimension(isDerivative)
	bytes, _ := json.Marshal(table)
	this.Ctx.WriteString(string(bytes))
}

// @Title GetDimensions
// @Summary 获取维度
// @Description 获取维度
// @router /dimension/postNew [post]
func (this *IndicatorController) NewDimensions() {
	view := indicator.DimensionView{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &view)
	dimension := indicator.Dimension{}
	id, err := dimension.PostNew(view)
	if err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: fmt.Sprintf("%d", id),
		}.String())
	} else {
		models.Logger.Error("%s", err.Error())
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	}
}

// @Title PostDimFactRelation
// @Summary 上传或更新事实表和维度白的关联
// @Description 上传或更新事实表和维度白的关联
// @router /dimension/postDimFactRelation [post]
func (this *IndicatorController) PostDimFactRelation() {
	relation := indicator.DimensionFactRelation{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &relation)
	err := relation.Post()
	if err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "操作成功",
		}.String())
	} else {
		models.Logger.Error("%s", err.Error())
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	}
}

// @Title GetDimFactRelations
// @Summary 获取维度表关联的事实表
// @Description 获取维度表关联的事实表
// @router /dimension/getDimFactRelations/:dimId [get]
func (this *IndicatorController) GetDimFactRelations() {
	dimId, _ := strconv.ParseInt(this.Ctx.Input.Params()[":dimId"], 10, 64)
	dim := indicator.DimensionFactRelation{}
	table := dim.GetAllFactTable(dimId)
	bytes, _ := json.Marshal(table)
	this.Ctx.WriteString(string(bytes))
}

// @Title DelDimFactRelation
// @Summary 删除事实表的关系
// @Description 删除事实表的关系
// @router /dimension/delDimFactRelation [post]
func (this *IndicatorController) DelDimFactRelation() {
	relation := indicator.DimensionFactRelation{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &relation)
	err := relation.DeleteById(relation.Dfid)
	if err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "删除成功",
		}.String())
	} else {
		models.Logger.Error("%s", err.Error())
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	}
}

// @Title GetAllMainDimension
// @Summary 获取所有主维度
// @Description 获取所有主维度
// @router /dimension/getAllMainDimension [get]
func (this *IndicatorController) GetAllMainDimension() {
	dimensions := []indicator.Dimension{}
	o := orm.NewOrm()
	_, _ = o.Raw("SELECT * FROM dimension WHERE is_derivative = ?", false).QueryRows(&dimensions)
	bytes, _ := json.Marshal(dimensions)
	this.Ctx.WriteString(string(bytes))
}