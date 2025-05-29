package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itxiaoma0610/sharddoc/database/engine"
	"github.com/itxiaoma0610/sharddoc/database/parser/ast"
)

type antlrParser struct {
	*ast.BaseSQLParserVisitor
}

func (a *antlrParser) Visit(tree antlr.ParseTree) interface{} {
	return a.VisitSql(tree.(*ast.SqlContext))
}
func (a *antlrParser) VisitSql(ctx *ast.SqlContext) interface{} {
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		switch c := child.(type) {
		case *ast.CreateTableStatementContext:
			return a.VisitCreateTableStatement(c)
		case *ast.InsertTableStatementContext:
			return a.VisitInsertTableStatement(c)
		}
	}
	return nil
}

func (a *antlrParser) VisitCreateTableStatement(ctx *ast.CreateTableStatementContext) interface{} {
	// CreatTableStatement return
	t := new(engine.TableDef)
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		switch c := child.(type) {
		case *ast.TableNameContext:
			t.Name = a.VisitTableName(c).(string)
		case *ast.ColumnDefinitionsContext:
			defs := a.VisitColumnDefinitions(c).(*engine.TableDef)
			t.Cols = defs.Cols
			t.Types = defs.Types
		case *ast.ColumnNameContext:
			t.Indexes = append(t.Indexes, []string{c.GetText()})
		case *ast.IndexDefinitionsContext:
			indexs := a.VisitIndexDefinitions(c).([][]string)
			if len(t.Indexes) == 0 {
				t.Indexes[0] = []string{}
			}
			t.Indexes = append(t.Indexes, indexs...)
		}
	}
	return t
}
func (a *antlrParser) VisitTableName(ctx *ast.TableNameContext) interface{} {
	return ctx.GetText()
}
func (a *antlrParser) VisitColumnDefinitions(ctx *ast.ColumnDefinitionsContext) interface{} {
	// return col []string{}, type []uint32{}
	t := &engine.TableDef{}
	defs := ctx.AllColumnDefinition()
	for _, def := range defs {
		t.Cols = append(t.Cols, def.GetStart().GetText())
		ctype := def.GetStop().GetText()
		switch ctype {
		case "INT64":
			t.Types = append(t.Types, engine.TYPE_INT64)
		case "BYTES":
			t.Types = append(t.Types, engine.TYPE_BYTES)
		default:
			return fmt.Errorf("unreachable type")
		}
	}
	return t
}

func (a *antlrParser) VisitIndexDefinitions(ctx *ast.IndexDefinitionsContext) interface{} {
	indexs := [][]string{}
	for _, indexdef := range ctx.AllIndexDefinition() {
		index := []string{}
		for _, iname := range indexdef.AllColumnName() {
			index = append(index, iname.GetText())
		}
		indexs = append(indexs, index)
	}
	return indexs
}
func (a *antlrParser) VisitInsertTableStatement(ctx *ast.InsertTableStatementContext) interface{} {
	// 	rec := Record{}
	// rec.AddInt64("ki1", int64(i)).AddStr("ks2", []byte("hello"))
	// rec.AddStr("s1", []byte("world")).AddInt64("i2", int64(i/2))
	// added := r.add("tbl_test", rec)
	// assert(added)
	rec := engine.Record{}
	for i := 0; i < len(ctx.ColumnInsertValues().AllColumnValue()); i++ {
		col := ctx.ColumnInsertNames().AllColumnName()[i].GetText()
		if ctx.ColumnInsertValues().AllColumnValue()[i].INTEGER() == nil {
			// value types bytes
			// val = engine.Value{Type: engine.TYPE_BYTES, Str: []byte(v)}
			v := ctx.ColumnInsertValues().AllColumnValue()[i].GetText()
			rec.AddStr(col, []byte(v))
		} else {
			v := ctx.ColumnInsertValues().AllColumnValue()[i].GetText()
			i, _ := strconv.ParseInt(v, 10, 64)
			rec.AddInt64(col, i)
			// val = engine.Value{Type: engine.TYPE_INT64, I64: -1}
		}
	}

	return rec
}

func TestCreateTable(t *testing.T) {

	creatTable := `
	CREATE TABLE users (
        id INT64,
        name BYTES,
        age INT64,
		height INT64,
		PRIMARY KEY (id),
        INDEX (age, height)
    );
	`
	input := antlr.NewInputStream(creatTable)
	lexer := ast.NewSQLLexer(input)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	p := ast.NewSQLParser(tokenStream)
	p.AddErrorListener(antlr.NewDefaultErrorListener())

	tree := p.Sql()
	if tree == nil {
		t.Errorf("parse tree is nil, check input or parser configuration")
	}
	v := new(antlrParser)
	tableDef := v.Visit(tree).(*engine.TableDef)
	fmt.Println(tableDef)
}
func TestInsertTable(t *testing.T) {
	InsertSql := `
	INSERT INTO users (id, name, age, height) 
	VALUES (1, 'Alice', 30, 170);
	`
	input := antlr.NewInputStream(InsertSql)
	lexer := ast.NewSQLLexer(input)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	p := ast.NewSQLParser(tokenStream)
	p.AddErrorListener(antlr.NewDefaultErrorListener())

	tree := p.Sql()
	if tree == nil {
		t.Errorf("parse tree is nil, check input or parser configuration")
	}
	v := new(antlrParser)
	tableDef := v.Visit(tree)
	fmt.Println(tableDef)
}
