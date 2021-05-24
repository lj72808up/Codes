package models

import (
	"AStar/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type BodongMenuBar struct {
	Mid    int64  `orm:"pk;auto"`
	Label  string `orm:"size(255)"`
	Pid    int64
	Sortid int64
	Img    string `orm:"size(255)"`
	Url    string `orm:"type(text)"`
}

type BodongMenuNode struct {
	Mid      int64            `json:"mid"`
	Label    string           `json:"label"`
	Children []BodongMenuNode `json:"children,omitempty"`
}

func (this *BodongMenuBar) GetFirstMenu(roles []string) []BodongMenuBar {
	o := orm.NewOrm()
	var menus []BodongMenuBar
	roleWildCard := utils.MakeRolePlaceHolder(roles)
	sql := fmt.Sprintf(`SELECT t1.* FROM bodong_menu_bar t1 
LEFT JOIN 
(SELECT menu_id FROM bodong_menu_mapping WHERE group_name in %s group by menu_id) t2
ON 
t1.mid = t2.menu_id
WHERE t1.pid = 0 AND t2.menu_id IS NOT NULL 
ORDER BY sortid asc`, roleWildCard)
	_, _ = o.Raw(sql, roles).QueryRows(&menus)
	return menus
}

func (this *BodongMenuBar) GetChildMenus(mid int64, roles []string) []BodongMenuBar {
	// todo
	o := orm.NewOrm()
	var menus []BodongMenuBar
	roleWildCard := utils.MakeRolePlaceHolder(roles)
	sql := fmt.Sprintf(`SELECT t1.* FROM 
(SELECT * FROM bodong_menu_bar WHERE pid = ?) t1
LEFT JOIN 
(SELECT * FROM bodong_menu_mapping WHERE group_name in %s) t2
ON t1.mid = t2.menu_id
WHERE t2.menu_id IS NOT NULL 
ORDER BY sortid asc`, roleWildCard)
	_, _ = o.Raw(sql, mid,  roles).QueryRows(&menus)
	return menus
}

type BodongMenuMapping struct {
	Id        int64  `orm:"pk;auto"`
	GroupName string `orm:"size(50)"`
	MenuId    int64
}

func UpdateBodongMapping(role string, menuIds []string) error {
	o := orm.NewOrm()
	_, _ = o.Raw("DELETE FROM bodong_menu_mapping WHERE group_name = ?", role).Exec()
	mappings := []BodongMenuMapping{}
	for _, menuIdStr := range menuIds {
		menuId, _ := strconv.ParseInt(menuIdStr, 10, 64)
		mapping := BodongMenuMapping{
			GroupName: role,
			MenuId:    menuId,
		}
		mappings = append(mappings, mapping)
	}
	_, err := o.InsertMulti(10, mappings)
	return err
}

func (this *BodongMenuBar) GetCheckedMenuIds(roles []string) []int64 {
	o := orm.NewOrm()
	sqlItems := []string{}
	for _, role := range roles {
		sql := fmt.Sprintf("SELECT menu_id FROM bodong_menu_mapping WHERE group_name = '%s'", role)
		sqlItems = append(sqlItems, sql)
	}
	finalSql := strings.Join(sqlItems, "\nunion\n")
	menuIds := []int64{}
	_, _ = o.Raw(finalSql).QueryRows(&menuIds)
	return menuIds
}

type BoDongTreeNode struct {
	Mid      int64            `json:"menuId"`
	Pid      int64            `json:"pid"`
	Label    string           `json:"label"`
	Children []BoDongTreeNode `json:"children,omitempty"`
}

func GetBoDongMenuTree() []BoDongTreeNode {
	o := orm.NewOrm()
	allNode := []BoDongTreeNode{}
	_, _ = o.Raw("SELECT * FROM bodong_menu_bar").QueryRows(&allNode)
	firstLevel := []BoDongTreeNode{} // 首层节点
	for _, node := range allNode {
		if node.Pid == 0 {
			firstLevel = append(firstLevel, node)
		}
	}

	for i, _ := range firstLevel {
		var children []BoDongTreeNode
		for _, node := range allNode {
			if node.Pid == firstLevel[i].Mid {
				children = append(children, node)
			}
		}
		firstLevel[i].Children = children
	}

	return firstLevel
}
