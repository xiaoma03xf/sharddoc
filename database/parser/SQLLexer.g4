lexer grammar SQLLexer;

// 关键字
CREATE : 'CREATE';
TABLE : 'TABLE';
PRIMARY : 'PRIMARY';
KEY : 'KEY';
INDEX : 'INDEX';

// 数据类型
INT64 : 'INT64';
BYTES : 'BYTES';

// 标识符（表名、列名等）
IDENTIFIER : [a-zA-Z_][a-zA-Z0-9_]*;

// 运算符、符号
LPAREN : '(';
RPAREN : ')';
COMMA : ',';
SEMICOLON : ';';

// 忽略空白符
WS : [ \t\r\n]+ -> skip;
