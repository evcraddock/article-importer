package tasks

type Image struct {
	Id          string `json:"_id"`
	FileName    string `json:"filename"`
	ContentType string `json:"contentType"`
}
