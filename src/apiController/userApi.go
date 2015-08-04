// Copyright 2015 Tuan.Pro. All rights reserved.

package apiController

import (
	"fmt"
	"log"
	"io"
    "io/ioutil"
    "appengine/urlfetch"
	"os"
	"time"
	"strings"
	"net/http"
	"crypto/rand"
	"crypto/sha1"
	"encoding/json"	
	"encoding/base64"
	"appengine"
	"appengine/user"
	"appengine/datastore"
	//"appengine/memcache"
	)

const saltSize = 16

func generateSalt(secret []byte) []byte {
	buf := make([]byte, saltSize, saltSize+sha1.Size)
	_, err := io.ReadFull(rand.Reader, buf)
	
	if err != nil {
		log.Println("random read failed: ", err)
		os.Exit(1)
	}
	
	hash := sha1.New()
	
	hash.Write(buf)
	hash.Write(secret)
	return hash.Sum(buf)
}

type User struct {
    UID             string		`json:"id"`
    Email           string		`json:"email"`
    Salt          	string		`json:"-"`
    DisplayName   	string		`json:"displayName"`
    CreateDate		time.Time	`json:"createdDate"`
    ModifiedDate	time.Time	`json:"modifiedDate"`
    ActiveFlag		bool		`json:"active"`
    Role			string		`json:"role"`
    LoginURL		string		`json:"loginURL"`
    LogoutURL		string		`json:"logoutURL"`
}

type Users []User

func loadCurrentUser(c appengine.Context) (User, error) {
    u := user.Current(c)
	
	var currentUser User

	if u == nil {
		currentUser = User{
			Email: "",
			DisplayName: "",
			ActiveFlag: false,
			Role: "Guest",
		}
	} else {
		q := datastore.NewQuery("UserTable").
			Filter("Email =", u.Email).
			Limit(1)
			
		var allUsers []User
	
		if _, err := q.GetAll(c, &allUsers); err != nil {
			return currentUser, err
		}
		
		for k, v := range allUsers {
			if v.Email == u.Email {
				currentUser = allUsers[k]
			}
		}
	}
	
	if currentUser.Role != "Guest" {
		if user.IsAdmin(c) {
			currentUser.Role = "SiteAdmin"
			if currentUser.Email == "" {
				currentUser.Email = u.Email
			}
			currentUser.ActiveFlag = true
		} else if currentUser.Email == u.Email && currentUser.Email != ""{
			currentUser.Role = "KnownUser"
		} else {
			currentUser.Email = u.Email
			currentUser.Role = "UnknownUser"
			currentUser.ActiveFlag = true
		}
	}
	
	currentUser.LoginURL, _ = user.LoginURL(c, "/admin/")
	currentUser.LogoutURL, _ = user.LogoutURL(c, "/admin/")

	return currentUser, nil
}	

func findUser(c appengine.Context, lookupID string) (User, error) {
	k := datastore.NewKey(c, "UserTable", lookupID, 0, nil)

	q := datastore.NewQuery("UserTable").
		Filter("Email =", lookupID).
		Limit(1)
	
	var foundUser User
	var foundUsers []User

	if err := datastore.Get(c, k, &foundUser); err != nil {
		if _, err1 := q.GetAll(c, &foundUsers); err1 != nil {
			return foundUser, err1
		}
		for k, v := range foundUsers {
			if v.Email == lookupID {
				foundUser = foundUsers[k]
			}
		}
		return foundUser, nil
	}
	
	return foundUser, nil
}

func userSave(c appengine.Context, user User) (error) {
	k := datastore.NewKey(c, "UserTable", user.UID, 0, nil)
	
	if _, err := datastore.Put(c, k, &user); err != nil {
		return err
	}

	return nil
}

func loadAllUsers(c appengine.Context) (Users, error) {
	q := datastore.NewQuery("UserTable")

	var users []User

	if _, err := q.GetAll(c, &users); err != nil {
		return users, err
	}

	return users, nil
}

func userGetNameString(c appengine.Context, userID string) (string, string, error) {

	user, err := findUser(c, userID)
	
	if err != nil {
		log.Println("Error finding user name: ", err)
		return "", "", err
	}
	
	return user.DisplayName, user.Email, nil
}

func UserGet(w http.ResponseWriter, r *http.Request, userReqID string) {
    c := appengine.NewContext(r)
	e := json.NewEncoder(w)
	
    userCurrent, err := loadCurrentUser(c)
	log.Println("GET /api/users: entered by", userCurrent.Role, userCurrent.Email)
	
	if err != nil {
		log.Println("GET /api/users: error loading current user", err)
		internalServerError(w, r)
		return
	}
		
 	if userReqID == "" {
 		log.Println("GET /api/users: success lookup user", userCurrent.Email)
	    e.Encode(&userCurrent)
 	} else {
		if userCurrent.Role != "SiteAdmin" {
			log.Println("GET /api/users: unauthorized user access by", userCurrent.Role, userCurrent.Email)
			log.Println("GET /api/users: unauthorized lookup of", userReqID)
			forbidden(w, r)
			return
		}
 		
 		if userReqID == "all"	{
 			users, err := loadAllUsers(c)

			if err != nil {
				log.Println("GET /api/users/all: error loading all users", err)
				internalServerError(w, r)
				return
			}			
			
			log.Println("GET /api/users/all: success lookup all")
			e.Encode(&users)
 		} else {
		    user, err := findUser(c, userReqID)
		
			if err != nil {
				log.Println(err)
				log.Println("GET /api/users/userReqID: error lookup failed on", userReqID)					
					notFound(w, r)
					return
				}

			if user.Email == "" {
				log.Println("GET /api/users/userReqID: user not found", userReqID)
				notFound(w, r)
				return
			}
			
 			log.Println("GET /api/users/userReqID: success lookup user", user.Email)			
			e.Encode(&user)
 		}
 	}
}

