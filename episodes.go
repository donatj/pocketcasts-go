package pocketcasts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (acon *AuthedConnection) GetPodcastEpisodes(UUID PodcastUUID) (*PodcastEpisodes, error) {
	req, err := http.NewRequest("GET", "https://cache.pocketcasts.com/podcast/full/"+string(UUID)+"/0/2/1000", nil)

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

	out := &PodcastEpisodes{}

	err = dec.Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type PodcastEpisodeUUID string

type PodcastEpisodes struct {
	EpisodeFrequency       string    `json:"episode_frequency"`
	EstimatedNextEpisodeAt time.Time `json:"estimated_next_episode_at"`
	HasSeasons             bool      `json:"has_seasons"`
	SeasonCount            int       `json:"season_count"`
	EpisodeCount           int       `json:"episode_count"`
	HasMoreEpisodes        bool      `json:"has_more_episodes"`
	Podcast                struct {
		URL         string `json:"url"`
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Audio       bool   `json:"audio"`
		ShowType    string `json:"show_type"`

		UUID     PodcastUUID `json:"uuid"`
		Episodes []struct {
			Title     string    `json:"title"`
			URL       string    `json:"url"`
			FileType  string    `json:"file_type"`
			FileSize  int       `json:"file_size"`
			Duration  int       `json:"duration"`
			Published time.Time `json:"published"`
			Type      string    `json:"type"`
			Season    int       `json:"season,omitempty"`
			Number    int       `json:"number,omitempty"`

			UUID PodcastEpisodeUUID `json:"uuid"`
		} `json:"episodes"`
	} `json:"podcast"`
}

func (acon *AuthedConnection) GetPodcastEpisodeStatuses(UUID PodcastUUID) (*PodcastEpisodeStatuses, error) {
	type reqPodcastEpisodeStatuses struct {
		UUID PodcastUUID `json:"uuid"`
	}

	sReq := reqPodcastEpisodeStatuses{
		UUID: UUID,
	}

	reqJSON, err := json.Marshal(sReq)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(reqJSON)

	req, err := http.NewRequest("POST", "https://api.pocketcasts.com/user/podcast/episodes", body)

	// Fetch Request
	resp, err := acon.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)

	out := &PodcastEpisodeStatuses{}

	err = dec.Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type EpisodePlayingStatus int

const (
	StatusUnplayed EpisodePlayingStatus = 1
	StatusStarted  EpisodePlayingStatus = 2
	StatusFinished EpisodePlayingStatus = 3
)

type PodcastEpisodeStatuses struct {
	Episodes []struct {
		UUID PodcastEpisodeUUID `json:"uuid"`

		PlayingStatus EpisodePlayingStatus `json:"playingStatus"`

		PlayedUpTo int  `json:"playedUpTo"`
		IsDeleted  bool `json:"isDeleted"`
		Starred    bool `json:"starred"`
		Duration   int  `json:"duration"`
	} `json:"episodes"`
}

func (acon *AuthedConnection) UpdateEpisodeStatus(episodeUUID PodcastEpisodeUUID, podcastUUID PodcastUUID, status EpisodePlayingStatus) error {
	type reqStatusUpdate struct {
		UUID    PodcastEpisodeUUID   `json:"uuid"`
		Podcast PodcastUUID          `json:"podcast"`
		Status  EpisodePlayingStatus `json:"status"`
	}

	statusReq := reqStatusUpdate{
		UUID:    episodeUUID,
		Podcast: podcastUUID,
		Status:  status,
	}

	reqJSON, err := json.Marshal(statusReq)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(reqJSON)

	req, err := http.NewRequest("POST", "https://api.pocketcasts.com/sync/update_episode", body)

	// Fetch Request
	resp, err := acon.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request error: %s", resp.Status)
	}

	return nil
}

func (acon *AuthedConnection) UpdateEpisodeArchive(episodeUUID PodcastEpisodeUUID, podcastUUID PodcastUUID, archive bool) error {
	type reqArchiveUpdatePodcast struct {
		UUID    PodcastEpisodeUUID `json:"uuid"`
		Podcast PodcastUUID        `json:"podcast"`
	}
	type reqArchiveUpdate struct {
		Episodes []reqArchiveUpdatePodcast `json:"episodes"`
		Archive  bool                      `json:"archive"`
	}

	archiveReq := reqArchiveUpdate{
		Episodes: []reqArchiveUpdatePodcast{
			{episodeUUID, podcastUUID},
		},
		Archive: archive,
	}

	reqJSON, err := json.Marshal(archiveReq)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(reqJSON)

	req, err := http.NewRequest("POST", "https://api.pocketcasts.com/sync/update_episodes_archive", body)

	// Fetch Request
	resp, err := acon.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request error: %s", resp.Status)
	}

	return nil
}
