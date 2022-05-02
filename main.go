package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Image struct {
	Url string
}

type Person struct {
	Name           string `json:"name"`
	AlternateName  string `json:"alternateName"`
	GivenName      string `json:"givenName"`
	AdditionalName string `json:"additionalName"`
	FamilyName     string `json:"familyName"`
}

type Organisation struct {
	Name string `json:"name"`
	Logo Image  `json:"logo"`
}

type BlogPosting struct {
	Context          string       `json:"@context"`
	Type             string       `json:"@type"`
	Name             string       `json:"name"`
	DateCreated      string       `json:"dateCreated"`
	DatePublished    string       `json:"datePublished"`
	DateModified     string       `json:"dateModified"`
	WordCount        int64        `json:"wordCount"`
	Author           Person       `json:"author"`
	Url              string       `json:"url"`
	MainEntityOfPage string       `json:"mainEntityOfPage"`
	InLanguage       string       `json:"inLanguage"`
	CopyrightHolder  Person       `json:"copyrightHolder"`
	Publisher        Organisation `json:"publisher"`
	Headline         string       `json:"headline"`
	License          string       `json:"license"`
	Image            Image        `json:"image"`
	SameAs           string       `json:"sameAs"`
	ArticleBody      string       `json:"articleBody"`
}

type ApiBlogPostingResponse struct {
	Data BlogPosting `json:"data"`
}

type Blog struct {
	Context   string        `json:"@context"`
	Type      string        `json:"@type"`
	BlogPosts []BlogPosting `json:"blogPosts"`
}

type ApiBlogResponse struct {
	Data Blog `json:"data"`
}

func formatDate(date string) string {
	dateHumanReadable, _ := time.Parse("2006-01-02T15:04:05-07:00", date)

	return dateHumanReadable.Format("02 January 2006")
}

func urlPath(urlString string) string {
	netUrl, _ := url.Parse(urlString)

	return netUrl.Path
}

func main() {
	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"formatDate": formatDate,
		"urlPath":    urlPath,
	})

	router.LoadHTMLGlob("./templates/*/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/posts", func(c *gin.Context) {
		content := ApiBlogResponse{}
		jsonErr := json.Unmarshal(getResponse("https://api.elliotjreed.com/blog/posts"), &content)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		posts := content.Data.BlogPosts

		c.HTML(http.StatusOK, "posts.html", gin.H{
			"posts": posts,
		})
	})

	router.GET("/blog/:date/:link", func(c *gin.Context) {
		url := "https://api.elliotjreed.com/blog/post/" + c.Param("date") + "/" + c.Param("link")

		content := ApiBlogPostingResponse{}
		jsonErr := json.Unmarshal(getResponse(url), &content)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		var buf bytes.Buffer
		blogPosting := content.Data
		markdownError := goldmark.Convert([]byte(blogPosting.ArticleBody), &buf)

		if markdownError != nil {
			log.Fatal(markdownError)
		}

		dateCreated := blogPosting.DateCreated
		c.HTML(http.StatusOK, "post.html", gin.H{
			"article":           template.HTML(buf.String()),
			"dateHumanReadable": formatDate(blogPosting.DateCreated),
			"date":              dateCreated,
			"headline":          blogPosting.Headline,
			"canonicalUrl":      "https://www.elliotjreed.com/blog/" + c.Param("date") + "/" + c.Param("link"),
			"wordCount":         blogPosting.WordCount,
		})
	})

	router.Run("0.0.0.0:98")
}

func getResponse(url string) []byte {
	client := http.Client{
		Timeout: time.Second * 10, // Timeout after 10 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body
}
