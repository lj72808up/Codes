package models

import (
	"AStar/protobuf"
	"AStar/utils"
	"encoding/csv"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"github.com/golang/protobuf/proto"
	"strconv"
	"strings"
	"time"
)

type ImportMetaData struct{}

func (this ImportMetaData) getMySql(dbName string, url string) {
	aliasName := fmt.Sprintf("%s-%d", dbName, time.Now().Unix())
	//dbName :=  "kl"
	//url := "reader:kunlun@tcp(10.134.10.113:3306)/kl?charset=utf8&loc=Asia%2FShanghai"
	// 1. 连接另一个数据源还能
	err := orm.RegisterDataBase(aliasName, "mysql", url, 0, 1) // SetMaxIdleConns, SetMaxOpenConns
	if err != nil {
		panic(err)
	}
	o1 := orm.NewOrm()
	_ = o1.Using(aliasName)
	o2 := orm.NewOrm()
	_ = o2.Using("default")
	var dbs []string
	now := utils.TimeFormat(time.Now())
	_, _ = o1.Raw("show databases").QueryRows(&dbs)

	var tables []string
	var fields []FieldMetaData
	var versions []VersionMetaData

	_, _ = o1.Raw("show tables").QueryRows(&tables) // .Exec()
	for _, tName := range tables {
		//if index>100{
		//	break
		//}
		// 1. 插入表
		table := TableMetaData{
			TableName:    tName,
			TableType:    "mysql",
			DatabaseName: dbName,
		}
		tid, err := o2.Insert(&table)
		if err != nil {
			panic(err)
		}
		// 2. 版本
		newVersion := GetDefaultNewVersion(tid)
		versions = append(versions, VersionMetaData{
			TableId:    tid,
			Version:    newVersion,
			CreateTime: now,
		})
		// 3. 字段
		var showFields []MysqlField
		_, _ = o1.Raw(fmt.Sprintf("desc %s.%s", dbName, tName)).QueryRows(&showFields)
		for i, _ := range showFields {
			fields = append(fields, FieldMetaData{
				TableId:   tid,
				FieldName: showFields[i].Field,
				FieldType: showFields[i].Type,
				Version:   newVersion,
			})
		}
	}
	_, _ = o2.InsertMulti(100, versions)
	_, _ = o2.InsertMulti(100, fields)

	Logger.Info("%v", dbs)
}

func (this ImportMetaData) getHive(param map[string]string) adtl_platform.HiveMeta {
	// 1. 访问接口, 拿到全量数据
	metaURL := utils.SparkRoute(ServerInfo.SparkHost, ServerInfo.SparkPort, "get_meta")
	req := httplib.Get(metaURL).SetTimeout(time.Hour, time.Hour)
	bytes, error := req.Bytes()
	if error != nil {
		panic(error.Error())
	}
	hiveMeta := &adtl_platform.HiveMeta{}
	_ = proto.Unmarshal(bytes, hiveMeta)

	// 2. 根据参数过滤出结果
	var filterDataBase []*adtl_platform.HiveDatabase
	if param["databaseName"] != "" {
		for _, db := range hiveMeta.Databases {
			if db.NameDatabase == param["databaseName"] {
				filterDataBase = append(filterDataBase, db)
			}
		}
		meta := adtl_platform.HiveMeta{
			Databases: filterDataBase,
		}
		return meta
	} else {
		return *hiveMeta
	}
}

