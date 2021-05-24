package models

import (
	pb "AStar/protobuf"
	"AStar/utils"
	"fmt"
	"strconv"
	"strings"
)

type SQLMaker struct {
	Db    string
	Fact  string
	Query pb.Query
}

// 使用左连接JOIN拼接SQL语句
// 返回值：sql字符串，查询所用表集合
func (this *SQLMaker) MkSQLOfLeftJoinByProto(sqlModel SqlModel) (string, []string) {

	var tableList []string

	for _, table := range this.getTableList() {
		tableList = append(tableList, this.Db+"."+table)
	}

	selectList := this.getSelectList()
	leftJoinList := this.getLeftJoinList(sqlModel)
	whereFilterList := this.getWhereFilterList()
	groupByList := this.getGroupByList()

	sqlElements := []string{
		"SELECT",
		strings.Join(selectList, ", "),
		"FROM",
		strings.Join(leftJoinList, " "),
		"WHERE",
		strings.Join(whereFilterList, " AND "),
	}

	if len(groupByList) > 0 {
		sqlElements = append(sqlElements, "GROUP BY")
		sqlElements = append(sqlElements, strings.Join(groupByList, ", "))
	}

	if len(sqlModel.OrderCols) > 0 {
		sqlElements = append(sqlElements, fmt.Sprintf("ORDER BY %s DESC", strings.Join(sqlModel.OrderCols, ",")))
	}

	if sqlModel.LimitNum > 0 {
		sqlElements = append(sqlElements, fmt.Sprintf("LIMIT %d", sqlModel.LimitNum))
	}

	sql := strings.Join(sqlElements, " ")
	sql = pv2ToAdpv(sql)

	return sql, tableList

}

func (this *SQLMaker) getTableList() []string {
	itemList := this.Query.Items
	// 利用map结构实现set的效果，存放不重复的table name
	tableMap := make(map[string]string)

	for _, input := range itemList {
		for _, fullName := range input.Names {
			tbName := strings.Split(fullName, ".")[0]
			tableMap[tbName] = tbName
		}
	}

	// 将tableMap的Key转换成数组
	var tableList []string
	for _, fullName := range tableMap {
		tableList = append(tableList, fullName)
	}

	return tableList
}

func pv2ToAdpv(sql string) string {

	if strings.Contains(sql, "dw_fact_ws_adpv") {
		sql = strings.Replace(sql, "dw_fact_ws_adpv.pv2", "dw_fact_ws_adpv.ad_pv", -1)
	}

	return sql
}

// getLeftJoinList得到左连接字句的列表
// 形如{fact, LEFT JOIN table1 on condition1, LEFT JOIN table2 on condition2...}
func (this *SQLMaker) getLeftJoinList(wsSQL SqlModel) []string {

	tableList := this.getTableList()

	var joinList []string

	// 添加事实表
	joinList = append(joinList, this.Db+"."+this.Fact)

	// 添加其它表并组成JOIN...ON...字句
	for _, table := range tableList {
		if table != this.Fact {
			joinList = append(joinList, " JOIN "+this.Db+"."+table)
			if relationship, ok := wsSQL.GetFactDimConnMap(this.Fact)[table]; ok {
				joinList = append(joinList, "ON "+relationship)
			}
		}
	}
	return joinList
}

// getSelectList得到查询中所有select的表
// 从Query的input中取出相关的表，如果有的话，依次添加udf，别名
// 该方法默认添加一个ad_pv的聚合值（待后续优化逻辑） // 2018.8.17 已优化
// 形如{table1, udf1(table2), table3 AS res3, udf2(table4, table5) AS res4...}
func (this *SQLMaker) getSelectList() []string {

	var selectList []string

	itemList := this.Query.Items
	for _, item := range itemList {

		if item.Usage == pb.QueryItemUsage_Select || item.Usage == pb.QueryItemUsage_Output {
			// 一个input元素里，所有的input的字段名要用逗号拼接在一起
			inputElement := strings.Join(item.Names, ", ")

			// 如果这个input带有udf，则添加udf
			if item.Func != "" {
				if strings.ToLower(item.Func) == "sum" {
					inputElement = fmt.Sprintf("CAST(%s as bigint)", inputElement)
				}
				// 方法中不包含 %s
				if !strings.Contains(item.Func, "%s") {
					inputElement = item.Func + "(" + inputElement + ")"
				} else {
					inputElement = fmt.Sprintf(item.Func, inputElement)
				}
			}
			// 如果这个input带有别名，则添加别名
			if item.Alias != "" {
				inputElement = inputElement + " AS " + item.Alias
			}

			selectList = append(selectList, inputElement)
		}
	}

	return selectList
}

