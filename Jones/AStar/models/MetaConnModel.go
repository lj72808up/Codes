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

type ConnectionManager struct {
	Dsid           int64  `orm:"pk;auto"`
	DatasourceName string `orm:"size(50)"`
	ConnType       string `orm:"size(20)"`
	User           string `orm:"size(50)"`
	Passwd         string `orm:"size(50)"`
	Host           string `orm:"size(50)"`
	Port           string `orm:"size(10)"`
	Database       string `orm:"size(50)"`
	Encode         string `orm:"size(20)"`
	Url            string `orm:"size(300)"`
	Owner          string `orm:"size(500)"`
	PrivilegeRoles string `orm:"-"`
	Redashid       int `orm:"size(20)"`
}

type GroupConnMapping struct {
	Gcid  int64 `orm:"pk;auto"`
	Dsid  int64
	Role  string `orm:"size(50)"`
	Valid bool   `orm:"default(true)"`
}

func (this *ConnectionManager) GetById(id int64) {
	o := orm.NewOrm()
	this.Dsid = id
	_ = o.Read(this)
}

//创建Redash连接
func (this *ConnectionManager) AddRedashConn() (bool,int){
	res_flag := false
	res_redashid := -1

	url := "http://inner-redash.adtech.sogou/api/data_sources"
	param := ""

	if this.ConnType == "clickhouse" {
		// CK HTTP Post 使用 9090

		param = fmt.Sprintf(`
			{
				"type": "%s",
				"options": {
					"url": "http://%s:%s",
					"password": "%s",
					"user": "%s",
					"dbname": "%s"
				},
				"name": "%s"
			}
		`, this.ConnType, this.Host, "9090",this.Passwd, this.User, this.Database, this.DatasourceName)
	} else if this.ConnType == "mysql" {
		param = fmt.Sprintf(`
			{
				"type": "%s",
				"options": {
					"host": "%s",
					"port": %s,
					"passwd": "%s",
					"user": "%s",
					"db": "%s"
				},
				"name": "%s"
			}
		`, this.ConnType, this.Host, this.Port,this.Passwd, this.User, this.Database, this.DatasourceName)
	} else {
		return true,0
	}

	fmt.Println(param)
	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body([]byte(param))

	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
	}
	var res map[string]interface{}
	_ = json.Unmarshal(bytes, &res)

	res_id, exist := res["id"]

	if exist {
		res_flag = exist
		res_redashid = int(res_id.(float64))
	} else {
		res_flag = exist
		res_redashid = -1
	}
	return res_flag,res_redashid
}


func (this *ConnectionManager) Add(roles []string, privilegeRoles []string) {
	//0. 先创建redash 连接
	suc_flag,redashid := this.AddRedashConn()
	if !suc_flag {
		panic("Create Redash DataSource Failed")
	}
	this.Redashid = redashid

	// 1. 插入连接
	o := orm.NewOrm()
	if this.ConnType == "mysql" {
		this.Url = utils.MakeJdbcUrl(this.Host, this.Port, this.Database, this.Encode)
	} else if this.ConnType == "clickhouse" {
		this.Url = utils.MakeClickHouseUrl(this.Host, this.Port, this.Database)
	}
	this.Owner = strings.Join(roles, ",") // owner使用逗号连接形成字符串
	dsId, err := o.Insert(this)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("增加connId:", dsId)
	// 2.增加权限,将 privilegeRoles 和 roles 加入到 GroupConnMapping 表中
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
	var connMappings []GroupConnMapping // 去重后的可使用表的角色
	for _, role := range roleSet {
		connMapping := GroupConnMapping{
			Dsid:  dsId,
			Role:  role,
			Valid: true,
		}
		connMappings = append(connMappings, connMapping)
	}
	_, _ = o.InsertMulti(100, connMappings)
}



//更新Redash连接
func (this *ConnectionManager) UpdateRedashConn() (bool){
	res_flag := false

	url :=fmt.Sprintf( "http://inner-redash.adtech.sogou/api/data_sources/%d" ,this.Redashid)
	param := ""

	if this.ConnType == "clickhouse" {
		param = fmt.Sprintf(`
			{
				"type": "%s",
				"options": {
					"url": "%s:%s",
					"password": "%s",
					"user": "%s",
					"dbname": "%s"
				},
				"name": "%s"
			}
		`, this.ConnType, this.Host, this.Port,this.Passwd, this.User, this.Database, this.DatasourceName)
	} else if this.ConnType == "mysql" {
		param = fmt.Sprintf(`
			{
				"type": "%s",
				"options": {
					"host": "%s",
					"port": %s,
					"passwd": "%s",
					"user": "%s",
					"db": "%s"
				},
				"name": "%s"
			}
		`, this.ConnType, this.Host, this.Port,this.Passwd, this.User, this.Database, this.DatasourceName)
	} else {
		return true
	}

	fmt.Println(param)
	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body([]byte(param))

	bytes, error := req.Bytes()
	if error != nil {
		Logger.Error(error.Error())
	}
	var res map[string]interface{}
	_ = json.Unmarshal(bytes, &res)

	_, exist := res["id"]

	if exist {
		res_flag = exist
	} else {
		res_flag = exist
	}
	return res_flag
}

