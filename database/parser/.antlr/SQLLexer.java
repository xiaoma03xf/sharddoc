// Generated from /home/dengjun/sharddoc/database/parser/SQLLexer.g4 by ANTLR 4.13.1
import org.antlr.v4.runtime.Lexer;
import org.antlr.v4.runtime.CharStream;
import org.antlr.v4.runtime.Token;
import org.antlr.v4.runtime.TokenStream;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.misc.*;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast", "CheckReturnValue", "this-escape"})
public class SQLLexer extends Lexer {
	static { RuntimeMetaData.checkVersion("4.13.1", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		CREATE=1, TABLE=2, PRIMARY=3, KEY=4, INDEX=5, INSERT=6, INTO=7, VALUES=8, 
		INT64=9, BYTES=10, IDENTIFIER=11, LPAREN=12, RPAREN=13, COMMA=14, SEMICOLON=15, 
		WS=16, INTEGER=17, STRING=18;
	public static String[] channelNames = {
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN"
	};

	public static String[] modeNames = {
		"DEFAULT_MODE"
	};

	private static String[] makeRuleNames() {
		return new String[] {
			"CREATE", "TABLE", "PRIMARY", "KEY", "INDEX", "INSERT", "INTO", "VALUES", 
			"INT64", "BYTES", "IDENTIFIER", "LPAREN", "RPAREN", "COMMA", "SEMICOLON", 
			"WS", "INTEGER", "STRING"
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


	public SQLLexer(CharStream input) {
		super(input);
		_interp = new LexerATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@Override
	public String getGrammarFileName() { return "SQLLexer.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public String[] getChannelNames() { return channelNames; }

	@Override
	public String[] getModeNames() { return modeNames; }

	@Override
	public ATN getATN() { return _ATN; }

	public static final String _serializedATN =
		"\u0004\u0000\u0012\u0089\u0006\uffff\uffff\u0002\u0000\u0007\u0000\u0002"+
		"\u0001\u0007\u0001\u0002\u0002\u0007\u0002\u0002\u0003\u0007\u0003\u0002"+
		"\u0004\u0007\u0004\u0002\u0005\u0007\u0005\u0002\u0006\u0007\u0006\u0002"+
		"\u0007\u0007\u0007\u0002\b\u0007\b\u0002\t\u0007\t\u0002\n\u0007\n\u0002"+
		"\u000b\u0007\u000b\u0002\f\u0007\f\u0002\r\u0007\r\u0002\u000e\u0007\u000e"+
		"\u0002\u000f\u0007\u000f\u0002\u0010\u0007\u0010\u0002\u0011\u0007\u0011"+
		"\u0001\u0000\u0001\u0000\u0001\u0000\u0001\u0000\u0001\u0000\u0001\u0000"+
		"\u0001\u0000\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002"+
		"\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0003\u0001\u0003\u0001\u0003"+
		"\u0001\u0003\u0001\u0004\u0001\u0004\u0001\u0004\u0001\u0004\u0001\u0004"+
		"\u0001\u0004\u0001\u0005\u0001\u0005\u0001\u0005\u0001\u0005\u0001\u0005"+
		"\u0001\u0005\u0001\u0005\u0001\u0006\u0001\u0006\u0001\u0006\u0001\u0006"+
		"\u0001\u0006\u0001\u0007\u0001\u0007\u0001\u0007\u0001\u0007\u0001\u0007"+
		"\u0001\u0007\u0001\u0007\u0001\b\u0001\b\u0001\b\u0001\b\u0001\b\u0001"+
		"\b\u0001\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001\t\u0001\n\u0001\n\u0005"+
		"\nf\b\n\n\n\f\ni\t\n\u0001\u000b\u0001\u000b\u0001\f\u0001\f\u0001\r\u0001"+
		"\r\u0001\u000e\u0001\u000e\u0001\u000f\u0004\u000ft\b\u000f\u000b\u000f"+
		"\f\u000fu\u0001\u000f\u0001\u000f\u0001\u0010\u0004\u0010{\b\u0010\u000b"+
		"\u0010\f\u0010|\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0005"+
		"\u0011\u0083\b\u0011\n\u0011\f\u0011\u0086\t\u0011\u0001\u0011\u0001\u0011"+
		"\u0000\u0000\u0012\u0001\u0001\u0003\u0002\u0005\u0003\u0007\u0004\t\u0005"+
		"\u000b\u0006\r\u0007\u000f\b\u0011\t\u0013\n\u0015\u000b\u0017\f\u0019"+
		"\r\u001b\u000e\u001d\u000f\u001f\u0010!\u0011#\u0012\u0001\u0000\u0005"+
		"\u0003\u0000AZ__az\u0004\u000009AZ__az\u0003\u0000\t\n\r\r  \u0001\u0000"+
		"09\u0001\u0000\'\'\u008d\u0000\u0001\u0001\u0000\u0000\u0000\u0000\u0003"+
		"\u0001\u0000\u0000\u0000\u0000\u0005\u0001\u0000\u0000\u0000\u0000\u0007"+
		"\u0001\u0000\u0000\u0000\u0000\t\u0001\u0000\u0000\u0000\u0000\u000b\u0001"+
		"\u0000\u0000\u0000\u0000\r\u0001\u0000\u0000\u0000\u0000\u000f\u0001\u0000"+
		"\u0000\u0000\u0000\u0011\u0001\u0000\u0000\u0000\u0000\u0013\u0001\u0000"+
		"\u0000\u0000\u0000\u0015\u0001\u0000\u0000\u0000\u0000\u0017\u0001\u0000"+
		"\u0000\u0000\u0000\u0019\u0001\u0000\u0000\u0000\u0000\u001b\u0001\u0000"+
		"\u0000\u0000\u0000\u001d\u0001\u0000\u0000\u0000\u0000\u001f\u0001\u0000"+
		"\u0000\u0000\u0000!\u0001\u0000\u0000\u0000\u0000#\u0001\u0000\u0000\u0000"+
		"\u0001%\u0001\u0000\u0000\u0000\u0003,\u0001\u0000\u0000\u0000\u00052"+
		"\u0001\u0000\u0000\u0000\u0007:\u0001\u0000\u0000\u0000\t>\u0001\u0000"+
		"\u0000\u0000\u000bD\u0001\u0000\u0000\u0000\rK\u0001\u0000\u0000\u0000"+
		"\u000fP\u0001\u0000\u0000\u0000\u0011W\u0001\u0000\u0000\u0000\u0013]"+
		"\u0001\u0000\u0000\u0000\u0015c\u0001\u0000\u0000\u0000\u0017j\u0001\u0000"+
		"\u0000\u0000\u0019l\u0001\u0000\u0000\u0000\u001bn\u0001\u0000\u0000\u0000"+
		"\u001dp\u0001\u0000\u0000\u0000\u001fs\u0001\u0000\u0000\u0000!z\u0001"+
		"\u0000\u0000\u0000#~\u0001\u0000\u0000\u0000%&\u0005C\u0000\u0000&\'\u0005"+
		"R\u0000\u0000\'(\u0005E\u0000\u0000()\u0005A\u0000\u0000)*\u0005T\u0000"+
		"\u0000*+\u0005E\u0000\u0000+\u0002\u0001\u0000\u0000\u0000,-\u0005T\u0000"+
		"\u0000-.\u0005A\u0000\u0000./\u0005B\u0000\u0000/0\u0005L\u0000\u0000"+
		"01\u0005E\u0000\u00001\u0004\u0001\u0000\u0000\u000023\u0005P\u0000\u0000"+
		"34\u0005R\u0000\u000045\u0005I\u0000\u000056\u0005M\u0000\u000067\u0005"+
		"A\u0000\u000078\u0005R\u0000\u000089\u0005Y\u0000\u00009\u0006\u0001\u0000"+
		"\u0000\u0000:;\u0005K\u0000\u0000;<\u0005E\u0000\u0000<=\u0005Y\u0000"+
		"\u0000=\b\u0001\u0000\u0000\u0000>?\u0005I\u0000\u0000?@\u0005N\u0000"+
		"\u0000@A\u0005D\u0000\u0000AB\u0005E\u0000\u0000BC\u0005X\u0000\u0000"+
		"C\n\u0001\u0000\u0000\u0000DE\u0005I\u0000\u0000EF\u0005N\u0000\u0000"+
		"FG\u0005S\u0000\u0000GH\u0005E\u0000\u0000HI\u0005R\u0000\u0000IJ\u0005"+
		"T\u0000\u0000J\f\u0001\u0000\u0000\u0000KL\u0005I\u0000\u0000LM\u0005"+
		"N\u0000\u0000MN\u0005T\u0000\u0000NO\u0005O\u0000\u0000O\u000e\u0001\u0000"+
		"\u0000\u0000PQ\u0005V\u0000\u0000QR\u0005A\u0000\u0000RS\u0005L\u0000"+
		"\u0000ST\u0005U\u0000\u0000TU\u0005E\u0000\u0000UV\u0005S\u0000\u0000"+
		"V\u0010\u0001\u0000\u0000\u0000WX\u0005I\u0000\u0000XY\u0005N\u0000\u0000"+
		"YZ\u0005T\u0000\u0000Z[\u00056\u0000\u0000[\\\u00054\u0000\u0000\\\u0012"+
		"\u0001\u0000\u0000\u0000]^\u0005B\u0000\u0000^_\u0005Y\u0000\u0000_`\u0005"+
		"T\u0000\u0000`a\u0005E\u0000\u0000ab\u0005S\u0000\u0000b\u0014\u0001\u0000"+
		"\u0000\u0000cg\u0007\u0000\u0000\u0000df\u0007\u0001\u0000\u0000ed\u0001"+
		"\u0000\u0000\u0000fi\u0001\u0000\u0000\u0000ge\u0001\u0000\u0000\u0000"+
		"gh\u0001\u0000\u0000\u0000h\u0016\u0001\u0000\u0000\u0000ig\u0001\u0000"+
		"\u0000\u0000jk\u0005(\u0000\u0000k\u0018\u0001\u0000\u0000\u0000lm\u0005"+
		")\u0000\u0000m\u001a\u0001\u0000\u0000\u0000no\u0005,\u0000\u0000o\u001c"+
		"\u0001\u0000\u0000\u0000pq\u0005;\u0000\u0000q\u001e\u0001\u0000\u0000"+
		"\u0000rt\u0007\u0002\u0000\u0000sr\u0001\u0000\u0000\u0000tu\u0001\u0000"+
		"\u0000\u0000us\u0001\u0000\u0000\u0000uv\u0001\u0000\u0000\u0000vw\u0001"+
		"\u0000\u0000\u0000wx\u0006\u000f\u0000\u0000x \u0001\u0000\u0000\u0000"+
		"y{\u0007\u0003\u0000\u0000zy\u0001\u0000\u0000\u0000{|\u0001\u0000\u0000"+
		"\u0000|z\u0001\u0000\u0000\u0000|}\u0001\u0000\u0000\u0000}\"\u0001\u0000"+
		"\u0000\u0000~\u0084\u0005\'\u0000\u0000\u007f\u0083\b\u0004\u0000\u0000"+
		"\u0080\u0081\u0005\'\u0000\u0000\u0081\u0083\u0005\'\u0000\u0000\u0082"+
		"\u007f\u0001\u0000\u0000\u0000\u0082\u0080\u0001\u0000\u0000\u0000\u0083"+
		"\u0086\u0001\u0000\u0000\u0000\u0084\u0082\u0001\u0000\u0000\u0000\u0084"+
		"\u0085\u0001\u0000\u0000\u0000\u0085\u0087\u0001\u0000\u0000\u0000\u0086"+
		"\u0084\u0001\u0000\u0000\u0000\u0087\u0088\u0005\'\u0000\u0000\u0088$"+
		"\u0001\u0000\u0000\u0000\u0006\u0000gu|\u0082\u0084\u0001\u0006\u0000"+
		"\u0000";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}