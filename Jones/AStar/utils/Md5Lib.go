package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/axgle/mahonia"
)

func GetMd5ForQuery(sqlStr string, session string, ts string) string {
	Md5Inst := md5.New()
	Md5Inst.Write([]byte(sqlStr + session + ts))
	return hex.EncodeToString(Md5Inst.Sum(nil))
}

func GetMd5(str string) string {
	Md5Inst := md5.New()
	Md5Inst.Write([]byte(str))
	return hex.EncodeToString(Md5Inst.Sum(nil))
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}