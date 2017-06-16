package tasks

//Image stores image data
type Image struct {
	ID          string `json:"_id"`
	FileName    string `json:"filename"`
	ContentType string `json:"contentType"`
}
