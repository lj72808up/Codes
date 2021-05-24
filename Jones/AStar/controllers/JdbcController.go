package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type JdbcController struct {
	BaseController
}

var JdbcNameSpace = "jdbc"

// @Title PutJdbcTemplate
// @Summary 上传sql模板
// @Description 上传sql模板
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/put [put]
func (this *JdbcController) PutJdbcTemplate() {
	//  { creator:"", title:"", template:"", param:"" }
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var template models.SqlTemplate
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &template)
	templateId,err := template.NewTemplate()
	if err!="" {
		this.Ctx.WriteString(fmt.Sprintf("fail:%s",err))
		return
	}

	if template.Pid == 0 {
		roles := models.Enforcer.GetRolesForUser(this.userName)
		for _, role := range roles {
			mapping := models.GroupWorksheetMapping{
				GroupName:   role,
				WorksheetId: fmt.Sprint(templateId),
				Valid:       true,
				CanModify:   true,
			}
			mapping.AddOneMapping()
		}
	}

	this.Ctx.WriteString("finish:" + fmt.Sprint(templateId))
}

// @Title GetSchedulerJob
// @Summary 查询一个模板
// @Description 查询一个模板
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/getById/:templateId [get]
func (this *JdbcController) GetJdbcTemplate() {
	param := this.Ctx.Input.Params()[":templateId"]
	jobId, err := strconv.ParseInt(param, 10, 64) // 转换成64位10进制数字
	if err != nil {
		this.Ctx.WriteString("参数转换异常")
		return
	}
	template := new(models.SqlTemplate) // 查找失败的话其Taskid会为0
	template.GetSqlTemplate(jobId)
	marshal, err := json.Marshal(template)
	this.Ctx.WriteString(string(marshal))
}

// @Title GetCanModify
// @Summary 查看用户所属角色是否有保存修改的权限
// @Description 查看用户所属角色是否有保存修改的权限
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/canModify/:templateId [get]
func (this *JdbcController) GetCanModify() {
	param := this.Ctx.Input.Params()[":templateId"]
	roles := models.Enforcer.GetRolesForUser(this.userName)
	workSheetId, err := strconv.ParseInt(param, 10, 64) // 转换成64位10进制数字
	if err != nil {
		this.Ctx.WriteString("参数转换异常")
		return
	}
	template := new(models.SqlTemplate) // 查找失败的话其Taskid会为0
	canModify := template.CanModifyById(workSheetId, roles)
	this.Ctx.WriteString(fmt.Sprint(canModify))
}

// @Title GetSchedulerJob
// @Summary 更新一个模板
// @Description 更新一个模板
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/updateJdbcTemplate [post]
func (this *JdbcController) UpdateJdbcTemplate() {

	//templateId, err := strconv.ParseInt(param, 10, 64) // 转换成64位10进制数字
	//fmt.Println(templateId)
	//if err != nil {
	//	this.Ctx.WriteString("参数转换异常")
	//	return
	//}

	var bodyparams map[string]string
	json.Unmarshal(this.Ctx.Input.RequestBody, &bodyparams)
	title := bodyparams["title"]
	template := bodyparams["template"]
	desc := bodyparams["description"]
	templateId, err := strconv.ParseInt(bodyparams["templateId"], 10, 64)
	if err != nil {
		this.Ctx.WriteString("参数转换异常")
		return
	}
	sqlTemplate := models.SqlTemplate{TemplateId: templateId}
	sqlTemplate.UpdateTemplate(title, template, desc,this.userName)

	models.Logger.Info(title, template, templateId)
	this.Ctx.WriteString("save success")
}

// @Title CreateQueryView
// @Summary 创建Query可视化
// @Description 创建Query可视化
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/createQueryView [post]
func (this *JdbcController) CreateQueryView() {
	var bodyparams map[string]string
	json.Unmarshal(this.Ctx.Input.RequestBody, &bodyparams)
	templateId, err := strconv.ParseInt(bodyparams["templateId"], 10, 64)
	if err != nil {
		this.Ctx.WriteString("参数转换异常")
		return
	}
	sqlTemplate := models.SqlTemplate{TemplateId: templateId}
	msg,redash_query_id := sqlTemplate.CreateQueryViewModel(this.userName)
	if len(msg) ==0{
		msg = "可视化创建成功!"
	}
	res, _ := json.Marshal(map[string]string{
		"msg": msg,
		"redash_query_id":strconv.Itoa(redash_query_id),
	})

	this.Ctx.WriteString(string(res))

}

