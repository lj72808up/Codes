package models

// 普通查询提交json
type FrontEndQueryJson struct {
	Dimensions []string                  `json:"dimensions"`
	Filters    []FrontEndQueryFilterJson `json:"filters"`
	Outputs    []string                  `json:"outputs"`
	Date       []string                  `json:"date"`
	Name       string                    `json:"name"`
	OrderCols  []string                  `json:"orderColumns"`
	LimitNum   string                    `json:"limitNum"`
}

type FrontEndQueryFilterJson struct {
	Name  string   `json:"name"`
	Value []string `json:"values"`
	Type  string   `josn:"type"`
}
