package test

import (
	"AStar/models"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2"
	"strings"
	"testing"
	"time"
	"unicode"
)

type IntervalVar struct {
	start string
	end   string
}

func TestExceptionSql(t *testing.T) {
	sql := "SELECT `哈哈啊` ,{{data.Start}}, {{包Id}}, {{age}}"
	param := make(map[string]interface{})
	param["data.Start"] = 20201103
	param["包Id"] = "zzzz"
	param["age"] = 13

	replaceKeys := findChineseAndDotaKeys(param)

	for k,v := range replaceKeys {
		sql = strings.Replace(sql, k, v.NewKey, -1)
		param[v.NewKey] = v.NewValue
		delete(param,k)
	}

	tpl, _ := pongo2.FromString(sql)
	fmt.Println("sql:"+sql)
	fmt.Printf("param:%v\n",param)
	out, _ := tpl.Execute(param)
	fmt.Println(out)
}

type TransParam struct {
	NewKey   string
	NewValue interface{}
}

func findChineseAndDotaKeys(m map[string]interface{}) map[string]TransParam {
	replaceKeys := make(map[string]TransParam)
	cnt := 0
	for k, _ := range m {
		if isChinese(k) {
			// 中文的改变变量名
			newKey := fmt.Sprintf("chineseVar%d", cnt)
			cnt = cnt + 1
			newValue := m[k]
			replaceKeys[k] = TransParam{
				NewKey:   newKey,
				NewValue: newValue,
			}
		}

		if strings.Contains(k, ".") {
			newKey := strings.Replace(k, ".", "", -1)
			newValue :=  m[k]
			replaceKeys[k] = TransParam{
				NewKey:   newKey,
				NewValue: newValue,
			}
		}
	}
	return replaceKeys
}

func isChinese(str string) bool {
	flag := false
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			flag = true
			break
		}
	}
	return flag
}

func TestReplace(t *testing.T) {
	res := strings.Replace("hahaha  haha.hehe","haha.hehe", "wc", -1)
	fmt.Println(res)
}

func TestPongoSQL(t *testing.T) {
	tmpl := `SELECT * FROM test1 WHERE date>='{{LastWeek.start}}' and date<='{{LastWeek.end}}' LIMIT {{cnt}}`
	tpl, _ := pongo2.FromString(tmpl)
	param := make(map[string]interface{})
	param["cnt"] = 13
	param["LastWeek"] = IntervalVar{
		start: "20190101",
		end:   "20200101",
	}
	out, _ := tpl.Execute(param)
	fmt.Println(out)
}
func TestPongo(t *testing.T) {
	// Compile the template first (i. e. creating the AST)
	tpl, err := pongo2.FromFile("airflow.tmpl")
	if err != nil {
		panic(err)
	}
	param := make(map[string]interface{})
	param["start_year"] = "2019"
	param["start_month"] = "11"
	param["start_day"] = "1"
	param["dagType"] = "day"
	param["delay"] = 1
	users := []string{
		"lj", "lj02",
	}
	emailUsers := fmt.Sprintf("\"%s\"", strings.Join(users, "\",\""))
	param["emailUsers"] = emailUsers
	param["cron_expression"] = "55 * * * *"
	param["dag_name"] = "Test_Dag"

	// task1
	taskParam := make(map[string]string)
	taskParam["cmd"] = "cd yuedu/ods/andriod; sh run.sh"
	env := map[string]string{
		"aaa": "1111",
		"bbb": "2222",
	}
	envRes, _ := json.Marshal(env)
	taskParam["env_json"] = string(envRes)
	taskParam["task_name"] = "task1"

	//task2
	taskParam2 := make(map[string]string)
	taskParam2["cmd"] = "cd /; sh run2.sh"
	env2 := map[string]string{
		"aaa2": "1111",
		"bbb2": "2222",
	}
	envRes2, _ := json.Marshal(env2)
	taskParam2["env_json"] = string(envRes2)
	taskParam2["task_name"] = "task2"

	param["tasks"] = []map[string]string{
		taskParam, taskParam2,
	}

	param["edges"] = []models.DagEdge{
		{
			SrcId: 111,
			DstId: 222,
		},
	}
	// Now you can render the template with the given
	// pongo2.Context how often you want to.
	out, err := tpl.Execute(param)
	if err != nil {
		panic(err)
	}
	fmt.Println(out) // Output: Hello Florian!

}

type AA struct {
	HH []BB
	XX []map[string]string
}
type BB struct {
	Name string `json:"name_a"`
}

func TestJson(t *testing.T) {
	v := AA{
		HH: []BB{
			BB{Name: "zhangsan",},
			BB{Name: "lisi",},
		},
		XX: []map[string]string{
			{
				"k1": "v1",
				"k2": "v2",
			},
			{
				"wolegequ": "wolegequ",
			},
		},
	}

	res, _ := json.Marshal(v)
	println(string(res))

	var newA AA
	_ = json.Unmarshal([]byte(`{"HH":[{"name_a":"zhangsan","Sex":"Man"},{"Name":"lisi"}],"XX":[{"k1":"v1","k2":"v2"},{"wolegequ":"wolegequ"}]}`), &newA)
	fmt.Println(newA)
}

func TestPrseDate(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2020-07-06")
	fmt.Println(date)

	layout := "2006-01-02"
	str := "2014-03-17"
	t2, _ := time.Parse(layout, str)
	fmt.Println(t2)
}

func TestDagJson(t *testing.T) {
	jsonStr := `{
    "edges": [
        {
            "src_node_id": 100,
            "src_output_idx": 0,
            "dst_node_id": 101,
            "dst_input_idx": 0,
            "id": 10
        }
    ],
    "nodes": [
        {
            "pos_x": 200.4642857142857,
            "pos_y": 98.71428571428571,
            "name": "task1",
            "rightClickEvent": [
                {
                    "label": "配置DAG调度",
                    "eventName": "editDAG"
                },
                {
                    "label": "配置任务",
                    "eventName": "editTask"
                }
            ],
            "id": 100,
            "in_ports": [
                0
            ],
            "out_ports": [
                0
            ],
            "image": "image1",
            "command": "ls /",
            "envs": [
                {
                    "env_key": "name",
                    "env_val": "zhangsan"
                }
            ]
        },
        {
            "pos_x": 110.46428571428571,
            "pos_y": 244.42857142857136,
            "name": "task2",
            "rightClickEvent": [
                {
                    "label": "配置DAG调度",
                    "eventName": "editDAG"
                },
                {
                    "label": "配置任务",
                    "eventName": "editTask"
                }
            ],
            "id": 101,
            "in_ports": [
                0
            ],
            "out_ports": [
                0
            ],
            "image": "image2",
            "command": "ls /boot",
            "envs": [
                {
                    "env_key": "sex",
                    "env_val": "male"
                }
            ]
        }
    ],
    "dagConfig": {
        "name": "test1",
        "startDate": "2020-07-03",
        "mailList": "liujie02",
        "cronExpression": "00 15 * * * ",
        "dagType": "dayLevel",
        "delay": 1
    }
}
`
	var view models.DagView
	_ = json.Unmarshal([]byte(jsonStr), &view)
	fmt.Println(view)
}
