package main

import (
	"github.com/ermos/annotation"
	"github.com/ermos/annotation/parser"
	"github.com/ermos/dotenv"
	"github.com/huetify/back/api"
	"github.com/huetify/back/internal/middleware"
	"github.com/huetify/back/internal/router"
	"log"
	"os"
)

func main() {
	// Get environment variable
	if err := dotenv.Parse(".env"); err != nil {
		log.Fatal(err)
	}
	// Check require env
	err := dotenv.Require(
		"HUETIFY_PORT",
		"HUETIFY_JWT_SECRET",
		"HUETIFY_DEBUG",
		"HUETIFY_DB_DRIVER",
		"HUETIFY_DB_HOST",
		"HUETIFY_DB_PORT",
		"HUETIFY_DB_USER",
		"HUETIFY_DB_PASSWORD",
		"HUETIFY_DB_NAME",
		)
	if err != nil {
		log.Fatal(err)
	}
	// Build API's Annotation
	if len(os.Args) > 1 && os.Args[1] == "build" {
		var routes []parser.API
		err = annotation.Fetch("api", &routes, parser.ToAPI)
		if err != nil {
			log.Fatal(err)
		}
		err = annotation.Save(routes, "router.json")
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
	// Router configuration
	router.SetDefaultMiddleware("HTTP", "DBStart", "DBStart")
	router.Serve(os.Getenv("HUETIFY_PORT"), "router.json", api.Handler{}, middleware.Handler{})
}
