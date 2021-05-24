package models

import (
	"AStar/utils"
	"container/list"
	"encoding/csv"
	"fmt"
	"github.com/astaxie/beego/orm"
	"io"
	"strconv"
	"strings"
	"time"
)

type FieldMetaData struct {
	Fid              int64  `orm:"pk;auto" vxe:"fid,hidden" json:"fid"`
	TableId          int64  `orm:"type(int)" vxe:"tableId,hidden" json:"tableId"`            // 所属表的id, 和TableMetaData的Tid对应
	FieldName        string `orm:"size(100)" vxe:"字段名称" json:"fieldName"`              // 字段名称
	FieldType        string `orm:"size(100)" vxe:"字段类型" json:"fieldType"`              // 字段类型
	Description      string `orm:"type(text)" vxe:"业务信息描述" json:"description"`         // 业务信息描述
	Version          string `orm:"default(50)" vxe:"版本" json:"version"`                // 字段的版本
	ModificationTime string `orm:"size(50)" vxe:"修改时间,hidden" json:"modificationTime"` // 修改时间
	IsPartitionField bool   `orm:"type(bool)" vxe:"分区字段,hidden" json:"isPartitionField"`     // 是否为分区字段
}

// 多字段索引
func (this *FieldMetaData) TableIndex() [][]string {
	return [][]string{
		[]string{"TableId", "Version", "IsPartitionField"},
	}
}

type VersionMetaData struct {
	Vid         int64  `orm:"pk;auto"`
	TableId     int64  `orm:"int"`       // 所属表的id, 和TableMetaData的Tid对应
	Version     string `orm:"size(100)"` // 字段的版本
	CreateTime  string `orm:"size(50)"`  // 创建时间
	Description string `orm:"size(255)"` // 版本描述
}

// 多字段索引
func (this *VersionMetaData) TableIndex() [][]string {
	return [][]string{
		[]string{"TableId", "CreateTime"},
	}
}

type RelationMetaData struct {
	Rid         int64  `orm:"pk;auto"`
	TableAId    int64  `orm:"int" field:"TableAId,drop"` // 表A的Id
	FieldAId    int64  `orm:"int" field:"FieldAId,drop"` // 表A的字段a
	VersionA    string `orm:"size(50)" field:"VersionA,drop"`
	TableBId    int64  `orm:"int" field:"表名"`   // 表B的id
	FieldBId    int64  `orm:"int" field:"字段名称"` // 表B的字段
	VersionB    string `orm:"size(50)" field:"版本"`
	Description string `orm:"size(50)" field:"关系描述"` // 描述信息
}

// view
type NewVersionData struct {
	TableId int64
	Version string
	Fields  []map[string]string
}

type TableRelation struct {
	Trid         int64 `orm:"pk;auto"`
	TableAId     int64 `orm:"int" field:"TableAId,drop"` // 表A的Id
	TableBId     int64 `orm:"int" field:"表名"`            // 表B的id
	ReferenceCnt int64 `orm:"int" field:"引用的字段个数,drop"`  // 表B的id
}

func (this *TableRelation) AddOne() {
	o := orm.NewOrm()
	var checkRelation TableRelation
	_ = o.Raw("SELECT * FROM table_relation WHERE table_a_id=? AND table_b_id=?", this.TableAId, this.TableBId).
		QueryRow(&checkRelation)
	if checkRelation.Trid == 0 {
		newRelation := TableRelation{
			TableAId:     this.TableAId,
			TableBId:     this.TableBId,
			ReferenceCnt: 1,  // 目前只导入了表级血缘, 以后导入字段血缘时, 这个应该是关联的字段个数
		}
		_, _ = o.Insert(&newRelation)
	}
}

type MetaDataImport struct {
	ImportId      int64
	ImportName    string
	ConnectionUrl string
}

func getConn(id int64) ConnectionManager {
	var connManager ConnectionManager
	connManager.GetById(id)
	return connManager
}

