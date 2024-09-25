package types

import "time"

type Article struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Keywords    string     `json:"keywords"`
	Body        string     `json:"body"`
	Published   *time.Time `json:"published"`
	Format      string     `json:"format"`
	Type        string     `json:"type"`
}
