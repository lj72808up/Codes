package Sql

import (
	"AStar/utils"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
	"strings"
)

type SemanticParser struct{}

func (this *SemanticParser) GetFields(sql string) []string {
	p := parser.New()
	sql = strings.Replace(sql,";","",-1)
	stmtNodes, _, err := p.Parse(sql, "", "")

	if err != nil {
		utils.Logger.Error("parse error:\n%v\n%s", err, sql)
		return []string{}
	}
	v := SelectFieldVisitor{}
	for _, stmtNode := range stmtNodes {
		stmtNode.Accept(&v)
	}
	return v.ColumnName
}

type SelectFieldVisitor struct {
	ColumnName []string
}

func (v *SelectFieldVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	root := in.(*ast.SelectStmt)
	for _, f := range root.Fields.Fields {
		originalName := f.Text()
		asName := f.AsName.L
		if originalName != "" {
			v.ColumnName = append(v.ColumnName, originalName)
		} else {
			v.ColumnName = append(v.ColumnName, asName)
		}
	}
	return in, true
}

func (v *SelectFieldVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}
