package indicator

import "github.com/astaxie/beego/orm"

type Qualifier struct {
	Qid           int64  `orm:"pk;auto"`
	QualifierType string `orm:"size(255)"`
	CName         string `orm:"size(255);unique"`
	EName         string `orm:"size(255);unique"`
}

func (this *Qualifier) GetTimeQualifier() []Qualifier {
	qualifiers := []Qualifier{}
	sql := "SELECT * FROM qualifier WHERE qualifier_type = '时间周期'"
	o := orm.NewOrm()
	_, _ = o.Raw(sql).QueryRows(&qualifiers)
	return qualifiers
}

func (this *Qualifier) Add() error {
	o := orm.NewOrm()
	_, err := o.Insert(this)
	return err
}
