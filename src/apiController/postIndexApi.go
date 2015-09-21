// Copyright 2015 Tuan.Pro. All rights reserved.

package apiController

import (
	"time"
	//"sort"
	"log"
	"net/http"
	//"strconv"
	//"html/template"
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
	//"appengine/memcache"
	"code.google.com/p/go-uuid/uuid"
)

type PostIndex struct {
	ID          string    `json:"id"`
	Name        string    `json:"postName"`
	AuthorID    string    `json:"-"`
	PostAuthor  Author    `json:"postAuthor"`
	CreatedDate time.Time `json:"createdDate"`
	ActiveFlag  bool      `json:"active"`
	StopFlag    bool      `json:"stopFlag"`
	Position    int       `json:"position"`
	PostDate    time.Time `json:"postDate"`
	StopDate    time.Time `json:"stopDate"`
	PostDateStr string    `json:"postDateStr"`
	StopDateStr string    `json:"stopDateStr"`
	DateLoc     string    `json:"dateLoc"`
	Displayed   bool      `json:"displayed"`
}

type Posts []PostIndex

func postIndexSave(c appengine.Context, post PostIndex, blogID string) error {
	blogKey := datastore.NewKey(c, "BlogIndex", blogID, 0, nil)
	k := datastore.NewKey(c, "PostIndex", post.ID, 0, blogKey)

	if _, err := datastore.Put(c, k, &post); err != nil {
		return err
	}

	return nil
}

