package pocketcasts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (acon *AuthedConnection) GetSubscribedPodcasts() (*SubscribedPodcasts, error) {
	body := strings.NewReader(`{"v":1}`)

	req, err := http.NewRequest("POST", "https://api.pocketcasts.com/user/podcast/list", body)

	// Fetch Request
	resp, err := acon.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)

	out := &SubscribedPodcasts{}

	err = dec.Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type PodcastUUID string

type SubscribedPodcasts struct {
	Podcasts []struct {
		UUID                     PodcastUUID `json:"uuid"`
		EpisodesSortOrder        int         `json:"episodesSortOrder"`
		AutoStartFrom            int         `json:"autoStartFrom"`
		Title                    string      `json:"title"`
		Author                   string      `json:"author"`
		Description              string      `json:"description"`
		URL                      string      `json:"url"`
		LastEpisodePublished     time.Time   `json:"lastEpisodePublished"`
		Unplayed                 bool        `json:"unplayed"`
		LastEpisodeUUID          string      `json:"lastEpisodeUuid"`
		LastEpisodePlayingStatus int         `json:"lastEpisodePlayingStatus"`
		LastEpisodeArchived      bool        `json:"lastEpisodeArchived"`
	} `json:"podcasts"`
}
