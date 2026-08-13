package pgxQueryParser

import (
	"regexp"
	"strings"

	"db-viewer/internal/engine/entities"
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

	stmtType := extractPostgresStatementType(sql)

	return &entities.SQLQueryEntity{
		RawSQL:        sql,
		Dialect:       entities.DialectPostgreSQL,
		StatementType: stmtType,
	}, nil
}

func extractPostgresStatementType(sql string) entities.StatementType {
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return entities.StatementUnknown
	}

	upperSQL := strings.ToUpper(trimmed)

	re := regexp.MustCompile(`^\s*(\w+)`)
	matches := re.FindStringSubmatch(upperSQL)
	if len(matches) < 2 {
		return entities.StatementUnknown
	}

	keyword := matches[1]
	switch keyword {
	case "SELECT":
		return entities.StatementSelect
	case "INSERT":
		return entities.StatementInsert
	case "UPDATE":
		return entities.StatementUpdate
	case "DELETE":
		return entities.StatementDelete
	case "CREATE":
		return entities.StatementCreate
	case "ALTER":
		return entities.StatementAlter
	case "DROP":
		return entities.StatementDrop
	default:
		return entities.StatementUnknown
	}
}