package user

import (
	"context"
	"github.com/huetify/back/internal/dbm"
	"github.com/huetify/back/internal/models/role"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(ctx context.Context, db *dbm.Instance, userID int) (u User, err error) {
	err = db.Get(ctx, &u, `
SELECT
	user_id,
	username,
	created_at,
	updated_at
FROM
	user
WHERE
	user_id = ?
	`,
		userID,
	)
	if err != nil {
		return
	}
	u.Roles, err = role.GetUserRoles(ctx, db, userID)
	return
}

func PostUser(ctx context.Context, db *dbm.Instance, username, password string) (u User, err error) {
	crypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	res, err := db.Exec(ctx, `INSERT user(username, password) VALUES(?, ?)`, username, string(crypted))
	if err != nil {
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		return
	}

	u = User{
		ID: int(id),
		Username: username,
	}

	return
}