package models

import (
	"AStar/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type LabelModel struct{}

type TopLabel struct {
	Lid       int64  `orm:"pk;auto"`
	Title     string `orm:"size(50)"`
	Class     string `orm:"size(50)"`
	NameSpace string `orn:"size(50)"` // 控制页面的读写, 数据库中要和controller中的namespace相同
	Pid       int64  `orm:"default(0)"`
	Index     string `orm:"size(50)"`
	Info      string `orm:"size(1024)"`
	SortId    int64
}

type NameSpaceWritable struct {
	Nid       int64  `orm:"pk;auto"`
	RoleName  string `orm:"size(50)"`
	NameSpace string `orm:"size(50)"`
	Lid       int64
}

type TopLabelCheckedView struct {
	Lid          int64
	Pid          int64
	Title        string
	Index        string
	NameSpace    string
	IsWritable   int64
	HasPrivilege bool
}

type PmLabelView struct {
	Lid    int64
	Pid    int64
	Title  string
	Class  string
	Ptitle string
	SortId int64
}

type RdLabelView struct {
	Lid       int64
	Title     string
	Index     string
	NameSpace string
}

type GroupLabelMapping struct {
	Glid      int64  `orm:"pk;auto"`
	GroupName string `orm:"size(255)"`
	LabelId   string `orm:"size(50)"`
	Valid     bool   `orm:"default(true)"`
}

// 多字段索引
func (u *GroupLabelMapping) TableIndex() [][]string {
	return [][]string{
		{"GroupName", "Valid"},
	}
}

// 多字段唯一键
func (u *GroupLabelMapping) TableUnique() [][]string {
	return [][]string{
		[]string{"GroupName", "LabelId"},
	}
}


// 获取所有顶层标签并标出是否有权限(2层拉平成一层)
func (LabelModel) GetAllTopLabel(group string) []TopLabelCheckedView {
	o := orm.NewOrm()
	var firstLevel []TopLabelCheckedView
	var secondLevel []TopLabelCheckedView

	firstSql := `
SELECT
	t1.lid,
	t1.pid,
	t1.title,
	t1.index,
	t1.name_space,
	case 
        when t3.lid is null THEN 2
        else 1 
    END as is_writable,
	case 
        when t2.label_id is null THEN false
        else true 
    END as has_privilege
FROM
	(SELECT * FROM top_label where pid=0) t1
LEFT JOIN
	(SELECT label_id FROM group_label_mapping WHERE valid = true and group_name = ? group by label_id) t2
ON 
	t1.lid = t2.label_id
LEFT JOIN 
	(SELECT role_name,name_space,lid FROM name_space_writable WHERE role_name = ? GROUP BY role_name,name_space,lid) t3
ON
	t1.lid = t3.lid
ORDER BY t1.lid asc
`
	_, _ = o.Raw(firstSql, group, group).QueryRows(&firstLevel)

	secondSql := `
SELECT
	t1.lid,
	t1.pid,
	t1.title,
	t1.index,
	t1.name_space,
	case 
        when t3.lid is null THEN 2
        else 1 
    END as is_writable,
	case 
        when t2.label_id is null THEN false
        else true 
    END as has_privilege
FROM
	(SELECT * FROM top_label where pid != 0) t1
LEFT JOIN
	(SELECT label_id FROM group_label_mapping WHERE valid = true and group_name = ? group by label_id) t2
ON 
	t1.lid = t2.label_id
LEFT JOIN 
	(SELECT role_name,name_space,lid FROM name_space_writable WHERE role_name = ? GROUP BY role_name,name_space,lid) t3
ON
	t1.lid = t3.lid
ORDER BY t1.lid asc
`
	_, _ = o.Raw(secondSql, group, group).QueryRows(&secondLevel)
	//拉平成同一层
	var nodes []TopLabelCheckedView
	for _, x := range firstLevel {
		nodes = append(nodes, x)
		for _, y := range secondLevel {
			if y.Pid == x.Lid {
				y.Title = fmt.Sprintf("%s - %s", x.Title, y.Title)
				nodes = append(nodes, y)
			}
		}
	}
	fmt.Println(nodes)
	return nodes
}

// 获取PM关注的标签属性title、class、pid
func (LabelModel) GetPMLabelAttr() []PmLabelView {
	o := orm.NewOrm()
	var firstLevel []PmLabelView
	firstSql := `
SELECT
	a.lid,
	a.pid,
	a.title,
	a.class,
IF
	( a.pid = 0, "根节点", b.title ) AS ptitle ,
	a.sort_id
FROM
	top_label a
	LEFT JOIN top_label b ON a.pid = b.lid
`
	_, _ = o.Raw(firstSql).QueryRows(&firstLevel)
	return firstLevel
}

// 更新PM关注的标签属性title、class、pid
func (LabelModel) UpdatePMLabelAttr(lid string, title string, class string, pid string) bool {
	flag := true
	o := orm.NewOrm()

	urlSql := fmt.Sprintf(" Update top_label set title='%s',`class`='%s',pid='%s' where lid='%s'", title, class, pid, lid)
	res, err := o.Raw(urlSql).Exec()
	if err != nil {
		fmt.Printf("UpdatePMLabelAttr get err:%s | res:%s", err, res)
		flag = false
	}
	return flag
}

// 获取RD关注的标签属性index|name_space
func (LabelModel) GetRDLabelAttr() []RdLabelView {
	o := orm.NewOrm()
	var firstLevel []RdLabelView

	urlSql := fmt.Sprintf("select lid,title,`index`,name_space from top_label")
	_, _ = o.Raw(urlSql).QueryRows(&firstLevel)
	return firstLevel
}

// 新增RD自定义索引和命名空间的 标签
func (LabelModel) AddRDLabel(title string, index string, namespace string) int {
	newLid := -1
	o := orm.NewOrm()

	urlSql := fmt.Sprintf(" INSERT into top_label (title,`index`,name_space,pid) VALUES('%s','%s','%s','0') ", title, index, namespace)
	res, err := o.Raw(urlSql).Exec()
	if err != nil {
		fmt.Printf("UpdatePMLabelAttr get err:%s | res:%s", err, res)
		newLid = -1
		return newLid
	}
	//插入成功， 则返回lid
	urlSql = fmt.Sprintf(" select Lid from top_label where title='%s' and `index`='%s' and name_space='%s' and pid='0' limit 1 ", title, index, namespace)
	err = o.Raw(urlSql).QueryRow(&newLid)
	if err != nil {
		fmt.Printf("UpdatePMLabelAttr get err:%s | res:%s", err, res)
		newLid = -1
	}
	return newLid
}

// 更新RD关注的标签属性index namespace
func (LabelModel) UpdateRDLabelAttr(lid string, index string, namespace string) bool {
	flag := true
	o := orm.NewOrm()

	urlSql := fmt.Sprintf(" Update top_label set `index`='%s',`name_space`='%s' where lid='%s'", index, namespace, lid)
	res, err := o.Raw(urlSql).Exec()
	if err != nil {
		fmt.Printf("UpdatePMLabelAttr get err:%s | res:%s", err, res)
		flag = false
	}
	return flag
}

type TopLabelView struct {
	Lid      int64
	Title    string
	Class    string
	Pid      int64
	Index    string
	Info     string
	Children []TopLabel
}

// 获取有权限的标签(分级)
func (LabelModel) GetCheckedTopLabel(roles []string) []TopLabelView {
	o := orm.NewOrm()
	var labels []TopLabelView
	roleWildCard := utils.MakeRolePlaceHolder(roles)
	firstSql := fmt.Sprintf(`
SELECT
	t2.lid,
	t2.title,
	t2.class,
	t2.pid,
	t2.index,
	t2.info
FROM
	(SELECT label_id FROM group_label_mapping WHERE valid = true and group_name in %s group by label_id) t1
LEFT JOIN 
	(SELECT * FROM top_label where pid=0) t2
on
 	t2.lid = t1.label_id
WHERE
 	t2.lid is not null
order by t2.sort_id asc`, roleWildCard)
	secondSql := fmt.Sprintf(`
SELECT
	t2.lid,
	t2.title,
	t2.class,
	t2.pid,
	t2.index,
	t2.info
FROM
	(SELECT label_id FROM group_label_mapping WHERE valid = true and group_name in %s group by label_id) t1
LEFT JOIN 
	(SELECT * FROM top_label where pid != 0) t2
on
 	t2.lid = t1.label_id
WHERE
	t2.lid is not null
order by t2.sort_id asc`, roleWildCard)

	var firstLevel []TopLabel
	var secondLevel []TopLabel
	_, _ = o.Raw(firstSql, roles).QueryRows(&firstLevel)
	_, _ = o.Raw(secondSql, roles).QueryRows(&secondLevel)

	for _, x := range firstLevel {
		var children []TopLabel
		for _, y := range secondLevel {
			if y.Pid == x.Lid {
				children = append(children, y)
			}
		}
		labels = append(labels, TopLabelView{
			Lid:      x.Lid,
			Title:    x.Title,
			Class:    x.Class,
			Pid:      x.Pid,
			Index:    x.Index,
			Info:     x.Info,
			Children: children,
		})
	}
	return labels
}

func (this *LabelModel) deleteLabelMapping(o orm.Ormer, groupName string) {
	// 1. 删除标签映射
	res, err := o.Raw("DELETE FROM group_label_mapping WHERE group_name = ?", groupName).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		Logger.Info("group worksheet mapping affected nums: ", num)
	} else {
		panic(err)
	}
	// 2. 删除角色对应的权限
	e := Enforcer
	filteredNamedPolicy := e.GetFilteredNamedPolicy("p", 0, groupName)
	// 删除用户组的所有权限, 但保留该用户组下与用户的映射
	for i, _ := range filteredNamedPolicy {
		groupName := filteredNamedPolicy[i][0]
		url := filteredNamedPolicy[i][1]
		action := filteredNamedPolicy[i][2]
		fmt.Println("Delete", groupName, url, action)
		e.RemovePolicy(groupName, url, action)
	}

	// 3. 删除只读映射表
	sql := "DELETE FROM name_space_writable WHERE role_name = ?"
	_, _ = o.Raw(sql, groupName).Exec()
}

