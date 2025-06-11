// Code generated from SQLParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package ast // SQLParser
import "github.com/antlr4-go/antlr/v4"

// BaseSQLParserListener is a complete listener for a parse tree produced by SQLParser.
type BaseSQLParserListener struct{}

var _ SQLParserListener = &BaseSQLParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSQLParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSQLParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSQLParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSQLParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSql is called when production sql is entered.
func (s *BaseSQLParserListener) EnterSql(ctx *SqlContext) {}

// ExitSql is called when production sql is exited.
func (s *BaseSQLParserListener) ExitSql(ctx *SqlContext) {}

// EnterCreateTableStatement is called when production createTableStatement is entered.
func (s *BaseSQLParserListener) EnterCreateTableStatement(ctx *CreateTableStatementContext) {}

// ExitCreateTableStatement is called when production createTableStatement is exited.
func (s *BaseSQLParserListener) ExitCreateTableStatement(ctx *CreateTableStatementContext) {}

// EnterTableName is called when production tableName is entered.
func (s *BaseSQLParserListener) EnterTableName(ctx *TableNameContext) {}

// ExitTableName is called when production tableName is exited.
func (s *BaseSQLParserListener) ExitTableName(ctx *TableNameContext) {}

// EnterColumnName is called when production columnName is entered.
func (s *BaseSQLParserListener) EnterColumnName(ctx *ColumnNameContext) {}

// ExitColumnName is called when production columnName is exited.
func (s *BaseSQLParserListener) ExitColumnName(ctx *ColumnNameContext) {}

// EnterColumnType is called when production columnType is entered.
func (s *BaseSQLParserListener) EnterColumnType(ctx *ColumnTypeContext) {}

// ExitColumnType is called when production columnType is exited.
func (s *BaseSQLParserListener) ExitColumnType(ctx *ColumnTypeContext) {}

// EnterColumnDefinitions is called when production columnDefinitions is entered.
func (s *BaseSQLParserListener) EnterColumnDefinitions(ctx *ColumnDefinitionsContext) {}

// ExitColumnDefinitions is called when production columnDefinitions is exited.
func (s *BaseSQLParserListener) ExitColumnDefinitions(ctx *ColumnDefinitionsContext) {}

// EnterColumnDefinition is called when production columnDefinition is entered.
func (s *BaseSQLParserListener) EnterColumnDefinition(ctx *ColumnDefinitionContext) {}

// ExitColumnDefinition is called when production columnDefinition is exited.
func (s *BaseSQLParserListener) ExitColumnDefinition(ctx *ColumnDefinitionContext) {}

// EnterIndexDefinitions is called when production indexDefinitions is entered.
func (s *BaseSQLParserListener) EnterIndexDefinitions(ctx *IndexDefinitionsContext) {}

// ExitIndexDefinitions is called when production indexDefinitions is exited.
func (s *BaseSQLParserListener) ExitIndexDefinitions(ctx *IndexDefinitionsContext) {}

// EnterIndexDefinition is called when production indexDefinition is entered.
func (s *BaseSQLParserListener) EnterIndexDefinition(ctx *IndexDefinitionContext) {}

// ExitIndexDefinition is called when production indexDefinition is exited.
func (s *BaseSQLParserListener) ExitIndexDefinition(ctx *IndexDefinitionContext) {}

// EnterInsertTableStatement is called when production insertTableStatement is entered.
func (s *BaseSQLParserListener) EnterInsertTableStatement(ctx *InsertTableStatementContext) {}

// ExitInsertTableStatement is called when production insertTableStatement is exited.
func (s *BaseSQLParserListener) ExitInsertTableStatement(ctx *InsertTableStatementContext) {}

// EnterColumnInsertNames is called when production columnInsertNames is entered.
func (s *BaseSQLParserListener) EnterColumnInsertNames(ctx *ColumnInsertNamesContext) {}

// ExitColumnInsertNames is called when production columnInsertNames is exited.
func (s *BaseSQLParserListener) ExitColumnInsertNames(ctx *ColumnInsertNamesContext) {}

