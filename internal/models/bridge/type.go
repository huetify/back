package bridge

import (
	"database/sql"
	"time"
)

type Instance struct {
	ID			int				`db:"bridge_id"   json:"bridge_id"`
	UID			sql.NullString	`db:"bridge_uid"  json:"bridge_uid"`
	IPAddr		string			`db:"ip_address"  json:"ip_address"`
	Token		sql.NullString	`db:"token"       json:"token"`
	CreatedAt	time.Time		`db:"created_at"  json:"created_at"`
}