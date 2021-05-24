package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/toolbox"
	"strconv"
	"time"
)

type SchedulerJobController struct {
	BaseController
}

type SchedulerJobMapping struct {
	JobName     string
	ToolBoxTask *toolbox.Task
}

// @Title PutSchedulerJob
// @Summary 添加一个定时任务
// @Description 添加一个定时任务
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /scheduler_job/put [put]
func (this *SchedulerJobController) PutSchedulerJob() {
	models.Logger.Info(string(this.Ctx.Input.RequestBody))
	var job models.SchedulerJob
	json.Unmarshal(this.Ctx.Input.RequestBody, &job)
	job.CreateTime = time.Now()
	job.LastExecTime = time.Now()
	jobId := job.NewJob()
	res := fmt.Sprintln(jobId)
	this.Ctx.WriteString("finish:" + res)
}

// @Title GetSchedulerJob
// @Summary 查询一个定时任务
// @Description 查询一个定时任务
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /scheduler_job/getById/:jobId [get]
func (this *SchedulerJobController) GetSchedulerJob() {
	param := this.Ctx.Input.Params()[":jobId"]
	jobId, err := strconv.ParseInt(param, 10, 64) // 转换成64位10进制数字
	if err != nil {
		this.Ctx.WriteString("参数转换异常")
		return
	}
	job := new(models.SchedulerJob) // 查找失败的话其Taskid会为0
	job.GetSchedulerJob(jobId)
	marshal, err := json.Marshal(job)
	this.Ctx.WriteString(string(marshal))
}

// @Title StartSchedulerJob
// @Summary 开启一个定时任务
// @Description 开启一个定时任务
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /scheduler_job/startById [post]
func (this *SchedulerJobController) StartSchedulerJob() {
	var data map[string]int64
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	jobId := data["jobId"]

	if err != nil {
		models.Logger.Error(err.Error())
		this.Ctx.WriteString("参数传输错误")
		return
	}
	job := new(models.SchedulerJob)
	job.GetSchedulerJob(jobId)

	taskFun := func() error {
		fmt.Println("========================================",time.Now().String(), "jobName: tk1")
		models.UpdateLastExecTime(jobId)
		return nil
	}
	jobName := utils.GetJobName(job.Owner, job.TaskId)

	log := fmt.Sprint("启动任务:", jobName, "定时时间:", job.CronExpress)
	models.Logger.Info(log)

	//task := toolbox.NewTask(jobName, job.SchedulerTime, taskFun)
	//toolbox.AddTask(jobName, task)
	models.ExecScheduleTask(jobId, jobName, job.CronExpress, taskFun)

	this.Ctx.WriteString(log)
	//toolbox.StartTask()
}

// 停止某个定时任务
// @Title StopSchedulerJob
// @Summary 关闭一个定时任务
// @Description 关闭一个定时任务
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /scheduler_job/stopJobById [post]
func (this *SchedulerJobController) StopSchedulerJob() {
	var data map[string]string
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	jobName := data["jobName"]

	if err != nil {
		models.Logger.Error(err.Error())
		this.Ctx.WriteString("参数传输错误")
		return
	}

	models.StopScheduleTask(jobName)
	this.Ctx.WriteString(fmt.Sprint(jobName,"已停止"))
	//toolbox.StartTask()
}