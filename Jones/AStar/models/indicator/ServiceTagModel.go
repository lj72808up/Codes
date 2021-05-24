package indicator

import (
	"AStar/utils"
	"github.com/astaxie/beego/orm"
	"time"
)

// 业务标签
type ServiceTag struct {
	Tid         int64  `orm:"pk;auto" field:"id,hidden"`
	TagName     string `orm:"size(255)" field:"业务标签名,hidden"`
	Description string `orm:"type(text)" field:"描述,hidden"`
	CreateTime  string `orm:"size(255)" field:"创建时间,hidden"`
	ModifyTime  string `orm:"size(255)" field:"最后修改时间,hidden"`
}

func (this *ServiceTag) AddOne(name string, desc string) error {
	o := orm.NewOrm()
	now := utils.TimeFormat(time.Now())
	serviceTag := ServiceTag{
		TagName:     name,
		Description: desc,
		CreateTime:  now,
		ModifyTime:  now,
	}
	_, err := o.Insert(&serviceTag)
	return err
}

func (this *ServiceTag) UpdateOne(tid int64, name string, desc string) error {
	o := orm.NewOrm()
	now := utils.TimeFormat(time.Now())
	sql := "UPDATE service_tag SET tag_name = ? , description = ? , modify_time = ? WHERE tid = ?"
	_, err := o.Raw(sql, name, desc, now, tid).Exec()
	return err
}

func (this *ServiceTag) GetAll() utils.TableData {
	o := orm.NewOrm()
	var tags []ServiceTag
	_, _ = o.Raw("SELECT * FROM service_tag").QueryRows(&tags)
	return this.transTableView(tags)
}

func (this *ServiceTag) transTableView(tags []ServiceTag) utils.TableData {
	var rows []map[string]interface{}
	for _, tag := range tags {
		tagRow := utils.GetSnakeStrMapOfStruct(tag)
		rows = append(rows, tagRow)
	}
	table := utils.TableData{Data: rows}
	table.GetSnakeTitles(ServiceTag{})
	return table
}

type ServiceTagView struct {
	Tid         int64
	Name        string
	Description string
}

func (this *ServiceTag) UpdateBatch(views []ServiceTagView) {
	for _, view := range views {
		if err := this.UpdateOne(view.Tid, view.Name, view.Description); err != nil {
			panic(err)
		}
	}
}
