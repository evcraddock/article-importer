package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/ericaro/frontmatter"
)

//Article represents article information
type Article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Banner      string    `json:"banner"`
	PublishDate time.Time `json:"publishDate"`
	DataSource  string    `json:"dataSource"`
	Author      string    `json:"author"`
	Categories  []string  `json:"categories"`
	Tags        []string  `json:"tags"`
	Content     string    `json:"content"`
}

//ImportArticle represents and article that can be marshalled to yaml
type ImportArticle struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	URL         string `yaml:"url"`
	Banner      string `yaml:"banner"`
	PublishDate string `yaml:"publishDate"`
	DataSource  string `yaml:"dataSource"`
	Author      string `yaml:"author"`
	Categories  string `yaml:"categories"`
	Tags        string `yaml:"tags"`
	Content     string `fm:"content" yaml:"-"`
}

func (articleTask *Task) saveMarkdownFile(article Article) error {

	filelocation := articleTask.articleLocation + article.DataSource
	fmt.Printf("Saving Markdown file to %s\n", filelocation)

	var importfile = &ImportArticle{
		article.ID,
		article.Title,
		article.URL,
		article.Banner,
		article.PublishDate.Format("01/02/2006"),
		article.DataSource,
		article.Author,
		strings.Join(article.Categories, ", "),
		strings.Join(article.Tags, ", "),
		article.Content,
	}

	data, err := frontmatter.Marshal(importfile)
	if err != nil {
		fmt.Printf("err! %s", err.Error())
	}

	err = ioutil.WriteFile(filelocation, data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

//SaveArticle saves input data as an article and backups to a local md file
func (articleTask *Task) SaveArticle(article *Article, bypassquestions bool) (*Article, error) {
	if articleTask.service.Username == "" {
		articleTask.service.Username = AskForStringValue("Username", "", true)
	}

	if articleTask.service.Password == "" {
		articleTask.service.Password = AskForStringValue("Password", "", true)
	}

	if articleTask.service.ServiceURL == "" {
		articleTask.service.ServiceURL = AskForStringValue("Service Url", "", true)
	}

	if articleTask.service.AuthKey == "" {
		log.Fatal("AuthKey environment variable must be set.")
	}

	if article.Title == "" || bypassquestions == false {
		article.Title = AskForStringValue("Article Title", article.Title, true)
	}

	if bypassquestions == false {
		article.PublishDate = AskForDateValue("Publish Date", article.PublishDate)
	}

	if article.URL == "" || bypassquestions == false {
		article.URL = AskForStringValue("Permalink", article.URL, true)
	}

	if article.Banner == "" || bypassquestions == false {
		for {
			imageFilePath := AskForStringValue("Banner Url", article.Banner, false)

			if imageFilePath != "" {
				b, err := articleTask.service.Upload("images", imageFilePath)

				if err != nil {
					fmt.Printf("Could not save images, please try again.\n")
					continue
				}

				img := &Image{}
				json.Unmarshal(b, img)
				article.Banner = articleTask.service.ServiceURL + "/images/" + img.ID
			}

			break
		}
	}

	if article.DataSource == "" || bypassquestions == false {
		article.DataSource = AskForStringValue("Data source", article.DataSource, false)
	}

	if article.Author == "" || bypassquestions == false {
		article.Author = AskForStringValue("Author Name", article.Author, true)
	}

	if bypassquestions == false {
		article.Categories = AskForCSV("Categories (csv)", article.Categories)
	}

	if bypassquestions == false {
		article.Tags = AskForCSV("Tags (csv)", article.Tags)
	}

	requestMethod := "POST"
	requestURL := "articles"

	if article.ID != "" {
		requestMethod = "PUT"
		requestURL = "articles/" + article.ID
	}

	err := articleTask.service.SendRequest(requestMethod, requestURL, article)
	articleTask.saveMarkdownFile(*article)

	return article, err
}

//UpdateArticle updates and existing article
func (articleTask *Task) UpdateArticle(bypassQuestions bool) (*Article, error) {

	article, err := articleTask.GetArticle()

	if err != nil {
		log.Fatal(err)
	}

	return articleTask.SaveArticle(article, bypassQuestions)
}

//LoadArticle loads an existing article
func (articleTask *Task) LoadArticle(bypassQuestions bool) (*Article, error) {
	fileName := AskForStringValue("Import File location", "", false)
	var article = &Article{
		Title:       "",
		PublishDate: time.Now(),
		URL:         "",
		Banner:      "",
		DataSource:  "",
		Author:      "",
	}

	importfilename := articleTask.articleLocation + fileName
	artfile, err := ioutil.ReadFile(importfilename)
	if err != nil {
		return articleTask.SaveArticle(article, false)
	}

	importfile := new(ImportArticle)
	err = frontmatter.Unmarshal(artfile, importfile)
	if err != nil {
		fmt.Printf("Error unmarshaling yaml file: %s", err.Error())
		return articleTask.SaveArticle(article, false)
	}

	if importfile.ID != "" {
		article.ID = importfile.ID
	}

	importPublishDate, err := time.Parse("01/02/2006", importfile.PublishDate)
	if err == nil {
		article.PublishDate = importPublishDate
	}

	article.Title = importfile.Title
	article.URL = importfile.URL
	article.Author = importfile.Author

	if importfile.Banner != "" {
		article.Banner = importfile.Banner
	}

	article.DataSource = fileName
	article.Categories, _ = getStringArray(importfile.Categories)
	article.Tags, _ = getStringArray(importfile.Tags)
	article.Content = importfile.Content

	return articleTask.SaveArticle(article, bypassQuestions)
}

//DeleteArticle deletes the specified article
func (articleTask *Task) DeleteArticle() (string, error) {
	id := AskForStringValue("Article Id", "", true)
	if articleTask.service.Username == "" {
		articleTask.service.Username = AskForStringValue("Username", "", true)
	}

	if articleTask.service.Password == "" {
		articleTask.service.Password = AskForStringValue("Password", "", true)
	}

	if articleTask.service.ServiceURL == "" {
		articleTask.service.ServiceURL = AskForStringValue("Service Url", "", true)
	}

	if articleTask.service.AuthKey == "" {
		log.Fatal("AuthKey environment variable must be set.")
	}

	requestURL := "articles/" + id

	return id, articleTask.service.SendRequest("DELETE", requestURL, nil)
}

//GetArticle gets an article by datasource
func (articleTask *Task) GetArticle() (*Article, error) {
	id := AskForStringValue("Article Id", "", true)

	var article = &Article{}
	err := articleTask.service.Get("articles", id, article)

	if err != nil {
		return article, err
	}

	return article, err
}
