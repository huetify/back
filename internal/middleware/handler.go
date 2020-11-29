package middleware

import (
	"context"
	"github.com/huetify/back/internal/manager"
	"log"
	"net/http"
	"reflect"
)

type Handler struct {}

type BeforeHandler struct {}

type AfterHandler struct {}

func (h Handler) Before (ctx context.Context, name string, m *manager.Manager, w http.ResponseWriter, r *http.Request) (status int, err error) {
	build := reflect.ValueOf(BeforeHandler{}).MethodByName(name)
	if build.IsValid() {
		rebuild := make([]reflect.Value, 4)
		rebuild[0] = reflect.ValueOf(ctx)
		rebuild[1] = reflect.ValueOf(m)
		rebuild[2] = reflect.ValueOf(w)
		rebuild[3] = reflect.ValueOf(r)

		result := build.Call(rebuild)

		if len(result) != 2 {
			log.Fatalf("middleware %s need to return (status int, err error)", name)
		}

		if result[1].Interface() != nil {
			return result[0].Interface().(int), result[1].Interface().(error)
		}
	}

	return 0, nil
}

func (h Handler) After (ctx context.Context, name string, m *manager.Manager, w http.ResponseWriter, r *http.Request) (status int, err error) {
	build := reflect.ValueOf(AfterHandler{}).MethodByName(name)
	if build.IsValid() {
		rebuild := make([]reflect.Value, 4)
		rebuild[0] = reflect.ValueOf(ctx)
		rebuild[1] = reflect.ValueOf(m)
		rebuild[2] = reflect.ValueOf(w)
		rebuild[3] = reflect.ValueOf(r)

		result := build.Call(rebuild)

		if len(result) != 2 {
			log.Fatalf("middleware %s need to return (status int, err error)", name)
		}

		if result[1].Interface() != nil {
			return result[0].Interface().(int), result[1].Interface().(error)
		}
	}

	return 0, nil
}
