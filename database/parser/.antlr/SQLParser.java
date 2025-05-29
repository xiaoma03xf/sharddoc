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
		INT64=9, BYTES=10, IDENTIFIER=11, LPAREN=12, RPAREN=13, COMMA=14, SEMICOLON=15, 
		WS=16, INTEGER=17, STRING=18;
	public static final int
		RULE_sql = 0, RULE_createTableStatement = 1, RULE_tableName = 2, RULE_columnName = 3, 
		RULE_columnType = 4, RULE_columnDefinitions = 5, RULE_columnDefinition = 6, 
		RULE_indexDefinitions = 7, RULE_indexDefinition = 8, RULE_insertTableStatement = 9, 
		RULE_columnInsertNames = 10, RULE_columnInsertValues = 11, RULE_columnValue = 12;
	private static String[] makeRuleNames() {
		return new String[] {
			"sql", "createTableStatement", "tableName", "columnName", "columnType", 
			"columnDefinitions", "columnDefinition", "indexDefinitions", "indexDefinition", 
			"insertTableStatement", "columnInsertNames", "columnInsertValues", "columnValue"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'CREATE'", "'TABLE'", "'PRIMARY'", "'KEY'", "'INDEX'", "'INSERT'", 
			"'INTO'", "'VALUES'", "'INT64'", "'BYTES'", null, "'('", "')'", "','", 
			"';'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, "CREATE", "TABLE", "PRIMARY", "KEY", "INDEX", "INSERT", "INTO", 
			"VALUES", "INT64", "BYTES", "IDENTIFIER", "LPAREN", "RPAREN", "COMMA", 
			"SEMICOLON", "WS", "INTEGER", "STRING"
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
		public SqlContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_sql; }
	}

	public final SqlContext sql() throws RecognitionException {
		SqlContext _localctx = new SqlContext(_ctx, getState());
		enterRule(_localctx, 0, RULE_sql);
		try {
			setState(28);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case CREATE:
				enterOuterAlt(_localctx, 1);
				{
				setState(26);
				createTableStatement();
				}
				break;
			case INSERT:
				enterOuterAlt(_localctx, 2);
				{
				setState(27);
				insertTableStatement();
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
			setState(30);
			match(CREATE);
			setState(31);
			match(TABLE);
			setState(32);
			tableName();
			setState(33);
			match(LPAREN);
			setState(34);
			columnDefinitions();
			setState(35);
			match(COMMA);
			setState(36);
			match(PRIMARY);
			setState(37);
			match(KEY);
			setState(38);
			match(LPAREN);
			setState(39);
			columnName();
			setState(40);
			match(RPAREN);
			setState(41);
			match(COMMA);
			setState(42);
			indexDefinitions();
			setState(43);
			match(RPAREN);
			setState(44);
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
			setState(46);
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
			setState(48);
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
			setState(50);
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
			setState(52);
			columnDefinition();
			setState(57);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,1,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(53);
					match(COMMA);
					setState(54);
					columnDefinition();
					}
					} 
				}
				setState(59);
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
			setState(60);
			columnName();
			setState(61);
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
			setState(63);
			indexDefinition();
			setState(68);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(64);
				match(COMMA);
				setState(65);
				indexDefinition();
				}
				}
				setState(70);
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
			setState(71);
			match(INDEX);
			setState(72);
			match(LPAREN);
			setState(73);
			columnName();
			setState(78);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(74);
				match(COMMA);
				setState(75);
				columnName();
				}
				}
				setState(80);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(81);
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
			setState(83);
			match(INSERT);
			setState(84);
			match(INTO);
			setState(85);
			tableName();
			setState(86);
			match(LPAREN);
			setState(87);
			columnInsertNames();
			setState(88);
			match(RPAREN);
			setState(89);
			match(VALUES);
			setState(90);
			match(LPAREN);
			setState(91);
			columnInsertValues();
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
			setState(94);
			columnName();
			setState(99);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(95);
				match(COMMA);
				setState(96);
				columnName();
				}
				}
				setState(101);
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
			setState(102);
			columnValue();
			setState(107);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(103);
				match(COMMA);
				setState(104);
				columnValue();
				}
				}
				setState(109);
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
			setState(110);
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

	public static final String _serializedATN =
		"\u0004\u0001\u0012q\u0002\u0000\u0007\u0000\u0002\u0001\u0007\u0001\u0002"+
		"\u0002\u0007\u0002\u0002\u0003\u0007\u0003\u0002\u0004\u0007\u0004\u0002"+
		"\u0005\u0007\u0005\u0002\u0006\u0007\u0006\u0002\u0007\u0007\u0007\u0002"+
		"\b\u0007\b\u0002\t\u0007\t\u0002\n\u0007\n\u0002\u000b\u0007\u000b\u0002"+
		"\f\u0007\f\u0001\u0000\u0001\u0000\u0003\u0000\u001d\b\u0000\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0002\u0001\u0002\u0001\u0003"+
		"\u0001\u0003\u0001\u0004\u0001\u0004\u0001\u0005\u0001\u0005\u0001\u0005"+
		"\u0005\u00058\b\u0005\n\u0005\f\u0005;\t\u0005\u0001\u0006\u0001\u0006"+
		"\u0001\u0006\u0001\u0007\u0001\u0007\u0001\u0007\u0005\u0007C\b\u0007"+
		"\n\u0007\f\u0007F\t\u0007\u0001\b\u0001\b\u0001\b\u0001\b\u0001\b\u0005"+
		"\bM\b\b\n\b\f\bP\t\b\u0001\b\u0001\b\u0001\t\u0001\t\u0001\t\u0001\t\u0001"+
		"\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001\n\u0001\n\u0001"+
		"\n\u0005\nb\b\n\n\n\f\ne\t\n\u0001\u000b\u0001\u000b\u0001\u000b\u0005"+
		"\u000bj\b\u000b\n\u000b\f\u000bm\t\u000b\u0001\f\u0001\f\u0001\f\u0000"+
		"\u0000\r\u0000\u0002\u0004\u0006\b\n\f\u000e\u0010\u0012\u0014\u0016\u0018"+
		"\u0000\u0002\u0001\u0000\t\n\u0001\u0000\u0011\u0012i\u0000\u001c\u0001"+
		"\u0000\u0000\u0000\u0002\u001e\u0001\u0000\u0000\u0000\u0004.\u0001\u0000"+
		"\u0000\u0000\u00060\u0001\u0000\u0000\u0000\b2\u0001\u0000\u0000\u0000"+
		"\n4\u0001\u0000\u0000\u0000\f<\u0001\u0000\u0000\u0000\u000e?\u0001\u0000"+
		"\u0000\u0000\u0010G\u0001\u0000\u0000\u0000\u0012S\u0001\u0000\u0000\u0000"+
		"\u0014^\u0001\u0000\u0000\u0000\u0016f\u0001\u0000\u0000\u0000\u0018n"+
		"\u0001\u0000\u0000\u0000\u001a\u001d\u0003\u0002\u0001\u0000\u001b\u001d"+
		"\u0003\u0012\t\u0000\u001c\u001a\u0001\u0000\u0000\u0000\u001c\u001b\u0001"+
		"\u0000\u0000\u0000\u001d\u0001\u0001\u0000\u0000\u0000\u001e\u001f\u0005"+
		"\u0001\u0000\u0000\u001f \u0005\u0002\u0000\u0000 !\u0003\u0004\u0002"+
		"\u0000!\"\u0005\f\u0000\u0000\"#\u0003\n\u0005\u0000#$\u0005\u000e\u0000"+
		"\u0000$%\u0005\u0003\u0000\u0000%&\u0005\u0004\u0000\u0000&\'\u0005\f"+
		"\u0000\u0000\'(\u0003\u0006\u0003\u0000()\u0005\r\u0000\u0000)*\u0005"+
		"\u000e\u0000\u0000*+\u0003\u000e\u0007\u0000+,\u0005\r\u0000\u0000,-\u0005"+
		"\u000f\u0000\u0000-\u0003\u0001\u0000\u0000\u0000./\u0005\u000b\u0000"+
		"\u0000/\u0005\u0001\u0000\u0000\u000001\u0005\u000b\u0000\u00001\u0007"+
		"\u0001\u0000\u0000\u000023\u0007\u0000\u0000\u00003\t\u0001\u0000\u0000"+
		"\u000049\u0003\f\u0006\u000056\u0005\u000e\u0000\u000068\u0003\f\u0006"+
		"\u000075\u0001\u0000\u0000\u00008;\u0001\u0000\u0000\u000097\u0001\u0000"+
		"\u0000\u00009:\u0001\u0000\u0000\u0000:\u000b\u0001\u0000\u0000\u0000"+
		";9\u0001\u0000\u0000\u0000<=\u0003\u0006\u0003\u0000=>\u0003\b\u0004\u0000"+
		">\r\u0001\u0000\u0000\u0000?D\u0003\u0010\b\u0000@A\u0005\u000e\u0000"+
		"\u0000AC\u0003\u0010\b\u0000B@\u0001\u0000\u0000\u0000CF\u0001\u0000\u0000"+
		"\u0000DB\u0001\u0000\u0000\u0000DE\u0001\u0000\u0000\u0000E\u000f\u0001"+
		"\u0000\u0000\u0000FD\u0001\u0000\u0000\u0000GH\u0005\u0005\u0000\u0000"+
		"HI\u0005\f\u0000\u0000IN\u0003\u0006\u0003\u0000JK\u0005\u000e\u0000\u0000"+
		"KM\u0003\u0006\u0003\u0000LJ\u0001\u0000\u0000\u0000MP\u0001\u0000\u0000"+
		"\u0000NL\u0001\u0000\u0000\u0000NO\u0001\u0000\u0000\u0000OQ\u0001\u0000"+
		"\u0000\u0000PN\u0001\u0000\u0000\u0000QR\u0005\r\u0000\u0000R\u0011\u0001"+
		"\u0000\u0000\u0000ST\u0005\u0006\u0000\u0000TU\u0005\u0007\u0000\u0000"+
		"UV\u0003\u0004\u0002\u0000VW\u0005\f\u0000\u0000WX\u0003\u0014\n\u0000"+
		"XY\u0005\r\u0000\u0000YZ\u0005\b\u0000\u0000Z[\u0005\f\u0000\u0000[\\"+
		"\u0003\u0016\u000b\u0000\\]\u0005\r\u0000\u0000]\u0013\u0001\u0000\u0000"+
		"\u0000^c\u0003\u0006\u0003\u0000_`\u0005\u000e\u0000\u0000`b\u0003\u0006"+
		"\u0003\u0000a_\u0001\u0000\u0000\u0000be\u0001\u0000\u0000\u0000ca\u0001"+
		"\u0000\u0000\u0000cd\u0001\u0000\u0000\u0000d\u0015\u0001\u0000\u0000"+
		"\u0000ec\u0001\u0000\u0000\u0000fk\u0003\u0018\f\u0000gh\u0005\u000e\u0000"+
		"\u0000hj\u0003\u0018\f\u0000ig\u0001\u0000\u0000\u0000jm\u0001\u0000\u0000"+
		"\u0000ki\u0001\u0000\u0000\u0000kl\u0001\u0000\u0000\u0000l\u0017\u0001"+
		"\u0000\u0000\u0000mk\u0001\u0000\u0000\u0000no\u0007\u0001\u0000\u0000"+
		"o\u0019\u0001\u0000\u0000\u0000\u0006\u001c9DNck";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}