// Code generated from SQLParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package ast // SQLParser
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type SQLParser struct {
	*antlr.BaseParser
}

var SQLParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func sqlparserParserInit() {
	staticData := &SQLParserParserStaticData
	staticData.LiteralNames = []string{
		"", "'CREATE'", "'TABLE'", "'PRIMARY'", "'KEY'", "'INDEX'", "'INSERT'",
		"'INTO'", "'VALUES'", "'SELECT'", "'FROM'", "'WHERE'", "'AND'", "'DELETE'",
		"'UPDATE'", "'SET'", "'INT64'", "'BYTES'", "'BETWEEN'", "", "'('", "')'",
		"','", "';'", "'*'",
	}
	staticData.SymbolicNames = []string{
		"", "CREATE", "TABLE", "PRIMARY", "KEY", "INDEX", "INSERT", "INTO",
		"VALUES", "SELECT", "FROM", "WHERE", "AND", "DELETE", "UPDATE", "SET",
		"INT64", "BYTES", "BETWEEN", "IDENTIFIER", "LPAREN", "RPAREN", "COMMA",
		"SEMICOLON", "STAR", "WS", "INTEGER", "STRING", "OP",
	}
	staticData.RuleNames = []string{
		"sql", "createTableStatement", "tableName", "columnName", "columnType",
		"columnDefinitions", "columnDefinition", "indexDefinitions", "indexDefinition",
		"insertTableStatement", "columnInsertNames", "columnInsertValues", "columnValue",
		"selectTableStatement", "conditions", "selectColumnNames", "condition",
		"comparisonCondition", "betweenCondition", "updateTableStatement", "setClauses",
		"setClause", "deleteTableStatement",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 28, 201, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 3, 0, 52, 8, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5,
		1, 5, 5, 5, 79, 8, 5, 10, 5, 12, 5, 82, 9, 5, 1, 6, 1, 6, 1, 6, 1, 7, 1,
		7, 1, 7, 5, 7, 90, 8, 7, 10, 7, 12, 7, 93, 9, 7, 1, 8, 1, 8, 1, 8, 1, 8,
		1, 8, 5, 8, 100, 8, 8, 10, 8, 12, 8, 103, 9, 8, 1, 8, 1, 8, 1, 9, 1, 9,
		1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1,
		10, 5, 10, 121, 8, 10, 10, 10, 12, 10, 124, 9, 10, 1, 11, 1, 11, 1, 11,
		5, 11, 129, 8, 11, 10, 11, 12, 11, 132, 9, 11, 1, 12, 1, 12, 1, 13, 1,
		13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 1, 14, 5, 14, 146,
		8, 14, 10, 14, 12, 14, 149, 9, 14, 3, 14, 151, 8, 14, 1, 15, 1, 15, 1,
		15, 1, 15, 5, 15, 157, 8, 15, 10, 15, 12, 15, 160, 9, 15, 3, 15, 162, 8,
		15, 1, 16, 1, 16, 3, 16, 166, 8, 16, 1, 17, 1, 17, 1, 17, 1, 17, 1, 18,
		1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1,
		19, 1, 20, 1, 20, 1, 20, 5, 20, 187, 8, 20, 10, 20, 12, 20, 190, 9, 20,
		1, 21, 1, 21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 0,
		0, 23, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34,
		36, 38, 40, 42, 44, 0, 2, 1, 0, 16, 17, 1, 0, 26, 27, 192, 0, 51, 1, 0,
		0, 0, 2, 53, 1, 0, 0, 0, 4, 69, 1, 0, 0, 0, 6, 71, 1, 0, 0, 0, 8, 73, 1,
		0, 0, 0, 10, 75, 1, 0, 0, 0, 12, 83, 1, 0, 0, 0, 14, 86, 1, 0, 0, 0, 16,
		94, 1, 0, 0, 0, 18, 106, 1, 0, 0, 0, 20, 117, 1, 0, 0, 0, 22, 125, 1, 0,
		0, 0, 24, 133, 1, 0, 0, 0, 26, 135, 1, 0, 0, 0, 28, 150, 1, 0, 0, 0, 30,
		161, 1, 0, 0, 0, 32, 165, 1, 0, 0, 0, 34, 167, 1, 0, 0, 0, 36, 171, 1,
		0, 0, 0, 38, 177, 1, 0, 0, 0, 40, 183, 1, 0, 0, 0, 42, 191, 1, 0, 0, 0,
		44, 195, 1, 0, 0, 0, 46, 52, 3, 2, 1, 0, 47, 52, 3, 18, 9, 0, 48, 52, 3,
		26, 13, 0, 49, 52, 3, 38, 19, 0, 50, 52, 3, 44, 22, 0, 51, 46, 1, 0, 0,
		0, 51, 47, 1, 0, 0, 0, 51, 48, 1, 0, 0, 0, 51, 49, 1, 0, 0, 0, 51, 50,
		1, 0, 0, 0, 52, 1, 1, 0, 0, 0, 53, 54, 5, 1, 0, 0, 54, 55, 5, 2, 0, 0,
		55, 56, 3, 4, 2, 0, 56, 57, 5, 20, 0, 0, 57, 58, 3, 10, 5, 0, 58, 59, 5,
		22, 0, 0, 59, 60, 5, 3, 0, 0, 60, 61, 5, 4, 0, 0, 61, 62, 5, 20, 0, 0,
		62, 63, 3, 6, 3, 0, 63, 64, 5, 21, 0, 0, 64, 65, 5, 22, 0, 0, 65, 66, 3,
		14, 7, 0, 66, 67, 5, 21, 0, 0, 67, 68, 5, 23, 0, 0, 68, 3, 1, 0, 0, 0,
		69, 70, 5, 19, 0, 0, 70, 5, 1, 0, 0, 0, 71, 72, 5, 19, 0, 0, 72, 7, 1,
		0, 0, 0, 73, 74, 7, 0, 0, 0, 74, 9, 1, 0, 0, 0, 75, 80, 3, 12, 6, 0, 76,
		77, 5, 22, 0, 0, 77, 79, 3, 12, 6, 0, 78, 76, 1, 0, 0, 0, 79, 82, 1, 0,
		0, 0, 80, 78, 1, 0, 0, 0, 80, 81, 1, 0, 0, 0, 81, 11, 1, 0, 0, 0, 82, 80,
		1, 0, 0, 0, 83, 84, 3, 6, 3, 0, 84, 85, 3, 8, 4, 0, 85, 13, 1, 0, 0, 0,
		86, 91, 3, 16, 8, 0, 87, 88, 5, 22, 0, 0, 88, 90, 3, 16, 8, 0, 89, 87,
		1, 0, 0, 0, 90, 93, 1, 0, 0, 0, 91, 89, 1, 0, 0, 0, 91, 92, 1, 0, 0, 0,
		92, 15, 1, 0, 0, 0, 93, 91, 1, 0, 0, 0, 94, 95, 5, 5, 0, 0, 95, 96, 5,
		20, 0, 0, 96, 101, 3, 6, 3, 0, 97, 98, 5, 22, 0, 0, 98, 100, 3, 6, 3, 0,
		99, 97, 1, 0, 0, 0, 100, 103, 1, 0, 0, 0, 101, 99, 1, 0, 0, 0, 101, 102,
		1, 0, 0, 0, 102, 104, 1, 0, 0, 0, 103, 101, 1, 0, 0, 0, 104, 105, 5, 21,
		0, 0, 105, 17, 1, 0, 0, 0, 106, 107, 5, 6, 0, 0, 107, 108, 5, 7, 0, 0,
		108, 109, 3, 4, 2, 0, 109, 110, 5, 20, 0, 0, 110, 111, 3, 20, 10, 0, 111,
		112, 5, 21, 0, 0, 112, 113, 5, 8, 0, 0, 113, 114, 5, 20, 0, 0, 114, 115,
		3, 22, 11, 0, 115, 116, 5, 21, 0, 0, 116, 19, 1, 0, 0, 0, 117, 122, 3,
		6, 3, 0, 118, 119, 5, 22, 0, 0, 119, 121, 3, 6, 3, 0, 120, 118, 1, 0, 0,
		0, 121, 124, 1, 0, 0, 0, 122, 120, 1, 0, 0, 0, 122, 123, 1, 0, 0, 0, 123,
		21, 1, 0, 0, 0, 124, 122, 1, 0, 0, 0, 125, 130, 3, 24, 12, 0, 126, 127,
		5, 22, 0, 0, 127, 129, 3, 24, 12, 0, 128, 126, 1, 0, 0, 0, 129, 132, 1,
		0, 0, 0, 130, 128, 1, 0, 0, 0, 130, 131, 1, 0, 0, 0, 131, 23, 1, 0, 0,
		0, 132, 130, 1, 0, 0, 0, 133, 134, 7, 1, 0, 0, 134, 25, 1, 0, 0, 0, 135,
		136, 5, 9, 0, 0, 136, 137, 3, 30, 15, 0, 137, 138, 5, 10, 0, 0, 138, 139,
		3, 4, 2, 0, 139, 140, 3, 28, 14, 0, 140, 27, 1, 0, 0, 0, 141, 142, 5, 11,
		0, 0, 142, 147, 3, 32, 16, 0, 143, 144, 5, 12, 0, 0, 144, 146, 3, 32, 16,
		0, 145, 143, 1, 0, 0, 0, 146, 149, 1, 0, 0, 0, 147, 145, 1, 0, 0, 0, 147,
		148, 1, 0, 0, 0, 148, 151, 1, 0, 0, 0, 149, 147, 1, 0, 0, 0, 150, 141,
		1, 0, 0, 0, 150, 151, 1, 0, 0, 0, 151, 29, 1, 0, 0, 0, 152, 162, 5, 24,
		0, 0, 153, 158, 3, 6, 3, 0, 154, 155, 5, 22, 0, 0, 155, 157, 3, 6, 3, 0,
		156, 154, 1, 0, 0, 0, 157, 160, 1, 0, 0, 0, 158, 156, 1, 0, 0, 0, 158,
		159, 1, 0, 0, 0, 159, 162, 1, 0, 0, 0, 160, 158, 1, 0, 0, 0, 161, 152,
		1, 0, 0, 0, 161, 153, 1, 0, 0, 0, 162, 31, 1, 0, 0, 0, 163, 166, 3, 34,
		17, 0, 164, 166, 3, 36, 18, 0, 165, 163, 1, 0, 0, 0, 165, 164, 1, 0, 0,
		0, 166, 33, 1, 0, 0, 0, 167, 168, 3, 6, 3, 0, 168, 169, 5, 28, 0, 0, 169,
		170, 3, 24, 12, 0, 170, 35, 1, 0, 0, 0, 171, 172, 3, 6, 3, 0, 172, 173,
		5, 18, 0, 0, 173, 174, 3, 24, 12, 0, 174, 175, 5, 12, 0, 0, 175, 176, 3,
		24, 12, 0, 176, 37, 1, 0, 0, 0, 177, 178, 5, 14, 0, 0, 178, 179, 3, 4,
		2, 0, 179, 180, 5, 15, 0, 0, 180, 181, 3, 40, 20, 0, 181, 182, 3, 28, 14,
		0, 182, 39, 1, 0, 0, 0, 183, 188, 3, 42, 21, 0, 184, 185, 5, 22, 0, 0,
		185, 187, 3, 42, 21, 0, 186, 184, 1, 0, 0, 0, 187, 190, 1, 0, 0, 0, 188,
		186, 1, 0, 0, 0, 188, 189, 1, 0, 0, 0, 189, 41, 1, 0, 0, 0, 190, 188, 1,
		0, 0, 0, 191, 192, 3, 6, 3, 0, 192, 193, 5, 28, 0, 0, 193, 194, 3, 24,
		12, 0, 194, 43, 1, 0, 0, 0, 195, 196, 5, 13, 0, 0, 196, 197, 5, 10, 0,
		0, 197, 198, 3, 4, 2, 0, 198, 199, 3, 28, 14, 0, 199, 45, 1, 0, 0, 0, 12,
		51, 80, 91, 101, 122, 130, 147, 150, 158, 161, 165, 188,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// SQLParserInit initializes any static state used to implement SQLParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewSQLParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func SQLParserInit() {
	staticData := &SQLParserParserStaticData
	staticData.once.Do(sqlparserParserInit)
}

// NewSQLParser produces a new parser instance for the optional input antlr.TokenStream.
func NewSQLParser(input antlr.TokenStream) *SQLParser {
	SQLParserInit()
	this := new(SQLParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &SQLParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "SQLParser.g4"

	return this
}

// SQLParser tokens.
const (
	SQLParserEOF        = antlr.TokenEOF
	SQLParserCREATE     = 1
	SQLParserTABLE      = 2
	SQLParserPRIMARY    = 3
	SQLParserKEY        = 4
	SQLParserINDEX      = 5
	SQLParserINSERT     = 6
	SQLParserINTO       = 7
	SQLParserVALUES     = 8
	SQLParserSELECT     = 9
	SQLParserFROM       = 10
	SQLParserWHERE      = 11
	SQLParserAND        = 12
	SQLParserDELETE     = 13
	SQLParserUPDATE     = 14
	SQLParserSET        = 15
	SQLParserINT64      = 16
	SQLParserBYTES      = 17
	SQLParserBETWEEN    = 18
	SQLParserIDENTIFIER = 19
	SQLParserLPAREN     = 20
	SQLParserRPAREN     = 21
	SQLParserCOMMA      = 22
	SQLParserSEMICOLON  = 23
	SQLParserSTAR       = 24
	SQLParserWS         = 25
	SQLParserINTEGER    = 26
	SQLParserSTRING     = 27
	SQLParserOP         = 28
)

// SQLParser rules.
const (
	SQLParserRULE_sql                  = 0
	SQLParserRULE_createTableStatement = 1
	SQLParserRULE_tableName            = 2
	SQLParserRULE_columnName           = 3
	SQLParserRULE_columnType           = 4
	SQLParserRULE_columnDefinitions    = 5
	SQLParserRULE_columnDefinition     = 6
	SQLParserRULE_indexDefinitions     = 7
	SQLParserRULE_indexDefinition      = 8
	SQLParserRULE_insertTableStatement = 9
	SQLParserRULE_columnInsertNames    = 10
	SQLParserRULE_columnInsertValues   = 11
	SQLParserRULE_columnValue          = 12
	SQLParserRULE_selectTableStatement = 13
	SQLParserRULE_conditions           = 14
	SQLParserRULE_selectColumnNames    = 15
	SQLParserRULE_condition            = 16
	SQLParserRULE_comparisonCondition  = 17
	SQLParserRULE_betweenCondition     = 18
	SQLParserRULE_updateTableStatement = 19
	SQLParserRULE_setClauses           = 20
	SQLParserRULE_setClause            = 21
	SQLParserRULE_deleteTableStatement = 22
)

// ISqlContext is an interface to support dynamic dispatch.
type ISqlContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CreateTableStatement() ICreateTableStatementContext
	InsertTableStatement() IInsertTableStatementContext
	SelectTableStatement() ISelectTableStatementContext
	UpdateTableStatement() IUpdateTableStatementContext
	DeleteTableStatement() IDeleteTableStatementContext

	// IsSqlContext differentiates from other interfaces.
	IsSqlContext()
}

type SqlContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySqlContext() *SqlContext {
	var p = new(SqlContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_sql
	return p
}

func InitEmptySqlContext(p *SqlContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_sql
}

func (*SqlContext) IsSqlContext() {}

func NewSqlContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SqlContext {
	var p = new(SqlContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_sql

	return p
}

func (s *SqlContext) GetParser() antlr.Parser { return s.parser }

func (s *SqlContext) CreateTableStatement() ICreateTableStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateTableStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateTableStatementContext)
}

func (s *SqlContext) InsertTableStatement() IInsertTableStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInsertTableStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInsertTableStatementContext)
}

