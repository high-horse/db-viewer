package pgxQueryParser

import (
	"db-viewer/internal/engine/entities"
	queryParaser "db-viewer/internal/engine/parser"

	pgquery "github.com/pganalyze/pg_query_go/v5"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(sql string) (*entities.SQLQueryEntity, error) {
	_, err := pgquery.Parse(sql)
	if err != nil {
		return nil, err
	}

	return &entities.SQLQueryEntity{
		RawSQL:        queryParaser.PaginatedSQLStmt(sql, 50),
		Dialect:       entities.DialectPostgreSQL,
		StatementType: queryParaser.ExtractStatementType(sql),
	}, nil
}