func AddNewVersionFieldData(data NewVersionData, totallyNew bool, tableId int64) {
	var fieldDatas []FieldMetaData
	now := utils.TimeFormat(time.Now())
	for _, field := range data.Fields {
		fieldDatas = append(fieldDatas, FieldMetaData{
			TableId:          data.TableId,
			FieldType:        field["fieldType"],
			FieldName:        field["fieldName"],
			Description:      field["description"],
			Version:          data.Version,
			ModificationTime: now,
		})
	}

	o := orm.NewOrm()
	if totallyNew {
		// 版本表插入信息
		versionData := VersionMetaData{
			TableId:    tableId,
			Version:    data.Version,
			CreateTime: now,
		}
		_, err1 := o.Insert(&versionData)
		if err1 != nil {
			panic(err1)
		}
	}
	// 字段表插入信息
	_, err := o.InsertMulti(100, fieldDatas)
	if err != nil {
		panic(err)
	}

	updateTableModifyTimeByField(now, data.TableId, o)
}

func GetFieldsAllVersion(tableId int64) []VersionMetaData {
	o := orm.NewOrm()
	var versions []VersionMetaData
	sql := "SELECT * FROM version_meta_data WHERE table_id=? ORDER BY create_time DESC"
	_, err := o.Raw(sql, tableId).QueryRows(&versions)
	if err != nil {
		panic(err)
	}
	return versions
}

func GetDefaultNewVersion(tableId int64) string {
	versions := GetFieldsAllVersion(tableId)
	if (versions == nil) || (len(versions) == 0) {
		//生成全新版本
		prefix := "v0"
		suffix := time.Now().Format("20060102-1504")
		return fmt.Sprintf("%s_%s", prefix, suffix)
	} else {
		//生成一个比最新版本号大一号的版本
		lastVersion := versions[0].Version
		lastPrefix := strings.Split(lastVersion, "_")[0]
		lastNumber, _ := strconv.ParseInt(lastPrefix[1:], 10, 64)
		curNumber := lastNumber + 1
		suffix := time.Now().Format("20060102-1504")
		return fmt.Sprintf("v%d_%s", curNumber, suffix)

	}

}

func GetTableNameById(tableId int64) string {
	o := orm.NewOrm()
	var tableNames []string
	sql := "SELECT table_name FROM table_meta_data WHERE tid=?"
	_, err := o.Raw(sql, tableId).QueryRows(&tableNames)
	if err != nil {
		panic(err)
	}
	return tableNames[0]
}

func GetFieldsListByVersion(tableId int64, version string) []FieldMetaData {
	o := orm.NewOrm()
	var fieldData []FieldMetaData
	sql := "SELECT * FROM field_meta_data WHERE table_id=? AND version=?"
	_, err := o.Raw(sql, tableId, version).QueryRows(&fieldData)
	if err != nil {
		panic(err)
	}
	return fieldData
}

func GetFieldsDataByVersion(tableId int64, version string) utils.VxeData {
	o := orm.NewOrm()
	var fieldData []FieldMetaData
	var err error
	var sql string
	if strings.ToLower(version) == "latest" {
		sql = `SELECT * FROM field_meta_data WHERE table_id=? AND version=(
		SELECT version FROM version_meta_data WHERE table_id=? ORDER BY create_time DESC limit 1
		) AND is_partition_field=false`
		_, err = o.Raw(sql, tableId, tableId).QueryRows(&fieldData)
	} else {
		sql = "SELECT * FROM field_meta_data WHERE table_id=? AND version=? AND is_partition_field=false"
		_, err = o.Raw(sql, tableId, version).QueryRows(&fieldData)
	}

	if err != nil {
		panic(err)
	}
	view := transTableView(fieldData)
	return view
}

func GetPartitionFieldsByVersion(tableId int64, version string) []FieldMetaData {
	o := orm.NewOrm()
	fieldData := make([]FieldMetaData, 0)
	var sql string
	var err error
	if strings.ToLower(version) == "latest" {
		sql = `SELECT * FROM field_meta_data WHERE table_id=? AND version=(
SELECT version FROM version_meta_data WHERE table_id=? ORDER BY create_time DESC limit 1
) AND is_partition_field=true`
		_, err = o.Raw(sql, tableId, tableId).QueryRows(&fieldData)
	} else {
		sql = "SELECT * FROM field_meta_data WHERE table_id=? AND version=? AND is_partition_field=true"
		_, err = o.Raw(sql, tableId, version).QueryRows(&fieldData)

	}
	if err != nil {
		panic(err)
	}
	return fieldData
}