func (s *SqlContext) SelectTableStatement() ISelectTableStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectTableStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectTableStatementContext)
}

func (s *SqlContext) UpdateTableStatement() IUpdateTableStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUpdateTableStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUpdateTableStatementContext)
}

func (s *SqlContext) DeleteTableStatement() IDeleteTableStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeleteTableStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeleteTableStatementContext)
}

func (s *SqlContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SqlContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SqlContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSql(s)
	}
}

func (s *SqlContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSql(s)
	}
}

func (s *SqlContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSql(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Sql() (localctx ISqlContext) {
	localctx = NewSqlContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, SQLParserRULE_sql)
	p.SetState(51)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserCREATE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(46)
			p.CreateTableStatement()
		}

	case SQLParserINSERT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(47)
			p.InsertTableStatement()
		}

	case SQLParserSELECT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(48)
			p.SelectTableStatement()
		}

	case SQLParserUPDATE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(49)
			p.UpdateTableStatement()
		}

	case SQLParserDELETE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(50)
			p.DeleteTableStatement()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICreateTableStatementContext is an interface to support dynamic dispatch.
type ICreateTableStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CREATE() antlr.TerminalNode
	TABLE() antlr.TerminalNode
	TableName() ITableNameContext
	AllLPAREN() []antlr.TerminalNode
	LPAREN(i int) antlr.TerminalNode
	ColumnDefinitions() IColumnDefinitionsContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	PRIMARY() antlr.TerminalNode
	KEY() antlr.TerminalNode
	ColumnName() IColumnNameContext
	AllRPAREN() []antlr.TerminalNode
	RPAREN(i int) antlr.TerminalNode
	IndexDefinitions() IIndexDefinitionsContext
	SEMICOLON() antlr.TerminalNode

	// IsCreateTableStatementContext differentiates from other interfaces.
	IsCreateTableStatementContext()
}

type CreateTableStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateTableStatementContext() *CreateTableStatementContext {
	var p = new(CreateTableStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_createTableStatement
	return p
}

func InitEmptyCreateTableStatementContext(p *CreateTableStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_createTableStatement
}

func (*CreateTableStatementContext) IsCreateTableStatementContext() {}

func NewCreateTableStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateTableStatementContext {
	var p = new(CreateTableStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_createTableStatement

	return p
}

func (s *CreateTableStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateTableStatementContext) CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserCREATE, 0)
}

func (s *CreateTableStatementContext) TABLE() antlr.TerminalNode {
	return s.GetToken(SQLParserTABLE, 0)
}

func (s *CreateTableStatementContext) TableName() ITableNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableNameContext)
}

func (s *CreateTableStatementContext) AllLPAREN() []antlr.TerminalNode {
	return s.GetTokens(SQLParserLPAREN)
}

func (s *CreateTableStatementContext) LPAREN(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserLPAREN, i)
}

func (s *CreateTableStatementContext) ColumnDefinitions() IColumnDefinitionsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnDefinitionsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnDefinitionsContext)
}

func (s *CreateTableStatementContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *CreateTableStatementContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *CreateTableStatementContext) PRIMARY() antlr.TerminalNode {
	return s.GetToken(SQLParserPRIMARY, 0)
}

func (s *CreateTableStatementContext) KEY() antlr.TerminalNode {
	return s.GetToken(SQLParserKEY, 0)
}

func (s *CreateTableStatementContext) ColumnName() IColumnNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *CreateTableStatementContext) AllRPAREN() []antlr.TerminalNode {
	return s.GetTokens(SQLParserRPAREN)
}

func (s *CreateTableStatementContext) RPAREN(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserRPAREN, i)
}

func (s *CreateTableStatementContext) IndexDefinitions() IIndexDefinitionsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIndexDefinitionsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIndexDefinitionsContext)
}

func (s *CreateTableStatementContext) SEMICOLON() antlr.TerminalNode {
	return s.GetToken(SQLParserSEMICOLON, 0)
}

func (s *CreateTableStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateTableStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CreateTableStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterCreateTableStatement(s)
	}
}

func (s *CreateTableStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitCreateTableStatement(s)
	}
}

