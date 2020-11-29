package user

import (
	"context"
	"github.com/huetify/back/internal/dbm"
	"github.com/huetify/back/internal/models/role"
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