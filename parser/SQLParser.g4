parser grammar SQLParser;

options {
  tokenVocab=SQLLexer;
}


sql 
  : createTableStatement 
  | insertTableStatement 
  | selectTableStatement
  | updateTableStatement
  | deleteTableStatement
  ; 

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
    : INTEGER    
    | STRING      
    ;


// Select 语句
selectTableStatement
  : SELECT selectColumnNames FROM tableName conditions
  ;

  conditions
  : (WHERE condition (AND condition)*)?
  ;

  selectColumnNames
  : STAR 
  | columnName (COMMA columnName)*
  ;

  // basic commpare and between compare
  condition
  : comparisonCondition
  | betweenCondition
  ;
  comparisonCondition
  : columnName OP columnValue
  ;
  betweenCondition
  : columnName BETWEEN columnValue AND columnValue
  ;

// update 语句
updateTableStatement
  : UPDATE tableName SET setClauses conditions
  ;
setClauses
  : setClause (COMMA setClause)*
  ;

setClause
  : columnName OP columnValue
  ;

// delete 语句
deleteTableStatement
  : DELETE FROM tableName conditions
  ;

