// Copyright 2015 Tuan.Pro. All rights reserved.

package apiController

import (
	//"time"
	//"sort"
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"net/http"
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
	SortMethod     string   `json:"sortMethod"`
}

type Blogs []BlogIndex

type Author struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
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

	c.Infof("POST /api/blogs: Position %v set to %v", setBlog.ID, newPosition)

	if blogs != nil {
		for k, v := range blogs {
			if v.Position == newPosition {
				//log.Println(blogs[k].ID, "has overlapping position")
				if err2 := blogIndexSetPosition(c, blogs[k].Position+k+1, blogs[k]); err2 != nil {
					c.Errorf("POST /api/blogs: Error setting %v position to %v", blogs[k].ID, err2)
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

	author, err := loadCurrentUser(c)
	c.Infof("POST /api/blogs: Entered by user: %v (%v)", author.Email, author.Role)

	if user.IsAdmin(c) == false {
		c.Warningf("GET /api/blogs: Unauthorized access by user: %v", author.Email)
		forbidden(w, r)
		return
	}

	if err != nil {
		c.Errorf("POST /api/blogs: Error loading user: %v", err)
		notFound(w, r)
		return
	}

	var blogIndexPost, blogIndexFound BlogIndex

	if err := d.Decode(&blogIndexPost); err != nil {
		c.Errorf("POST /api/blogs: Error decoding blog post: %v", err)
		internalServerError(w, r)
		return
	}

	blogID := uuid.New()

	if len(blogIndexPost.ID) != 0 {
		blogIndexFound, _ = blogIndexLoad(c, blogIndexPost.ID)

		if blogIndexFound.ID != blogIndexPost.ID {
			c.Errorf("POST /api/blogs: Error finding blogID: %v", blogIndexPost.ID)
			notFound(w, r)
			return
		}

		blogID = blogIndexFound.ID
		c.Infof("POST /api/blogs: Editing blogID: %v", blogID)
	} else {
		c.Infof("POST /api/blogs: Creating New blogID: %v", blogID)
	}

	var blogAuthorsID []string

	for i, _ := range blogIndexPost.BlogAuthors {

		author := blogIndexPost.BlogAuthors[i]
		user, err := findUser(c, author.Email)

		if err != nil || user.UID == "" {
			c.Errorf("POST /api/blogs: Error author not found: %v", author.Email)
		} else {
			if !stringInSlice(user.UID, blogAuthorsID) {
				blogAuthorsID = append(blogAuthorsID, user.UID)
			} else {
				c.Infof("POST /api/blogs: Duplicate author: %v", author.Email)
			}
		}
	}

	if len(blogAuthorsID) == 0 {
		blogAuthorsID = append(blogAuthorsID, author.UID)
	}

	setNewPosition := false

	if blogIndexPost.Position == 0 {

		c.Infof("POST /api/blogs: Position set to 0")
		lastPosition, err := blogIndexLastPosition(c)

		if err != nil {
			c.Errorf("POST /api/blogs: Error getting last index: %v", err)
			internalServerError(w, r)
		}
		blogIndexPost.Position = lastPosition + 10
	} else if blogIndexPost.Position != blogIndexFound.Position {
		setNewPosition = true
	}

	if blogIndexPost.SortMethod == "" {
		if blogIndexFound.SortMethod == "" {
			blogIndexPost.SortMethod = "1"
		} else {
			blogIndexPost.SortMethod = blogIndexFound.SortMethod
		}
	}

	blogIndex := BlogIndex{
		ID:             blogID,
		Name:           blogIndexPost.Name,
		AuthorsID:      blogAuthorsID,
		CommentsAllow:  blogIndexPost.CommentsAllow,
		CommentsReview: blogIndexPost.CommentsReview,
		ActiveFlag:     blogIndexPost.ActiveFlag,
		Position:       blogIndexPost.Position,
		SortMethod:     blogIndexPost.SortMethod,
	}

	if setNewPosition {
		if err := blogIndexSetPosition(c, blogIndexPost.Position, blogIndex); err != nil {
			c.Errorf("POST /api/blogs: Error saving blog: %v", err)
			internalServerError(w, r)
			return
		}
	} else {
		if err := blogIndexSave(c, blogIndex); err != nil {
			c.Errorf("POST /api/blogs: Error saving blog: %v", err)
			internalServerError(w, r)
			return
		}
	}

	c.Infof("POST /api/blogs: Exited succesfully")
	e.Encode(&blogIndex)
}

func BlogsIndexGet(w http.ResponseWriter, r *http.Request, blogID string) {
	c := appengine.NewContext(r)
	e := json.NewEncoder(w)

	userCurrent, err := loadCurrentUser(c)
	c.Infof("GET /api/blogs/%v: Entered by user: %v (%v)", blogID, userCurrent.Email, userCurrent.Role)

	if err != nil {
		c.Errorf("GET /api/blogs: Error loading user: %v", err)
		notFound(w, r)
		return
	}

	notFoundError := ErrorJson{
		Message: "No Blogs Found",
	}

	if user.IsAdmin(c) == false {
		c.Warningf("GET /api/blogs: Unauthorized access by user: %v", userCurrent.Email)
		forbidden(w, r)
		return
	}

	if blogID == "all" {
		blogsIndex, err := blogsIndexLoadAll(c)

		if err != nil {
			c.Errorf("GET /api/blogs/all: Error loading blogs: %v", err)
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
		}

		if blogsIndex == nil {
			c.Infof("GET /api/blogs/all: No Blogs Found")
			e.Encode(&notFoundError)
		} else {
			c.Infof("GET /api/blogs/all: Exited successfully")
			e.Encode(&blogsIndex)
		}
	} else if blogID == "new" {
		var blog BlogIndex

		author := Author{
			Name:  userCurrent.DisplayName,
			Email: userCurrent.Email,
		}

		blog.BlogAuthors = append(blog.BlogAuthors, author)
		blog.SortMethod = "1"

		c.Infof("GET /api/blogs/new: Exited successfully")

		e.Encode(&blog)
	} else {
		blogIndex, err := blogIndexLoad(c, blogID)

		if err != nil {
			c.Errorf("GET /api/blogs/%v: Error lookup blog: %v", blogID, err)
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

		if blogIndex.ID != blogID {
			c.Infof("GET /api/blogs/%v: No Blog Found", blogID)
			e.Encode(&notFoundError)
		} else {
			c.Infof("GET /api/blogs/%v: Exited successfully", blogID)
			e.Encode(&blogIndex)
		}
	}
}
