package credentials

import (
	"context"
	"github.com/huetify/back/internal/dbm"
	"golang.org/x/crypto/bcrypt"
)

func CheckCredentials(ctx context.Context, db *dbm.Instance, username, password string) (userID int, err error) {
	var result struct{
		id 				int		`db:"user_id"`
		hashedPassword 	string	`db:"password"`
	}
	err = db.Get(ctx, &result, `
SELECT
	password
FROM
	user
WHERE
	username = ?
	`,
		username,
	)
	if err != nil {
		return
	}

	userID = result.id

	err = bcrypt.CompareHashAndPassword([]byte(result.hashedPassword), []byte(password))
	return
}
