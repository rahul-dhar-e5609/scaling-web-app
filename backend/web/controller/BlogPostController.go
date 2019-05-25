package controller

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/IAmRDhar/scaling-web-app/backend/entity"
	"github.com/IAmRDhar/scaling-web-app/backend/web/model"
)

type BlogPostController struct {
	blogListTemplate *template.Template
	blogTemplate     *template.Template
}

func (c *BlogPostController) showBlogList(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	cacheKey := url.QueryEscape(r.URL.RequestURI())
	resp, ok := getFromCache(cacheKey)
	if ok {
		io.Copy(w, resp)
		resp.Close()
		return
	}

	posts, err := model.GetLastPosts(3)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	titles, err := model.GetLastPostTitles(10)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	context := map[string]interface{}{
		"posts":  posts,
		"titles": titles,
	}

	buf := bytes.Buffer{}
	c.blogListTemplate.Execute(&buf, context)
	data := buf.Bytes()
	w.Write(data)
	go savingToCache(cacheKey, int64(24*time.Hour), data[:])

}

func (c *BlogPostController) showBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	cacheKey := url.QueryEscape(r.URL.RequestURI())
	resp, ok := getFromCache(cacheKey)
	if ok {
		io.Copy(w, resp)
		resp.Close()
		return
	}

	matches := postPath.FindStringSubmatch(r.URL.Path)

	//no need to check for error since regex guarantees an integer value
	postID, _ := strconv.Atoi(matches[1])

	post, err := model.GetPostById(postID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	titles, err := model.GetLastPostTitles(10)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	context := map[string]interface{}{
		"post":   post,
		"titles": titles,
	}

	buf := bytes.Buffer{}
	c.blogListTemplate.Execute(&buf, context)
	data := buf.Bytes()
	w.Write(data)
	go savingToCache(cacheKey, int64(24*time.Hour), data[:])
}

func (c *BlogPostController) createBlogPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	now := time.Now()
	post := &entity.BlogPost{
		ContentItem: entity.ContentItem{
			Subject:     r.FormValue("subject"),
			Body:        r.FormValue("body"),
			Author:      nil,
			Comments:    []entity.Comment{},
			CreatedDate: &now,
			PublishDate: nil,
			IsPublished: false,
		},
	}

	post, err := model.CreateBlogPost(post)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("Location", "/posts/"+strconv.Itoa(post.ID))
	w.WriteHeader(http.StatusSeeOther)
}

func (c *BlogPostController) updateBlogPost(w http.ResponseWriter, r *http.Request) {

	matches := postPath.FindStringSubmatch(r.URL.Path)

	//no need to check for error since regex guarantees an integer value
	postID, _ := strconv.Atoi(matches[1])

	r.ParseForm()
	now := time.Now()
	post := &entity.BlogPost{
		ContentItem: entity.ContentItem{
			ID:          postID,
			Subject:     r.FormValue("subject"),
			Body:        r.FormValue("body"),
			Author:      nil,
			Comments:    []entity.Comment{},
			CreatedDate: &now,
			PublishDate: nil,
			IsPublished: false,
		},
	}

	post, err := model.UpdateBlogPost(post)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("Location", "/posts/"+strconv.Itoa(post.ID))
	w.WriteHeader(http.StatusSeeOther)
	go invalideCacheEntry(url.QueryEscape(r.RequestURI))
}
