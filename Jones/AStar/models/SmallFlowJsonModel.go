package models

// 小流量查询json
type SmallFlowQueryJson struct {
	TaskId    string   `json:"task_id"`
	DateRange []string `json:"date_range"`
	IsMobilQq bool     `json:"is_mobil_qq"`
	IsSame    bool     `json:"is_same"`
	Type      int      `json:"type"`
	CustomId  []string `json:"custom_id"`
	IsAvg     bool     `json:"is_avg"`
	IsCustom  bool     `json:"is_custom"`
	Query     string   `json:"query"`
	Download  string   `json:"download"`
}

// 小流量任务提交JSON
type FlowDef struct {
	Pos   string `json:"pos"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type ExtendReserver struct {
	Rshift string `json:"rshift"`
	And    string `json:"and"`
	Value  string `json:"value"`
}

type StyleReserver struct {
	Pos    string `json:"pos"`
	Rshift string `json:"rshift"`
	And    string `json:"and"`
	Value  string `json:"value"`
}

type ClickFlag struct {
	FValue string `json:"f_value"`
	Rshift string `json:"rshift"`
	Value  string `json:"value"`
}

func (this *ExtendReserver) IsEmpty() bool {
	if this.Value == "" && this.And == "" && this.Rshift == "" {
		return true
	}
	return false
}

func (this *StyleReserver) IsEmpty() bool {
	if this.Pos == "" && this.Value == "" && this.And == "" && this.Rshift == "" {
		return true
	}
	return false
}

func (this *ClickFlag) IsEmpty() bool {
	if this.FValue == "" && this.Value == "" && this.Rshift == "" {
		return true
	}
	return false
}
