package pocketcasts

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (acon *AuthedConnection) EpisodesFindByPodcast(UUID string) (*FindByPodcast, error) {
	data := FindByPodcastRequest{
		UUID: UUID,
		Page: 1,
		Sort: 3,
	}

	bodyData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(bodyData)

	req, err := http.NewRequest("POST", "https://play.pocketcasts.com/web/episodes/find_by_podcast.json", body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// Fetch Request
	resp, err := acon.Client.Do(req)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)

	out := &FindByPodcast{}

	err = dec.Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type FindByPodcastRequest struct {
	UUID string `json:"uuid"`
	Page int    `json:"page"`
	Sort int    `json:"sort"`
}

type FindByPodcast struct {
	Status    string `json:"status"`
	Token     string `json:"token"`
	Copyright string `json:"copyright"`
	Result    struct {
		Episodes []struct {
			ID            interface{} `json:"id"`
			UUID          string      `json:"uuid"`
			URL           string      `json:"url"`
			Title         string      `json:"title"`
			PublishedAt   string      `json:"published_at"`
			Duration      string      `json:"duration"`
			FileType      string      `json:"file_type"`
			Size          int         `json:"size"`
			PlayingStatus int         `json:"playing_status"`
			PlayedUpTo    int         `json:"played_up_to"`
			IsDeleted     PKBoolean   `json:"is_deleted"`
			Starred       PKBoolean   `json:"starred"`
			IsVideo       PKBoolean   `json:"is_video"`
		} `json:"episodes"`
		Total int `json:"total"`
	} `json:"result"`
}
