// Copyright 2015 Tuan.Pro. All rights reserved.

package app

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "apiController"
    )


func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", 404)
    fmt.Fprintf(w, "Hmm... can't seem to find the page you were looking for.")
	return
}

func index(w http.ResponseWriter, r *http.Request) {
    log.Println("Loading file")
    http.Handle("/",http.FileServer(http.Dir("static")))
}

func init() {
    r := mux.NewRouter()
    
    apiGet := r.PathPrefix("/api").Methods("GET").Subrouter()
    apiGet.HandleFunc("/{view}", apiController.ApiGetHandler)
    apiGet.HandleFunc("/{view}/{key}", apiController.ApiGetHandler)    

    apiPost := r.PathPrefix("/api").Methods("POST").Subrouter()
    apiPost.HandleFunc("/{view}", apiController.ApiPostHandler)
    apiPost.HandleFunc("/{view}/{key}", apiController.ApiPostHandler)    

    apiDelete := r.PathPrefix("/api").Methods("DELETE").Subrouter()
    apiDelete.HandleFunc("/{view}/{key}", apiController.ApiDeleteHandler)    

    r.PathPrefix("/admin/").Handler(http.StripPrefix("/admin/", http.FileServer(http.Dir("./static/admin"))))
    r.PathPrefix("/admin").Handler(http.RedirectHandler("/admin/", 301))

    r.PathPrefix("/blog").Handler(http.StripPrefix("/blog", http.FileServer(http.Dir("./static"))))

    r.PathPrefix("/foundation/").Handler(http.StripPrefix("/foundation/", http.FileServer(http.Dir("./static/foundation"))))

    r.HandleFunc("/{.path:.*}", cloudAdminHandler).Methods("GET")
    r.HandleFunc("/{.path:.*}", cloudAdminPostHandler).Methods("POST")

    http.Handle("/", r)
}