func (this *ConnectionManager) Update(selfRoles []string) {
	o := orm.NewOrm()

	// 0. 更新redash 连接
	sql := fmt.Sprintf("SELECT redashid from connection_manager where dsid=%d ;",this.Dsid )
	var redashid int
	o.Raw(sql).QueryRow(&redashid)

	var suc_flag bool
	if redashid == 0{
		suc_flag,redashid = this.AddRedashConn()
		this.Redashid = redashid
	} else {
		this.Redashid = redashid
		suc_flag = this.UpdateRedashConn()
	}
	if !suc_flag {
		panic("Update Redash DataSource Failed")
	}



	if this.ConnType == "mysql" {
		this.Url = utils.MakeJdbcUrl(this.Host, this.Port, this.Database, this.Encode)
	} else if this.ConnType == "clickhouse" {
		this.Url = utils.MakeClickHouseUrl(this.Host, this.Port, this.Database)
	}
	if _, err := o.Update(this); err != nil {
		panic(err)
	}

	privilegeRoles := strings.Split(this.PrivilegeRoles, ",")
	uniqRoles := make(map[string]string)
	for _, r := range privilegeRoles {  // 页面勾选的角色
		if r != "" {
			uniqRoles[r] = ""
		}
	}
	for _, r := range selfRoles {       // 自己的角色
		uniqRoles[r] = ""
	}
	var roleSet []string // 去重后的集合
	for k, _ := range uniqRoles {
		roleSet = append(roleSet, k)
	}
	_, _ = o.Raw("UPDATE group_conn_mapping SET valid=false WHERE dsid=?", this.Dsid).Exec()
	var connMappings []GroupConnMapping // 去重后的可使用表的角色
	for _, role := range roleSet {
		connMapping := GroupConnMapping{
			Dsid:  this.Dsid,
			Role:  role,
			Valid: true,
		}
		connMappings = append(connMappings, connMapping)
	}
	_, _ = o.InsertMulti(100, connMappings)
}

func (this *ConnectionManager) TestConn(connType string, user string, passwd string, database string, host string, port string) error {
	var err error
	switch connType {
	case "mysql":
		aliasName := fmt.Sprintf("%s-%d", database, time.Now().Unix())
		url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, passwd, host, port, database)
		fmt.Printf("数据源: %s\n", url)

		// SetMaxIdleConns, SetMaxOpenConns: 最大空闲0,最大连接1,让orm执行一次查询后就断开连接
		connErr := orm.RegisterDataBase(aliasName, "mysql", url, 0, 1)
		if connErr != nil {
			panic(connErr)
		}
		o1 := orm.NewOrm()
		_ = o1.Using(aliasName)

		// 1. 执行查询
		sql := fmt.Sprintf("SELECT 1")
		var res int64
		err = o1.Raw(sql).QueryRow(&res)
	case "clickhouse":
		url := fmt.Sprintf("tcp://%s:%s?user=%s&password=%s&debug=true", host, port, user, passwd)
		dbConn, err := sqlx.Open("clickhouse", url)
		if err != nil {
			fmt.Println(err)
			return err
		}
		var res []int
		err = dbConn.Select(&res, "select 1")
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return err
}

func GetAllConn(connType string) []ConnectionManager {
	var conns []ConnectionManager
	builder, _ := orm.NewQueryBuilder("mysql")
	builder.Select("*").From("connection_manager")
	if connType != "" {
		builder.Where("conn_type = ?")
	}
	sql := builder.String()
	fmt.Println(sql)
	o := orm.NewOrm()
	if connType != "" {
		_, _ = o.Raw(sql, connType).QueryRows(&conns)
	} else {
		_, _ = o.Raw(sql).QueryRows(&conns)
	}
	return conns
}

// owner是以","分隔的角色名数组
func GetConnByOwner(connType string, ownerRoles []string) []ConnectionManager {
	var conns []ConnectionManager
	var sqlItems []string
	for _, role := range ownerRoles {
		sql := fmt.Sprintf(`SELECT t1.* FROM connection_manager t1
LEFT JOIN
(SELECT * FROM group_conn_mapping WHERE role = '%s' AND valid = TRUE) t2 
ON t1.dsid = t2.dsid
WHERE t2.dsid IS NOT NULL `, role)
		if connType != "" {
			sql = fmt.Sprintf("%s AND t1.conn_type = '%s'", sql, connType)
		}
		sqlItems = append(sqlItems, sql)
	}
	sql := strings.Join(sqlItems, "\nUNION\n")

	o := orm.NewOrm()
	_, _ = o.Raw(sql).QueryRows(&conns)
	return conns
}

// owner是以","分隔的角色名数组
func GetConnOnlyByOwner(connType string, ownerRoles []string) []ConnectionManager {
	var conns []ConnectionManager
	var sqlItems []string
	for _, role := range ownerRoles {
		sql := fmt.Sprintf("SELECT * FROM connection_manager WHERE find_in_set('%s',OWNER) > 0", role)
		if connType != "" {
			sql = fmt.Sprintf("%s AND conn_type = '%s'", sql, connType)
		}
		sqlItems = append(sqlItems, sql)
	}
	sql := strings.Join(sqlItems, "\nUNION\n")

	o := orm.NewOrm()
	_, _ = o.Raw(sql).QueryRows(&conns)
	return conns
}

func GetConnByPrivilege(connType string, roles []string) []ConnectionManager {
	conns := make([]ConnectionManager, 0)
	rolesWildCard := utils.MakeRolePlaceHolder(roles)
	o := orm.NewOrm()
	sql := fmt.Sprintf(`SELECT t1.*  FROM
(SELECT * from connection_manager WHERE conn_type = ?) t1
LEFT JOIN (SELECT dsid FROM group_conn_mapping WHERE role in %s AND valid = 1 group by dsid)t2
ON t1.dsid = t2.dsid
WHERE t2.dsid is not null
`, rolesWildCard)
	_, _ = o.Raw(sql, connType, roles).QueryRows(&conns)
	return conns
}

func GetPrivilegeRolesByConnId(dsId int64) []string {
	o := orm.NewOrm()
	roles := make([]string, 0)
	_, _ = o.Raw("SELECT role FROM group_conn_mapping WHERE dsid=? and valid=true group by role", dsId).QueryRows(&roles)
	return roles
}
