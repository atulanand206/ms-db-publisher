package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/atulanand206/ms-db-publisher/objects"
)

const (
	contentTypeKey             = "content-type"
	contentTypeApplicationJson = "application/json"
	acceptKey                  = "Accept"
	applicationJson            = "application/json"
)

func GetUser(username string) (objects.User, error) {
	client := &http.Client{}
	hostname := os.Getenv("USERS_URL")
	endpoint := "/users/username/"
	url := "http://" + hostname + endpoint + username
	request, err := http.NewRequest("GET", url, nil)
	var ob objects.User
	if err != nil {
		return ob, err
	}
	request.Header.Add(acceptKey, applicationJson)
	request.Header.Add(contentTypeKey, applicationJson)
	response, err := client.Do(request)
	if err != nil {
		return ob, err
	}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&ob)
	if err != nil {
		return ob, err
	}
	return ob, nil
}

func UpdateUser(userId string, user objects.UserRequest) (objects.User, error) {
	client := &http.Client{}
	hostname := os.Getenv("USERS_URL")
	endpoint := "/user/username/"
	url := "http://" + hostname + endpoint + userId
	requestByte, _ := json.Marshal(user)
	requestReader := bytes.NewReader(requestByte)
	request, err := http.NewRequest("POST", url, requestReader)
	var ob objects.User
	if err != nil {
		return ob, err
	}
	request.Header.Add(acceptKey, applicationJson)
	request.Header.Add(contentTypeKey, applicationJson)
	response, err := client.Do(request)
	if err != nil {
		return ob, err
	}
	fmt.Println(response)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&ob)
	if err != nil {
		return ob, err
	}
	return ob, nil
}
