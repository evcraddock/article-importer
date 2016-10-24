package service

import (
	"bytes"
	"encoding/json"
	"encoding/base64"
	//"fmt"
	"net/http"
	//"net/http/httputil"
	"log"
	"github.com/evcraddock/article-importer/config"
)

type HttpService struct {
	ServiceUrl				string
	AuthKey					string
	Username				string
	Password				string
}

type AuthBody struct {
	access_token			string
}

type User struct {
	Id 				string 			`json:"id"`
	Name 			string 			`json:"name"`
	Picture 		string 			`json:"picture"`
	Email 			string 			`json:"email"`
}

type AuthUser struct {
	Token 			string 			`json:"token"`
	User 			User 			`json:"user"`
}

func NewHttpService(settings *config.Settings) * HttpService {
	svc := &HttpService{
		settings.ServiceUrl,
		settings.AuthKey,
		settings.UserName,
		settings.Password,
	}

	return svc
}

func (this *HttpService) GetJson(endpoint string, id string, target interface{}) error {
	url := this.ServiceUrl + "/" + endpoint + "/" + id

    r, err := http.Get(url)
    if err != nil {
        return err
    }

    defer r.Body.Close()
    return json.NewDecoder(r.Body).Decode(target)
}

func (this *HttpService) PostJson(endpoint string, target interface{}) error {
	currentUser, err := this.getUserToken()
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(target)
	if err !=nil {
		log.Fatal(err)
	}

	url := this.ServiceUrl + "/" + endpoint

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer " + currentUser.Token)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    res, err := client.Do(req)
    
    defer res.Body.Close()

    err = json.NewDecoder(res.Body).Decode(target)

    return err
}

func (this *HttpService) getUserToken() (*AuthUser, error) {
	log.Println("Getting User Token")

	authBody := &AuthBody{
		this.AuthKey,
	}

	b, err := json.Marshal(authBody)
	if err !=nil {
		log.Println("Error Marshaling")
		log.Fatal(err)
	}

	authstring := basicAuth(this.Username, this.Password)
	serviceUrl := this.ServiceUrl + "/auth?access_token=" + this.AuthKey 

	req, err := http.NewRequest("POST", serviceUrl, bytes.NewReader(b))
	req.Header.Set("Authorization", "Basic " + authstring)
    req.Header.Set("Content-Type", "application/json")

 //    dump, err := httputil.DumpRequestOut(req, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("%q", dump)

    client := &http.Client{}
    res, err := client.Do(req)

    if err != nil {
    	log.Fatal(err)
    }
    
    defer res.Body.Close()

 //    dumpr, err := httputil.DumpResponse(res, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("%q", dumpr)

    authUser := &AuthUser{}

    err = json.NewDecoder(res.Body).Decode(authUser)

    return authUser, err
}

func basicAuth(username, password string) string {
  auth := username + ":" + password
   return base64.StdEncoding.EncodeToString([]byte(auth))
}

