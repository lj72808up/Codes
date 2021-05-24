package SmallFlow

import (
	"AStar/models"
	"strconv"
	"fmt"
)

func CalcABTestResult(resA models.SmallFlowRes, resB models.SmallFlowRes) models.SmallFlowRes{
	return models.SmallFlowRes{
		Id:         0,
		Dt:         resA.Dt,
		Tag:        "A/B-1",
		Pw:         resA.Pw,
		SumPvPage:  cal(resA.SumPvPage, resB.SumPvPage),
		SumPvAd:    cal(resA.SumPvAd, resB.SumPvAd),
		SumCharge:  cal(resA.SumCharge, resB.SumCharge),
		SumClick:   cal(resA.SumClick, resB.SumClick),
		AvgCtrPage: cal(resA.AvgCtrPage, resB.AvgCtrPage),
		AvgCtrAd:   cal(resA.AvgCtrAd, resB.AvgCtrAd),
		AvgAcp:     cal(resA.AvgAcp, resB.AvgAcp),
		AvgRpmPage: cal(resA.AvgRpmPage, resB.AvgRpmPage),
		CustomId:   resA.CustomId,
	}
}

func cal(a string, b string) string{

	var af, _ = strconv.ParseFloat(a, 64)

	var bf, _ = strconv.ParseFloat(b, 64)

	var res = (af/bf -1)*100.0
	return fmt.Sprintf("%.3f", res) + "%"
}