func (s *CreateTableStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitCreateTableStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) CreateTableStatement() (localctx ICreateTableStatementContext) {
	localctx = NewCreateTableStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, SQLParserRULE_createTableStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(53)
		p.Match(SQLParserCREATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(54)
		p.Match(SQLParserTABLE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(55)
		p.TableName()
	}
	{
		p.SetState(56)
		p.Match(SQLParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(57)
		p.ColumnDefinitions()
	}
	{
		p.SetState(58)
		p.Match(SQLParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(59)
		p.Match(SQLParserPRIMARY)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(60)
		p.Match(SQLParserKEY)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(61)
		p.Match(SQLParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(62)
		p.ColumnName()
	}
	{
		p.SetState(63)
		p.Match(SQLParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(64)
		p.Match(SQLParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(65)
		p.IndexDefinitions()
	}
	{
		p.SetState(66)
		p.Match(SQLParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(67)
		p.Match(SQLParserSEMICOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITableNameContext is an interface to support dynamic dispatch.
type ITableNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode

	// IsTableNameContext differentiates from other interfaces.
	IsTableNameContext()
}

type TableNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTableNameContext() *TableNameContext {
	var p = new(TableNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_tableName
	return p
}

func InitEmptyTableNameContext(p *TableNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_tableName
}

func (*TableNameContext) IsTableNameContext() {}

func NewTableNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TableNameContext {
	var p = new(TableNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tableName

	return p
}

func (s *TableNameContext) GetParser() antlr.Parser { return s.parser }

func (s *TableNameContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SQLParserIDENTIFIER, 0)
}

func (s *TableNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TableNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TableNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterTableName(s)
	}
}

func (s *TableNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitTableName(s)
	}
}

func (s *TableNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitTableName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TableName() (localctx ITableNameContext) {
	localctx = NewTableNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SQLParserRULE_tableName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		p.Match(SQLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnNameContext is an interface to support dynamic dispatch.
type IColumnNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode

	// IsColumnNameContext differentiates from other interfaces.
	IsColumnNameContext()
}

type ColumnNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnNameContext() *ColumnNameContext {
	var p = new(ColumnNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnName
	return p
}

func InitEmptyColumnNameContext(p *ColumnNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnName
}

func (*ColumnNameContext) IsColumnNameContext() {}

func NewColumnNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnNameContext {
	var p = new(ColumnNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_columnName

	return p
}

func (s *ColumnNameContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnNameContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SQLParserIDENTIFIER, 0)
}

func (s *ColumnNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterColumnName(s)
	}
}

func (s *ColumnNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitColumnName(s)
	}
}

func (s *ColumnNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitColumnName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ColumnName() (localctx IColumnNameContext) {
	localctx = NewColumnNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, SQLParserRULE_columnName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(71)
		p.Match(SQLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnTypeContext is an interface to support dynamic dispatch.
type IColumnTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INT64() antlr.TerminalNode
	BYTES() antlr.TerminalNode

	// IsColumnTypeContext differentiates from other interfaces.
	IsColumnTypeContext()
}

type ColumnTypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnTypeContext() *ColumnTypeContext {
	var p = new(ColumnTypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnType
	return p
}

func InitEmptyColumnTypeContext(p *ColumnTypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnType
}

func (*ColumnTypeContext) IsColumnTypeContext() {}

func NewColumnTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnTypeContext {
	var p = new(ColumnTypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_columnType

	return p
}

func (s *ColumnTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnTypeContext) INT64() antlr.TerminalNode {
	return s.GetToken(SQLParserINT64, 0)
}

func (s *ColumnTypeContext) BYTES() antlr.TerminalNode {
	return s.GetToken(SQLParserBYTES, 0)
}

func (s *ColumnTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterColumnType(s)
	}
}

func (s *ColumnTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitColumnType(s)
	}
}

func (s *ColumnTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitColumnType(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ColumnType() (localctx IColumnTypeContext) {
	localctx = NewColumnTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SQLParserRULE_columnType)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(73)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserINT64 || _la == SQLParserBYTES) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnDefinitionsContext is an interface to support dynamic dispatch.
type IColumnDefinitionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllColumnDefinition() []IColumnDefinitionContext
	ColumnDefinition(i int) IColumnDefinitionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsColumnDefinitionsContext differentiates from other interfaces.
	IsColumnDefinitionsContext()
}

type ColumnDefinitionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnDefinitionsContext() *ColumnDefinitionsContext {
	var p = new(ColumnDefinitionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnDefinitions
	return p
}

func InitEmptyColumnDefinitionsContext(p *ColumnDefinitionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnDefinitions
}

func (*ColumnDefinitionsContext) IsColumnDefinitionsContext() {}

func NewColumnDefinitionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnDefinitionsContext {
	var p = new(ColumnDefinitionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_columnDefinitions

	return p
}

func (s *ColumnDefinitionsContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnDefinitionsContext) AllColumnDefinition() []IColumnDefinitionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IColumnDefinitionContext); ok {
			len++
		}
	}

	tst := make([]IColumnDefinitionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IColumnDefinitionContext); ok {
			tst[i] = t.(IColumnDefinitionContext)
			i++
		}
	}

	return tst
}

func (s *ColumnDefinitionsContext) ColumnDefinition(i int) IColumnDefinitionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnDefinitionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnDefinitionContext)
}

func (s *ColumnDefinitionsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *ColumnDefinitionsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *ColumnDefinitionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnDefinitionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnDefinitionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterColumnDefinitions(s)
	}
}

func (s *ColumnDefinitionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitColumnDefinitions(s)
	}
}

func (s *ColumnDefinitionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitColumnDefinitions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ColumnDefinitions() (localctx IColumnDefinitionsContext) {
	localctx = NewColumnDefinitionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SQLParserRULE_columnDefinitions)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(75)
		p.ColumnDefinition()
	}
	p.SetState(80)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(76)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(77)
				p.ColumnDefinition()
			}

		}
		p.SetState(82)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnDefinitionContext is an interface to support dynamic dispatch.
type IColumnDefinitionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ColumnName() IColumnNameContext
	ColumnType() IColumnTypeContext

	// IsColumnDefinitionContext differentiates from other interfaces.
	IsColumnDefinitionContext()
}

type ColumnDefinitionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnDefinitionContext() *ColumnDefinitionContext {
	var p = new(ColumnDefinitionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnDefinition
	return p
}

func InitEmptyColumnDefinitionContext(p *ColumnDefinitionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnDefinition
}

func (*ColumnDefinitionContext) IsColumnDefinitionContext() {}

func NewColumnDefinitionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnDefinitionContext {
	var p = new(ColumnDefinitionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_columnDefinition

	return p
}

func (s *ColumnDefinitionContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnDefinitionContext) ColumnName() IColumnNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *ColumnDefinitionContext) ColumnType() IColumnTypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnTypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnTypeContext)
}

func (s *ColumnDefinitionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnDefinitionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnDefinitionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterColumnDefinition(s)
	}
}

func (s *ColumnDefinitionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitColumnDefinition(s)
	}
}

func (s *ColumnDefinitionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitColumnDefinition(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ColumnDefinition() (localctx IColumnDefinitionContext) {
	localctx = NewColumnDefinitionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SQLParserRULE_columnDefinition)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(83)
		p.ColumnName()
	}
	{
		p.SetState(84)
		p.ColumnType()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIndexDefinitionsContext is an interface to support dynamic dispatch.
type IIndexDefinitionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIndexDefinition() []IIndexDefinitionContext
	IndexDefinition(i int) IIndexDefinitionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsIndexDefinitionsContext differentiates from other interfaces.
	IsIndexDefinitionsContext()
}

type IndexDefinitionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIndexDefinitionsContext() *IndexDefinitionsContext {
	var p = new(IndexDefinitionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_indexDefinitions
	return p
}

func InitEmptyIndexDefinitionsContext(p *IndexDefinitionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_indexDefinitions
}

func (*IndexDefinitionsContext) IsIndexDefinitionsContext() {}

func NewIndexDefinitionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IndexDefinitionsContext {
	var p = new(IndexDefinitionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_indexDefinitions

	return p
}

func (s *IndexDefinitionsContext) GetParser() antlr.Parser { return s.parser }

func (s *IndexDefinitionsContext) AllIndexDefinition() []IIndexDefinitionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IIndexDefinitionContext); ok {
			len++
		}
	}

	tst := make([]IIndexDefinitionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IIndexDefinitionContext); ok {
			tst[i] = t.(IIndexDefinitionContext)
			i++
		}
	}

	return tst
}

func (s *IndexDefinitionsContext) IndexDefinition(i int) IIndexDefinitionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIndexDefinitionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIndexDefinitionContext)
}

