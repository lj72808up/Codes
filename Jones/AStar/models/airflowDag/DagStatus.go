package airflowDag

type DagStatus int

const (
	NEW        DagStatus = 1 // 刚刚创建
	ESTABLISH  DagStatus = 2 // 创建完毕
	BUILDFAIL DagStatus = 3 // 创建失败
	REMOVING   DagStatus = 8 // 删除中
	REMOVED    DagStatus = 9 // 已删除
)
