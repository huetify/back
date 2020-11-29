package middleware

import (
	"context"
	"github.com/huetify/back/internal/dbm"
	"github.com/huetify/back/internal/manager"
	"net/http"
	"os"
)

func (h BeforeHandler) DBStart (ctx context.Context, m *manager.Manager, w http.ResponseWriter, r *http.Request) (status int, err error) {
	i, err := dbm.Conn(
		ctx,
		os.Getenv("HUETIFY_DB_NAME"),
		os.Getenv("HUETIFY_DB_DRIVER"),
		os.Getenv("HUETIFY_DB_USER"),
		os.Getenv("HUETIFY_DB_PASSWORD"),
		os.Getenv("HUETIFY_DB_HOST"),
		os.Getenv("HUETIFY_DB_PORT"),
	)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	m.SetDB(i)

	return 0, nil
}

func (h AfterHandler) DBClose (ctx context.Context, m *manager.Manager, w http.ResponseWriter, r *http.Request) (status int, err error) {
	_ = m.DB.Commit()

	err = m.DB.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}
