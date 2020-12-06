package bridge

import (
	"time"
)

type Instance struct {
	ID			int			`db:"bridge_id"   json:"bridge_id"`
	Name 		string		`db:"name"        json:"name"`
	UID			string		`db:"bridge_uid"  json:"bridge_uid"`
	IPAddr		string		`db:"ip_address"  json:"ip_address"`
	Token		string		`db:"token"       json:"token"`
	CreatedAt	time.Time	`db:"created_at"  json:"created_at"`
}