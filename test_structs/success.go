package test_structs

import "time"

// tagit
type User struct {
	Id             int       `json:"id" db:"id"`
	OrganizationId int       `json:"organization_id" db:"organization_id"`
	Email          string    `json:"email" db:"email"`
	FirstName      string    `json:"first_name" db:"first_name"`
	LastName       string    `json:"last_name" db:"last_name"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type Organization struct {
	Id          int
	Name        string
	Domain      string
	Description string
}

// tagit
type Company struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Domain      string `json:"domain" db:"domain"`
	Description string `json:"description" db:"description"`
}
