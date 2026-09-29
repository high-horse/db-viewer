package mysqlQueryParser

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
		RawSQL:        sql,
		Dialect:       entities.DialectMySQL,
		StatementType: queryParaser.ExtractStatementType(sql),
	}, nil
}
