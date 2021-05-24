package models

import (
	"github.com/astaxie/beego/orm"
)

// 小流量任务的struct
type
	PlatformDoc struct {
		Id           int
		LastDate     string
		Submitter    string
		Content      string `orm:"type(text)"`
		Html 		 string `orm:"type(text)"`
	}

type
	SubmitDocJson struct {
		Submitter    string
		Content      string
		Html 		 string
	}

func (this PlatformDoc) UpdateContent(lastDate string, submitter string, content string, html string) {
	this.LastDate = lastDate
	this.Submitter = submitter
	this.Content = content
	this.Html = html
	_, err := orm.NewOrm().Update(&this,"LastDate","Submitter","Content","Html")
	if err != nil {
		Logger.Error(err.Error())
	}
}


