package utils

import (
	"time"
)

//GetDiffDays 获取两天之间的天数
func GetDiffDays(start string, end string) (int, time.Time, time.Time) {
	dateStart, _ := time.Parse("20060102", start)
	dateEnd, _ := time.Parse("20060102", end)
	diffDays := int(dateEnd.Sub(dateStart).Hours() / 24)
	return diffDays, dateStart, dateEnd
}

// 包含首尾的日期
func GetDateRange(start string, end string) []string {
	var dateRange []string
	days, dtStart, _ := GetDiffDays(start, end)
	for i := 0; i <= days; i++ {
		dateRange = append(dateRange, DateFormatWithoutLine(dtStart.AddDate(0, 0, i)))
	}
	return dateRange
}

//DateFormat 日期格式化
func DateFormat(time time.Time) string {
	return time.Format("2006-01-02")
}

//DateFormatWithoutLine 日期格式化
func DateFormatWithoutLine(time time.Time) string {
	return time.Format("20060102")
}

//TimeFormat 日期格式化
func TimeFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

// 获取日期间隔中的全部日期字符串数组
func GetDaysAmong(start string, end string) []string {
	var diffDays, dateStart, dateEnd = GetDiffDays(start, end)
	var DaysAmongArr = []string{}
	var oneDay, _ = time.ParseDuration("24h")
	for i := 0; i < diffDays; i++ {
		DaysAmongArr = append(DaysAmongArr, DateFormatWithoutLine(dateStart))
		dateStart = dateStart.Add(oneDay)
	}
	DaysAmongArr = append(DaysAmongArr, DateFormatWithoutLine(dateEnd))
	return DaysAmongArr
}

func TransDateByStandard(timeExpression string, fromFormat string, toFormat string) string {
	internal, _ := time.Parse(fromFormat, timeExpression)
	return internal.Format(toFormat)
}