// getGroupByList得到查询中应该被group by的表
// 从Query的input中取出字段名
// 形如{tb1.field1, tb2.field2....}
func (this *SQLMaker) getGroupByList() []string {

	itemList := this.Query.Items

	// 利用map结构实现set的效果，存放所有的table name
	groupByMap := make(map[string]string)

	// 从input中取出的full name，groupByMap
	for _, item := range itemList {
		if item.Usage == pb.QueryItemUsage_Select {
			if (item.Alias == "`季度`") || (item.Alias == "`月份`") {
				groupByElement := strings.Join(item.Names, ",")
				groupByElement = fmt.Sprintf(item.Func, groupByElement)
				groupByMap[groupByElement] = groupByElement
			} else {
				groupByElement := strings.Join(item.Names, ",")
				if item.Func != "" {
					groupByElement = item.Func + "(" + groupByElement + ")"
				}
				groupByMap[groupByElement] = groupByElement
			}
		}
	}

	// groupByMap
	var groupByList []string
	for _, groupBy := range groupByMap {
		groupByList = append(groupByList, groupBy)
	}

	return groupByList
}

// getWhereFilterList取出WHERE语句中的filter部分
// 从Query的filter中取出值进行拼接
// 支持的操作有：IN NI(NOT IN) GT LT GE LE BE(BETWEEN)
// 形如{fact.dt=20010101, tb1.field1 > 100...}
func (this *SQLMaker) getWhereFilterList() []string {

	itemList := this.Query.Items

	var whereList []string

	// where字句的第二部分是过滤语句，需要对操作类型进行相应的判断，转换成对应的where

	for _, item := range itemList {

		if item.Usage == pb.QueryItemUsage_Where {

			// 一个filter元素里，所有的filter的字段名要用逗号拼接在一起
			filterElement := strings.Join(item.Names, ", ")

			// 如果这个input带有udf，则添加udf
			if (item.Func != "") && (item.FilterOp != pb.FilterOperation_QW_AM) && (item.FilterOp != pb.FilterOperation_QW_FM) {
				filterElement = item.Func + "(" + filterElement + ")"
			}

			// 定义转换后的where子句
			var filterClause string

			switch item.FilterOp {
			case pb.FilterOperation_EQ:
				filterClause = filterElement + " = " + "'" + item.Values[0] + "'"
			case pb.FilterOperation_NE:
				filterClause = filterElement + " <> " + "'" + item.Values[0] + "'"
			case pb.FilterOperation_IN:
				filterClause = filterElement + " IN ('" + strings.Join(item.Values, "', '") + "')"
			case pb.FilterOperation_NI:
				filterClause = filterElement + " NOT IN ('" + strings.Join(item.Values, "', '") + "')"
			case pb.FilterOperation_GT:
				filterClause = filterElement + " > " + "'" + item.Values[0] + "'"
			case pb.FilterOperation_LT:
				filterClause = filterElement + " < " + "'" + item.Values[0] + "'"
			case pb.FilterOperation_GE:
				filterClause = filterElement + " >= " + "'" + item.Values[0] + "'"
			case pb.FilterOperation_LE:
				filterClause = filterElement + " <= " + "'" + item.Values[0] + "'"
			// {{.qwFunc}}(query_keyword, {{.matchType}})=1
			case pb.FilterOperation_QW_AM:
				filterElement = item.Func + "(" + filterElement + ",0)"
				filterClause = filterElement + " = 1 "
			case pb.FilterOperation_QW_FM:
				filterElement = item.Func + "(" + filterElement + ",1)"
				filterClause = filterElement + " = 1 "
			case pb.FilterOperation_BE: // 只支持date字段，对非dt使用BE操作符目前会报错

				begin := item.Values[0]
				end := item.Values[1]
				subDays, startDate, _ := utils.GetDiffDays(begin, end)

				var queryDayClauses []string
				for day := 0; day <= subDays; day++ {
					dateStr := utils.DateFormatWithoutLine(startDate.AddDate(0, 0, day))
					//queryDayClauses = append(queryDayClauses, this.Fact+".dt"+" = "+"'"+dateStr+"'")
					queryDayClauses = append(queryDayClauses, dateStr)
				}

				filterClause = this.Fact + ".dt in ('" + strings.Join(queryDayClauses, "','") + "')"
			case pb.FilterOperation_ST:
				var startWithClauses []string
				for _, str := range item.Values {
					strLen := strconv.Itoa(len(str))
					startWithClauses = append(startWithClauses, "substr("+filterElement+",0,"+strLen+")='"+str+"'")
				}
				filterClause = "(" + strings.Join(startWithClauses, " OR ") + ")"

			}
			whereList = append(whereList, filterClause)
		}
	}

	return whereList
}
