package db

import (
	"context"
	"database/sql"
	"time"
)

type QueryHistoryEntity struct {
	ID        int `json:"id"`
	ConnectionId string `json:"connection_id"`
	DatabaseName string    `json:"databaseName"`
	QueryText     string `json:"query_text"`
	ExecutedAt    int `json:"executed_at"`
	Duration    int `json:"duration"`
	Status        string `json:"status"`
}

type HistoryRepository struct {
	localDB *sql.DB
}


func NewHistoryRepository(localDB *sql.DB) *HistoryRepository {
	return &HistoryRepository{localDB: localDB}
}

func (r *HistoryRepository)  Log(
	ctx context.Context, entry QueryHistoryEntity,
) error {
	q := `
			INSERT INTO query_history (connection_id, database_name, query_text, duration_ms, status)
			VALUES (?, ?, ?, ?, ?)
		`
		_, err := r.localDB.ExecContext(ctx, q, 
			entry.ConnectionId, 
			entry.DatabaseName, 
			entry.QueryText, 
			entry.Duration, 
			entry.Status,
		)
		return err
}

func (r *HistoryRepository) GetHistory(
	ctx context.Context,
	limit int,
	since time.Time, 
) ([]QueryHistoryEntity, error) {
	q := `
		SELECT * FROM query_history
		WHERE executed_at >= ?
		ORDER BY executed_at DESC
		LIMIT ?
	`
	rows, err := r.localDB.QueryContext(ctx, q, int(since.Unix()), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var history []QueryHistoryEntity
	for rows.Next() {
		var entry QueryHistoryEntity
		if err := rows.Scan(&entry.ID, &entry.ConnectionId, &entry.DatabaseName, &entry.QueryText, &entry.ExecutedAt, &entry.Duration, &entry.Status); err != nil {
			return nil, err
		}
		history = append(history, entry)
	}
	return history, nil
}