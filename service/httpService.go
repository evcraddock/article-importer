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
	"os"
	"strconv"

	"github.com/evcraddock/article-importer/config"
)

type HttpService struct {
	ServiceUrl string
	AuthKey    string
	Username   string
	Password   string
}

type AuthBody struct {
	access_token string
}

type User struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

type AuthUser struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func NewHttpService(settings config.Authorization) *HttpService {
	svc := &HttpService{
		settings.ServiceUrl,
		settings.AuthKey,
		settings.UserName,
		settings.Password,
	}

	return svc
}

func (httpService *HttpService) GetJson(endpoint string, id string, target interface{}) error {
	url := httpService.ServiceUrl + "/" + endpoint + "/" + id

	r, err := http.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func (httpService *HttpService) Upload(endpoint, filename string) ([]byte, error) {
	url := httpService.ServiceUrl + "/" + endpoint

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	filewriter, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(filewriter, file)
	if err != nil {
		return nil, err
	}

	filewriter, err = writer.CreateFormField("key")
	if err != nil {
		return nil, err
	}

	_, err = filewriter.Write([]byte("KEY"))
	if err != nil {
		return nil, err
	}

	writer.Close()
	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

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

func (httpService *HttpService) SendRequest(verb string, endpoint string, target interface{}) error {
	url := httpService.ServiceUrl + "/" + endpoint

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

	if err != nil {
		fmt.Printf("Error sending request: %s", err.Error())
	}

	return err
}

func (httpService *HttpService) getUserToken() (*AuthUser, error) {
	authstring := basicAuth(httpService.Username, httpService.Password)
	serviceUrl := httpService.ServiceUrl + "/auth?access_token=" + httpService.AuthKey

	req, err := http.NewRequest("POST", serviceUrl, nil)
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
