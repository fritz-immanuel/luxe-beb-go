package models

import (
	"luxe-beb-go/library/types"
)

type UserBulk struct {
	ID       string `json:"ID" db:"id"`
	Name     string `json:"Name" db:"name" validate:"required"`
	Email    string `json:"Email" db:"email"`
	Username string `json:"Username" db:"username"`
	Password string `json:"Password" db:"password" validate:"required"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`
}

type User struct {
	ID       string `json:"ID" db:"id"`
	Name     string `json:"Name" db:"name" validate:"required"`
	Email    string `json:"Email" db:"email"`
	Username string `json:"Username" db:"username"`
	Password string `json:"Password" db:"password" validate:"required"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`
}

type UserLogin struct {
	ID    string `json:"ID" db:"id"`
	Name  string `json:"Name" db:"name" validate:"required"`
	Token string `json:"Token"`
	Email string `json:"Email" db:"email" validate:"required"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`
}

type FindAllUserParams struct {
	FindAllParams types.FindAllParams
	Name          string
	Email         string
	Username      string
	Password      string
}
