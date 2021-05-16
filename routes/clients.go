package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	net "github.com/atulanand206/go-network"
	"github.com/atulanand206/ms-db-publisher/objects"
)

// Method to interact with the users client for getting the user from username.
func GetUser(username string, authorization string) (objects.User, error) {
	client := &http.Client{}
	prefix := os.Getenv("HOST_PREFIX")
	hostname := os.Getenv("USERS_URL")
	endpoint := "/users/username/"
	url := prefix + hostname + endpoint + username
	request, err := http.NewRequest(http.MethodGet, url, nil)
	var ob objects.User
	if err != nil {
		return ob, err
	}
	// Add the request headers to the request.
	request.Header.Set(net.Authorization, authorization)
	request.Header.Set(net.Accept, net.ApplicationJson)
	request.Header.Set(net.ContentTypeKey, net.ApplicationJson)
	// Trigger the network request.
	response, err := client.Do(request)
	if err != nil {
		return ob, err
	}
	decoder := json.NewDecoder(response.Body)
	// Decode the user from the response body.
	err = decoder.Decode(&ob)
	if err != nil {
		return ob, err
	}
	return ob, nil
}

// Method to interact with the users client for updating the user using userId.
func UpdateUser(userId string, user objects.UserRequest, authorization string) (objects.User, error) {
	client := &http.Client{}
	prefix := os.Getenv("HOST_PREFIX")
	hostname := os.Getenv("USERS_URL")
	endpoint := "/user/username/"
	url := prefix + hostname + endpoint + userId
	// Create the request body to update user.
	requestByte, _ := json.Marshal(user)
	requestReader := bytes.NewReader(requestByte)
	request, err := http.NewRequest(http.MethodPost, url, requestReader)
	var ob objects.User
	if err != nil {
		return ob, err
	}
	// Add the request headers to the request.
	request.Header.Set(net.Authorization, authorization)
	request.Header.Set(net.Accept, net.ApplicationJson)
	request.Header.Set(net.ContentTypeKey, net.ApplicationJson)
	// Trigger the network request.
	response, err := client.Do(request)
	if err != nil {
		return ob, err
	}
	fmt.Println(response)
	decoder := json.NewDecoder(response.Body)
	// Decode the user from the response body.
	err = decoder.Decode(&ob)
	if err != nil {
		return ob, err
	}
	return ob, nil
}
