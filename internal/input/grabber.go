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

type UrlLoadingFunc func() (string, error)

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

func (g *Grabber) LoadPlaylistContent() (string, error) {
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

func (g *Grabber) GrabURLS(ctx context.Context, loadingFunc UrlLoadingFunc) UrlsChan {
	ch := make(UrlsChan)
	ticker := time.Tick(g.conf.UpdateEvery)
	go func() {
		defer close(ch)
		for loop := true; loop; {
			select {
			case <-ctx.Done():
				loop = false
			case <-ticker:
				log.Println("Starting grab URLS")
				content, err := loadingFunc()
				if err != nil {
					log.Println(err)
					continue
				}
				links := extractVideoLinks(content)

				log.Printf("Extracted %d links", len(links))

				if reflect.DeepEqual(g.lastGrabbed, links) {
					log.Printf("Continue")
					continue
				}
				g.lastGrabbed = links

				for _, link := range links {
					ch <- link
				}

			}

		}
	}()

	return ch

}
