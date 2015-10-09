// Copyright 2015 Tuan.Pro. All rights reserved.

package apiController

import (
	"time"
	//"sort"
	//"log"
	"net/http"
	//"strconv"
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
	"html/template"
	//"appengine/memcache"
	"code.google.com/p/go-uuid/uuid"
)

type Entry struct {
	ID           string        `json:"id"`
	Text         template.HTML `datastore:",noindex", json:"text"`
	AuthorID     string        `json:"-"`
	PostAuthor   Author        `json:"postAuthor"`
	CreatedDate  time.Time     `json:"createdDate"`
	SavedDate    time.Time     `json:"savedDate`
	FinishedFlag bool          `json:"finished"`
}

type Entries []Entry

func entrySave(c appengine.Context, entry Entry, postID string) error {
	postKey := datastore.NewKey(c, "PostIndex", postID, 0, nil)
	k := datastore.NewKey(c, "Entry", entry.ID, 0, postKey)

	if _, err := datastore.Put(c, k, &entry); err != nil {
		return err
	}

	return nil
}

func entriesLoadAll(c appengine.Context, postID string) (Entries, error) {
	postKey := datastore.NewKey(c, "PostIndex", postID, 0, nil)
	q := datastore.NewQuery("Entry").
		Ancestor(postKey).
		Order("SavedDate")

	var entries Entries

	if _, err := q.GetAll(c, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func entryLoad(c appengine.Context, entryID string, postID string) (Entry, error) {
	postKey := datastore.NewKey(c, "PostIndex", postID, 0, nil)
	k := datastore.NewKey(c, "Entry", entryID, 0, postKey)

	var entry Entry

	if err := datastore.Get(c, k, &entry); err != nil {
		return entry, err
	}

	return entry, nil
}

func entryLastPost(c appengine.Context) (Entry, error) {
	q := datastore.NewQuery("PostIndex").
		Filter("FinishedFlag =", true).
		Order("SavedDate").Limit(1)

	var entries Entries
	var entry Entry

	_, err := q.GetAll(c, &entries)

	if err != nil {
		return entry, err
	}

	if entries == nil {
		return entry, nil
	} else {
		return entry, nil
	}
}

func EntryPost(w http.ResponseWriter, r *http.Request, postID string) {
	c := appengine.NewContext(r)
	d := json.NewDecoder(r.Body)
	e := json.NewEncoder(w)

	author, err := loadCurrentUser(c)
	c.Infof("POST /api/entries/%v: Entered by user: %v (%v)", postID, author.Email, author.Role)

	if user.IsAdmin(c) == false {
		c.Warningf("POST /api/entries/%v: Unauthorized access by user: %v", postID, author.Email)
		forbidden(w, r)
		return
	}

	if err != nil {
		c.Errorf("POST /api/entries/%v: Error loading user: %v", postID, err)
		notFound(w, r)
		return
	}

	var entryPost, entryFound Entry

	if err := d.Decode(&entryPost); err != nil {
		c.Errorf("POST /api/entries/%v: Error decoding entry post: %v", postID, err)
		internalServerError(w, r)
		return
	}

	entryID := uuid.New()
	entryCreatedDate := time.Now().UTC()
	entrySavedDate := time.Now().UTC()

	if len(entryPost.ID) != 0 {
		entryFound, _ = entryLoad(c, entryPost.ID, postID)

		if entryFound.ID != entryPost.ID {
			c.Errorf("POST /api/entries/%v: Error finding entryID: %v", postID, entryPost.ID)
			notFound(w, r)
			return
		}

		if entryFound.FinishedFlag {
			c.Errorf("POST /api/entries/%v: Entry finalized, error saving entryID %v")
			forbidden(w, r)
			return
		}

		entryID = entryFound.ID

		c.Infof("POST /api/entries/%v: Editing entryID: %v", postID, entryID)
		entryCreatedDate = entryFound.CreatedDate
	} else {
		c.Infof("POST /api/entries/%v: Creating New entryID: %v", postID, entryID)
	}

	entry := Entry{
		ID:           entryID,
		Text:         entryPost.Text,
		AuthorID:     author.UID,
		CreatedDate:  entryCreatedDate,
		SavedDate:    entrySavedDate,
		FinishedFlag: entryPost.FinishedFlag,
	}

	if err := entrySave(c, entry, postID); err != nil {
		c.Errorf("POST /api/entries/%v: Error saving entry: %v", postID, err)
		internalServerError(w, r)
		return
	}

	c.Infof("POST /api/entries/%v: Exited succesfully", postID)
	e.Encode(&entry)
}

func EntryGet(w http.ResponseWriter, r *http.Request, postID string, entryID string) {
	c := appengine.NewContext(r)
	e := json.NewEncoder(w)

	userCurrent, err := loadCurrentUser(c)
	c.Infof("GET /api/entries/%v/%v: Entered by user: %v (%v)", postID, entryID, userCurrent.Email, userCurrent.Role)

	if err != nil {
		c.Errorf("GET /api/entries/%v/%v: Error loading entry: %v", postID, entryID, err)
		notFound(w, r)
		return
	}

	notFoundError := ErrorJson{
		Message: "No Entries Found",
	}

	if user.IsAdmin(c) == false {
		c.Warningf("GET /api/entries/%v/%v: Unauthorized access by user: %v", postID, entryID, userCurrent.Email)
		forbidden(w, r)
		return
	}

	if entryID == "all" {

		entries, err := entriesLoadAll(c, postID)

		if err != nil {
			c.Errorf("GET /api/entries/%v/all: Error loading posts: %v", postID, err)
			internalServerError(w, r)
		}

		for i, _ := range entries {

			namesID := entries[i].AuthorID

			nameString, emailString, err := userGetNameString(c, namesID)

			if err != nil {
				nameString = "Name Not Found"
				emailString = "Email Not Found"
			}

			author := Author{
				Name:  nameString,
				Email: emailString,
			}

			entries[i].PostAuthor = author

		}

		if entries == nil {
			c.Infof("GET /api/entries/%v/all: No Entries Found", postID)
			e.Encode(&notFoundError)
		} else {
			c.Infof("GET /api/entries/%v/all: Exited successfully", postID)
			e.Encode(&entries)
		}
	} else if entryID == "new" {

		var entry Entry

		author := Author{
			Name:  userCurrent.DisplayName,
			Email: userCurrent.Email,
		}

		entry.PostAuthor = author

		c.Infof("GET /api/entries/%v/new: Exited successfully", postID)
		e.Encode(&entry)
	} else {
		entry, err := entryLoad(c, entryID, postID)

		if err != nil {
			c.Errorf("GET /api/entries/%v/%v: Error lookup entry: %v", postID, entryID, err)
		}

		namesID := entry.AuthorID

		nameString, emailString, err := userGetNameString(c, namesID)

		if err != nil {
			nameString = "Name Not Found"
			emailString = "Email Not Found"
		}

		author := Author{
			Name:  nameString,
			Email: emailString,
		}

		entry.PostAuthor = author

		if entry.ID != entryID {
			c.Infof("GET /api/entries/%v/%v: No Entry Found", postID, entryID)
			e.Encode(&notFoundError)
		} else {
			c.Infof("GET /api/entries/%v/%v: Exited successfully", postID, entryID)
			e.Encode(&entry)
		}
	}
}
