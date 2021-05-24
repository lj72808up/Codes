package models

import (
	"AStar/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/jmoiron/sqlx"
	_ "github.com/kshvakov/clickhouse"
	"strings"
	"time"
)

type CKModel struct{}
type CKField struct {
	TtlExpression     string `db:"ttl_expression"`
	DefaultType       string `db:"default_type"`
	Name              string `db:"name"`
	DefaultExpression string `db:"default_expression"`
	Comment           string `db:"comment"`
	Type              string `db:"type"`
	CodecExpression   string `db:"codec_expression"`
}

func (this *CKModel) GetFields() bool {
	o := orm.NewOrm()

	//url := "tcp://ck.adtech.sogou:9000?&user=default&password=6lYaUiFi"
	url := beego.AppConfig.String("clickhouse")
	conn, err := sqlx.Open("clickhouse", url)
	if err != nil {
		panic(err.Error())
	}

	var dbs []string
	if err := conn.Select(&dbs, "show databases"); err != nil {
		panic(err.Error())
		return false
	}

	fmt.Println(dbs)
	_ = conn.Close()

	var fields []FieldMetaData
	var versions []VersionMetaData

	for _, db := range dbs {
		if db == "system" {
			continue
		}
		var tables []string
		dbConn, err := sqlx.Open("clickhouse", fmt.Sprintf("%s&database=%s", url, db))
		if err != nil {
			panic(err.Error())
			return false
		}
		_ = dbConn.Select(&tables, "show tables")
		fmt.Println(tables)
		now := utils.TimeFormat(time.Now())
		var selectFields []CKField
		for _, tbl := range tables {
			// insert table
			table := TableMetaData{
				TableName:    tbl,
				TableType:    "clickhouse",
				DatabaseName: db,
			}
			tid, _ := o.Insert(&table)
			// get version
			newVersion := GetDefaultNewVersion(tid)
			versions = append(versions, VersionMetaData{
				TableId:    tid,
				Version:    newVersion,
				CreateTime: now,
			})

			dbConn.Select(&selectFields, fmt.Sprintf("desc %s", tbl))
			for i, _ := range selectFields {
				fields = append(fields, FieldMetaData{
					TableId:          tid,
					FieldName:        selectFields[i].Name,
					FieldType:        selectFields[i].Type,
					Description:      "",
					Version:          newVersion,
					ModificationTime: "",
				})
			}
		}

		_ = dbConn.Close()
	}

	_, _ = o.InsertMulti(100, versions)
	_, _ = o.InsertMulti(100, fields)

	return true

}

func (this *CKModel) GetFieldsSingle(param map[string]string, host string, port string, user string, passwd string, importId int64, roles []string) (bool, int64) {
	o := orm.NewOrm()
	_ = o.Using("default")

	dbName := param["databaseName"]
	tName := param["tableName"]

	//url := "tcp://ck.adtech.sogou:9000?&user=default&password=6lYaUiFi"
	url := fmt.Sprintf("tcp://%s:%s?user=%s&password=%s&database=%s", host, port, user, passwd, dbName)

	if dbName == "" || tName == "" {
		return false, 0
	}

	var fields []FieldMetaData

	now := utils.TimeFormat(time.Now())
	selectFields := this.syncCkField(url, tName)

	if (selectFields == nil) || (len(selectFields) == 0) {
		return false, 0
	}
	// insert table
	table := TableMetaData{
		TableName:    tName,
		TableType:    "clickhouse",
		DatabaseName: dbName,
		Identifier:   GetTableIdentifier("clickhouse", dbName, tName),
		Dsid:         importId,
		Owner:        strings.Join(roles, ","),
	}
	tid, _ := o.Insert(&table)
	// get version
	newVersion := GetDefaultNewVersion(tid)
	version := VersionMetaData{
		TableId:    tid,
		Version:    newVersion,
		CreateTime: now,
	}
	_, _ = o.Insert(&version)

	for i, _ := range selectFields {
		fields = append(fields, FieldMetaData{
			TableId:          tid,
			FieldName:        selectFields[i].Name,
			FieldType:        selectFields[i].Type,
			Description:      "",
			Version:          newVersion,
			ModificationTime: "",
		})
	}

	_, _ = o.InsertMulti(100, fields)
	return true, tid
}

