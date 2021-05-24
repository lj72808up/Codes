package airflowDag

type DagListView struct {
	Code int         `json:"code"`
	Data []DagByRest `json:"data"`
	Msg  string      `json:"msg"`
}

type DagByRest struct {
	DagId            string `json:"dag_id"`
	IsPaused         string `json:"is_paused"`
	Owner            string `json:"owner"`
	ScheduleInterval string `json:"schedule_interval"`
}
