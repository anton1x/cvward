package input

import (
	"context"
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

func TestGrabber_GrabURLS(t *testing.T) {
	conf := &GrabberConf{
		PlaylistURL: "",
		UpdateEvery: 1,
	}
	grabber := NewGrabber(conf, &log.Logger{})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	ch := grabber.GrabURLS(ctx, mockLoadPlaylistContent)

	customCh := make(chan string)
	defer close(customCh)
	go func() {
		customCh <- mockURL
	}()

	a, b := <-ch, <-customCh

	if a != b {
		t.Errorf("channels content not equal %v != %v", a, b)
	}

	if grabber.lastGrabbed[len(grabber.lastGrabbed)-1] != mockURL {
		t.Errorf("last grabbed not stored")
	}

	cancel()

}
