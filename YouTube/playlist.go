package youtube

import (
	"google.golang.org/api/googleapi"
)

// Playlist is used to identify uploaded videos for a channel
type Playlist struct {
	ID            string
	pageNextToken string
	hasNextPage   bool
}

// NewPlaylist receives a Channel struct and creates a new Playlist instance
// based on its parameter
func NewPlaylist(channel Channel) (*Playlist, error) {
	return &Playlist{
		ID:            channel.UploadPlaylistID,
		pageNextToken: "",
		hasNextPage:   true,
	}, nil
}

// HasNextPage indicate if a playlist has more videos or not
func (p *Playlist) HasNextPage() bool {
	return p.hasNextPage
}

// GetNextVideos obtains the list of videos from the current page
func (p *Playlist) GetNextVideos() ([]Video, error) {
	if !p.hasNextPage {
		return []Video{}, nil
	}

	call := svc.PlaylistItems.List([]string{"snippet"}).
		PlaylistId(p.ID).
		MaxResults(50).
		Fields(
			googleapi.Field("nextPageToken"),
			googleapi.Field("items/snippet(publishedAt,title,channelId,description,position,resourceId)"),
		)
	if p.pageNextToken != "" {
		call = call.PageToken(p.pageNextToken)
	}
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	videos := []Video{}
	for _, item := range resp.Items {
		video, err := NewVideo(
			item.Snippet.PublishedAt,
			item.Snippet.Title,
			item.Snippet.Description,
			item.Snippet.ResourceId.VideoId,
			item.Snippet.ChannelId,
			int(item.Snippet.Position),
		)
		if err != nil {
			return nil, err
		}

		videos = append(videos, video)
	}
	p.pageNextToken = resp.NextPageToken
	if resp.NextPageToken == "" {
		p.hasNextPage = false
	}
	return videos, nil
}
