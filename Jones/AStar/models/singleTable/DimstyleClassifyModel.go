package singleTable

import (
	"github.com/astaxie/beego/orm"
)

type DimStyleClassify struct {
	Id              int64
	AdStyleId       int64
	ThirdStyleName  string
	SecondStyleName string
	FirstStyleName  string
}

func GetDimStyleClassifySet() []DimStyleClassify {
	o1 := orm.NewOrm()
	_ = o1.Using("DimStyleClassify")
	res := []DimStyleClassify{}
	_, _ = o1.Raw("SELECT * FROM dim_style_classify").QueryRows(&res)
	return res
}

func AddDimStyleClassify(dim DimStyleClassify) {
	o1 := orm.NewOrm()
	_ = o1.Using("DimStyleClassify")
	_, _ = o1.Raw(`
INSERT INTO dim_style_classify(ad_style_id,third_style_name,second_style_name,first_style_name)
VALUES (?,?,?,?);
`,dim.AdStyleId,dim.ThirdStyleName,dim.SecondStyleName,dim.FirstStyleName).Exec()
}

func UpdateDimStyleClassifySet(dim DimStyleClassify) {
	o1 := orm.NewOrm()
	_ = o1.Using("DimStyleClassify")
	_, _ = o1.Raw(`UPDATE dim_style_classify SET
ad_style_id = ?,
third_style_name = ?,
second_style_name = ?,
first_style_name = ?
WHERE id = ?
`, dim.AdStyleId, dim.ThirdStyleName, dim.SecondStyleName, dim.FirstStyleName, dim.Id).Exec()
}

func DeleteDimStyleClassifySet(id int64) {
	o1 := orm.NewOrm()
	_ = o1.Using("DimStyleClassify")
	_, _ = o1.Raw("DELETE FROM dim_style_classify WHERE id = ?",id).Exec()
}
