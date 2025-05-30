package parser

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itxiaoma0610/sharddoc/database/engine"
	"github.com/itxiaoma0610/sharddoc/database/parser/ast"
)

type SQLParser struct {
	*ast.BaseSQLParserVisitor
}

func (a *SQLParser) Visit(tree antlr.ParseTree) interface{} {
	return a.VisitSql(tree.(*ast.SqlContext))
}
func (a *SQLParser) VisitSql(ctx *ast.SqlContext) interface{} {
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		switch c := child.(type) {
		case *ast.CreateTableStatementContext:
			return a.VisitCreateTableStatement(c)
		case *ast.InsertTableStatementContext:
			return a.VisitInsertTableStatement(c)
		case *ast.SelectTableStatementContext:
			return a.VisitSelectTableStatement(c)
		}
	}
	return nil
}

func (a *SQLParser) VisitCreateTableStatement(ctx *ast.CreateTableStatementContext) interface{} {
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
func (a *SQLParser) VisitTableName(ctx *ast.TableNameContext) interface{} {
	return ctx.GetText()
}
func (a *SQLParser) VisitColumnDefinitions(ctx *ast.ColumnDefinitionsContext) interface{} {
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

func (a *SQLParser) VisitIndexDefinitions(ctx *ast.IndexDefinitionsContext) interface{} {
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
func (a *SQLParser) VisitInsertTableStatement(ctx *ast.InsertTableStatementContext) interface{} {
	// 	rec := Record{}

	rec := &engine.Record{}
	for i := 0; i < len(ctx.ColumnInsertValues().AllColumnValue()); i++ {
		col := ctx.ColumnInsertNames().AllColumnName()[i].GetText()
		if ctx.ColumnInsertValues().AllColumnValue()[i].INTEGER() == nil {
			// value types bytes
			v := ctx.ColumnInsertValues().AllColumnValue()[i].GetText()
			v = strings.Trim(v, "'")
			rec.AddStr(col, []byte(v))
		} else {
			v := ctx.ColumnInsertValues().AllColumnValue()[i].GetText()
			i, _ := strconv.ParseInt(v, 10, 64)
			rec.AddInt64(col, i)
		}
	}

	return rec

}

type CondType int

const (
	SIG_PARE          CondType = iota
	SIG_BEWTEEN                // age between 18 and 24
	PARE_LINK_PARE             // age ==18 and height > 24
	PARE_LINK_BETWEEN          //age = 18 and height between 170 and 175
	UNSUPPORTED_TYPE
)
const (
	MIN_NAME = "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
	MAX_NAME = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

func conditionType(cond ...ast.IConditionContext) CondType {
	if len(cond) == 3 {
		return UNSUPPORTED_TYPE
	}
	if len(cond) == 1 {
		if condition, ok := (cond[0]).(*ast.ConditionContext); ok {
			if condition.BetweenCondition() != nil {
				return SIG_BEWTEEN
			} else {
				return SIG_PARE
			}
		}
	}
	if condition, ok := (cond[1]).(*ast.ConditionContext); ok {
		if condition.BetweenCondition() != nil {
			return PARE_LINK_BETWEEN
		} else {
			return PARE_LINK_PARE
		}
	}
	return UNSUPPORTED_TYPE
}

var Op2Cmp = map[string]int{
	">":  engine.CMP_GT,
	">=": engine.CMP_GE,
	"<":  engine.CMP_LT,
	"<=": engine.CMP_LE,
}

func (v *SQLParser) VisitSelectTableStatement(ctx *ast.SelectTableStatementContext) interface{} {
	// todo current support select sync, 暂时只支持两个条件
	// select * from table where age > 18 and age < 40 相同键位
	// select * from table where age = 18 and height = 175
	// select * from table where age = 18 and height bewteen 170 and 175
	scan := &engine.Scanner{}
	tablename := ctx.TableName().GetText()

	condtype := conditionType(ctx.AllCondition()...)
	switch condtype {
	case SIG_PARE:
		_ = selectPare(scan, ctx)
	case SIG_BEWTEEN:
		selectBewteen(scan, ctx)
	case PARE_LINK_PARE:
		selectPareLinkPare(scan, ctx)
	case PARE_LINK_BETWEEN:
		fmt.Println("pare link between")
	case UNSUPPORTED_TYPE:
		return fmt.Errorf("unsupported type")
	}
	fmt.Println(tablename)
	fmt.Println(scan)
	return scan
}

func selectPare(scan *engine.Scanner, ctx *ast.SelectTableStatementContext) error {
	cond := ctx.AllCondition()[0].ComparisonCondition()
	op := cond.OP().GetText()
	colname := cond.ColumnName().GetText()

	recold := engine.Record{}
	if cond.ColumnValue().INTEGER() != nil {
		val, _ := strconv.ParseInt(cond.ColumnValue().INTEGER().GetText(), 10, 64)
		recold.AddInt64(colname, val)
	} else {
		recold.AddStr(colname, []byte(cond.ColumnValue().GetText()))
	}

	recnew := engine.Record{}
	switch op {
	case ">":
		// age > 18
		if cond.ColumnValue().INTEGER() != nil {
			recnew.AddInt64(colname, math.MaxInt64/2)
		} else {
			// name 最大长度32, 规定编码规则
			recnew.AddStr(colname, []byte(MAX_NAME))
		}
		scan.Cmp1 = engine.CMP_GT
		scan.Cmp2 = engine.CMP_LE
		scan.Key1 = recold
		scan.Key2 = recnew
		return nil
	case ">=":
		// age >= 18
		if cond.ColumnValue().INTEGER() != nil {
			recnew.AddInt64(colname, math.MaxInt64/2)
		} else {
			// name 最大长度32, 规定编码规则
			recnew.AddStr(colname, []byte(MAX_NAME))
		}
		scan.Cmp1 = engine.CMP_GE
		scan.Cmp2 = engine.CMP_LE
		scan.Key1 = recold
		scan.Key2 = recnew
		return nil
	case "=":
		if cond.ColumnValue().INTEGER() != nil {
			val, _ := strconv.ParseInt(cond.ColumnValue().INTEGER().GetText(), 10, 64)
			recnew.AddInt64(colname, val)
		} else {
			// name 最大长度32, 规定编码规则
			recnew.AddStr(colname, []byte(cond.ColumnValue().GetText()))
		}
		scan.Cmp1 = engine.CMP_GE
		scan.Cmp2 = engine.CMP_LE
		scan.Key1 = recold
		scan.Key2 = recnew
		return nil
	case "<":
		// -inf < age < 18
		if cond.ColumnValue().INTEGER() != nil {
			recnew.AddInt64(colname, math.MinInt64/2)
		} else {
			// name 最大长度32, 规定编码规则
			recnew.AddStr(colname, []byte(MIN_NAME))
		}
		scan.Cmp1 = engine.CMP_GE
		scan.Cmp2 = engine.CMP_LT
		scan.Key1 = recnew
		scan.Key2 = recold
		return nil
	case "<=":
		// -inf < age < 18
		if cond.ColumnValue().INTEGER() != nil {
			recnew.AddInt64(colname, math.MinInt64/2)
		} else {
			// name 最大长度32, 规定编码规则
			recnew.AddStr(colname, []byte(MIN_NAME))
		}
		scan.Cmp1 = engine.CMP_GE
		scan.Cmp2 = engine.CMP_LE
		scan.Key1 = recnew
		scan.Key2 = recold
		return nil
	}
	return nil
}
func selectBewteen(scan *engine.Scanner, ctx *ast.SelectTableStatementContext) {
	cond := ctx.AllCondition()[0].BetweenCondition()
	colname := cond.ColumnName().GetText()
	recl := engine.Record{}
	recr := engine.Record{}

	vl := cond.ColumnValue(0).GetText()
	vr := cond.ColumnValue(1).GetText()

	if cond.ColumnValue(0).INTEGER() != nil {
		vall, _ := strconv.ParseInt(vl, 10, 64)
		valr, _ := strconv.ParseInt(vr, 10, 64)
		recl.AddInt64(colname, vall)
		recr.AddInt64(colname, valr)
		return
	}
	recl.AddStr(colname, []byte(vl))
	recr.AddStr(colname, []byte(vr))
	return
}
func condInt64Data(data *ast.ColumnValueContext) int64 {
	if data.INTEGER() == nil {
		panic("data type err")
	}
	val, _ := strconv.ParseInt(data.GetText(), 10, 64)
	return val
}
func condStrData(data *ast.ColumnValueContext) string {
	if data.STRING() == nil {
		panic("data type err")
	}
	return data.GetText()
}
func selectPareLinkPare(scan *engine.Scanner, ctx *ast.SelectTableStatementContext) error {
	// 两种情况, todo 暂时规定左大右小
	// select * from table where age > 18 and age < 40 相同键位
	// select * from table where age = 18 and height = 175
	cond1 := ctx.AllCondition()[0].ComparisonCondition()
	cond2 := ctx.AllCondition()[1].ComparisonCondition()

	recold := engine.Record{}
	recnew := engine.Record{}

	if cond1.ColumnName().GetText() == cond2.ColumnName().GetText() {
		// select * from table where age > 18 and age < 40 相同键位
		scan.Cmp1 = Op2Cmp[cond1.OP().GetText()]
		scan.Cmp2 = Op2Cmp[cond2.OP().GetText()]
		colname := cond1.ColumnName().GetText()
		// 左边仅允许使用 > || >= , 右边仅允许 < || <=
		if cond1.ColumnValue().INTEGER() != nil {
			vall := condInt64Data(cond1.ColumnValue().(*ast.ColumnValueContext))
			valr := condInt64Data(cond2.ColumnValue().(*ast.ColumnValueContext))

			recold.AddInt64(colname, vall)
			recold.AddInt64(colname, valr)
			scan.Key1, scan.Key2 = recold, recnew
			return nil
		} else {
			recold.AddStr(colname, []byte(cond1.ColumnValue().GetText()))
			recnew.AddStr(colname, []byte(cond2.ColumnValue().GetText()))
			scan.Key1, scan.Key2 = recold, recnew
			return nil
		}
	} else {
		// select * from table where age = 18 and height = 175
		// 左边age仅允许为=,右边任意
		if cond1.OP().GetText() != "=" {
			return fmt.Errorf("not supported operation")
		}

		col1name := cond1.ColumnName().GetText()
		cond1Info := condData(ctx.AllCondition()[0].(*ast.ConditionContext))
		if cond1Info.datatype == engine.TYPE_INT64 {
			recold.AddInt64(col1name, cond1Info.i64)
			recnew.AddInt64(col1name, cond1Info.i64)
		} else {
			recold.AddStr(col1name, cond1Info.bytes)
			recnew.AddStr(col1name, cond1Info.bytes)
		}

		// 处理复合右边
		c2Info := condData(ctx.AllCondition()[1].(*ast.ConditionContext))

		switch c2Info.op {
		case ">", ">=":
			// height > 175
			if c2Info.datatype == engine.TYPE_INT64 {
				recold.AddInt64(c2Info.colname, c2Info.i64)
				recnew.AddInt64(c2Info.colname, math.MaxInt64/2)
			} else {
				recold.AddStr(c2Info.colname, c2Info.bytes)
				recnew.AddStr(c2Info.colname, []byte(MAX_NAME))
			}
			scan.Cmp1, scan.Cmp2 = engine.CMP_GT, engine.CMP_LE
			scan.Key1, scan.Key2 = recold, recnew

			if c2Info.op == ">=" {
				scan.Cmp1, scan.Cmp2 = engine.CMP_GE, engine.CMP_LE
			}
			return nil
		case "=":
			if c2Info.datatype == engine.TYPE_INT64 {
				recold.AddInt64(c2Info.colname, c2Info.i64)
				recnew.AddInt64(c2Info.colname, c2Info.i64)
			} else {
				recold.AddStr(c2Info.colname, c2Info.bytes)
				recnew.AddStr(c2Info.colname, c2Info.bytes)
			}
			scan.Cmp1, scan.Cmp2 = engine.CMP_GE, engine.CMP_LE
			scan.Key1, scan.Key2 = recold, recnew
			return nil
		case "<", "<=":
			if c2Info.datatype == engine.TYPE_INT64 {
				recold.AddInt64(c2Info.colname, c2Info.i64)
				recnew.AddInt64(c2Info.colname, c2Info.i64)
			} else {
				recold.AddStr(c2Info.colname, c2Info.bytes)
				recnew.AddStr(c2Info.colname, c2Info.bytes)
			}
			scan.Cmp1, scan.Cmp2 = engine.CMP_GE, engine.CMP_LT
			scan.Key1, scan.Key2 = recold, recnew

			if c2Info.op == "<=" {
				scan.Cmp1, scan.Cmp2 = engine.CMP_GE, engine.CMP_LE
			}
			return nil
		}
	}
	return nil
}

type condDataInfo struct {
	colname  string
	datatype uint32
	op       string // >, <, =, >=, <= , between
	// compare 数据
	i64   int64
	bytes []byte

	// between 数据
	i64range   []int64
	bytesrange []string
}

func condData(cond *ast.ConditionContext) *condDataInfo {
	c := &condDataInfo{}
	if cond.ComparisonCondition() != nil {
		compareCond := cond.ComparisonCondition()
		c.colname = compareCond.ColumnName().GetText()
		c.op = compareCond.OP().GetText()
		if compareCond.ColumnValue().INTEGER() != nil {
			c.datatype = engine.TYPE_INT64
			val := compareCond.ColumnValue().GetText()
			v, _ := strconv.ParseInt(val, 10, 64)
			c.i64 = v
		} else {
			c.datatype = engine.TYPE_BYTES
			c.bytes = []byte(compareCond.ColumnValue().GetText())
		}
	} else {
		// between condition
		c.op = "between"
		betweenCond := cond.BetweenCondition()
		c.colname = betweenCond.ColumnName().GetText()
		if betweenCond.ColumnValue(0).INTEGER() != nil {
			c.datatype = engine.TYPE_INT64
			lv, _ := strconv.ParseInt(betweenCond.ColumnValue(0).GetText(), 10, 64)
			rv, _ := strconv.ParseInt(betweenCond.ColumnValue(1).GetText(), 10, 64)
			c.i64range = append(c.i64range, lv, rv)
		} else {
			c.datatype = engine.TYPE_BYTES
			lv := betweenCond.ColumnValue(0).GetText()
			rv := betweenCond.ColumnValue(1).GetText()
			c.bytesrange = append(c.bytesrange, lv, rv)
		}
	}
	return c
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
		return
	}
	v := new(SQLParser)
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
	v := new(SQLParser)
	tableDef := v.Visit(tree).(*engine.Record)
	for i := 0; i < len(tableDef.Vals); i++ {
		if tableDef.Vals[i].Type == engine.TYPE_INT64 {
			fmt.Println("col:", tableDef.Cols[i], "val: ", tableDef.Vals[i].I64)
		} else {
			fmt.Println("col:", tableDef.Cols[i], "val: ", string(tableDef.Vals[i].Str))
		}
	}
}

func TestSelectTable(t *testing.T) {
	// SelectSql1 := `
	// SELECT name, id FROM users WHERE age = 18 AND height BETWEEN 170 AND 175;
	// `
	SelectSql2 := `
	SELECT name, id FROM users WHERE age > 18 AND age <= 25;
	`
	input := antlr.NewInputStream(SelectSql2)
	lexer := ast.NewSQLLexer(input)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	p := ast.NewSQLParser(tokenStream)
	p.AddErrorListener(antlr.NewDefaultErrorListener())

	tree := p.Sql()
	if tree == nil {
		t.Errorf("parse tree is nil, check input or parser configuration")
	}
	v := new(SQLParser)
	v.Visit(tree)
}
