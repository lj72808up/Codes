package controllers

import (
	"AStar/models"
	"AStar/models/airflowDag"
	"encoding/json"
	"fmt"
	"strconv"
)

type AirflowController struct {
	BaseController
}

var AirflowNameSpace = "airflow"

// @Title GetAllAirflow
// @Summary 获取所有图表信息
// @Description 获取所有图表信息
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /airflowall [get]
func (this *AirflowController) GetAllAirflow() {
	checkUrl := fmt.Sprintf("%s/%s_%s", RootNameSpace,AirflowNameSpace, "authority")
	models.Logger.Info(checkUrl)
	flag := models.Enforcer.Enforce(this.userName, checkUrl, "POST")

	resInfo := models.GetAllDagInfo(flag, this.userName)
	jsonData, _ := json.Marshal(resInfo)
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}

// @Title SwitchDag
// @Summary 改变Dag开关标志位
// @Description 改变Dag开关标志位
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /switchdag [post]
func (this *AirflowController) SwitchDag() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	dagName := attr["name"]
	isOpen := attr["isopen"]
	suc_flag := models.ChangeDagSwitchAPI(dagName, isOpen)

	//API更新成功之后更新本地DB
	if suc_flag == true {
		suc_flag = models.ChangeDagSwitchDB(dagName, isOpen)
	}

	param, _ := json.Marshal(map[string]bool{
		"flag": suc_flag,
	})

	this.Ctx.WriteString(string(param))
}

// @Title DeleteDag
// @Summary 删除Dag
// @Description 删除Dag
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /deletedag [post]
func (this *AirflowController) DeleteDag() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	dagName := attr["name"]

	//DeleteRemoteDag File
	models.DeleteDagFile(dagName)

	//更新本地DB为删除中
	flag := models.UpdateDagStatus(dagName, int(airflowDag.REMOVING))

	/*	//调用API删除Dag
		flag := models.DeleteDagApi(dagName)

		//以上步骤成功，则更新数据库状态为已删除
		if(flag == true){
			flag = models.UpdateDagStatus(dagName,int(airflowDag.REMOVED))
		}*/

	param, _ := json.Marshal(map[string]bool{
		"flag": flag,
	})
	this.Ctx.WriteString(string(param))
}



// @Title MarkDagState
// @Summary 改变Dag指定日期状态
// @Description Dag指定日期 清除标记、置为成功、置为失败
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /markdagstate [post]
func (this *AirflowController) MarkDagState() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	dagName := attr["name"]
	runTime := attr["execution_date"]
	newStatus := attr["new_status"]

	flag := models.MarkDagApi(dagName,runTime,newStatus)

	if(flag == true){
		flag = models.UpdateDagRunStatus(dagName,runTime,newStatus)
	}

	param, _ := json.Marshal(map[string]string{
		"flag" : strconv.FormatBool(flag),
		"newStatus":newStatus,
	})
	this.Ctx.WriteString(string(param))
}



// @Title MarkTaskState
// @Summary 改变Task指定日期状态
// @Description Task指定日期 清除标记、置为成功、置为失败
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /markTaskState [post]
func (this *AirflowController) MarkTaskState() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	dagName := attr["dag_name"]
	taskId := attr["task_id"]
	runTime := attr["execution_date"]
	newStatus := attr["new_status"]

	flag := models.MarkTaskApi(dagName,taskId,runTime,newStatus)
	param, _ := json.Marshal(map[string]string{
		"flag" : strconv.FormatBool(flag),
		"newStatus":newStatus,
	})
	this.Ctx.WriteString(string(param))
}


// @Title GetAllTemplateName
// @Summary 查询所有子节点
// @Description 查询所有子节点
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/getAllTemplate/:tasktype [get]
func (this *AirflowController) GetAllTemplate() {
	tasktype := this.Ctx.Input.Params()[":tasktype"]
	template := new(models.SqlTemplate) // 查找失败的话其Taskid会为0
	roles := models.Enforcer.GetRolesForUser(this.userName)
	children := template.GetAllTemplate(roles,tasktype,this.userName)
	marshal, _ := json.Marshal(children)
	this.Ctx.WriteString(string(marshal))
}

// @Title GetSchedulerJob
// @Summary 查询一个模板
// @Description 查询一个模板
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /template/getById/:templateId [get]
func (this *AirflowController) GetJdbcTemplate() {
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

// @Title 获取自身角色可以看到的连接
// @Summary 获取自身角色可以看到的连接
// @Description 获取自身角色可以看到的连接
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /GetConnByPrivilege/:connType [get]
func (this *AirflowController) GetConnByPrivilege() {
	connType := this.Ctx.Input.Params()[":connType"]
	roles := models.Enforcer.GetRolesForUser(this.userName)
	conns := models.GetConnByPrivilege(connType, roles)
	res, _ := json.Marshal(conns)
	this.Ctx.WriteString(string(res))
}
