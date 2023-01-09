package input

import (
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
