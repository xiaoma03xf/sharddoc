// Code generated from SQLParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package ast // SQLParser
import "github.com/antlr4-go/antlr/v4"

// SQLParserListener is a complete listener for a parse tree produced by SQLParser.
type SQLParserListener interface {
	antlr.ParseTreeListener

	// EnterSql is called when entering the sql production.
	EnterSql(c *SqlContext)

	// EnterCreateTableStatement is called when entering the createTableStatement production.
	EnterCreateTableStatement(c *CreateTableStatementContext)

	// EnterTableName is called when entering the tableName production.
	EnterTableName(c *TableNameContext)

	// EnterColumnName is called when entering the columnName production.
	EnterColumnName(c *ColumnNameContext)

	// EnterColumnType is called when entering the columnType production.
	EnterColumnType(c *ColumnTypeContext)

	// EnterColumnDefinitions is called when entering the columnDefinitions production.
	EnterColumnDefinitions(c *ColumnDefinitionsContext)

	// EnterColumnDefinition is called when entering the columnDefinition production.
	EnterColumnDefinition(c *ColumnDefinitionContext)

	// EnterIndexDefinitions is called when entering the indexDefinitions production.
	EnterIndexDefinitions(c *IndexDefinitionsContext)

	// EnterIndexDefinition is called when entering the indexDefinition production.
	EnterIndexDefinition(c *IndexDefinitionContext)

	// EnterInsertTableStatement is called when entering the insertTableStatement production.
	EnterInsertTableStatement(c *InsertTableStatementContext)

	// EnterColumnInsertNames is called when entering the columnInsertNames production.
	EnterColumnInsertNames(c *ColumnInsertNamesContext)

	// EnterColumnInsertValues is called when entering the columnInsertValues production.
	EnterColumnInsertValues(c *ColumnInsertValuesContext)

	// EnterColumnValue is called when entering the columnValue production.
	EnterColumnValue(c *ColumnValueContext)

	// EnterSelectTableStatement is called when entering the selectTableStatement production.
	EnterSelectTableStatement(c *SelectTableStatementContext)

	// EnterSelectColumnNames is called when entering the selectColumnNames production.
	EnterSelectColumnNames(c *SelectColumnNamesContext)

	// EnterCondition is called when entering the condition production.
	EnterCondition(c *ConditionContext)

	// EnterComparisonCondition is called when entering the comparisonCondition production.
	EnterComparisonCondition(c *ComparisonConditionContext)

	// EnterBetweenCondition is called when entering the betweenCondition production.
	EnterBetweenCondition(c *BetweenConditionContext)

	// ExitSql is called when exiting the sql production.
	ExitSql(c *SqlContext)

	// ExitCreateTableStatement is called when exiting the createTableStatement production.
	ExitCreateTableStatement(c *CreateTableStatementContext)

	// ExitTableName is called when exiting the tableName production.
	ExitTableName(c *TableNameContext)

	// ExitColumnName is called when exiting the columnName production.
	ExitColumnName(c *ColumnNameContext)

	// ExitColumnType is called when exiting the columnType production.
	ExitColumnType(c *ColumnTypeContext)

	// ExitColumnDefinitions is called when exiting the columnDefinitions production.
	ExitColumnDefinitions(c *ColumnDefinitionsContext)

	// ExitColumnDefinition is called when exiting the columnDefinition production.
	ExitColumnDefinition(c *ColumnDefinitionContext)

	// ExitIndexDefinitions is called when exiting the indexDefinitions production.
	ExitIndexDefinitions(c *IndexDefinitionsContext)

	// ExitIndexDefinition is called when exiting the indexDefinition production.
	ExitIndexDefinition(c *IndexDefinitionContext)

	// ExitInsertTableStatement is called when exiting the insertTableStatement production.
	ExitInsertTableStatement(c *InsertTableStatementContext)

	// ExitColumnInsertNames is called when exiting the columnInsertNames production.
	ExitColumnInsertNames(c *ColumnInsertNamesContext)

	// ExitColumnInsertValues is called when exiting the columnInsertValues production.
	ExitColumnInsertValues(c *ColumnInsertValuesContext)

	// ExitColumnValue is called when exiting the columnValue production.
	ExitColumnValue(c *ColumnValueContext)

	// ExitSelectTableStatement is called when exiting the selectTableStatement production.
	ExitSelectTableStatement(c *SelectTableStatementContext)

	// ExitSelectColumnNames is called when exiting the selectColumnNames production.
	ExitSelectColumnNames(c *SelectColumnNamesContext)

	// ExitCondition is called when exiting the condition production.
	ExitCondition(c *ConditionContext)

	// ExitComparisonCondition is called when exiting the comparisonCondition production.
	ExitComparisonCondition(c *ComparisonConditionContext)

	// ExitBetweenCondition is called when exiting the betweenCondition production.
	ExitBetweenCondition(c *BetweenConditionContext)
}
