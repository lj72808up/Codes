package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/saintfish/chardet"
	"os"
	"strconv"
	"strings"
	"time"
)

type FieldsMetaController struct {
	BaseController
}

// @Title PutFieldsMetaData
// @Summary 新增一个版本的字段元数据
// @Description 新增一个版本的字段元数据
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/newVersion/put [post]
func (this *FieldsMetaController) PutFieldsMetaData() {
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)

	tableId, _ := strconv.ParseInt(params["TableId"], 10, 64)
	totallyNew, _ := strconv.ParseBool(params["TotallyNew"]) // string 转bool

	var fields []map[string]string
	_ = json.Unmarshal([]byte(params["Fields"]), &fields)

	newVersionData := models.NewVersionData{
		TableId: tableId,
		Version: params["Version"],
		Fields:  fields,
	}
	models.AddNewVersionFieldData(newVersionData, totallyNew, tableId)
	this.Ctx.WriteString("success")
}

type VersionTableName struct {
	TableName string
	Versions  []models.VersionMetaData
}

// @Title GetFieldAllVersion
// @Summary 查看字段元数据的所有版本
// @Description 查看字段元数据的所有版本
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/version/get/:tableId [get]
func (this *FieldsMetaController) GetFieldAllVersion() {
	param := this.Ctx.Input.Params()[":tableId"]
	tableId, _ := strconv.ParseInt(param, 10, 64)
	versions := models.GetFieldsAllVersion(tableId)
	tableName := models.GetTableNameById(tableId)
	view := VersionTableName{
		TableName: tableName,
		Versions:  versions,
	}
	res, _ := json.Marshal(view)
	this.Ctx.WriteString(string(res))
}

// @Title GetDefaultNewVersion
// @Summary 增加全新版本时,获取一个默认版本
// @Description 增加全新版本时,获取一个默认版本
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/version/getNewVersion/:tableId [get]
func (this *FieldsMetaController) GetDefaultNewVersion() {
	param := this.Ctx.Input.Params()[":tableId"]
	tableId, _ := strconv.ParseInt(param, 10, 64)
	newVersion := models.GetDefaultNewVersion(tableId)
	this.Ctx.WriteString(newVersion)
}

// @Title GetFieldsMetaDataByVersion
// @Summary 查看表某一个版本的字段元数据
// @Description 查看表某一个版本的字段元数据
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/getByVersion [post]
func (this *FieldsMetaController) GetFieldsMetaDataByVersion() {
	// { “tableId”:””, “version”:”” }
	var param map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &param)
	tableId, _ := strconv.ParseInt(param["tableId"], 10, 64)
	vxeData := models.GetFieldsDataByVersion(tableId, param["version"])
	this.Ctx.WriteString(vxeData.ToString())
}

// @Title GetPartitionFieldsByVersion
// @Summary 获取分区字段信息
// @Description 获取分区字段信息
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/getPartitionByVersion [post]
func (this *FieldsMetaController) GetPartitionFieldsByVersion() {
	// { “tableId”:””, “version”:”” }
	var param map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &param)
	tableId, _ := strconv.ParseInt(param["tableId"], 10, 64)
	partitionFields := models.GetPartitionFieldsByVersion(tableId, param["version"])
	res, _ := json.Marshal(partitionFields)
	this.Ctx.WriteString(string(res))
}

// @Title GetFieldsListByVersion
// @Summary 查看表某一个版本的字段元数据
// @Description 查看表某一个版本的字段元数据
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/getListByVersion [post]
func (this *FieldsMetaController) GetFieldsListByVersion() {
	// { “tableId”:””, “version”:”” }
	var param map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &param)
	tableId, _ := strconv.ParseInt(param["tableId"], 10, 64)
	fieldDataView := models.GetFieldsListByVersion(tableId, param["version"])
	res, _ := json.Marshal(fieldDataView)
	this.Ctx.WriteString(string(res))
}

// @Title UpdateFieldsByVersion
// @Summary 修改某个版本的某一个字段的元数据
// @Description 修改某个版本的某一个字段的元数据
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/updateByVersion [put]
func (this *FieldsMetaController) UpdateFieldsByVersion() {
	//  { “tableId”:,
	// 	    “fieldId”:””,
	//		“fieldName”:””,
	//		"fieldType":"",
	// 		"description":""   // 被修改的字段属性 ]
	//		"version":""
	// }
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var data models.FieldMetaData
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	tableId := data.TableId
	now := utils.TimeFormat(time.Now())
	data.ModificationTime = now
	models.UpdateFieldsByVersion(data, now, tableId)
	this.Ctx.WriteString("")
}

// @Title GetTableTypes
// @Summary 查询所有表类型
// @Description 查询所有表类型
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /tableTypes/get [get]
func (this *FieldsMetaController) GetTableTypes() {
	tTypes := models.GetAllTableTypes()
	var types []string
	for _, t := range tTypes {
		types = append(types, t.TypeName)
	}
	res, _ := json.Marshal(types)
	this.Ctx.WriteString(string(res))
}

