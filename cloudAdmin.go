// Copyright 2015 Tuan.Pro. All rights reserved.

package app

import (
    "io"
    "log"
    "html"
    "strings"
    "net/http"
    "appengine"
    "appengine/urlfetch"
    )

func copyHeader(dst, src http.Header) {
    for k, w := range src {
        for _, v := range w {
            dst.Add(k, v)
        }
    }
}
    
func copyResponse(r *http.Response, w http.ResponseWriter) {
    copyHeader(w.Header(), r.Header)
    w.WriteHeader(r.StatusCode)
    io.Copy(w, r.Body)
}
    
func cloudAdminHandler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    
    path := html.EscapeString(r.URL.Path)
    query := html.EscapeString(r.URL.RawQuery)
    
    views := strings.Split(path, "/")

    switch views[1] {
        case "cloudadmin" : path = ""
        case "instances" :
        case "datastore" :
        case "datastore-indexes" :
        case "datastore-stats" :
        case "console" :
        case "memcache" :
        case "blobstore" :
        case "taskqueue" :
        case "cron" : 
        case "xmpp" : 
        case "mail" : 
        case "search" : 
        default : path = "notfound"
    }

    log.Println("CloudAdmin GET View:", views[1])
    log.Println("CloudAdmin GET Path:", path)
    log.Println("CloudAdmin GET Query:", query)

    resp, err := client.Get("http://0.0.0.0:8000" + path + "?" + query)
    if err != nil {
        notFound(w, r)
        log.Println("CloudAdmin GET Error:", err.Error())
        return
    }
    defer resp.Body.Close()
    copyResponse(resp, w)
}

func cloudAdminPostHandler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    
    path := html.EscapeString(r.URL.Path)
    query := html.EscapeString(r.URL.RawQuery)
    views := strings.Split(path, "/")        

    switch views[1] {
        case "cloudadmin" : path = ""
        case "instances" :
        case "datastore" :
        case "datastore-indexes" :
        case "datastore-stats" :
        case "console" :
        case "memcache" :
        case "blobstore" :
        case "taskqueue" :
        case "cron" : 
        case "xmpp" : 
        case "mail" : 
        case "search" : 
        default : path = "notfound"
    }

    log.Println("CloudAdmin POST View:", views[1])
    log.Println("CloudAdmin POST Path:", path)
    log.Println("CloudAdmin POST Query:", query)
    
    resp, err := client.Post("http://0.0.0.0:8000" + path + "?" + query,
        "application/x-www-form-urlencoded",
        io.Reader(r.Body))
    if err != nil {
        notFound(w, r)
        log.Println("CloudAdmin POST Error:", err.Error())
        return
    }
    defer resp.Body.Close()
    copyResponse(resp, w)
}
