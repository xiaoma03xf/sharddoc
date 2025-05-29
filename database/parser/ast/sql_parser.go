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
		"'INTO'", "'VALUES'", "'INT64'", "'BYTES'", "", "'('", "')'", "','",
		"';'",
	}
	staticData.SymbolicNames = []string{
		"", "CREATE", "TABLE", "PRIMARY", "KEY", "INDEX", "INSERT", "INTO",
		"VALUES", "INT64", "BYTES", "IDENTIFIER", "LPAREN", "RPAREN", "COMMA",
		"SEMICOLON", "WS", "INTEGER", "STRING",
	}
	staticData.RuleNames = []string{
		"sql", "createTableStatement", "tableName", "columnName", "columnType",
		"columnDefinitions", "columnDefinition", "indexDefinitions", "indexDefinition",
		"insertTableStatement", "columnInsertNames", "columnInsertValues", "columnValue",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 18, 113, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 1, 0, 1, 0, 3, 0, 29, 8, 0, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 5, 5,
		56, 8, 5, 10, 5, 12, 5, 59, 9, 5, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 5,
		7, 67, 8, 7, 10, 7, 12, 7, 70, 9, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 5, 8,
		77, 8, 8, 10, 8, 12, 8, 80, 9, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1,
		9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 5, 10, 98,
		8, 10, 10, 10, 12, 10, 101, 9, 10, 1, 11, 1, 11, 1, 11, 5, 11, 106, 8,
		11, 10, 11, 12, 11, 109, 9, 11, 1, 12, 1, 12, 1, 12, 0, 0, 13, 0, 2, 4,
		6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 0, 2, 1, 0, 9, 10, 1, 0, 17, 18,
		105, 0, 28, 1, 0, 0, 0, 2, 30, 1, 0, 0, 0, 4, 46, 1, 0, 0, 0, 6, 48, 1,
		0, 0, 0, 8, 50, 1, 0, 0, 0, 10, 52, 1, 0, 0, 0, 12, 60, 1, 0, 0, 0, 14,
		63, 1, 0, 0, 0, 16, 71, 1, 0, 0, 0, 18, 83, 1, 0, 0, 0, 20, 94, 1, 0, 0,
		0, 22, 102, 1, 0, 0, 0, 24, 110, 1, 0, 0, 0, 26, 29, 3, 2, 1, 0, 27, 29,
		3, 18, 9, 0, 28, 26, 1, 0, 0, 0, 28, 27, 1, 0, 0, 0, 29, 1, 1, 0, 0, 0,
		30, 31, 5, 1, 0, 0, 31, 32, 5, 2, 0, 0, 32, 33, 3, 4, 2, 0, 33, 34, 5,
		12, 0, 0, 34, 35, 3, 10, 5, 0, 35, 36, 5, 14, 0, 0, 36, 37, 5, 3, 0, 0,
		37, 38, 5, 4, 0, 0, 38, 39, 5, 12, 0, 0, 39, 40, 3, 6, 3, 0, 40, 41, 5,
		13, 0, 0, 41, 42, 5, 14, 0, 0, 42, 43, 3, 14, 7, 0, 43, 44, 5, 13, 0, 0,
		44, 45, 5, 15, 0, 0, 45, 3, 1, 0, 0, 0, 46, 47, 5, 11, 0, 0, 47, 5, 1,
		0, 0, 0, 48, 49, 5, 11, 0, 0, 49, 7, 1, 0, 0, 0, 50, 51, 7, 0, 0, 0, 51,
		9, 1, 0, 0, 0, 52, 57, 3, 12, 6, 0, 53, 54, 5, 14, 0, 0, 54, 56, 3, 12,
		6, 0, 55, 53, 1, 0, 0, 0, 56, 59, 1, 0, 0, 0, 57, 55, 1, 0, 0, 0, 57, 58,
		1, 0, 0, 0, 58, 11, 1, 0, 0, 0, 59, 57, 1, 0, 0, 0, 60, 61, 3, 6, 3, 0,
		61, 62, 3, 8, 4, 0, 62, 13, 1, 0, 0, 0, 63, 68, 3, 16, 8, 0, 64, 65, 5,
		14, 0, 0, 65, 67, 3, 16, 8, 0, 66, 64, 1, 0, 0, 0, 67, 70, 1, 0, 0, 0,
		68, 66, 1, 0, 0, 0, 68, 69, 1, 0, 0, 0, 69, 15, 1, 0, 0, 0, 70, 68, 1,
		0, 0, 0, 71, 72, 5, 5, 0, 0, 72, 73, 5, 12, 0, 0, 73, 78, 3, 6, 3, 0, 74,
		75, 5, 14, 0, 0, 75, 77, 3, 6, 3, 0, 76, 74, 1, 0, 0, 0, 77, 80, 1, 0,
		0, 0, 78, 76, 1, 0, 0, 0, 78, 79, 1, 0, 0, 0, 79, 81, 1, 0, 0, 0, 80, 78,
		1, 0, 0, 0, 81, 82, 5, 13, 0, 0, 82, 17, 1, 0, 0, 0, 83, 84, 5, 6, 0, 0,
		84, 85, 5, 7, 0, 0, 85, 86, 3, 4, 2, 0, 86, 87, 5, 12, 0, 0, 87, 88, 3,
		20, 10, 0, 88, 89, 5, 13, 0, 0, 89, 90, 5, 8, 0, 0, 90, 91, 5, 12, 0, 0,
		91, 92, 3, 22, 11, 0, 92, 93, 5, 13, 0, 0, 93, 19, 1, 0, 0, 0, 94, 99,
		3, 6, 3, 0, 95, 96, 5, 14, 0, 0, 96, 98, 3, 6, 3, 0, 97, 95, 1, 0, 0, 0,
		98, 101, 1, 0, 0, 0, 99, 97, 1, 0, 0, 0, 99, 100, 1, 0, 0, 0, 100, 21,
		1, 0, 0, 0, 101, 99, 1, 0, 0, 0, 102, 107, 3, 24, 12, 0, 103, 104, 5, 14,
		0, 0, 104, 106, 3, 24, 12, 0, 105, 103, 1, 0, 0, 0, 106, 109, 1, 0, 0,
		0, 107, 105, 1, 0, 0, 0, 107, 108, 1, 0, 0, 0, 108, 23, 1, 0, 0, 0, 109,
		107, 1, 0, 0, 0, 110, 111, 7, 1, 0, 0, 111, 25, 1, 0, 0, 0, 6, 28, 57,
		68, 78, 99, 107,
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
	SQLParserINT64      = 9
	SQLParserBYTES      = 10
	SQLParserIDENTIFIER = 11
	SQLParserLPAREN     = 12
	SQLParserRPAREN     = 13
	SQLParserCOMMA      = 14
	SQLParserSEMICOLON  = 15
	SQLParserWS         = 16
	SQLParserINTEGER    = 17
	SQLParserSTRING     = 18
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
)

