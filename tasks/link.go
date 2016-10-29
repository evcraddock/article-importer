package tasks

import (
	"log"
)

type Link struct {
	Id				string			`json:"id"`
	Title 			string			`json:"title"`
	LinkTitle		string			`json:"linkTitle"`
	Url 			string			`json:"url"`
	Banner			string 			`json:"banner"`
	Categories		[]string		`json:"categories"`
	Tags			[]string 		`json:"tags"`
}

func (this *Task) saveLink(link *Link) (*Link, error) {
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

	link.Title = AskForStringValue("Title", link.Title)
	link.LinkTitle = AskForStringValue("Link Title", link.LinkTitle)
	link.Url = AskForStringValue("Permalink", link.Url)
	link.Banner = AskForStringValue("Banner Url", link.Banner)
	link.Categories = AskForCsv("Categories (csv)", link.Categories)
	link.Tags = AskForCsv("Tags (csv)", link.Tags)
	
	requestMethod := "POST"
	requestUrl := "links"

	err := this.service.SendRequest(requestMethod, requestUrl, link)

	return link, err
}

func (this *Task) CreateNewLink() (*Link, error) {
	var link *Link = &Link{
		Title: "",
		LinkTitle: "",
		Url: "",
		Banner: "",
	}

	return this.saveLink(link)
}

func (this *Task) DeleteLink() (string, error) {
	id := AskForStringValue("Link Id", "")
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

	requestUrl := "links/" + id

	return id, this.service.SendRequest("DELETE", requestUrl, nil)
}
