package youtube

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/google/go-cmp/cmp"
)

func NewMockConfiguration(t *testing.T) *ucfg.Config {
	// GET api_key from os env
	key, ok := os.LookupEnv("API_KEY")
	if !ok {
		t.Fatal("the API_KEY environment variable must be set to access youtube data api v3 before testing")
	}
	cfg, err := yaml.NewConfig([]byte(fmt.Sprintf("api_key: %v", key)), ucfg.PathSep("."))
	if err != nil {
		t.Fatalf("error creating a new configuration: %v", err.Error())
	}
	return cfg
}

func TestNewChannel(t *testing.T) {
	// First create a new service
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := NewMockConfiguration(t)
	err := NewService(ctx, cfg)
	if err != nil {
		t.Fatalf("error creating a new youtube service: %v", err.Error())
	}

	testCases := []struct {
		testName        string
		channelURL      string
		expectedChannel *Channel
	}{
		{
			"goodURL",
			"https://www.youtube.com/channel/UCG0SzK_t4-Ylf1yZq9Xmi_g",
			&Channel{
				Title:            "Freeme NCS Music",
				URL:              "https://www.youtube.com/channel/UCG0SzK_t4-Ylf1yZq9Xmi_g",
				ID:               "UCG0SzK_t4-Ylf1yZq9Xmi_g",
				UploadPlaylistID: "UUG0SzK_t4-Ylf1yZq9Xmi_g",
			},
		},
		// Channel with bad url
		{
			"badURL",
			"https://www.youtube.com/channel/UCG0SzK_t4-Ylf19Xmi_g",
			&Channel{
				Title:            "Freeme NCS Music",
				URL:              "https://www.youtube.com/channel/UCG0SzK_t4-Ylf1yZq9Xmi_g",
				ID:               "UCG0SzK_t4-Ylf1yZq9Xmi_g",
				UploadPlaylistID: "UUG0SzK_t4-Ylf1yZq9Xmi_g",
			},
		},
	}

	test := testCases[0]
	t.Run(test.testName, func(t *testing.T) {
		channel, err := NewChannel(test.channelURL)
		if err != nil {
			t.Fatalf("error creating a channel from a valid url: %v\n", err)
		}

		if !cmp.Equal(test.expectedChannel, channel) {
			t.Fatalf("channel %+v and testChannel %+v don't contain the same underlaying values", channel, test.expectedChannel)
		}
	})

	test = testCases[1]
	t.Run(test.testName, func(t *testing.T) {
		_, err := NewChannel(test.channelURL)
		if err == nil {
			t.Fatalf("created a channel from an invalid url: %v\n", err)
		}
	})
}

