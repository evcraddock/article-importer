package tasks

import (
	//"fmt"
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

func (this *Task) CreateNewArticle() (*Article, error) {
	title := AskForStringValue("Article Title", "")
	publishDate := AskForDateValue("Publish Date")
	articleUrl := AskForStringValue("Permalink", "")
	bannerUrl := AskForStringValue("Banner Url", "/images/articles/bronco_stadium.jpg")
	dataSource := AskForStringValue("Data source", "")
	author := AskForStringValue("Author Name", "")
	categories := AskForStringValue("Categories (csv)", "")
	tags := AskForStringValue("Tags (csv)", "")

	var article *Article = &Article{
		Title: title,
		PublishDate: publishDate,
		Url: articleUrl,
		Banner: bannerUrl,
		DataSource: dataSource,
		Author: author,
	}

	r := csv.NewReader(strings.NewReader(categories))
	article.Categories, _ = r.Read()

	r = csv.NewReader(strings.NewReader(tags))
	article.Tags, _ = r.Read()	

	err := this.service.PostJson("articles", article)
	return article, err 
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