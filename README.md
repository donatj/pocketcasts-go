# pocketcasts-go

[![GoDoc](https://godoc.org/github.com/donatj/pocketcasts-go?status.svg)](https://godoc.org/github.com/donatj/pocketcasts-go)

Go connector for Pocket Casts Private API

## Caution

Pocket Casts does not have a public API.  This implements their **private API** thus there is zero guarantee this will not suddenly stop working.

Also as their API results may change, I can make no guarantee there will not be major changes to the interfaces.

I try to keep this up to date, but I make no promises.

I have no affiliation with Pocket Casts, I just like their product.

I have implmented the features I use, if you need something else, feel free to open an issue or PR. Features of the API are generally easy to implement.

## Usage

```go
package main

import (
	"log"

	pocketcasts "github.com/donatj/pocketcasts-go"
)

func main() {
	authedConn, err := pocketcasts.NewConnection().Authenticate("username", "password")
	if err != nil {
		log.Fatal(err)
	}

	// Get all subscribed podcasts
	podcasts, err := authedConn.GetSubscribedPodcasts()
	if err != nil {
		log.Fatal(err)
	}

	// Get all episodes for a podcast
	episodes, err := authedConn.GetPodcastEpisodes(podcasts.Podcasts[0].UUID)
	if err != nil {
		log.Fatal(err)
	}

	// Get all episode statuses for a podcast (playback position, archived, completion etc)
	statuses, err := authedConn.GetPodcastEpisodeStatuses(podcasts.Podcasts[0].UUID)
	if err != nil {
		log.Fatal(err)
	}
}
```
