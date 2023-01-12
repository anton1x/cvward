package input

import (
	"blob/internal/pkg/helper"
	"context"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

type Grabber struct {
	conf        *GrabberConf
	logger      *log.Logger
	lastGrabbed []string
	ticker      *time.Ticker
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
		ticker: time.NewTicker(conf.UpdateEvery),
	}
}

func extractVideoLinks(playlistContent string) []string {

	rexp := regexp.MustCompile(`https?://[^\s]*`)
	res := rexp.FindAllString(playlistContent, 100)

	return res
}

func (g *Grabber) LoadPlaylistContent() (string, error) {
	resp, err := http.Get(g.conf.PlaylistURL)

	if err != nil || resp == nil {
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
	go func() {
		defer close(ch)
		defer g.ticker.Stop()
		for loop := true; loop; {
			select {
			case <-ctx.Done():
				loop = false
			case <-g.ticker.C:
				log.Println("Starting grab URLS")
				content, err := loadingFunc()
				if err != nil {
					log.Println(err)
					continue
				}
				links := extractVideoLinks(content)

				log.Printf("Extracted %d links", len(links))
				newLinks := helper.SliceDiff(g.lastGrabbed, links)
				if len(g.lastGrabbed) > 0 && len(newLinks) == 0 {
					log.Printf("Continue")
					continue
				}
				g.lastGrabbed = links

				for _, link := range newLinks {
					ch <- link
				}

			}

		}
	}()

	return ch

}