func (s *IndexDefinitionsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *IndexDefinitionsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *IndexDefinitionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IndexDefinitionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IndexDefinitionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterIndexDefinitions(s)
	}
}

func (s *IndexDefinitionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitIndexDefinitions(s)
	}
}

func (s *IndexDefinitionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitIndexDefinitions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) IndexDefinitions() (localctx IIndexDefinitionsContext) {
	localctx = NewIndexDefinitionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SQLParserRULE_indexDefinitions)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(86)
		p.IndexDefinition()
	}
	p.SetState(91)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(87)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(88)
			p.IndexDefinition()
		}

		p.SetState(93)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIndexDefinitionContext is an interface to support dynamic dispatch.
type IIndexDefinitionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INDEX() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AllColumnName() []IColumnNameContext
	ColumnName(i int) IColumnNameContext
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsIndexDefinitionContext differentiates from other interfaces.
	IsIndexDefinitionContext()
}

type IndexDefinitionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIndexDefinitionContext() *IndexDefinitionContext {
	var p = new(IndexDefinitionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_indexDefinition
	return p
}

func InitEmptyIndexDefinitionContext(p *IndexDefinitionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_indexDefinition
}

func (*IndexDefinitionContext) IsIndexDefinitionContext() {}

func NewIndexDefinitionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IndexDefinitionContext {
	var p = new(IndexDefinitionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_indexDefinition

	return p
}

func (s *IndexDefinitionContext) GetParser() antlr.Parser { return s.parser }

func (s *IndexDefinitionContext) INDEX() antlr.TerminalNode {
	return s.GetToken(SQLParserINDEX, 0)
}

func (s *IndexDefinitionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(SQLParserLPAREN, 0)
}

func (s *IndexDefinitionContext) AllColumnName() []IColumnNameContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IColumnNameContext); ok {
			len++
		}
	}

	tst := make([]IColumnNameContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IColumnNameContext); ok {
			tst[i] = t.(IColumnNameContext)
			i++
		}
	}

	return tst
}

func (s *IndexDefinitionContext) ColumnName(i int) IColumnNameContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *IndexDefinitionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(SQLParserRPAREN, 0)
}

func (s *IndexDefinitionContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *IndexDefinitionContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *IndexDefinitionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IndexDefinitionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IndexDefinitionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterIndexDefinition(s)
	}
}

func (s *IndexDefinitionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitIndexDefinition(s)
	}
}

func (s *IndexDefinitionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitIndexDefinition(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) IndexDefinition() (localctx IIndexDefinitionContext) {
	localctx = NewIndexDefinitionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SQLParserRULE_indexDefinition)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(94)
		p.Match(SQLParserINDEX)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(95)
		p.Match(SQLParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(96)
		p.ColumnName()
	}
	p.SetState(101)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(97)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(98)
			p.ColumnName()
		}

		p.SetState(103)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(104)
		p.Match(SQLParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IInsertTableStatementContext is an interface to support dynamic dispatch.
type IInsertTableStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INSERT() antlr.TerminalNode
	INTO() antlr.TerminalNode
	TableName() ITableNameContext
	AllLPAREN() []antlr.TerminalNode
	LPAREN(i int) antlr.TerminalNode
	ColumnInsertNames() IColumnInsertNamesContext
	AllRPAREN() []antlr.TerminalNode
	RPAREN(i int) antlr.TerminalNode
	VALUES() antlr.TerminalNode
	ColumnInsertValues() IColumnInsertValuesContext

	// IsInsertTableStatementContext differentiates from other interfaces.
	IsInsertTableStatementContext()
}

type InsertTableStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInsertTableStatementContext() *InsertTableStatementContext {
	var p = new(InsertTableStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_insertTableStatement
	return p
}

func InitEmptyInsertTableStatementContext(p *InsertTableStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_insertTableStatement
}

func (*InsertTableStatementContext) IsInsertTableStatementContext() {}

func NewInsertTableStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InsertTableStatementContext {
	var p = new(InsertTableStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_insertTableStatement

	return p
}

func (s *InsertTableStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *InsertTableStatementContext) INSERT() antlr.TerminalNode {
	return s.GetToken(SQLParserINSERT, 0)
}

func (s *InsertTableStatementContext) INTO() antlr.TerminalNode {
	return s.GetToken(SQLParserINTO, 0)
}

func (s *InsertTableStatementContext) TableName() ITableNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableNameContext)
}

func (s *InsertTableStatementContext) AllLPAREN() []antlr.TerminalNode {
	return s.GetTokens(SQLParserLPAREN)
}

func (s *InsertTableStatementContext) LPAREN(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserLPAREN, i)
}

func (s *InsertTableStatementContext) ColumnInsertNames() IColumnInsertNamesContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnInsertNamesContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnInsertNamesContext)
}

func (s *InsertTableStatementContext) AllRPAREN() []antlr.TerminalNode {
	return s.GetTokens(SQLParserRPAREN)
}

func (s *InsertTableStatementContext) RPAREN(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserRPAREN, i)
}

func (s *InsertTableStatementContext) VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserVALUES, 0)
}

func (s *InsertTableStatementContext) ColumnInsertValues() IColumnInsertValuesContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnInsertValuesContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnInsertValuesContext)
}

func (s *InsertTableStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InsertTableStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *InsertTableStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterInsertTableStatement(s)
	}
}

func (s *InsertTableStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitInsertTableStatement(s)
	}
}

func (s *InsertTableStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitInsertTableStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) InsertTableStatement() (localctx IInsertTableStatementContext) {
	localctx = NewInsertTableStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SQLParserRULE_insertTableStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(106)
		p.Match(SQLParserINSERT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(107)
		p.Match(SQLParserINTO)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(108)
		p.TableName()
	}
	{
		p.SetState(109)
		p.Match(SQLParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(110)
		p.ColumnInsertNames()
	}
	{
		p.SetState(111)
		p.Match(SQLParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(112)
		p.Match(SQLParserVALUES)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(113)
		p.Match(SQLParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(114)
		p.ColumnInsertValues()
	}
	{
		p.SetState(115)
		p.Match(SQLParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnInsertNamesContext is an interface to support dynamic dispatch.
type IColumnInsertNamesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllColumnName() []IColumnNameContext
	ColumnName(i int) IColumnNameContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsColumnInsertNamesContext differentiates from other interfaces.
	IsColumnInsertNamesContext()
}

type ColumnInsertNamesContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnInsertNamesContext() *ColumnInsertNamesContext {
	var p = new(ColumnInsertNamesContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnInsertNames
	return p
}

func InitEmptyColumnInsertNamesContext(p *ColumnInsertNamesContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnInsertNames
}

func (*ColumnInsertNamesContext) IsColumnInsertNamesContext() {}

func NewColumnInsertNamesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnInsertNamesContext {
	var p = new(ColumnInsertNamesContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_columnInsertNames

	return p
}

func (s *ColumnInsertNamesContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnInsertNamesContext) AllColumnName() []IColumnNameContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IColumnNameContext); ok {
			len++
		}
	}

	tst := make([]IColumnNameContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IColumnNameContext); ok {
			tst[i] = t.(IColumnNameContext)
			i++
		}
	}

	return tst
}

func (s *ColumnInsertNamesContext) ColumnName(i int) IColumnNameContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *ColumnInsertNamesContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *ColumnInsertNamesContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *ColumnInsertNamesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnInsertNamesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnInsertNamesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterColumnInsertNames(s)
	}
}

func (s *ColumnInsertNamesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitColumnInsertNames(s)
	}
}

