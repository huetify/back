package api

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/huetify/back/internal/manager"
	"github.com/huetify/back/internal/models/credentials"
	"github.com/huetify/back/internal/models/user"
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
	userID, err := credentials.CheckCredentials(
		ctx,
		m.DB,
		m.Payload["username"].(string),
		m.Payload["password"].(string),
		)
	if err != nil {
		return response.Error(m.DB, w, http.StatusBadRequest, 3, "username or password is incorrect")
	}

	usr, err := user.GetUser(ctx, m.DB, userID)
	if err != nil {
		return response.Error(m.DB, w, http.StatusInternalServerError, 4, err)
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"issuer": "huetify/api",
		"id": usr.ID,
		"username": usr.Username,
		"roles": usr.Roles,
		"exp": time.Now().Add(time.Minute * time.Duration(120)).Unix(),
		})
	token, err := t.SignedString([]byte(os.Getenv("HUETIFY_JWT_SECRET")))
	if err != nil {
		return response.Error(m.DB, w, http.StatusInternalServerError, 5, err)
	}

	return response.Success(m.DB, w, http.StatusCreated, CredentialsToken{ Token: token, TokenType: "bearer" })
}