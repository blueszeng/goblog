// Copyright 2015 Tuan.Pro. All rights reserved.

package apiController

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"appengine/user"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	//"log"
	"net/http"
	"os"
	"strings"
	"time"
	//"appengine/memcache"
)

const saltSize = 16

func generateSalt(secret []byte) []byte {
	buf := make([]byte, saltSize, saltSize+sha1.Size)
	_, err := io.ReadFull(rand.Reader, buf)

	if err != nil {
		//log.Println("random read failed: ", err)
		os.Exit(1)
	}

	hash := sha1.New()

	hash.Write(buf)
	hash.Write(secret)
	return hash.Sum(buf)
}

type User struct {
	UID          string    `json:"id"`
	Email        string    `json:"email"`
	Salt         string    `json:"-"`
	DisplayName  string    `json:"displayName"`
	CreateDate   time.Time `json:"createdDate"`
	ModifiedDate time.Time `json:"modifiedDate"`
	ActiveFlag   bool      `json:"active"`
	Role         string    `json:"role"`
	LoginURL     string    `json:"loginURL"`
	LogoutURL    string    `json:"logoutURL"`
}

type Users []User

func loadCurrentUser(c appengine.Context) (User, error) {
	u := user.Current(c)

	var currentUser User

	//log.Println("looking for user", u)

	if u == nil {
		currentUser = User{
			Email:       "no email",
			DisplayName: "",
			ActiveFlag:  false,
			Role:        "Guest",
		}
	} else {
		//log.Println(u.ID)
		userEmail := strings.ToLower(u.Email)

		q := datastore.NewQuery("UserTable").
			Filter("Email =", userEmail).
			Limit(1)

		var allUsers []User

		if _, err := q.GetAll(c, &allUsers); err != nil {
			return currentUser, err
		}

		for k, v := range allUsers {
			if v.Email == userEmail {
				currentUser = allUsers[k]
			}
		}
	}

	if currentUser.Email == "" && currentUser.Role != "" {
		currentUser.Email = u.Email
	}

	if user.IsAdmin(c) {
		currentUser.Role = "SiteAdmin"
		currentUser.Email = u.Email
		currentUser.ActiveFlag = true
	}

	if currentUser.Role == "" {
		if user.IsAdmin(c) {
			currentUser.Role = "SiteAdmin"
			if currentUser.Email == "" {
				currentUser.Email = u.Email
			}
			currentUser.ActiveFlag = true
		} else if currentUser.Email == u.Email && currentUser.Email != "" {
			currentUser.Role = "KnownUser"
		} else if currentUser.UID == "" {
			currentUser.Email = u.Email
			currentUser.ActiveFlag = true
			currentUser.Role = "Unregistered"
		} else {
			currentUser.Email = u.Email
			currentUser.Role = "UnknownUser"
			currentUser.ActiveFlag = true
		}
	}

	currentUser.LoginURL, _ = user.LoginURL(c, "/admin/")
	currentUser.LogoutURL, _ = user.LogoutURL(c, "/admin/")

	//c.Infof("%v", currentUser.Email)
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

func userSave(c appengine.Context, user User) error {
	k := datastore.NewKey(c, "UserTable", user.UID, 0, nil)

	if _, err := datastore.Put(c, k, &user); err != nil {
		return err
	}

	return nil
}

func userDelete(c appengine.Context, user User) error {
	k := datastore.NewKey(c, "UserTable", user.UID, 0, nil)

	if err := datastore.Delete(c, k); err != nil {
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
		//log.Println("Error finding user name: ", err)
		return "", "", err
	}

	return user.DisplayName, user.Email, nil
}

func UserGet(w http.ResponseWriter, r *http.Request, userReqID string) {
	c := appengine.NewContext(r)
	e := json.NewEncoder(w)

	userCurrent, err := loadCurrentUser(c)
	c.Infof("GET /api/users/%v: Entered by user: %v (%v)", userReqID, userCurrent.Email, userCurrent.Role)

	if err != nil {
		c.Errorf("GET /api/users/%v: Error loading current user: %v", userReqID, err)
		internalServerError(w, r)
		return
	}

	if userReqID == "" {
		c.Infof("GET /api/users/%v: Exited successfully", userReqID)
		//log.Println(userCurrent)
		e.
			Encode(&userCurrent)
	} else {
		if userCurrent.Role != "SiteAdmin" {
			c.Warningf("GET /api/users/%v: Unauthorized access by user: %v", userReqID, userCurrent.Email)
			forbidden(w, r)
			return
		}

		if userReqID == "all" {
			users, err := loadAllUsers(c)

			if err != nil {
				c.Errorf("GET /api/users/all: Error loading all users: %v", err)
				internalServerError(w, r)
				return
			}
			c.Infof("GET /api/users/all: Exited successfully")
			e.Encode(&users)
		} else {
			user, err := findUser(c, userReqID)

			if err != nil {
				c.Errorf("GET /api/users/%v: Error loading user: %v", userReqID, err)
				notFound(w, r)
				return
			}

			if user.Email == "" {
				c.Warningf("GET /api/users/%v: User not found", userReqID)
				notFound(w, r)
				return
			}

			c.Infof("GET /api/users/%v: Exited successfully", userReqID)
			//log.Println(user)
			e.Encode(&user)
		}
	}
}

func UserDelete(w http.ResponseWriter, r *http.Request, userID string) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	e := json.NewEncoder(w)

	if u == nil {
		unauthorized(w, r)
		return
	}

	if user.IsAdmin(c) == false {
		forbidden(w, r)
		return
	}

	if userID == "" {
		notFound(w, r)
		return
	}

	userDelete, err := findUser(c, userID)

	if err != nil {
		notFound(w, r)
		return
	}

	if userDelete.UID != userID {
		notFound(w, r)
		return
	}

	userDelete.ModifiedDate = time.Now()
	userDelete.ActiveFlag = false

	if err := userSave(c, userDelete); err != nil {
		//log.Println("Error saving user: ", err)
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
	c.Infof("POST /api/users: Entered by user: %v (%v)", userCurrent.Email, userCurrent.Role)

	if err != nil {
		c.Errorf("POST /api/users: Error loading current user: %v", err)
		internalServerError(w, r)
		return
	}

	if userCurrent.Role == "" {
		c.Warningf("POST /api/users: Unauthorized access by user: %v", userCurrent.Email)
		unauthorized(w, r)
		return
	}

	var userPost, userEdited User

	if err := d.Decode(&userPost); err != nil {
		c.Errorf("POST /api/users: Error decoding user post: %v", err)
		internalServerError(w, r)
		return
	}

	var userSalt string
	var userCreate time.Time
	var userName string
	var userID string
	var userRole string

	hasUserInfo := true
	cleanupOldUser := false

	if len(userPost.Email) == 0 || userPost.Email == userCurrent.Email {
		userEdited = userCurrent
		c.Infof("POST /api/users: Editing user: %v", userEdited.Email)
	} else {
		userEdited, _ = findUser(c, userPost.Email)

		if userEdited.Email == "" {
			hasUserInfo = false
			userEdited.Email = userPost.Email
		}
	}

	currentHashID := sha1.New()
	io.WriteString(currentHashID, string(u.ID)+userEdited.Salt)

	currentUserID := base64.URLEncoding.EncodeToString(currentHashID.Sum(nil))
	currentUserID = strings.ToLower(strings.TrimRight(currentUserID, "="))

	if currentUserID != userEdited.UID {
		if userCurrent.Role == "New" && userCurrent.Email == userEdited.Email {
			cleanupOldUser = true
			c.Infof("POST /api/users: Setting New user: %v", userEdited.Email)
		} else if userCurrent.Role == "SiteAdmin" {
			c.Infof("POST /api/users: Admin Editing user: %v", userEdited.Email)
		} else {
			c.Warningf("POST /api/users: Unauthorized access by user: %v", userCurrent.Email)
			c.Warningf("POST /api/users: Unauthorized editing of: %v", userPost.Email)
			forbidden(w, r)
			return
		}
	}

	if userEdited.ActiveFlag == false && hasUserInfo {
		c.Warningf("POST /api/users: User Deactivited: %v", userEdited.Email)
		forbidden(w, r)
		return
	}

	if userEdited.Salt == "" || userEdited.Role == "New" {
		userSalt = string(base64.URLEncoding.EncodeToString(generateSalt([]byte(u.Email))))
		userCreate = time.Now()
		userName = userPost.DisplayName

		hashID := sha1.New()
		io.WriteString(hashID, string(u.ID)+userSalt)

		userID = base64.URLEncoding.EncodeToString(hashID.Sum(nil))
		userID = strings.ToLower(strings.TrimRight(userID, "="))
	} else {
		userSalt = userEdited.Salt
		userCreate = userEdited.CreateDate
		userID = userEdited.UID

		if len(userPost.DisplayName) == 0 {
			userName = userEdited.DisplayName
		} else {
			userName = userPost.DisplayName
		}
	}

	if !hasUserInfo {
		userRole = "New"
	} else {
		userRole = "KnownUser"
	}

	userUpdate := User{
		UID:          userID,
		Email:        strings.ToLower(userEdited.Email),
		Salt:         userSalt,
		DisplayName:  userName,
		CreateDate:   userCreate,
		ModifiedDate: time.Now(),
		Role:         userRole,
		ActiveFlag:   true}

	if err := userSave(c, userUpdate); err != nil {
		c.Errorf("POST /api/users: Error saving user: %v", err)
		internalServerError(w, r)
		return
	}

	if cleanupOldUser {
		if err := userDelete(c, userEdited); err != nil {
			c.Errorf("POST /api/users: Error cleaning up old user: %v", err)
			internalServerError(w, r)
			return
		}
	}

	userUpdate.LoginURL, _ = user.LoginURL(c, "/admin/")
	userUpdate.LogoutURL, _ = user.LogoutURL(c, "/admin/")

	c.Infof("POST /api/users: Exited successfully")
	e.Encode(&userUpdate)
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
		//log.Println(err)
		fmt.Fprintf(w, "Sorry, something went wrong")
		return
	}

	defer resp.Body.Close()

	contents, err1 := ioutil.ReadAll(resp.Body)

	if err1 != nil {
		//log.Println(err)
		fmt.Fprintf(w, "Sorry, something went wrong")
		return
	}

	fmt.Fprintf(w, string(contents))
}

func UserLookupGet(w http.ResponseWriter, r *http.Request, userEmail string) {
	c := appengine.NewContext(r)
	e := json.NewEncoder(w)

	notFoundError := ErrorJson{
		Message: "No Blogs Found",
	}

	userCurrent, err := loadCurrentUser(c)
	c.Infof("GET /api/userlookup/%v: Entered by user: %v (%v)", userEmail, userCurrent.Email, userCurrent.Role)

	if err != nil {
		c.Errorf("GET /api/userlookup/%v: Error loading current user: %v", userEmail, err)
		internalServerError(w, r)
		return
	}

	if userCurrent.Role != "SiteAdmin" {
		c.Warningf("GET /api/userlookup/%v: Unauthorized access by user: %v", userEmail, userCurrent.Email)
		forbidden(w, r)
		return
	}

	user, err := findUser(c, userEmail)

	if err != nil {
		c.Errorf("GET /api/userlookup/%v: Error finding user name: ", userEmail, err)
		e.Encode(&notFoundError)
	}

	userInfo := Author{
		Name:  user.DisplayName,
		Email: user.Email,
	}

	c.Infof("GET /api/userlookup/%v: Exited successfully", userEmail)
	e.Encode(&userInfo)
}
