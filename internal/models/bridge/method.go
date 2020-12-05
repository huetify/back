package bridge

import (
	"context"
	"github.com/ermos/hue"
	"github.com/huetify/back/internal/dbm"
)

func SetBridgesStore(ctx context.Context, db *dbm.Instance) error {
	var bis []Instance

	err := db.GetAll(ctx, &bis, `SELECT * FROM bridge`)
	if err != nil {
		return err
	}

	for _, bi := range bis {
		if !bi.Token.Valid {
			continue
		}

		b := hue.Conn(bi.IPAddr, hue.BridgeOptions{
			Debug: hue.DebugNone,
			Token: bi.Token.String,
		})

		if err := b.Fetch.Bridge(); err != nil {
			return err
		}

		Store[bi.ID] = b
	}

	return nil
}