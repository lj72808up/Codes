package indicator

import (
	"AStar/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
)

type DictionaryItem struct {
	ItemName      string `field:"资产名称"`
	OriginId      int64  `field:"originId,hidden"`
	TypeName      string `field:"资产类型"`
	Description   string `field:"描述"`
	BusinessLabel string `field:"业务标签"`
}

// 按照查询词和业务标签搜索数据字典
func (this *Indicator) GetDictionaryItems(searchWord string, businessLabel string, assetType string) utils.TableData {
	// todo sql 加入 资产类型筛选
	o := orm.NewOrm()
	indexItems := []DictionaryItem{}
	var conditionStat string
	if searchWord != "" {
		conditionStat = fmt.Sprintf("WHERE (item_name LIKE '%%%s%%' OR description LIKE '%%%s%%')", searchWord, searchWord)
	}
	if businessLabel != "" {
		labelCondition := fmt.Sprintf("business_label LIKE '%%%s%%'", businessLabel)
		if conditionStat != "" {
			conditionStat = fmt.Sprintf("%s AND %s", conditionStat, labelCondition)
		} else {
			conditionStat = fmt.Sprintf("WHERE %s", labelCondition)
		}
	}

	// 原子指标
	predefineSql := `SELECT 
	item_name, origin_id, type_name, business_label, 
	group_concat(concat('事实表: ',TABLE_NAME)  SEPARATOR '\n') AS description
FROM (
	SELECT 
		t1.name AS item_name, 
		t1.id AS origin_id,
		'原子指标' AS type_name,
		t1.business_label,
		concat(t3.database_name,'.',t3.table_name) AS TABLE_NAME
	FROM 
	( SELECT * FROM indicator WHERE index_type='predefine') t1
	LEFT JOIN 
		indicator_fact_relation t2
	ON t1.id = t2.indicator_id
	LEFT JOIN 
		table_meta_data t3
	ON t2.fact_id = t3.tid
) t4
GROUP BY 
item_name, origin_id, type_name, business_label`
	combineSql := `SELECT id AS origin_id, 
	NAME AS item_name, 
	business_label, 
	group_concat(concat('依赖指标: ',dependency_name)  SEPARATOR '\n') AS description,
	'计算指标' AS type_name
FROM 
(
	SELECT t3.id, t3.NAME, t3.business_label, t3.dependency_id, t4.name AS dependency_name
	FROM (
		SELECT t1.id, t1.name, t1.business_label, t2.dependency_id FROM 
	( SELECT * FROM indicator WHERE index_type = 'combine' ) t1
	LEFT JOIN 
	indicator_relation t2
	ON t1.id = t2.indicator_id
) t3
LEFT JOIN 
	indicator t4
ON t3.dependency_id = t4.id
) t5
GROUP BY id, NAME, business_label`
	dimensionSql := `SELECT 
	t1.name AS item_name,
	t1.id AS origin_id,
	'维度' AS type_name,
	concat('维度主表: ',t2.database_name,'.',t2.table_name,'\n','主表数据源: ',t3.datasource_name, '\n','维度主键: ',t1.pk_field) AS description,
	t1.business_label
FROM dimension t1 
LEFT JOIN 
	table_meta_data t2
ON t1.primary_table_id = t2.tid
LEFT JOIN 
	connection_manager t3
ON t2.dsid = t3.dsid`
	selectStat := `SELECT item_name, origin_id, type_name, business_label,description FROM (
	%s
) p1
%s`
	sql1 := fmt.Sprintf(selectStat, predefineSql, conditionStat)
	sql2 := fmt.Sprintf(selectStat, combineSql, conditionStat)
	sql3 := fmt.Sprintf(selectStat, dimensionSql, conditionStat)
	sql := ""

	if assetType == "" {
		sql = fmt.Sprintf(`%s 
UNION 
%s 
UNION
%s`, sql1, sql2, sql3)
	} else if assetType == "predefine" {
		sql = sql1
	} else if assetType == "combine" {
		sql = sql2
	} else if assetType == "dimension" {
		sql = sql3
	} else {
		utils.Logger.Error("没有对应资产类型: %s", assetType)
	}

	_, _ = o.Raw(sql).QueryRows(&indexItems)
	return transView(indexItems)
}

func transView(views []DictionaryItem) utils.TableData {
	var rows []map[string]interface{}
	for _, view := range views {
		tagRow := utils.GetSnakeStrMapOfStruct(view)
		rows = append(rows, tagRow)
	}
	table := utils.TableData{Data: rows}
	table.GetSnakeTitles(DictionaryItem{})
	return table
}
