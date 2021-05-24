package models

import (
	"AStar/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type MysqlModel struct{}

type MysqlField struct {
	Field string `orm:"column(Field)"`
	Type  string `orm:"column(Type)"`
}

func (this *MysqlModel) DiffField(table TableMetaData, version string) {
	// todo mysql的
	o := orm.NewOrm()
	tid := table.Tid

	var conn ConnectionManager
	_ = o.Raw("SELECT * FROM connection_manager WHERE dsid = ?", table.Dsid).QueryRow(&conn)

	jdbcUrl := conn.Url
	user := conn.User
	password := conn.Passwd
	tbName := table.TableName
	dbName := table.DatabaseName
	onlineFields := this.SyncMysqlFields(jdbcUrl, user, password, dbName, tbName)

	var recordFields []FieldMetaData
	_, dbErr := o.Raw("SELECT * FROM field_meta_data WHERE table_id = ? AND version = ?", tid, version).QueryRows(&recordFields)
	if dbErr != nil {
		Logger.Error("从数据库查看(%d)字段失败: %s", tid, dbErr.Error())
		panic(fmt.Sprintf("从数据库查看(%d)字段失败", tid))
	}

	// diff field
	if len(*onlineFields) > 0 {
		this.UpdateFields(onlineFields, &recordFields)
		/*increased, deleted := this.diffFields(onlineFields, &recordFields)
		HandleOnlineRecordFieldDiff(increased, deleted)*/
	} else {
		Logger.Error("没有从CK中查到(tid:%d)的字段", tid)
		return
	}
}

func (this *MysqlModel) SyncMysqlFields(jdbcUrl string, user string, passwd string, dbName string, tbName string) *[]MysqlField {

	ip := strings.Split(strings.Replace(jdbcUrl, "jdbc:mysql://", "", -1), "/")[0]
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, passwd, ip, dbName)
	Logger.Info("数据源: %s\n", url)

	aliasName := fmt.Sprintf("%s-%d", dbName, time.Now().Unix())
	err := orm.RegisterDataBase(aliasName, "mysql", url, 0, 1) // SetMaxIdleConns, SetMaxOpenConns
	if err != nil {
		panic(err)
	}
	o1 := orm.NewOrm()
	_ = o1.Using(aliasName)

	var onlineFields []MysqlField
	_, _ = o1.Raw(fmt.Sprintf("desc %s.%s", dbName, tbName)).QueryRows(&onlineFields)
	return &onlineFields
}

func (this *MysqlModel) ImportFieldsSingleTable(conn ConnectionManager, dbName string, tbName string, roles []string) (bool, int64) {
	var fields []FieldMetaData
	jdbcUrl := conn.Url
	user := conn.User
	passwd := conn.Passwd

	o := orm.NewOrm()
	now := utils.TimeFormat(time.Now())

	showFields := this.SyncMysqlFields(jdbcUrl, user, passwd, dbName, tbName)
	if (showFields != nil) && (len(*showFields) != 0) {
		// 1. 表
		table := TableMetaData{
			TableName:    tbName,
			TableType:    "mysql",
			DatabaseName: dbName,
			Identifier:   GetTableIdentifier("mysql", dbName, tbName),
			Dsid:         conn.Dsid,
			Owner:        strings.Join(roles, ","),
		}
		tid, _ := o.Insert(&table)
		// 2. 版本
		newVersion := GetDefaultNewVersion(tid)
		version := VersionMetaData{
			TableId:    tid,
			Version:    newVersion,
			CreateTime: now,
		}
		_, _ = o.Insert(&version)
		// 3. 字段
		for i, _ := range *showFields {
			fields = append(fields, FieldMetaData{
				TableId:   tid,
				FieldName: (*showFields)[i].Field,
				FieldType: (*showFields)[i].Type,
				Version:   newVersion,
			})
		}
		_, _ = o.InsertMulti(100, fields)
		return true, tid
	} else {
		return false, -1
	}
}

func (this *MysqlModel) diffFields(onlineFields *[]MysqlField, recordFields *[]FieldMetaData) (*[]FieldMetaData, *[]FieldMetaData) {
	onlineFieldsMap := make(map[string]string)
	recordFieldsMap := make(map[string]string)

	for _, f := range *onlineFields {
		onlineFieldsMap[f.Field] = f.Type
	}

	for _, f := range *recordFields {
		recordFieldsMap[f.FieldName] = f.FieldType
	}

	var increased []FieldMetaData
	var deleted []FieldMetaData

	if len(*recordFields) == 0 {
		Logger.Error("数据库中未找到表匹配的字段")
		return &[]FieldMetaData{}, &[]FieldMetaData{}
	}

	tid := (*recordFields)[0].TableId
	version := (*recordFields)[0].Version
	now := utils.TimeFormat(time.Now())
	// online 比 record 新增的字段
	for _, onlineField := range *onlineFields {
		if recordFieldsMap[onlineField.Field] == "" {
			increased = append(increased, FieldMetaData{
				TableId:          tid,
				FieldName:        onlineField.Field,
				FieldType:        onlineField.Type,
				Description:      "",
				Version:          version,
				ModificationTime: now,
				IsPartitionField: false,
			})
		}
	}

	// online 比 record 已删除的字段
	for _, recordField := range *recordFields {
		if onlineFieldsMap[recordField.FieldName] == "" {
			deleted = append(deleted, recordField)
		}
	}

	return &increased, &deleted
}

func (this *MysqlModel) UpdateFields(onlineFields *[]MysqlField, recordFields *[]FieldMetaData) {
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
		desc := descMapping[f.Field]

		paddingDescFields = append(paddingDescFields, FieldMetaData{
			TableId:          tableId,
			FieldName:        f.Field,
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
	// 插入新的列
	_, _ = o.InsertMulti(100, paddingDescFields)
}