// ISqlContext is an interface to support dynamic dispatch.
type ISqlContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CreateTableStatement() ICreateTableStatementContext
	InsertTableStatement() IInsertTableStatementContext

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
	p.SetState(28)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserCREATE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(26)
			p.CreateTableStatement()
		}

	case SQLParserINSERT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(27)
			p.InsertTableStatement()
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
		p.SetState(30)
		p.Match(SQLParserCREATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(31)
		p.Match(SQLParserTABLE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(32)
		p.TableName()
	}
	{
		p.SetState(33)
		p.Match(SQLParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(34)
		p.ColumnDefinitions()
	}
	{
		p.SetState(35)
		p.Match(SQLParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(36)
		p.Match(SQLParserPRIMARY)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(37)
		p.Match(SQLParserKEY)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(38)
		p.Match(SQLParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(39)
		p.ColumnName()
	}
	{
		p.SetState(40)
		p.Match(SQLParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(41)
		p.Match(SQLParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(42)
		p.IndexDefinitions()
	}
	{
		p.SetState(43)
		p.Match(SQLParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(44)
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
		p.SetState(46)
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
		p.SetState(48)
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
		p.SetState(50)
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
		p.SetState(52)
		p.ColumnDefinition()
	}
	p.SetState(57)
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
				p.SetState(53)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(54)
				p.ColumnDefinition()
			}

		}
		p.SetState(59)
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
		p.SetState(60)
		p.ColumnName()
	}
	{
		p.SetState(61)
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
		p.SetState(63)
		p.IndexDefinition()
	}
	p.SetState(68)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
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
			p.IndexDefinition()
		}

		p.SetState(70)
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
		p.SetState(71)
		p.Match(SQLParserINDEX)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(72)
		p.Match(SQLParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(73)
		p.ColumnName()
	}
	p.SetState(78)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(74)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(75)
			p.ColumnName()
		}

		p.SetState(80)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(81)
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
		p.SetState(83)
		p.Match(SQLParserINSERT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(84)
		p.Match(SQLParserINTO)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(85)
		p.TableName()
	}
	{
		p.SetState(86)
		p.Match(SQLParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(87)
		p.ColumnInsertNames()
	}
	{
		p.SetState(88)
		p.Match(SQLParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(89)
		p.Match(SQLParserVALUES)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(90)
		p.Match(SQLParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(91)
		p.ColumnInsertValues()
	}
	{
		p.SetState(92)
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
		p.SetState(94)
		p.ColumnName()
	}
	p.SetState(99)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(95)
			p.Match(SQLParserCOMMA)
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
		p.SetState(102)
		p.ColumnValue()
	}
	p.SetState(107)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(103)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(104)
			p.ColumnValue()
		}

		p.SetState(109)
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
		p.SetState(110)
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
