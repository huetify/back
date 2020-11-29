package role

import (
	"context"
	"github.com/huetify/back/internal/dbm"
)

func GetUserRoles(ctx context.Context, db *dbm.Instance, userID int) (roles []string, err error) {
	err = db.GetAll(ctx, &roles, `
SELECT
	name
FROM
	user_role
JOIN
	role ON role.role_id = user_role.role_id
WHERE
	user_role.user_id = ?
	`,
	userID,
	)
	return
}