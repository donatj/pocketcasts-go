package pocketcasts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const endpointSignin = "https://play.pocketcasts.com/users/sign_in"

type Connection struct{}

type AuthedConnection struct {
	Client *http.Client
	*Connection
}

var ErrorInvalidUsernameOrPassword = errors.New("invalid username or password")

func (con *Connection) Authenticate(email, password string) (*AuthedConnection, error) {
	type authRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Scope    string `json:"scope"`
	}

	type authSuccess struct {
		Token string `json:"token"`
		UUID  string `json:"uuid"`
	}

	aReq := authRequest{
		Email:    email,
		Password: password,
		Scope:    "webplayer",
	}

	reqJSON, err := json.Marshal(aReq)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(reqJSON)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://api.pocketcasts.com/user/login", body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrorInvalidUsernameOrPassword
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error: %s", resp.Status)
	}

	respSuccess := authSuccess{}
	err = json.Unmarshal(respBody, &respSuccess)
	if err != nil {
		return nil, err
	}

	client.Transport = tokenTransport{
		Token: respSuccess.Token,
	}

	return &AuthedConnection{
		Client:     client,
		Connection: con,
	}, nil
}

type tokenTransport struct {
	Token string
}

func (t tokenTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("Authorization", "Bearer "+t.Token)

	if r.Header.Get("Content-Type") == "" {
		r.Header.Add("Content-Type", "application/json")
	}

	resp, err := http.DefaultTransport.RoundTrip(r)

	return resp, err
}
