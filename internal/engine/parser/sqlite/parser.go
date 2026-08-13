package sqliteQueryParser

import (
	"db-viewer/internal/engine/entities"
	queryParaser "db-viewer/internal/engine/parser"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(sql string) (*entities.SQLQueryEntity, error) {
	return &entities.SQLQueryEntity{
		RawSQL:        queryParaser.PaginatedSQLStmt(sql, 50),
		Dialect:       entities.DialectSQLite,
		StatementType: queryParaser.ExtractStatementType(sql),
	}, nil
}
