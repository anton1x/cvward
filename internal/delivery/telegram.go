package delivery

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"os"
)

type Telegram struct {
	Api              *tgbotapi.BotAPI
	VideoSubscribers []*tgbotapi.Chat
}

func NewTelegram() (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI("1762087955:AAHZtUQRWuXYS_XiLOZFCtZG-Y_u9oXbKuQ")
	if err != nil {
		return nil, err
	}

	return &Telegram{Api: bot}, nil
}

func (t *Telegram) HandleSubscribers() error {
	ch, err := t.Api.GetUpdatesChan(tgbotapi.NewUpdate(0))
	if err != nil {
		return err
	}
	go func() {
		for upd := range ch {
			t.VideoSubscribers = append(t.VideoSubscribers, upd.Message.Chat)
			log.Println("appended sub", upd.Message.Chat.UserName)
		}
	}()

	return nil
}

func (t *Telegram) SendToSubs(ch chan string) {
	go func(ch chan string) {
		for fpath := range ch {
			f, err := ioutil.ReadFile(fpath)
			if err != nil {
				continue
			}

			tgFile := tgbotapi.FileBytes{
				Name:  fpath,
				Bytes: f,
			}

			for _, sub := range t.VideoSubscribers {
				vupload := tgbotapi.NewVideoUpload(sub.ID, tgFile)
				_, err := t.Api.Send(vupload)
				if err != nil {
					log.Println("error sending fpath", fpath, err)
				}
			}
			//f.Close()
			_ = os.RemoveAll(fpath)

		}
	}(ch)

}