func (this *LabelModel) addLabel(o orm.Ormer, groupName string, worksheetIds []string, writableUrls []string) {
	// 1. 增加标签映射
	var mappings []GroupLabelMapping
	for _, id := range worksheetIds {
		mapping := GroupLabelMapping{
			GroupName: groupName,
			LabelId:   id,
			Valid:     true,
		}
		mappings = append(mappings, mapping)
	}
	_, err2 := o.InsertMulti(100, mappings)
	if err2 != nil {
		panic(err2)
	}

	// 2. 增加casbin权限
	e := Enforcer
	for _, uid := range worksheetIds {
		uidByInt, _ := strconv.ParseInt(uid, 10, 64)
		label := TopLabel{Lid: uidByInt}
		_ = o.Read(&label)

		nameSpace := utils.GenNameSpace("/jones", label.NameSpace)
		casbinURL := fmt.Sprintf("%s/%s", nameSpace, "*")
		action := "*"
		e.AddPolicy(groupName, casbinURL, action)
		fmt.Println("Add Policy:", groupName, casbinURL, action)
	}

	// 3. 增加可写权限
	for _, uid := range writableUrls {
		uidByInt, _ := strconv.ParseInt(uid, 10, 64)
		label := TopLabel{Lid: uidByInt}
		_ = o.Read(&label)

		casbinURL := fmt.Sprintf("%s/%s_%s", "/jones", label.NameSpace, "authority")
		action := "POST"
		e.AddPolicy(groupName, casbinURL, action)
		fmt.Println("Add Policy:", groupName, casbinURL, action)

		// 插入可写表
		writables := NameSpaceWritable{
			RoleName:  groupName,
			NameSpace: label.NameSpace,
			Lid:       uidByInt,
		}

		if _, err := o.Insert(&writables); err != nil {
			panic(err)
		}
	}
}

