package mysqlQueryParser

import (
	"regexp"
	"strings"

	"db-viewer/internal/engine/entities"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(sql string) (*entities.SQLQueryEntity, error) {
	stmtType := extractMySQLStatementType(sql)

	return &entities.SQLQueryEntity{
		RawSQL:        sql,
		Dialect:       entities.DialectMySQL,
		StatementType: stmtType,
	}, nil
}

func extractMySQLStatementType(sql string) entities.StatementType {
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