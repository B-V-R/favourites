/*
 * Copyright (c) 2020 BVR (Vighneswar Rao Bojja)
 * This file is subject to the terms and conditions defined in file 'LICENSE'.
 *
 */

package app

import (
	"log"
	"net/http"
	"os"
	"vighnesh.org/favourate/database"
	"vighnesh.org/favourate/handler"
	"vighnesh.org/favourate/server"
	"vighnesh.org/favourate/session"
)

// App interface for application
type App interface {
	Start()
}

type app struct {
	service  server.Server
	router   handler.Router
	database database.Database
	log      *log.Logger
}

// App creates a new App
func New() App {
	app := app{}

	server := server.New()
	database := database.New()
	log := log.New(os.Stdout, "favourite-app:\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err := database.CreateTables(); err != nil {
		log.Fatal(err)
	}

	app.database = database
	router := handler.New(app.database)

	app.service = server
	app.router = router

	app.log = log

	return &app
}

func (app *app) registerRoutes() {
	app.service.Register("/", func(writer http.ResponseWriter, request *http.Request) {
		if session.IsAuthenticated(writer, request) {
			app.router.Favorites(writer, request)
			return
		}
		handler.HomePage(writer, request)
		return
	})
	app.service.Register("/page/images", func(writer http.ResponseWriter, request *http.Request) {
		if session.IsAuthenticated(writer, request) {
			app.router.Favorites(writer, request)
			return
		}
		handler.HomePage(writer, request)
		return
	})
	app.service.Register("/page/signup", func(writer http.ResponseWriter, request *http.Request) {
		if session.IsAuthenticated(writer, request) {
			app.router.Favorites(writer, request)
			return
		}
		handler.SignUpPage(writer, request)
		return
	})
	app.service.Register("/page/signin", func(writer http.ResponseWriter, request *http.Request) {
		if session.IsAuthenticated(writer, request) {
			app.router.Favorites(writer, request)
			return
		}
		handler.SignInPage(writer, request)
		return
	})
	app.service.Register("/signup", app.router.SignUp)
	app.service.Register("/signin", app.router.SignIn)
	app.service.Register("/signout", app.router.SignOut)
	app.service.Register("/favourite", app.router.Favourite)
	app.service.Register("/favourites", app.router.Favorites)

	app.log.Println("Routes registered")
}

// Start register routes and starts service
func (app *app) Start() {
	app.registerRoutes()
	app.service.Start(":80")
}
