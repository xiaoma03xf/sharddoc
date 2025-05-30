// Generated from /home/dengjun/sharddoc/database/parser/SQLParser.g4 by ANTLR 4.13.1
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast", "CheckReturnValue"})
public class SQLParser extends Parser {
	static { RuntimeMetaData.checkVersion("4.13.1", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		CREATE=1, TABLE=2, PRIMARY=3, KEY=4, INDEX=5, INSERT=6, INTO=7, VALUES=8, 
		SELECT=9, FROM=10, WHERE=11, AND=12, INT64=13, BYTES=14, BETWEEN=15, IDENTIFIER=16, 
		LPAREN=17, RPAREN=18, COMMA=19, SEMICOLON=20, STAR=21, WS=22, INTEGER=23, 
		STRING=24, OP=25;
	public static final int
		RULE_sql = 0, RULE_createTableStatement = 1, RULE_tableName = 2, RULE_columnName = 3, 
		RULE_columnType = 4, RULE_columnDefinitions = 5, RULE_columnDefinition = 6, 
		RULE_indexDefinitions = 7, RULE_indexDefinition = 8, RULE_insertTableStatement = 9, 
		RULE_columnInsertNames = 10, RULE_columnInsertValues = 11, RULE_columnValue = 12, 
		RULE_selectTableStatement = 13, RULE_selectColumnNames = 14, RULE_condition = 15, 
		RULE_comparisonCondition = 16, RULE_betweenCondition = 17;
	private static String[] makeRuleNames() {
		return new String[] {
			"sql", "createTableStatement", "tableName", "columnName", "columnType", 
			"columnDefinitions", "columnDefinition", "indexDefinitions", "indexDefinition", 
			"insertTableStatement", "columnInsertNames", "columnInsertValues", "columnValue", 
			"selectTableStatement", "selectColumnNames", "condition", "comparisonCondition", 
			"betweenCondition"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'CREATE'", "'TABLE'", "'PRIMARY'", "'KEY'", "'INDEX'", "'INSERT'", 
			"'INTO'", "'VALUES'", "'SELECT'", "'FROM'", "'WHERE'", "'AND'", "'INT64'", 
			"'BYTES'", "'BETWEEN'", null, "'('", "')'", "','", "';'", "'*'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, "CREATE", "TABLE", "PRIMARY", "KEY", "INDEX", "INSERT", "INTO", 
			"VALUES", "SELECT", "FROM", "WHERE", "AND", "INT64", "BYTES", "BETWEEN", 
			"IDENTIFIER", "LPAREN", "RPAREN", "COMMA", "SEMICOLON", "STAR", "WS", 
			"INTEGER", "STRING", "OP"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}

	@Override
	public String getGrammarFileName() { return "SQLParser.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public SQLParser(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@SuppressWarnings("CheckReturnValue")
	public static class SqlContext extends ParserRuleContext {
		public CreateTableStatementContext createTableStatement() {
			return getRuleContext(CreateTableStatementContext.class,0);
		}
		public InsertTableStatementContext insertTableStatement() {
			return getRuleContext(InsertTableStatementContext.class,0);
		}
		public SelectTableStatementContext selectTableStatement() {
			return getRuleContext(SelectTableStatementContext.class,0);
		}
		public SqlContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_sql; }
	}

	public final SqlContext sql() throws RecognitionException {
		SqlContext _localctx = new SqlContext(_ctx, getState());
		enterRule(_localctx, 0, RULE_sql);
		try {
			setState(39);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case CREATE:
				enterOuterAlt(_localctx, 1);
				{
				setState(36);
				createTableStatement();
				}
				break;
			case INSERT:
				enterOuterAlt(_localctx, 2);
				{
				setState(37);
				insertTableStatement();
				}
				break;
			case SELECT:
				enterOuterAlt(_localctx, 3);
				{
				setState(38);
				selectTableStatement();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class CreateTableStatementContext extends ParserRuleContext {
		public TerminalNode CREATE() { return getToken(SQLParser.CREATE, 0); }
		public TerminalNode TABLE() { return getToken(SQLParser.TABLE, 0); }
		public TableNameContext tableName() {
			return getRuleContext(TableNameContext.class,0);
		}
		public List<TerminalNode> LPAREN() { return getTokens(SQLParser.LPAREN); }
		public TerminalNode LPAREN(int i) {
			return getToken(SQLParser.LPAREN, i);
		}
		public ColumnDefinitionsContext columnDefinitions() {
			return getRuleContext(ColumnDefinitionsContext.class,0);
		}
		public List<TerminalNode> COMMA() { return getTokens(SQLParser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(SQLParser.COMMA, i);
		}
		public TerminalNode PRIMARY() { return getToken(SQLParser.PRIMARY, 0); }
		public TerminalNode KEY() { return getToken(SQLParser.KEY, 0); }
		public ColumnNameContext columnName() {
			return getRuleContext(ColumnNameContext.class,0);
		}
		public List<TerminalNode> RPAREN() { return getTokens(SQLParser.RPAREN); }
		public TerminalNode RPAREN(int i) {
			return getToken(SQLParser.RPAREN, i);
		}
		public IndexDefinitionsContext indexDefinitions() {
			return getRuleContext(IndexDefinitionsContext.class,0);
		}
		public TerminalNode SEMICOLON() { return getToken(SQLParser.SEMICOLON, 0); }
		public CreateTableStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_createTableStatement; }
	}

	public final CreateTableStatementContext createTableStatement() throws RecognitionException {
		CreateTableStatementContext _localctx = new CreateTableStatementContext(_ctx, getState());
		enterRule(_localctx, 2, RULE_createTableStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(41);
			match(CREATE);
			setState(42);
			match(TABLE);
			setState(43);
			tableName();
			setState(44);
			match(LPAREN);
			setState(45);
			columnDefinitions();
			setState(46);
			match(COMMA);
			setState(47);
			match(PRIMARY);
			setState(48);
			match(KEY);
			setState(49);
			match(LPAREN);
			setState(50);
			columnName();
			setState(51);
			match(RPAREN);
			setState(52);
			match(COMMA);
			setState(53);
			indexDefinitions();
			setState(54);
			match(RPAREN);
			setState(55);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class TableNameContext extends ParserRuleContext {
		public TerminalNode IDENTIFIER() { return getToken(SQLParser.IDENTIFIER, 0); }
		public TableNameContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_tableName; }
	}

	public final TableNameContext tableName() throws RecognitionException {
		TableNameContext _localctx = new TableNameContext(_ctx, getState());
		enterRule(_localctx, 4, RULE_tableName);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(57);
			match(IDENTIFIER);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ColumnNameContext extends ParserRuleContext {
		public TerminalNode IDENTIFIER() { return getToken(SQLParser.IDENTIFIER, 0); }
		public ColumnNameContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_columnName; }
	}

	public final ColumnNameContext columnName() throws RecognitionException {
		ColumnNameContext _localctx = new ColumnNameContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_columnName);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(59);
			match(IDENTIFIER);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ColumnTypeContext extends ParserRuleContext {
		public TerminalNode INT64() { return getToken(SQLParser.INT64, 0); }
		public TerminalNode BYTES() { return getToken(SQLParser.BYTES, 0); }
		public ColumnTypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_columnType; }
	}

	public final ColumnTypeContext columnType() throws RecognitionException {
		ColumnTypeContext _localctx = new ColumnTypeContext(_ctx, getState());
		enterRule(_localctx, 8, RULE_columnType);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(61);
			_la = _input.LA(1);
			if ( !(_la==INT64 || _la==BYTES) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ColumnDefinitionsContext extends ParserRuleContext {
		public List<ColumnDefinitionContext> columnDefinition() {
			return getRuleContexts(ColumnDefinitionContext.class);
		}
		public ColumnDefinitionContext columnDefinition(int i) {
			return getRuleContext(ColumnDefinitionContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(SQLParser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(SQLParser.COMMA, i);
		}
		public ColumnDefinitionsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_columnDefinitions; }
	}

	public final ColumnDefinitionsContext columnDefinitions() throws RecognitionException {
		ColumnDefinitionsContext _localctx = new ColumnDefinitionsContext(_ctx, getState());
		enterRule(_localctx, 10, RULE_columnDefinitions);
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(63);
			columnDefinition();
			setState(68);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,1,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(64);
					match(COMMA);
					setState(65);
					columnDefinition();
					}
					} 
				}
				setState(70);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,1,_ctx);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ColumnDefinitionContext extends ParserRuleContext {
		public ColumnNameContext columnName() {
			return getRuleContext(ColumnNameContext.class,0);
		}
		public ColumnTypeContext columnType() {
			return getRuleContext(ColumnTypeContext.class,0);
		}
		public ColumnDefinitionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_columnDefinition; }
	}

	public final ColumnDefinitionContext columnDefinition() throws RecognitionException {
		ColumnDefinitionContext _localctx = new ColumnDefinitionContext(_ctx, getState());
		enterRule(_localctx, 12, RULE_columnDefinition);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(71);
			columnName();
			setState(72);
			columnType();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class IndexDefinitionsContext extends ParserRuleContext {
		public List<IndexDefinitionContext> indexDefinition() {
			return getRuleContexts(IndexDefinitionContext.class);
		}
		public IndexDefinitionContext indexDefinition(int i) {
			return getRuleContext(IndexDefinitionContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(SQLParser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(SQLParser.COMMA, i);
		}
		public IndexDefinitionsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_indexDefinitions; }
	}

	public final IndexDefinitionsContext indexDefinitions() throws RecognitionException {
		IndexDefinitionsContext _localctx = new IndexDefinitionsContext(_ctx, getState());
		enterRule(_localctx, 14, RULE_indexDefinitions);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(74);
			indexDefinition();
			setState(79);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(75);
				match(COMMA);
				setState(76);
				indexDefinition();
				}
				}
				setState(81);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class IndexDefinitionContext extends ParserRuleContext {
		public TerminalNode INDEX() { return getToken(SQLParser.INDEX, 0); }
		public TerminalNode LPAREN() { return getToken(SQLParser.LPAREN, 0); }
		public List<ColumnNameContext> columnName() {
			return getRuleContexts(ColumnNameContext.class);
		}
		public ColumnNameContext columnName(int i) {
			return getRuleContext(ColumnNameContext.class,i);
		}
		public TerminalNode RPAREN() { return getToken(SQLParser.RPAREN, 0); }
		public List<TerminalNode> COMMA() { return getTokens(SQLParser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(SQLParser.COMMA, i);
		}
		public IndexDefinitionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_indexDefinition; }
	}

	public final IndexDefinitionContext indexDefinition() throws RecognitionException {
		IndexDefinitionContext _localctx = new IndexDefinitionContext(_ctx, getState());
		enterRule(_localctx, 16, RULE_indexDefinition);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(82);
			match(INDEX);
			setState(83);
			match(LPAREN);
			setState(84);
			columnName();
			setState(89);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(85);
				match(COMMA);
				setState(86);
				columnName();
				}
				}
				setState(91);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(92);
			match(RPAREN);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class InsertTableStatementContext extends ParserRuleContext {
		public TerminalNode INSERT() { return getToken(SQLParser.INSERT, 0); }
		public TerminalNode INTO() { return getToken(SQLParser.INTO, 0); }
		public TableNameContext tableName() {
			return getRuleContext(TableNameContext.class,0);
		}
		public List<TerminalNode> LPAREN() { return getTokens(SQLParser.LPAREN); }
		public TerminalNode LPAREN(int i) {
			return getToken(SQLParser.LPAREN, i);
		}
		public ColumnInsertNamesContext columnInsertNames() {
			return getRuleContext(ColumnInsertNamesContext.class,0);
		}
		public List<TerminalNode> RPAREN() { return getTokens(SQLParser.RPAREN); }
		public TerminalNode RPAREN(int i) {
			return getToken(SQLParser.RPAREN, i);
		}
		public TerminalNode VALUES() { return getToken(SQLParser.VALUES, 0); }
		public ColumnInsertValuesContext columnInsertValues() {
			return getRuleContext(ColumnInsertValuesContext.class,0);
		}
		public InsertTableStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_insertTableStatement; }
	}

	public final InsertTableStatementContext insertTableStatement() throws RecognitionException {
		InsertTableStatementContext _localctx = new InsertTableStatementContext(_ctx, getState());
		enterRule(_localctx, 18, RULE_insertTableStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(94);
			match(INSERT);
			setState(95);
			match(INTO);
			setState(96);
			tableName();
			setState(97);
			match(LPAREN);
			setState(98);
			columnInsertNames();
			setState(99);
			match(RPAREN);
			setState(100);
			match(VALUES);
			setState(101);
			match(LPAREN);
			setState(102);
			columnInsertValues();
			setState(103);
			match(RPAREN);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ColumnInsertNamesContext extends ParserRuleContext {
		public List<ColumnNameContext> columnName() {
			return getRuleContexts(ColumnNameContext.class);
		}
		public ColumnNameContext columnName(int i) {
			return getRuleContext(ColumnNameContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(SQLParser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(SQLParser.COMMA, i);
		}
		public ColumnInsertNamesContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_columnInsertNames; }
	}

	public final ColumnInsertNamesContext columnInsertNames() throws RecognitionException {
		ColumnInsertNamesContext _localctx = new ColumnInsertNamesContext(_ctx, getState());
		enterRule(_localctx, 20, RULE_columnInsertNames);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(105);
			columnName();
			setState(110);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(106);
				match(COMMA);
				setState(107);
				columnName();
				}
				}
				setState(112);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ColumnInsertValuesContext extends ParserRuleContext {
		public List<ColumnValueContext> columnValue() {
			return getRuleContexts(ColumnValueContext.class);
		}
		public ColumnValueContext columnValue(int i) {
			return getRuleContext(ColumnValueContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(SQLParser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(SQLParser.COMMA, i);
		}
		public ColumnInsertValuesContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_columnInsertValues; }
	}

	public final ColumnInsertValuesContext columnInsertValues() throws RecognitionException {
		ColumnInsertValuesContext _localctx = new ColumnInsertValuesContext(_ctx, getState());
		enterRule(_localctx, 22, RULE_columnInsertValues);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(113);
			columnValue();
			setState(118);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(114);
				match(COMMA);
				setState(115);
				columnValue();
				}
				}
				setState(120);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ColumnValueContext extends ParserRuleContext {
		public TerminalNode INTEGER() { return getToken(SQLParser.INTEGER, 0); }
		public TerminalNode STRING() { return getToken(SQLParser.STRING, 0); }
		public ColumnValueContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_columnValue; }
	}

	public final ColumnValueContext columnValue() throws RecognitionException {
		ColumnValueContext _localctx = new ColumnValueContext(_ctx, getState());
		enterRule(_localctx, 24, RULE_columnValue);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(121);
			_la = _input.LA(1);
			if ( !(_la==INTEGER || _la==STRING) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class SelectTableStatementContext extends ParserRuleContext {
		public TerminalNode SELECT() { return getToken(SQLParser.SELECT, 0); }
		public SelectColumnNamesContext selectColumnNames() {
			return getRuleContext(SelectColumnNamesContext.class,0);
		}
		public TerminalNode FROM() { return getToken(SQLParser.FROM, 0); }
		public TableNameContext tableName() {
			return getRuleContext(TableNameContext.class,0);
		}
		public TerminalNode WHERE() { return getToken(SQLParser.WHERE, 0); }
		public List<ConditionContext> condition() {
			return getRuleContexts(ConditionContext.class);
		}
		public ConditionContext condition(int i) {
			return getRuleContext(ConditionContext.class,i);
		}
		public List<TerminalNode> AND() { return getTokens(SQLParser.AND); }
		public TerminalNode AND(int i) {
			return getToken(SQLParser.AND, i);
		}
		public SelectTableStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_selectTableStatement; }
	}

	public final SelectTableStatementContext selectTableStatement() throws RecognitionException {
		SelectTableStatementContext _localctx = new SelectTableStatementContext(_ctx, getState());
		enterRule(_localctx, 26, RULE_selectTableStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(123);
			match(SELECT);
			setState(124);
			selectColumnNames();
			setState(125);
			match(FROM);
			setState(126);
			tableName();
			setState(136);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==WHERE) {
				{
				setState(127);
				match(WHERE);
				setState(128);
				condition();
				setState(133);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while (_la==AND) {
					{
					{
					setState(129);
					match(AND);
					setState(130);
					condition();
					}
					}
					setState(135);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class SelectColumnNamesContext extends ParserRuleContext {
		public TerminalNode STAR() { return getToken(SQLParser.STAR, 0); }
		public List<ColumnNameContext> columnName() {
			return getRuleContexts(ColumnNameContext.class);
		}
		public ColumnNameContext columnName(int i) {
			return getRuleContext(ColumnNameContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(SQLParser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(SQLParser.COMMA, i);
		}
		public SelectColumnNamesContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_selectColumnNames; }
	}

	public final SelectColumnNamesContext selectColumnNames() throws RecognitionException {
		SelectColumnNamesContext _localctx = new SelectColumnNamesContext(_ctx, getState());
		enterRule(_localctx, 28, RULE_selectColumnNames);
		int _la;
		try {
			setState(147);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case STAR:
				enterOuterAlt(_localctx, 1);
				{
				setState(138);
				match(STAR);
				}
				break;
			case IDENTIFIER:
				enterOuterAlt(_localctx, 2);
				{
				setState(139);
				columnName();
				setState(144);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while (_la==COMMA) {
					{
					{
					setState(140);
					match(COMMA);
					setState(141);
					columnName();
					}
					}
					setState(146);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ConditionContext extends ParserRuleContext {
		public ComparisonConditionContext comparisonCondition() {
			return getRuleContext(ComparisonConditionContext.class,0);
		}
		public BetweenConditionContext betweenCondition() {
			return getRuleContext(BetweenConditionContext.class,0);
		}
		public ConditionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_condition; }
	}

	public final ConditionContext condition() throws RecognitionException {
		ConditionContext _localctx = new ConditionContext(_ctx, getState());
		enterRule(_localctx, 30, RULE_condition);
		try {
			setState(151);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,10,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(149);
				comparisonCondition();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(150);
				betweenCondition();
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ComparisonConditionContext extends ParserRuleContext {
		public ColumnNameContext columnName() {
			return getRuleContext(ColumnNameContext.class,0);
		}
		public TerminalNode OP() { return getToken(SQLParser.OP, 0); }
		public ColumnValueContext columnValue() {
			return getRuleContext(ColumnValueContext.class,0);
		}
		public ComparisonConditionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_comparisonCondition; }
	}

	public final ComparisonConditionContext comparisonCondition() throws RecognitionException {
		ComparisonConditionContext _localctx = new ComparisonConditionContext(_ctx, getState());
		enterRule(_localctx, 32, RULE_comparisonCondition);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(153);
			columnName();
			setState(154);
			match(OP);
			setState(155);
			columnValue();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class BetweenConditionContext extends ParserRuleContext {
		public ColumnNameContext columnName() {
			return getRuleContext(ColumnNameContext.class,0);
		}
		public TerminalNode BETWEEN() { return getToken(SQLParser.BETWEEN, 0); }
		public List<ColumnValueContext> columnValue() {
			return getRuleContexts(ColumnValueContext.class);
		}
		public ColumnValueContext columnValue(int i) {
			return getRuleContext(ColumnValueContext.class,i);
		}
		public TerminalNode AND() { return getToken(SQLParser.AND, 0); }
		public BetweenConditionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_betweenCondition; }
	}

	public final BetweenConditionContext betweenCondition() throws RecognitionException {
		BetweenConditionContext _localctx = new BetweenConditionContext(_ctx, getState());
		enterRule(_localctx, 34, RULE_betweenCondition);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(157);
			columnName();
			setState(158);
			match(BETWEEN);
			setState(159);
			columnValue();
			setState(160);
			match(AND);
			setState(161);
			columnValue();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static final String _serializedATN =
		"\u0004\u0001\u0019\u00a4\u0002\u0000\u0007\u0000\u0002\u0001\u0007\u0001"+
		"\u0002\u0002\u0007\u0002\u0002\u0003\u0007\u0003\u0002\u0004\u0007\u0004"+
		"\u0002\u0005\u0007\u0005\u0002\u0006\u0007\u0006\u0002\u0007\u0007\u0007"+
		"\u0002\b\u0007\b\u0002\t\u0007\t\u0002\n\u0007\n\u0002\u000b\u0007\u000b"+
		"\u0002\f\u0007\f\u0002\r\u0007\r\u0002\u000e\u0007\u000e\u0002\u000f\u0007"+
		"\u000f\u0002\u0010\u0007\u0010\u0002\u0011\u0007\u0011\u0001\u0000\u0001"+
		"\u0000\u0001\u0000\u0003\u0000(\b\u0000\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0002\u0001\u0002\u0001\u0003\u0001\u0003\u0001"+
		"\u0004\u0001\u0004\u0001\u0005\u0001\u0005\u0001\u0005\u0005\u0005C\b"+
		"\u0005\n\u0005\f\u0005F\t\u0005\u0001\u0006\u0001\u0006\u0001\u0006\u0001"+
		"\u0007\u0001\u0007\u0001\u0007\u0005\u0007N\b\u0007\n\u0007\f\u0007Q\t"+
		"\u0007\u0001\b\u0001\b\u0001\b\u0001\b\u0001\b\u0005\bX\b\b\n\b\f\b[\t"+
		"\b\u0001\b\u0001\b\u0001\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001"+
		"\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001\n\u0001\n\u0001\n\u0005\nm\b"+
		"\n\n\n\f\np\t\n\u0001\u000b\u0001\u000b\u0001\u000b\u0005\u000bu\b\u000b"+
		"\n\u000b\f\u000bx\t\u000b\u0001\f\u0001\f\u0001\r\u0001\r\u0001\r\u0001"+
		"\r\u0001\r\u0001\r\u0001\r\u0001\r\u0005\r\u0084\b\r\n\r\f\r\u0087\t\r"+
		"\u0003\r\u0089\b\r\u0001\u000e\u0001\u000e\u0001\u000e\u0001\u000e\u0005"+
		"\u000e\u008f\b\u000e\n\u000e\f\u000e\u0092\t\u000e\u0003\u000e\u0094\b"+
		"\u000e\u0001\u000f\u0001\u000f\u0003\u000f\u0098\b\u000f\u0001\u0010\u0001"+
		"\u0010\u0001\u0010\u0001\u0010\u0001\u0011\u0001\u0011\u0001\u0011\u0001"+
		"\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0000\u0000\u0012\u0000\u0002"+
		"\u0004\u0006\b\n\f\u000e\u0010\u0012\u0014\u0016\u0018\u001a\u001c\u001e"+
		" \"\u0000\u0002\u0001\u0000\r\u000e\u0001\u0000\u0017\u0018\u009d\u0000"+
		"\'\u0001\u0000\u0000\u0000\u0002)\u0001\u0000\u0000\u0000\u00049\u0001"+
		"\u0000\u0000\u0000\u0006;\u0001\u0000\u0000\u0000\b=\u0001\u0000\u0000"+
		"\u0000\n?\u0001\u0000\u0000\u0000\fG\u0001\u0000\u0000\u0000\u000eJ\u0001"+
		"\u0000\u0000\u0000\u0010R\u0001\u0000\u0000\u0000\u0012^\u0001\u0000\u0000"+
		"\u0000\u0014i\u0001\u0000\u0000\u0000\u0016q\u0001\u0000\u0000\u0000\u0018"+
		"y\u0001\u0000\u0000\u0000\u001a{\u0001\u0000\u0000\u0000\u001c\u0093\u0001"+
		"\u0000\u0000\u0000\u001e\u0097\u0001\u0000\u0000\u0000 \u0099\u0001\u0000"+
		"\u0000\u0000\"\u009d\u0001\u0000\u0000\u0000$(\u0003\u0002\u0001\u0000"+
		"%(\u0003\u0012\t\u0000&(\u0003\u001a\r\u0000\'$\u0001\u0000\u0000\u0000"+
		"\'%\u0001\u0000\u0000\u0000\'&\u0001\u0000\u0000\u0000(\u0001\u0001\u0000"+
		"\u0000\u0000)*\u0005\u0001\u0000\u0000*+\u0005\u0002\u0000\u0000+,\u0003"+
		"\u0004\u0002\u0000,-\u0005\u0011\u0000\u0000-.\u0003\n\u0005\u0000./\u0005"+
		"\u0013\u0000\u0000/0\u0005\u0003\u0000\u000001\u0005\u0004\u0000\u0000"+
		"12\u0005\u0011\u0000\u000023\u0003\u0006\u0003\u000034\u0005\u0012\u0000"+
		"\u000045\u0005\u0013\u0000\u000056\u0003\u000e\u0007\u000067\u0005\u0012"+
		"\u0000\u000078\u0005\u0014\u0000\u00008\u0003\u0001\u0000\u0000\u0000"+
		"9:\u0005\u0010\u0000\u0000:\u0005\u0001\u0000\u0000\u0000;<\u0005\u0010"+
		"\u0000\u0000<\u0007\u0001\u0000\u0000\u0000=>\u0007\u0000\u0000\u0000"+
		">\t\u0001\u0000\u0000\u0000?D\u0003\f\u0006\u0000@A\u0005\u0013\u0000"+
		"\u0000AC\u0003\f\u0006\u0000B@\u0001\u0000\u0000\u0000CF\u0001\u0000\u0000"+
		"\u0000DB\u0001\u0000\u0000\u0000DE\u0001\u0000\u0000\u0000E\u000b\u0001"+
		"\u0000\u0000\u0000FD\u0001\u0000\u0000\u0000GH\u0003\u0006\u0003\u0000"+
		"HI\u0003\b\u0004\u0000I\r\u0001\u0000\u0000\u0000JO\u0003\u0010\b\u0000"+
		"KL\u0005\u0013\u0000\u0000LN\u0003\u0010\b\u0000MK\u0001\u0000\u0000\u0000"+
		"NQ\u0001\u0000\u0000\u0000OM\u0001\u0000\u0000\u0000OP\u0001\u0000\u0000"+
		"\u0000P\u000f\u0001\u0000\u0000\u0000QO\u0001\u0000\u0000\u0000RS\u0005"+
		"\u0005\u0000\u0000ST\u0005\u0011\u0000\u0000TY\u0003\u0006\u0003\u0000"+
		"UV\u0005\u0013\u0000\u0000VX\u0003\u0006\u0003\u0000WU\u0001\u0000\u0000"+
		"\u0000X[\u0001\u0000\u0000\u0000YW\u0001\u0000\u0000\u0000YZ\u0001\u0000"+
		"\u0000\u0000Z\\\u0001\u0000\u0000\u0000[Y\u0001\u0000\u0000\u0000\\]\u0005"+
		"\u0012\u0000\u0000]\u0011\u0001\u0000\u0000\u0000^_\u0005\u0006\u0000"+
		"\u0000_`\u0005\u0007\u0000\u0000`a\u0003\u0004\u0002\u0000ab\u0005\u0011"+
		"\u0000\u0000bc\u0003\u0014\n\u0000cd\u0005\u0012\u0000\u0000de\u0005\b"+
		"\u0000\u0000ef\u0005\u0011\u0000\u0000fg\u0003\u0016\u000b\u0000gh\u0005"+
		"\u0012\u0000\u0000h\u0013\u0001\u0000\u0000\u0000in\u0003\u0006\u0003"+
		"\u0000jk\u0005\u0013\u0000\u0000km\u0003\u0006\u0003\u0000lj\u0001\u0000"+
		"\u0000\u0000mp\u0001\u0000\u0000\u0000nl\u0001\u0000\u0000\u0000no\u0001"+
		"\u0000\u0000\u0000o\u0015\u0001\u0000\u0000\u0000pn\u0001\u0000\u0000"+
		"\u0000qv\u0003\u0018\f\u0000rs\u0005\u0013\u0000\u0000su\u0003\u0018\f"+
		"\u0000tr\u0001\u0000\u0000\u0000ux\u0001\u0000\u0000\u0000vt\u0001\u0000"+
		"\u0000\u0000vw\u0001\u0000\u0000\u0000w\u0017\u0001\u0000\u0000\u0000"+
		"xv\u0001\u0000\u0000\u0000yz\u0007\u0001\u0000\u0000z\u0019\u0001\u0000"+
		"\u0000\u0000{|\u0005\t\u0000\u0000|}\u0003\u001c\u000e\u0000}~\u0005\n"+
		"\u0000\u0000~\u0088\u0003\u0004\u0002\u0000\u007f\u0080\u0005\u000b\u0000"+
		"\u0000\u0080\u0085\u0003\u001e\u000f\u0000\u0081\u0082\u0005\f\u0000\u0000"+
		"\u0082\u0084\u0003\u001e\u000f\u0000\u0083\u0081\u0001\u0000\u0000\u0000"+
		"\u0084\u0087\u0001\u0000\u0000\u0000\u0085\u0083\u0001\u0000\u0000\u0000"+
		"\u0085\u0086\u0001\u0000\u0000\u0000\u0086\u0089\u0001\u0000\u0000\u0000"+
		"\u0087\u0085\u0001\u0000\u0000\u0000\u0088\u007f\u0001\u0000\u0000\u0000"+
		"\u0088\u0089\u0001\u0000\u0000\u0000\u0089\u001b\u0001\u0000\u0000\u0000"+
		"\u008a\u0094\u0005\u0015\u0000\u0000\u008b\u0090\u0003\u0006\u0003\u0000"+
		"\u008c\u008d\u0005\u0013\u0000\u0000\u008d\u008f\u0003\u0006\u0003\u0000"+
		"\u008e\u008c\u0001\u0000\u0000\u0000\u008f\u0092\u0001\u0000\u0000\u0000"+
		"\u0090\u008e\u0001\u0000\u0000\u0000\u0090\u0091\u0001\u0000\u0000\u0000"+
		"\u0091\u0094\u0001\u0000\u0000\u0000\u0092\u0090\u0001\u0000\u0000\u0000"+
		"\u0093\u008a\u0001\u0000\u0000\u0000\u0093\u008b\u0001\u0000\u0000\u0000"+
		"\u0094\u001d\u0001\u0000\u0000\u0000\u0095\u0098\u0003 \u0010\u0000\u0096"+
		"\u0098\u0003\"\u0011\u0000\u0097\u0095\u0001\u0000\u0000\u0000\u0097\u0096"+
		"\u0001\u0000\u0000\u0000\u0098\u001f\u0001\u0000\u0000\u0000\u0099\u009a"+
		"\u0003\u0006\u0003\u0000\u009a\u009b\u0005\u0019\u0000\u0000\u009b\u009c"+
		"\u0003\u0018\f\u0000\u009c!\u0001\u0000\u0000\u0000\u009d\u009e\u0003"+
		"\u0006\u0003\u0000\u009e\u009f\u0005\u000f\u0000\u0000\u009f\u00a0\u0003"+
		"\u0018\f\u0000\u00a0\u00a1\u0005\f\u0000\u0000\u00a1\u00a2\u0003\u0018"+
		"\f\u0000\u00a2#\u0001\u0000\u0000\u0000\u000b\'DOYnv\u0085\u0088\u0090"+
		"\u0093\u0097";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}