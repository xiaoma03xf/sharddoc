lexer grammar SQLLexer;

// 关键字
CREATE : 'CREATE';
TABLE : 'TABLE';
PRIMARY : 'PRIMARY';
KEY : 'KEY';
INDEX : 'INDEX';
INSERT : 'INSERT';
INTO : 'INTO';
VALUES : 'VALUES';
SELECT : 'SELECT';
FROM : 'FROM';
WHERE : 'WHERE';
AND : 'AND';
DELETE : 'DELETE';
UPDATE : 'UPDATE';
SET : 'SET';
// 数据类型
INT64 : 'INT64';
BYTES : 'BYTES';
BETWEEN : 'BETWEEN';

// 标识符（表名、列名等）
IDENTIFIER : [a-zA-Z_][a-zA-Z0-9_]*;

// 运算符、符号
LPAREN : '(';
RPAREN : ')';
COMMA : ',';
SEMICOLON : ';';
STAR : '*';

// 忽略空白符
WS : [ \t\r\n]+ -> skip;



// for add || update
INTEGER : [0-9]+;        // 匹配整数（如 1, 30）
STRING : '\'' (~'\'' | '\'\'')* '\'';  // 匹配字符串（如 'Alice'）


OP : '=' | '>' | '<' | '>=' | '<=' | '!=';
