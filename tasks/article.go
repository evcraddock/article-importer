package tasks

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ericaro/frontmatter"
)

//Article represents article information
type Article struct {
	ID          string    `json:"id" `
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Images      []string  `json:"images"`
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
	ID          string   `yaml:"id,omitempty"`
	Title       string   `yaml:"title"`
	URL         string   `yaml:"url"`
	Images      []string `yaml:"images"`
	Banner      string   `yaml:"banner"`
	PublishDate string   `yaml:"publishDate"`
	DataSource  string   `yaml:"dataSource"`
	Author      string   `yaml:"author"`
	Categories  []string `yaml:"categories"`
	Tags        []string `yaml:"tags"`
	Content     string   `fm:"content" yaml:"-"`
}

//HugoArticle represents and article that can be marshalled to yaml
type HugoArticle struct {
	Title      string   `yaml:"title"`
	URL        string   `yaml:"url"`
	Banner     string   `yaml:"banner"`
	Date       string   `yaml:"date"`
	Author     string   `yaml:"author"`
	Categories []string `yaml:"categories"`
	Tags       []string `yaml:"tags"`
	Layout     string   `yaml:"layout"`
	Content    string   `fm:"content" yaml:"-"`
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
func (articleTask *Task) GetArticle(id string) (*Article, error) {
	if id == "" {
		id = AskForStringValue("Article Id", "", true)
	}

	var article = &Article{}
	err := articleTask.service.Get("articles", id, article)

	if err != nil {
		return article, err
	}

	return article, err
}

//LoadArticle loads an existing article
func (articleTask *Task) LoadArticle(fileName string, bypassQuestions bool) (*Article, error) {
	if fileName == "" {
		fileName = AskForStringValue("Import File location", "", false)
	}

	var article = &Article{
		Title:       "",
		PublishDate: time.Now(),
		URL:         "",
		Banner:      "",
		DataSource:  "",
		Author:      "",
	}

	artfile, err := ioutil.ReadFile(fileName)

	if err != nil || len(artfile) == 0 {
		return nil, fmt.Errorf("Could not open file")
	}

	importfile := new(ImportArticle)
	err = frontmatter.Unmarshal(artfile, importfile)
	if err != nil {
		msg := fmt.Errorf("Error unmarshaling yaml file: %s", err.Error())
		return nil, msg
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

	article.Banner = importfile.Banner
	article.DataSource = fileName
	article.Categories = importfile.Categories
	article.Tags = importfile.Tags
	article.Images = importfile.Images
	article.Content = importfile.Content

	return articleTask.SaveArticle(article, bypassQuestions)
}

//ImportArticles imports list of articles in path
func (articleTask *Task) ImportArticles(filedir string) error {
	if filedir == "" {
		filedir = AskForStringValue("Import File or Folder", "", false)
	}

	isdir, err := isDirectory(filedir)

	if err != nil {
		return err
	}

	if isdir {
		files, err := ioutil.ReadDir(filedir)
		if err != nil {
			return err
		}

		for _, f := range files {

			filename := f.Name()
			extension := filepath.Ext(filename)
			// fmt.Printf("filename: %s ext: %s \n ", filename, extension)

			if extension == ".md" {
				importfilepath := filedir + "/" + filename
				fmt.Printf("importing file: %s \n ", importfilepath)

				_, err := articleTask.ImportArticle(importfilepath)
				if err != nil {
					fmt.Printf("error: %s \n ", err.Error())
					return err
				}
			}
		}
	} else {
		_, err := articleTask.ImportArticle(filedir)
		return err
	}

	return nil
}

//ImportArticle loads an existing article
func (articleTask *Task) ImportArticle(fileName string) (*Article, error) {
	if fileName == "" {
		fileName = AskForStringValue("Import File location", "", false)
	}

	var article = &Article{
		Title:       "",
		PublishDate: time.Now(),
		URL:         "",
		Banner:      "",
		DataSource:  "",
		Author:      "",
	}

	artfile, err := ioutil.ReadFile(fileName)

	if err != nil || len(artfile) == 0 {
		return nil, fmt.Errorf("Could not open file")
	}

	importfile := new(HugoArticle)
	err = frontmatter.Unmarshal(artfile, importfile)
	if err != nil {
		msg := fmt.Errorf("Error unmarshaling yaml file: %s", err.Error())
		return nil, msg
	}

	articleurl := GetFileName(importfile.URL, "/")

	articlepath := filepath.Dir(fileName)
	newarticlepath := articlepath + "/" + articleurl

	if _, err := os.Stat(newarticlepath); os.IsNotExist(err) {
		err = os.Mkdir(newarticlepath, 0755)
		if err != nil {
			msg := fmt.Errorf("Error creating directory: %s /\n ", err.Error())
			return nil, msg
		}
	}

	importPublishDate, err := time.Parse("2006-01-02", importfile.Date)
	if err == nil {
		article.PublishDate = importPublishDate
	}

	article.Title = importfile.Title
	article.URL = articleurl + ".md"
	article.Author = importfile.Author

	if importfile.Banner != "" {
		article.Banner = GetFileName(importfile.Banner, "/")
	}

	if article.Banner != "" {
		article.Images = []string{article.Banner}
	}

	article.DataSource = newarticlepath + "/" + article.URL

	for _, cat := range importfile.Categories {
		newcat := strings.ToLower(cat)
		article.Categories = append(article.Categories, newcat)
	}

	for _, tag := range importfile.Tags {
		newtag := strings.ToLower(tag)
		article.Tags = append(article.Tags, newtag)
	}

	article.Content = importfile.Content

	articleTask.saveMarkdownFile(*article)

	if _, err := os.Stat(fileName); err == nil {
		fmt.Printf("Removing file: %s \n", fileName)
		err = os.Remove(fileName)
		if err != nil {
			msg := fmt.Errorf("Error Deleting import yaml file: %s \n ", err.Error())
			return nil, msg
		}
	}

	return article, err
}

//SaveArticle saves input data as an article and backups to a local md file
func (articleTask *Task) SaveArticle(article *Article, bypassquestions bool) (*Article, error) {
	if articleTask.service.Username == "" {
		articleTask.service.Username = AskForStringValue("Username", "", true)
	}

	if articleTask.service.Password == "" {
		articleTask.service.Password = AskForHiddenStringValue("Password", "", true)
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

	if bypassquestions == false {
		article.Banner = AskForStringValue("Banner Image FileName", article.Banner, false)
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
		_, err := articleTask.GetArticle(article.ID)

		if err != nil {
			requestMethod = "POST"
			article.ID = ""
		}
	}

	err := articleTask.service.SendRequest(requestMethod, requestURL, article)

	if err != nil {
		fmt.Printf("Unable to Save File, %s \n", err.Error())
		return article, err
	}

	imageEndPoint := fmt.Sprintf("images/%v", article.ID)
	datasourcePath := filepath.Dir(article.DataSource)

	for _, imageFilePath := range article.Images {
		imagepath := datasourcePath + "/" + imageFilePath
		strfile := strings.Split(imageFilePath, "/")
		filename := strfile[len(strfile)-1]

		imageLink := articleTask.service.ServiceURL + "/" + imageEndPoint + "/" + filename
		if !articleTask.service.ResolveLink(imageLink) {
			_, err := articleTask.service.Upload(imageEndPoint, imagepath)
			if err != nil {
				fmt.Printf("Could not save images %v, please try again. %v \n", imageLink, err.Error())
				continue
			}
		}
	}

	articleTask.saveMarkdownFile(*article)

	return article, err
}

//UpdateArticles updates and articles in a folder
func (articleTask *Task) UpdateArticles(filedir string, bypassQuestions bool) error {
	if filedir == "" {
		filedir = AskForStringValue("Import File or Folder", "", false)
	}

	subDirToSkip := []string{".git", ".DS_Store"}
	err := filepath.Walk(filedir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", filedir, err)
			return err
		}

		if info.IsDir() && contains(subDirToSkip, info.Name()) {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}

		if info.IsDir() == false {
			filename := info.Name()
			extension := filepath.Ext(filename)

			if extension == ".md" {
				fmt.Printf("updating file: %s \n ", path)

				_, err := articleTask.LoadArticle(path, bypassQuestions)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (articleTask *Task) saveMarkdownFile(article Article) error {
	filelocation := article.DataSource

	var importfile = &ImportArticle{
		article.ID,
		article.Title,
		article.URL,
		article.Images,
		article.Banner,
		article.PublishDate.Format("01/02/2006"),
		article.DataSource,
		article.Author,
		article.Categories,
		article.Tags,
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