func (this *CKModel) syncCkField(ckURL string, tName string) []CKField {
	dbConn, err := sqlx.Open("clickhouse", ckURL)
	if err != nil {
		panic(err.Error())
	}
	var selectFields []CKField
	_ = dbConn.Select(&selectFields, fmt.Sprintf("desc %s", tName))
	_ = dbConn.Close()
	return selectFields
}

func (this *CKModel) DiffField(table TableMetaData, tableVersion string) {
	o := orm.NewOrm()
	tid := table.Tid

	var conn ConnectionManager
	_ = o.Raw("SELECT * FROM connection_manager WHERE dsid = ?", table.Dsid).QueryRow(&conn)

	host := conn.Host
	port := conn.Port
	user := conn.User
	passwd := conn.Passwd

	tableName := table.TableName
	dbName := table.DatabaseName
	url := fmt.Sprintf("tcp://%s:%s?user=%s&password=%s&database=%s", host, port, user, passwd, dbName)

	// online field
	onlineFields := this.syncCkField(url, tableName)

	var recordFields []FieldMetaData
	_, dbErr := o.Raw("SELECT * FROM field_meta_data WHERE table_id = ? AND version = ?", tid, tableVersion).QueryRows(&recordFields)
	if dbErr != nil {
		Logger.Error("从数据库查看(%d)字段失败: %s", tid, dbErr.Error())
		panic(fmt.Sprintf("从数据库查看(%d)字段失败", tid))
	}

	// diff field
	if len(onlineFields) > 0 {
		this.UpdateFields(&onlineFields, &recordFields)
		//increased, deleted := diffFields(onlineFields, recordFields)
		//HandleOnlineRecordFieldDiff(&increased, &deleted)
	} else {
		Logger.Error("没有从CK中查到(tid:%d)的字段", tid)
		return
	}
}

func diffFields(onlineFields []CKField, recordFields []FieldMetaData) ([]FieldMetaData, []FieldMetaData) {
	onlineFieldsMap := make(map[string]string)
	recordFieldsMap := make(map[string]string)

	for _, f := range onlineFields {
		onlineFieldsMap[f.Name] = f.Type
	}

	for _, f := range recordFields {
		recordFieldsMap[f.FieldName] = f.FieldType
	}

	var increased []FieldMetaData
	var deleted []FieldMetaData

	if len(recordFields) == 0 {
		Logger.Error("数据库中未找到表匹配的字段")
		return []FieldMetaData{}, []FieldMetaData{}
	}

	tid := recordFields[0].TableId
	version := recordFields[0].Version
	now := utils.TimeFormat(time.Now())
	// online 比 record 新增的字段
	for _, onlineField := range onlineFields {
		if recordFieldsMap[onlineField.Name] == "" {
			increased = append(increased, FieldMetaData{
				TableId:          tid,
				FieldName:        onlineField.Name,
				FieldType:        onlineField.Type,
				Description:      "",
				Version:          version,
				ModificationTime: now,
				IsPartitionField: false,
			})
		}
	}

	// online 比 record 已删除的字段
	for _, recordField := range recordFields {
		if onlineFieldsMap[recordField.FieldName] == "" {
			deleted = append(deleted, recordField)
		}
	}

	return increased, deleted
}

func (this *CKModel) UpdateFields(onlineFields *[]CKField, recordFields *[]FieldMetaData) {
	tableId := (*recordFields)[0].TableId
	version := (*recordFields)[0].Version
	now := utils.TimeFormat(time.Now())

	descMapping := make(map[string]string)
	for _, f := range *recordFields {
		descMapping[f.FieldName] = f.Description
	}

	// 插入recordFields
	paddingDescFields := make([]FieldMetaData, 0)
	for _, f := range *onlineFields {
		desc := f.Comment
		if desc == "" {
			desc = descMapping[f.Name]
			f.Comment = desc // 取已经存在的注释为字段的注释
		}
		paddingDescFields = append(paddingDescFields, FieldMetaData{
			TableId:          tableId,
			FieldName:        f.Name,
			FieldType:        f.Type,
			Description:      desc,
			Version:          version,
			ModificationTime: now,
			IsPartitionField: false,
		})
	}

	// 删除onlineFields
	o := orm.NewOrm()
	for _, f := range *recordFields {
		_, _ = o.Raw("DELETE FROM field_meta_data WHERE fid = ?", f.Fid).Exec()
	}
	// 插入新的字段
	_, _ = o.InsertMulti(100, paddingDescFields)
}
