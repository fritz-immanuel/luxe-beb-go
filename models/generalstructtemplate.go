package models

type IDNameTemplate struct {
	ID   string `json:"ID" db:"id"`
	Name string `json:"Name" db:"name"`
}
