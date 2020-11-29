package api

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/huetify/back/internal/manager"
	"github.com/huetify/back/internal/response"
	"net/http"
	"os"
	"time"
)

type CredentialsToken struct {
	Token 		string `json:"access_token"`
	TokenType 	string `json:"token_type"`
}

/*
	@Route("POST", "/credentials")
	@Payload("username", string)
	@Payload("password", string)
*/
func (Handler) PostCredentials(ctx context.Context, m *manager.Manager, w http.ResponseWriter) bool {
	if m.Payload["username"].(string) != "admin" {
		return response.Error(m.DB, w, http.StatusBadRequest, 3, "username is incorrect")
	}

	if m.Payload["password"].(string) != "admin" {
		return response.Error(m.DB, w, http.StatusBadRequest, 4, "password is incorrect")
	}

	var roles []string
	for _, role := range []string{ "member" } {
		roles = append(roles, role)
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"roles": roles,
		"exp": time.Now().Add(time.Minute * time.Duration(1)).Unix(),
		})
	token, err := t.SignedString([]byte(os.Getenv("HUETIFY_JWT_SECRET")))
	if err != nil {
		return response.Error(m.DB, w, http.StatusInternalServerError, 5, err)
	}

	return response.Success(m.DB, w, http.StatusCreated, CredentialsToken{ Token: token, TokenType: "bearer" })
}

/*
	@Route("GET", "/credentials")
	@Auth(["admin"])
*/
func (Handler) CheckCredentials(ctx context.Context, m *manager.Manager, w http.ResponseWriter) bool {
	return response.Success(m.DB, w, http.StatusCreated, CredentialsToken{ TokenType: "bearer" })
}