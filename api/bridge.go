package api

import (
	"context"
	"github.com/huetify/back/internal/manager"
	"github.com/huetify/back/internal/models/bridge"
	"github.com/huetify/back/internal/response"
	"net/http"
)

/*
	@Route("GET", "/bridges/instance")
	@auth(["admin"])
*/
func (Handler) GetBridgesInstance(ctx context.Context, m *manager.Manager, w http.ResponseWriter) bool {
	bi, err := bridge.GetBridgesInstance(ctx, m.DB)
	if err != nil {
		return response.Error(m.DB, w, http.StatusInternalServerError, 3, err)
	}

	return response.Success(m.DB, w, http.StatusOK, bi)
}

/*
	@Route("GET", "/bridges/discover")
	@auth(["admin"])
*/
func (Handler) BridgesDiscover(ctx context.Context, m *manager.Manager, w http.ResponseWriter) bool {
	bridges, err := bridge.DiscoverBridges(ctx)
	if err != nil {
		return response.Error(m.DB, w, http.StatusBadRequest, 3, err)
	}

	return response.Success(m.DB, w, http.StatusOK, bridges)
}

/*
	@Route("POST", "/bridges")
	@auth(["admin"])
	@Payload("ip_addr", string)
*/
func (Handler) PostBridge(ctx context.Context, m *manager.Manager, w http.ResponseWriter) bool {
	b, err := bridge.SetBridge(ctx, m.DB, m.Payload["ip_addr"].(string))
	if err != nil {
		return response.Error(m.DB, w, http.StatusInternalServerError, 3, err)
	}

	return response.Success(m.DB, w, http.StatusCreated, b)
}