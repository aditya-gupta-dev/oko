package api

import (
	"context"

	"github.com/aditya-gupta-dev/oko/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeSong struct {
	Title   string
	VideoId string
}

func SearchSongYoutube(query string, maxResults int64) ([]YoutubeSong, error) {
	var apikey string = config.GetYoutubeApiKey()
	ctx := context.Background()

	service, err := youtube.NewService(ctx, option.WithAPIKey(apikey))
	if err != nil {
		return nil, err
	}

	call := service.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(maxResults)

	var songs []YoutubeSong = make([]YoutubeSong, maxResults)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	for _, item := range response.Items {
		song := YoutubeSong{
			Title:   item.Snippet.Title,
			VideoId: item.Id.VideoId,
		}
		songs = append(songs, song)
	}

	return songs, nil
}
