package utils

import (
	"strings"
)

type TableTitle struct {
	Prop     string `json:"prop"`
	Label    string `json:"label"`
	Sortable bool   `json:"sortable"`
	Hidden   bool   `json:"hidden"`
}

type TableData struct {
	Cols []TableTitle        `json:"cols"`
	Data []map[string]interface{} `json:"data"`
}

const tagName = "field"

func (this * TableData)GetSnakeTitles(structValue interface{}) {
	var res []TableTitle
	sb := GetFieldName(structValue)
	st := GetFieldTagName(structValue, tagName)
	for _, title := range sb {
		var diyField = strings.Split(st[title], ",")
		snakeTitle := SnakeString(title)
		if len(diyField) > 1 {
			var showType = diyField[1]
			if showType == "hidden" {
				res = append(res, TableTitle{snakeTitle, diyField[0], true, true})
			} else if showType == "drop" {
				continue
			} else  {
				res = append(res, TableTitle{snakeTitle, diyField[0], true, false})
			}

		} else {
			res = append(res, TableTitle{snakeTitle, diyField[0], true, false})
		}
	}
	this.Cols = res
}

func (this * TableData)GetSnakeTitlesWithHidden(structValue interface{}, showMap map[string]bool) {
	var res []TableTitle
	sb := GetFieldName(structValue)
	st := GetFieldTagName(structValue, tagName)
	for _, title := range sb {
		var diyField = strings.Split(st[title], ",")
		snakeTitle := SnakeString(title)
		if len(diyField) > 1 {
			var showType = diyField[1]
			if showType == "hidden" {
				res = append(res, TableTitle{snakeTitle, diyField[0], true, !showMap[snakeTitle]})
			} else if showType == "drop" {
				continue
			} else {
				res = append(res, TableTitle{snakeTitle, diyField[0], true, false})
			}

		} else {
			res = append(res, TableTitle{snakeTitle, diyField[0], true, false})
		}
	}
	this.Cols = res
}

func (this * TableData)GetSnakeTitlesWithHiddenForMuti(structValue interface{}, showMap map[string]bool, idx int) {
	var res []TableTitle
	sb := GetFieldName(structValue)
	st := GetFieldTagName(structValue, tagName)
	for _, title := range sb {
		var diyField = strings.Split(st[title], ",")
		var showTitle = strings.Split(diyField[0], ";")
		snakeTitle := SnakeString(title)
		if len(diyField) > 1 {
			var showType = diyField[1]
			if showType == "hidden" {
				res = append(res, TableTitle{snakeTitle, showTitle[0], true, !showMap[snakeTitle]})
			} else if showType == "drop" {
				continue
			} else {
				res = append(res, TableTitle{snakeTitle, showTitle[idx], true, false})
			}

		} else {
			res = append(res, TableTitle{snakeTitle, showTitle[idx], true, false})
		}
	}
	this.Cols = res
}
