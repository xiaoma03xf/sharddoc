// Generated from /home/dengjun/sharddoc/parser/SQLParser.g4 by ANTLR 4.13.1
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
		SELECT=9, FROM=10, WHERE=11, AND=12, DELETE=13, UPDATE=14, SET=15, INT64=16, 
		BYTES=17, BETWEEN=18, IDENTIFIER=19, LPAREN=20, RPAREN=21, COMMA=22, SEMICOLON=23, 
		STAR=24, WS=25, INTEGER=26, STRING=27, OP=28;
	public static final int
		RULE_sql = 0, RULE_createTableStatement = 1, RULE_tableName = 2, RULE_columnName = 3, 
		RULE_columnType = 4, RULE_columnDefinitions = 5, RULE_columnDefinition = 6, 
		RULE_indexDefinitions = 7, RULE_indexDefinition = 8, RULE_insertTableStatement = 9, 
		RULE_columnInsertNames = 10, RULE_columnInsertValues = 11, RULE_columnValue = 12, 
		RULE_selectTableStatement = 13, RULE_conditions = 14, RULE_selectColumnNames = 15, 
		RULE_condition = 16, RULE_comparisonCondition = 17, RULE_betweenCondition = 18, 
		RULE_updateTableStatement = 19, RULE_setClauses = 20, RULE_setClause = 21, 
		RULE_deleteTableStatement = 22;
	private static String[] makeRuleNames() {
		return new String[] {
			"sql", "createTableStatement", "tableName", "columnName", "columnType", 
			"columnDefinitions", "columnDefinition", "indexDefinitions", "indexDefinition", 
			"insertTableStatement", "columnInsertNames", "columnInsertValues", "columnValue", 
			"selectTableStatement", "conditions", "selectColumnNames", "condition", 
			"comparisonCondition", "betweenCondition", "updateTableStatement", "setClauses", 
			"setClause", "deleteTableStatement"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'CREATE'", "'TABLE'", "'PRIMARY'", "'KEY'", "'INDEX'", "'INSERT'", 
			"'INTO'", "'VALUES'", "'SELECT'", "'FROM'", "'WHERE'", "'AND'", "'DELETE'", 
			"'UPDATE'", "'SET'", "'INT64'", "'BYTES'", "'BETWEEN'", null, "'('", 
			"')'", "','", "';'", "'*'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, "CREATE", "TABLE", "PRIMARY", "KEY", "INDEX", "INSERT", "INTO", 
			"VALUES", "SELECT", "FROM", "WHERE", "AND", "DELETE", "UPDATE", "SET", 
			"INT64", "BYTES", "BETWEEN", "IDENTIFIER", "LPAREN", "RPAREN", "COMMA", 
			"SEMICOLON", "STAR", "WS", "INTEGER", "STRING", "OP"
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
		public UpdateTableStatementContext updateTableStatement() {
			return getRuleContext(UpdateTableStatementContext.class,0);
		}
		public DeleteTableStatementContext deleteTableStatement() {
			return getRuleContext(DeleteTableStatementContext.class,0);
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
			setState(51);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case CREATE:
				enterOuterAlt(_localctx, 1);
				{
				setState(46);
				createTableStatement();
				}
				break;
			case INSERT:
				enterOuterAlt(_localctx, 2);
				{
				setState(47);
				insertTableStatement();
				}
				break;
			case SELECT:
				enterOuterAlt(_localctx, 3);
				{
				setState(48);
				selectTableStatement();
				}
				break;
			case UPDATE:
				enterOuterAlt(_localctx, 4);
				{
				setState(49);
				updateTableStatement();
				}
				break;
			case DELETE:
				enterOuterAlt(_localctx, 5);
				{
				setState(50);
				deleteTableStatement();
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
			setState(53);
			match(CREATE);
			setState(54);
			match(TABLE);
			setState(55);
			tableName();
			setState(56);
			match(LPAREN);
			setState(57);
			columnDefinitions();
			setState(58);
			match(COMMA);
			setState(59);
			match(PRIMARY);
			setState(60);
			match(KEY);
			setState(61);
			match(LPAREN);
			setState(62);
			columnName();
			setState(63);
			match(RPAREN);
			setState(64);
			match(COMMA);
			setState(65);
			indexDefinitions();
			setState(66);
			match(RPAREN);
			setState(67);
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
			setState(69);
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
			setState(71);
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
			setState(73);
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
			setState(75);
			columnDefinition();
			setState(80);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,1,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(76);
					match(COMMA);
					setState(77);
					columnDefinition();
					}
					} 
				}
				setState(82);
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
			setState(83);
			columnName();
			setState(84);
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
			setState(86);
			indexDefinition();
			setState(91);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(87);
				match(COMMA);
				setState(88);
				indexDefinition();
				}
				}
				setState(93);
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
			setState(94);
			match(INDEX);
			setState(95);
			match(LPAREN);
			setState(96);
			columnName();
			setState(101);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(97);
				match(COMMA);
				setState(98);
				columnName();
				}
				}
				setState(103);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(104);
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
			setState(106);
			match(INSERT);
			setState(107);
			match(INTO);
			setState(108);
			tableName();
			setState(109);
			match(LPAREN);
			setState(110);
			columnInsertNames();
			setState(111);
			match(RPAREN);
			setState(112);
			match(VALUES);
			setState(113);
			match(LPAREN);
			setState(114);
			columnInsertValues();
			setState(115);
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
			setState(117);
			columnName();
			setState(122);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(118);
				match(COMMA);
				setState(119);
				columnName();
				}
				}
				setState(124);
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
			setState(125);
			columnValue();
			setState(130);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(126);
				match(COMMA);
				setState(127);
				columnValue();
				}
				}
				setState(132);
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
			setState(133);
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
		public ConditionsContext conditions() {
			return getRuleContext(ConditionsContext.class,0);
		}
		public SelectTableStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_selectTableStatement; }
	}

	public final SelectTableStatementContext selectTableStatement() throws RecognitionException {
		SelectTableStatementContext _localctx = new SelectTableStatementContext(_ctx, getState());
		enterRule(_localctx, 26, RULE_selectTableStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(135);
			match(SELECT);
			setState(136);
			selectColumnNames();
			setState(137);
			match(FROM);
			setState(138);
			tableName();
			setState(139);
			conditions();
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
	public static class ConditionsContext extends ParserRuleContext {
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
		public ConditionsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_conditions; }
	}

	public final ConditionsContext conditions() throws RecognitionException {
		ConditionsContext _localctx = new ConditionsContext(_ctx, getState());
		enterRule(_localctx, 28, RULE_conditions);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(150);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==WHERE) {
				{
				setState(141);
				match(WHERE);
				setState(142);
				condition();
				setState(147);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while (_la==AND) {
					{
					{
					setState(143);
					match(AND);
					setState(144);
					condition();
					}
					}
					setState(149);
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
		enterRule(_localctx, 30, RULE_selectColumnNames);
		int _la;
		try {
			setState(161);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case STAR:
				enterOuterAlt(_localctx, 1);
				{
				setState(152);
				match(STAR);
				}
				break;
			case IDENTIFIER:
				enterOuterAlt(_localctx, 2);
				{
				setState(153);
				columnName();
				setState(158);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while (_la==COMMA) {
					{
					{
					setState(154);
					match(COMMA);
					setState(155);
					columnName();
					}
					}
					setState(160);
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
		enterRule(_localctx, 32, RULE_condition);
		try {
			setState(165);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,10,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(163);
				comparisonCondition();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(164);
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
		enterRule(_localctx, 34, RULE_comparisonCondition);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(167);
			columnName();
			setState(168);
			match(OP);
			setState(169);
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
		enterRule(_localctx, 36, RULE_betweenCondition);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(171);
			columnName();
			setState(172);
			match(BETWEEN);
			setState(173);
			columnValue();
			setState(174);
			match(AND);
			setState(175);
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
	public static class UpdateTableStatementContext extends ParserRuleContext {
		public TerminalNode UPDATE() { return getToken(SQLParser.UPDATE, 0); }
		public TableNameContext tableName() {
			return getRuleContext(TableNameContext.class,0);
		}
		public TerminalNode SET() { return getToken(SQLParser.SET, 0); }
		public SetClausesContext setClauses() {
			return getRuleContext(SetClausesContext.class,0);
		}
		public ConditionsContext conditions() {
			return getRuleContext(ConditionsContext.class,0);
		}
		public UpdateTableStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_updateTableStatement; }
	}

	public final UpdateTableStatementContext updateTableStatement() throws RecognitionException {
		UpdateTableStatementContext _localctx = new UpdateTableStatementContext(_ctx, getState());
		enterRule(_localctx, 38, RULE_updateTableStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(177);
			match(UPDATE);
			setState(178);
			tableName();
			setState(179);
			match(SET);
			setState(180);
			setClauses();
			setState(181);
			conditions();
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
	public static class SetClausesContext extends ParserRuleContext {
		public List<SetClauseContext> setClause() {
			return getRuleContexts(SetClauseContext.class);
		}
		public SetClauseContext setClause(int i) {
			return getRuleContext(SetClauseContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(SQLParser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(SQLParser.COMMA, i);
		}
		public SetClausesContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_setClauses; }
	}

	public final SetClausesContext setClauses() throws RecognitionException {
		SetClausesContext _localctx = new SetClausesContext(_ctx, getState());
		enterRule(_localctx, 40, RULE_setClauses);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(183);
			setClause();
			setState(188);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(184);
				match(COMMA);
				setState(185);
				setClause();
				}
				}
				setState(190);
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
	public static class SetClauseContext extends ParserRuleContext {
		public ColumnNameContext columnName() {
			return getRuleContext(ColumnNameContext.class,0);
		}
		public TerminalNode OP() { return getToken(SQLParser.OP, 0); }
		public ColumnValueContext columnValue() {
			return getRuleContext(ColumnValueContext.class,0);
		}
		public SetClauseContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_setClause; }
	}

	public final SetClauseContext setClause() throws RecognitionException {
		SetClauseContext _localctx = new SetClauseContext(_ctx, getState());
		enterRule(_localctx, 42, RULE_setClause);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(191);
			columnName();
			setState(192);
			match(OP);
			setState(193);
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
	public static class DeleteTableStatementContext extends ParserRuleContext {
		public TerminalNode DELETE() { return getToken(SQLParser.DELETE, 0); }
		public TerminalNode FROM() { return getToken(SQLParser.FROM, 0); }
		public TableNameContext tableName() {
			return getRuleContext(TableNameContext.class,0);
		}
		public ConditionsContext conditions() {
			return getRuleContext(ConditionsContext.class,0);
		}
		public DeleteTableStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_deleteTableStatement; }
	}

	public final DeleteTableStatementContext deleteTableStatement() throws RecognitionException {
		DeleteTableStatementContext _localctx = new DeleteTableStatementContext(_ctx, getState());
		enterRule(_localctx, 44, RULE_deleteTableStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(195);
			match(DELETE);
			setState(196);
			match(FROM);
			setState(197);
			tableName();
			setState(198);
			conditions();
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
		"\u0004\u0001\u001c\u00c9\u0002\u0000\u0007\u0000\u0002\u0001\u0007\u0001"+
		"\u0002\u0002\u0007\u0002\u0002\u0003\u0007\u0003\u0002\u0004\u0007\u0004"+
		"\u0002\u0005\u0007\u0005\u0002\u0006\u0007\u0006\u0002\u0007\u0007\u0007"+
		"\u0002\b\u0007\b\u0002\t\u0007\t\u0002\n\u0007\n\u0002\u000b\u0007\u000b"+
		"\u0002\f\u0007\f\u0002\r\u0007\r\u0002\u000e\u0007\u000e\u0002\u000f\u0007"+
		"\u000f\u0002\u0010\u0007\u0010\u0002\u0011\u0007\u0011\u0002\u0012\u0007"+
		"\u0012\u0002\u0013\u0007\u0013\u0002\u0014\u0007\u0014\u0002\u0015\u0007"+
		"\u0015\u0002\u0016\u0007\u0016\u0001\u0000\u0001\u0000\u0001\u0000\u0001"+
		"\u0000\u0001\u0000\u0003\u00004\b\u0000\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0002\u0001\u0002\u0001\u0003\u0001\u0003\u0001"+
		"\u0004\u0001\u0004\u0001\u0005\u0001\u0005\u0001\u0005\u0005\u0005O\b"+
		"\u0005\n\u0005\f\u0005R\t\u0005\u0001\u0006\u0001\u0006\u0001\u0006\u0001"+
		"\u0007\u0001\u0007\u0001\u0007\u0005\u0007Z\b\u0007\n\u0007\f\u0007]\t"+
		"\u0007\u0001\b\u0001\b\u0001\b\u0001\b\u0001\b\u0005\bd\b\b\n\b\f\bg\t"+
		"\b\u0001\b\u0001\b\u0001\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001"+
		"\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001\n\u0001\n\u0001\n\u0005\ny\b"+
		"\n\n\n\f\n|\t\n\u0001\u000b\u0001\u000b\u0001\u000b\u0005\u000b\u0081"+
		"\b\u000b\n\u000b\f\u000b\u0084\t\u000b\u0001\f\u0001\f\u0001\r\u0001\r"+
		"\u0001\r\u0001\r\u0001\r\u0001\r\u0001\u000e\u0001\u000e\u0001\u000e\u0001"+
		"\u000e\u0005\u000e\u0092\b\u000e\n\u000e\f\u000e\u0095\t\u000e\u0003\u000e"+
		"\u0097\b\u000e\u0001\u000f\u0001\u000f\u0001\u000f\u0001\u000f\u0005\u000f"+
		"\u009d\b\u000f\n\u000f\f\u000f\u00a0\t\u000f\u0003\u000f\u00a2\b\u000f"+
		"\u0001\u0010\u0001\u0010\u0003\u0010\u00a6\b\u0010\u0001\u0011\u0001\u0011"+
		"\u0001\u0011\u0001\u0011\u0001\u0012\u0001\u0012\u0001\u0012\u0001\u0012"+
		"\u0001\u0012\u0001\u0012\u0001\u0013\u0001\u0013\u0001\u0013\u0001\u0013"+
		"\u0001\u0013\u0001\u0013\u0001\u0014\u0001\u0014\u0001\u0014\u0005\u0014"+
		"\u00bb\b\u0014\n\u0014\f\u0014\u00be\t\u0014\u0001\u0015\u0001\u0015\u0001"+
		"\u0015\u0001\u0015\u0001\u0016\u0001\u0016\u0001\u0016\u0001\u0016\u0001"+
		"\u0016\u0001\u0016\u0000\u0000\u0017\u0000\u0002\u0004\u0006\b\n\f\u000e"+
		"\u0010\u0012\u0014\u0016\u0018\u001a\u001c\u001e \"$&(*,\u0000\u0002\u0001"+
		"\u0000\u0010\u0011\u0001\u0000\u001a\u001b\u00c0\u00003\u0001\u0000\u0000"+
		"\u0000\u00025\u0001\u0000\u0000\u0000\u0004E\u0001\u0000\u0000\u0000\u0006"+
		"G\u0001\u0000\u0000\u0000\bI\u0001\u0000\u0000\u0000\nK\u0001\u0000\u0000"+
		"\u0000\fS\u0001\u0000\u0000\u0000\u000eV\u0001\u0000\u0000\u0000\u0010"+
		"^\u0001\u0000\u0000\u0000\u0012j\u0001\u0000\u0000\u0000\u0014u\u0001"+
		"\u0000\u0000\u0000\u0016}\u0001\u0000\u0000\u0000\u0018\u0085\u0001\u0000"+
		"\u0000\u0000\u001a\u0087\u0001\u0000\u0000\u0000\u001c\u0096\u0001\u0000"+
		"\u0000\u0000\u001e\u00a1\u0001\u0000\u0000\u0000 \u00a5\u0001\u0000\u0000"+
		"\u0000\"\u00a7\u0001\u0000\u0000\u0000$\u00ab\u0001\u0000\u0000\u0000"+
		"&\u00b1\u0001\u0000\u0000\u0000(\u00b7\u0001\u0000\u0000\u0000*\u00bf"+
		"\u0001\u0000\u0000\u0000,\u00c3\u0001\u0000\u0000\u0000.4\u0003\u0002"+
		"\u0001\u0000/4\u0003\u0012\t\u000004\u0003\u001a\r\u000014\u0003&\u0013"+
		"\u000024\u0003,\u0016\u00003.\u0001\u0000\u0000\u00003/\u0001\u0000\u0000"+
		"\u000030\u0001\u0000\u0000\u000031\u0001\u0000\u0000\u000032\u0001\u0000"+
		"\u0000\u00004\u0001\u0001\u0000\u0000\u000056\u0005\u0001\u0000\u0000"+
		"67\u0005\u0002\u0000\u000078\u0003\u0004\u0002\u000089\u0005\u0014\u0000"+
		"\u00009:\u0003\n\u0005\u0000:;\u0005\u0016\u0000\u0000;<\u0005\u0003\u0000"+
		"\u0000<=\u0005\u0004\u0000\u0000=>\u0005\u0014\u0000\u0000>?\u0003\u0006"+
		"\u0003\u0000?@\u0005\u0015\u0000\u0000@A\u0005\u0016\u0000\u0000AB\u0003"+
		"\u000e\u0007\u0000BC\u0005\u0015\u0000\u0000CD\u0005\u0017\u0000\u0000"+
		"D\u0003\u0001\u0000\u0000\u0000EF\u0005\u0013\u0000\u0000F\u0005\u0001"+
		"\u0000\u0000\u0000GH\u0005\u0013\u0000\u0000H\u0007\u0001\u0000\u0000"+
		"\u0000IJ\u0007\u0000\u0000\u0000J\t\u0001\u0000\u0000\u0000KP\u0003\f"+
		"\u0006\u0000LM\u0005\u0016\u0000\u0000MO\u0003\f\u0006\u0000NL\u0001\u0000"+
		"\u0000\u0000OR\u0001\u0000\u0000\u0000PN\u0001\u0000\u0000\u0000PQ\u0001"+
		"\u0000\u0000\u0000Q\u000b\u0001\u0000\u0000\u0000RP\u0001\u0000\u0000"+
		"\u0000ST\u0003\u0006\u0003\u0000TU\u0003\b\u0004\u0000U\r\u0001\u0000"+
		"\u0000\u0000V[\u0003\u0010\b\u0000WX\u0005\u0016\u0000\u0000XZ\u0003\u0010"+
		"\b\u0000YW\u0001\u0000\u0000\u0000Z]\u0001\u0000\u0000\u0000[Y\u0001\u0000"+
		"\u0000\u0000[\\\u0001\u0000\u0000\u0000\\\u000f\u0001\u0000\u0000\u0000"+
		"][\u0001\u0000\u0000\u0000^_\u0005\u0005\u0000\u0000_`\u0005\u0014\u0000"+
		"\u0000`e\u0003\u0006\u0003\u0000ab\u0005\u0016\u0000\u0000bd\u0003\u0006"+
		"\u0003\u0000ca\u0001\u0000\u0000\u0000dg\u0001\u0000\u0000\u0000ec\u0001"+
		"\u0000\u0000\u0000ef\u0001\u0000\u0000\u0000fh\u0001\u0000\u0000\u0000"+
		"ge\u0001\u0000\u0000\u0000hi\u0005\u0015\u0000\u0000i\u0011\u0001\u0000"+
		"\u0000\u0000jk\u0005\u0006\u0000\u0000kl\u0005\u0007\u0000\u0000lm\u0003"+
		"\u0004\u0002\u0000mn\u0005\u0014\u0000\u0000no\u0003\u0014\n\u0000op\u0005"+
		"\u0015\u0000\u0000pq\u0005\b\u0000\u0000qr\u0005\u0014\u0000\u0000rs\u0003"+
		"\u0016\u000b\u0000st\u0005\u0015\u0000\u0000t\u0013\u0001\u0000\u0000"+
		"\u0000uz\u0003\u0006\u0003\u0000vw\u0005\u0016\u0000\u0000wy\u0003\u0006"+
		"\u0003\u0000xv\u0001\u0000\u0000\u0000y|\u0001\u0000\u0000\u0000zx\u0001"+
		"\u0000\u0000\u0000z{\u0001\u0000\u0000\u0000{\u0015\u0001\u0000\u0000"+
		"\u0000|z\u0001\u0000\u0000\u0000}\u0082\u0003\u0018\f\u0000~\u007f\u0005"+
		"\u0016\u0000\u0000\u007f\u0081\u0003\u0018\f\u0000\u0080~\u0001\u0000"+
		"\u0000\u0000\u0081\u0084\u0001\u0000\u0000\u0000\u0082\u0080\u0001\u0000"+
		"\u0000\u0000\u0082\u0083\u0001\u0000\u0000\u0000\u0083\u0017\u0001\u0000"+
		"\u0000\u0000\u0084\u0082\u0001\u0000\u0000\u0000\u0085\u0086\u0007\u0001"+
		"\u0000\u0000\u0086\u0019\u0001\u0000\u0000\u0000\u0087\u0088\u0005\t\u0000"+
		"\u0000\u0088\u0089\u0003\u001e\u000f\u0000\u0089\u008a\u0005\n\u0000\u0000"+
		"\u008a\u008b\u0003\u0004\u0002\u0000\u008b\u008c\u0003\u001c\u000e\u0000"+
		"\u008c\u001b\u0001\u0000\u0000\u0000\u008d\u008e\u0005\u000b\u0000\u0000"+
		"\u008e\u0093\u0003 \u0010\u0000\u008f\u0090\u0005\f\u0000\u0000\u0090"+
		"\u0092\u0003 \u0010\u0000\u0091\u008f\u0001\u0000\u0000\u0000\u0092\u0095"+
		"\u0001\u0000\u0000\u0000\u0093\u0091\u0001\u0000\u0000\u0000\u0093\u0094"+
		"\u0001\u0000\u0000\u0000\u0094\u0097\u0001\u0000\u0000\u0000\u0095\u0093"+
		"\u0001\u0000\u0000\u0000\u0096\u008d\u0001\u0000\u0000\u0000\u0096\u0097"+
		"\u0001\u0000\u0000\u0000\u0097\u001d\u0001\u0000\u0000\u0000\u0098\u00a2"+
		"\u0005\u0018\u0000\u0000\u0099\u009e\u0003\u0006\u0003\u0000\u009a\u009b"+
		"\u0005\u0016\u0000\u0000\u009b\u009d\u0003\u0006\u0003\u0000\u009c\u009a"+
		"\u0001\u0000\u0000\u0000\u009d\u00a0\u0001\u0000\u0000\u0000\u009e\u009c"+
		"\u0001\u0000\u0000\u0000\u009e\u009f\u0001\u0000\u0000\u0000\u009f\u00a2"+
		"\u0001\u0000\u0000\u0000\u00a0\u009e\u0001\u0000\u0000\u0000\u00a1\u0098"+
		"\u0001\u0000\u0000\u0000\u00a1\u0099\u0001\u0000\u0000\u0000\u00a2\u001f"+
		"\u0001\u0000\u0000\u0000\u00a3\u00a6\u0003\"\u0011\u0000\u00a4\u00a6\u0003"+
		"$\u0012\u0000\u00a5\u00a3\u0001\u0000\u0000\u0000\u00a5\u00a4\u0001\u0000"+
		"\u0000\u0000\u00a6!\u0001\u0000\u0000\u0000\u00a7\u00a8\u0003\u0006\u0003"+
		"\u0000\u00a8\u00a9\u0005\u001c\u0000\u0000\u00a9\u00aa\u0003\u0018\f\u0000"+
		"\u00aa#\u0001\u0000\u0000\u0000\u00ab\u00ac\u0003\u0006\u0003\u0000\u00ac"+
		"\u00ad\u0005\u0012\u0000\u0000\u00ad\u00ae\u0003\u0018\f\u0000\u00ae\u00af"+
		"\u0005\f\u0000\u0000\u00af\u00b0\u0003\u0018\f\u0000\u00b0%\u0001\u0000"+
		"\u0000\u0000\u00b1\u00b2\u0005\u000e\u0000\u0000\u00b2\u00b3\u0003\u0004"+
		"\u0002\u0000\u00b3\u00b4\u0005\u000f\u0000\u0000\u00b4\u00b5\u0003(\u0014"+
		"\u0000\u00b5\u00b6\u0003\u001c\u000e\u0000\u00b6\'\u0001\u0000\u0000\u0000"+
		"\u00b7\u00bc\u0003*\u0015\u0000\u00b8\u00b9\u0005\u0016\u0000\u0000\u00b9"+
		"\u00bb\u0003*\u0015\u0000\u00ba\u00b8\u0001\u0000\u0000\u0000\u00bb\u00be"+
		"\u0001\u0000\u0000\u0000\u00bc\u00ba\u0001\u0000\u0000\u0000\u00bc\u00bd"+
		"\u0001\u0000\u0000\u0000\u00bd)\u0001\u0000\u0000\u0000\u00be\u00bc\u0001"+
		"\u0000\u0000\u0000\u00bf\u00c0\u0003\u0006\u0003\u0000\u00c0\u00c1\u0005"+
		"\u001c\u0000\u0000\u00c1\u00c2\u0003\u0018\f\u0000\u00c2+\u0001\u0000"+
		"\u0000\u0000\u00c3\u00c4\u0005\r\u0000\u0000\u00c4\u00c5\u0005\n\u0000"+
		"\u0000\u00c5\u00c6\u0003\u0004\u0002\u0000\u00c6\u00c7\u0003\u001c\u000e"+
		"\u0000\u00c7-\u0001\u0000\u0000\u0000\f3P[ez\u0082\u0093\u0096\u009e\u00a1"+
		"\u00a5\u00bc";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}