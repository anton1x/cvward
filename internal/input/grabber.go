package input

import (
	"context"
	"io"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"time"
)

type Grabber struct {
	conf        *GrabberConf
	logger      *log.Logger
	lastGrabbed []string
}

type UrlsChan chan string

type GrabberConf struct {
	PlaylistURL string
	UpdateEvery time.Duration
}

func NewGrabber(conf *GrabberConf, l *log.Logger) *Grabber {
	return &Grabber{
		conf:   conf,
		logger: l,
	}
}

func extractVideoLinks(playlistContent string) []string {

	rexp := regexp.MustCompile(`https?://[^\s]*`)
	res := rexp.FindAllString(playlistContent, 100)

	return res
}

func (g *Grabber) loadPlaylistContent() (string, error) {
	resp, err := http.Get(g.conf.PlaylistURL)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(res), err
}

func (g *Grabber) GrabURLS(ctx context.Context) UrlsChan {
	ch := make(UrlsChan)
	ticker := time.Tick(g.conf.UpdateEvery)
	go func() {
		defer close(ch)
		for loop := true; loop; {
			select {
			case <-ticker:
				content, err := g.loadPlaylistContent()
				if err != nil {
					continue
				}
				links := extractVideoLinks(content)

				if !reflect.DeepEqual(g.lastGrabbed, links) {
					continue
				}
				g.lastGrabbed = links

				for _, link := range links {
					ch <- link
				}
			case <-ctx.Done():
				loop = false
			}

		}
	}()

	return ch

}
