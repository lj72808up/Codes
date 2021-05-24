package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
)

var Logger *logs.BeeLogger

type VueResponse struct {
	Res  bool
	Info string
}

type SparkServerResponse struct {
	Flag bool   `json:"flag"`
	Msg  string `json:"msg"`
}

func (this VueResponse) String() string {
	jsonStr, _ := json.Marshal(this)
	return string(jsonStr)
}

func MakeUniqSet(items ...string) []string {
	maps := make(map[string]string)
	for _, r := range items {
		maps[r] = ""
	}
	var roleSet []string // 去重后的集合
	for k, _ := range maps {
		roleSet = append(roleSet, k)
	}
	return roleSet
}

func GenNameSpace(root string, child string) string {
	return fmt.Sprintf("%s/%s", root, child)
}

func HasValue(arr []string, value string) bool {
	var flag = false
	for _, v := range arr {
		if v == value {
			flag = true
			break
		}
	}
	return flag
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// camel string, xx_yy to XxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	flag, num := true, len(s)-1
	for i := 0; i <= num; i++ {
		d := s[i]
		if d == '_' {
			flag = true
			continue
		} else if flag {
			if d >= 'a' && d <= 'z' {
				d = d - 32
			}
			flag = false
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 执行shell命令
func ExecShell(commandName string, params []string) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(commandName, params...)
	newEnv := append(os.Environ())
	cmd.Env = newEnv
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("cmd.Start() failed with '%s'\n", err)
	}
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		fmt.Printf("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Printf("std out:%s\nstderr:%s\n", outStr, errStr)
}

func HasBom(head string) bool {
	if head[0] == 239 && head[1] == 187 && head[2] == 191 {
		return true
	} else {
		return false
	}
}

func ExecHiveCmd(cmd string) {
	hiveCommand := "/usr/lib/hive/apache-hive-1.2.1-bin/bin/beeline"
	param1 := `-n`
	param2 := `adtd_platform`
	param3 := `-p`
	param4 := `adtd_platform`
	param5 := `-u`
	param6 := `'jdbc:hive2://rsync.master09.saturn.hadoop.js.ted:10005/default;principal=hive/hive@SATURN.HADOOP.COM'` // 加上引号

	ExecShell(hiveCommand, []string{param1, param2, param3, param4, param5, param6, "-e", cmd})
}

func LoadQueryWords(hdfsDir string, localFile string) {

	ExecShell("hdfs", []string{"dfs", "-mkdir", "-p", hdfsDir})
	ExecShell("hdfs", []string{"dfs", "-put", "-f", localFile, hdfsDir})
	ExecHiveCmd(fmt.Sprintf(`"Msck repair table default.query_word"`))

}

func MakeJdbcUrl(host string, port string, database string, encode string) string {
	url := fmt.Sprintf("jdbc:mysql://%s:%s/%s?useSSL=false", host, port, database)
	if encode != "" {
		url = fmt.Sprintf("%s&useUnicode=true&characterEncoding=%s", url, encode)
	}
	return url
}

func MakeClickHouseUrl(host string, port string, database string) string {
	// 8123:http, 9000:tcp
	url := fmt.Sprintf("jdbc:clickhouse://%s:%d/%s", host, 8123, database)
	return url
}

func TransTableIdentifier(tableIdentifier string) (string, string) {
	db := "datacenter"
	tb := tableIdentifier

	splits := strings.Split(tableIdentifier, ".")
	if len(splits) > 1 {
		db = splits[0]
		tb = splits[1]
	}
	return db, tb
}

// 产生形如 (?,?,?) 的通配符序列
func MakeRolePlaceHolder(roles []string) string {
	var roleParams []string
	for range roles {
		roleParams = append(roleParams, "?")
	}
	return fmt.Sprintf("(%s)", strings.Join(roleParams, ","))
}

func IsChinese(str string) bool {
	flag := false
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			flag = true
			break
		}
	}
	return flag
}

func MakeString(arr []string, joiner string) string {
	length := len(arr)
	if length == 0 {
		return ""
	} else {
		end := length - 1
		res := arr[0]
		for i := 1; i <= end; i++ {
			res = res + joiner + arr[i]
		}
		return res
	}
}

func GetServerDBInfo() (string, string, string, string) {
	url := beego.AppConfig.String("datasourcename")
	r, _ := regexp.Compile("^(.*):(.*)@tcp\\((.*)\\)/(.*)\\?(.*)$")
	params := r.FindStringSubmatch(url)

	user := params[1]
	passwd := params[2]
	hostAndPort := params[3]
	dataBase := params[4]

	return user,passwd,hostAndPort,dataBase
}
