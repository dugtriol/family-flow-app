package entity

import "database/sql"

type User struct {
	Id       string         `pgdb:"id"`
	Name     string         `pgdb:"name"`
	Email    string         `pgdb:"email"`
	Password string         `pgdb:"password"`
	Role     string         `pgdb:"role"`
	FamilyId sql.NullString `pgdb:"family_id"`
}
