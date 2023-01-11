package input

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestParseVideoLinks(t *testing.T) {
	//conf := &GrabberConf{
	//	PlaylistURL: "",
	//	UpdateEvery: 10 * time.Second,
	//}
	//grabber := NewGrabber(conf, &log.Logger{})

	mockPlaylistContent, err := os.ReadFile("./testdata/index.m3u8")

	if err != nil {
		t.Fatal(err)
	}

	links := extractVideoLinks(string(mockPlaylistContent))

	log.Println(links)

	if len(links) != 3 {
		t.Errorf("assummed 3 links, got %d", len(links))
	}

}

const mockURL = "https://example.com"

func mockLoadPlaylistContent() (string, error) {
	return mockURL, nil
}

func mockLoadPlaylistContentWithErr() (string, error) {
	return "", fmt.Errorf("test error")
}

func TestGrabber_GrabURLS(t *testing.T) {
	conf := &GrabberConf{
		PlaylistURL: "",
		UpdateEvery: 1,
	}
	grabber := NewGrabber(conf, &log.Logger{})

	tests := []struct {
		name   string
		f      UrlLoadingFunc
		want   UrlsChan
		values []string
	}{
		{
			name:   "success load",
			f:      mockLoadPlaylistContent,
			values: []string{mockURL},
		},
		{
			name:   "error on load",
			f:      mockLoadPlaylistContentWithErr,
			values: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			ch := grabber.GrabURLS(ctx, test.f)

			for _, v := range test.values {
				chValue := <-ch
				if v != chValue {
					t.Errorf("channels content not equal %v != %v", v, chValue)
				}
			}

			cancel()
		})

	}

}