func transTableView(fieldDatas []FieldMetaData) utils.VxeData {
	var vxe utils.VxeData
	editableCols := map[string]bool{
		"description": true,
	}
	jsonTagMapping := vxe.MakeColumns(FieldMetaData{}, editableCols)

	var dataList []map[string]interface{}
	for i, _ := range fieldDatas {
		dataList = append(dataList, vxe.MakeData(fieldDatas[i], &jsonTagMapping))
	}

	vxe.DataList = dataList
	return vxe
}

func UpdateFieldsByVersion(data FieldMetaData, now string, tableId int64) {
	o := orm.NewOrm()
	_, err := o.Update(data) // 更新字段信息
	if err != nil {
		panic(err)
	}
	updateTableModifyTimeByField(now, tableId, o)
}

// 更新表的元数据修改时间
func updateTableModifyTimeByField(now string, tableId int64, o orm.Ormer) {
	fmt.Println("更新表的修改时间:", tableId, "===>", now)
	//o := orm.NewOrm()
	//res, err := o.Raw("UPDATE user SET name = ?", "your").Exec()
}

func AddFieldRelation(tableAId int64, fieldAId int64, tableAVersion string, tableBId int64, fieldBId int64, tableBVersion string, description string) {
	data := RelationMetaData{
		TableAId:    tableAId,
		FieldAId:    fieldAId,
		VersionA:    tableAVersion,
		TableBId:    tableBId,
		FieldBId:    fieldBId,
		VersionB:    tableBVersion,
		Description: description,
	}
	o := orm.NewOrm()
	_, err := o.Insert(&data)
	if err != nil {
		panic(err)
	}

	var relation TableRelation
	_ = o.Raw("SELECT * FROM table_relation WHERE table_a_id=? AND table_b_id=?", tableAId, tableBId).
		QueryRow(&relation)
	if relation.Trid == 0 {
		newRelation := TableRelation{
			TableAId:     tableAId,
			TableBId:     tableBId,
			ReferenceCnt: 1,
		}
		_, _ = o.Insert(&newRelation)
	} else {
		newRelation := TableRelation{
			Trid:         relation.Trid,
			TableAId:     tableAId,
			TableBId:     tableBId,
			ReferenceCnt: relation.ReferenceCnt + 1,
		}
		_, _ = o.Update(&newRelation)
	}
}

type RelationView struct {
	Rid         int64  `field:"rid,hidden"`
	TableName   string `field:"关联表"`
	FieldName   string `field:"关联字段"`
	FieldType   string `field:"字段类型"`
	Version     string `field:"版本"`
	Description string `field:"关系描述"`
}

func GetFieldRelation(fieldId int64) utils.TableData {
	var relation []RelationView
	o := orm.NewOrm()
	sql := `SELECT 
	t1.rid,
	CONCAT(t3.table_type,': ',t3.database_name,'.',t3.table_name) as table_name,
	t2.field_name,
	t2.field_type,
	t1.version_b as version,
	t1.description
FROM
(SELECT rid,table_b_id,field_b_id,version_b, description FROM relation_meta_data WHERE field_a_id = ?) t1
LEFT JOIN 
(SELECT fid,field_name,field_type FROM field_meta_data) t2
on t1.field_b_id = t2.fid
LEFT JOIN
(SELECT tid,table_type,database_name, table_name FROM table_meta_data) t3
ON t1.table_b_id = t3.tid`
	_, err := o.Raw(sql, fieldId).QueryRows(&relation)
	if err != nil {
		panic(err)
	}
	// transView
	var bodys []map[string]interface{}
	for _, data := range relation {
		// 转成map
		tmpData := utils.GetSnakeStrMapOfStruct(data)
		bodys = append(bodys, tmpData)
	}

	// 返回数据

	tsb := utils.TableData{Data: bodys}
	tsb.GetSnakeTitles(RelationView{})

	return tsb
}

func DeleteFieldById(fieldId int64) {
	o := orm.NewOrm()
	_, err := o.Raw("DELETE FROM field_meta_data WHERE fid=?", fieldId).Exec()
	if err != nil {
		panic(err)
	}
}

