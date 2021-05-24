package test

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {

	url := "http://inner-redash.adtech.sogou/api/queries"
	query := `"--我是注释  \n SELECT t.server_time AS dj ,t.user_ip AS ip ,t.accountid AS khid FROM charge.bill_ie_log_{{date}} t limit {{cnt}}"`

	param := fmt.Sprintf(`
{
    "query":%s,
    "data_source_id":35,
    "name":"test_bill_query",
    "options":{
        "parameters":[
            {
                "name":"date",
                "title":"日期",
                "global":false,
                "value":20200918,
                "type":"number",
                "locals":[

                ]
            },
            {
				"name":"cnt",
                "title":"条数",
                "global":false,
                "value":120,
                "type":"number",
                "locals":[

                ]
			}
        ]
    }
}
`, query)

	fmt.Println(param)
	req := httplib.Post(url).SetTimeout(time.Second*5, time.Second*5).Body([]byte(param))
	bytes, err := req.Bytes()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(bytes))
}


