package pocketcasts

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (acon *AuthedConnection) GetSubscribedPodcasts() (*SubscribedPodcasts, error) {
	req, err := http.NewRequest("POST", "https://play.pocketcasts.com/web/podcasts/all.json", nil)

	// Fetch Request
	resp, err := acon.Client.Do(req)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)

	out := &SubscribedPodcasts{}

	err = dec.Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type SubscribedPodcasts struct {
	Podcasts []struct {
		ID                int    `json:"id"`
		UUID              string `json:"uuid"`
		URL               string `json:"url"`
		Title             string `json:"title"`
		Description       string `json:"description"`
		ThumbnailURL      string `json:"thumbnail_url"`
		Author            string `json:"author"`
		EpisodesSortOrder int    `json:"episodes_sort_order"`
	} `json:"podcasts"`
	App struct {
		UserVersionCode int     `json:"userVersionCode"`
		VersionCode     int     `json:"versionCode"`
		VersionName     float64 `json:"versionName"`
		VersionSummary  string  `json:"versionSummary"`
	} `json:"app"`
}
