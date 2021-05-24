package indicator

import (
	"AStar/utils"
	"github.com/astaxie/beego/orm"
)

type Dimension struct {
	Id              int64  `orm:"pk;auto"`
	Name            string `orm:"size(50)"` // 维度名字
	AliasName       string `orm:"size(150)"`
	PrimaryTableId  int64  `orm:"int"`      // 维度主表的表Id
	PkField         string `orm:"size(50)"` // 维度主表的主键
	ValField        string // 维度值字段
	AssistInfo      string // 辅助信息字段
	BusinessLabel   string `orm:"size(255)"`  // 业务标签
	Description     string `orm:"type(text)"` // 维度的描述
	IsDerivative    bool   // 是否是衍生维度
	MainDimensionId int64  // 如果是衍生维度, 则有一个唯一的主维度关联
}

type DimensionView struct {
	Id                     int64
	Name                   string // 维度名字
	AliasName              string // 英文别名
	PrimaryTableId         int64  // 维度主表ID
	PrimaryTableName       string //  如下信息是冗余字段, 在PrimaryTableId中医科取到
	PrimaryDatabaseName    string // 维度主表库名
	PrimaryTableDataSource string // 维度主表的数据源名称
	PkField                string // 这个维度在主表中的字段名
	ValField               string // 维度值字段
	AssistInfo             string // 辅助信息字段
	Description            string // 维度描述
	BusinessLabel          string
	// form表单属性
	DataSourceType string
	DataSource     int64
	IsDerivative    bool   // 是否是衍生维度
	MainDimensionId int64  // 如果是衍生维度, 则有一个唯一的主维度关联
}

type DimensionFactRelation struct {
	Dfid       int64  `orm:"pk;auto"`
	DimId      int64  `orm:"int"`
	FactId     int64  `orm:"int"`
	ForeignKey string `orm:"size(50)"`
}

func (this *Dimension) GetAllDimension(isDerivative bool) utils.TableData {
	o := orm.NewOrm()
	views := []DimensionView{}
	sql := `
SELECT * FROM (
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
	ON t2.dsid = t3.dsid) p
WHERE is_derivative = ?
ORDER BY id
DESC
`
	_, _ = o.Raw(sql, isDerivative).QueryRows(&views)
	return this.transView(views)
}

func (this *Dimension) transView(views []DimensionView) utils.TableData {
	var rows []map[string]interface{}
	for _, view := range views {
		tagRow := utils.GetSnakeStrMapOfStruct(view)
		rows = append(rows, tagRow)
	}
	table := utils.TableData{Data: rows}
	table.GetSnakeTitles(DimensionView{})
	return table
}

func (this *Dimension) PostNew(view DimensionView) (int64, error) {
	this.Name = view.Name
	this.AliasName = view.AliasName
	this.PrimaryTableId = view.PrimaryTableId
	this.PkField = view.PkField
	this.BusinessLabel = view.BusinessLabel
	this.Description = view.Description
	this.ValField = view.ValField
	this.AssistInfo = view.AssistInfo
	this.MainDimensionId = view.MainDimensionId
	this.IsDerivative = view.IsDerivative
	if view.Id == 0 {
		_, err := orm.NewOrm().Insert(this)
		return this.Id, err
	} else {
		// id不为0, 则更新已存在的
		this.Id = view.Id
		_, err := orm.NewOrm().Update(this, "Name", "AliasName", "PrimaryTableId", "PkField",
			"BusinessLabel", "Description", "ValField", "AssistInfo","MainDimensionId")
		return this.Id, err
	}
}

func (this *DimensionFactRelation) Post() error {
	o := orm.NewOrm()
	var err error
	if this.Dfid == 0 {
		// 插入
		_, err = o.Insert(this)
	} else {
		// 更新
		_, err = o.Update(this, "DimId", "FactId", "ForeignKey")
	}
	if err != nil {
		utils.Logger.Error("%s", err.Error())
	}
	return err
}

type DimensionFactRelationView struct {
	Dfid         int64  `field:"dfid,hidden"`
	DimId        int64  `field:"dimId,hidden"`
	FactId       int64  `field:"factid,hidden"`
	FactDatabase string `field:"事实表库名"`
	FacTableName string `field:"事实表表名"`
	ForeignKey   string `field:"外键"`
}

func (this *DimensionFactRelation) GetAllFactTable(dimId int64) utils.TableData {
	o := orm.NewOrm()
	sql := `SELECT t1.dfid, t1.dim_id, t1.fact_id, t1.foreign_key,
	t2.database_name AS fact_database,t2.table_name AS fac_table_name
FROM 
(SELECT * FROM dimension_fact_relation WHERE dim_id=?) t1
LEFT JOIN 
table_meta_data t2 
ON 
t1.fact_id=t2.tid
`
	var relations []DimensionFactRelationView
	_, _ = o.Raw(sql, dimId).QueryRows(&relations)
	return this.transView(relations)
}

func (this *DimensionFactRelation) transView(views []DimensionFactRelationView) utils.TableData {
	var rows []map[string]interface{}
	for _, view := range views {
		tagRow := utils.GetSnakeStrMapOfStruct(view)
		rows = append(rows, tagRow)
	}
	table := utils.TableData{Data: rows}
	table.GetSnakeTitles(DimensionFactRelationView{})
	return table
}

func (this *DimensionFactRelation) DeleteById(dfId int64) error {
	sql := "DELETE FROM dimension_fact_relation WHERE dfid = ?"
	_, err := orm.NewOrm().Raw(sql, dfId).Exec()
	return err
}
