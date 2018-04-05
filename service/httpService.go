package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"strings"

	"net/url"

	"github.com/evcraddock/article-importer/config"
)

//HTTPService information about an http service
type HTTPService struct {
	ServiceURL string
	AuthKey    string
	Username   string
	Password   string
}

//AuthBody authorization body information
type AuthBody struct {
	accessToken string
}

//User user object
type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

//AuthUser authorization user
type AuthUser struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

//NewHTTPService creates a new HTTPService
func NewHTTPService(settings config.Authorization) *HTTPService {
	svc := &HTTPService{
		settings.ServiceURL,
		settings.AuthKey,
		settings.UserName,
		settings.Password,
	}

	return svc
}

//Get returns a json payload
func (httpService *HTTPService) Get(endpoint string, id string, target interface{}) error {
	serviceURL := httpService.ServiceURL + "/" + endpoint + "/" + id

	r, err := http.Get(serviceURL)
	if err != nil {
		return err
	}

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

//ResolveLink checks the status of a link
func (httpService *HTTPService) ResolveLink(link string) bool {
	_, err := url.Parse(link)
	if err != nil {
		return false
	}

	r, err := http.Get(link)
	if err != nil {
		return false
	}

	return r.StatusCode == http.StatusOK
}

//Upload uploads and image to a service
func (httpService *HTTPService) Upload(endpoint, filename string) ([]byte, error) {
	url := httpService.ServiceURL + "/" + endpoint

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	bfile := make([]byte, 512)
	if err != nil {
		return nil, err
	}

	ctype := http.DetectContentType(bfile)

	if ctype == "application/octet-stream" {
		carr := strings.Split(filename, ".")
		ctype = "image/" + carr[len(carr)-1]
	}

	defer file.Close()

	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="%s"`, filename))
	header.Set("Content-Type", ctype)

	filewriter, err := writer.CreatePart(header)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(filewriter, file)
	if err != nil {
		return nil, err
	}

	writer.Close()
	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, err
	}

	contenttype := writer.FormDataContentType()
	req.Header.Set("Content-Type", contenttype)

	currentUser, err := httpService.getUserToken()
	if err != nil {
		log.Printf("Can't get user token: %s", err.Error())
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+currentUser.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		err = fmt.Errorf("Unable to save file: statusCode %s", res.Status)
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}

//SendRequest sends an http request
func (httpService *HTTPService) SendRequest(verb string, endpoint string, target interface{}) error {
	url := httpService.ServiceURL + "/" + endpoint

	currentUser, err := httpService.getUserToken()
	if err != nil {
		log.Fatal(err)
	}

	var req *http.Request

	if target != nil {
		b, err := json.Marshal(target)
		if err != nil {
			log.Fatal(err)
		}

		req, err = http.NewRequest(verb, url, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(verb, url, nil)
	}

	req.Header.Set("Authorization", "Bearer "+currentUser.Token)

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: Status Code - " + strconv.Itoa(res.StatusCode))
	}

	defer res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	if target != nil {
		err = json.NewDecoder(res.Body).Decode(target)
	}

	return err
}

func (httpService *HTTPService) getUserToken() (*AuthUser, error) {
	authstring := basicAuth(httpService.Username, httpService.Password)
	serviceURL := httpService.ServiceURL + "/auth?access_token=" + httpService.AuthKey

	req, err := http.NewRequest("POST", serviceURL, nil)
	req.Header.Set("Authorization", "Basic "+authstring)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 201 {
		err = errors.New("Unable to get user token: Status Code - " + strconv.Itoa(res.StatusCode))
		return nil, err
	}

	if err != nil {
		log.Fatal(err)
	}

	authUser := &AuthUser{}
	err = json.NewDecoder(res.Body).Decode(authUser)

	return authUser, err
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
