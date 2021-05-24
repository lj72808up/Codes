package models

import (
	"AStar/utils"
	"github.com/astaxie/beego/toolbox"
	"time"
)

func airflowTask() error {
	t := utils.TimeFormat(time.Now())
	Logger.Info("check airflow task (creating or deleting) at [%s]", t)
	dag := DagAirFlow{}
	dag.UpdateStatus()
	return nil
}

func AddGlobalTask() {
	airflowTask := toolbox.NewTask("airflowTask",
		"0 0/1 * * * *",
		airflowTask)
	toolbox.AddTask("airflowTask", airflowTask)
	toolbox.StartTask()
}
