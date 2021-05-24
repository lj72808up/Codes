package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type DomainLvpi struct {
	Did           int64  `orm:"pk;auto" field:"id,hidden"`
	Lid			  int64
	Name          string `orm:"size(255)"` // 域名或绿皮名
	Administrator string `orm:"size(255)"` // 管理员
	Pid           int64  // 绿皮层级
	// 参照 top_label 的跳转逻辑
	Index        string `orm:"type(text)"` // 页面跳转url
	Info         string `orm:"type(text)"` // 具体的绿皮地址
	UserCanWatch string `orm:"size(255)"`  // 可见绿皮的用户
	Valid	     int32
}

func (u *DomainLvpi) TableUnique() [][]string {
	return [][]string{
		[]string{"Name", "Pid"},
	}
}

type LvpiMapping struct {
	Lid    int64  `orm:"pk;auto" field:"id,hidden"`
	UserId string `orm:"size(255)"` // 用户名
	LvpiId int64
	Pid    int64 // 冗余字段, 是二级绿皮的父目录
	Valid  bool
}

func (u *LvpiMapping) TableUnique() [][]string {
	return [][]string{
		[]string{"UserId", "LvpiId"},
	}
}

type DomainLvpiView struct {
	Did      int64
	Name     string
	Pid      int64
	Index    string
	Info     string
	Children []DomainLvpi
}

// 获取有权限的标签(分级)
func GetCheckedDomainLvpi(uid string) []DomainLvpiView {
	firstDomains := GetDomainForUSer(uid)
	res := []DomainLvpiView{}
	for _, domain := range firstDomains {
		secondLvpis := GetLvpiForUser(uid, domain.Did)
		firstView := DomainLvpiView{
			Did:      domain.Did,
			Name:     domain.Name,
			Pid:      domain.Pid,
			Index:    domain.Index,
			Info:     domain.Info,
			Children: secondLvpis,
		}
		res = append(res, firstView)
	}

	return res
}

func GetAllDomain() []DomainLvpi {
	o := orm.NewOrm()
	domains := []DomainLvpi{}
	_, _ = o.Raw("SELECT * FROM domain_lvpi WHERE pid = 0").QueryRows(&domains)
	return domains
}

// 获取可见二级绿皮
func GetLvpiForUser(userId string, pid int64) []DomainLvpi {
	o := orm.NewOrm()
	lvpis := []DomainLvpi{}
	_, _ = o.Raw(`SELECT t2.* FROM 
(SELECT lvpi_id FROM lvpi_mapping WHERE user_id=? AND pid=? AND valid=TRUE) t1
LEFT JOIN 
domain_lvpi t2
ON 
t2.did = t1.lvpi_id`, userId, pid).QueryRows(&lvpis)
	return lvpis
}

// 获取可见一级域
func GetDomainForUSer(userId string) []DomainLvpi {
	o := orm.NewOrm()
	lvpis := []DomainLvpi{}
	_, _ = o.Raw(`SELECT t2.* FROM
	(SELECT lvpi_id FROM lvpi_mapping WHERE user_id=? AND pid=0 AND valid=TRUE) t1
	LEFT JOIN
	domain_lvpi t2
	ON
	t2.did = t1.lvpi_id`, userId).QueryRows(&lvpis)
	return lvpis

}

func GetDomainByAdministrator(userId string) []DomainLvpi {
	o := orm.NewOrm()
	res := []DomainLvpi{}
	sql := fmt.Sprintf("SELECT * FROM domain_lvpi WHERE find_in_set('%s',Administrator) > 0", userId)
	_, _ = o.Raw(sql).QueryRows(&res)
	return res
}

func IsAdministrator(userId string) bool {
	// 是否是有管理域的管理员
	adminisDomains := GetDomainByAdministrator(userId)
	return (adminisDomains != nil) && (len(adminisDomains) != 0)
}

// 增加绿皮id和uid映射
func AddUserMapping(mappings []LvpiMapping) {
	o := orm.NewOrm()
	for _, mappng := range mappings {
		_, _ = o.Raw("UPDATE lvpi_mapping SET valid=false WHERE user_id = ?", mappng.UserId).Exec()
		_, _ = o.Raw(`
REPLACE INTO lvpi_mapping(user_id, lvpi_id, valid) VALUES(?,?,true)`, mappng.UserId, mappng.LvpiId).Exec()

	}
}

func GetDomainLvpiById(did int64) DomainLvpi {
	o := orm.NewOrm()
	lvpi := DomainLvpi{}
	_ = o.Raw("SELECT * FROM domain_lvpi WHERE did = ?", did).QueryRow(&lvpi)
	return lvpi
}

func (this *DomainLvpi) Add() {
	o := orm.NewOrm()
	_, _ = o.Insert(this)
}

// 获取一级域下所有绿皮
func GetLvpiByPid(pid int64) []DomainLvpi {
	o := orm.NewOrm()
	lvpi := []DomainLvpi{}
	_, _ = o.Raw("SELECT * FROM domain_lvpi WHERE pid = ?", pid).QueryRows(&lvpi)
	return lvpi
}
