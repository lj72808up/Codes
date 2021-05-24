package controllers

import (
	"AStar/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type LvpiAuthController struct {
	BaseController
}

// @Title 查看所有域
// @Summary 查看所有域
// @Description 查看所有域
// @Success 200
// @Accept json
// @router /getAllDomain [get]
func (this *LvpiAuthController) GetAllDomain() {
	domains := models.GetAllDomain()
	bytes, _ := json.Marshal(domains)
	this.Ctx.WriteString(string(bytes))
}

// @Title 获取有权管理的一级域
// @Summary 获取有权管理的一级域
// @Description 获取有权管理的一级域
// @Success 200
// @Accept json
// @router /getDomainBySelfAdmin [get]
func (this *LvpiAuthController) GetDomainBySelfAdmin() {
	firstDomains := models.GetDomainByAdministrator(this.userName)
	bytes, _ := json.Marshal(firstDomains)
	this.Ctx.WriteString(string(bytes))
}

// @Title 获取域名下的二级绿皮
// @Summary 获取域名下的二级绿皮
// @Description 获取域名下的二级绿皮
// @Success 200
// @Accept json
// @router /getLvpiBySuperior/:pid [get]
func (this *LvpiAuthController) GetLvpiBySuperior() {
	pid, _ := strconv.ParseInt(this.Ctx.Input.Params()[":pid"], 10, 64)
	lvpis := models.GetLvpiByPid(pid)
	bytes, _ := json.Marshal(lvpis)
	this.Ctx.WriteString(string(bytes))
}

// @Title 获取可见的一级绿皮和二级绿皮
// @Summary 获取可见的一级绿皮和二级绿皮
// @Description 获取可见的一级绿皮和二级绿皮
// @Success 200
// @Accept json
// @router /getFirstDomainLvpiVisible [get]
func (this *LvpiAuthController) GetFirstDomainLvpiVisible() {
	pid, _ := strconv.ParseInt(this.Ctx.Input.Params()[":pid"], 10, 64)
	lvpis := models.GetLvpiByPid(pid)
	bytes, _ := json.Marshal(lvpis)
	this.Ctx.WriteString(string(bytes))
}

// 左侧菜单兰
// @Title 获取可见的一级绿皮和二级绿皮
// @Summary 获取可见的一级绿皮和二级绿皮
// @Description 获取可见的一级绿皮和二级绿皮
// @Success 200
// @Accept json
// @router /getlvpiTree [get]
func (this *LvpiAuthController) GetlvpiTree() {
	//todo 查 canUserWatch 是否包含自己? 包含的为二级绿皮, 一级域名为pid
	userId := this.userName
	views := []lvpiTreeView{}
	secondLvpis := []models.DomainLvpi{}
	o := orm.NewOrm()
	_, _ = o.Raw("SELECT * FROM domain_lvpi WHERE find_in_set(?,user_can_watch) > 0", userId).QueryRows(&secondLvpis)
	firstDomains := []models.DomainLvpi{}
	firstSql := `
SELECT t1.* FROM domain_lvpi t1 LEFT JOIN 
(SELECT pid FROM domain_lvpi WHERE find_in_set(?,user_can_watch) > 0 GROUP BY pid) t2
ON 
t1.did = t2.pid 
WHERE t2.pid IS NOT NULL 
`
	_, _ = o.Raw(firstSql, userId).QueryRows(&firstDomains)
	// todo 组建树
	for i, first := range firstDomains {
		children := []models.DomainLvpi{}
		for _, second := range secondLvpis {
			if second.Pid == first.Did && second.Valid == 1 {
				children = append(children, second)
			}
		}
		view := lvpiTreeView{
			Did:           first.Did,
			Name:          first.Name,
			Pid:           0,
			Index:         fmt.Sprint(i+1),
			Info:          first.Info,
			Children:      children,
		}
		views = append(views, view)
	}

	bytes,_ := json.Marshal(views)
	this.Ctx.WriteString(string(bytes))
}

type lvpiTreeView struct {
	Did           int64
	Name          string
	Administrator string
	Pid           int64
	Index         string
	Info          string
	Children      []models.DomainLvpi
}

//删除
// @Title 更新绿皮权限映射
// @Summary 更新绿皮权限映射
// @Description 更新绿皮权限映射
// @Success 200
// @Accept json
// @router /getFirstDomainLvpiVisible/uid [post]
func (this *LvpiAuthController) PostLvpiMapping() {
	mappings := []models.LvpiMapping{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &mappings)
	models.AddUserMapping(mappings)
}

// @Title 添加二级绿皮
// @Summary  添加二级绿皮
// @Description  添加二级绿皮
// @Accept json
// @router /postSecondLvpi [post]
func (this *LvpiAuthController) PostSecondLvpi() {
	var lvpis []models.DomainLvpi
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &lvpis)
	o := orm.NewOrm()
	_, _ = o.Raw("DELETE FROM domain_lvpi WHERE pid = ?", lvpis[0].Pid).Exec()
	for _, lvpi := range lvpis {
		fmt.Println(lvpi)
		_, _ = o.Raw("REPLACE INTO domain_lvpi(name,pid,`index`,info, user_can_watch, valid, lid) VALUE (?,?,?,?,?,?,?)",
			lvpi.Name, lvpi.Pid, lvpi.Index, lvpi.Info, lvpi.UserCanWatch, lvpi.Valid, lvpi.Lid).Exec()
	}
	this.Ctx.WriteString("添加完成")
}

// @Title 添加二级绿皮
// @Summary  添加二级绿皮
// @Description  添加二级绿皮
// @Accept json
// @router /postSecondLvpiUpdate [post]
func (this *LvpiAuthController) PostSecondLvpiUpdate() {
	var lvpis []models.DomainLvpi
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &lvpis)
	o := orm.NewOrm()
	for _, lvpi := range lvpis {
		fmt.Println(lvpi)
		_, _ = o.Raw("REPLACE INTO domain_lvpi(name,pid,`index`,info, user_can_watch, valid, lid) VALUE (?,?,?,?,?,?,?)",
			lvpi.Name, lvpi.Pid, lvpi.Index, lvpi.Info, lvpi.UserCanWatch, lvpi.Valid, lvpi.Lid).Exec()
	}
	this.Ctx.WriteString("添加完成")
}


// @Title 查看域下包含的绿皮
// @Summary  查看域下包含的绿皮
// @Description  查看域下包含的绿皮
// @Accept json
// @router /getLvpiInDomain/:did [post]
func (this *LvpiAuthController) GetLvpiInDomain() {
	did, _ := strconv.ParseInt(this.Ctx.Input.Params()[":did"], 10, 64)
	lvpis := []models.DomainLvpi{}
	o:= orm.NewOrm()
	_, _ = o.Raw("Select * from domain_lvpi WHERE pid = ?", did).QueryRows(&lvpis)
	bytes, _ := json.Marshal(lvpis)
	this.Ctx.WriteString(string(bytes))
}

// @Title 查看所有的绿皮
// @Summary  查看所有的绿皮
// @Description  查看所有的绿皮
// @Accept json
// @router /getLvpis/ [get]
func (this *LvpiAuthController) GetLvpis() {
	lvpis := []models.DomainLvpi{}
	o:= orm.NewOrm()
	_, _ = o.Raw("Select * from domain_lvpi WHERE 1 = 1").QueryRows(&lvpis)
	bytes, _ := json.Marshal(lvpis)
	this.Ctx.WriteString(string(bytes))
}