func (this ImportMetaData) insertHiveField(data adtl_platform.HiveMeta, user string) {
	o := orm.NewOrm()
	now := utils.TimeFormat(time.Now())
	for _, db := range data.Databases {
		var fields []FieldMetaData
		var versions []VersionMetaData
		for _, t := range db.Talbes {
			compress := strings.Split(t.ExpendInfo["storageProperties"], ",")[0]
			compress = "[" + strings.Replace(strings.Replace(compress, "[", "", -1), "]", "", -1) + "]"
			table := TableMetaData{
				TableName:        t.NameTable,
				TableType:        "hive",
				DatabaseName:     db.NameDatabase,
				Creator:          t.ExpendInfo["owner"],
				CreateTime:       t.ExpendInfo["createdTime"],
				Modifier:         "",
				ModificationTime: "",
				Description:      "",
				Location:         t.ExpendInfo["location"],
				CompressMode:     compress,
			}
			tid, err := o.Insert(&table)
			if err != nil {
				panic(err)
			}
			newVersion := GetDefaultNewVersion(tid)
			versions = append(versions, VersionMetaData{
				TableId:    tid,
				Version:    newVersion,
				CreateTime: now,
			})

			for _, f := range t.IntegratedField {
				fields = append(fields, FieldMetaData{
					TableId:          tid,
					FieldName:        f.FName,
					FieldType:        f.FType,
					Description:      f.FComment,
					Version:          newVersion,
					ModificationTime: "",
				})
			}
		}
		_, _ = o.InsertMulti(100, versions)
		_, _ = o.InsertMulti(100, fields)
	}
}
func (this ImportMetaData) importHdfsSingle(user string, param map[string]string, roles []string) (bool, int64) {
	content := param["hdfsImportStr"]
	csvReader := csv.NewReader(strings.NewReader(content))

	o := orm.NewOrm()
	now := utils.TimeFormat(time.Now())
	// 1. table
	table := TableMetaData{
		TableName:    param["tableName"],
		TableType:    param["tableType"],
		DatabaseName: param["databaseName"],
		Location:     param["location"],
		Creator:      user,
		CreateTime:   now,
		Description:  "",
		CompressMode: "",
		Identifier:   GetTableIdentifier(param["tableType"], param["databaseName"], param["tableName"]),
		Owner:        strings.Join(roles, ","),
	}
	tid, err := o.Insert(&table)
	if err != nil {
		panic(err)
	}

	// 2. 版本
	newVersion := GetDefaultNewVersion(tid)
	v := VersionMetaData{
		TableId:    tid,
		Version:    newVersion,
		CreateTime: now,
	}
	_, _ = o.Insert(&v)

	var fields []FieldMetaData
	for {
		recoder, err := csvReader.Read()
		if err != nil {
			fmt.Println(err) // Will break on EOF.
			break
		}

		fields = append(fields, FieldMetaData{
			TableId:          tid,
			FieldName:        recoder[0],
			FieldType:        recoder[1],
			Description:      recoder[2],
			Version:          newVersion,
			ModificationTime: now,
		})
	}
	_, err = o.InsertMulti(100, fields)
	fmt.Println(err)
	return err == nil, tid

}

func (this ImportMetaData) importHiveSingle(param map[string]string, roles []string) (bool, int64) {
	dbName := param["databaseName"]
	tbName := param["tableName"]
	connId := param["importId"]

	hiveModel := HiveModel{}
	hiveTable := hiveModel.GetHiveProtoMeta(dbName, tbName)

	o := orm.NewOrm()
	now := utils.TimeFormat(time.Now())

	var tid int64 = 0
	if (dbName != "") && (tbName != "") {
		tid = this.InsertTableFromMeta(dbName, hiveTable, o, roles, connId)
		newVersion := this.InsertVersionFromMeta(tid, now, o)
		this.InsertFieldFromMeta(tid, newVersion, hiveTable, o)

	}

	return true, tid
}

func (this ImportMetaData) InsertTableFromMeta(dbName string,
	tb *adtl_platform.HiveTable,
	o orm.Ormer,
	roles []string,
	connId string) int64 {
	compress := tb.ExpendInfo["storageProperties"]
	dsId, _ := strconv.ParseInt(connId, 10, 64)
	table := TableMetaData{
		TableName:    tb.NameTable,
		TableType:    "hive",
		DatabaseName: dbName,
		Creator:      tb.ExpendInfo["owner"],
		CreateTime:   tb.ExpendInfo["createdTime"],
		Location:     tb.ExpendInfo["location"],
		CompressMode: compress,
		Identifier:   GetTableIdentifier("hive", dbName, tb.NameTable),
		Owner:        strings.Join(roles, ","),
		Dsid:         dsId,
	}
	tid, err := o.Insert(&table)
	if err != nil {
		panic(err)
	}
	return tid
}

func (this ImportMetaData) InsertVersionFromMeta(tid int64, nowTime string, o orm.Ormer) string {
	newVersion := GetDefaultNewVersion(tid)
	v := VersionMetaData{
		TableId:    tid,
		Version:    newVersion,
		CreateTime: nowTime,
	}
	_, _ = o.Insert(&v)
	return newVersion
}

func (this ImportMetaData) InsertFieldFromMeta(tid int64, newVersion string, tb *adtl_platform.HiveTable, o orm.Ormer) {
	hiveModel := HiveModel{}
	fields := hiveModel.ParseFields(tid, newVersion, tb)
	_, _ = o.InsertMulti(100, fields)
}

// jdbcUrl string, user string, passwd string, importId int64
func (this ImportMetaData) importMysqlSingle(conn ConnectionManager, dbName string, tbName string, roles []string) (bool, int64) {
	if dbName == "" {
		return false, 0
	}

	mysqlModel := MysqlModel{}
	return mysqlModel.ImportFieldsSingleTable(conn, dbName, tbName, roles)
}

func HandleOnlineRecordFieldDiff(increased *[]FieldMetaData, deleted *[]FieldMetaData) {
	o := orm.NewOrm()
	if *increased != nil {
		if _, err := o.InsertMulti(100, increased); err != nil {
			Logger.Error("%s", err.Error())
		}
	}
	if *deleted != nil {
		for _, field := range *deleted {
			o.Raw("DELETE FROM field_mata_data WHERE fid = ?", field.Fid)
		}
	}
}
