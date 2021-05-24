package test

import (
	"AStar/models"
	"AStar/utils"
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/axgle/mahonia"
	"github.com/saintfish/chardet"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestMd5(t *testing.T) {
	now := time.Now()
	data := md5.Sum([]byte(fmt.Sprint(now)))
	aaa := fmt.Sprintf("%s_%s_%x", "liujie02", utils.DateFormatWithoutLine(now), data)
	fmt.Println(aaa)
}

func TestAutoEncoding(t *testing.T) {
	b, _ := ioutil.ReadFile("/Users/liujie02/Desktop/baidu2.csv")
	str := string(b)

	println(utils.HasBom(str))

	detector := chardet.NewTextDetector()
	result, _ := detector.DetectBest(b)
	println(result.Charset)

	var resStr string
	//coding:="latin1"
	coding := "windows-1251"
	if result.Charset != "UTF-8" {
		resStr = utils.ConvertToString(string(b), coding, "UTF-8")
	} else {
		resStr = string(b)
	}

	println(resStr)

	convertString := mahonia.NewDecoder(coding).ConvertString(string(b))
	println(convertString)

}

// 4: 反射列表数据
func TestReflect(t *testing.T) {
	//mapAttr := utils.GetFieldTagName(models.FieldMetaData{},"json")
	//fmt.Println(mapAttr)
	// 创建表内容
	type AA struct {
		Name string `vxe:"name"`
		Age  int64  `vxe:"age"`
	}
	datas := []AA{
		{Name: "zhangsan", Age: 24}, {Name: "lisi", Age: 25},
	}

	var util utils.VxeData
	editableMap := map[string]bool{
		"name": true,
	}
	jsonTagMapping := util.MakeColumns(AA{}, editableMap)

	var dataList []map[string]interface{}
	for i, _ := range datas {
		dataList = append(dataList, util.MakeData(datas[i], &jsonTagMapping))
	}

	util.DataList = dataList
	println(util.ToString())
}

// 5. 解析csv
func TestCsv(t *testing.T) {
	content := `
zhangsan,23
"li,si",24
"wang,""wu""",25
`
	// r can be any io.Reader, including a file.
	csvReader := csv.NewReader(strings.NewReader(content))
	// Set comment character to '#'.
	//csvReader.Comment = '#'
	for {
		recoder, err := csvReader.Read()
		if err != nil {
			fmt.Println(err)
			// Will break on EOF.
			break
		}
		fmt.Printf("姓名:%s 年龄:%s\n", recoder[0], recoder[1])
	}
}

// 6. panic recover
func TestPanic(t *testing.T) {
	f := f()
	fmt.Println("返回", f)
}

func f() int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in f %s\n", r) // panic中赋予的事字符串
		}
	}()
	g(1)
	return 99
}

func g(i int) {
	panic(fmt.Sprintf("%v", i))

}

// 6. test mysqlConn
func TestOrmConn(t *testing.T) {
	aliasName := "default"
	url := "yuhancheng:yuhancheng@tcp(10.139.36.81:3306)/adtl_test"
	fmt.Printf("数据源: %s\n", url)

	// SetMaxIdleConns, SetMaxOpenConns: 最大空闲0,最大连接1,让orm执行一次查询后就断开连接
	connErr := orm.RegisterDataBase(aliasName, "mysql", url, 0, 1)
	if connErr != nil {
		panic(connErr)
	}
	o1 := orm.NewOrm()
	_ = o1.Using(aliasName)

	// 1. 执行查询
	sql := fmt.Sprintf("SELECT * FROM result_full limit 1")
	var res []orm.Params
	_, _ = o1.Raw(sql).Values(&res)

	a := 99
	fmt.Println(a)
}

// 7. testmap
func TestMap(t *testing.T) {
	m := map[string]string{
		"Thu, 09 Jul 2020 01:45:00 GMT": "3",
		"aa":                            "1",
		"bb":                            "2",
	}
	fmt.Println(m)
}

func TestDagList(t *testing.T) {
	dag := models.DagAirFlow{}
	dag.UpdateStatus()
}

func Connect(retry int) {
	if retry <= 0 {
		panic("重试完毕, 不能再重试了")
	}
	retry = retry - 1
	fmt.Println("开始重试, 还剩 ", retry)
	Connect(retry)
}

func TestConnect(t *testing.T) {
	Connect(3)
}
