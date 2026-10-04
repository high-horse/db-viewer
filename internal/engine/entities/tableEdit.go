package entities

import "encoding/json"

type TableRef struct {
	ConnectionID string `json:"connectionId"`
	Name         string `json:"name"`
	Schema       string `json:"schema"`
	Database     string `json:"database"`
}
type TableEditInfo struct {
	Table     TableRef            `json:"table"`
	Driver    string              `json:"driver"`
	Columns   []InspectColumnInfo `json:"columns"`
	Keys      []string            `json:"keys"`
	CanInsert bool                `json:"canInsert"`
	CanModify bool                `json:"canModify"`
	Reason    string              `json:"reason"`
}
type RowChange struct {
	Operation string                     `json:"operation"`
	Keys      map[string]json.RawMessage `json:"keys"`
	Values    map[string]json.RawMessage `json:"values"`
}
type TableChanges struct {
	Cursor  string      `json:"cursor"`
	Table   TableRef    `json:"table"`
	Changes []RowChange `json:"changes"`
}
type TableSaveResult struct {
	Applied int    `json:"applied"`
	Error   string `json:"error"`
}
