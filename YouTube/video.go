package youtube

import (
	"fmt"
	"os"
	"path/filepath"

	ytdl "github.com/kkdai/youtube"
)

type Video struct {
	PublishedAt string
	Title       string
	Description string
	Position    int
	VideoID     string
	ChannelID   string
}

func NewVideo(
	PublishedAt, Title, Description, VideoID, ChannelID string,
	Position int,
) (Video, error) {
	return Video{
		PublishedAt: PublishedAt,
		Title:       Title,
		Description: Description,
		Position:    Position,
		VideoID:     VideoID,
		ChannelID:   ChannelID,
	}, nil
}

func (v Video) Download() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	currentDir := fmt.Sprintf("%v/Videos", exPath)

	y := ytdl.NewYoutube(false, false)
	err = y.DecodeURL(fmt.Sprintf("https://www.youtube.com/watch?v=%v", v.VideoID))
	if err != nil {
		return "", err
	}

	err = y.StartDownload(
		currentDir,
		fmt.Sprintf("%s-%s.mp4", v.ChannelID, v.VideoID),
		"medium",
		0,
	)

	return fmt.Sprintf("%s/%s", currentDir, fmt.Sprintf("%s-%s.mp4", v.ChannelID, v.VideoID)), err
}
