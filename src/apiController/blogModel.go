// Copyright 2014 Tuan.Pro. All rights reserved.

package apiController

import (
	"time"
	//"sort"
	//"net/http"
	"html/template"
	//"appengine"
	//"appengine/datastore"
	//"appengine/memcache"
	)




type Entry struct {
	Id				string
	BlogId          string
	Tags			[]string	
	PublishDate		time.Time
	Publish			bool
	Deleted			bool
	CommentsOn		bool
}

type Content struct {
	Id				string
	EntryId			string
	Title			string
	ModifiedDate	time.Time
	Author          string
	Body			template.HTML	`datastore:",noindex"`
}
	
type TagIndex struct {
	Tag				string
	BlogId			[]string
}

type Comment struct {
	Text			string
	Author			string
	Date			time.Time
	Email			string
	Publish         bool
	Deleted			bool
}
	

/*
func contentFindAll(w http.ResponseWriter, r *http.Request, contentType string, isArchive bool) (Contents, error) {
	c := appengine.NewContext(r)
		
	q := datastore.NewQuery(contentType).
		Filter("Archive =", isArchive).
		Order("DateEnd")
		
	var contents Contents
	
	_, err := q.GetAll(c, &contents)
		if err != nil {
		return nil, err
	}
	
	contents_sorted := make(Contents, 0, len(contents))
	for _, d := range contents {
		contents_sorted = append(contents_sorted, d)
	}
			
	sort.Sort(contents_sorted)
	
	//now := time.Now().Local()
	
	for i, content := range contents_sorted {

		
		contents_sorted[i].EndDate = content.DateEnd.Format("01/02")
		contents_sorted[i].StartDate = content.DateStart.Format("01/02")
		
		if content.DateStart.Format("01/02") == content.DateEnd.Format("01/02") {
			if content.DateStart.Format("03:04pm") == "12:01am" {
				contents_sorted[i].StartTime = content.DateStart.Format("01/02") + " All Day"
			} else {
				contents_sorted[i].StartTime = content.DateStart.Format("01/02 03:04pm")
			}
			
			if content.DateStart == content.DateEnd || content.DateEnd.Format("03:04pm") == "11:59pm" {
				contents_sorted[i].EndTime = ""
			} else {
				contents_sorted[i].EndTime = content.DateEnd.Format("03:04pm")
			}
			
		} else {
			if content.DateStart.Format("03:04pm") == "12:01am" {
				contents_sorted[i].StartTime = content.DateStart.Format("01/02") + " All Day"
			} else {
				contents_sorted[i].StartTime = content.DateStart.Format("01/02 03:04pm")
			}
			if content.DateEnd.Format("03:04pm") == "11:59pm" {
				contents_sorted[i].EndTime = content.DateEnd.Format("01/02") + " All Day"
			} else {
				contents_sorted[i].EndTime = content.DateEnd.Format("01/02 03:04pm")
			}
		}
		
		//if content.DateEnd.Before(now) {
		//	contents_sorted[i].Publish = false
		//}
	}
	
	return contents_sorted, nil
}

func contentLoad(w http.ResponseWriter, r *http.Request, contentType string, key string) (Content, error) {
	c := appengine.NewContext(r)
		
	var content Content
		
	findKey := datastore.NewKey(c, contentType, key, 0, nil)
	
	err := datastore.Get(c, findKey, &content)
		if err != nil {
		return content, err
	}
	
	return content, nil
}
*/


/*
func (p Contents) Len() int {
	return len(p)
}

func (p Contents) Less(i, j int) bool {
	return p[i].DateStart.Before(p[j].DateStart)
}

func (p Contents) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
*/