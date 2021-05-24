package models

import (
	"github.com/astaxie/beego/httplib"
	"time"
	"github.com/tidwall/gjson"
	"AStar/utils"
)

type KylinInfoModel struct {
	KylinURL    string
	KylinUser   string
	KylinPasswd string
}

// USAGE : if ! models.KylinInfo.CheckKylin(fact.Kylin, dateFilter.Values[0], dateFilter.Values[1]
func (this *KylinInfoModel) CheckKylin(cube string, dateStart string, dateEnd string) bool {
	req := httplib.Get(this.KylinURL + "/cubes/" + cube).SetTimeout(3*time.Second, 3*time.Second).SetBasicAuth(this.KylinUser, this.KylinPasswd)
	bytes, error := req.Bytes()
	if error != nil {
		return false
	}

	seg := gjson.GetBytes(bytes, "segments")
	//fmt.Println("cube", cube)
	//fmt.Println("segments", seg)

	dateMap := make(map[string]int)
	for _, arr := range seg.Array() {
		date := utils.DateFormatWithoutLine(time.Unix(gjson.Get(arr.Raw, "date_range_start").Int()/1000, 0))
		status := gjson.Get(arr.Raw, "status").Str
		if status == "READY" {
			dateMap[date] = 1
		}
	}

	subDays, startDate, _ := utils.GetDiffDays(dateStart, dateEnd)

	for day := 0; day <= subDays; day ++ {
		dateStr := utils.DateFormatWithoutLine(startDate.AddDate(0, 0, day))
		_, isKey := dateMap[dateStr]
		if !isKey {
			return false
		}
	}
	return true
}
