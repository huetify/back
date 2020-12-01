package credentials

import (
	"context"
	"github.com/huetify/back/internal/dbm"
	"golang.org/x/crypto/bcrypt"
)

type result struct{
	Id 				int		`db:"user_id"`
	HashedPassword 	string	`db:"password"`
}

func CheckCredentials(ctx context.Context, db *dbm.Instance, username, password string) (userID int, err error) {
	var r result
	err = db.Get(ctx, &r, `
SELECT
	user_id,
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

	userID = r.Id

	err = bcrypt.CompareHashAndPassword([]byte(r.HashedPassword), []byte(password))
	return
}
