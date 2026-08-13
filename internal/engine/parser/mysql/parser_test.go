package mysqlQueryParser

import (
	"testing"

	"db-viewer/internal/engine/entities"
)

func TestMySQLParser(t *testing.T) {
	p := NewParser()
	tests := []struct {
		sql  string
		want entities.StatementType
	}{
		{"SELECT * FROM users", entities.StatementSelect},
		{"INSERT INTO users VALUES (1, 'john')", entities.StatementInsert},
		{"UPDATE users SET name='jane' WHERE id=1", entities.StatementUpdate},
		{"DELETE FROM users WHERE id=1", entities.StatementDelete},
		{"CREATE TABLE users (id INT)", entities.StatementCreate},
		{"ALTER TABLE users ADD COLUMN age INT", entities.StatementAlter},
		{"DROP TABLE users", entities.StatementDrop},
	}

	for _, tt := range tests {
		testName := tt.sql
		if len(testName) > 30 {
			testName = testName[:30]
		}
		t.Run(testName, func(t *testing.T) {
			result, err := p.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if result.StatementType != tt.want {
				t.Errorf("got %v, want %v", result.StatementType, tt.want)
			}
			if result.Dialect != entities.DialectMySQL {
				t.Errorf("dialect: got %v, want %v", result.Dialect, entities.DialectMySQL)
			}
		})
	}
}
