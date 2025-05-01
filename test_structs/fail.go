package test_structs

import "time"

// tagit
type UserFail struct {
	Id             int       `json:"id" db:"id"`
	OrganizationId int       `json:"organization_id" db:"organization_id"`
	Email          string    `json:"email" db:"email"`
	FirstName      string    `json:"first_name" db:"first_name"`
	LastName       string    `json:"last_name" db:"last_name"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// tagit
type OrganizationFail struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Domain      string `json:"domain" db:"domain"`
	Description string `json:"description" db:"description"`
}

// tagit
type CompanyFail struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Domain      string `json:"domain" db:"domain"`
	Description string `json:"description" db:"description"`
}