// EnterColumnInsertValues is called when production columnInsertValues is entered.
func (s *BaseSQLParserListener) EnterColumnInsertValues(ctx *ColumnInsertValuesContext) {}

// ExitColumnInsertValues is called when production columnInsertValues is exited.
func (s *BaseSQLParserListener) ExitColumnInsertValues(ctx *ColumnInsertValuesContext) {}

// EnterColumnValue is called when production columnValue is entered.
func (s *BaseSQLParserListener) EnterColumnValue(ctx *ColumnValueContext) {}

// ExitColumnValue is called when production columnValue is exited.
func (s *BaseSQLParserListener) ExitColumnValue(ctx *ColumnValueContext) {}

// EnterSelectTableStatement is called when production selectTableStatement is entered.
func (s *BaseSQLParserListener) EnterSelectTableStatement(ctx *SelectTableStatementContext) {}

// ExitSelectTableStatement is called when production selectTableStatement is exited.
func (s *BaseSQLParserListener) ExitSelectTableStatement(ctx *SelectTableStatementContext) {}

// EnterConditions is called when production conditions is entered.
func (s *BaseSQLParserListener) EnterConditions(ctx *ConditionsContext) {}

// ExitConditions is called when production conditions is exited.
func (s *BaseSQLParserListener) ExitConditions(ctx *ConditionsContext) {}

// EnterSelectColumnNames is called when production selectColumnNames is entered.
func (s *BaseSQLParserListener) EnterSelectColumnNames(ctx *SelectColumnNamesContext) {}

// ExitSelectColumnNames is called when production selectColumnNames is exited.
func (s *BaseSQLParserListener) ExitSelectColumnNames(ctx *SelectColumnNamesContext) {}

// EnterCondition is called when production condition is entered.
func (s *BaseSQLParserListener) EnterCondition(ctx *ConditionContext) {}

// ExitCondition is called when production condition is exited.
func (s *BaseSQLParserListener) ExitCondition(ctx *ConditionContext) {}

// EnterComparisonCondition is called when production comparisonCondition is entered.
func (s *BaseSQLParserListener) EnterComparisonCondition(ctx *ComparisonConditionContext) {}

// ExitComparisonCondition is called when production comparisonCondition is exited.
func (s *BaseSQLParserListener) ExitComparisonCondition(ctx *ComparisonConditionContext) {}

// EnterBetweenCondition is called when production betweenCondition is entered.
func (s *BaseSQLParserListener) EnterBetweenCondition(ctx *BetweenConditionContext) {}

// ExitBetweenCondition is called when production betweenCondition is exited.
func (s *BaseSQLParserListener) ExitBetweenCondition(ctx *BetweenConditionContext) {}

// EnterUpdateTableStatement is called when production updateTableStatement is entered.
func (s *BaseSQLParserListener) EnterUpdateTableStatement(ctx *UpdateTableStatementContext) {}

// ExitUpdateTableStatement is called when production updateTableStatement is exited.
func (s *BaseSQLParserListener) ExitUpdateTableStatement(ctx *UpdateTableStatementContext) {}

// EnterSetClauses is called when production setClauses is entered.
func (s *BaseSQLParserListener) EnterSetClauses(ctx *SetClausesContext) {}

// ExitSetClauses is called when production setClauses is exited.
func (s *BaseSQLParserListener) ExitSetClauses(ctx *SetClausesContext) {}

// EnterSetClause is called when production setClause is entered.
func (s *BaseSQLParserListener) EnterSetClause(ctx *SetClauseContext) {}

// ExitSetClause is called when production setClause is exited.
func (s *BaseSQLParserListener) ExitSetClause(ctx *SetClauseContext) {}

// EnterDeleteTableStatement is called when production deleteTableStatement is entered.
func (s *BaseSQLParserListener) EnterDeleteTableStatement(ctx *DeleteTableStatementContext) {}

// ExitDeleteTableStatement is called when production deleteTableStatement is exited.
func (s *BaseSQLParserListener) ExitDeleteTableStatement(ctx *DeleteTableStatementContext) {}
