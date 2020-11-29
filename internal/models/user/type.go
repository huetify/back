package user

import "time"

type User struct {
	ID			int					`json:"user_id"     db:"user_id"`
	Username 	string				`json:"username"    db:"username"`
	Roles 		[]string			`json:"roles"`
	CreatedAt 	time.Time			`json:"created_at"  db:"created_at"`
	UpdatedAt 	time.Time			`json:"updated_at"  db:"updated_at"`
}