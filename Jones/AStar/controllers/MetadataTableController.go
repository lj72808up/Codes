package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type MetadataTableController struct {
	BaseController
}

var MetaDataNameSpace = "metaData"

// @Title 更新当前版本元数据
// @Summary 更新当前版本元数据
// @Description 更新当前版本元数据
// @Success 200 string "suc num"
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /updateTableByVersion [post]
func (this *MetadataTableController) UpdateTableByVersion() {
	var param map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &param)

	version := param["version"]
	tableId := param["tableId"]

	defer func() {
		if r := recover(); r != nil {
			utils.PrintStackTrace()
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()

	o := orm.NewOrm()
	var table models.TableMetaData
	_ = o.Raw("SELECT * FROM table_meta_data WHERE tid = ?", tableId).QueryRow(&table)

	switch table.TableType {
	case "hive":
		model := models.HiveModel{}
		model.DiffField(table, version)
	case "mysql":
		model := models.MysqlModel{}
		model.DiffField(table, version)
	case "clickhouse":
		model := models.CKModel{}
		model.DiffField(table, version)
	}

	this.Ctx.WriteString(utils.VueResponse{
		Res:  true,
		Info: "更新成功",
	}.String())
}

// @Title AddTableMetaData
// @Summary 添加表级元数据
// @Description 添加表级元数据
// @Param para body models.FrontMetaAddTable true "表数据元信息"
// @Success 200 string "suc num"
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /addtable [post]
func (this *MetadataTableController) AddTableMetaData() {
	var add_tb_info models.FrontMetaAddTable
	json.Unmarshal(this.Ctx.Input.RequestBody, &add_tb_info)

	var meta_tb models.TableMetaData
	meta_tb.TableName = add_tb_info.TableName
	meta_tb.DatabaseName = add_tb_info.DatabaseName
	meta_tb.Creator = add_tb_info.Creator
	meta_tb.Modifier = add_tb_info.Modifier
	meta_tb.CreateTime = utils.TimeFormat(time.Now())
	meta_tb.ModificationTime = utils.TimeFormat(time.Now())

	jobId := meta_tb.InsertMetaTable()
	res := fmt.Sprintln(jobId)
	this.Ctx.WriteString("finish:" + res)
}

// @Title GetAllTableMetaData
// @Summary 查看元数据列表汇总
// @Description 查看元数据列表汇总
// @Success 200 {object} []models.FrontMetaGetAllTableRes
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getalltable [get]
func (this *MetadataTableController) GetAllTableMetaData() {
	var tbs []models.FrontMetaGetAllTableRes
	orm.NewOrm().Raw("select table_name,database_name from table_meta_data").QueryRows(&tbs)

	marshal, err := json.Marshal(tbs)
	if err != nil {
		this.Ctx.WriteString("参数转换异常")
		return
	}
	this.Ctx.WriteString(string(marshal))
}

// @Title GetAllDBMetaData
// @Summary 查看存在的所有数据库
// @Description 查看存在的所有数据库
// @Success 200 array[string]
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getalldb [get]
func (this *MetadataTableController) GetAllDBMetaData() {
	var dbs []string
	orm.NewOrm().Raw("select distinct database_name from table_meta_data").QueryRows(&dbs)

	marshal, err := json.Marshal(dbs)
	if err != nil {
		this.Ctx.WriteString("参数转换异常")
		return
	}
	this.Ctx.WriteString(string(marshal))
}

// @Title GetTablesOfDb
// @Param dbname path string true "database_name"
// @Summary 查看某个数据库下的所有表
// @Description  查看某个数据库下的所有表
// @Success 200 string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getdbtable/:dbname [get]
func (this *MetadataTableController) GetTablesOfDb() {
	dbname := this.Ctx.Input.Params()[":dbname"]
	var tbs []string
	orm.NewOrm().Raw("select distinct table_name from table_meta_data where database_name='" + dbname + "'").QueryRows(&tbs)

	marshal, err := json.Marshal(tbs)
	if err != nil {
		this.Ctx.WriteString("参数转换异常")
		return
	}
	this.Ctx.WriteString(string(marshal))
}

// @Title GetTableByConditions
// @Summary 根据条件查询表名
// @Description 根据条件查询表名
// @Success 200 array[string]
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getTableByConditions [post]
func (this *MetadataTableController) GetTableByConditions() {
	var param map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &param)
	var meta models.TableMetaData
	tables := meta.GetTablesByConditions(param, this.userName)
	res, _ := json.Marshal(tables)
	this.Ctx.WriteString(string(res))
}

// @Title PostTables
// @Summary 批量添加表级元数据
// @Description 批量添加表级元数据
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /postTables [post]
func (this *MetadataTableController) PostTables() {
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var params []map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)

	if models.CheckDuplicateTable(params[0]["tableType"], params[0]["database"], params[0]["tableName"]) {
		response := utils.VueResponse{
			Res:  false,
			Info: "该表的元数据已经存在,请确保类型,库,表名称不同",
		}
		this.Ctx.WriteString(response.String())
	} else {
		var meta models.TableMetaData
		meta.AddTables(params, this.userName)
		response := utils.VueResponse{
			Res:  true,
			Info: "表创建成功",
		}
		this.Ctx.WriteString(response.String())
	}

}

