package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

type VxeColumn struct {
	Field      string      `json:"field"`
	Title      string      `json:"title"`
	EditRender *EditRender `json:"editRender,omitempty"`
	Visible    bool        `json:"visible"`
}

type VxeData struct {
	TableColumn []VxeColumn              `json:"tableColumn"`
	DataList    []map[string]interface{} `json:"dataList"`
}

type EditRender struct {
	Name string `json:"name"`
}

const VXE_TAG string = "vxe"

func (this *VxeData) ToString() string {
	res, _ := json.Marshal(this)
	return string(res)
}

func (this *VxeData) getJsonKeyName() {

}

func (this *VxeData) MakeColumns(emptyType interface{}, editableCols map[string]bool) map[string]string {
	var cols []VxeColumn
	fieldNames := GetFieldName(emptyType)
	vxeFieldTagMapping := GetFieldTagName(emptyType, VXE_TAG)
	jsonTagMapping := GetFieldTagName(emptyType, "json")

	for _, fieldName := range fieldNames {
		vxeTag := vxeFieldTagMapping[fieldName]
		vxeFieldName := strings.Split(jsonTagMapping[fieldName], ",")[0]
		tags := strings.Split(vxeTag, ",")

		canEdit := editableCols[vxeFieldName]
		var col VxeColumn
		if canEdit {
			col = VxeColumn{
				Field:      vxeFieldName,
				Title:      tags[0],
				EditRender: &EditRender{Name: "input"},
				Visible:    true,
			}
		} else {
			col = VxeColumn{
				Field:      vxeFieldName,
				Title:      tags[0],
				EditRender: nil,
				Visible:    true,
			}
			//col.EditRender = nil
		}

		if len(tags) > 1 {
			if strings.Contains(vxeTag, "drop") { // tag中包含drop
				continue
			}
			if strings.Contains(vxeTag, "hidden") { // 是否可见
				col.Visible = false
			}
		}
		cols = append(cols, col)
	}
	this.TableColumn = cols
	return jsonTagMapping
}

func (this *VxeData) MakeData(obj interface{}, jsonFieldMapping *map[string]string) map[string]interface{} {
	var resMap = make(map[string]interface{}) // k:属性名,v:值
	var sType = reflect.TypeOf(obj)
	var sValue = reflect.ValueOf(obj)
	var numsOfFields = sType.NumField()
	for i := 0; i < numsOfFields; i++ {
		//snakeName := utils.SnakeString(sType.Field(i).Name)
		snakeName := strings.Split((*jsonFieldMapping)[sType.Field(i).Name], ",")[0]
		tagVal := sType.Field(i).Tag.Get(VXE_TAG)
		if strings.Contains(tagVal, "drop") {
			continue
		}
		resMap[snakeName] = sValue.Field(i).Interface()
	}
	return resMap
}
