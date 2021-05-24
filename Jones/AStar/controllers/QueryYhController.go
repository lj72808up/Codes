package controllers

import (
	"AStar/models"
	"time"
	"github.com/golang/protobuf/proto"
	"encoding/json"
	pb "AStar/protobuf"
	"AStar/utils"
)

type QueryYhController struct {
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
// @router /jone/combination_yh [post]
func (this *QueryYhController) QuerySubmit() {

	// 单独取一个sqlModel过来，方便修改
	sqlModel := models.YhSQLModel

	// 解析前端发来的json得到后面需要的值

	var queryJSON = models.FrontEndQueryJson{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &queryJSON)
	if err != nil {
		models.Logger.Error(err.Error())
		this.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	var dimensions = queryJSON.Dimensions
	var filters = queryJSON.Filters
	var date = queryJSON.Date
	var outputs = queryJSON.Outputs
	var queryTaskName = queryJSON.Name

	// 添加瀚海渠道的判断 到 filters中
	filters = append(filters, models.FrontEndQueryFilterJson{
		Name:  "RET",
		Value: []string{"0"},
		Type:  "IN",
	})

	filters = append(filters, models.FrontEndQueryFilterJson{
		Name:  "SID",
		Value: []string{"6","7","205","207"},
		Type:  "ST",
	})

	// 事实表唯一，名字目前的名字是：
	// ods_ad_click
	// 省略判断事实表的逻辑
	fact := sqlModel.FactTbs["ods_ad_charge_click"]

	// 引擎唯一，但是保留引擎判断的逻辑
	// 保存所有遇到的engine，用map模拟set结构
	engineSet := make(map[pb.QueryEngine]int)

	var QueryItemsList []*pb.QueryItem
	var keyMapWithFact, _ = sqlModel.GetKeyQueryMap(fact.Name)

	// 这些值是用于给QueryInfo存储的
	var dimListForQueryInfo []string
	var outputListForQueryInfo []string
	var filterListForQueryInfo []string

	for _, item := range dimensions {
		if _, ok := keyMapWithFact[item]; ok {
			queryItem := keyMapWithFact[item]
			queryItem.Usage = pb.QueryItemUsage_Select
			QueryItemsList = append(QueryItemsList, &queryItem)

			dim := item
			dimListForQueryInfo = append(dimListForQueryInfo, dim)

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
			} else if filter.Type == "ST" {
				queryItem.FilterOp = pb.FilterOperation_ST
			}
			QueryItemsList = append(QueryItemsList, &queryItem)

			filterName := filter.Name
			filterListForQueryInfo = append(filterListForQueryInfo, filterName)

			_, isKey := engineSet[queryItem.Engine] // 添加engine
			if !isKey {
				engineSet[queryItem.Engine] = 1
			}
		}
	}
	// 添加时间
	dateQueryMap, _ := sqlModel.GetKeyQueryMap(fact.Name)
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

	////// 进行引擎的选择
	var queryEngine pb.QueryEngine
	_, isSpark := engineSet[pb.QueryEngine_spark]
	if isSpark {
		queryEngine = pb.QueryEngine_spark
	} else {
		queryEngine = pb.QueryEngine_kylin
	}

	if !models.KylinInfo.CheckKylin(fact.Kylin, dateFilter.Values[0], dateFilter.Values[1]) {
		queryEngine = pb.QueryEngine_spark
	}

	nowTime := time.Now()
	// 生成Md5结果，作为QueryInfo的QueryId
	Md5Res := utils.GetMd5ForQuery(sqlStr, this.userName, nowTime.String())

	// 生成QueryInfo
	queryInfo := pb.QueryInfo{
		QueryId:    Md5Res,                 // query id
		Sql:        []string{sqlStr},       // sql statement
		Session:    this.userName,          // session
		Engine:     queryEngine,            // 引擎
		Dimensions: dimListForQueryInfo,    // 维度表头
		Outputs:    outputListForQueryInfo, // 输出表头
		Filters:    filterListForQueryInfo, // 过滤字段
		Status:     pb.QueryStatus_Waiting, // 状态
		LastUpdate: nowTime.UnixNano() / 1e6,         // 最后更新时间戳
		CreateTime: nowTime.UnixNano() / 1e6,         // 创建时间戳
		Tables:     tables,                 // 查询所用表
		Type:       "yh",                   // 查询类别"银河搜索查询"
	} // 将QueryInfo入库MySql

	var queryJSONStr = string(this.Ctx.Input.RequestBody)
	models.QuerySubmitLog{}.InsertQuerySubmitLog(queryInfo, queryTaskName, queryJSONStr)
	// 编码QueryInfo
	queryBytes, _ := proto.Marshal(&queryInfo)
	// 得到MQ
	models.RabbitMQModel{}.Publish(queryBytes)
	// 给前端反应
	this.Data["json"] = queryInfo.QueryId
	this.ServeJSON()

}
