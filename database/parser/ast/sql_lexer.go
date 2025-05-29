// Code generated from SQLLexer.g4 by ANTLR 4.13.2. DO NOT EDIT.

package ast

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type SQLLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var SQLLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func sqllexerLexerInit() {
	staticData := &SQLLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
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
		"CREATE", "TABLE", "PRIMARY", "KEY", "INDEX", "INSERT", "INTO", "VALUES",
		"INT64", "BYTES", "IDENTIFIER", "LPAREN", "RPAREN", "COMMA", "SEMICOLON",
		"WS", "INTEGER", "STRING",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 18, 137, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
		1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4,
		1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6,
		1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8,
		1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 5, 10, 102,
		8, 10, 10, 10, 12, 10, 105, 9, 10, 1, 11, 1, 11, 1, 12, 1, 12, 1, 13, 1,
		13, 1, 14, 1, 14, 1, 15, 4, 15, 116, 8, 15, 11, 15, 12, 15, 117, 1, 15,
		1, 15, 1, 16, 4, 16, 123, 8, 16, 11, 16, 12, 16, 124, 1, 17, 1, 17, 1,
		17, 1, 17, 5, 17, 131, 8, 17, 10, 17, 12, 17, 134, 9, 17, 1, 17, 1, 17,
		0, 0, 18, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19,
		10, 21, 11, 23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17, 35, 18, 1,
		0, 5, 3, 0, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97,
		122, 3, 0, 9, 10, 13, 13, 32, 32, 1, 0, 48, 57, 1, 0, 39, 39, 141, 0, 1,
		1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9,
		1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0,
		17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0,
		0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0,
		0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0, 0, 1, 37, 1, 0, 0, 0, 3, 44, 1, 0,
		0, 0, 5, 50, 1, 0, 0, 0, 7, 58, 1, 0, 0, 0, 9, 62, 1, 0, 0, 0, 11, 68,
		1, 0, 0, 0, 13, 75, 1, 0, 0, 0, 15, 80, 1, 0, 0, 0, 17, 87, 1, 0, 0, 0,
		19, 93, 1, 0, 0, 0, 21, 99, 1, 0, 0, 0, 23, 106, 1, 0, 0, 0, 25, 108, 1,
		0, 0, 0, 27, 110, 1, 0, 0, 0, 29, 112, 1, 0, 0, 0, 31, 115, 1, 0, 0, 0,
		33, 122, 1, 0, 0, 0, 35, 126, 1, 0, 0, 0, 37, 38, 5, 67, 0, 0, 38, 39,
		5, 82, 0, 0, 39, 40, 5, 69, 0, 0, 40, 41, 5, 65, 0, 0, 41, 42, 5, 84, 0,
		0, 42, 43, 5, 69, 0, 0, 43, 2, 1, 0, 0, 0, 44, 45, 5, 84, 0, 0, 45, 46,
		5, 65, 0, 0, 46, 47, 5, 66, 0, 0, 47, 48, 5, 76, 0, 0, 48, 49, 5, 69, 0,
		0, 49, 4, 1, 0, 0, 0, 50, 51, 5, 80, 0, 0, 51, 52, 5, 82, 0, 0, 52, 53,
		5, 73, 0, 0, 53, 54, 5, 77, 0, 0, 54, 55, 5, 65, 0, 0, 55, 56, 5, 82, 0,
		0, 56, 57, 5, 89, 0, 0, 57, 6, 1, 0, 0, 0, 58, 59, 5, 75, 0, 0, 59, 60,
		5, 69, 0, 0, 60, 61, 5, 89, 0, 0, 61, 8, 1, 0, 0, 0, 62, 63, 5, 73, 0,
		0, 63, 64, 5, 78, 0, 0, 64, 65, 5, 68, 0, 0, 65, 66, 5, 69, 0, 0, 66, 67,
		5, 88, 0, 0, 67, 10, 1, 0, 0, 0, 68, 69, 5, 73, 0, 0, 69, 70, 5, 78, 0,
		0, 70, 71, 5, 83, 0, 0, 71, 72, 5, 69, 0, 0, 72, 73, 5, 82, 0, 0, 73, 74,
		5, 84, 0, 0, 74, 12, 1, 0, 0, 0, 75, 76, 5, 73, 0, 0, 76, 77, 5, 78, 0,
		0, 77, 78, 5, 84, 0, 0, 78, 79, 5, 79, 0, 0, 79, 14, 1, 0, 0, 0, 80, 81,
		5, 86, 0, 0, 81, 82, 5, 65, 0, 0, 82, 83, 5, 76, 0, 0, 83, 84, 5, 85, 0,
		0, 84, 85, 5, 69, 0, 0, 85, 86, 5, 83, 0, 0, 86, 16, 1, 0, 0, 0, 87, 88,
		5, 73, 0, 0, 88, 89, 5, 78, 0, 0, 89, 90, 5, 84, 0, 0, 90, 91, 5, 54, 0,
		0, 91, 92, 5, 52, 0, 0, 92, 18, 1, 0, 0, 0, 93, 94, 5, 66, 0, 0, 94, 95,
		5, 89, 0, 0, 95, 96, 5, 84, 0, 0, 96, 97, 5, 69, 0, 0, 97, 98, 5, 83, 0,
		0, 98, 20, 1, 0, 0, 0, 99, 103, 7, 0, 0, 0, 100, 102, 7, 1, 0, 0, 101,
		100, 1, 0, 0, 0, 102, 105, 1, 0, 0, 0, 103, 101, 1, 0, 0, 0, 103, 104,
		1, 0, 0, 0, 104, 22, 1, 0, 0, 0, 105, 103, 1, 0, 0, 0, 106, 107, 5, 40,
		0, 0, 107, 24, 1, 0, 0, 0, 108, 109, 5, 41, 0, 0, 109, 26, 1, 0, 0, 0,
		110, 111, 5, 44, 0, 0, 111, 28, 1, 0, 0, 0, 112, 113, 5, 59, 0, 0, 113,
		30, 1, 0, 0, 0, 114, 116, 7, 2, 0, 0, 115, 114, 1, 0, 0, 0, 116, 117, 1,
		0, 0, 0, 117, 115, 1, 0, 0, 0, 117, 118, 1, 0, 0, 0, 118, 119, 1, 0, 0,
		0, 119, 120, 6, 15, 0, 0, 120, 32, 1, 0, 0, 0, 121, 123, 7, 3, 0, 0, 122,
		121, 1, 0, 0, 0, 123, 124, 1, 0, 0, 0, 124, 122, 1, 0, 0, 0, 124, 125,
		1, 0, 0, 0, 125, 34, 1, 0, 0, 0, 126, 132, 5, 39, 0, 0, 127, 131, 8, 4,
		0, 0, 128, 129, 5, 39, 0, 0, 129, 131, 5, 39, 0, 0, 130, 127, 1, 0, 0,
		0, 130, 128, 1, 0, 0, 0, 131, 134, 1, 0, 0, 0, 132, 130, 1, 0, 0, 0, 132,
		133, 1, 0, 0, 0, 133, 135, 1, 0, 0, 0, 134, 132, 1, 0, 0, 0, 135, 136,
		5, 39, 0, 0, 136, 36, 1, 0, 0, 0, 6, 0, 103, 117, 124, 130, 132, 1, 6,
		0, 0,
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

// SQLLexerInit initializes any static state used to implement SQLLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewSQLLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func SQLLexerInit() {
	staticData := &SQLLexerLexerStaticData
	staticData.once.Do(sqllexerLexerInit)
}

// NewSQLLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewSQLLexer(input antlr.CharStream) *SQLLexer {
	SQLLexerInit()
	l := new(SQLLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &SQLLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "SQLLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// SQLLexer tokens.
const (
	SQLLexerCREATE     = 1
	SQLLexerTABLE      = 2
	SQLLexerPRIMARY    = 3
	SQLLexerKEY        = 4
	SQLLexerINDEX      = 5
	SQLLexerINSERT     = 6
	SQLLexerINTO       = 7
	SQLLexerVALUES     = 8
	SQLLexerINT64      = 9
	SQLLexerBYTES      = 10
	SQLLexerIDENTIFIER = 11
	SQLLexerLPAREN     = 12
	SQLLexerRPAREN     = 13
	SQLLexerCOMMA      = 14
	SQLLexerSEMICOLON  = 15
	SQLLexerWS         = 16
	SQLLexerINTEGER    = 17
	SQLLexerSTRING     = 18
)
