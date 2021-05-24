package test

import (
	"AStar/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"regexp"
	"testing"
)

func TestOrm(t *testing.T) {
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	_ = orm.RegisterDataBase("default", "mysql", "adtd:noSafeNoWork2016@tcp(adtd.datacenter.rds.sogou:3306)/datacenter?charset=utf8&loc=Asia%2FShanghai")
	orm.RegisterModel(new(models.DagAirFlow))

	o := orm.NewOrm()
	//var dags []airflowDag.DagAirFlow
	var dags []map[string]interface{}
	_, _ = o.Raw("select * from dag_air_flow").QueryRows(&dags)
	fmt.Println(dags)
}

func TestUrl(t *testing.T) {
	//url := "yuhancheng:yuhancheng@tcp(10.139.36.81:3306)/adtl_test?charset=utf8&loc=Asia%2FShanghai"
	url := beego.AppConfig.String("datasourcename")
	r, _ := regexp.Compile("^(.*):(.*)@tcp\\((.*)\\)/(.*)\\?(.*)$")
	params := r.FindStringSubmatch(url)

	user := params[1]
	passwd := params[2]
	hostPort := params[3]
	dataBase := params[4]

	jdbcUrl := fmt.Sprintf("jdbc:mysql://%s/%s?user=%s&password=%s&useUnicode=true&characterEncoding=UTF-8",
		hostPort, dataBase, user, passwd)
	fmt.Println(jdbcUrl)
}
