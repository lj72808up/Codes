package indicator

import (
	"AStar/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

const PREDEFINE = "predefine"
const COMBINE = "combine"

type Indicator struct {
	Id            int64  `orm:"pk;auto" field:"id,hidden" json:"id"`
	Name          string `orm:"size(200);unique"  field:"指标名" json:"name"`
	FieldName     string `orm:"size(200);unique"  field:"字段名称" json:"fieldName"`
	Meaning       string `orm:"type(text)"  field:"指标含义" json:"meaning"`
	IndexType     string `orm:"size(200)"  field:"指标类型,hidden" json:"indexType"` // PREDEFINE 或 COMBINE
	CreateTime    string `orm:"size(200)"  field:"创建时间" json:"createTime"`
	ModifyTime    string `orm:"size(200)"  field:"最后修改时间" json:"modifyTime"`
	Expression    string `orm:"type(text)"  field:"指标表达式,hidden" json:"expression"`
	Aggregation   string `orm:"size(200)"  field:"聚合方式" json:"aggregation"`
	BusinessLabel string `orm:"size(200)"  field:"业务标签" json:"businessLabel"` // 业务标签 (逗号分隔的标签字符列表)
	Qualifiers    string `orm:"size(255)" field:"修饰词" json:"qualifiers"`
	TimeQualifier int64  `orm:"size(255)" field:"时间修饰词" json:"timeQualifier"`
}

type IndicatorRelation struct {
	Id           int64 `orm:"pk;auto"`
	IndicatorId  int64 `orm:"type(int)"`
	DependencyId int64 `orm:"type(int)"`
}

// 多字段索引
func (u *IndicatorRelation) TableIndex() [][]string {
	return [][]string{
		{"IndicatorId"},
	}
}

func (this *Indicator) AddOne(indexView View) (int64, error) {
	now := utils.TimeFormat(time.Now())
	this.Name = indexView.Name
	this.FieldName = indexView.FieldName
	this.Aggregation = indexView.Aggregation
	this.Meaning = indexView.Meaning
	this.IndexType = indexView.IndexType
	this.CreateTime = now
	this.ModifyTime = now
	this.Expression = indexView.Expression
	this.BusinessLabel = indexView.TransBusinessLabels()
	this.Qualifiers = utils.MakeString(indexView.Qualifiers, ",")
	this.TimeQualifier = indexView.TimeQualifier

	o := orm.NewOrm()
	err := fmt.Errorf("")
	if _, err = o.Insert(this); err != nil {
		utils.Logger.Error("%v", err)
	}
	// beego insert动作会设置生成的pk
	uniqIds := map[int64]bool{}
	if len(indexView.Dependencies) > 0 {
		relations := []IndicatorRelation{}
		for _, denpendencyId := range indexView.Dependencies {
			if !uniqIds[denpendencyId] {
				relations = append(relations, IndicatorRelation{
					IndicatorId:  this.Id,
					DependencyId: denpendencyId,
				})
				uniqIds[denpendencyId] = true
			}

		}
		_, _ = o.Raw("DELETE FROM indicator_relation WHERE indicator_id = ?", this.Id).Exec()
		_, err = o.InsertMulti(100, &relations)
	}

	return this.Id, err
}

type View struct {
	Id             int64
	Name           string
	FieldName      string
	Meaning        string
	IndexType      string
	Expression     string
	Aggregation    string
	BusinessLabels []string
	Dependencies   []int64
	Qualifiers     []string
	TimeQualifier  int64
}

func (this *View) TransBusinessLabels() string {
	return strings.Join(this.BusinessLabels, ",")
}

func (this *Indicator) UpdateOne(view View) error {
	sql := `UPDATE indicator 
SET name=?, field_name = ?, meaning=?, index_type=?, business_label = ?, aggregation = ?,
	expression=?, qualifiers=?, time_qualifier = ?
WHERE id = ?`
	o := orm.NewOrm()
	_, err := o.Raw(sql, view.Name, view.FieldName, view.Meaning, view.IndexType, view.TransBusinessLabels(), view.Aggregation,
		view.Expression, utils.MakeString(view.Qualifiers, ","),view.TimeQualifier,
		view.Id).Exec()

	// beego insert动作会设置生成的pk
	uniqIds := map[int64]bool{}
	if len(view.Dependencies) > 0 {
		relations := []IndicatorRelation{}
		for _, denpendencyId := range view.Dependencies {
			if !uniqIds[denpendencyId] {
				relations = append(relations, IndicatorRelation{
					IndicatorId:  view.Id,
					DependencyId: denpendencyId,
				})
				uniqIds[denpendencyId] = true
			}
		}
		_, _ = o.Raw("DELETE FROM indicator_relation WHERE indicator_id = ?", view.Id).Exec()
		_, err = o.InsertMulti(100, &relations)
	}

	return err
}

func (this *Indicator) UpdateBatch(views []View) error {
	for _, view := range views {
		if err := this.UpdateOne(view); err != nil {
			return err
		}
	}
	return nil
}

func (this *Indicator) GetIndicator(indexType string) (utils.TableData, error) {
	o := orm.NewOrm()
	var indices []Indicator
	sql := "SELECT * FROM indicator WHERE index_type = ?"
	if _, err := o.Raw(sql, indexType).QueryRows(&indices); err == nil {
		return this.transTableView(indices), nil
	} else {
		utils.Logger.Error("%v", err)
		return utils.TableData{}, err
	}
}

func (this *Indicator) GetCombineIndicator() (utils.TableData, error) {
	return this.GetIndicator(COMBINE)
}

func (this *Indicator) GetPreDefineIndicator() (utils.TableData, error) {
	return this.GetIndicator(PREDEFINE)
}

func (this *Indicator) transTableView(tags []Indicator) utils.TableData {
	var rows []map[string]interface{}
	for _, tag := range tags {
		tagRow := utils.GetSnakeStrMapOfStruct(tag)
		rows = append(rows, tagRow)
	}
	table := utils.TableData{Data: rows}
	table.GetSnakeTitles(Indicator{})
	return table
}

func GetDependencyIds(indexId int) []int {
	o := orm.NewOrm()
	ids := make([]int, 0)
	_, _ = o.Raw("SELECT dependency_id from indicator_relation where indicator_id = ? group by dependency_id", indexId).QueryRows(&ids)
	return ids
}

type IndicatorFactRelation struct {
	Ifid            int64  `orm:"pk;auto" field:"id,hidden"`
	IndicatorId     int64  `field:"indicatorId"`
	FactId          int64  `field:"factId"`
	AssociatedField string `orm:"size(200)" field:"associatedField"`
}

type IFRelationView struct {
	Ifid            int64
	IndicatorId     int64
	FactId          int64
	FactName        string
	DatabaseName    string
	FactType        string
	AssociatedField string
}

func (this *IndicatorFactRelation) GetFactRelations(indicatorId int) utils.TableData {
	views := []IFRelationView{}
	sql := `
SELECT t1.*, t2.table_name as fact_name, t2.database_name as database_name, t2.table_type as fact_type FROM indicator_fact_relation t1
LEFT JOIN 
table_meta_data t2
ON t1.fact_id = t2.tid
WHERE t1.indicator_id = ?
`
	_, _ = orm.NewOrm().Raw(sql, indicatorId).QueryRows(&views)
	return this.transRelationView(views)
}

func (this *IndicatorFactRelation) transRelationView(views []IFRelationView) utils.TableData {
	var rows []map[string]interface{}
	for _, view := range views {
		tagRow := utils.GetSnakeStrMapOfStruct(view)
		rows = append(rows, tagRow)
	}
	table := utils.TableData{Data: rows}
	table.GetSnakeTitles(IFRelationView{})
	return table
}

func (this *IndicatorFactRelation) Add() error {
	o := orm.NewOrm()
	_, err := o.Insert(this)
	if err != nil {
		utils.Logger.Error("%s", err.Error())
	}
	return err
}

func (this *Indicator) DelteById(id int64) error {
	o := orm.NewOrm()
	_, err := o.Raw("DELETE FROM indicator WHERE id = ?", id).Exec()
	if err != nil {
		utils.Logger.Error("%s", err.Error())
	}
	return err
}
