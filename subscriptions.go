package pocketcasts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GetSubscribedPodcasts returns a list of podcasts the user is subscribed to
// as well as the users folders.
func (acon *AuthedConnection) GetSubscribedPodcasts() (*SubscribedPodcasts, error) {
	body := strings.NewReader(`{"v":1}`)

	req, err := http.NewRequest("POST", "https://api.pocketcasts.com/user/podcast/list", body)
	if err != nil {
		return nil, err
	}

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

type FolderUUID string

type SubscribedPodcasts struct {
	Podcasts []struct {
		UUID                     PodcastUUID        `json:"uuid"`
		EpisodesSortOrder        int                `json:"episodesSortOrder"`
		AutoStartFrom            int                `json:"autoStartFrom"`
		Title                    string             `json:"title"`
		Author                   string             `json:"author"`
		Description              string             `json:"description"`
		URL                      string             `json:"url"`
		LastEpisodePublished     time.Time          `json:"lastEpisodePublished"`
		Unplayed                 bool               `json:"unplayed"`
		LastEpisodeUUID          PodcastEpisodeUUID `json:"lastEpisodeUuid"`
		LastEpisodePlayingStatus int                `json:"lastEpisodePlayingStatus"`
		LastEpisodeArchived      bool               `json:"lastEpisodeArchived"`
		AutoSkipLast             int                `json:"autoSkipLast"`
		FolderUUID               FolderUUID         `json:"folderUuid"`
		SortPosition             int                `json:"sortPosition"`
		DateAdded                time.Time          `json:"dateAdded"`
	} `json:"podcasts"`
	Folders []struct {
		FolderUUID       FolderUUID `json:"folderUuid"`
		Name             string     `json:"name"`
		Color            int        `json:"color"`
		SortPosition     int        `json:"sortPosition"`
		PodcastsSortType int        `json:"podcastsSortType"`
		DateAdded        time.Time  `json:"dateAdded"`
	} `json:"folders"`
}
