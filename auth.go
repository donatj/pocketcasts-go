package pocketcasts

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const endpointSignin = "https://play.pocketcasts.com/users/sign_in"

type Connection struct {
}

type AuthedConnection struct {
	Client *http.Client
	*Connection
}

var ErrorInvalidUsernameOrPassword = errors.New("invalid username or password")

func (con *Connection) Authenticate(username, password string) (*AuthedConnection, error) {
	params := url.Values{}
	params.Set("[user]email", username)
	params.Set("[user]password", password)
	body := bytes.NewBufferString(params.Encode())

	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: jar,
	}

	// Create request
	req, err := http.NewRequest("POST", "https://play.pocketcasts.com/users/sign_in", body)
	if err != nil {
		return nil, err
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if strings.Contains(strings.ToLower(string(respBody)), "invalid email or password") {
		return nil, ErrorInvalidUsernameOrPassword
	}

	return &AuthedConnection{
		Client:     client,
		Connection: con,
	}, nil
}

type ConvertibleBoolean bool

func (bit ConvertibleBoolean) UnmarshalJSON(data []byte) error {
	asString := string(data)
	if asString == "1" || asString == "true" {
		bit = true
	} else if asString == "0" || asString == "false" {
		bit = false
	} else {
		return errors.New(fmt.Sprintf("Boolean unmarshal error: invalid input %s", asString))
	}
	return nil
}
