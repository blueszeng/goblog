// Copyright 2014 Tuan.Pro. All rights reserved.

package apiController

import (
    "log"
    "fmt"    
	"strings"
	"net/http"
	"appengine"
    "appengine/user"	
	"github.com/gorilla/mux"
	)
	
func HoldOn(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    u := user.Current(c)

    //log.Println("User :", u)
    if u != nil {
    	fmt.Fprintf(w, "User Information \n")
    	fmt.Fprintf(w, " User Id: %v \n", u.ID)	
    	fmt.Fprintf(w, " User Name: %v \n", u)
    	fmt.Fprintf(w, " User Email: %v \n", u.Email)
    	fmt.Fprintf(w, " Administrator: %v \n", u.Admin)
    	fmt.Fprintf(w, " AuthDomain: %v \n", u.AuthDomain)	
    	fmt.Fprintf(w, " Federated Identity: %v \n", u.FederatedIdentity)	
    	fmt.Fprintf(w, " Federated Provider: %v \n", u.FederatedProvider)	

    	
    } else {
		log.Println("Not Logged In")
		unauthorized(w, r)
    }
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
	
	log.Println("API GET:", view)
	
	switch view {
		case "users": UserGet(w, r, key)
		case "loginpage": LoginPageHtml(w, r)
		case "blogs": BlogsIndexGet(w, r, key)
		case "test": HoldOn(w, r)
		default: notFound(w, r)
	}
}

func ApiPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	view := strings.ToLower(vars["view"])

	log.Println("API POST:", view)
	
	switch view {
		case "users": UserPost(w, r)
		case "blogs": BlogIndexPost(w, r)
		default: notFound(w, r)
	}
}


func ApiDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	view := strings.ToLower(vars["view"])
	key := strings.ToLower(vars["key"])
	
	switch view {
		case "users": UserDelete(w, r, key)
		default: notFound(w, r)
	}
}