func DeleteFieldRelation(rid int64) {
	o := orm.NewOrm()
	relationMetaData := RelationMetaData{
		Rid: rid,
	}
	_ = o.Read(&relationMetaData)
	_, err := o.Raw("DELETE FROM relation_meta_data WHERE rid=?", rid).Exec()
	_, err = o.Raw("UPDATE table_relation SET reference_cnt = reference_cnt - 1 WHERE table_a_id=? AND table_b_id=?",
		relationMetaData.TableAId, relationMetaData.TableBId).Exec()
	if err != nil {
		panic(err)
	}
}

type TableBView struct {
	TableBId int64
	FieldBId int64
	VersionB string
}

func InsertRelationByFile(content string, tableAId int64, fieldAId int64, tableAVersion string) {
	r := csv.NewReader(strings.NewReader(content))
	o := orm.NewOrm()
	var batchRes []RelationMetaData
	relationMap := make(map[int64]int64)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		// csv文件6列
		tableType := record[0]   // 表类型
		dataBase := record[1]    // 数据库
		tableName := record[2]   // 表名
		fieldName := record[3]   // 字段名
		version := record[4]     // 字段版本
		description := record[5] // 描述

		var tableB TableBView
		sql := `
SELECT 
	t2.tid as table_b_id,
	t1.fid as field_b_id, 
	t1.version as version_b
FROM
	( SELECT fid,table_id,version FROM field_meta_data WHERE field_name = ? AND version = ? ) t1 
left join 
	( SELECT tid FROM table_meta_data WHERE
      table_type = ? AND database_name = ? AND table_name = ?
	) t2
ON t1.table_id = t2.tid
limit 1`
		err = o.Raw(sql, fieldName, version, tableType, dataBase, tableName).QueryRow(&tableB)
		if err != nil {
			panic(err.Error())
		}

		batchRes = append(batchRes, RelationMetaData{
			TableAId:    tableAId,
			FieldAId:    fieldAId,
			VersionA:    tableAVersion,
			TableBId:    tableB.TableBId,
			FieldBId:    tableB.FieldBId,
			VersionB:    tableB.VersionB,
			Description: description,
		})

		if relationMap[tableB.TableBId] == 0 {
			relationMap[tableB.TableBId] = 1
		} else {
			relationMap[tableB.TableBId] = relationMap[tableB.TableBId] + 1
		}
	}

	_, err := o.InsertMulti(100, batchRes)
	if err != nil {
		panic(err)
	}

	for tableBId, cnt := range relationMap {
		var relation TableRelation
		_ = o.Raw("SELECT * FROM table_relation WHERE table_a_id=? AND table_b_id=?", tableAId, tableBId).
			QueryRow(&relation)
		if relation.Trid == 0 {
			newRelation := TableRelation{
				TableAId:     tableAId,
				TableBId:     tableBId,
				ReferenceCnt: cnt,
			}
			_, _ = o.Insert(&newRelation)
		} else {
			newRelation := TableRelation{
				Trid:         relation.Trid,
				TableAId:     tableAId,
				TableBId:     tableBId,
				ReferenceCnt: relation.ReferenceCnt + cnt,
			}
			_, _ = o.Update(&newRelation)
		}
	}
}

