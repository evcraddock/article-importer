package tasks

import (
	"log"
)

//Link stores link information
type Link struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	LinkTitle  string   `json:"linkTitle"`
	URL        string   `json:"url"`
	Banner     string   `json:"banner"`
	Categories []string `json:"categories"`
	Tags       []string `json:"tags"`
}

func (linkTask *Task) saveLink(link *Link) (*Link, error) {
	if linkTask.service.Username == "" {
		linkTask.service.Username = AskForStringValue("Username", "", true)
	}

	if linkTask.service.Password == "" {
		linkTask.service.Password = AskForStringValue("Password", "", true)
	}

	if linkTask.service.ServiceURL == "" {
		linkTask.service.ServiceURL = AskForStringValue("Service Url", "", true)
	}

	if linkTask.service.AuthKey == "" {
		log.Fatal("AuthKey environment variable must be set.")
	}

	link.Title = AskForStringValue("Title", link.Title, true)
	link.LinkTitle = AskForStringValue("Link Title", link.LinkTitle, true)
	link.URL = AskForStringValue("Permalink", link.URL, true)
	link.Banner = AskForStringValue("Banner Url", link.Banner, false)
	link.Categories = AskForCSV("Categories (csv)", link.Categories)
	link.Tags = AskForCSV("Tags (csv)", link.Tags)

	requestMethod := "POST"
	requestURL := "links"

	err := linkTask.service.SendRequest(requestMethod, requestURL, link)

	return link, err
}

//CreateNewLink creates new Link
func (linkTask *Task) CreateNewLink() (*Link, error) {
	var link = &Link{
		Title:     "",
		LinkTitle: "",
		URL:       "",
		Banner:    "",
	}

	return linkTask.saveLink(link)
}

//DeleteLink deletes a specific link
func (linkTask *Task) DeleteLink() (string, error) {
	id := AskForStringValue("Link Id", "", true)
	if linkTask.service.Username == "" {
		linkTask.service.Username = AskForStringValue("Username", "", true)
	}

	if linkTask.service.Password == "" {
		linkTask.service.Password = AskForStringValue("Password", "", true)
	}

	if linkTask.service.ServiceURL == "" {
		linkTask.service.ServiceURL = AskForStringValue("Service Url", "", true)
	}

	if linkTask.service.AuthKey == "" {
		log.Fatal("AuthKey environment variable must be set.")
	}

	requestURL := "links/" + id

	return id, linkTask.service.SendRequest("DELETE", requestURL, nil)
}
