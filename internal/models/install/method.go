package install

import (
	"context"
	"errors"
	"github.com/huetify/back/internal/dbm"
)

func IsAvailable(ctx context.Context, db *dbm.Instance) error {
	var b bool
	err := db.Get(ctx, &b, `
SELECT
	count(1)
FROM
	user
	`)
	if err != nil {
		return err
	}
	if b {
		return errors.New("not available")
	}
	return nil
}

func SetConfig(ctx context.Context, db *dbm.Instance, name, language string, analytics bool) (err error) {
	_, err = db.Exec(ctx, `UPDATE app SET name = ?, language = ?, analytics = ?`, name, language, analytics)
	return
}
