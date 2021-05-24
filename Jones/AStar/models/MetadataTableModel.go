package models

import (
	"AStar/utils"
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type TableMetaData struct {
	Tid              int64  `orm:"pk;auto" field:"tid,hidden"`
	TableName        string `orm:"size(255)" field:"表名"`  //表名
	TableType        string `orm:"size(255)" field:"表类型"` //表类型
	DatabaseName     string `orm:"size(255)" field:"所属库"` //所属数据库
	Location         string `orm:"size(255)" field:"文件存储位置,hidden"`
	Creator          string `orm:"size(100)" field:"创建者"`           //创建者
	CreateTime       string `orm:"size(100)" field:"创建时间"`          //创建时间
	Modifier         string `orm:"size(100)" field:"修改者,hidden"`    //修改者
	ModificationTime string `orm:"size(100)" field:"上次修改时间,hidden"` //修改时间
	Description      string `orm:"type(text)" field:"描述"`           // 表的描述信息
	CompressMode     string `orm:"size(255)" field:"压缩方式,hidden"`
	Identifier       string `orm:"size(100)" field:"标识符,hidden"`
	Dsid             int64  `orm:"size(11)" field:"数据源ID,hidden"`
	Owner            string `orm:"size(50)" field:"所有者,hidden"`
	PrivilegeRoles   string `orm:"-" field:"权限赋予,hidden"`
}

type GroupMetatableMapping struct {
	Gwid      int64  `orm:"pk;auto"`
	GroupName string `orm:"size(50)"`
	TableId   string `orm:"size(50)"`
	Valid     int    `orm:"default(1)"`
	ReadFlag  int64  `orm:"size(1)"`
	ExecFlag  int64  `orm:"size(1)"`
}

//func (u *GroupMetatableMapping) TableIndex() [][]string {
//	return [][]string{
//		[]string{"GroupName", "TableId"},
//	}
//}
func (u *GroupMetatableMapping) TableUnique() [][]string {
	return [][]string{
		[]string{"GroupName", "TableId"},
	}
}

type FrontMetaAddTable struct {
	TableName    string `json:"tableName"`    //表名
	TableType    string `json:"tableType"`    //表类型
	DatabaseName string `json:"databaseName"` //数据库名
	Creator      string `json:"creator"`      //创建者
	Modifier     string `json:"modifier"`     //修改者
}

type FrontMetaGetAllTableRes struct {
	TableName    string `json:"tablename"    orm:"column(table_name)"`    //表名
	DatabaseName string `json:"databasename" orm:"column(database_name)"` //数据库名
}

type TableTypes struct {
	Tid       int64  `orm:"pk;auto"`
	TypeName  string `orm:"size(255)"`
	TypeColor string `orm:"size(255)"`
}

type TreeViewNode struct {
	Id           int
	NodeName     string
	NodeChildren []TreeViewNode
}

func (this *TableMetaData) GetTableById() {
	o := orm.NewOrm()
	if err := o.Read(this); err != nil {
		panic(err)
	}

}

type MetaTableAccessView struct {
	Tid        int64
	TableName  string
	ReadAccess bool
	ExecAccess bool
}

type MetaTableFullInfo struct {
	Dsid           int
	ConnName       string
	DatasourceName string
	DatabaseName   string
	TableName      string
	Tid            int
}

// update table_meta_data set identifier= MD5(concat_ws('-',table_type,database_name,table_name))
func GetTableIdentifier(tableType string, databaseName string, tableName string) string {
	str := fmt.Sprintf("%s-%s-%s", tableType, databaseName, tableName)
	data := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", data)
}

func CheckDuplicateTable(tableType string, databaseName string, tableName string) bool {
	tid := GetTableId(tableType, databaseName, tableName)
	return tid != 0 // 不等于0是重复
}

func GetTableId(tableType string, databaseName string, tableName string) int64 {
	checkIdentifier := GetTableIdentifier(tableType, databaseName, tableName)
	o := orm.NewOrm()
	var t TableMetaData
	_ = o.Raw("select * from table_meta_data where identifier = ?", checkIdentifier).QueryRow(&t)
	return t.Tid
}

func GetAllTableTypesMap() map[string]string {
	types := GetAllTableTypes()
	maps := make(map[string]string)
	for _, tt := range types {
		maps[tt.TypeName] = tt.TypeColor
	}
	return maps
}

func GetAllTableTypes() []TableTypes {
	o := orm.NewOrm()
	var types []TableTypes
	sql := "SELECT * FROM table_types"
	_, err := o.Raw(sql).QueryRows(&types)
	if err != nil {
		panic(err)
	}
	return types
}

func (this TableMetaData) InsertMetaTable() int64 {
	taskId, err := orm.NewOrm().Insert(&this)
	if err != nil {
		Logger.Error(err.Error())
	}
	return taskId
}

func (this *TableMetaData) GetTableCount(param map[string]string, userName string) int64 {
	o := orm.NewOrm()
	var totalCnt int64
	roles := Enforcer.GetRolesForUser(userName)
	flag := hasCondition(param)
	var roleSqlList []string
	var paramList []interface{}
	for _, roleName := range roles {
		//var cnt int64
		if !flag {
			sql := `SELECT COUNT(1) FROM (SELECT * FROM table_meta_data t1
left join 
(select table_id from group_metatable_mapping where group_name=? and (read_flag=1 or exec_flag=1)) t2
on t1.tid = t2.table_id
where t2.table_id is not null) p1`
			roleSqlList = append(roleSqlList, sql)
			paramList = append(paramList, roleName)
			//_ = o.Raw(sql, roleName).QueryRow(&cnt)
			//totalCnt = totalCnt + cnt
		} else {
			conditionSql := " WHERE"
			for k, v := range param {
				if k == "curPage" || k == "pageSize" || v == "" {
					continue
				}
				if k == "tableNameLike" {
					conditionSql = fmt.Sprintf("%s (table_name like '%%%s%%' or database_name like '%%%s%%' or description like '%%%s%%') AND",
						conditionSql, v, v, v)
				} else {
					conditionSql = fmt.Sprintf("%s %s='%s' AND", conditionSql, k, v)
				}
			}
			conditionSql = conditionSql[0 : len(conditionSql)-4]
			sql := fmt.Sprintf(`SELECT count(1) FROM (SELECT * FROM table_meta_data t1
left join 
(select table_id from group_metatable_mapping where group_name=? and (read_flag=1 or exec_flag=1)) t2
on t1.tid = t2.table_id
where t2.table_id is not null) t3 %s `, conditionSql)
			roleSqlList = append(roleSqlList, sql)
			paramList = append(paramList, roleName)
			//_ = o.Raw(sql, roleName).QueryRow(&cnt)
			//totalCnt = totalCnt + cnt
		}
	}
	uniqSql := strings.Join(roleSqlList, "\nunion\n")
	_ = o.Raw(uniqSql, paramList).QueryRow(&totalCnt)
	return totalCnt
}

// 查询所有类型
func (this *TableMetaData) GetAllTableTypes() []string {
	o := orm.NewOrm()
	var types []string
	sql := "SELECT table_type FROM table_meta_data GROUP BY table_type"
	_, err := o.Raw(sql).QueryRows(&types)
	if err != nil {
		panic(err)
	}
	return types
}

// 根据类型查数据库
func (this *TableMetaData) GetDataBaseByType(tableType string) []string {
	o := orm.NewOrm()
	var dataBases []string
	sql := "SELECT database_name FROM table_meta_data WHERE table_type=? GROUP BY database_name"
	_, err := o.Raw(sql, tableType).QueryRows(&dataBases)
	if err != nil {
		panic(err)
	}
	return dataBases
}

// 根据数据库查表
func (this *TableMetaData) GetTablesByDataBase(dateBase string, tableType string) []TableMetaData {
	o := orm.NewOrm()
	var tables []TableMetaData
	sql := "SELECT * FROM table_meta_data WHERE database_name=? and table_type=?"
	_, err := o.Raw(sql, dateBase, tableType).QueryRows(&tables)
	if err != nil {
		panic(err)
	}
	return tables
}

func hasCondition(conditions map[string]string) bool {
	count := 0
	for _, v := range conditions {
		if v != "" {
			count++
		}
	}
	return count > 2
}

type TableMetaDataShow struct {
	Tid              int64  `field:"tid,hidden"`
	TableName        string `field:"表名"`  //表名
	TableType        string `field:"表类型"` //表类型
	DatasourceName   string `field:"数据源"`
	DatabaseName     string `field:"所属库"` //所属数据库
	Location         string `field:"文件存储位置,hidden"`
	Creator          string `field:"创建者"`           //创建者
	CreateTime       string `field:"创建时间"`          //创建时间
	Modifier         string `field:"修改者,hidden"`    //修改者
	ModificationTime string `field:"上次修改时间,hidden"` //修改时间
	Description      string `field:"描述"`            // 表的描述信息
	CompressMode     string `field:"压缩方式,hidden"`
	Identifier       string `field:"标识符,hidden"`
	Dsid             int64  `field:"数据源ID,hidden"`
	Owner            string `field:"所有者,hidden"`
	IsOwner          bool   `field:"当前用户是否为所有者,hidden"`
	PrivilegeRoles   string `field:"权限赋予,hidden"`
}

// 根据表类型和数据库查询存在什么表
func (this *TableMetaData) GetTablesByConditions(conditions map[string]string, userName string) utils.TableData {

	pageSize, _ := strconv.ParseInt(conditions["pageSize"], 10, 64)
	curPage, _ := strconv.ParseInt(conditions["curPage"], 10, 64)
	curPage = curPage - 1 // 前台控件的page号从1开始

	//o := orm.NewOrm()
	var totallTables []TableMetaDataShow

	roles := Enforcer.GetRolesForUser(userName)
	flag := hasCondition(conditions)
	var roleSqlList []string
	var paramList []interface{}
	for _, roleName := range roles {
		//var tables []TableMetaDataShow
		if !flag {
			sql := `SELECT t1.*, t3.datasource_name FROM table_meta_data t1
left join 
(select table_id from group_metatable_mapping where group_name=? and (read_flag=1 or exec_flag=1)) t2
on t1.tid = t2.table_id
LEFT JOIN 
connection_manager t3 ON t1.dsid = t3.dsid
where t2.table_id is not null
`
			roleSqlList = append(roleSqlList, sql)
			paramList = append(paramList, roleName)
		} else {
			conditionSql := " WHERE"
			for k, v := range conditions {
				if k == "curPage" || k == "pageSize" || v == "" {
					continue
				}
				if k == "tableNameLike" {
					conditionSql = fmt.Sprintf("%s (table_name like '%%%s%%' or database_name like '%%%s%%' or description like '%%%s%%') AND",
						conditionSql, v, v, v)
				} else {
					conditionSql = fmt.Sprintf("%s %s='%s' AND", conditionSql, k, v)
				}
			}
			conditionSql = conditionSql[0 : len(conditionSql)-4]
			sql := fmt.Sprintf(`SELECT * FROM (SELECT t1.*, p1.datasource_name FROM table_meta_data t1
left join 
(select table_id from group_metatable_mapping where group_name=? and (read_flag=1 or exec_flag=1)) t2
on t1.tid = t2.table_id
LEFT JOIN 
connection_manager p1 ON t1.dsid = p1.dsid
where t2.table_id is not null) t3 %s 
`, conditionSql)
			roleSqlList = append(roleSqlList, sql)
			paramList = append(paramList, roleName)
		}
	}
	limitSql := "\norder by tid desc limit ? offset ?"
	unionSql := strings.Join(roleSqlList, "\n union \n")

	culminateSql := fmt.Sprintf(`SELECT *, %s as is_owner FROM (%s) pp1`, MakeIsOwnerSql(roles), unionSql+limitSql)

	utils.Logger.Info("表元数据最终查询的sql:\n%s", culminateSql)

	if _, err := orm.NewOrm().Raw(culminateSql, paramList, pageSize, pageSize*curPage).QueryRows(&totallTables); err != nil {
		panic(err.Error())
	}

	// transView
	var bodys []map[string]interface{}
	for _, data := range totallTables {
		// 转成map
		tmpData := utils.GetSnakeStrMapOfStruct(data)
		bodys = append(bodys, tmpData)
	}

	// 返回数据

	tsb := utils.TableData{Data: bodys}
	tsb.GetSnakeTitles(TableMetaDataShow{})

	return tsb
}

func MakeIsOwnerSql(roles []string) string {
	findSqlItem := []string{}
	for _, role := range roles {
		findSqlItem = append(findSqlItem, fmt.Sprintf("find_in_set('%s',OWNER)", role))
	}
	findSql := strings.Join(findSqlItem, " + ")
	findSql = fmt.Sprintf("(%s) > 0", findSql)
	return findSql
}

func (this *TableMetaData) AddTables(params []map[string]string, user string) {
	o := orm.NewOrm()
	now := utils.TimeFormat(time.Now())
	roles := Enforcer.GetRolesForUser(user)
	for _, param := range params {
		dsid, _ := strconv.ParseInt(param["importId"], 10, 64)
		t := TableMetaData{
			TableName:        param["tableName"],
			TableType:        param["tableType"],
			DatabaseName:     param["database"],
			Creator:          user,
			CreateTime:       now,
			Modifier:         user,
			ModificationTime: now,
			Description:      param["description"],
			Location:         param["location"],
			CompressMode:     param["compressMode"],
			Dsid:             dsid,
			Owner:            strings.Join(roles, ","),
		}
		tid, err := o.Insert(&t)
		if err != nil {
			panic(err)
		}

		// 创建表时增加一个默认版本
		prefix := "v0"
		suffix := time.Now().Format("20060102-1504")
		initVersion := fmt.Sprintf("%s_%s", prefix, suffix)

		versionData := VersionMetaData{
			TableId:    tid,
			Version:    initVersion,
			CreateTime: now,
		}
		_, _ = o.Insert(&versionData)

		// 加入默认权限
		if tid != 0 {
			privilegeRoles := strings.Split(param["privilegeRoles"], ",")
			uniqRoles := make(map[string]string)
			for _, r := range privilegeRoles {
				if r != "" {
					uniqRoles[r] = ""
				}
			}
			for _, r := range roles {
				uniqRoles[r] = ""
			}

			var roleSet []string // 去重后的集合
			for k, _ := range uniqRoles {
				roleSet = append(roleSet, k)
			}
			AddTablePrivilegeMapping(tid, roleSet)
		}
	}
}

func (this *TableMetaData) DeleteTable(tid int64) {
	o := orm.NewOrm()
	_, err := o.Raw("DELETE FROM table_meta_data WHERE tid=?", tid).Exec()
	if err != nil {
		panic(err)
	}
	_, err = o.Raw("DELETE FROM field_meta_data WHERE table_id=?", tid).Exec()
	if err != nil {
		panic(err)
	}
	_, _ = o.Raw("DELETE FROM version_meta_data WHERE table_id=?", tid).Exec()
	_, _ = o.Raw("DELETE FROM group_metatable_mapping WHERE table_id=?", tid).Exec()
}

func (this *TableMetaData) UpdateTable(roles []string, privileges []string) {
	//Update 默认更新所有的字段
	o := orm.NewOrm()
	_, err := o.Update(this, "TableName",
		"TableType",
		"DatabaseName",
		"Modifier",
		"ModificationTime",
		"Description",
		"Location",
		"CompressMode",
		"Creator",
		"CreateTime")
	if err != nil {
		panic(err.Error())
	}
	uniqRoles := make(map[string]string)
	for _, r := range privileges {
		if r != "" {
			uniqRoles[r] = ""
		}
	}
	for _, r := range roles {
		uniqRoles[r] = ""
	}

	var roleSet []string // 去重后的集合
	for k, _ := range uniqRoles {
		roleSet = append(roleSet, k)
	}
	UpdateTablePrivilegeMapping(this.Tid, roleSet)
}

func (this *TableMetaData) GetConnectionId(connType string) []MetaDataImport {
	var imports []MetaDataImport
	conns := GetAllConn(connType)
	for i, _ := range conns {
		imports = append(imports, MetaDataImport{
			ImportId:      conns[i].Dsid,
			ImportName:    conns[i].DatasourceName,
			ConnectionUrl: conns[i].Url,
		})
	}
	return imports
}

// 查询用户下有权限的数据表
func (this TableMetaData) GetCheckedMetaTable(groupName string) []MetaTableAccessView {
	o := orm.NewOrm()
	var views []MetaTableAccessView
	sql := `
			SELECT
				t1.tid,
				t1.table_name,
			CASE
					
					WHEN t2.table_id IS NULL THEN
					0 ELSE t2.read_flag 
					END AS read_access,
			CASE
					
					WHEN t2.table_id IS NULL THEN
					0 ELSE t2.exec_flag 
					END AS exec_access 
			FROM
				( SELECT tid, table_name FROM table_meta_data ) t1
				LEFT JOIN 
				( SELECT table_id, read_flag, exec_flag FROM group_metatable_mapping WHERE  group_name = ? ) t2 
				ON t1.tid = t2.table_id 
			ORDER BY
				t1.tid ASC
`
	_, err := o.Raw(sql, groupName).QueryRows(&views)
	if err != nil {
		panic(err.Error())
	}
	return views
}

// 查询元数据表的 连接类型，连接名，数据库名，表名
func (this TableMetaData) GetMetaTableFullInfoList(role string) []MetaTableFullInfo {
	o := orm.NewOrm()
	var views []MetaTableFullInfo
	sql := `
			SELECT 
				cm.dsid,
				cm.conn_type AS conn_name,
				cm.datasource_name,
				IFNULL(tmd.database_name,'') AS database_name,
				IFNULL(tmd.table_name,'') AS table_name,
				IFNULL(tmd.tid,0) AS tid 
			FROM
				connection_manager cm
			LEFT JOIN table_meta_data tmd
			on  tmd.dsid = cm.dsid
			ORDER BY conn_type,datasource_name,database_name DESC


	`
	_, err := o.Raw(sql).QueryRows(&views)
	if err != nil {
		panic(err.Error())
	}
	return views
}

//更新角色和工单列表
// worksheetids: 有权限的tid列表, curType: 权限类型
func (this TableMetaData) UpdateMetaTableMapping(groupName string, worksheetIds []string, curType string) {
	o := orm.NewOrm()
	update_column := ""
	keep_column := ""
	if curType == "read" {
		update_column = "read_flag"
		keep_column = "exec_flag"
	} else {
		update_column = "exec_flag"
		keep_column = "read_flag"
	}

	sql := fmt.Sprintf(" update group_metatable_mapping SET %s = false WHERE group_name = '%s' ", update_column, groupName)
	_, err := o.Raw(sql).Exec()
	if err == nil {
		Logger.Info("update group_metatable_mapping")
	} else {
		panic(err)
	}

	for _, tid := range worksheetIds {
		int_tid, _ := strconv.Atoi(tid)
		if int_tid >= 1000001 {
			continue
		}
		var keep_column_value int

		//当replace into 的数据不存在时， 要keep的字段的数据select会为null， 默认保留keep_column_value 的原始值
		//设置为 0 就是当 null 的时候采用0 插入要保留的字段值
		// 当不为 null 有该条记录的时候， 采用原始记录的值。
		keep_column_value = 0

		sql = fmt.Sprintf("  select %s from  group_metatable_mapping where group_name = '%s' and table_id = '%s' ", keep_column, groupName, tid)
		_ = o.Raw(sql).QueryRow(&keep_column_value)

		sql = fmt.Sprintf("  Replace into group_metatable_mapping(group_name,table_id,%s,%s)  values('%s','%s',1,%d)  ", update_column, keep_column, groupName, tid, keep_column_value)
		_, err := o.Raw(sql).Exec()
		if err == nil {
			Logger.Info("Replace into group_metatable_mapping ")
		} else {
			panic(err)
		}
	}

}

func AddTablePrivilegeMapping(tid int64, roles []string) {
	tableMappings := make([]GroupMetatableMapping, 0)
	for _, role := range roles {
		mapping := GroupMetatableMapping{
			GroupName: role,
			TableId:   fmt.Sprintf("%d", tid),
			Valid:     1,
			ReadFlag:  1,
			ExecFlag:  1,
		}
		tableMappings = append(tableMappings, mapping)
	}

	o := orm.NewOrm()
	_, _ = o.InsertMulti(100, tableMappings)
}

func UpdateTablePrivilegeMapping(tid int64, roles []string) {
	o := orm.NewOrm()
	_, _ = o.Raw("DELETE FROM group_metatable_mapping WHERE table_id = ?", fmt.Sprintf("%d", tid)).Exec()
	AddTablePrivilegeMapping(tid, roles)
}

func GetTablePrivilegeByTid(tid int64) []string {
	o := orm.NewOrm()
	mappings := make([]GroupMetatableMapping, 0)
	sql := "SELECT * FROM group_metatable_mapping WHERE table_id = ? AND read_flag = 1"
	_, _ = o.Raw(sql, fmt.Sprintf("%d", tid)).QueryRows(&mappings)
	roles := make([]string, 0)
	for _, mapping := range mappings {
		roles = append(roles, mapping.GroupName)
	}
	return roles
}
