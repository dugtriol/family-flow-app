package entity

import `database/sql`

type User struct {
	Id        string          `pgdb:"id" json:"id"`
	Name      string          `pgdb:"name" json:"name"`
	Email     string          `pgdb:"email" json:"email"`
	Password  string          `pgdb:"password" json:"password"`
	Role      string          `pgdb:"role" json:"role"`
	FamilyId  sql.NullString  `pgdb:"family_id" json:"family_id" swaggerignore:"true"`
	Latitude  sql.NullFloat64 `pgdb:"latitude" json:"latitude" swaggerignore:"true"`
	Longitude sql.NullFloat64 `pgdb:"longitude" json:"longitude" swaggerignore:"true"`
}
