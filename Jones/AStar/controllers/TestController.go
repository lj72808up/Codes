package controllers

import (
	"AStar/models"
	"AStar/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
)

type TestController struct {
	BaseController
}

var TestNameSpace = "test"

// @router /testGet/:id [get]
func (this *TestController) GetAll() {
	//idStr := this.Ctx.Input.Param(":id")
	user := "zhangsan"
	addDomainForUser(user, "lvpixxx")
	//println(getDomainForUser(user))

	models.Enforcer.GetAllRoles()
	//e := models.Enforcer
	//e.AddPolicy("ciciGroup", "/jones/cici/*", "*")
	//e.AddRoleForUserInDomain("zhangsan", "dummyGroup", "lvpi2")  //  AddNamedGroupingPolicy
	//res := e.GetFilteredNamedGroupingPolicy("g",0,"zhangsan")

	//m := e.GetModel()
	//fmt.Println(getDomainForUser("zhangsan"))
	//fmt.Println(models.Enforcer.GetRolesForUser(user))

	fmt.Println(models.Enforcer.GetFilteredGroupingPolicy(0))
	fmt.Println(getAllDomain())
}

func getAllDomain() []string {
	e := models.Enforcer
	res := e.GetFilteredGroupingPolicy(3)
	domains := []string{}
	for _, info := range res {
		if len(info) >= 3 {
			domains = append(domains, info[2])
		}
	}
	return utils.MakeUniqSet(domains...)
}

func getDomainForUser(userName string) []string {
	e := models.Enforcer
	res := e.GetFilteredNamedGroupingPolicy("g", 0, userName)
	domains := []string{}
	for _, info := range res {
		if len(info) >= 3 {
			domains = append(domains, info[2])
		}
	}
	return utils.MakeUniqSet(domains...)
}

func addDomainForUser(userName string, domain string) {
	e := models.Enforcer
	e.AddRoleForUserInDomain(userName, "dummyGroup", domain)
}

// @router /testPost [post]
func (this *TestController) PostAll() {
	file, header, _ := this.GetFile("file")
	name := this.GetString("name")
	bytes, _ := utils.ReadBytes(file, header.Size)
	fmt.Println(string(bytes))

	fmt.Println(file)
	fmt.Println(header.Filename)
	fmt.Println(header.Size)
	fmt.Println(name)
}

// @router /testGetPermission [get]
func (this *TestController) GetCanWrite() {
	e := models.Enforcer
	roles := e.GetRolesForUser("zhangsan2")

	fmt.Println(roles)
	o := orm.NewOrm()
	var rows []string

	placeholder := makeRolePlaceHolder(roles)
	sql := fmt.Sprintf("SELECT v1 FROM casbin_rule WHERE v0 in %s and v2=?", placeholder)

	println(sql)
	_, _ = o.Raw(sql, roles, "*").QueryRows(&rows)
	for _, row := range rows {
		println(row)
	}
}

// 返回 (?,?,?) 的形式
func makeRolePlaceHolder(roles []string) string {
	var roleParams []string
	for range roles {
		roleParams = append(roleParams, "?")
	}
	return fmt.Sprintf("(%s)", strings.Join(roleParams, ","))
}

func (this *TestController) testCasbin() {
	// subject和role  ==>  用户组
	// object ==> url
	// action  ==> Get / Post

	e := models.Enforcer
	subjects := e.GetAllSubjects() // 获取所有subject(用户组)
	fmt.Println(subjects)

	allNamedSubjects := e.GetAllNamedSubjects("p") //根据p_type展示subject, 不能填写g, 因为g只配置了3项
	fmt.Println(allNamedSubjects)

	// 获取object (相当于获取配置的url)
	// select v1 from casbin_rule where p_type='p' group by v1
	allObjects := e.GetAllObjects()
	fmt.Println(allObjects)
	fmt.Println(len(allObjects))

	// 获取所有action (所有不同中的操作)
	// select v2 from casbin_rule where p_type='p' group by v2
	allActions := e.GetAllActions()
	fmt.Println(allActions)
	fmt.Println(len(allActions))

	// 获取所有角色 (这里痛subject)
	allRoles := e.GetAllRoles()
	fmt.Println(allRoles)
	fmt.Println(len(allRoles))

	// 获取所有策略 (用户组, url, 读写)的匹配对
	policy := e.GetPolicy()
	fmt.Println(policy)
	fmt.Println(len(policy))

	// 根据条件获取策略
	// select * from casbin_rule where p_type='p' and v0='adminGroup'
	filteredNamedPolicy := e.GetFilteredNamedPolicy("p", 0, "jdbcCheckGroup")
	fmt.Println(filteredNamedPolicy)
	fmt.Println(len(filteredNamedPolicy))

	// 获取角色下的用户列表
	users := e.GetUsersForRole("jdbcCheckGroup")
	fmt.Println(users)
	fmt.Println(len(users))

	// 增加一个用户组 (并赋予一个权限)
	areRulesAdded := e.AddPolicy("jdbcCheckGroup", "/jones/metaData/*", "*")
	fmt.Printf("增加页面权限:%t\n", areRulesAdded)

	readWriteAdded := e.AddPolicy("jdbcCheckGroup", "/jones/metaData_authority", "GET")
	fmt.Printf("增加读写权限:%t\n", readWriteAdded)

	// 删除用户组的所有权限, 但保留该用户组下与用户的映射
	filteredNamedPolicy2 := e.GetFilteredNamedPolicy("p", 0, "hahaGroup")
	for i, _ := range filteredNamedPolicy2 {
		groupName := filteredNamedPolicy2[i][0]
		url := filteredNamedPolicy2[i][1]
		action := filteredNamedPolicy2[i][2]
		fmt.Println("Delete", groupName, url, action)
		e.RemovePolicy(groupName, url, action)
	}
}
