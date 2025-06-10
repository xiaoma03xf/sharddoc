package engine

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/xiaoma03xf/sharddoc/engine/parser/ast"
	"github.com/xiaoma03xf/sharddoc/lib/logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type SQLParser struct {
	*ast.BaseSQLParserVisitor
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
		case *ast.UpdateTableStatementContext:
			return a.VisitUpdateTableStatement(c)
		case *ast.DeleteTableStatementContext:
			return a.VisitDeleteTableStatement(c)
		}
	}
	return nil
}

func (a *SQLParser) VisitCreateTableStatement(ctx *ast.CreateTableStatementContext) interface{} {
	// CreatTableStatement return
	t := new(TableDef)
	t.Name = ctx.TableName().GetText()

	alldefs := ctx.ColumnDefinitions().AllColumnDefinition()
	for _, def := range alldefs {
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

	// 添加主键索引
	t.Indexes = append(t.Indexes, []string{ctx.ColumnName().GetText()})
	// 添加其他索引
	for _, indexdef := range ctx.IndexDefinitions().AllIndexDefinition() {
		index := []string{}
		for _, iname := range indexdef.AllColumnName() {
			index = append(index, iname.GetText())
		}
		t.Indexes = append(t.Indexes, index)
	}

	return t
}

func (a *SQLParser) VisitInsertTableStatement(ctx *ast.InsertTableStatementContext) interface{} {
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
	// 返回插入数据和表名
	return &InsertRes{TableName: ctx.TableName().GetText(), Rec: rec}
}

type CondType int

const (
	SIG_PARE          CondType = iota
	SIG_BEWTEEN                // age between 18 and 24
	PARE_LINK_PARE             // age ==18 and height > 24
	PARE_LINK_BETWEEN          //age = 18 and height between 170 and 175
	UNSUPPORTED_TYPE
)

const MIN_NAME = ""
const MAX_NAME = "\xff\xff\xff\xff\xff\xff\xff\xff" // 足够长的最大字节

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

func conditionType(cond ...ast.IConditionContext) CondType {
	if len(cond) >= 3 {
		return UNSUPPORTED_TYPE
	}

	if len(cond) == 1 {
		if cond[0].BetweenCondition() != nil {
			return SIG_BEWTEEN
		} else {
			return SIG_PARE
		}
	} else if len(cond) == 2 {
		if cond[1].ComparisonCondition() != nil {
			return PARE_LINK_PARE
		}
		return PARE_LINK_BETWEEN
	}
	return UNSUPPORTED_TYPE
}

func (v *SQLParser) VisitConditions(ctx *ast.ConditionsContext) interface{} {
	var err error
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

	return scan
}

func selectPare(scan *Scanner, ctx *ast.ConditionsContext) error {
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
func selectBewteen(scan *Scanner, ctx *ast.ConditionsContext) error {
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

func selectPareLinkPare(scan *Scanner, ctx *ast.ConditionsContext) error {
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
func selectPareLinkBetween(scan *Scanner, ctx *ast.ConditionsContext) error {
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
func (v *SQLParser) VisitSelectTableStatement(ctx *ast.SelectTableStatementContext) interface{} {
	// todo current support select sync, 暂时只支持两个条件
	// select * from table where age > 18 and age < 40 相同键位
	// select * from table where age = 18 and height = 175
	// select * from table where age = 18 and height bewteen 170 and 175

	scan := v.VisitConditions(ctx.Conditions().(*ast.ConditionsContext)).(*Scanner)

	var selectField []string
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

func toCamel(s string) string {
	return cases.Title(language.English).String(s)
}

// update 语句
func (v *SQLParser) VisitUpdateTableStatement(ctx *ast.UpdateTableStatementContext) interface{} {
	// 根据条件查出所有 Record 记录
	var err error
	for _, setclause := range ctx.SetClauses().AllSetClause() {
		if setclause.OP().GetText() != "=" {
			err = fmt.Errorf("updata OP not equal")
		}
	}

	setmp := make(map[string]string)
	scan := v.VisitConditions(ctx.Conditions().(*ast.ConditionsContext)).(*Scanner)
	allsets := ctx.SetClauses().AllSetClause()
	for _, set := range allsets {
		setmp[set.ColumnName().GetText()] = set.ColumnValue().GetText()
	}
	return &UpdateRes{
		TableName: ctx.TableName().GetText(),
		UpdateMp:  setmp,
		Scan:      scan,
		Err:       err,
	}
}

func (v *SQLParser) VisitDeleteTableStatement(ctx *ast.DeleteTableStatementContext) interface{} {
	return &DelRes{
		TableName: ctx.TableName().GetText(),
		Scan:      v.VisitConditions(ctx.Conditions().(*ast.ConditionsContext)).(*Scanner),
	}
}

// func (db *DB) Begin()
// select 语句, 查询不需要建立事务
// db.Raw("SELECT name, id FROM tbl_test WHERE age = 18").Scan(&[]User{})
func (db *DB) Raw(sql string) (q *QueryResult) {
	q = &QueryResult{}
	selectInfo := VisitTree(sql).(*SelectInfo)

	// Since this is a read-only operation, we can use a transaction
	// without committing it, as we only need its snapshot for reading
	tx := &DBTX{}
	db.Begin(tx)

	err := tx.Scan(selectInfo.TableName, selectInfo.Scan)
	if err != nil {
		db.Abort(tx)
		return &QueryResult{nil, err}
	}
	err = db.Commit(tx)
	if err != nil {
		q.Err = err
	}
	q.Recs = reduceSelectData(selectInfo.Scan)
	return
}

// 增删改,建表
func (db *DB) Exec(sql string) error {
	tx := DBTX{}
	db.Begin(&tx)

	execres := VisitTree(sql)
	switch res := execres.(type) {
	case *TableDef:
		// 创表语句
		if err := tx.TableNew(res); err != nil {
			logger.Warn(sql, "create table err:", err)
			db.Abort(&tx)
		}
	case *InsertRes:
		if _, err := tx.Insert(res.TableName, *res.Rec); err != nil {
			logger.Warn(sql, "insert data err:", err)
			db.Abort(&tx)
			return fmt.Errorf("insert data %v, err:%v", res.Rec, err)
		}
	case *UpdateRes:
		if res.Err != nil {
			logger.Warn(sql, "update data err:", res.Err)
			db.Abort(&tx)
			return res.Err
		}
		// 1.根据条件查询出需要更新的记录
		if res.Scan == nil {
			logger.Warn(sql, "update data err: scan is nil")
			db.Abort(&tx)
			return fmt.Errorf("update data err: scan is nil")
		}
		if err := tx.Scan(res.TableName, res.Scan); err != nil {
			logger.Warn(sql, "update data err:", err)
			db.Abort(&tx)
		}
		Rec := reduceSelectData(res.Scan)

		// 2.更新记录
		for _, rec := range Rec {
			// 遍历每一条需要更新的值
			for col, newval := range res.UpdateMp {
				// 依次更新查询记录中的值
				for i, v := range rec.Cols {
					if v == col {
						switch rec.Vals[i].Type {
						case TYPE_INT64:
							newvalInt, _ := strconv.ParseInt(newval, 10, 64)
							rec.Vals[i].I64 = newvalInt
						case TYPE_BYTES:
							rec.Vals[i].Str = []byte(newval)
						default:
							db.Abort(&tx)
							return fmt.Errorf("unsupported type for update")
						}
					}
				}
			}
			_, err := tx.Update(res.TableName, rec)
			if err != nil {
				logger.Warn(sql, "update data err:", err)
				db.Abort(&tx)
				return fmt.Errorf("updata rec %v, err: %v", rec, err)
			}
		}
	case *DelRes:
		// 1.根据条件查询出需要更新的记录
		if res.Scan == nil {
			logger.Warn(sql, "update data err: scan is nil")
			db.Abort(&tx)
			return fmt.Errorf("update data err: scan is nil")
		}
		if err := tx.Scan(res.TableName, res.Scan); err != nil {
			logger.Warn(sql, "update data err:", err)
			db.Abort(&tx)
		}
		Rec := reduceSelectData(res.Scan)
		for _, rec := range Rec {
			_, err := tx.Delete(res.TableName, rec)
			if err != nil {
				logger.Warn(sql, "delete data err:", err)
				db.Abort(&tx)
				return fmt.Errorf("delete rec %v, err: %v", rec, err)
			}
		}
	default:
		return fmt.Errorf("unsupported sql type")
	}

	return db.Commit(&tx)
}
func BuildInsertSQL(table string, columns []string, values []interface{}) string {
	var cols string
	var vals []string
	for i, val := range values {
		cols += columns[i]
		if i != len(values)-1 {
			cols += ", "
		}
		vals = append(vals, formatValue(val))
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", table, cols, strings.Join(vals, ", "))
}
func formatValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		// 注意转义字符串中间的引号
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	default:
		return "NULL"
	}
}
