package pocketcasts

import (
	"encoding/json"
	"net/http"
)

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

type PodcastUUID string

type SubscribedPodcasts struct {
	Podcasts []struct {
		ID                int         `json:"id"`
		UUID              PodcastUUID `json:"uuid"`
		URL               string      `json:"url"`
		Title             string      `json:"title"`
		Description       string      `json:"description"`
		ThumbnailURL      string      `json:"thumbnail_url"`
		Author            string      `json:"author"`
		EpisodesSortOrder int         `json:"episodes_sort_order"`
	} `json:"podcasts"`

	App struct {
		UserVersionCode int     `json:"userVersionCode"`
		VersionCode     int     `json:"versionCode"`
		VersionName     float64 `json:"versionName"`
		VersionSummary  string  `json:"versionSummary"`
	} `json:"app"`
}