// @Title GetDataBaseByType
// @Summary 查询该类型下的所有数据库
// @Description 查询该类型下的所有数据库
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getDataBaseByType/:tableType [get]
func (this *FieldsMetaController) GetDataBaseByType() {
	tableType := this.Ctx.Input.Params()[":tableType"]
	metaData := models.TableMetaData{}

	databases := metaData.GetDataBaseByType(tableType)
	res, _ := json.Marshal(databases)
	this.Ctx.WriteString(string(res))
}

// @Title GetTablesByDataBase
// @Summary 查询该数据库下的所有表
// @Description 查询该数据库下的所有表
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getTablesByDataBase [post]
func (this *FieldsMetaController) GetTablesByDataBase() {
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	dataBase := params["database"]
	tableType := params["tableType"]

	metaData := models.TableMetaData{}
	tables := metaData.GetTablesByDataBase(dataBase, tableType)
	res, _ := json.Marshal(tables)
	this.Ctx.WriteString(string(res))
}

// @Title PutFieldsRelation
// @Summary 配置某个字段的关联关系
// @Description 配置某个字段的关联关系
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/relation/update [post]
func (this *FieldsMetaController) PutFieldsRelation() {
	//  { “tableAId”:””
	//    “fieldAId”:””
	//    “tableBId”:””
	//    “fieldBId”:””
	//    “description”:””  // 关系描述
	//  }
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var param map[string]string
	json.Unmarshal(this.Ctx.Input.RequestBody, &param)

	tableAId, _ := strconv.ParseInt(param["tableAId"], 10, 64)
	tableBId, _ := strconv.ParseInt(param["tableBId"], 10, 64)
	fieldAId, _ := strconv.ParseInt(param["fieldAId"], 10, 64)
	fieldBId, _ := strconv.ParseInt(param["fieldBId"], 10, 64)
	tableAVersion := param["versionA"]
	tableBVersion := param["versionB"]

	models.AddFieldRelation(tableAId, fieldAId, tableAVersion, tableBId, fieldBId, tableBVersion, param["description"])
	this.Ctx.WriteString("success")
}

// @Title PutFieldsRelation
// @Summary 配置某个字段的关联关系
// @Description 配置某个字段的关联关系
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/relation/get/:fieldId [get]
func (this *FieldsMetaController) GetFieldsRelation() {
	//  { “tableAId”:””
	//    “fieldAId”:””
	//    “tableBId”:””
	//    “fieldBId”:””
	//    “description”:””  // 关系描述
	//  }
	param := this.Ctx.Input.Params()[":fieldId"]
	fieldId, _ := strconv.ParseInt(param, 10, 64)

	relation := models.GetFieldRelation(fieldId)
	res, _ := json.Marshal(relation)
	this.Ctx.WriteString(string(res))
}

// @Title DeleteField
// @Summary 删除某个字段
// @Description 删除某个字段
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/delete [post]
func (this *FieldsMetaController) DeleteField() {
	//  { “fieldId”:””}
	var param map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &param)
	fieldId, _ := strconv.ParseInt(param["fieldId"], 10, 64)
	models.DeleteFieldById(fieldId)
	this.Ctx.WriteString("success")
}

// @Title DeleteFieldsRelation
// @Summary 删除某个字段的关联关系
// @Description 删除某个字段的关联关系
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /fields/relation/delete [post]
func (this *FieldsMetaController) DeleteFieldsRelation() {
	//  { "rid":"" }
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var param map[string]int64
	json.Unmarshal(this.Ctx.Input.RequestBody, &param)

	models.DeleteFieldRelation(param["rid"])
	this.Ctx.WriteString("success")
}

// @Title 上传csv文件, 批量插入列之间的关系
// @Summary 上传csv文件, 批量插入列之间的关系
// @Description 上传csv文件, 批量插入列之间的关系
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /upload/relation/file [post]
func (this *FieldsMetaController) UploadFieldRelation() {
	// 接收文件
	uploadRoot := fmt.Sprintf("%s/%s/%s", models.ServerInfo.UploadRoot, this.userName,
		utils.DateFormatWithoutLine(time.Now()))
	utils.MkDir(uploadRoot)
	f, h, _ := this.GetFile("file")
	//fmt.Println(this.GetString("username", "")) // 表单中的其他参数
	tableAId, _ := this.GetInt64("tableAId", -1)
	fieldAId, _ := this.GetInt64("fieldAId", -1)
	tableAVersion := this.GetString("tableAVersion", "")

	toFileName := fmt.Sprintf("%s_%d", this.userName, time.Now().Unix())
	toFile := fmt.Sprintf("%s/%s", uploadRoot, toFileName)
	fmt.Printf("upload %s to %s\n", h.Filename, toFile)
	defer func() {
		f.Close()
		if r := recover(); r != nil {
			this.Ctx.WriteString(fmt.Sprintf("捕获到的错误：%s", r))
		} else {
			this.Ctx.WriteString("success")
		}
	}()

	bytes, _ := utils.ReadBytes(f, h.Size)
	detector := chardet.NewTextDetector()
	result, _ := detector.DetectBest(bytes)

	var resStr string
	if result.Charset != "UTF-8" {
		resStr = utils.ConvertToString(string(bytes), result.Charset, "UTF-8")
	} else {
		resStr = string(bytes)
	}

	fileData := strings.Replace(resStr, "\r", "", -1)
	bom := utils.HasBom(fileData)
	fmt.Println("fields中包含不可见字符:", bom)
	if bom {
		fileData = fileData[3:] // 文件头的bom占3字符(要判断是否包含开头的bom)
	}

	// 解析文件并插入
	models.InsertRelationByFile(fileData, tableAId, fieldAId, tableAVersion)

	os.RemoveAll(toFile)
	fmt.Println("已删除文件:", toFile)
	//this.Ctx.WriteString("success")
}