func TestPlaylist(t *testing.T) {
	// First create a new service
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := NewMockConfiguration(t)
	err := NewService(ctx, cfg)
	if err != nil {
		t.Fatalf("error creating a new youtube service: %v", err.Error())
	}

	testCases := []struct {
		testName       string
		channel        *Channel
		expectedVideos []string
	}{
		{
			"chanelWithGoodPlaylistID",
			&Channel{
				UploadPlaylistID: "UUdgUTNVvHrcJq7K9nyEQ6qg",
			},
			[]string{
				"Legendary washing machine - washing of my black socks collection",
				"Amica AWN510D vaskemaskine",
				"RELAXING NOISE COLLECTION - a wind turbine noise - COMMERCIALS ONLY BEFORE VIDEO!",
				"Amica Waschmaschine tägliche Wäsche daglig tvätt daglig vask vaskemaskine tvättmaskin",
				"100% REAL SEA sound for REALXING and CALMING sea noise on the beach",
				"ASMR washing machine NIGHT SHOT stylish laundry",
				"Washing machine INTERIOR during washing 1600 rpm and Haineken beer :)",
				"Drum Lamp and Synthetics fabrics cycle washing machine program night shot",
				"Manually adding water to washing machine",
				"RELAXING CYCLE - soaking",
				"Intensive wash & preliminary washing",
				"Rinsing and centrifuging of Samsung washing machine",
				"Samsung s9 slow motion washing machine",
				"Washing machine spin 2000 rpm relaxing sound",
				"Cotton 60 cycle white program washing machine 2019",
				"New Year's Eve for nerds and introvertics! 😀 Manual wash program Gorenje wa 65205",
				"Christmas cycle - washing machine spin",
				"Like a BOSS - GORENJE WA 65205 - delicate program",
				"Baby care cycle Samsung washing machine white noise sleep",
				"Washing machine inside aliexpress endocsope camera 720p",
				"Outdoor Care Washing machine Samsung Jacket wash",
				"Recomendation for washing machine fans!",
				"Christmas Spin Cycle - The shortest program in samsung washing machine",
				"Super Eco Wash cycle washer samsung washing machine",
				"White noise generator Farel electric blower heater",
				"Washing machine spin cycle 12H 💪🔥",
				"Washing at an accelerated pace Compression from 1 hour to 3 minutes Jeans program",
				"Washing jeans in washer 👖",
				"Sleep sounds Air filter sound",
				"Daily wash washing machine in samsung washer",
				"Newborn baby sleep █■█ █ ▀█▀ Hair dryer sound - stereo best quality 👍",
				"Working white noise (fresh and tasty) - hair dryer",
				"CLICK ME =) white noise",
				"Remington AC5011 Hair dryer - Relax put a baby to sleep - 2H white noise",
				"THE LONGEST program for COTTON washer. Samsung washing machine",
				"White noise in home - relaxing sound",
				"Grass mowing Lawn mower White noise Tondeuse à gazon çim biçme makinesi",
				"Sound of rain / Geräusch von Regen / le son de la pluie / yağmur sesi / medicamento para dormir",
				"Cooker hood Mastercook 727 60 -  white noise machine :)",
				"Hand wash program Samsung wasching machine cuci tangan",
				"Samsung washing machine eco drum clean program - wash wasching machine demo",
				"XpressJa freestyle NOKIA - snow",
				"Takam vs Rekowski - last round  (KO) / ostatnia runda",
				"Mofeta - Tylicz",
				"Gocław 24h - dzień z życia osiedla",
				"Witold Pilecki, Kazimierz Piechowski, Eduard Schulte - film dokumentalny",
				"Samsung Eco Bubble WF60F4E4W2W - program for synthetic fabrics",
				"ICBM - Civilization III - atom bomb",
				"Samsung washing machine melody",
				"Additional rinsing - washing machine - two cameras",
				"Samsung WF60F4E4W2W ecobubble washing machine/pralka/vaskemaskine/rentadora",
			},
		},
	}

	test := testCases[0]
	t.Run(test.testName, func(t *testing.T) {
		playlist, err := NewPlaylist(*test.channel)
		if err != nil {
			t.Fatalf("error creating a new playlist: %v", err.Error())
		}

		// Iterate over all videos
		videos := []string{}
		for playlist.HasNextPage() {
			pageVideos, err := playlist.GetNextVideos()
			if err != nil {
				t.Fatal(err)
			}
			switch {
			case len(pageVideos) > 0:
				for _, video := range pageVideos {
					videos = append(videos, video.Title)
				}
			default:
				break
			}
		}

		if !cmp.Equal(videos, test.expectedVideos) {
			t.Fatal("video are not the same as expected in the playlist")
		}
	})
}

func TestDownload(t *testing.T) {
	// First create a new service
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := NewMockConfiguration(t)
	err := NewService(ctx, cfg)
	if err != nil {
		t.Fatalf("error creating a new youtube service: %v", err.Error())
	}

	testVideo, err := NewVideo(
		"2016-01-20T19:42:25Z",
		"ICBM - Civilization III - atom bomb",
		"ICBM - Civilization III / Civilization 3 - ICBM",
		"ZgyEGAukuE0",
		"UCdgUTNVvHrcJq7K9nyEQ6qg",
		47,
	)
	if err != nil {
		t.Fatalf("error creating a new video: %v", err.Error())
	}

	_, err = testVideo.Download()
	if err != nil {
		t.Fatalf("error downloading video: %v", err.Error())
	}
}
