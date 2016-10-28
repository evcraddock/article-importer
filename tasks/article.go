package tasks

import (
	//"fmt"
	"log"
	"encoding/csv"
	"strings"
	"time"
)

type Article struct {
	Id				string			`json:"id"`
	Title			string			`json:"title"`
	Url 			string			`json:"url"`
	Banner			string 			`json:"banner"`
	PublishDate		time.Time 		`json:"publishDate"`
	DataSource		string			`json:"dataSource"`
	Author			string			`json:"author"`
	Categories		[]string		`json:"categories"`
	Tags			[]string 		`json:"tags"`
}

func (this *Task) SaveArticle(article *Article) (*Article, error) {
	if this.service.Username == "" {
		this.service.Username = AskForStringValue("Username", "")
	}

	if this.service.Password == "" {
		this.service.Password = AskForStringValue("Password", "")
	}

	if this.service.ServiceUrl == "" {
		this.service.ServiceUrl = AskForStringValue("Service Url", "")
	}

	if this.service.AuthKey == "" {
		log.Fatal("AuthKey environment variable must be set.")
	}

	article.Title = AskForStringValue("Article Title", article.Title)
	article.PublishDate = AskForDateValue("Publish Date", article.PublishDate)
	article.Url = AskForStringValue("Permalink", article.Url)
	article.Banner = AskForStringValue("Banner Url", article.Banner)
	article.DataSource = AskForStringValue("Data source", article.DataSource)
	article.Author = AskForStringValue("Author Name", article.Author)

	current_categories := strings.Join(article.Categories, ", ")
	current_tags := strings.Join(article.Tags, ", ")

	categories := AskForStringValue("Categories (csv)", current_categories)
	tags := AskForStringValue("Tags (csv)", current_tags)

	r := csv.NewReader(strings.NewReader(categories))
	article.Categories, _ = r.Read()

	r = csv.NewReader(strings.NewReader(tags))
	article.Tags, _ = r.Read()	

	requestMethod := "POST"
	requestUrl := "articles"

	if article.Id != "" {
		requestMethod = "PUT"
		requestUrl = "articles/" + article.Id
	}

	err := this.service.SendRequest(requestMethod, requestUrl, article)

	return article, err
}

func (this *Task) UpdateArticle() (*Article, error) {
	
	article, err := this.GetArticle()

	if err != nil {
		log.Fatal(err)
	}

	return this.SaveArticle(article)
}

func (this *Task) CreateNewArticle() (*Article, error) {
	var article *Article = &Article{
		Title: "",
		PublishDate: time.Now(),
		Url: "",
		Banner: "/images/articles/bronco_stadium.jpg",
		DataSource: "",
		Author: "",
	}

	return this.SaveArticle(article)
}

func (this *Task) DeleteArticle() (string, error) {
	id := AskForStringValue("Article Id", "")
	if this.service.Username == "" {
		this.service.Username = AskForStringValue("Username", "")
	}

	if this.service.Password == "" {
		this.service.Password = AskForStringValue("Password", "")
	}

	if this.service.ServiceUrl == "" {
		this.service.ServiceUrl = AskForStringValue("Service Url", "")
	}

	if this.service.AuthKey == "" {
		log.Fatal("AuthKey environment variable must be set.")
	}

	requestUrl := "articles/" + id

	return id, this.service.SendRequest("DELETE", requestUrl, nil)
}

func (this *Task) GetArticle() (*Article, error) {
	id := AskForStringValue("Article Id", "")

	var article *Article = &Article{}
	err := this.service.GetJson("articles", id, article)

	if err != nil {
		return article, err
	}

	return article, err
}