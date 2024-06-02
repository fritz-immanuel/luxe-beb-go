package models

import (
	"luxe-beb-go/library/types"
)

type BrandBulk struct {
	ID   string `json:"ID" db:"id"`
	Name string `json:"Name" db:"name" validate:"required"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`
}

type Brand struct {
	ID   string `json:"ID" db:"id"`
	Name string `json:"Name" db:"name" validate:"required"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`
}

type FindAllBrandParams struct {
	FindAllParams types.FindAllParams
}