type NodeView struct {
	Id      string `json:"-"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Name    string `json:"name"`
	Level   int    `json:"-"`
	TableId string `json:"tableId"`
}
type EdgeView struct {
	Source string `json:"source"`
	Target string `json:"target"`
}
type TableView struct {
	Id        string
	Name      string
	TableType string
}

type GraphRelationView struct {
	FieldAName    string `json:"field_a_name" vxe:"字段A"`
	FieldAVersion string `json:"field_a_version" vxe:"版本A"`
	FieldBName    string `json:"field_b_name" vxe:"字段B"`
	FieldBVersion string `json:"field_b_version" vxe:"版本B"`
	Description   string `json:"description" vxe:"描述"`
}

func GetGraphRelation(sourceName string, targetName string) utils.VxeData {
	o := orm.NewOrm()
	var source TableMetaData
	var target TableMetaData
	_ = o.Raw(`SELECT * FROM table_meta_data WHERE table_name=? limit 1`, sourceName).QueryRow(&source)
	_ = o.Raw(`SELECT * FROM table_meta_data WHERE table_name=? limit 1`, targetName).QueryRow(&target)
	var views []GraphRelationView

	upSql := `SELECT 
	t2.field_name as field_a_name, 
	t1.version_a as field_a_version, 
	t3.field_name as field_b_name, 
	t1.version_b as field_b_version,
	t1.description as description
FROM
(SELECT * FROM relation_meta_data WHERE table_a_id=? AND table_b_id=?) t1
LEFT JOIN
	field_meta_data t2 ON t1.field_a_id=t2.fid
LEFT JOIN
	field_meta_data t3 ON t1.field_b_id=t3.fid
`
	var upView []GraphRelationView
	_, _ = o.Raw(upSql, source.Tid, target.Tid).QueryRows(&upView)

	downSql := `SELECT 
	t2.field_name as field_b_name, 
	t1.version_a as field_b_version, 
	t3.field_name as field_a_name, 
	t1.version_b as field_a_version,
	t1.description as description
FROM
(SELECT * FROM relation_meta_data WHERE table_a_id=? AND table_b_id=?) t1
LEFT JOIN
	field_meta_data t2 ON t1.field_a_id=t2.fid
LEFT JOIN
	field_meta_data t3 ON t1.field_b_id=t3.fid
`
	var downView []GraphRelationView
	_, _ = o.Raw(downSql, target.Tid, source.Tid).QueryRows(&downView)
	views = append(views, upView...)
	views = append(views, downView...)

	var vxe utils.VxeData
	jsonTagMapping := vxe.MakeColumns(GraphRelationView{}, nil)

	var dataList []map[string]interface{}
	for i, _ := range views {
		dataList = append(dataList, vxe.MakeData(views[i], &jsonTagMapping))
	}

	vxe.DataList = dataList
	return vxe
}

func GetGraph(tid int64) map[string]interface{} {
	//colorMap := GetAllTableTypesMap()
	o := orm.NewOrm()
	startX := 10 // 每轮加200
	startY := 0  // 每次加60, 每轮重置

	minLevel := 0 // 上游左移的最深层次

	x := startX
	y := startY

	//sqlDownStream := `
	sqlUpStream := `
SELECT t2.* FROM
(SELECT table_b_id as id FROM table_relation WHERE table_a_id = ? AND reference_cnt>0) t1
LEFT JOIN
(SELECT tid as id,table_name as name,table_type FROM table_meta_data) t2
ON 
t1.id = t2.id`
	//sqlUpStream := `
	sqlDownStream := `
SELECT t2.* FROM
(SELECT table_a_id as id FROM table_relation WHERE table_b_id = ? AND reference_cnt>0) t1
LEFT JOIN
(SELECT tid as id,table_name as name,table_type FROM table_meta_data) t2
ON 
t1.id = t2.id`

	// 查询自身的表名
	var selfTable TableView
	_ = o.Raw("SELECT tid as id,table_name as name,table_type FROM table_meta_data WHERE tid=?", tid).QueryRow(&selfTable)

	nodeSet := make(map[string]bool) // 不要重复加入节点

	// 上游计算
	var upNodes []NodeView
	var upEdges []EdgeView
	tmpQueue := list.New()
	// upStream加入自身
	tmpQueue.PushBack(NodeView{ // enqueue
		Id:      fmt.Sprintf("node_%d", tid),
		X:       x,
		Y:       y,
		Name:    selfTable.Name,
		Level:   0,
		TableId: fmt.Sprint(tid),
	})
	x = x - 200
	for {
		if tmpQueue.Len() == 0 {
			break
		}

		firstEle := tmpQueue.Front() // First element
		// 查号当前节点的上游节点, 并把上游节点加入tmpQueue, 把当前节点弹出加入upNodes
		curNode := firstEle.Value.(NodeView)
		var ups []TableView
		_, _ = o.Raw(sqlUpStream, strings.Split(curNode.Id, "_")[1]).QueryRows(&ups)

		// (3)tmpQueue删除当前节点, nodes加入节点
		tmpQueue.Remove(firstEle) // Dequeue
		if nodeSet[curNode.Id] == false {
			upNodes = append(upNodes, curNode)
			nodeSet[curNode.Id] = true

			// (1)当前节点的所有上游节点加入tmpQueue
			// (2)所有上有节点和当前节点的边加入upEdges
			for _, t := range ups {
				tmpQueue.PushBack(NodeView{
					Id:      fmt.Sprintf("node_%s", t.Id),
					X:       x,
					Y:       y,
					Name:    t.Name,
					Level:   curNode.Level + 1,
					TableId: fmt.Sprint(t.Id),
				})
				y = y + 60

				upEdges = append(upEdges, EdgeView{
					Source: t.Name,       // fmt.Sprintf("node_%s", t.Id),
					Target: curNode.Name, //curNode.Id,
				})
			}

			x = x - 200
			y = 0

		}

		if curNode.Level > minLevel {
			minLevel = curNode.Level
		}
	}

	// 下游
	x = startX
	y = startY
	var downNodes []NodeView
	var downEdges []EdgeView
	// downstream加入自身
	tmpQueue.PushBack(NodeView{ // enqueue
		Id:   fmt.Sprintf("node_%d", tid),
		X:    x,
		Y:    y,
		Name: selfTable.Name, // "自身",
	})
	x = x + 200
	for {
		if tmpQueue.Len() == 0 {
			break
		}
		firstEle := tmpQueue.Front() // First element
		curNode := firstEle.Value.(NodeView)
		var downs []TableView
		_, _ = o.Raw(sqlDownStream, strings.Split(curNode.Id, "_")[1]).QueryRows(&downs)

		// (3)tmpQueue删除当前节点, nodes加入节点
		tmpQueue.Remove(firstEle) // Dequeue
		if (nodeSet[curNode.Id] == false) || (curNode.Id == fmt.Sprintf("node_%s", selfTable.Id)) {
			if curNode.Id != fmt.Sprintf("node_%s", selfTable.Id) {
				downNodes = append(downNodes, curNode)
				nodeSet[curNode.Id] = true
			}

			// (1)当前节点的所有下游节点加入tmpQueue
			// (2)所有下游节点和当前节点的边加入downEdges
			for _, t := range downs { // (1)当前节点的上游节点
				tmpQueue.PushBack(NodeView{
					Id: fmt.Sprintf("node_%s", t.Id),
					X:  x,
					Y:  y,
					//Height:    fmt.Sprint(50),
					//Width:     fmt.Sprint(120),
					Name: t.Name,
					//ClassName: colorMap[t.TableType],
					TableId: fmt.Sprint(t.Id),
				})
				y = y + 60

				downEdges = append(downEdges, EdgeView{
					Source: curNode.Name, //curNode.Id,
					Target: t.Name,       //fmt.Sprintf("node_%s", t.Id),
				})
			}

			x = x + 200
			y = 0
		}

	}

	// 已通过nodeSet去重, 下游节点不必再次删除自身
	// downNodes = downNodes[1:]

	var nodes []NodeView
	nodes = append(upNodes, downNodes...)
	var edges []EdgeView
	edges = append(upEdges, downEdges...)

	//所有节点的x要向右移动
	for index, _ := range nodes {
		nodes[index].X = nodes[index].X + minLevel*200
	}

	res := make(map[string]interface{})
	res["data"] = nodes
	res["links"] = edges
	return res
}

