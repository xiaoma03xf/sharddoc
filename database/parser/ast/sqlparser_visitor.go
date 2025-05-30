// Code generated from SQLParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package ast // SQLParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by SQLParser.
type SQLParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SQLParser#sql.
	VisitSql(ctx *SqlContext) interface{}

	// Visit a parse tree produced by SQLParser#createTableStatement.
	VisitCreateTableStatement(ctx *CreateTableStatementContext) interface{}

	// Visit a parse tree produced by SQLParser#tableName.
	VisitTableName(ctx *TableNameContext) interface{}

	// Visit a parse tree produced by SQLParser#columnName.
	VisitColumnName(ctx *ColumnNameContext) interface{}

	// Visit a parse tree produced by SQLParser#columnType.
	VisitColumnType(ctx *ColumnTypeContext) interface{}

	// Visit a parse tree produced by SQLParser#columnDefinitions.
	VisitColumnDefinitions(ctx *ColumnDefinitionsContext) interface{}

	// Visit a parse tree produced by SQLParser#columnDefinition.
	VisitColumnDefinition(ctx *ColumnDefinitionContext) interface{}

	// Visit a parse tree produced by SQLParser#indexDefinitions.
	VisitIndexDefinitions(ctx *IndexDefinitionsContext) interface{}

	// Visit a parse tree produced by SQLParser#indexDefinition.
	VisitIndexDefinition(ctx *IndexDefinitionContext) interface{}

	// Visit a parse tree produced by SQLParser#insertTableStatement.
	VisitInsertTableStatement(ctx *InsertTableStatementContext) interface{}

	// Visit a parse tree produced by SQLParser#columnInsertNames.
	VisitColumnInsertNames(ctx *ColumnInsertNamesContext) interface{}

	// Visit a parse tree produced by SQLParser#columnInsertValues.
	VisitColumnInsertValues(ctx *ColumnInsertValuesContext) interface{}

	// Visit a parse tree produced by SQLParser#columnValue.
	VisitColumnValue(ctx *ColumnValueContext) interface{}

	// Visit a parse tree produced by SQLParser#selectTableStatement.
	VisitSelectTableStatement(ctx *SelectTableStatementContext) interface{}

	// Visit a parse tree produced by SQLParser#selectColumnNames.
	VisitSelectColumnNames(ctx *SelectColumnNamesContext) interface{}

	// Visit a parse tree produced by SQLParser#condition.
	VisitCondition(ctx *ConditionContext) interface{}

	// Visit a parse tree produced by SQLParser#comparisonCondition.
	VisitComparisonCondition(ctx *ComparisonConditionContext) interface{}

	// Visit a parse tree produced by SQLParser#betweenCondition.
	VisitBetweenCondition(ctx *BetweenConditionContext) interface{}
}
