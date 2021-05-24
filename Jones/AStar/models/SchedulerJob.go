package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"time"
)

type SchedulerJob struct {
	TaskId       int64     `orm:"pk;auto"` // tag中加入json:"-", 会在转换json时忽略该字段
	Owner        string    `orm:"index"`
	CreateTime   time.Time `orm:"type(timestamp)"`
	CronExpress  string    `orm:"size(50)"`
	LastExecTime time.Time `orm:"type(timestamp)"`
}

func (this SchedulerJob) NewJob() int64 {
	taskId, err := orm.NewOrm().Insert(&this)
	if err != nil {
		Logger.Error(err.Error())
	}
	return taskId
}

func (this *SchedulerJob) GetSchedulerJob(jobId int64) {
	o := orm.NewOrm()
	//user := User{Id: 1}
	err := o.Raw("SELECT * FROM scheduler_job WHERE task_id=?", jobId).QueryRow(this)
	if err == orm.ErrNoRows {
		Logger.Error("查询不到")
	} else if err == orm.ErrMissPK {
		Logger.Error("找不到主键")
	} else {
		Logger.Info("owner:",this.Owner, "createTime",this.CreateTime)
	}
}

func UpdateLastExecTime(jobId int64) {
	o := orm.NewOrm()
	job := SchedulerJob{TaskId: jobId}

	t := time.Now()

	if readErr := o.Read(&job, "TaskId"); readErr == nil {
		job.LastExecTime = t
		_, err := o.Update(&job)
		if err != nil {
			Logger.Error(err.Error())
		}
	} else {
		Logger.Error(readErr.Error())
	}
}

func ExecScheduleTask(jobId int64, jobName string, cronExpress string, taskFun toolbox.TaskFunc) {
	task := toolbox.NewTask(jobName, cronExpress, taskFun)
	toolbox.AddTask(jobName, task)
}

func StopScheduleTask(jobName string) {
	toolbox.DeleteTask(jobName)
}
