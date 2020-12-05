package bridge

import (
	"context"
	"errors"
	"fmt"
	"github.com/ermos/hue"
	"github.com/huetify/back/internal/dbm"
	"strings"
)

func DiscoverBridges(ctx context.Context) ([]hue.DiscoverModel, error) {
	dms, err := hue.DiscoverAll()
	if err != nil {
		return dms, err
	}

	var n []hue.DiscoverModel
	l:for _, d := range dms {
		for _, b := range Store {
			if strings.ToLower(b.Config.Bridgeid) == strings.ToLower(d.ID) {
				continue l
			}
		}
		n = append(n, d)
	}

	if len(n) == 0 {
		return dms, errors.New("sorry, we don't found new bridges, check if huetify is on the same network")
	}

	return n, nil
}

func SetBridge(ctx context.Context, db *dbm.Instance, IPAddr string) error {
	for _, b := range Store {
		if b.IPAddr == IPAddr {
			return fmt.Errorf("%s is already logged", IPAddr)
		}
	}

	b := hue.Conn(IPAddr, hue.BridgeOptions{
		Debug: hue.DebugNone,
	})

	if err := b.Fetch.Bridge(); err != nil {
		return err
	}

	res, err := db.Exec(ctx, `INSERT INTO bridge(bridge_uid, ip_address, token) VALUES(?, ?, ?)`, b.Config.Bridgeid, IPAddr, b.Token)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	Store[int(id)] = b

	return nil
}