func (s *ColumnInsertNamesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitColumnInsertNames(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ColumnInsertNames() (localctx IColumnInsertNamesContext) {
	localctx = NewColumnInsertNamesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SQLParserRULE_columnInsertNames)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(117)
		p.ColumnName()
	}
	p.SetState(122)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(118)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(119)
			p.ColumnName()
		}

		p.SetState(124)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnInsertValuesContext is an interface to support dynamic dispatch.
type IColumnInsertValuesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllColumnValue() []IColumnValueContext
	ColumnValue(i int) IColumnValueContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsColumnInsertValuesContext differentiates from other interfaces.
	IsColumnInsertValuesContext()
}

type ColumnInsertValuesContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnInsertValuesContext() *ColumnInsertValuesContext {
	var p = new(ColumnInsertValuesContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnInsertValues
	return p
}

func InitEmptyColumnInsertValuesContext(p *ColumnInsertValuesContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnInsertValues
}

func (*ColumnInsertValuesContext) IsColumnInsertValuesContext() {}

func NewColumnInsertValuesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnInsertValuesContext {
	var p = new(ColumnInsertValuesContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_columnInsertValues

	return p
}

func (s *ColumnInsertValuesContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnInsertValuesContext) AllColumnValue() []IColumnValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IColumnValueContext); ok {
			len++
		}
	}

	tst := make([]IColumnValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IColumnValueContext); ok {
			tst[i] = t.(IColumnValueContext)
			i++
		}
	}

	return tst
}

func (s *ColumnInsertValuesContext) ColumnValue(i int) IColumnValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnValueContext)
}

func (s *ColumnInsertValuesContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *ColumnInsertValuesContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *ColumnInsertValuesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnInsertValuesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnInsertValuesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterColumnInsertValues(s)
	}
}

func (s *ColumnInsertValuesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitColumnInsertValues(s)
	}
}

func (s *ColumnInsertValuesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitColumnInsertValues(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ColumnInsertValues() (localctx IColumnInsertValuesContext) {
	localctx = NewColumnInsertValuesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SQLParserRULE_columnInsertValues)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(125)
		p.ColumnValue()
	}
	p.SetState(130)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(126)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(127)
			p.ColumnValue()
		}

		p.SetState(132)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnValueContext is an interface to support dynamic dispatch.
type IColumnValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INTEGER() antlr.TerminalNode
	STRING() antlr.TerminalNode

	// IsColumnValueContext differentiates from other interfaces.
	IsColumnValueContext()
}

type ColumnValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnValueContext() *ColumnValueContext {
	var p = new(ColumnValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnValue
	return p
}

func InitEmptyColumnValueContext(p *ColumnValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnValue
}

func (*ColumnValueContext) IsColumnValueContext() {}

func NewColumnValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnValueContext {
	var p = new(ColumnValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_columnValue

	return p
}

func (s *ColumnValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnValueContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(SQLParserINTEGER, 0)
}

func (s *ColumnValueContext) STRING() antlr.TerminalNode {
	return s.GetToken(SQLParserSTRING, 0)
}

func (s *ColumnValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterColumnValue(s)
	}
}

func (s *ColumnValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitColumnValue(s)
	}
}

func (s *ColumnValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitColumnValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ColumnValue() (localctx IColumnValueContext) {
	localctx = NewColumnValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, SQLParserRULE_columnValue)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(133)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserINTEGER || _la == SQLParserSTRING) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISelectTableStatementContext is an interface to support dynamic dispatch.
type ISelectTableStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SELECT() antlr.TerminalNode
	SelectColumnNames() ISelectColumnNamesContext
	FROM() antlr.TerminalNode
	TableName() ITableNameContext
	Conditions() IConditionsContext

	// IsSelectTableStatementContext differentiates from other interfaces.
	IsSelectTableStatementContext()
}

type SelectTableStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectTableStatementContext() *SelectTableStatementContext {
	var p = new(SelectTableStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_selectTableStatement
	return p
}

func InitEmptySelectTableStatementContext(p *SelectTableStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_selectTableStatement
}

func (*SelectTableStatementContext) IsSelectTableStatementContext() {}

func NewSelectTableStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectTableStatementContext {
	var p = new(SelectTableStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_selectTableStatement

	return p
}

func (s *SelectTableStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectTableStatementContext) SELECT() antlr.TerminalNode {
	return s.GetToken(SQLParserSELECT, 0)
}

func (s *SelectTableStatementContext) SelectColumnNames() ISelectColumnNamesContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectColumnNamesContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectColumnNamesContext)
}

func (s *SelectTableStatementContext) FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserFROM, 0)
}

func (s *SelectTableStatementContext) TableName() ITableNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableNameContext)
}

func (s *SelectTableStatementContext) Conditions() IConditionsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConditionsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConditionsContext)
}

func (s *SelectTableStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectTableStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectTableStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSelectTableStatement(s)
	}
}

func (s *SelectTableStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSelectTableStatement(s)
	}
}

func (s *SelectTableStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSelectTableStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) SelectTableStatement() (localctx ISelectTableStatementContext) {
	localctx = NewSelectTableStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, SQLParserRULE_selectTableStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(135)
		p.Match(SQLParserSELECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(136)
		p.SelectColumnNames()
	}
	{
		p.SetState(137)
		p.Match(SQLParserFROM)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(138)
		p.TableName()
	}
	{
		p.SetState(139)
		p.Conditions()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConditionsContext is an interface to support dynamic dispatch.
type IConditionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHERE() antlr.TerminalNode
	AllCondition() []IConditionContext
	Condition(i int) IConditionContext
	AllAND() []antlr.TerminalNode
	AND(i int) antlr.TerminalNode

	// IsConditionsContext differentiates from other interfaces.
	IsConditionsContext()
}

type ConditionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConditionsContext() *ConditionsContext {
	var p = new(ConditionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_conditions
	return p
}

func InitEmptyConditionsContext(p *ConditionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_conditions
}

func (*ConditionsContext) IsConditionsContext() {}

func NewConditionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionsContext {
	var p = new(ConditionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_conditions

	return p
}

func (s *ConditionsContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionsContext) WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserWHERE, 0)
}

func (s *ConditionsContext) AllCondition() []IConditionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IConditionContext); ok {
			len++
		}
	}

	tst := make([]IConditionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IConditionContext); ok {
			tst[i] = t.(IConditionContext)
			i++
		}
	}

	return tst
}

func (s *ConditionsContext) Condition(i int) IConditionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConditionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConditionContext)
}

func (s *ConditionsContext) AllAND() []antlr.TerminalNode {
	return s.GetTokens(SQLParserAND)
}

func (s *ConditionsContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserAND, i)
}

func (s *ConditionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConditionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConditionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterConditions(s)
	}
}

func (s *ConditionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitConditions(s)
	}
}

