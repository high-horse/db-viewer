package sqliteQueryParser

import (
	"db-lens/internal/engine/entities"
	queryParaser "db-lens/internal/engine/parser"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(sql string) (*entities.SQLQueryEntity, error) {
	return &entities.SQLQueryEntity{
		RawSQL:        sql,
		Dialect:       entities.DialectSQLite,
		StatementType: queryParaser.ExtractStatementType(sql),
	}, nil
}
