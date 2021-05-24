package models

var RouteMap = map[string]string{
	"DWController_QuerySQL":      "query_sql",
	"DWController_QuerySQLTask":  "query_custom",
	"DWController_Meta":          "get_meta",
	"TaskController_Limit":       "get_limit",
	"SmallFlowController_Submit": "small_flow",
}

type DBConfigModel struct {
	DriverName     string
	DataSourceName string
}

type RedisConfigModel struct {
	RedisHost string
	RedisPort string
	RedisPwd  string
}

type TriggerDagModel struct {
	TriggerPollingInterval int
	WorkerNum              int
}

type ServerInfoModel struct {
	SparkHost       string
	PaganiHost      string
	SparkPort       int
	RabbitMQ        string
	MQName          string
	RsyncPath       string
	PaganiRsyncPath string
	UploadRoot      string
}

type AirflowInfoModel struct {
	AirflowApiUrl string
}