// @Title 上传csv文件, 批量插入列之间的关系
// @Summary 上传csv文件, 批量插入列之间的关系
// @Description 上传csv文件, 批量插入列之间的关系
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /graph/:tableId [get]
func (this *FieldsMetaController) Graph() {
	param := this.Ctx.Input.Params()[":tableId"]
	tableId, _ := strconv.ParseInt(param, 10, 64)
	graph := models.GetGraph(tableId)
	marshal, _ := json.Marshal(graph)
	this.Ctx.WriteString(string(marshal))
}

// @Title 查看图中2表的关系
// @Summary 查看图中2表的关系
// @Description 查看图中2表的关系
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /graph/relation/check [post]
func (this *FieldsMetaController) CheckTablePairRelation() {
	var param map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &param)
	sourceName := param["source"]
	targetName := param["target"]
	vxeData := models.GetGraphRelation(sourceName, targetName)
	fmt.Println(vxeData.ToString())
	this.Ctx.WriteString(vxeData.ToString())
}

// @Title 自动导入元数据
// @Summary 自动导入元数据
// @Description 自动导入元数据
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /import [post]
func (this *FieldsMetaController) GenerateFieldAuto() {
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var param map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &param)
	models.GenerateFieldByDefinition(param, this.userName)
	this.Ctx.WriteString("success")
}

// @Title 自动导入元数据(单表)
// @Summary 自动导入元数据(单表)
// @Description 自动导入元数据(单表)
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /import/single [post]
func (this *FieldsMetaController) GenerateFieldAutoSingle() {
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var param map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &param)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()
	// 去重
	if models.CheckDuplicateTable(param["tableType"], param["databaseName"], param["tableName"]) {
		response := utils.VueResponse{
			Res:  false,
			Info: "该表的元数据已经存在,请确保类型,库,表名称不同",
		}
		this.Ctx.WriteString(response.String())
	} else {
		flag := models.GenerateFieldByDefSingle(param, this.userName)
		response := utils.VueResponse{
			Res:  flag,
			Info: "",
		}
		this.Ctx.WriteString(response.String())
	}
}

// @Title 更新版本描述
// @Summary 更新版本描述
// @Description 更新版本描述
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /version/update [post]
func (this *FieldsMetaController) UpdateVersionDesc() {
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var newVersion models.VersionMetaData
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &newVersion)
	models.UpdateVersionData(&newVersion)
	this.Ctx.WriteString("success")
}

// @Title 更新字段
// @Summary 更新字段
// @Description 更新字段
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /field/update [post]
func (this *FieldsMetaController) UpdateFields() {
	//updateStr := `{"name":"f1","value":"111"}`
	//param := fmt.Sprintf(`[{"fid":"1","update":"%s"}]`, updateStr)
	var updates []models.FieldMetaData
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &updates)
	var resp utils.VueResponse
	defer func() {
		if err := recover(); err != nil {
			utils.PrintStackTrace()
			resp = utils.VueResponse{
				Res:  false,
				Info: fmt.Sprint(err),
			}
			this.Ctx.WriteString(resp.String())
		}
	}()
	models.ExecUpdateFields(&updates)
	resp = utils.VueResponse{
		Res:  true,
		Info: "更新成功",
	}
	this.Ctx.WriteString(resp.String())
}

// @Title 更新分区字段
// @Summary 更新分区字段
// @Description 更新分区字段
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /field/updatePartition [post]
func (this *FieldsMetaController) UpdatePartitionField() {
	var update models.FieldMetaData
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &update)
	var resp utils.VueResponse
	defer func() {
		if err := recover(); err != nil {
			utils.PrintStackTrace()
			resp = utils.VueResponse{
				Res:  false,
				Info: fmt.Sprint(err),
			}
			this.Ctx.WriteString(resp.String())
		}
	}()
	models.ExecUpdatePartitionField(&update)
	resp = utils.VueResponse{
		Res:  true,
		Info: "更新成功",
	}
	this.Ctx.WriteString(resp.String())
}