func postsIndexLoadAll(c appengine.Context, blogID string) (Posts, error) {
	blogKey := datastore.NewKey(c, "BlogIndex", blogID, 0, nil)
	q := datastore.NewQuery("PostIndex").
		Ancestor(blogKey).
		Order("Name")

	var posts Posts

	if _, err := q.GetAll(c, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func postIndexLoad(c appengine.Context, postID string, blogID string) (PostIndex, error) {
	blogKey := datastore.NewKey(c, "BlogIndex", blogID, 0, nil)
	k := datastore.NewKey(c, "PostIndex", postID, 0, blogKey)

	var post PostIndex

	if err := datastore.Get(c, k, &post); err != nil {
		return post, err
	}

	return post, nil
}

func postIndexLastPosition(c appengine.Context) (int, error) {
	q := datastore.NewQuery("PostIndex").
		Project("Position").
		Order("-Position").Limit(1)

	var posts Posts
	_, err := q.GetAll(c, &posts)

	if err != nil {
		return 0, err
	}
	log.Println(posts)

	if posts == nil {
		return 0, nil
	} else {
		return posts[0].Position, nil
	}
}

func postIndexSetPosition(c appengine.Context, newPosition int, setPost PostIndex, blogID string) error {
	q := datastore.NewQuery("PostIndex").
		Filter("Position =", newPosition)

	var posts Posts

	_, err := q.GetAll(c, &posts)

	if err != nil {
		return err
	}

	setPost.Position = newPosition

	if err1 := postIndexSave(c, setPost, blogID); err1 != nil {
		return err1
	}
	log.Println("Post Position setting", setPost.ID, "to position", newPosition)

	if posts != nil {
		for k, v := range posts {
			if v.Position == newPosition {
				log.Println(posts[k].ID, "has overlapping position")
				if err2 := postIndexSetPosition(c, posts[k].Position+k+1, posts[k], blogID); err2 != nil {
					log.Println("Error setting postion for ", posts[k].ID, err2)
				}
			}
		}
	}
	return nil
}

func PostIndexPost(w http.ResponseWriter, r *http.Request, blogID string) {
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

	var postIndexPost, postIndexFound PostIndex

	if err := d.Decode(&postIndexPost); err != nil {
		log.Println("Error decoding postIndex post: ", err)
		internalServerError(w, r)
		return
	}

	postID := uuid.New()
	postCreatedDate := time.Now().UTC()
	postPostDate := time.Now().UTC()
	postStopDate := time.Now().UTC()

	if len(postIndexPost.ID) != 0 {
		postIndexFound, _ = postIndexLoad(c, postIndexPost.ID, blogID)

		if postIndexFound.ID != postIndexPost.ID {
			log.Println("Error finding post ID")
			notFound(w, r)
			return
		}
		postID = postIndexFound.ID
		postCreatedDate = postIndexFound.CreatedDate
	}

	setNewPosition := false

	if postIndexPost.Position == 0 {
		log.Println("No Position Passed")
		lastPosition, err := postIndexLastPosition(c)

		if err != nil {
			log.Println("Error getting last index", err)
			internalServerError(w, r)
		}
		postIndexPost.Position = lastPosition + 10
	} else if postIndexPost.Position != postIndexFound.Position {
		setNewPosition = true
	}

	log.Println(postIndexPost.PostDateStr)

	loc, _ := time.LoadLocation("America/Chicago")

	if postIndexPost.PostDateStr != "" {
		postPostDate, _ = time.ParseInLocation("01-02-2006", postIndexPost.PostDateStr, loc)
	}

	if postIndexPost.StopDateStr != "" {
		postStopDate, _ = time.ParseInLocation("01-02-2006", postIndexPost.StopDateStr, loc)
	}

	postIndex := PostIndex{
		ID:          postID,
		Name:        postIndexPost.Name,
		AuthorID:    author.UID,
		ActiveFlag:  postIndexPost.ActiveFlag,
		StopFlag:    postIndexPost.StopFlag,
		Position:    postIndexPost.Position,
		CreatedDate: postCreatedDate,
		PostDate:    postPostDate,
		StopDate:    postStopDate,
	}

	if setNewPosition {
		if err := postIndexSetPosition(c, postIndexPost.Position, postIndex, blogID); err != nil {
			log.Println("Error saving user: ", err)
			internalServerError(w, r)
			return
		}
	} else {
		if err := postIndexSave(c, postIndex, blogID); err != nil {
			log.Println("Error saving user: ", err)
			internalServerError(w, r)
			return
		}
	}

	e.Encode(&postIndex)
}

func PostsIndexGet(w http.ResponseWriter, r *http.Request, blogID string, postID string) {
	c := appengine.NewContext(r)
	e := json.NewEncoder(w)

	log.Println("GET /api/posts entered")
	notFoundError := ErrorJson{
		Message: "No Posts Found",
	}

	if user.IsAdmin(c) == false {
		forbidden(w, r)
		return
	}

	if postID == "all" {

		postsIndex, err := postsIndexLoadAll(c, blogID)

		if err != nil {
			log.Println("Error loading posts: ", err)
			internalServerError(w, r)
		}

		for i, _ := range postsIndex {

			namesID := postsIndex[i].AuthorID

			nameString, emailString, err := userGetNameString(c, namesID)

			if err != nil {
				nameString = "Name Not Found"
				emailString = "Email Not Found"
			}

			author := Author{
				Name:  nameString,
				Email: emailString,
			}

			postsIndex[i].PostAuthor = author

			if postsIndex[i].PostDate.Before(time.Now()) && postsIndex[i].ActiveFlag {
				if postsIndex[i].StopFlag && postsIndex[i].StopDate.Before(time.Now()) {
					postsIndex[i].Displayed = false
				} else {
					postsIndex[i].Displayed = true
				}
			} else {
				postsIndex[i].Displayed = false
			}
		}

		//log.Println(blogsIndex[i].Authors)

		if postsIndex == nil {
			e.Encode(&notFoundError)
		} else {
			e.Encode(&postsIndex)
		}
	} else if postID == "new" {

		var post PostIndex
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

		post.PostAuthor = author

		e.Encode(&post)
	} else {
		log.Println("PostID = ", postID)
		postIndex, err := postIndexLoad(c, postID, blogID)

		if err != nil {
			log.Println("GET /api/posts error", err)

		}

		namesID := postIndex.AuthorID

		nameString, emailString, err := userGetNameString(c, namesID)

		if err != nil {
			nameString = "Name Not Found"
			emailString = "Email Not Found"
		}

		author := Author{
			Name:  nameString,
			Email: emailString,
		}

		postIndex.PostAuthor = author

		postIndex.PostDateStr = postIndex.PostDate.Format("01-02-2006")
		postIndex.StopDateStr = postIndex.StopDate.Format("01-02-2006")

		if postIndex.ID != postID {
			e.Encode(&notFoundError)
		} else {
			e.Encode(&postIndex)
		}
	}
}
