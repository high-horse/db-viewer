package queryParaser

import (
	"db-viewer/internal/engine/entities"
	"fmt"
	"regexp"
	"strings"
)

type Parser interface {
	Parse(query string) (*entities.SQLQueryEntity, error)
}

func ExtractStatementType(sql string) entities.StatementType {
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

func stripTrailingSemicolon(sql string) string {
	sql = strings.TrimSpace(sql)
	for strings.HasSuffix(sql, ";") {
		sql = strings.TrimSuffix(sql, ";")
		sql = strings.TrimSpace(sql)
	}
	return sql
}

func WrapSqlStmt(sql string) string {
	sql = stripTrailingSemicolon(sql)
	return "SELECT * FROM (" + sql + ") AS subquery "
}

func CountSQLStmt(sql string) string {
	sql = stripTrailingSemicolon(sql)
	return "SELECT COUNT(*) FROM (" + sql + ") AS subquery"
}

func PaginatedSQLStmt(sql string, limit int) string {
	if limit <= 0 {
		limit = 50
	}
	limit = 50
	sql = stripTrailingSemicolon(sql)
	return fmt.Sprintf(
		"SELECT * FROM (%s) AS subquery LIMIT %d",
		sql,
		limit,
	)
}