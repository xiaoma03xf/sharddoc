// Code generated from SQLParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package ast // SQLParser
import "github.com/antlr4-go/antlr/v4"

type BaseSQLParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSQLParserVisitor) VisitSql(ctx *SqlContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitCreateTableStatement(ctx *CreateTableStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitTableName(ctx *TableNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitColumnName(ctx *ColumnNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitColumnType(ctx *ColumnTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitColumnDefinitions(ctx *ColumnDefinitionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitColumnDefinition(ctx *ColumnDefinitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitIndexDefinitions(ctx *IndexDefinitionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitIndexDefinition(ctx *IndexDefinitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitInsertTableStatement(ctx *InsertTableStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitColumnInsertNames(ctx *ColumnInsertNamesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitColumnInsertValues(ctx *ColumnInsertValuesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitColumnValue(ctx *ColumnValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitSelectTableStatement(ctx *SelectTableStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitSelectColumnNames(ctx *SelectColumnNamesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitCondition(ctx *ConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitComparisonCondition(ctx *ComparisonConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitBetweenCondition(ctx *BetweenConditionContext) interface{} {
	return v.VisitChildren(ctx)
}
