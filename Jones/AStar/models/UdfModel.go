package models

import (
	"AStar/utils"
	"github.com/astaxie/beego/orm"
)

type UdfManager struct {
	FuncId      int    `orm:"pk;auto" vxe:"tid,hidden" json:"funcId"`
	FuncName    string `orm:"size(255)" vxe:"函数名" json:"funcName"`
	Description string `orm:"type(text);null" vxe:"功能描述" json:"description"`
	Parameters  string `orm:"type(text);null" vxe:"参数类型" json:"parameters"`
	Example     string `orm:"type(text);null" vxe:"sql样例" json:"example"`
	Known       bool   `orm:"default(false)" vxe:"已配置描述,hidden" json:"known"` // sparkServer 启动后, 对比数据库中的方法, 比较得出新加入的方法. 添加描述后才变成已知方法
	Builtin     bool   `orm:"default(false)" vxe:"builtIn,drop"`              // 是否是内置函数
}

func (this *UdfManager) AddDescription(funcId int, desc string) {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE udf_manager SET description = ? AND known = true WHERE FuncId = ?", desc, funcId).Exec()
	if err != nil {
		Logger.Error("%s", err.Error())
		panic(err.Error())
	}
}

func (this *UdfManager) MarkKnown(funcIds []int) {
	o := orm.NewOrm()
	for funcId, _ := range funcIds {
		_, err := o.Raw("UPDATE udf_manager SET known = true WHERE FuncId = ?", funcId).Exec()
		if err != nil {
			Logger.Error("%s", err.Error())
			panic(err.Error())
		}
	}
}

func (this *UdfManager) GetUdfList(known bool) utils.VxeData {
	o := orm.NewOrm()
	var funcList []UdfManager
	_, _ = o.Raw("SELECT * FROM udf_manager WHERE known = ? and builtin = false", known).QueryRows(&funcList)
	return this.transTableView(funcList)
}

func (this *UdfManager) transTableView(udfData []UdfManager) utils.VxeData {
	var vxe utils.VxeData
	editableCols := map[string]bool{
		"description": true,
	}
	jsonTagMapping := vxe.MakeColumns(UdfManager{}, editableCols)

	var dataList []map[string]interface{}
	for i, _ := range udfData {
		dataList = append(dataList, vxe.MakeData(udfData[i], &jsonTagMapping))
	}

	vxe.DataList = dataList
	return vxe
}

func (this *UdfManager) UpdateDescription(udfs []UdfManager) error {
	o := orm.NewOrm()
	for _, udf := range udfs {
		sql := "UPDATE udf_manager SET description=? , known = true, parameters = ?, example = ? WHERE func_name=? "
		_, err := o.Raw(sql, udf.Description, udf.Parameters, udf.Example, udf.FuncName).Exec()
		if err != nil {
			Logger.Error(err.Error())
			return err
		}
	}
	return nil
}
