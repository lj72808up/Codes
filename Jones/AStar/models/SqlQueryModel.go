package models

import (
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type SqlQueryModel struct {
	SqlHint      map[string][]string
	FunctionUDFs []string
}

type HintModel struct {
	TableName string
	FieldList string
}

// 获取元数据下的所有提示
func (this *SqlQueryModel) GetSqlHint(dsType string, dsId int64) {
	o := orm.NewOrm()
	var hints []HintModel
	if dsType == "mysql" {
		sql := `
select
  concat(t4.database_name,'.',t4.table_name) as table_name, t3.field_list
from
  (SELECT * FROM table_meta_data WHERE table_type=? AND dsid=?) t4
  left join (
    select
      t1.table_id, GROUP_CONCAT(t2.field_name) as field_list
    from
      (
        select
          table_id, max(create_time)
        from
          version_meta_data
        group by
          table_id
      ) t1
      left join field_meta_data t2 on t1.table_id = t2.table_id
    group by
      t1.table_id
  ) t3 on t3.table_id = t4.tid`
		_, _ = o.Raw(sql, dsType, dsId).QueryRows(&hints)

	} else {
		sql := `
select
  concat(t4.database_name,'.',t4.table_name) as table_name, t3.field_list
from
  (SELECT * FROM table_meta_data WHERE table_type=?) t4
  left join (
    select
      t1.table_id, GROUP_CONCAT(t2.field_name) as field_list
    from
      (
        select
          table_id, max(create_time)
        from
          version_meta_data
        group by
          table_id
      ) t1
      left join field_meta_data t2 on t1.table_id = t2.table_id
    group by
      t1.table_id
  ) t3 on t3.table_id = t4.tid`
		_, _ = o.Raw(sql, dsType).QueryRows(&hints)
	}

	hintMap := make(map[string][]string)

	for i, _ := range hints {
		tbHints := strings.Split(hints[i].FieldList, ",")
		tbName := hints[i].TableName
		if tbName != "" {
			hintMap[tbName] = tbHints
		}
	}
	fmt.Println(hints)
	this.SqlHint = hintMap
}

// 获取所有的udf函数
func (this *SqlQueryModel) GetSqlUDF() {
	metaURL := utils.SparkRoute(ServerInfo.SparkHost, ServerInfo.SparkPort, "get_sparkUDF")
	req := httplib.Get(metaURL).SetTimeout(time.Hour, time.Hour)
	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
	}
	var res []string
	_ = json.Unmarshal(bytes, &res)
	this.FunctionUDFs = res
}

func (this *SqlQueryModel) GetMetaTreeByPrivilege(dsType string, dsId int64, userName string) MetaTreeNode {
	o := orm.NewOrm()
	var res []TableMetaData
	roles := Enforcer.GetRolesForUser(userName)
	roleWildCard := utils.MakeRolePlaceHolder(roles)
	sql := ""
	if (dsType == "mysql") || (dsType == "clickhouse") {
		sql = "SELECT tid,table_name,database_name FROM table_meta_data WHERE table_type=? AND dsid=?"
	} else {
		sql = "SELECT tid,table_name,database_name FROM table_meta_data WHERE table_type=?"
	}
	sql = fmt.Sprintf(`SELECT t1.* FROM
(%s) t1
LEFT JOIN 
(select table_id from group_metatable_mapping where group_name in %s AND (read_flag =1 or exec_flag=1) AND valid=1 group by table_id) t2
ON t1.tid = t2.table_id where t2.table_id is not null`, sql, roleWildCard)
	if (dsType == "mysql") || (dsType == "clickhouse") {
		_, _ = o.Raw(sql, dsType, dsId, roles).QueryRows(&res)
	} else {
		_, _ = o.Raw(sql, dsType, roles).QueryRows(&res)
	}

	tablesMap := make(map[string][]MetaTreeNode)
	for i, _ := range res {
		dbName := res[i].DatabaseName
		tbName := res[i].TableName
		tid := res[i].Tid
		tablesMap[dbName] = append(tablesMap[dbName], MetaTreeNode{
			Label:     tbName,
			TableName: fmt.Sprintf("%s.%s", dbName, tbName),
			TableId:   tid,
		})
	}

	metaTree := MetaTreeNode{Label: "meta"}
	for db, thirdLevel := range tablesMap {
		node := MetaTreeNode{
			Label:    db,
			Children: thirdLevel,
		}
		metaTree.Children = append(metaTree.Children, node)
	}
	return metaTree
}

func (this *SqlQueryModel) GetExploreData(tableName string, count int64, dsType string, dsId int64) utils.VxeData {
	switch dsType {
	case "hive":
		return exploreHiveData(tableName, count)
	case "mysql":
		return exploreMysqlData(tableName, count, dsId)
	case "clickhouse":
		return exploreCKData(tableName, count, dsId)
	default:
		return utils.VxeData{}
	}
}

