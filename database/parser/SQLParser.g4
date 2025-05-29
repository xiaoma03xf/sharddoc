parser grammar SQLParser;

options {
  tokenVocab=SQLLexer;
}


sql : createTableStatement | insertTableStatement ; 

// CREATE TABLE 语句
createTableStatement : 
    CREATE TABLE tableName LPAREN     
    columnDefinitions COMMA
    PRIMARY KEY LPAREN columnName RPAREN COMMA    
    indexDefinitions         
    RPAREN SEMICOLON;          

tableName : IDENTIFIER;
columnName : IDENTIFIER;
columnType : INT64 | BYTES;

columnDefinitions : columnDefinition (COMMA columnDefinition)*;
columnDefinition : columnName columnType;

indexDefinitions : indexDefinition (COMMA indexDefinition)*;
indexDefinition : INDEX LPAREN columnName (COMMA columnName)* RPAREN;


// Insert 语句
insertTableStatement
  : INSERT INTO tableName LPAREN columnInsertNames RPAREN VALUES LPAREN columnInsertValues RPAREN 
  ;

columnInsertNames
  : columnName (COMMA columnName)*
  ;
columnInsertValues
  : columnValue (COMMA columnValue)*
  ;
columnValue
    : INTEGER       // 对应 INT64 列
    | STRING        // 对应 BYTES 列
    ;