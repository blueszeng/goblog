// Copyright 2015 Tuan.Pro. All rights reserved.

package apiController

import (
	//"time"
	//"sort"
	"log"
	"net/http"
	//"html/template"
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
	//"appengine/memcache"
	"code.google.com/p/go-uuid/uuid"
)

type BlogIndex struct {
	ID             string   `json:"id"`
	Name           string   `json:"blogName"`
	AuthorsID      []string `json:"-"`
	BlogAuthors    []Author `json:"blogAuthors"`
	CommentsAllow  bool     `json:"commentsAllow"`
	CommentsReview bool     `json:"commentsReview"`
	ActiveFlag     bool     `json:"active"`
	Position       int      `json:"position"`
}

type Blogs []BlogIndex

type Author struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func blogIndexSave(c appengine.Context, blog BlogIndex) error {
	k := datastore.NewKey(c, "BlogIndex", blog.ID, 0, nil)

	if _, err := datastore.Put(c, k, &blog); err != nil {
		return err
	}

	return nil
}

func blogsIndexLoadAll(c appengine.Context) (Blogs, error) {
	q := datastore.NewQuery("BlogIndex").
		Order("Name")

	var blogs Blogs

	if _, err := q.GetAll(c, &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

func blogIndexLoad(c appengine.Context, blogID string) (BlogIndex, error) {
	k := datastore.NewKey(c, "BlogIndex", blogID, 0, nil)

	var blog BlogIndex

	if err := datastore.Get(c, k, &blog); err != nil {
		return blog, err
	}

	return blog, nil
}

func blogIndexLastPosition(c appengine.Context) (int, error) {
	q := datastore.NewQuery("BlogIndex").
		Project("Position").
		Order("-Position").Limit(1)

	var blogs Blogs
	_, err := q.GetAll(c, &blogs)

	if err != nil {
		return 0, err
	}
	log.Println(blogs)

	if blogs == nil {
		return 0, nil
	} else {
		return blogs[0].Position, nil
	}
}

func blogIndexSetPosition(c appengine.Context, newPosition int, setBlog BlogIndex) error {
	q := datastore.NewQuery("BlogIndex").
		Filter("Position =", newPosition)

	var blogs Blogs

	_, err := q.GetAll(c, &blogs)

	if err != nil {
		return err
	}

	setBlog.Position = newPosition

	if err1 := blogIndexSave(c, setBlog); err1 != nil {
		return err1
	}
	log.Println("Blog Position setting", setBlog.ID, "to position", newPosition)

	if blogs != nil {
		for k, v := range blogs {
			if v.Position == newPosition {
				log.Println(blogs[k].ID, "has overlapping position")
				if err2 := blogIndexSetPosition(c, blogs[k].Position+k+1, blogs[k]); err2 != nil {
					log.Println("Error setting postion for ", blogs[k].ID, err2)
				}
			}
		}
	}
	return nil
}

func BlogIndexPost(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	d := json.NewDecoder(r.Body)
	e := json.NewEncoder(w)

	if user.IsAdmin(c) == false {
		forbidden(w, r)
		return
	}

	author, err := loadCurrentUser(c)

	if err != nil {
		log.Println("No User Found: ", err)
		notFound(w, r)
		return
	}

	var blogIndexPost, blogIndexFound BlogIndex

	if err := d.Decode(&blogIndexPost); err != nil {
		log.Println("Error decoding blog index post: ", err)
		internalServerError(w, r)
		return
	}

	blogID := uuid.New()

	if len(blogIndexPost.ID) != 0 {
		blogIndexFound, _ = blogIndexLoad(c, blogIndexPost.ID)

		if blogIndexFound.ID != blogIndexPost.ID {
			log.Println("Error finding blog ID")
			notFound(w, r)
			return
		}
		blogID = blogIndexFound.ID
	}

	var blogAuthorsID []string

	for i, _ := range blogIndexPost.BlogAuthors {

		author := blogIndexPost.BlogAuthors[i]
		user, err := findUser(c, author.Email)

		if err != nil || user.UID == "" {
			log.Println("not saved", author.Email)
		} else {
			if !stringInSlice(user.UID, blogAuthorsID) {
				log.Println("saved", author.Email)
				blogAuthorsID = append(blogAuthorsID, user.UID)
			} else {
				log.Println("duplicate", author.Email)
			}
		}
	}

	if len(blogAuthorsID) == 0 {
		blogAuthorsID = append(blogAuthorsID, author.UID)
	}

	setNewPosition := false

	if blogIndexPost.Position == 0 {
		log.Println("No Position Passed")
		lastPosition, err := blogIndexLastPosition(c)

		if err != nil {
			log.Println("Error getting last index", err)
			internalServerError(w, r)
		}
		blogIndexPost.Position = lastPosition + 10
	} else if blogIndexPost.Position != blogIndexFound.Position {
		setNewPosition = true
	}

	blogIndex := BlogIndex{
		ID:             blogID,
		Name:           blogIndexPost.Name,
		AuthorsID:      blogAuthorsID,
		CommentsAllow:  blogIndexPost.CommentsAllow,
		CommentsReview: blogIndexPost.CommentsReview,
		ActiveFlag:     blogIndexPost.ActiveFlag,
		Position:       blogIndexPost.Position,
	}

	if setNewPosition {
		if err := blogIndexSetPosition(c, blogIndexPost.Position, blogIndex); err != nil {
			log.Println("Error saving user: ", err)
			internalServerError(w, r)
			return
		}
	} else {
		if err := blogIndexSave(c, blogIndex); err != nil {
			log.Println("Error saving user: ", err)
			internalServerError(w, r)
			return
		}
	}

	e.Encode(&blogIndex)
}

func BlogsIndexGet(w http.ResponseWriter, r *http.Request, blogID string) {
	c := appengine.NewContext(r)
	e := json.NewEncoder(w)

	log.Println("GET /api/blogs entered")
	notFoundError := ErrorJson{
		Message: "No Blogs Found",
	}

	if user.IsAdmin(c) == false {
		forbidden(w, r)
		return
	}

	if blogID == "all" {

		blogsIndex, err := blogsIndexLoadAll(c)

		if err != nil {
			log.Println("Error loading blogs: ", err)
			internalServerError(w, r)
		}

		for i, _ := range blogsIndex {

			namesID := blogsIndex[i].AuthorsID

			for _, j := range namesID {
				nameString, emailString, err := userGetNameString(c, j)

				if err != nil {
					nameString = "Name Not Found"
					emailString = "Email Not Found"
				}

				author := Author{
					Name:  nameString,
					Email: emailString,
				}

				blogsIndex[i].BlogAuthors = append(blogsIndex[i].BlogAuthors, author)
			}

			//log.Println(blogsIndex[i].Authors)
		}

		if blogsIndex == nil {
			e.Encode(&notFoundError)
		} else {
			e.Encode(&blogsIndex)
		}
	} else if blogID == "new" {

		var blog BlogIndex
		currentUser, err := loadCurrentUser(c)

		if err != nil {
			log.Println(err)
			internalServerError(w, r)
			return
		}

		author := Author{
			Name:  currentUser.DisplayName,
			Email: currentUser.Email,
		}

		blog.BlogAuthors = append(blog.BlogAuthors, author)
		e.Encode(&blog)
	} else {
		log.Println("BlogID = ", blogID)
		blogIndex, err := blogIndexLoad(c, blogID)

		if err != nil {
			log.Println("GET /api/blogs error", err)
		}

		namesID := blogIndex.AuthorsID

		for _, j := range namesID {
			nameString, emailString, err := userGetNameString(c, j)

			if err != nil {
				nameString = "Name Not Found"
				emailString = "Email Not Found"
			}

			author := Author{
				Name:  nameString,
				Email: emailString,
			}

			blogIndex.BlogAuthors = append(blogIndex.BlogAuthors, author)
		}

		e.Encode(&blogIndex)
	}
}
