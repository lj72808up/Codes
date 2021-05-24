package utils

import (
	"fmt"
)

func GetJobName(owner string, taskId int64) string {
	jobId := fmt.Sprint(taskId)
	return owner + "-" + jobId
}