package sqlExecutor

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	"db-lens/internal/engine/entities"
)

// PageableSQL accepts a single SELECT or read-only CTE shape. Commands, locking
// reads, SELECT INTO, and data-modifying CTEs must keep execute-once semantics.
// This is a conservative capability check, not SQL validation or a security boundary.
func PageableSQL(query string) (string, bool) {
	query = strings.TrimSpace(query)
	var words []string
	end := len(query)
	terminated := false
	for i := 0; i < len(query); {
		if unicode.IsSpace(rune(query[i])) {
			i++
			continue
		}
		if strings.HasPrefix(query[i:], "--") {
			if n := strings.IndexByte(query[i:], '\n'); n >= 0 {
				i += n + 1
				continue
			}
			break
		}
		if strings.HasPrefix(query[i:], "/*") {
			// Nested/versioned comments and dollar quoting are deliberately unsupported.
			n := strings.Index(query[i+2:], "*/")
			if n < 0 || strings.Contains(query[i+2:i+2+n], "/*") || strings.HasPrefix(query[i:], "/*!") {
				return "", false
			}
			i += n + 4
			continue
		}
		if terminated {
			return "", false
		}
		ch := query[i]
		if ch == ';' {
			end = i
			terminated = true
			i++
			continue
		}
		if ch == '$' || ch == '#' || ch == '@' {
			return "", false
		}
		if ch == '\'' || ch == '"' || ch == '`' || ch == '[' {
			quote := ch
			if ch == '[' {
				quote = ']'
			}
			i++
			closed := false
			for i < len(query) {
				if query[i] == '\\' {
					return "", false
				}
				if query[i] == quote {
					if i+1 < len(query) && query[i+1] == quote {
						i += 2
						continue
					}
					i++
					closed = true
					break
				}
				i++
			}
			if !closed {
				return "", false
			}
			continue
		}
		if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_' {
			start := i
			for i < len(query) && (query[i] >= 'a' && query[i] <= 'z' || query[i] >= 'A' && query[i] <= 'Z' || query[i] >= '0' && query[i] <= '9' || query[i] == '_') {
				i++
			}
			words = append(words, strings.ToUpper(query[start:i]))
			continue
		}
		i++
	}
	if len(words) == 0 || words[0] != "SELECT" && words[0] != "WITH" {
		return "", false
	}
	for _, word := range words {
		switch word {
		case "INSERT", "UPDATE", "DELETE", "MERGE", "REPLACE", "INTO", "FOR", "LOCK", "CALL", "EXEC", "EXECUTE", "CREATE", "DROP", "ALTER", "PRAGMA", "ATTACH", "DETACH":
			return "", false
		}
	}
	return strings.TrimSpace(query[:end]), true
}

// OpenRead counts the result, then opens a row stream at the requested page.
// OFFSET is used only for random access; subsequent Next requests reuse the stream.
func (p *Pager) OpenRead(ctx context.Context, db *sql.DB, connection, dialect string, input entities.QueryInput) (*entities.QueryResult, error) {
	base, ok := PageableSQL(input.Query)
	if !ok {
		return nil, fmt.Errorf("this statement does not support counted page navigation or sorting")
	}
	size := input.PageSize
	if size == 0 {
		size = 100
	}
	if size < 1 || size > 500 {
		return nil, fmt.Errorf("page size must be between 1 and 500")
	}
	if input.SortColumn < 0 {
		return nil, fmt.Errorf("invalid sort column")
	}
	direction := strings.ToUpper(input.SortDirection)
	if input.SortColumn > 0 && direction != "ASC" && direction != "DESC" {
		return nil, fmt.Errorf("invalid sort direction")
	}
	wrapped := "(\n" + base + "\n) AS _dbv_result"
	var total int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+wrapped).Scan(&total); err != nil {
		return nil, fmt.Errorf("count result rows: %w", err)
	}
	page := int64(input.Page)
	lastPage := int64(1)
	if total > 0 {
		lastPage = (total-1)/int64(size) + 1
	}
	if page < 1 {
		page = 1
	}
	if page > lastPage {
		page = lastPage
	}
	offset := (page - 1) * int64(size)
	query := base
	if input.SortColumn > 0 {
		// Ordinals support duplicate output labels and cannot inject identifiers.
		query = fmt.Sprintf("SELECT * FROM %s ORDER BY %d %s", wrapped, input.SortColumn, direction)
	}
	if offset > 0 {
		if input.SortColumn == 0 {
			query = "SELECT * FROM " + wrapped
		}
		switch dialect {
		case "mysql":
			query += fmt.Sprintf("\nLIMIT 18446744073709551615 OFFSET %d", offset)
		case "sqlite":
			query += fmt.Sprintf("\nLIMIT -1 OFFSET %d", offset)
		default:
			query += fmt.Sprintf("\nOFFSET %d", offset)
		}
	}
	return p.OpenPage(ctx, db, connection, query, size, offset, total)
}