// @Title DeleteTableById
// @Summary 删除某个表
// @Description 删除某个表
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /table/delete/:tid [get]
func (this *MetadataTableController) DeleteTableById() {
	param := this.Ctx.Input.Params()[":tid"]
	tableId, _ := strconv.ParseInt(param, 10, 64)
	var meta models.TableMetaData
	meta.DeleteTable(tableId)
	this.Ctx.WriteString("success")
}

// @Title UpdateTableById
// @Summary 更新表的元数据
// @Description 更新表的元数据
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /table/update [post]
func (this *MetadataTableController) UpdateTableById() {
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var table models.TableMetaData
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &table)

	table.Modifier = this.userName
	table.ModificationTime = utils.TimeFormat(time.Now())

	privileges := strings.Split(table.PrivilegeRoles, ",")
	roles := models.Enforcer.GetRolesForUser(this.userName)
	table.UpdateTable(roles, privileges)
	this.Ctx.WriteString("success")
}

// @Title GetTableNumber
// @Summary 查看有多少个表
// @Description 查看有多少个表
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /table/getNums [post]
func (this *MetadataTableController) GetTableNumber() {
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)

	model := models.TableMetaData{}
	count := model.GetTableCount(params, this.userName)
	res, _ := json.Marshal(map[string]int64{
		"count": count,
	})
	this.Ctx.WriteString(string(res))
}

// @Title GetAllMysqlImports
// @Summary 获取所有mysql的导入id
// @Description 获取所有mysql的导入id
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /tableImport/getConnections [get]
func (this *MetadataTableController) GetConnectionIds() {
	var meta models.TableMetaData
	connType := this.GetString("connType")
	importIds := meta.GetConnectionId(connType)
	res, _ := json.Marshal(importIds)
	this.Ctx.WriteString(string(res))
}

// @Title GetTablesById
// @Summary 根据Id查看某个表
// @Description  根据Id查看某个表
// @Success 200 string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getTableById/:tid [get]
func (this *MetadataTableController) GetTablesById() {
	param := this.Ctx.Input.Params()[":tid"]
	tid, _ := strconv.ParseInt(param, 10, 64)
	tableMetaData := models.TableMetaData{
		Tid: tid,
	}
	tableMetaData.GetTableById()
	res, _ := json.Marshal(tableMetaData)
	this.Ctx.WriteString(string(res))
}

// @Title GetOwnerEdit
// @Summary 查看是否可以编辑
// @Description 查看是否可以编辑
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getOwnerEdit [get]
func (this *MetadataTableController) GetOwnerEdit() {
	// 通过casbin,判断/jones/jdbc_authority (可写权限的角色要配置上这个url,没有可写权限的可以不管)
	tid, _ := strconv.ParseInt(this.GetString("tid", ""), 10, 64)
	roles := models.Enforcer.GetRolesForUser(this.userName)
	checkUrl := fmt.Sprintf("%s/%s_%s", RootNameSpace, MetaDataNameSpace, "authority")
	models.Logger.Info(checkUrl)
	flag1 := models.Enforcer.Enforce(this.userName, checkUrl, "POST")
	var flag2 bool
	o := orm.NewOrm()
	sql := fmt.Sprintf("SELECT %s FROM table_meta_data WHERE tid = ?", models.MakeIsOwnerSql(roles))
	_ = o.Raw(sql, tid).QueryRow(&flag2)
	res, _ := json.Marshal(map[string]bool{
		"readOnly": !(flag1 && flag2),
	})
	this.Ctx.WriteString(string(res))
}

// @Title GetEdit
// @Summary 查看是否可以编辑
// @Description 查看是否可以编辑
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getEdit [get]
func (this *MetadataTableController) GetEdit() {

	// 通过casbin,判断/jones/jdbc_authority (可写权限的角色要配置上这个url,没有可写权限的可以不管)
	method := "POST"
	checkUrl := fmt.Sprintf("%s/%s_%s", RootNameSpace, MetaDataNameSpace, "authority")
	models.Logger.Info(checkUrl)
	//fmt.Printf("check url: %s\n",checkUrl)
	flag := models.Enforcer.Enforce(this.userName, checkUrl, method)
	res, _ := json.Marshal(map[string]bool{
		"readOnly": !flag,
	})
	this.Ctx.WriteString(string(res))
}

// @Title GetTablePrivilegeByTid
// @Summary 查看表有哪些权限
// @Description 查看表有哪些权限
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getTablePrivilegeByTid/:tid [get]
func (this *MetadataTableController) GetTablePrivilegeByTid() {
	param := this.Ctx.Input.Params()[":tid"]
	tid, _ := strconv.ParseInt(param, 10, 64)
	roles := models.GetTablePrivilegeByTid(tid)
	res, _ := json.Marshal(roles)
	this.Ctx.WriteString(string(res))
}
