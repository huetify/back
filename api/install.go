package api

import (
	"context"
	"github.com/huetify/back/internal/manager"
	"github.com/huetify/back/internal/models/install"
	"github.com/huetify/back/internal/models/role"
	"github.com/huetify/back/internal/models/user"
	"github.com/huetify/back/internal/response"
	"github.com/huetify/back/internal/utils"
	"net/http"
)

/*
	@Route("POST", "/install")
	@Payload("username", string)
	@Payload("password", string)
	@Payload("title", string)
	@Payload("language", string)
	@Payload("analytics", bool)
*/
func (Handler) Install(ctx context.Context, m *manager.Manager, w http.ResponseWriter) bool {
	err := utils.StringLength("username", m.Payload["username"].(string), 3, 45)
	if err != nil {
		return response.Error(m.DB, w, http.StatusBadRequest, 3, err)
	}

	err = utils.StringLength("password", m.Payload["password"].(string), 6, 32)
	if err != nil {
		return response.Error(m.DB, w, http.StatusBadRequest, 4, err)
	}

	err = utils.StringLength("title", m.Payload["title"].(string), 3, 32)
	if err != nil {
		return response.Error(m.DB, w, http.StatusBadRequest, 5, err)
	}

	if m.Payload["language"].(string) == "" {
		return response.Error(m.DB, w, http.StatusBadRequest, 6, "you need to choose a language")
	}

	u, err := user.PostUser(ctx, m.DB, m.Payload["username"].(string), m.Payload["password"].(string))
	if err != nil {
		return response.Error(m.DB, w, http.StatusInternalServerError, 7, err)
	}

	err = role.PostUserRole(ctx, m.DB, u.ID, role.Admin)
	if err != nil {
		return response.Error(m.DB, w, http.StatusInternalServerError, 7, err)
	}

	err = install.SetConfig(ctx, m.DB, m.Payload["title"].(string), m.Payload["language"].(string), m.Payload["analytics"].(bool))
	if err != nil {
		return response.Error(m.DB, w, http.StatusInternalServerError, 8, err)
	}

	return response.NoContent(m.DB, w)
}

/*
	@Route("GET", "/install")
*/
func (Handler) CheckInstall(ctx context.Context, m *manager.Manager, w http.ResponseWriter) bool {
	err := install.IsAvailable(ctx, m.DB)
	if err != nil {
		return response.Error(m.DB, w, http.StatusBadRequest, 3, err)
	}

	return response.NoContent(m.DB, w)
}