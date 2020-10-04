/*
 * Copyright (c) 2020 BVR (Vighneswar Rao Bojja)
 * This file is subject to the terms and conditions defined in file 'LICENSE'.
 *
 */

package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"vighnesh.org/favourate/database"
	"vighnesh.org/favourate/database/schema"
	"vighnesh.org/favourate/security"
	"vighnesh.org/favourate/session"
)

// Router interface for handler
type Router interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
	Favourite(w http.ResponseWriter, r *http.Request)
	Favorites(w http.ResponseWriter, r *http.Request)
}

type router struct {
	log      *log.Logger
	database database.Database
}

// SignUp register user to application
func (router *router) SignUp(w http.ResponseWriter, r *http.Request) {
	router.log.Println("request received at endpoint: SignUp")

	if session.IsAuthenticated(w, r) {
		router.Favorites(w, r)
		return
	}

	user := &schema.User{}

	r.ParseForm()

	user.Email = r.FormValue("email")
	if username := r.FormValue("username"); username != "" {
		user.User = username
	} else {
		SignInPage(w, r)
		return
	}

	if pwd := r.FormValue("password"); pwd != "" {
		pwdHash, _ := security.HashPassword(pwd)
		user.Password = pwdHash
	} else {
		SignInPage(w, r)
		return
	}

	router.database.Save(user)

	router.log.Println("User registered: ", user.User)

	router.log.Println("Forwarding request to sign in page: ", user.User)
	SignInPage(w, r)
	return
}

// SignIn authenticates user credentials
func (router *router) SignIn(w http.ResponseWriter, r *http.Request) {
	router.log.Println("request received at endpoint: SignIn")
	user := &schema.User{}

	r.ParseForm()

	if username := r.FormValue("username"); username != "" {
		user.User = username
	} else {
		router.log.Println("username is empty redirecting signin page")
		SignInPage(w, r)
		return
	}

	if pwd := r.FormValue("password"); pwd != "" {
		user.Password = pwd
	} else {
		router.log.Println("password is empty redirecting signin page")
		SignInPage(w, r)
		return
	}

	dbUser := router.database.User(user.User)

	if security.CheckPasswordHash(user.Password, dbUser.Password) {
		router.log.Println("Sign In success for user ", user.User)
		session.CreateSession(w, r)
		router.Favorites(w, r)
		return
	} else {
		router.log.Println("Sign In failed for user ", user.User, "redirecting signin page")
		SignInPage(w, r)
		return
	}

}

// SignIn authenticates user credentials
func (router *router) SignOut(w http.ResponseWriter, r *http.Request) {
	router.log.Println("request received at endpoint: SignOut")
	if session.IsAuthenticated(w, r) {
		router.database.Delete(&schema.Favourite{User: session.GetUser(w, r)})
		session.SignOut(w, r)
		router.log.Println("sign out completed redirecting to home page")
	} else {
		router.log.Println("Not signed in to sign out, redirecting to home page")
	}

	HomePage(w, r)
	return
}

// Favourite stores user favourites
func (router *router) Favourite(w http.ResponseWriter, r *http.Request) {
	router.log.Println("request received at endpoint: Favourite")

	favourite := schema.Favourite{}

	r.ParseForm()

	favourite.User = session.GetUser(w, r)
	favourite.Favourite = r.FormValue("favourite")

	favourites := router.database.Favourites(favourite.User)

	router.log.Println(favourites)

	for _, dbFavourite := range favourites {
		if favourite.Favourite == dbFavourite.Favourite {
			router.log.Println("Already in your favourites list")
			router.Favorites(w, r)
			return
		}
	}

	router.database.Favourite(&favourite)
	router.log.Println("Added to your favourites list")
	//ImagesPage(w,r)
	router.Favorites(w, r)
	return
}

// Favorites list user favourites
func (router *router) Favorites(w http.ResponseWriter, r *http.Request) {
	router.log.Println("request received at endpoint: Favorites")

	user := schema.User{}

	body := r.Body
	data, err := ioutil.ReadAll(body)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	json.Unmarshal(data, user)

	favourites := router.database.Favourites(user.User)
	FavouritesPage(w, favourites)
}

func New(database database.Database) Router {
	return &router{
		log:      log.New(os.Stdout, "favourite-app-handler:\t", log.Ldate|log.Ltime|log.Lshortfile),
		database: database,
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	file, _ := ioutil.ReadFile("./html/home.html")
	fmt.Fprint(w, string(file))
	return
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	file, _ := ioutil.ReadFile("./html/signup.html")
	fmt.Fprint(w, string(file))
}

func SignInPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	file, _ := ioutil.ReadFile("./html/signin.html")
	fmt.Fprint(w, string(file))
}
func FavouritesPage(w http.ResponseWriter, favourites []schema.Favourite) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	file, _ := ioutil.ReadFile("./html/images.html")
	fileStr := string(file)

	for _, favourite := range favourites {
		if strings.Contains(fileStr, "id=\""+favourite.Favourite+"\"") {
			fileStr = strings.ReplaceAll(fileStr, "id=\""+favourite.Favourite+"\"", "id=\""+favourite.Favourite+"\""+"style=\"background-color: #008CBA\"")
		}
	}
	fmt.Fprint(w, fileStr)
}