// @Title SQLCreateQueryView
// @Summary SQL查询中心通过自定义SQL创建Query可视化
// @Description SQL查询中心通过自定义SQL创建Query可视化
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /sql/createQueryView [post]
func (this *JdbcController) SQLCreateQueryView() {
	sqlTemplate := models.SqlTemplate{TemplateId: 0}
	sqlTemplate.Creator = this.userName
	sqlTemplate.Title = this.GetString("name")
	sqlTemplate.Template = this.GetString("sql")
	sqlTemplate.Dsid,_ = strconv.ParseInt(this.GetString("dsId"),10,64)

	msg,redash_query_id := sqlTemplate.SQLCreateQueryViewModel()
	if len(msg) ==0{
		msg = "可视化创建成功!"
	}
	res, _ := json.Marshal(map[string]string{
		"msg": msg,
		"redash_query_id":strconv.Itoa(redash_query_id),
	})

	this.Ctx.WriteString(string(res))

}


// @Title LvpiCreateQueryView
// @Summary 一键绿皮功能
// @Description 一键绿皮功能
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /sql/createLvpiView [post]
func (this *JdbcController) LvpiCreateQueryView() {
	tablename 		:= this.GetString("tablename")
	encode 			:= this.GetString("encode")
	date_pianyi,_ 	:= this.GetBool("date_pianyi",false)
	tableid 		:= this.GetString("tableid")
	title 			:= this.GetString("name")
	dsid,_ 			:= strconv.ParseInt(this.GetString("dsId"),10,64)

	sql := fmt.Sprintf("select * from %s where date='{{基准日期}}' or date='{{对照日期}}'; ",tablename)

	o := orm.NewOrm()
	err := o.Raw("SELECT * from field_meta_data where LOWER(field_name)='date' and table_id=?", tableid).QueryRow(this)
	if err == orm.ErrNoRows {
		res, _ := json.Marshal(map[string]string{
			"msg": "所选表没有date字段",
			"redash_query_id":strconv.Itoa(0),
		})
		this.Ctx.WriteString(string(res))
		return
	}

	sqlTemplate := models.SqlTemplate{TemplateId: 0}
	sqlTemplate.Creator 	= this.userName
	sqlTemplate.Title 		= title
	sqlTemplate.Template 	= sql
	sqlTemplate.Dsid  		= dsid

	msg,redash_query_id := sqlTemplate.LvpiCreateQueryViewModel(encode,date_pianyi)
	if len(msg) ==0{
		msg = "可视化创建成功!"
	}
	res, _ := json.Marshal(map[string]string{
		"msg": msg,
		"redash_query_id":strconv.Itoa(redash_query_id),
	})

	this.Ctx.WriteString(string(res))

}

