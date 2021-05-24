// SqlModel.go定义了如下内容：
//    1. 用于存储组装Sql的结构体
//    2. 组装sql过程中需要用到的方法
package models

import (
	"fmt"
	"strings"
	"sort"
	"github.com/astaxie/beego/config"
	pb "AStar/protobuf"
	"errors"
	"AStar/utils"
)

//
// SqlModel是储存Sql的最上层模型，SqlModel将查询Sql分成四个部分：
//
//    1. Db 代表查询应用的数据库
//
//    2. FactTbs 代表查询所用到的所有事实表（通常是FACT表，也可能是ODS表）
//       map 映射关系: name -> SqlFactTableModel
//       将查询所用的表定义成数组保存多个值的原因是，某一种查询可能需要多个事实表的支持
//       具体最终的查询要使用哪个事实表，根据所选择的具体查询维度来决定
//       也可能对于某种冲突的查询，不存在合适的事实表
//       由于事实表存在不确定性，在文字量中使用到的事实表均使用_fact来标注
//       在具体使用时需要对_fact进行替换
//
//    3. DimTbs 代表查询的所用到的所有维度表
//       map 映射关系: name -> SqlDimTableModel
//
//    4. Fields 代表查询所用到的所有字段。这些字段可能来自事实表，也可能来自维度表
//       map 映射关系: key -> SqlFieldModel
//
// SqlModel的设计原因及考虑的点：
//    1. 完整的查询需要多张事实表的支持，根据查询的内容不同，使用的事实表也不同
//    2. 每一次查询仅使用一张事实表，在某一种查询被多个事实表响应的时候，使用优先级最高的事实表
//    3. 事实表的外键和维度表的主键的名称保持一致
//    4. 所有事实表与维度的表的连接均使用Join
//

type SqlModel struct {
	Db        string
	FactTbs   map[string]SqlFactTableModel
	DimTbs    map[string]SqlDimTableModel
	Fields    map[string]SqlFieldModel
	OrderCols []string
	LimitNum  int64
}

//
// SqlFactTableModel在Sql模型中存储事实表，SqlFactTable将事实表分成五个部分：
//
//    1. Name 代表事实表的表名
//
//    2. Fields 代表事实表所包含的所有字段的名称
//       该字段名称应该与Fields Model中定义的Key名称应该相同，详见SqlFieldModel
//
//    3. Priority 代表该事实表的使用优先级
//       在多个事实表都能完成同一种查询的时候，优先使用优先级高的事实表进行查询
//
//    4. Comment 存储该事实表的中文名称
//
//    5. Kylin 记录该事实表在Kylin中的表名
//
type SqlFactTableModel struct {
	Name     string
	Fields   []string
	Priority int
	Comment  string
	Kylin    string
}

//
// SqlDimTableModel在Sql模型中存储维度表，SqlDimTable将维度表分成三个部分：
//
//    1. Name 代表维度表的表名
//
//    2. Statement 代表该维度表与事实表做JOIN时，ON关键字后面的条件语句
//
//    3. Comment 存储该维度表的中文名称
//
//
type SqlDimTableModel struct {
	Name      string
	Statement string
	Comment   string
}

//
// SqlFieldModel在Sql模型中存储某个字段，这里的字段是广义的“字段”，是sql语句当中逗号隔断的一个查询单元
// 如Func(field1, field2) AS RESULT，被视作一个“字段”。SqlField将字段分为六个部分：
//
//    1. key 代表字段的唯一标识
//       注，在同一个sql模型当中，所有字段的key应该是不相同的
//
//    2. alias 代表字段的别名
//        注，别名用在拼接sql的时候，作为AS 后面的字段别名出现
//        通常别名alias和key可以设置为相同的值，因为在一个sql语句中通常也不应该出现别名相同的两个字段
//
//    3. Names 代表一个字段所包含的名称。(逗号分隔单个字符串值)
//         注，因为udf可能会涉及多个字段，如Func(field1, field2)，这里的Names就是"field1,field2"
//
//    4. FuncName 代表要应用的方法的名称，没有时设为空字符串""
//
//    5. Engine 代表这个字段适用的查询引擎。
//         注，查询引擎比方说hive，spark，kylin，是与字段相关的特性。一个sql语句最终适用哪个引擎来执行，
//         实际上是由它包含的全部字段来决定的
//
//    6. Comment 代表这个字段的中文名称
//
type SqlFieldModel struct {
	Key      string
	Alias    string
	Names    []string
	FuncName string
	Engine   int
	Comment  string
}

// SqlFactList是为了存储有序的Fact表设计的线性结构，实现下面的方法以实现自定义排序
type SqlFactList []SqlFactTableModel

func (Fact SqlFactList) Len() int {
	return len(Fact)
}
func (Fact SqlFactList) Less(i, j int) bool {
	return Fact[i].Priority > Fact[j].Priority
}
func (Fact SqlFactList) Swap(i, j int) {
	Fact[i], Fact[j] = Fact[j], Fact[i]
}