func UserDelete(w http.ResponseWriter, r *http.Request, userID string) {
	c := appengine.NewContext(r)
    u := user.Current(c)
	e := json.NewEncoder(w)
	
    if (u == nil) {
    	unauthorized(w, r)
 		return
    }
    
	if user.IsAdmin(c) == false {
		forbidden(w, r)
		return
	}
	
	if (userID == "") {
		notFound(w, r)
		return
	}
	
	userDelete, err := findUser(c, userID)
	
	if (err != nil) {
		notFound(w, r)
		return
	}
	
	if (userDelete.UID != userID) {
		notFound(w, r)
		return
	}
	
	userDelete.ModifiedDate = time.Now()
	userDelete.ActiveFlag = false
	
	if err := userSave(c, userDelete); err != nil {
		log.Println("Error saving user: ", err)
		internalServerError(w, r)
		return
	}
	
	e.Encode(&userDelete)
}

func UserPost(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    u := user.Current(c)
	d := json.NewDecoder(r.Body)
	e := json.NewEncoder(w)

    userCurrent, err := loadCurrentUser(c)
	log.Println("POST /api/users: entered by", userCurrent.Role, userCurrent.Email)

	if err != nil {
		log.Println("POST /api/users: error loading current user", err)
		internalServerError(w, r)
		return
	}

    if (userCurrent.UID == "") {
		log.Println("POST /api/users: error unauthorized access")
    	unauthorized(w, r)
 		return
    }

	var userPost, userEdited, userOld User

	if err := d.Decode(&userPost); err != nil {
		log.Println("POST /api/users: error decoding user post", err)
		internalServerError(w, r)
		return
	}

	var userSalt string
	var userCreate time.Time
	var userName string
	
	if len(userPost.Email) == 0 {
		log.Println("POST /api/users: no email in request body")
		userEdited = userCurrent
	} else {
		userEdited, _ = findUser(c, userPost.Email)
		
		if userEdited.Email == "" {
			log.Println("POST /api/users: user not found", userPost.Email)
			notFound(w, r)
			return
		}
	}
	
	if (userEdited.Email != "") && (userOld.ActiveFlag == false) {
		log.Println("POST /api/users: error user deactivated", userOld.Email)
		forbidden(w, r)
		return
	}
		
	if (userEdited.Salt == "") {
		userSalt = string(base64.URLEncoding.EncodeToString(generateSalt([]byte(u.Email))))
		userCreate = time.Now()
		userName = userPost.DisplayName
	
	} else {
		userSalt = userOld.Salt
		userCreate = userOld.CreateDate
		
		if len(userPost.DisplayName) == 0 {
			userName = userOld.DisplayName
		} else {
			userName = userPost.DisplayName
		}
	}

	
	hashID := sha1.New()
	io.WriteString(hashID, string(u.ID) + userSalt)

	userID := base64.URLEncoding.EncodeToString(hashID.Sum(nil)) 
	
	userID = strings.ToLower(strings.TrimRight(userID, "="))
	
	if (userOld.UID != "" && userOld.UID != userID) {
		unauthorized(w, r)
		return
	}

	user := User {
		UID: userID,
		Email: u.Email,
		Salt: userSalt,
		DisplayName: userName,
		CreateDate: userCreate,
		ModifiedDate: time.Now(),
		ActiveFlag: true}

	if err := userSave(c, user); err != nil {
		log.Println("Error saving user: ", err)
		internalServerError(w, r)
		return
	}
	
	e.Encode(&user)
}	

func LoginPageHtml(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	
	loginPageURL, _ := user.LoginURL(c, "/admin/")
	
	if !strings.HasPrefix(loginPageURL, "http") {
		loginPageURL = "https://goblog-geoct826.c9.io/" + loginPageURL
	} 
	
	loginPageURL = loginPageURL + "&output=embed"
	
	resp, err := client.Get(loginPageURL)

	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Sorry, something went wrong")
		return
	}
	
	defer resp.Body.Close()
	
	contents, err1 := ioutil.ReadAll(resp.Body)
	
	if err1 != nil {
		log.Println(err)
		fmt.Fprintf(w, "Sorry, something went wrong")
		return
	}		
	
	fmt.Fprintf(w, string(contents))
}