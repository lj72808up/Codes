package test

import (
	"AStar/utils/Sql"
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"testing"

	_ "github.com/pingcap/parser/test_driver"
)
type visitor struct{}

func (v *visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	fmt.Printf("%T\n", in)
	root := in.(*ast.SelectStmt)
	for _,f := range root.Fields.Fields{
		fmt.Println(f.Text()=="")
		fmt.Println(f.AsName)
		fmt.Println("xxx")
	}
	return in, true
}

func (v *visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func TestEx(t *testing.T) {
	p := parser.New()

	sql := `SELECT dws_wsbidding_jonesad_di.dt AS RQ, name, dws_wsbidding_jonesad_di.account_id AS KHID FROM adtl_biz.dws_wsbidding_jonesad_di`
	stmtNodes, _, err := p.Parse(sql, "", "")

	if err != nil {
		fmt.Printf("parse error:\n%v\n%s", err, sql)
		return
	}
	for _, stmtNode := range stmtNodes {
		v := visitor{}
		stmtNode.Accept(&v)
	}
}

func TestEx2 (t *testing.T){
	parser := Sql.SemanticParser{}
	sql := `SELECT dws_wsbidding_jonesad_di.dt AS RQ, name, dws_wsbidding_jonesad_di.account_id AS KHID FROM adtl_biz.dws_wsbidding_jonesad_di`
	fmt.Println(parser.GetFields(sql))
}