// InitSQL方法对sqlModel进行初始化
// 参数：
//    保存sql配置的json文件名
// 返回值：
//    sqlModel 对象实例
//    error 错误
func (sqlModel SqlModel) InitSQL(jsonFile string) (SqlModel, error) {

	sqlModel.FactTbs = make(map[string]SqlFactTableModel)
	sqlModel.DimTbs = make(map[string]SqlDimTableModel)
	sqlModel.Fields = make(map[string]SqlFieldModel)

	// 解json
	conf, err := config.NewConfig("json", jsonFile)

	if err != nil {
		return sqlModel, err
	}

	// 填充sqlmodel.Db
	sqlModel.Db = conf.String("db")

	// 填充SqlModel.FactTbs
	factsRaw, _ := conf.DIY("fact_tables")

	for _, factRaw := range factsRaw.([]interface{}) {

		fact := factRaw.(map[string]interface{})

		var fields []string
		for _, field := range fact["fields"].([]interface{}) { // 所循环的是原始的未作类型断言的fact[fields]里的值
			fields = append(fields, field.(string))
		}

		var name = fact["name"].(string)
		sqlModel.FactTbs[name] = SqlFactTableModel{
			Name:     name,
			Fields:   fields,
			Priority: int(fact["priority"].(float64)),
			Comment:  fact["comment"].(string),
			Kylin:    fact["kylin"].(string),
		}
	}

	// 填充SqlModel.DimTbs

	dimsRaw, _ := conf.DIY("dim_tables")

	for _, dimRaw := range dimsRaw.([]interface{}) {

		dim := dimRaw.(map[string]interface{})

		var name = dim["name"].(string)

		sqlModel.DimTbs[name] = SqlDimTableModel{
			Name:      name,
			Statement: dim["conn_stat"].(string),
			Comment:   dim["comment"].(string),
		}

	}

	// 填充SqlModel.Fields
	fieldsRaw, _ := conf.DIY("fields")
	for _, fieldRaw := range fieldsRaw.([]interface{}) {
		field := fieldRaw.(map[string]interface{})

		var key = field["key"].(string)

		sqlModel.Fields[key] = SqlFieldModel{
			Key:      key,
			Alias:    field["alias"].(string),
			Names:    strings.Split(field["names"].(string), ";"),
			FuncName: field["func"].(string),
			Engine:   int(field["engine"].(float64)),
			Comment:  field["comment"].(string),
		}
	}
	return sqlModel, nil
}

// 创建一个按照优先级降级排序的fact表的list
// 该list中存储的元素为SqlFactTableModel的对象
func (sqlModel SqlModel) InitFactsPriorityOrderedList() []SqlFactTableModel {

	var facts SqlFactList
	for _, fact := range sqlModel.FactTbs {
		facts = append(facts, fact)
	}
	sort.Sort(facts)
	return facts
}

// 根据传入的fact表名，创建一个包含 维度表表名-> 维度表join语句 映射的map
// 参数：
//    fact 事实表名
//    在sql的json格式的配置文件当中，事实表使用的是_fact代替，因为要使用哪个事实表
//    是根据field的选择情况动态确定的。所以在确定了事实表以后，再能调用这个方法
//    并且需要对连接语句中的_fact进行替换，替换成为真正的事实表
func (sqlModel SqlModel) GetFactDimConnMap(fact string) map[string]string {
	connectionMap := make(map[string]string)
	for _, dim := range sqlModel.DimTbs {
		connectionMap[dim.Name] = strings.Replace(dim.Statement, "_fact", fact, -1)
	}

	return connectionMap
}

// 根据传入的事实表，返回一个包含key到protobuf.QueryItem映射的map
func (sqlModel SqlModel) GetKeyQueryMap(fact string) (map[string]pb.QueryItem,map[string]string) {
	keyQueryItemMap := make(map[string]pb.QueryItem)
	commentMap := make(map[string]string)
	for _, field := range sqlModel.Fields {
		var engine pb.QueryEngine
		if field.Engine == 1 {
			engine = pb.QueryEngine_spark
		} else if field.Engine == 2 {
			engine = pb.QueryEngine_kylin
		} else {
			engine = pb.QueryEngine_notspecified
		}
		var names []string
		for _, name := range field.Names {
			newName := strings.Replace(name, "_fact", fact, -1)
			names = append(names, newName)
		}
		keyQueryItemMap[field.Key] = pb.QueryItem{
			Alias:  field.Alias,
			Names:  names,
			Func:   field.FuncName,
			Engine: engine,
		}
		commentMap[field.Key] = field.Comment
	}
	return keyQueryItemMap, commentMap
}

// 根据选择的字段的key的名称，判断使用哪个事实表
func (sqlModel SqlModel) DecideFactTb(queryArr []string) (SqlFactTableModel, error) {

	var ansFact SqlFactTableModel

	// 没有维度传入的时候，不进行

	if len(queryArr) == 0 {
		return ansFact, errors.New("请进行页面维度的勾选")
	}

	var isValid = true
	for _, fact := range WsFactProiOrdList {
		isValid = true
		for _, query := range queryArr {
			factFields := fact.Fields
			if !utils.HasValue(factFields, query) {
				fmt.Printf("%s 里没有 %s\n", fact.Name, query)
				isValid = false
			}
		}

		if isValid {
			ansFact = fact
			break
		}
	}

	if isValid {
		return ansFact, nil
	}
	return ansFact, errors.New("页面维度勾选不正确。例如勾选了客户维度的PV1或PV3。")
}

// 根据选择的字段的key的名称，判断使用哪个事实表
func (sqlModel SqlModel) DecideFactTbPC(queryArr []string) (SqlFactTableModel, error) {

	var ansFact SqlFactTableModel

	// 没有维度传入的时候，不进行

	if len(queryArr) == 0 {
		return ansFact, errors.New("请进行页面维度的勾选")
	}

	var isValid = true
	for _, fact := range PcFactProiOrdList {
		isValid = true
		for _, query := range queryArr {
			factFields := fact.Fields
			if !utils.HasValue(factFields, query) {
				fmt.Printf("%s 里没有 %s\n", fact.Name, query)
				isValid = false
			}
		}

		if isValid {
			ansFact = fact
			break
		}
	}

	if isValid {
		return ansFact, nil
	}
	return ansFact, errors.New("页面维度勾选不正确。例如勾选了客户维度的PV1或PV3。")
}
