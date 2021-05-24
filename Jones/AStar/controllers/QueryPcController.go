package controllers

import (
	"AStar/models"
	pb "AStar/protobuf"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

type QueryPcController struct {
	BaseController
}

// @Title SubmitJonesQuery
// @Summary 提交Jones任务
// @Param para body models.FrontEndQueryJson true "任务信息"
// @Description 提交Jones任务
// @Success 200 string "queryId"
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /jone/combination_pc [post]
func (this *QueryPcController) QuerySubmit() {
	typeName := this.GetString("typeName", "PC搜索查询")
	var qwFileContent = this.GetString("qwFile", "")
	fmt.Println(qwFileContent)
	var queryStr = this.GetString("queryJSON", "") // queryJSON.Dimensions

	// 单独取一个sqlModel过来，方便修改
	// 取pcSqlModel
	sqlModel := models.PcSQLModel
	// 解析前端发来的json得到后面需要的值

	var queryJSON = models.FrontEndQueryJson{}
	err := json.Unmarshal([]byte(queryStr), &queryJSON)
	if err != nil {
		models.Logger.Error(err.Error() + "JSON解释失败")
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte(err.Error() + "JSON解释失败"))
		return
	}

	partition := ""
	if qwFileContent != "" {
		partition = LoadQWPartitionAfterTranscoding(qwFileContent, this.userName)
		fmt.Printf("query word partition:%s \n ", partition)
	}

	var dimensions = queryJSON.Dimensions
	var filters = queryJSON.Filters
	var date = queryJSON.Date
	var outputs = queryJSON.Outputs
	var queryTaskName = queryJSON.Name
	var orderCols = queryJSON.OrderCols
	limitNum,_ := strconv.ParseInt(queryJSON.LimitNum,10,64)

	sqlModel.LimitNum = limitNum
	sqlModel.OrderCols = orderCols

	selects := dimensions
	for _, value := range outputs {
		selects = append(selects, value)
	}

	// 保存所有遇到的engine，用map模拟set结构
	engineSet := make(map[pb.QueryEngine]int)

	// 取出所有用到的字段的名字（使用的是KEY概念）
	var queryUsingArr []string
	for _, item := range selects {
		queryUsingArr = append(queryUsingArr, item)
	}
	for _, filter := range filters {
		queryUsingArr = append(queryUsingArr, filter.Name)
	}

	// 得到要使用的事实表
	fact, err := sqlModel.DecideFactTbPC(queryUsingArr)
	if err != nil {
		// 所选择的指标和维度产生了冲突，不能得到有意义的结果，终止方法
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	var QueryItemsList []*pb.QueryItem
	keyMapWithFact, commentMap := sqlModel.GetKeyQueryMap(fact.Name)

	// 这些值是用于给QueryInfo存储的
	var dimListForQueryInfo []string
	var dimListMq []string
	var outputListForQueryInfo []string
	var outputMqList []string
	var filterListForQueryInfo []string

	for _, item := range dimensions {
		if _, ok := keyMapWithFact[item]; ok {
			queryItem := keyMapWithFact[item]
			queryItem.Usage = pb.QueryItemUsage_Select
			QueryItemsList = append(QueryItemsList, &queryItem)

			dim := item
			dimListForQueryInfo = append(dimListForQueryInfo, dim)
			dimListMq = append(dimListMq, commentMap[item])

			_, isKey := engineSet[queryItem.Engine] // 添加engine
			if !isKey {
				engineSet[queryItem.Engine] = 1
			}
		}
	}
	for _, item := range outputs {
		if _, ok := keyMapWithFact[item]; ok {
			queryItem := keyMapWithFact[item]
			queryItem.Usage = pb.QueryItemUsage_Output
			QueryItemsList = append(QueryItemsList, &queryItem)
			_, isKey := engineSet[queryItem.Engine]
			if !isKey {
				engineSet[queryItem.Engine] = 1
			}
			output := item
			outputListForQueryInfo = append(outputListForQueryInfo, output)
			outputMqList = append(outputMqList, commentMap[item])
		}
	}

	for _, filter := range filters {
		if _, ok := keyMapWithFact[filter.Name]; ok {
			queryItem := keyMapWithFact[filter.Name]
			queryItem.Usage = pb.QueryItemUsage_Where
			queryItem.Values = filter.Value
			if filter.Type == "IN" {
				if len(filter.Value) == 1 {
					queryItem.FilterOp = pb.FilterOperation_EQ
				} else {
					queryItem.FilterOp = pb.FilterOperation_IN
				}

			} else if filter.Type == "NOTIN" {
				queryItem.FilterOp = pb.FilterOperation_NI
			} else if filter.Type == "QW_AM" {
				queryItem.FilterOp = pb.FilterOperation_QW_AM
				queryItem.Func = partition
				queryItem.Values = append(queryItem.Values, partition)
			} else if filter.Type == "QW_FM" {
				queryItem.FilterOp = pb.FilterOperation_QW_FM
				queryItem.Func = partition
				queryItem.Values = append(queryItem.Values, partition)

			}

			QueryItemsList = append(QueryItemsList, &queryItem)

			//filterName := filter.Name
			//filterListForQueryInfo = append(filterListForQueryInfo, filterName)

			_, isKey := engineSet[queryItem.Engine] // 添加engine
			if !isKey {
				engineSet[queryItem.Engine] = 1
			}
		}
	}

	filterInMQ := make(map[string]interface{})
	filterInMQ["name"] = "queryWordPartition"
	filterInMQ["type"] = "QW_PART"
	filterInMQ["value"] = partition
	filterRes, _ := json.Marshal(filterInMQ)

	filterListForQueryInfo = append(filterListForQueryInfo, string(filterRes))

	// 添加时间
	dateQueryMap,_ := sqlModel.GetKeyQueryMap(fact.Name)
	dateFilter := pb.QueryItem{
		Names:    dateQueryMap["RQ"].Names,
		FilterOp: pb.FilterOperation_BE,
		Usage:    pb.QueryItemUsage_Where,
		Values:   []string{date[0], date[1]},
		Func:     "",
	}

	QueryItemsList = append(QueryItemsList, &dateFilter)

	qsQuery := pb.Query{Items: QueryItemsList}

	sqlMaker := models.SQLMaker{Db: sqlModel.Db, Fact: fact.Name, Query: qsQuery}

	sqlStr, tables := sqlMaker.MkSQLOfLeftJoinByProto(sqlModel)

	models.Logger.Debug(this.userName + "\t" + sqlStr)
	models.Logger.Info("%s",sqlStr)

	////// 进行引擎的选择
	var queryEngine pb.QueryEngine
	_, isSpark := engineSet[pb.QueryEngine_hive]
	if isSpark {
		queryEngine = pb.QueryEngine_hive
	} else {
		queryEngine = pb.QueryEngine_kylin
	}

/*	if !models.KylinInfo.CheckKylin(fact.Kylin, dateFilter.Values[0], dateFilter.Values[1]) {
		queryEngine = pb.QueryEngine_hive
	}*/

	nowTime := time.Now()
	// 生成Md5结果，作为QueryInfo的QueryId
	Md5Res := utils.GetMd5ForQuery(sqlStr, this.userName, nowTime.String())

	// 生成QueryInfo
	queryInfo := pb.QueryInfo{
		QueryId:    Md5Res,                   // query id
		Sql:        []string{sqlStr},         // sql statement
		Session:    this.userName,            // session
		Engine:     queryEngine,              // 引擎
		Dimensions: dimListForQueryInfo,      // 维度表头
		Outputs:    outputListForQueryInfo,   // 输出表头
		Filters:    filterListForQueryInfo,   // 过滤字段
		Status:     pb.QueryStatus_Waiting,   // 状态
		LastUpdate: nowTime.UnixNano() / 1e6, // 最后更新时间戳
		CreateTime: nowTime.UnixNano() / 1e6, // 创建时间戳
		Tables:     tables,                   // 查询所用表
		Type:       typeName,                 // 查询类别"PC搜索查询"
	} // 将QueryInfo入库MySql

	queryJson := make(map[string]interface{})
	queryJson["dimensions"] = queryJSON.Dimensions
	repairFilters := queryJSON.Filters
	for i, _ := range repairFilters {
		typeName := repairFilters[i].Type
		if (typeName == "QW_AM") || (typeName == "QW_FM") {
			repairFilters[i].Value = []string{partition}
		}
	}

	models.QuerySubmitLog{}.InsertQuerySubmitLog(queryInfo, queryTaskName, queryStr)

	// 得到MQ
	queryInfo.Dimensions = dimListMq
	queryInfo.Outputs = outputMqList
	queryBytes, _ := proto.Marshal(&queryInfo)
	models.RabbitMQModel{}.Publish(queryBytes)
	// 给前端反应
	this.Data["json"] = queryInfo.QueryId
	this.ServeJSON()

}