func (s *ConditionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitConditions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Conditions() (localctx IConditionsContext) {
	localctx = NewConditionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, SQLParserRULE_conditions)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(150)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserWHERE {
		{
			p.SetState(141)
			p.Match(SQLParserWHERE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(142)
			p.Condition()
		}
		p.SetState(147)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserAND {
			{
				p.SetState(143)
				p.Match(SQLParserAND)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(144)
				p.Condition()
			}

			p.SetState(149)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISelectColumnNamesContext is an interface to support dynamic dispatch.
type ISelectColumnNamesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STAR() antlr.TerminalNode
	AllColumnName() []IColumnNameContext
	ColumnName(i int) IColumnNameContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsSelectColumnNamesContext differentiates from other interfaces.
	IsSelectColumnNamesContext()
}

type SelectColumnNamesContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectColumnNamesContext() *SelectColumnNamesContext {
	var p = new(SelectColumnNamesContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_selectColumnNames
	return p
}

func InitEmptySelectColumnNamesContext(p *SelectColumnNamesContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_selectColumnNames
}

func (*SelectColumnNamesContext) IsSelectColumnNamesContext() {}

func NewSelectColumnNamesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectColumnNamesContext {
	var p = new(SelectColumnNamesContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_selectColumnNames

	return p
}

func (s *SelectColumnNamesContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectColumnNamesContext) STAR() antlr.TerminalNode {
	return s.GetToken(SQLParserSTAR, 0)
}

func (s *SelectColumnNamesContext) AllColumnName() []IColumnNameContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IColumnNameContext); ok {
			len++
		}
	}

	tst := make([]IColumnNameContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IColumnNameContext); ok {
			tst[i] = t.(IColumnNameContext)
			i++
		}
	}

	return tst
}

func (s *SelectColumnNamesContext) ColumnName(i int) IColumnNameContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *SelectColumnNamesContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *SelectColumnNamesContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *SelectColumnNamesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectColumnNamesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectColumnNamesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSelectColumnNames(s)
	}
}

func (s *SelectColumnNamesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSelectColumnNames(s)
	}
}

func (s *SelectColumnNamesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSelectColumnNames(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) SelectColumnNames() (localctx ISelectColumnNamesContext) {
	localctx = NewSelectColumnNamesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, SQLParserRULE_selectColumnNames)
	var _la int

	p.SetState(161)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserSTAR:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(152)
			p.Match(SQLParserSTAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case SQLParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(153)
			p.ColumnName()
		}
		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(154)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(155)
				p.ColumnName()
			}

			p.SetState(160)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConditionContext is an interface to support dynamic dispatch.
type IConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ComparisonCondition() IComparisonConditionContext
	BetweenCondition() IBetweenConditionContext

	// IsConditionContext differentiates from other interfaces.
	IsConditionContext()
}

type ConditionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConditionContext() *ConditionContext {
	var p = new(ConditionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_condition
	return p
}

func InitEmptyConditionContext(p *ConditionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_condition
}

func (*ConditionContext) IsConditionContext() {}

func NewConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionContext {
	var p = new(ConditionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_condition

	return p
}

func (s *ConditionContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionContext) ComparisonCondition() IComparisonConditionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonConditionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComparisonConditionContext)
}

func (s *ConditionContext) BetweenCondition() IBetweenConditionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBetweenConditionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBetweenConditionContext)
}

func (s *ConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConditionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterCondition(s)
	}
}

func (s *ConditionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitCondition(s)
	}
}

func (s *ConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitCondition(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Condition() (localctx IConditionContext) {
	localctx = NewConditionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, SQLParserRULE_condition)
	p.SetState(165)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(163)
			p.ComparisonCondition()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(164)
			p.BetweenCondition()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IComparisonConditionContext is an interface to support dynamic dispatch.
type IComparisonConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ColumnName() IColumnNameContext
	OP() antlr.TerminalNode
	ColumnValue() IColumnValueContext

	// IsComparisonConditionContext differentiates from other interfaces.
	IsComparisonConditionContext()
}

type ComparisonConditionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonConditionContext() *ComparisonConditionContext {
	var p = new(ComparisonConditionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_comparisonCondition
	return p
}

func InitEmptyComparisonConditionContext(p *ComparisonConditionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_comparisonCondition
}

func (*ComparisonConditionContext) IsComparisonConditionContext() {}

func NewComparisonConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonConditionContext {
	var p = new(ComparisonConditionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_comparisonCondition

	return p
}

func (s *ComparisonConditionContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonConditionContext) ColumnName() IColumnNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *ComparisonConditionContext) OP() antlr.TerminalNode {
	return s.GetToken(SQLParserOP, 0)
}

func (s *ComparisonConditionContext) ColumnValue() IColumnValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnValueContext)
}

func (s *ComparisonConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonConditionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterComparisonCondition(s)
	}
}

func (s *ComparisonConditionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitComparisonCondition(s)
	}
}

func (s *ComparisonConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitComparisonCondition(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ComparisonCondition() (localctx IComparisonConditionContext) {
	localctx = NewComparisonConditionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, SQLParserRULE_comparisonCondition)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(167)
		p.ColumnName()
	}
	{
		p.SetState(168)
		p.Match(SQLParserOP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(169)
		p.ColumnValue()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBetweenConditionContext is an interface to support dynamic dispatch.
type IBetweenConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ColumnName() IColumnNameContext
	BETWEEN() antlr.TerminalNode
	AllColumnValue() []IColumnValueContext
	ColumnValue(i int) IColumnValueContext
	AND() antlr.TerminalNode

	// IsBetweenConditionContext differentiates from other interfaces.
	IsBetweenConditionContext()
}

type BetweenConditionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBetweenConditionContext() *BetweenConditionContext {
	var p = new(BetweenConditionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_betweenCondition
	return p
}

func InitEmptyBetweenConditionContext(p *BetweenConditionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_betweenCondition
}

func (*BetweenConditionContext) IsBetweenConditionContext() {}

func NewBetweenConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BetweenConditionContext {
	var p = new(BetweenConditionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_betweenCondition

	return p
}

func (s *BetweenConditionContext) GetParser() antlr.Parser { return s.parser }

func (s *BetweenConditionContext) ColumnName() IColumnNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *BetweenConditionContext) BETWEEN() antlr.TerminalNode {
	return s.GetToken(SQLParserBETWEEN, 0)
}

func (s *BetweenConditionContext) AllColumnValue() []IColumnValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IColumnValueContext); ok {
			len++
		}
	}

	tst := make([]IColumnValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IColumnValueContext); ok {
			tst[i] = t.(IColumnValueContext)
			i++
		}
	}

	return tst
}

func (s *BetweenConditionContext) ColumnValue(i int) IColumnValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnValueContext)
}

func (s *BetweenConditionContext) AND() antlr.TerminalNode {
	return s.GetToken(SQLParserAND, 0)
}

func (s *BetweenConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BetweenConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BetweenConditionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterBetweenCondition(s)
	}
}

func (s *BetweenConditionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitBetweenCondition(s)
	}
}

