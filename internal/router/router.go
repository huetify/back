package router

import (
	"encoding/json"
	"fmt"
	"github.com/ermos/annotation/parser"
	"github.com/huetify/back/internal/manager"
	"github.com/huetify/back/internal/response"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

var withPrefix = true

var withFront = struct{
	Enable 	bool
	Static 	string
	Web 	string
	Upload 	string
}{
	Enable: false,
}

func UsePrefix (b bool) {
	withPrefix = b
}

func EnableFront (staticDir, webDir, uploadDir string) {
	withFront.Static = staticDir
	withFront.Web = webDir
	withFront.Upload = uploadDir
	withFront.Enable = true
}

func Serve (port, routeLocation string, controllerHandler, middlewareHandler interface{}) {
	router := httprouter.New()

	annotations := setRoutes(routeLocation)
	for _, route := range annotations {
		for _, r := range route.Routes {
			var routePath string
			if withPrefix {
				routePath += "/api"
			}
			if route.Version != "" {
				routePath += fmt.Sprintf("/v%s%s", route.Version, r.Route)
			}else{
				routePath += r.Route
			}

			switch strings.ToLower(r.Method) {
			case "get":
				router.GET(routePath, call(route, controllerHandler, middlewareHandler))
			case "post":
				router.POST(routePath, call(route, controllerHandler, middlewareHandler))
			case "put":
				router.PUT(routePath, call(route, controllerHandler, middlewareHandler))
			case "patch":
				router.PATCH(routePath, call(route, controllerHandler, middlewareHandler))
			case "delete":
				router.DELETE(routePath, call(route, controllerHandler, middlewareHandler))
			}
		}
	}


	if withFront.Enable {
		// Web
		router.NotFound = http.HandlerFunc(_front)
		// Static
		router.Handler("GET", "/static/*filepath", _noDirListing(http.StripPrefix("/static/", http.FileServer(http.Dir(withFront.Static)))))
		// Upload
		router.Handler("GET", "/upload/*filepath", _noDirListing(http.StripPrefix("/upload/", http.FileServer(http.Dir(withFront.Upload)))))
	}

	printHeader(port)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		AllowedMethods: []string{"GET","PUT","POST","DELETE","PATCH","OPTIONS"},
		AllowCredentials: false,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(router)


	log.Fatal(http.ListenAndServe(":" + port, handler))
}

func _front (w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, fmt.Sprintf("%s/index.html", withFront.Web))
}

func _noDirListing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

var defaultMiddleware []string

func SetDefaultMiddleware (mw ...string) {
	for _, name := range mw {
		defaultMiddleware = append(defaultMiddleware, name)
	}
}

func printHeader(port string) {
	var content string
	var separator string

	if os.Getenv("api_name") != "" {
		content += fmt.Sprintf("%s's ", os.Getenv("api_name"))
	}

	content += fmt.Sprintf("API currently running on port \033[1m%s\033[0m..", port)

	for i := 0; i < len(content)-6; i++ {
		separator += "-"
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Printf("%s--\n| ", separator)
	fmt.Print(content)
	fmt.Printf(" |\n--%s\n", separator)
}

func setRoutes(location string) (annotations []parser.API){
	file, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &annotations)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func call(route parser.API, handler interface{}, middlewareHandler interface{}) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	build := reflect.ValueOf(handler).MethodByName(route.Controller)
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		m := manager.New(r, route, ps)

		status, err := m.CheckRequest(r)
		if err != nil {
			response.Error(nil, w, status, 1, err)
			return
		}

		var middlewares []string
		middlewares = append(middlewares, defaultMiddleware...)
		for _, md := range route.Middleware.Before {
			middlewares = append(middlewares, md)
		}
		for _, md := range route.Middleware.After {
			middlewares = append(middlewares, md)
		}

		// Before middleware

		beforeMw := reflect.ValueOf(middlewareHandler).MethodByName("Before")

		for _, middleware := range middlewares {
			beforeMwRebuild := make([]reflect.Value, 5)
			beforeMwRebuild[0] = reflect.ValueOf(r.Context())
			beforeMwRebuild[1] = reflect.ValueOf(middleware)
			beforeMwRebuild[2] = reflect.ValueOf(m)
			beforeMwRebuild[3] = reflect.ValueOf(w)
			beforeMwRebuild[4] = reflect.ValueOf(r)

			mwResult := beforeMw.Call(beforeMwRebuild)

			if len(mwResult) != 0 && mwResult[1].Interface() != nil {
				response.Error(nil, w, mwResult[0].Interface().(int), 2, mwResult[1].Interface().(error))
				return
			}
		}

		// Controller

		rebuild := make([]reflect.Value, 3)
		rebuild[0] = reflect.ValueOf(r.Context())
		rebuild[1] = reflect.ValueOf(m)
		rebuild[2] = reflect.ValueOf(w)
		_ = build.Call(rebuild)

		// After middleware

		afterMw := reflect.ValueOf(middlewareHandler).MethodByName("After")

		for _, middleware := range middlewares {
			afterMwRebuild := make([]reflect.Value, 5)
			afterMwRebuild[0] = reflect.ValueOf(r.Context())
			afterMwRebuild[1] = reflect.ValueOf(middleware)
			afterMwRebuild[2] = reflect.ValueOf(m)
			afterMwRebuild[3] = reflect.ValueOf(w)
			afterMwRebuild[4] = reflect.ValueOf(r)

			mwResult := afterMw.Call(afterMwRebuild)

			if len(mwResult) != 0 && mwResult[1].Interface() != nil {
				fmt.Println(mwResult[1].Interface().(error))
			}
		}
	}
}
