package main

import (
	_ "Beego_demo/routers"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main2() {
	beego.Run()
}










type User struct {
	Id          int
	Name        string
	Profile     *Profile   `orm:"rel(one)"` // OneToOne relation
	Post    	[]*Post `orm:"reverse(many)"` // 设置一对多的反向关系
}

type Profile struct {
	Id          int
	Age         int16
	User        *User   `orm:"reverse(one)"` // 设置一对一反向关系(可选)
}

type Post struct {
	Id    int
	Title string
	User  *User  `orm:"rel(fk)"`	//设置一对多关系
}
func init2() {
	// 参数4(可选)  设置最大空闲连接
	// 参数5(可选)  设置最大数据库连接 (go >= 1.2)
	maxIdle := 30
	maxConn := 30

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "lj:123456@tcp(localhost:3306)/test_db?charset=utf8", maxIdle, maxConn)
	// 需要在init中注册定义的model
	orm.RegisterModel(new(User), new(Profile), new(Post))

	/////////////////////自动建表////////////////////////////////
	// 数据库别名
	name := "default"
	// drop table 后再建表
	force := false
	// 打印执行过程
	verbose := true
	// 遇到错误立即返回
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
	/////////////////////////////////////////////////////
}

func main1(){
	o := orm.NewOrm()
	o.Using("test_db") // 默认使用 default，你可以指定为其他数据库

	profile := new(Profile)
	profile.Age = 30

	user := new(User)
	user.Profile = profile
	user.Name = "slene"

	fmt.Println(o.Insert(profile))
	fmt.Println(o.Insert(user))
}


func main(){
	taskFun := func() error { fmt.Println(time.Now().String(),"tk1"); return nil }
	tk1 := toolbox.NewTask("tk1", "5 * * * * *", taskFun)
	toolbox.AddTask("tk1", tk1)
	toolbox.StartTask()
	defer toolbox.StopTask()
	beego.Run()
}

