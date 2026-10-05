package sqlExecutor

import (
	"context"
	"testing"

	"db-lens/internal/engine/entities"
)

func TestReadPagesCountJumpAndGlobalSort(t *testing.T) {
	p, db := testPager(t)
	ctx := context.Background()
	if _, err := db.Exec("CREATE TABLE items (id INTEGER PRIMARY KEY, score INTEGER); WITH RECURSIVE n(x) AS (VALUES(1) UNION ALL SELECT x+1 FROM n WHERE x<250) INSERT INTO items SELECT x, 251-x FROM n"); err != nil {
		t.Fatal(err)
	}
	input := entities.QueryInput{Query: "SELECT id, score FROM items ORDER BY id; -- keep comment", PageSize: 100, Page: 1}
	first, err := p.OpenRead(ctx, db, "test", "sqlite", input)
	if err != nil {
		t.Fatal(err)
	}
	if !first.CanNavigate || first.TotalRows == nil || *first.TotalRows != 250 || first.Rows[0][0] != int64(1) {
		t.Fatalf("bad first page: %+v", first)
	}
	second, err := p.Fetch(first.Cursor, "test", 2)
	if err != nil || second.StartRow != 101 || *second.TotalRows != 250 {
		t.Fatal("stream next/count failed", err)
	}
	p.Close(first.Cursor, "test")
	input.Page = 3
	last, err := p.OpenRead(ctx, db, "test", "sqlite", input)
	if err != nil || last.StartRow != 201 || len(last.Rows) != 50 || last.HasMore || last.Rows[0][0] != int64(201) {
		t.Fatal("direct last page failed", err)
	}
	input.Page = 2
	back, err := p.OpenRead(ctx, db, "test", "sqlite", input)
	if err != nil || back.StartRow != 101 {
		t.Fatal("backwards offset jump failed", err)
	}
	last, err = p.Fetch(back.Cursor, "test", 3)
	if err != nil || last.StartRow != 201 {
		t.Fatal("stream sequence after jump failed", err)
	}
	input.Page = 1
	input.SortColumn = 2
	input.SortDirection = "asc"
	sorted, err := p.OpenRead(ctx, db, "test", "sqlite", input)
	if err != nil {
		t.Fatal(err)
	}
	// The lowest scores originate beyond the first unsorted page.
	if sorted.Rows[0][0] != int64(250) || sorted.Rows[99][0] != int64(151) {
		t.Fatal("sort only applied to first page")
	}
	p.Close(sorted.Cursor, "test")
	input.Query = "SELECT id, score FROM items ORDER BY id LIMIT 205"
	input.Page = 3
	last, err = p.OpenRead(ctx, db, "test", "sqlite", input)
	if err != nil || *last.TotalRows != 205 || len(last.Rows) != 5 || last.Rows[0][0] != int64(5) {
		t.Fatal("original limit not preserved with global sort", err)
	}
}

func TestReadPageBoundaryAndValidation(t *testing.T) {
	p, db := testPager(t)
	for _, query := range []string{"SELECT 1 AS value WHERE false", "SELECT 1 AS value"} {
		result, err := p.OpenRead(context.Background(), db, "test", "sqlite", entities.QueryInput{Query: query, Page: 999, PageSize: 100})
		if err != nil || result.StartRow != 1 || result.HasMore || result.TotalRows == nil {
			t.Fatal("empty/single page failed", err)
		}
	}
	for _, input := range []entities.QueryInput{
		{Query: "DELETE FROM items RETURNING id"},
		{Query: "SELECT 1", SortColumn: 1, SortDirection: "desc; DROP TABLE items"},
		{Query: "SELECT 1", SortColumn: -1},
		{Query: "SELECT 1", SortColumn: 2, SortDirection: "asc"},
		{Query: "SELECT 1", PageSize: 501},
	} {
		if _, err := p.OpenRead(context.Background(), db, "test", "sqlite", input); err == nil {
			t.Fatalf("accepted invalid input: %+v", input)
		}
	}
}

func TestPageableSQL(t *testing.T) {
	for _, query := range []string{
		"SELECT * FROM items; -- comment", "/* hi */ SELECT 'a;b'", "SELECT \"update\" FROM t",
		"WITH n AS (SELECT 1) SELECT * FROM n", "SELECT 'it''s fine'", "SELECT `odd;name` FROM t",
	} {
		if _, ok := PageableSQL(query); !ok {
			t.Errorf("rejected read: %s", query)
		}
	}
	for _, query := range []string{
		"INSERT INTO t VALUES(1) RETURNING id", "SELECT 1; DELETE FROM t", "SELECT * INTO temp FROM t",
		"WITH deleted AS (DELETE FROM t RETURNING *) SELECT * FROM deleted", "SELECT * FROM t FOR UPDATE",
		"SELECT * FROM t LOCK IN SHARE MODE", "SELECT 1 /*! INTO OUTFILE 'x' */", "SELECT $$text$$",
		"SELECT 'unterminated", "/* unterminated", "SELECT @value := 1", "SELECT 1 /* nested /* comment */ */",
	} {
		if _, ok := PageableSQL(query); ok {
			t.Errorf("accepted non-pageable statement: %s", query)
		}
	}
}