func exploreHiveData(tableName string, count int64) utils.VxeData {
	sql := fmt.Sprintf("SELECT * FROM %s LIMIT %d", tableName, count)
	fmt.Println(sql)
	url := utils.SparkRoute(ServerInfo.SparkHost, ServerInfo.SparkPort, "quickly_query")
	param, _ := json.Marshal(map[string]string{
		"sql": sql,
	})
	req := httplib.Post(url).SetTimeout(time.Hour, time.Hour).Body(param)
	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
	}
	var res []map[string]interface{}
	_ = json.Unmarshal(bytes, &res)

	cols := make([]utils.VxeColumn, 0)
	for k, _ := range res[0] {
		cols = append(cols, utils.VxeColumn{
			Field:   k,
			Title:   k,
			Visible: true,
		})
	}

	return utils.VxeData{
		TableColumn: cols,
		DataList:    res,
	}
}

func exploreMysqlData(tableName string, count int64, dsId int64) utils.VxeData {

	var conn ConnectionManager
	conn.GetById(dsId)

	aliasName := fmt.Sprintf("%s-%d", tableName, time.Now().Unix())
	ip := strings.Split(strings.Replace(conn.Url, "jdbc:mysql://", "", -1), "/")[0]
	dbName := strings.Split(tableName, ".")[0]
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s", conn.User, conn.Passwd, ip, dbName)
	encoderMapping := map[string]string{
		"cp1252": "latin1",
		"utf-8":  "utf8",
	}
	if conn.Encode != "" {
		encoder := encoderMapping[strings.ToLower(conn.Encode)]
		if encoder == "" {
			encoder = conn.Encode
			fmt.Printf("没有找到'%s'对应的tcp charset, 使用原版输入'%s'\n", encoder, encoder)
		}
		url = fmt.Sprintf("%s?charset=%s", url, encoder) //conn.Encode)
	}

	fmt.Printf("数据源: %s\n", url)

	connErr := orm.RegisterDataBase(aliasName, "mysql", url, 0, 1) // SetMaxIdleConns, SetMaxOpenConns
	if connErr != nil {
		panic(connErr)
	}
	o1 := orm.NewOrm()
	_ = o1.Using(aliasName)

	// 1. 执行查询
	sql := fmt.Sprintf("SELECT * FROM %s LIMIT %d" /*strings.Join(colNames,","),*/, tableName, count)
	var res []orm.Params
	_, _ = o1.Raw(sql).Values(&res)
	// 2. 封装vxe
	cols := make([]utils.VxeColumn, 0)
	for k, _ := range res[0] {
		cols = append(cols, utils.VxeColumn{
			Field:   k,
			Title:   k,
			Visible: true,
		})
	}

	var vxeRes []map[string]interface{}
	for i, _ := range res {
		map1 := res[i]
		vxeRes = append(vxeRes, map1)
	}

	return utils.VxeData{
		TableColumn: cols,
		DataList:    vxeRes,
	}
}

func exploreCKData(tableName string, count int64, dsId int64) utils.VxeData {
	var conn ConnectionManager
	conn.GetById(dsId)

	url := fmt.Sprintf("tcp://%s:%s?user=%s&password=%s&debug=false", conn.Host, conn.Port, conn.User, conn.Passwd)
	dbConn, err := sqlx.Open("clickhouse", url)
	if err != nil {
		panic(err)
	}

	sql := fmt.Sprintf("SELECT * FROM %s LIMIT %d", tableName, count)
	rows, err := dbConn.Query(sql)
	if err != nil {
		panic(err)
	}

	// 开始将row转化成vxeRes
	cols, _ := rows.Columns()
	fieldsCacheInOnwRow := make([]interface{}, len(cols))
	for index, _ := range fieldsCacheInOnwRow {
		var a interface{}
		fieldsCacheInOnwRow[index] = &a
	}

	var vxeRes []map[string]interface{}
	for rows.Next() {
		_ = rows.Scan(fieldsCacheInOnwRow...)
		item := make(map[string]interface{})
		for i, field := range fieldsCacheInOnwRow {
			item[cols[i]] = *(field.(*interface{}))
		}
		vxeRes = append(vxeRes, item)
	}
	_ = dbConn.Close()

	// 2. 封装vxe
	vxeCols := make([]utils.VxeColumn, 0)
	for _, col := range cols {
		vxeCols = append(vxeCols, utils.VxeColumn{
			Field:   col,
			Title:   col,
			Visible: true,
		})
	}

	return utils.VxeData{
		TableColumn: vxeCols,
		DataList:    vxeRes,
	}
}