// @Title PutJdbcTemplate
// @Summary 运行模板sql语句
// @Description 上传sql模板
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/run [post]
func (this *JdbcController) RunJdbcTemplate() {
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	dsid, err := strconv.ParseInt(attr["dsid"], 10, 64)
	if err != nil {
		panic(err.Error())
	}

	sql := attr["sql"]
	template := new(models.SqlTemplate)

	template.Title = attr["queryName"]
	template.Dsid = dsid
	qwFileContent := attr["queryWords"]

	var param []map[string]string
	_ = json.Unmarshal([]byte(attr["param"]), &param)

	partition := ""
	if qwFileContent != "" {
		qwFileContent = strings.Replace(qwFileContent, "[|]|\"", "", -1)
		partition = LoadQWPartitionAfterTranscoding(qwFileContent, this.userName)
		fmt.Printf("query word partition:%s \n ", partition)
		for i, _ := range param {
			if param[i]["type"] == "qw_partition" {
				param[i]["value"] = partition
			}
		}
		fmt.Println("after repair param:", param)

	}

	defer func() {
		if r := recover(); r != nil {
			models.Logger.Error("%v", r)
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()
	//根据前台传来的sql和param进行模板替换, below
	runSql := template.GenerateSql(sql, param, "dag")
	queryId := template.PublishMq(sql, dsid, this.userName, runSql, attr["param"], partition, attr["queryName"])

	templateId, _ := strconv.ParseInt(attr["templateId"], 10, 64)
	template.TemplateId = templateId
	template.AddTemplateHeat(queryId, this.userName)
	this.Ctx.WriteString(utils.VueResponse{
		Res:  true,
		Info: queryId,
	}.String())
}

// @Title GetAllTemplateName
// @Summary 查询所有子节点
// @Description 查询所有子节点
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/getChildNodes/:templateId [get]
func (this *JdbcController) GetChildNodes() {
	template := new(models.SqlTemplate) // 查找失败的话其Taskid会为0
	param := this.Ctx.Input.Params()[":templateId"]
	templateId, err := strconv.ParseInt(param, 10, 64) // 转换成64位10进制数字
	if err != nil {
		this.Ctx.WriteString("参数转换异常")
		return
	}
	children := template.GetChildren(templateId)
	marshal, _ := json.Marshal(children)
	this.Ctx.WriteString(string(marshal))
}

// @Title GetAllTemplateName
// @Summary 查询所有子节点
// @Description 查询所有子节点
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/getAllTemplate [get]
func (this *JdbcController) GetAllTemplate() {
	template := new(models.SqlTemplate) // 查找失败的话其Taskid会为0
	roles := models.Enforcer.GetRolesForUser(this.userName)
	children := template.GetAllTemplate(roles,"all",this.userName)
	marshal, _ := json.Marshal(children)
	this.Ctx.WriteString(string(marshal))
}

// @Title GetAllDataSources
// @Summary 查询所有数据源
// @Description 查询所有数据源
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/getAllDataSources [get]
func (this *JdbcController) GetAllDataSources() {
	template := new(models.SqlTemplate)
	connType := this.GetString("connType") // 连接类型
	source := template.GetAllDataSource(connType)
	marshal, _ := json.Marshal(source)
	this.Ctx.WriteString(string(marshal))
}

// @Title GetEdit
// @Summary 查看是否可以编辑模板
// @Description 查看是否可以编辑模板, 如果是某个用户只能编辑, 则在LimitedAccessElement表中加入username和"/template/getEdit"
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/getEdit [get]
func (this *JdbcController) GetEdit() {
	//old: 判断是否是jdbcCheckGroup
	//rolesForUser := models.Enforcer.GetRolesForUser(this.userName)
	//flag := false
	//for _, role := range rolesForUser {
	//	if role == "jdbcCheckGroup" {
	//		flag = true
	//	}
	//}

	// 通过casbin,判断/jones/jdbc_authority (可写权限的角色要配置上这个url,没有可写权限的可以不管)
	method := "POST"
	checkUrl := fmt.Sprintf("%s/%s_%s", RootNameSpace, JdbcNameSpace, "authority")
	models.Logger.Info(checkUrl)
	//fmt.Printf("check url: %s\n",checkUrl)
	flag := models.Enforcer.Enforce(this.userName, checkUrl, method)
	res, _ := json.Marshal(map[string]bool{
		"readOnly": !flag,
	})
	this.Ctx.WriteString(string(res))
}

// @Title GetViewAccess
// @Summary 判断用户所在组是否有可视化权限
// @Description 判断用户所在组是否有可视化权限
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/getviewaccess [get]
func (this *JdbcController) GetViewAccess() {
	template := new(models.SqlTemplate) // 查找失败的话其Taskid会为0
	roles := models.Enforcer.GetRolesForUser(this.userName)
	viewaccess := template.GetUserViewAccess(roles)
	res, _ := json.Marshal(map[string]bool{
		"viewaccess": viewaccess,
	})
	this.Ctx.WriteString(string(res))

}


// @Title GetJumpAccess
// @Summary 判断用户是否有在数据工单页面跳转可视化的权限
// @Description 判断用户是否有在数据工单页面跳转可视化的权限
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/getjumpaccess [post]
func (this *JdbcController) GetJumpAccess() {
	var bodyparams map[string]string
	json.Unmarshal(this.Ctx.Input.RequestBody, &bodyparams)
	templateId, err := strconv.ParseInt(bodyparams["templateId"], 10, 64)
	if err != nil {
		this.Ctx.WriteString("参数转换异常")
		return
	}
	sqlTemplate := models.SqlTemplate{TemplateId: templateId}
	jumpaccess := sqlTemplate.GetUserJumpAccess(this.userName)

	res, _ := json.Marshal(map[string]bool{
		"jumpaccess": jumpaccess,
	})
	this.Ctx.WriteString(string(res))

}

// @Title DeleteTemplateById
// @Summary 删除工单
// @Description 删除工单
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/deleteTemplateById [post]
func (this *JdbcController) DeleteTemplateById() {
	var bodyparams map[string]string
	json.Unmarshal(this.Ctx.Input.RequestBody, &bodyparams)
	templateId,_ := strconv.ParseInt(bodyparams["templateId"], 10, 64)

	sqlTemplate := models.SqlTemplate{}
	err := sqlTemplate.DeleteById(templateId)

	if err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	}
}

// @Title ChangeTemplateAffiliation
// @Summary 更改工单所属目录
// @Description 更改工单所属目录
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/ChangeTemplateAffiliation [post]
func (this *JdbcController) ChangeTemplateAffiliation() {
	var bodyparams map[string]string
	json.Unmarshal(this.Ctx.Input.RequestBody, &bodyparams)
	templateId,_ := strconv.ParseInt(bodyparams["templateId"], 10, 64)
	folderId,_ := strconv.ParseInt(bodyparams["folderId"], 10, 64)

	sqlTemplate := models.SqlTemplate{}
	err := sqlTemplate.ChangeAffiliation(templateId, folderId)

	if err == nil {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "",
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	}
}
