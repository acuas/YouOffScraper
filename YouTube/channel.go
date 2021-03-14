package youtube

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"google.golang.org/api/googleapi"
)

// Channel contains information about a YouTube channel
type Channel struct {
	Title            string
	URL              string
	ID               string
	UploadPlaylistID string
}

// NewChannel receive an url and based on it, creates a new Channel instance.
// In case of failure, an error is returned in second parameter, otherwise
// this parameter is nil
func NewChannel(url string) (*Channel, error) {
	// Get the ID of the channel
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// TODO: make a reader that doesn't read the entire content of the request
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("\"browseEndpoint\":{\"browseId\":\"(.+?)\"")
	matches := re.FindStringSubmatch(string(body))
	if len(matches) != 2 {
		return nil, fmt.Errorf("cannot extract the id from url: %v", url)
	}
	id := matches[1]

	// Get info about the channel
	call := svc.Channels.List([]string{"snippet", "contentDetails"}).
		Id(id).
		Fields(googleapi.Field("items(snippet(title),contentDetails/relatedPlaylists/uploads)"))
	ytResponse, err := call.Do()
	if err != nil {
		return nil, err
	}
	if len(ytResponse.Items) != 1 {
		return nil, fmt.Errorf("number of items returned by the API is different than 1, so maybe the channel doesn't exists")
	}

	return &Channel{
		Title:            ytResponse.Items[0].Snippet.Title,
		UploadPlaylistID: ytResponse.Items[0].ContentDetails.RelatedPlaylists.Uploads,
		URL:              url,
		ID:               id,
	}, nil
}
