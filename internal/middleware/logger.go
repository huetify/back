package middleware

import (
	"context"
	"fmt"
	"github.com/huetify/back/internal/manager"
	"net/http"
	"time"
)

func (h BeforeHandler) HTTP (ctx context.Context, m *manager.Manager, w http.ResponseWriter, r *http.Request) (status int, err error) {
	if r.Method != "OPTIONS"  {
		ipAddress := getIPAddr(r)
		fmt.Printf("[%s] %s \"%s %s\"\n",
			time.Now().Format("02/01/2006T15:04:05+0000"),
			ipAddress,
			r.Method,
			r.RequestURI,
		)
	}
	return 0, nil
}

func getIPAddr(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Real-Ip")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
