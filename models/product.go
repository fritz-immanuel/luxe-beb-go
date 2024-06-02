package models

import (
	"luxe-beb-go/library/types"
)

type ProductBulk struct {
	ID          string  `json:"ID" db:"id"`
	Code        string  `json:"Code" db:"code"`
	Name        string  `json:"Name" db:"name"`
	Price       float64 `json:"Price" db:"price"`
	BrandID     string  `json:"BrandID" db:"brand_id"`
	CategoryID  string  `json:"CategoryID" db:"category_id"`
	Description string  `json:"Description" db:"description"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`
}

type Product struct {
	ID          string  `json:"ID" db:"id"`
	Code        string  `json:"Code" db:"code"`
	Name        string  `json:"Name" db:"name" validate:"required"`
	Price       float64 `json:"Price" db:"price" validate:"required"`
	BrandID     string  `json:"BrandID" db:"brand_id" validate:"required"`
	CategoryID  string  `json:"CategoryID" db:"category_id" validate:"required"`
	Description string  `json:"Description" db:"description" validate:"required"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`

	Brand    *IDNameTemplate `json:"Brand"`
	Category *IDNameTemplate `json:"Category"`
}

type FindAllProductParams struct {
	FindAllParams types.FindAllParams
	Code          string
	Name          string
	BrandID       string
	CategoryID    string
}
