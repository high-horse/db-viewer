package detector

import "strings"

type SQLite struct{}

func (d SQLite) Detect(statement string) StatementKind {
	trimmed := strings.TrimSpace(statement)
	if trimmed == "" {
		return KindExec
	}

	upper := strings.ToUpper(removeSqlComments(trimmed))

	fields := strings.Fields(upper)
	if len(fields) == 0 {
		return KindUnknown
	}
	firstWord := fields[0]

	switch firstWord {
	case "INSERT", "UPDATE", "DELETE", "REPLACE",
		"CREATE", "DROP", "ALTER", "TRUNCATE",
		"GRANT", "REVOKE", "BEGIN", "COMMIT", "ROLLBACK":
		return KindExec

	case "PRAGMA":
		if(strings.Contains(upper,"=")){
			return KindExec
		}
		return KindQuery
	default:
		// PRAGMA is ambiguous (some pragmas return rows, some don't) — stay safe.
		return KindQuery
	}
}

func removeSqlComments(stmt string) string {
	var source string
	for {
		source = strings.TrimSpace(stmt)

		if strings.HasPrefix(source, "--") {
			idx := strings.Index(source, "\n")
			if idx == -1 {
				idx = len(source)
			}
			stmt = source[idx+1:]
			continue
		}

		if strings.HasPrefix(source, "/*") {
			idx := strings.Index(source, "*/")
			if idx == -1 {
				idx = len(source)
			}
			stmt = source[idx+2:]
			continue
		}

		break
	}
	return source
}
