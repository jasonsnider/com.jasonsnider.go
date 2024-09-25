package types

import (
	"database/sql"
)

type Article struct {
	ID          string         `json:"id"`
	Slug        string         `json:"slug"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Keywords    sql.NullString `json:"keywords"`
	Body        sql.NullString `json:"body"`
	Published   sql.NullTime   `json:"published"`
	Format      sql.NullString `json:"format"`
	Type        sql.NullString `json:"type"`
}
