package test

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/jmoiron/sqlx"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	connType := ""
	fmt.Println(models.GetAllConn(connType))
}


func TestCK(t *testing.T) {
	url := "tcp://ck1111111.adtech.sogou:9000?user=default&password=6lYaUiFi&debug=true"
	dbConn, err := sqlx.Open("clickhouse", url) // fmt.Sprintf("%s&database=%s", url, dbName))
	if err != nil {
		panic(err.Error())
	}
	var res int
	sql := "select 1"
	err = dbConn.Select(&res, sql)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
}

func Test1(t *testing.T) {
	fmt.Println(c()) // 输出2
}

func c() (i int) {
	defer func() { i++ }()
	return 1
}

func TestDate(t *testing.T) {
	str := `Tue, 30 Jun 2020 01:45:00 GMT`
	format := time.RFC1123
	time1, _ := time.Parse(format, str)
	fmt.Println(utils.TimeFormat(time1))

	str = "2020-06-29T01:45:00+08:00"
	format = time.RFC3339
	time2, _ := time.Parse(format, str)
	fmt.Println(utils.TimeFormat(time2))
}

func TestReverseArray(t *testing.T) {
	arr := []int{
		1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7,
	}

	batch := 3

	i := 0
	j := len(arr) - batch
	cnt := 1 // 轮数

	for i < j {
		for increase := 0; increase < batch; increase++ {
			tmp := arr[i]
			arr[i] = arr[j]
			arr[j] = tmp
			i = i + 1
			j = j + 1
		}
		cnt = cnt + 1
		j = len(arr) - batch*cnt
	}

	fmt.Println(arr)
}

func TestPostDagDelete(t *testing.T) {
	url := "http://api.airflow.adtech.sogou/api/experimental/dag_delete"
	param, _ := json.Marshal(map[string]string{
		"dag_id": "aaaa",
	})

	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body(param)
	bytes, err := req.Bytes()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(bytes))
}


func TestRunDate(t *testing.T) {
	start := "20200101"
	end := "20200202"
	println(strings.Join(utils.GetDateRange(start, end),","))
}

func TestIntervalRunn (t *testing.T) {
	maxRunner := 4
	jobList := []int{1,2,3,4,5,6,7,8,9}

	// 用信号同步量, 每次只能运行3个
	var goSync sync.WaitGroup
	for i := 0; i < len(jobList); {
		if (i + maxRunner) > len(jobList) {
			for j := i; j < len(jobList); j++ {
				// 执行函数  print(jobList[j])
				goSync.Add(1)
				go dummy_run(jobList[j], &goSync)
			}
		} else {
			for j := i; j < i+maxRunner; j++ {
				// 执行函数  print(jobList[j])
				goSync.Add(1)
				go dummy_run(jobList[j],  &goSync)
			}
		}
		i = i + maxRunner
		goSync.Wait()
		println()
	}
}

func dummy_run(val int, goSync *sync.WaitGroup) {
	print("开始 ")
	print (val)
	print(",")
	defer goSync.Done()
	time.Sleep(15 * time.Second)
}