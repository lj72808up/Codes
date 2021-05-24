package models

import (
	adtl_platform "AStar/protobuf"
	"AStar/utils"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"github.com/golang/protobuf/proto"
	"strings"
	"time"
)

type MetaTreeNode struct {
	Label     string         `json:"label"`
	TableName string         `json:"tableName"`
	TableId   int64          `json:"tid"`
	Children  []MetaTreeNode `json:"children,omitempty"`
}

type HiveModel struct{}

func (this *HiveModel) DiffField(table TableMetaData, version string) {
	tid := table.Tid
	o := orm.NewOrm()

	hDatabase := table.DatabaseName
	hTable := table.TableName
	protoTable := this.GetHiveProtoMeta(hDatabase, hTable)
	onlineFields := this.ParseFields(tid, version, protoTable)
	
	var recordFields []FieldMetaData
	_, dbErr := o.Raw("SELECT * FROM field_meta_data WHERE table_id = ? AND version = ?", tid, version).QueryRows(&recordFields)
	if dbErr != nil {
		Logger.Error("从数据库查看(%d)字段失败: %s", tid, dbErr.Error())
		panic(fmt.Sprintf("从数据库查看(%d)字段失败", tid))
	}

	// diff field
	if len(*onlineFields) > 0{
		this.UpdateFields(onlineFields, &recordFields)
		/*increased, deleted := this.diffFields(onlineFields, &recordFields)
		HandleOnlineRecordFieldDiff(increased, deleted)*/
	} else {
		Logger.Error("没有从CK中查到(tid:%d)的字段", tid)
		return
	}
}

func (this *HiveModel) ParseFields(tid int64, newVersion string, tb *adtl_platform.HiveTable) *[]FieldMetaData {
	var fields []FieldMetaData
	var partitionMap = make(map[string]bool)
	partitionCols := strings.Split(tb.ExpendInfo["partitions"], ",")
	for _, col := range partitionCols {
		partitionMap[col] = true
	}

	for _, f := range tb.IntegratedField {
		var fName = f.FName
		isPartition := false
		if partitionMap[fName] {
			isPartition = true
		}
		fields = append(fields, FieldMetaData{
			TableId:          tid,
			FieldName:        f.FName,
			FieldType:        f.FType,
			Description:      f.FComment,
			Version:          newVersion,
			ModificationTime: "",
			IsPartitionField: isPartition,
		})
	}
	return &fields
}

func (this *HiveModel) GetHiveProtoMeta(hiveDB string, hTable string) *adtl_platform.HiveTable {
	metaURL := utils.SparkRoute(ServerInfo.SparkHost, ServerInfo.SparkPort, "get_singleMeta")
	metaURL = fmt.Sprintf("%s?dbName=%s&tbName=%s", metaURL, hiveDB, hTable)
	req := httplib.Get(metaURL).SetTimeout(time.Hour, time.Hour)
	bytes, err := req.Bytes()
	if err != nil {
		panic(err.Error())
	}
	hiveTable := &adtl_platform.HiveTable{}
	_ = proto.Unmarshal(bytes, hiveTable)
	return hiveTable
}

func (this *HiveModel) diffFields(onlineFields *[]FieldMetaData, recordFields *[]FieldMetaData) (*[]FieldMetaData, *[]FieldMetaData) {
	onlineFieldsMap := make(map[string]string)
	recordFieldsMap := make(map[string]string)

	for _, f := range *onlineFields {
		onlineFieldsMap[f.FieldName] = f.FieldType
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
		if recordFieldsMap[onlineField.FieldName] == "" {
			increased = append(increased, FieldMetaData{
				TableId:          tid,
				FieldName:        onlineField.FieldName,
				FieldType:        onlineField.FieldType,
				Description:      onlineField.Description,
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




func (this *HiveModel) UpdateFields(onlineFields *[]FieldMetaData, recordFields *[]FieldMetaData) {
	descMapping := make(map[string]string)
	for _,f := range *recordFields {
		descMapping[f.FieldName] = f.Description
	}

	// 删除onlineFields
	o := orm.NewOrm()
	for _,f := range *recordFields {
		_, _ = o.Raw("DELETE FROM field_meta_data WHERE fid = ?", f.Fid).Exec()
	}

	// 插入recordFields
	paddingDescFields := make([]FieldMetaData,0)
	for _,f := range *onlineFields {
		if f.Description == "" {
			f.Description = descMapping[f.FieldName]
		}
		paddingDescFields = append(paddingDescFields, f)
	}
	_, _ = o.InsertMulti(100, paddingDescFields)
}
