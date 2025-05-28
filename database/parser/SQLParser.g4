parser grammar SQLParser;

options {
  tokenVocab=SQLLexer;
}

// 根规则，表示整个 SQL 语句
sql : createTableStatement;

// CREATE TABLE 语句
createTableStatement : CREATE TABLE tableName LPAREN columnDefinitions indexDefinitions? RPAREN SEMICOLON;

// 表名
tableName : IDENTIFIER;

// 列定义
columnDefinitions : columnDefinition (COMMA columnDefinition)*;

// 单列定义
columnDefinition : columnName columnType;

// 列名
columnName : IDENTIFIER;

// 列类型
columnType : INT64 | BYTES | VARCHAR | TEXT | DATE | FLOAT;

// 索引定义
indexDefinitions : PRIMARY KEY LPAREN columnName RPAREN
                 | INDEX LPAREN columnName (COMMA columnName)* RPAREN;
