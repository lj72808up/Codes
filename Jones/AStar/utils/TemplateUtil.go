package utils

import (
	"bytes"
	"fmt"
	"github.com/flosch/pongo2"
	"text/template"
	"time"
)

type SqlUtil struct{}
type IntervalVar struct {
	start string
	end   string
}

func (this *SqlUtil) SubstituteGo(sql string, param map[string]string) string {
	var doc bytes.Buffer
	tmplName := fmt.Sprint(time.Now().Unix())
	t := template.New(tmplName)
	t, _ = t.Parse(sql)
	err := t.Execute(&doc, param)
	if err != nil {
		panic(err.Error())
	}
	runSql := doc.String()
	return runSql
}

func (this *SqlUtil) SubstituteJinja(tmplateSql string, param map[string]interface{}) string {
	tpl, _ := pongo2.FromString(tmplateSql)
	Logger.Info("模板: %v",tmplateSql)
	Logger.Info("参数: %v",param)
	out, err := tpl.Execute(param)
	if err != nil {
		panic(err)
	}
	return out
}

func (this *SqlUtil) PaddingDagInternalVarFromUserParam(userParam map[string]interface{}) map[string]interface{} {
	intervalVars := []string{
		"Last7days", "LastWeek", "ThisWeek", "LastMonth", "ThisMonth",
	}
	momentVars := []string{
		"yyyyMMdd", "Yesterday", "DayBeforeYesterday","Today",
	}
	param := make(map[string]interface{})
	for k, v := range userParam {
		param[k] = v
	}

	dayPlaceHolder := "20060101"
	hourPlaceHolder := "1"

	// 填充小时变量
	param["HH"] = hourPlaceHolder
	// 填充间隔变量
	for _, variable := range intervalVars {
		param[variable] = IntervalVar{
			start: dayPlaceHolder,
			end:   dayPlaceHolder,
		}
	}

	// 填充瞬时变量
	for _, variable := range momentVars {
		param[variable] = dayPlaceHolder
	}

	return param
}
