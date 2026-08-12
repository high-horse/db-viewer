package queryParaser

import "db-viewer/internal/engine/entities"

type Parser interface {
	Parse(query string) (*entities.SQLQueryEntity, error)
}