func AddDefaultaCasbin(groupName string) {
	e := Enforcer
	defaultURLs := []string{
		"/jones/auth/labels/*", "/jones/auth/routes", "/jones/redash/timestamp", "/jones/auth/roles", "/jones/jone/*",
		"/jones/metaData/*","/jones/task_god/*","/jones/mobileDashboard/*","/jones/bodong/*","/jones/indicator/*","/jones/redash/*",
	}
	for _, url := range defaultURLs {
		action := "*"
		e.AddPolicy(groupName, url, action)
		fmt.Println("Add Policy:", groupName, url, action)
	}
}

// 更新用户的标签权限
func (this *LabelModel) AddLabelMapping(groupName string, worksheetIds []string, writableUrls []string) {
	o := orm.NewOrm()
	// 1. 清除标签映射及对应的casbin权限,只读权限
	this.deleteLabelMapping(o, groupName)

	// 2. 增加新的标签映射,casbin权限,只读权限
	this.addLabel(o, groupName, worksheetIds, writableUrls)

	// 3. 增加一些默认的acsbin权限
	AddDefaultaCasbin(groupName)

	// 4. 增加默认的标签
	defaultLabels := []string{"需求管理"}
	AddLabelsForUser([]string{groupName},defaultLabels)
}

func AddLabelsForUser(roles []string, labelNames []string) {
	e := Enforcer
	for _, label := range labelNames {
		label := getLabelIdByName(label)
		for _, role := range roles {
			if label.Pid != 0 {
				// 父级目录没有命名空间
				replaceLabelMapping(role, fmt.Sprintf("%d", label.Pid))
			}
			if label.Lid != 0 {
				// 子目录有命名空间
				nameSpace := utils.GenNameSpace("/jones", label.NameSpace)
				casbinURL := fmt.Sprintf("%s/%s", nameSpace, "*")
				action := "*"
				e.AddPolicy(role, casbinURL, action)
				replaceLabelMapping(role, fmt.Sprintf("%d", label.Lid))
			}
		}
	}
}

func getLabelIdByName(labelName string) TopLabel {
	o := orm.NewOrm()
	label := TopLabel{}
	_ = o.Raw("SELECT * FROM top_label WHERE title = ?", labelName).QueryRow(&label)
	return label
}

func replaceLabelMapping(role string, labelId string) {
	o := orm.NewOrm()
	_, _ = o.Raw("REPLACE INTO group_label_mapping(group_name, label_id, valid) VALUES(?,?,?)", role, labelId, true).Exec()
}


func (this *LabelModel) AddLabelToMultiGroup(groupNameList string, newLid int64) {
	o := orm.NewOrm()
	var mappings []GroupLabelMapping
	lid := strconv.FormatInt(newLid, 10)

	for _, groupName := range strings.Split(groupNameList, ",") {
		mapping := GroupLabelMapping{
			GroupName: groupName,
			LabelId:   lid,
			Valid:     true,
		}
		mappings = append(mappings, mapping)
	}
	_, err2 := o.InsertMulti(100, mappings)
	if err2 != nil {
		panic(err2)
	}
}
