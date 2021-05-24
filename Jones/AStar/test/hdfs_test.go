package test

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/colinmarc/hdfs"
	"testing"
)

func TestHDFS(t *testing.T) {
	client, _ := hdfs.New("rsync.master01.saturn.hadoop.js.ted:8020")
	file, err := client.Open("/user/adtd_platform/search/odin/qs/20200819/326852679f1c9284dd7f689066aa5606")

	if err != nil {
		panic(err)
	}
	buf := make([]byte, 59)
	file.ReadAt(buf, 48847)

	fmt.Println(string(buf))

}

func TestBase64(t *testing.T) {
	data := "哈哈abc123!?$*&()'-=@~呵呵"
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)

	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	fmt.Println(string(sDec))
}