//单表自动导入
func GenerateFieldByDefSingle(param map[string]string, user string) bool {
	importMetaData := ImportMetaData{}
	roles := Enforcer.GetRolesForUser(user)
	var tid int64
	isSuccess := false

	duplicate := CheckDuplicateTable(param["tableType"], param["databaseName"], param["tableName"])
	if duplicate {
		return false
	}

	switch param["tableType"] {
	case "mysql":
		importId, _ := strconv.ParseInt(param["importId"], 10, 64)
		conn := getConn(importId)
		dbName := param["databaseName"]
		tName := param["tableName"]
		isSuccess, tid = importMetaData.importMysqlSingle(conn, dbName, tName, roles)
	case "hive":
		isSuccess, tid = importMetaData.importHiveSingle(param, roles)
	case "hdfs":
		isSuccess, tid = importMetaData.importHdfsSingle(user, param, roles)
	case "clickhouse":
		importId, _ := strconv.ParseInt(param["importId"], 10, 64)
		conn := getConn(importId)
		ck := CKModel{}
		isSuccess, tid = ck.GetFieldsSingle(param, conn.Host, conn.Port, conn.User, conn.Passwd, importId, roles)
	default:
	}
	if tid != 0 {
		privilegeRoles := strings.Split(param["privilegeRoles"], ",")
		uniqRoles := make(map[string]string)
		for _, r := range privilegeRoles {
			if r != "" {
				uniqRoles[r] = ""
			}
		}
		for _,r := range roles {
			uniqRoles[r] = ""
		}
		var roleSet []string // 去重后的集合
		for k, _ := range uniqRoles {
			roleSet = append(roleSet, k)
		}
		AddTablePrivilegeMapping(tid, roleSet)
	}
	return isSuccess
}

