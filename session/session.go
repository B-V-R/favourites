/*
 * Copyright (c) 2020 BVR (Vighneswar Rao Bojja)
 * This file is subject to the terms and conditions defined in file 'LICENSE'.
 *
 */

package session

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func GetUser(w http.ResponseWriter, r *http.Request) string {
	session, _ := store.Get(r, "favourites-cookie")
	if session.IsNew {
		return r.FormValue("username")
	}
	return fmt.Sprintf("%v", session.Values["user"])
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "favourites-cookie")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false
	}
	return true
}

func CreateSession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "favourites-cookie")
	session.Values["authenticated"] = true
	session.Values["user"] = r.FormValue("username")
	session.Save(r, w)
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "favourites-cookie")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Values["user"] = nil
	session.Save(r, w)
}
