// Copyright 2014 Tuan.Pro. All rights reserved.

package apiController

import (
	//"appengine"
	//"appengine/user"
	//"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type ErrorJson struct {
	Message string `json:"error"`
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	message := "{\"msg\": \"Unauthenticated User. Please log in and try again.\"}"
	http.Error(w, message, 401)
	return
}

func forbidden(w http.ResponseWriter, r *http.Request) {
	message := "{\"msg\": \"Bummer. You don't have permission.\"}"
	http.Error(w, message, 403)
	return
}

func notFound(w http.ResponseWriter, r *http.Request) {
	message := "{\"msg\": \"What you are looking for is not here.\"}"
	http.Error(w, message, 404)
	return
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	message := "{\"msg\": \"Oh no. Something majorly went wrong.\"}"
	http.Error(w, message, 500)
	return
}

func ApiGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	view := strings.ToLower(vars["view"])
	key := strings.ToLower(vars["key"])

	log.Println("API GET:", view, "With Key:", key)

	switch view {
	case "users":
		UserGet(w, r, key)
	case "userlookup":
		UserLookupGet(w, r, key)
	case "loginpage":
		LoginPageHtml(w, r)
	case "blogs":
		log.Println("view = blogs")
		if key == "" {
			BlogsIndexGet(w, r, "")
		} else {
			BlogsIndexGet(w, r, key)
		}
	default:
		log.Println("view not found")
		notFound(w, r)
	}
}

func ApiPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	view := strings.ToLower(vars["view"])

	log.Println("API POST:", view)

	switch view {
	case "users":
		UserPost(w, r)
	case "blogs":
		BlogIndexPost(w, r)
	default:
		notFound(w, r)
	}
}

func ApiDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	view := strings.ToLower(vars["view"])
	key := strings.ToLower(vars["key"])

	switch view {
	case "users":
		UserDelete(w, r, key)
	default:
		notFound(w, r)
	}
}
