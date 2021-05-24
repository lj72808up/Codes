package models

// 给前端动态获取选框里的选项的模板结构体
type FrontFilterModel struct {
	Id          int
	Date        string
	FiltersJson string `orm:"type(text)"`
	PcJson string `orm:"type(text)"`
	GalaxyJson string `orm:"type(text)"`
}