// 自动导入元数据
func GenerateFieldByDefinition(param map[string]string, user string) {
	importMetaData := ImportMetaData{}
	switch param["tableType"] {
	case "mysql":
		fmt.Println("mysql")
		//importId, _ := strconv.ParseInt(param["importId"], 10, 64)
		//mysqlImportData := getConn(importId)
		//importMetaData.getMySql(mysqlImportData., mysqlImportData.ConnectionUrl)
	case "hive":
		fmt.Println("hive")
		hiveMetaData := importMetaData.getHive(param)
		importMetaData.insertHiveField(hiveMetaData, user)
	case "hdfs":
		fmt.Println("hdfs")
	case "clickhouse":
		fmt.Println("clickhouse")
		getClickHouse()
	default:
		fmt.Println("未知类型:", param["tableType"])
	}
}

func getClickHouse() {
	ck := CKModel{}
	ck.GetFields()
}

func UpdateVersionData(data *VersionMetaData) {
	o := orm.NewOrm()
	_, _ = o.Raw("UPDATE version_meta_data set description=? WHERE table_id=? AND version=?",
		data.Description,
		data.TableId,
		data.Version).Exec()
}

func ExecUpdateFields(updates *[]FieldMetaData) {
	o := orm.NewOrm()
	// 更新hive
	execUpdateCmd(updates, o)
	// 更新数据库
	now := utils.TimeFormat(time.Now())
	for i, _ := range *updates {
		(*updates)[i].ModificationTime = now
		if _, err := o.Update(&(*updates)[i], "Description", "ModificationTime"); err != nil {
			fmt.Println(err)
			panic(err.Error())
		}
	}
}

func ExecUpdatePartitionField(update *FieldMetaData) {
	o := orm.NewOrm()
	// 更新数据库
	now := utils.TimeFormat(time.Now())
	(*update).ModificationTime = now
	if _, err := o.Update(update, "Description", "ModificationTime"); err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
}

func execUpdateCmd(updates *[]FieldMetaData, o orm.Ormer) {
	var sqlBatch []string
	for i, _ := range *updates {
		fieldData := (*updates)[i]
		t := TableMetaData{
			Tid: fieldData.TableId,
		}
		if err := o.Read(&t); err != nil {
			if err == orm.ErrNoRows {
				fmt.Printf("查询不到table: %d\n", fieldData.TableId)
			} else if err == orm.ErrMissPK {
				fmt.Println("找不到主键")
			}
			panic(err)
		}

		switch t.TableType {
		case "hive":
			database := t.DatabaseName
			tableName := t.TableName
			oldColName := fieldData.FieldName
			typeName := fieldData.FieldType
			comment := fieldData.Description

			sql := fmt.Sprintf("ALTER TABLE %s.%s change column %s %s %s comment '%s'",
				database, tableName, oldColName, oldColName, typeName, comment)
			sqlBatch = append(sqlBatch, sql)
			sql = fmt.Sprintf("%s;", strings.Join(sqlBatch, ";"))

			fmt.Println(sql)
			utils.ExecHiveCmd(sql)

		default:
		}
	}

}