func (s *BetweenConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitBetweenCondition(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) BetweenCondition() (localctx IBetweenConditionContext) {
	localctx = NewBetweenConditionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, SQLParserRULE_betweenCondition)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(171)
		p.ColumnName()
	}
	{
		p.SetState(172)
		p.Match(SQLParserBETWEEN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(173)
		p.ColumnValue()
	}
	{
		p.SetState(174)
		p.Match(SQLParserAND)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(175)
		p.ColumnValue()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUpdateTableStatementContext is an interface to support dynamic dispatch.
type IUpdateTableStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UPDATE() antlr.TerminalNode
	TableName() ITableNameContext
	SET() antlr.TerminalNode
	SetClauses() ISetClausesContext
	Conditions() IConditionsContext

	// IsUpdateTableStatementContext differentiates from other interfaces.
	IsUpdateTableStatementContext()
}

type UpdateTableStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUpdateTableStatementContext() *UpdateTableStatementContext {
	var p = new(UpdateTableStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_updateTableStatement
	return p
}

func InitEmptyUpdateTableStatementContext(p *UpdateTableStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_updateTableStatement
}

func (*UpdateTableStatementContext) IsUpdateTableStatementContext() {}

func NewUpdateTableStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UpdateTableStatementContext {
	var p = new(UpdateTableStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_updateTableStatement

	return p
}

func (s *UpdateTableStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *UpdateTableStatementContext) UPDATE() antlr.TerminalNode {
	return s.GetToken(SQLParserUPDATE, 0)
}

func (s *UpdateTableStatementContext) TableName() ITableNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableNameContext)
}

func (s *UpdateTableStatementContext) SET() antlr.TerminalNode {
	return s.GetToken(SQLParserSET, 0)
}

func (s *UpdateTableStatementContext) SetClauses() ISetClausesContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISetClausesContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISetClausesContext)
}

func (s *UpdateTableStatementContext) Conditions() IConditionsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConditionsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConditionsContext)
}

func (s *UpdateTableStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UpdateTableStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UpdateTableStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterUpdateTableStatement(s)
	}
}

func (s *UpdateTableStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitUpdateTableStatement(s)
	}
}

func (s *UpdateTableStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitUpdateTableStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) UpdateTableStatement() (localctx IUpdateTableStatementContext) {
	localctx = NewUpdateTableStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, SQLParserRULE_updateTableStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(177)
		p.Match(SQLParserUPDATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(178)
		p.TableName()
	}
	{
		p.SetState(179)
		p.Match(SQLParserSET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(180)
		p.SetClauses()
	}
	{
		p.SetState(181)
		p.Conditions()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISetClausesContext is an interface to support dynamic dispatch.
type ISetClausesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSetClause() []ISetClauseContext
	SetClause(i int) ISetClauseContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsSetClausesContext differentiates from other interfaces.
	IsSetClausesContext()
}

type SetClausesContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySetClausesContext() *SetClausesContext {
	var p = new(SetClausesContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_setClauses
	return p
}

func InitEmptySetClausesContext(p *SetClausesContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_setClauses
}

func (*SetClausesContext) IsSetClausesContext() {}

func NewSetClausesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetClausesContext {
	var p = new(SetClausesContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_setClauses

	return p
}

func (s *SetClausesContext) GetParser() antlr.Parser { return s.parser }

func (s *SetClausesContext) AllSetClause() []ISetClauseContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISetClauseContext); ok {
			len++
		}
	}

	tst := make([]ISetClauseContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISetClauseContext); ok {
			tst[i] = t.(ISetClauseContext)
			i++
		}
	}

	return tst
}

func (s *SetClausesContext) SetClause(i int) ISetClauseContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISetClauseContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISetClauseContext)
}

func (s *SetClausesContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *SetClausesContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *SetClausesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetClausesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SetClausesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSetClauses(s)
	}
}

func (s *SetClausesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSetClauses(s)
	}
}

func (s *SetClausesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSetClauses(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) SetClauses() (localctx ISetClausesContext) {
	localctx = NewSetClausesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, SQLParserRULE_setClauses)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(183)
		p.SetClause()
	}
	p.SetState(188)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(184)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(185)
			p.SetClause()
		}

		p.SetState(190)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISetClauseContext is an interface to support dynamic dispatch.
type ISetClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ColumnName() IColumnNameContext
	OP() antlr.TerminalNode
	ColumnValue() IColumnValueContext

	// IsSetClauseContext differentiates from other interfaces.
	IsSetClauseContext()
}

type SetClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySetClauseContext() *SetClauseContext {
	var p = new(SetClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_setClause
	return p
}

func InitEmptySetClauseContext(p *SetClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_setClause
}

func (*SetClauseContext) IsSetClauseContext() {}

func NewSetClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetClauseContext {
	var p = new(SetClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_setClause

	return p
}

func (s *SetClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *SetClauseContext) ColumnName() IColumnNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *SetClauseContext) OP() antlr.TerminalNode {
	return s.GetToken(SQLParserOP, 0)
}

func (s *SetClauseContext) ColumnValue() IColumnValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnValueContext)
}

func (s *SetClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SetClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSetClause(s)
	}
}

func (s *SetClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSetClause(s)
	}
}

func (s *SetClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSetClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) SetClause() (localctx ISetClauseContext) {
	localctx = NewSetClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, SQLParserRULE_setClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(191)
		p.ColumnName()
	}
	{
		p.SetState(192)
		p.Match(SQLParserOP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(193)
		p.ColumnValue()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDeleteTableStatementContext is an interface to support dynamic dispatch.
type IDeleteTableStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DELETE() antlr.TerminalNode
	FROM() antlr.TerminalNode
	TableName() ITableNameContext
	Conditions() IConditionsContext

	// IsDeleteTableStatementContext differentiates from other interfaces.
	IsDeleteTableStatementContext()
}

type DeleteTableStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDeleteTableStatementContext() *DeleteTableStatementContext {
	var p = new(DeleteTableStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_deleteTableStatement
	return p
}

func InitEmptyDeleteTableStatementContext(p *DeleteTableStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_deleteTableStatement
}

func (*DeleteTableStatementContext) IsDeleteTableStatementContext() {}

func NewDeleteTableStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeleteTableStatementContext {
	var p = new(DeleteTableStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_deleteTableStatement

	return p
}

func (s *DeleteTableStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *DeleteTableStatementContext) DELETE() antlr.TerminalNode {
	return s.GetToken(SQLParserDELETE, 0)
}

func (s *DeleteTableStatementContext) FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserFROM, 0)
}

func (s *DeleteTableStatementContext) TableName() ITableNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableNameContext)
}

func (s *DeleteTableStatementContext) Conditions() IConditionsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConditionsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConditionsContext)
}

func (s *DeleteTableStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeleteTableStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeleteTableStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterDeleteTableStatement(s)
	}
}

func (s *DeleteTableStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitDeleteTableStatement(s)
	}
}

func (s *DeleteTableStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitDeleteTableStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) DeleteTableStatement() (localctx IDeleteTableStatementContext) {
	localctx = NewDeleteTableStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, SQLParserRULE_deleteTableStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(195)
		p.Match(SQLParserDELETE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(196)
		p.Match(SQLParserFROM)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(197)
		p.TableName()
	}
	{
		p.SetState(198)
		p.Conditions()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
