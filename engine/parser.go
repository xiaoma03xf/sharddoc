package engine

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itxiaoma0610/sharddoc/engine/parser/ast"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	t := new(TableDef)
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		switch c := child.(type) {
		case *ast.TableNameContext:
			t.Name = a.VisitTableName(c).(string)
		case *ast.ColumnDefinitionsContext:
			defs := a.VisitColumnDefinitions(c).(*TableDef)
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
	t := &TableDef{}
	defs := ctx.AllColumnDefinition()
	for _, def := range defs {
		t.Cols = append(t.Cols, def.GetStart().GetText())
		ctype := def.GetStop().GetText()
		switch ctype {
		case "INT64":
			t.Types = append(t.Types, TYPE_INT64)
		case "BYTES":
			t.Types = append(t.Types, TYPE_BYTES)
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

	rec := &Record{}
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
	">":  CMP_GT,
	">=": CMP_GE,
	"<":  CMP_LT,
	"<=": CMP_LE,
}

type SelectInfo struct {
	TableName   string
	SelectField []string
	Scan        *Scanner
}

func (v *SQLParser) VisitSelectTableStatement(ctx *ast.SelectTableStatementContext) interface{} {
	// todo current support select sync, 暂时只支持两个条件
	// select * from table where age > 18 and age < 40 相同键位
	// select * from table where age = 18 and height = 175
	// select * from table where age = 18 and height bewteen 170 and 175
	var err error
	var selectField []string
	scan := &Scanner{}

	condtype := conditionType(ctx.AllCondition()...)
	switch condtype {
	case SIG_PARE:
		err = selectPare(scan, ctx)
	case SIG_BEWTEEN:
		err = selectBewteen(scan, ctx)
	case PARE_LINK_PARE:
		err = selectPareLinkPare(scan, ctx)
	case PARE_LINK_BETWEEN:
		err = selectPareLinkBetween(scan, ctx)
	case UNSUPPORTED_TYPE:
		return fmt.Errorf("unsupported type")
	}
	if err != nil {
		return fmt.Errorf("build scanner err")
	}

	if ctx.SelectColumnNames().STAR() != nil {
		selectField = append(selectField, "*")
	} else {
		for _, col := range ctx.SelectColumnNames().AllColumnName() {
			selectField = append(selectField, col.GetText())
		}
	}
	return &SelectInfo{
		TableName:   ctx.TableName().GetText(),
		SelectField: selectField,
		Scan:        scan,
	}
}

func selectPare(scan *Scanner, ctx *ast.SelectTableStatementContext) error {
	cond := ctx.AllCondition()[0].ComparisonCondition()
	op := cond.OP().GetText()
	colname := cond.ColumnName().GetText()

	recold := Record{}
	if cond.ColumnValue().INTEGER() != nil {
		val, _ := strconv.ParseInt(cond.ColumnValue().INTEGER().GetText(), 10, 64)
		recold.AddInt64(colname, val)
	} else {
		recold.AddStr(colname, []byte(cond.ColumnValue().GetText()))
	}

	recnew := Record{}
	switch op {
	case ">":
		// age > 18
		if cond.ColumnValue().INTEGER() != nil {
			recnew.AddInt64(colname, math.MaxInt64/2)
		} else {
			// name 最大长度32, 规定编码规则
			recnew.AddStr(colname, []byte(MAX_NAME))
		}
		scan.Cmp1 = CMP_GT
		scan.Cmp2 = CMP_LE
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
		scan.Cmp1 = CMP_GE
		scan.Cmp2 = CMP_LE
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
		scan.Cmp1 = CMP_GE
		scan.Cmp2 = CMP_LE
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
		scan.Cmp1 = CMP_GE
		scan.Cmp2 = CMP_LT
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
		scan.Cmp1 = CMP_GE
		scan.Cmp2 = CMP_LE
		scan.Key1 = recnew
		scan.Key2 = recold
		return nil
	}
	return nil
}
func selectBewteen(scan *Scanner, ctx *ast.SelectTableStatementContext) error {
	condInfo := condData(ctx.AllCondition()[0].(*ast.ConditionContext))
	recl := Record{}
	recr := Record{}
	if condInfo.datatype == TYPE_BYTES {
		if err := leftLessThanRight(string(condInfo.bytesrange[0]),
			string(condInfo.bytesrange[1])); err != nil {
			return fmt.Errorf("check syntax err, left should less than right")
		}
		recl.AddStr(condInfo.colname, []byte(condInfo.bytesrange[0]))
		recr.AddStr(condInfo.colname, []byte(condInfo.bytesrange[1]))
	} else {
		if err := leftLessThanRight(int64(condInfo.i64range[0]),
			int64(condInfo.i64range[1])); err != nil {
			return fmt.Errorf("check syntax err, left should less than right")
		}
		recl.AddInt64(condInfo.colname, condInfo.i64range[0])
		recr.AddInt64(condInfo.colname, condInfo.i64range[1])
	}
	scan.Cmp1 = CMP_GE
	scan.Cmp2 = CMP_LE
	scan.Key1 = recl
	scan.Key2 = recr
	return nil
}
func condInt64Data(data *ast.ColumnValueContext) int64 {
	if data.INTEGER() == nil {
		panic("data type err")
	}
	val, _ := strconv.ParseInt(data.GetText(), 10, 64)
	return val
}

func selectPareLinkPare(scan *Scanner, ctx *ast.SelectTableStatementContext) error {
	// 两种情况, todo 暂时规定左大右小
	// select * from table where age > 18 and age < 40 相同键位
	// select * from table where age = 18 and height = 175
	allCond := ctx.AllCondition()
	cond1Info := condData(allCond[0].(*ast.ConditionContext))
	cond2Info := condData(allCond[1].(*ast.ConditionContext))

	recold := Record{}
	recnew := Record{}

	if cond1Info.colname == cond2Info.colname {
		// select * from table where age > 18 and age < 40 相同键位
		// 左边仅允许使用 > || >= , 右边仅允许 < || <=
		if !(cond1Info.op == ">" || cond1Info.op == ">=") ||
			!(cond2Info.op == "<" || cond2Info.op == "<=") {
			return fmt.Errorf("syntax err, left hand only support > || >=, right hand only support < || <=")
		}
		scan.Cmp1 = Op2Cmp[cond1Info.op]
		scan.Cmp2 = Op2Cmp[cond2Info.op]
		if cond1Info.datatype == TYPE_INT64 {
			recold.AddInt64(cond1Info.colname, cond1Info.i64)
			recnew.AddInt64(cond2Info.colname, cond2Info.i64)
			scan.Key1, scan.Key2 = recold, recnew
			return nil
		} else {
			recold.AddStr(cond1Info.colname, cond1Info.bytes)
			recnew.AddStr(cond2Info.colname, cond2Info.bytes)
			scan.Key1, scan.Key2 = recold, recnew
			return nil
		}
	} else {
		// select * from table where age = 18 and height = 175
		// 左边age仅允许为=,右边任意
		if cond1Info.op != "=" {
			return fmt.Errorf("not supported operation")
		}
		if cond1Info.datatype == TYPE_INT64 {
			recold.AddInt64(cond1Info.colname, cond1Info.i64)
			recnew.AddInt64(cond1Info.colname, cond1Info.i64)
		} else {
			recold.AddStr(cond1Info.colname, cond1Info.bytes)
			recnew.AddStr(cond1Info.colname, cond1Info.bytes)
		}

		switch cond2Info.op {
		case ">", ">=":
			// height > 175
			if cond2Info.datatype == TYPE_INT64 {
				recold.AddInt64(cond2Info.colname, cond2Info.i64)
				recnew.AddInt64(cond2Info.colname, math.MaxInt64/2)
			} else {
				recold.AddStr(cond2Info.colname, cond2Info.bytes)
				recnew.AddStr(cond2Info.colname, []byte(MAX_NAME))
			}
			scan.Cmp1, scan.Cmp2 = CMP_GT, CMP_LE
			if cond2Info.op == ">=" {
				scan.Cmp1, scan.Cmp2 = CMP_GE, CMP_LE
			}
			scan.Key1, scan.Key2 = recold, recnew
			return nil
		case "=":
			if cond2Info.datatype == TYPE_INT64 {
				recold.AddInt64(cond2Info.colname, cond2Info.i64)
				recnew.AddInt64(cond2Info.colname, cond2Info.i64)
			} else {
				recold.AddStr(cond2Info.colname, cond2Info.bytes)
				recnew.AddStr(cond2Info.colname, cond2Info.bytes)
			}
			scan.Cmp1, scan.Cmp2 = CMP_GE, CMP_LE
			scan.Key1, scan.Key2 = recold, recnew
			return nil
		case "<", "<=":
			if cond2Info.datatype == TYPE_INT64 {
				recold.AddInt64(cond2Info.colname, math.MinInt64/2)
				recnew.AddInt64(cond2Info.colname, cond2Info.i64)
			} else {
				recold.AddStr(cond2Info.colname, []byte(MIN_NAME))
				recnew.AddStr(cond2Info.colname, cond2Info.bytes)
			}
			scan.Cmp1, scan.Cmp2 = CMP_GE, CMP_LT
			if cond2Info.op == "<=" {
				scan.Cmp1, scan.Cmp2 = CMP_GE, CMP_LE
			}
			scan.Key1, scan.Key2 = recold, recnew
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
			c.datatype = TYPE_INT64
			val := compareCond.ColumnValue().GetText()
			v, _ := strconv.ParseInt(val, 10, 64)
			c.i64 = v
		} else {
			c.datatype = TYPE_BYTES
			c.bytes = []byte(compareCond.ColumnValue().GetText())
		}
	} else {
		// between condition
		c.op = "between"
		betweenCond := cond.BetweenCondition()
		c.colname = betweenCond.ColumnName().GetText()
		if betweenCond.ColumnValue(0).INTEGER() != nil {
			c.datatype = TYPE_INT64
			lv, _ := strconv.ParseInt(betweenCond.ColumnValue(0).GetText(), 10, 64)
			rv, _ := strconv.ParseInt(betweenCond.ColumnValue(1).GetText(), 10, 64)
			c.i64range = append(c.i64range, lv, rv)
		} else {
			c.datatype = TYPE_BYTES
			lv := betweenCond.ColumnValue(0).GetText()
			rv := betweenCond.ColumnValue(1).GetText()
			c.bytesrange = append(c.bytesrange, lv, rv)
		}
	}
	return c
}
func leftLessThanRight(data ...any) error {
	if len(data) < 2 {
		return fmt.Errorf("data length less than 2, check syntax")
	}
	switch data[0].(type) {
	case int64:
		if data[0].(int64) >= data[1].(int64) {
			return fmt.Errorf("check syntax err, left should less than right")
		}
	case string:
		if data[0].(string) >= data[1].(string) {
			return fmt.Errorf("check syntax err, left should less than right")
		}
	default:
		return fmt.Errorf("not supported type")
	}
	return nil
}
func selectPareLinkBetween(scan *Scanner, ctx *ast.SelectTableStatementContext) error {
	// select * from table where age =18 and height between 170 and 175
	allCond := ctx.AllCondition()

	// 左边字段名等于右边或左边操作不等于 "="
	if allCond[0].ComparisonCondition().ColumnName().GetText() == allCond[1].BetweenCondition().ColumnName().GetText() ||
		allCond[0].ComparisonCondition().OP().GetText() != "=" {
		return fmt.Errorf("syntax error: %v", ctx.GetText())
	}
	recold := Record{}
	recnew := Record{}
	cond1Info := condData(allCond[0].(*ast.ConditionContext))
	if cond1Info.datatype == TYPE_BYTES {
		recold.AddStr(cond1Info.colname, cond1Info.bytes)
		recnew.AddStr(cond1Info.colname, cond1Info.bytes)
	} else {
		recold.AddInt64(cond1Info.colname, cond1Info.i64)
		recnew.AddInt64(cond1Info.colname, cond1Info.i64)
	}

	// 校验between语句的合法性
	cond2Info := condData(allCond[1].(*ast.ConditionContext))
	if cond2Info.datatype == TYPE_BYTES {
		err := leftLessThanRight(string(cond2Info.bytesrange[0]), string(cond2Info.bytesrange[1]))
		if err != nil {
			return fmt.Errorf("between syntax err:%v", allCond[1].GetText())
		}
		recold.AddStr(cond1Info.colname, []byte(cond2Info.bytesrange[0]))
		recnew.AddStr(cond1Info.colname, []byte(cond2Info.bytesrange[1]))
	} else {
		err := leftLessThanRight(int64(cond2Info.i64range[0]), int64(cond2Info.i64range[1]))
		if err != nil {
			return fmt.Errorf("between syntax err:%v", allCond[1].GetText())
		}
		recold.AddInt64(cond2Info.colname, cond2Info.i64range[0])
		recnew.AddInt64(cond2Info.colname, cond2Info.i64range[1])
	}

	scan.Cmp1 = CMP_GE
	scan.Cmp2 = CMP_LE
	scan.Key1 = recold
	scan.Key2 = recnew

	return nil
}
func reduceSelectData(scan *Scanner) []Record {
	got := []Record{}
	for scan.Valid() {
		rec := Record{}
		scan.Deref(&rec)
		got = append(got, rec)
		scan.Next()
	}
	return got
}

func toCamel(s string) string {
	return cases.Title(language.English).String(s)
}

func VisitTree(sql string) interface{} {
	input := antlr.NewInputStream(sql)
	lexer := ast.NewSQLLexer(input)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	p := ast.NewSQLParser(tokenStream)
	p.AddErrorListener(antlr.NewDefaultErrorListener())

	tree := p.Sql()
	if tree == nil {
		panic("parse tree is nil, check input or parser configuration")
	}
	v := new(SQLParser)
	return v.Visit(tree)
}

func mapStructFieldsOnce(ptr any) map[string]reflect.Value {
	val := reflect.ValueOf(ptr).Elem() // 获取结构体值
	typ := val.Type()

	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" { // 非导出字段跳过
			continue
		}
		fieldMap[strings.ToLower(field.Name)] = val.Field(i)
	}
	return fieldMap
}
func convertValue(v Value) (any, error) {
	switch v.Type {
	case TYPE_INT64:
		return int(v.I64), nil
	case TYPE_BYTES:
		return string(v.Str), nil
	default:
		return nil, fmt.Errorf("unsupported value type: %d", v.Type)
	}
}
func scanRecordToStruct(record Record, ptr any, fields map[string]reflect.Value) error {
	for i, col := range record.Cols {
		f, ok := fields[strings.ToLower(col)]
		if !ok {
			continue
		}
		val, err := convertValue(record.Vals[i])
		if err != nil {
			return err
		}
		rv := reflect.ValueOf(val)
		if rv.Type().AssignableTo(f.Type()) {
			f.Set(rv)
		} else if rv.Type().ConvertibleTo(f.Type()) {
			f.Set(rv.Convert(f.Type()))
		} else {
			return fmt.Errorf("cannot assign %v to field %v", rv.Type(), f.Type())
		}
	}
	return nil
}
func scanRecordsToStructs(records []Record, ptrs []any) error {
	fields := mapStructFieldsOnce(ptrs[0])
	for i, record := range records {
		if i >= len(ptrs) {
			break
		}
		if err := scanRecordToStruct(record, ptrs[i], fields); err != nil {
			return err
		}
	}
	return nil
}

