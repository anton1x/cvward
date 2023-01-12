package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"os"
)

type Handler struct {
	Command     string
	Description string
	HandleFunc  HandleFunc
}

type HandleFunc func(*Telegram, tgbotapi.Update)

type Telegram struct {
	Api              *tgbotapi.BotAPI
	VideoSubscribers map[tgbotapi.Chat]struct{}
	handlers         map[string]*Handler
	Config           *Config
}

type Config struct {
	Token string `json:"token" yaml:"token"`
}

func NewTelegram(conf *Config) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		Api:              bot,
		VideoSubscribers: make(map[tgbotapi.Chat]struct{}),
		handlers:         make(map[string]*Handler),
	}, nil
}

var (
	ErrCmdAlreadyRegistered = errors.New("cmd already registered")
)

func (t *Telegram) AddCommandHandler(cmd string, desc string, fn HandleFunc) error {
	cmdHandler := &Handler{
		Command:     cmd,
		Description: desc,
		HandleFunc:  fn,
	}

	if _, exist := t.handlers[cmd]; exist {
		return ErrCmdAlreadyRegistered
	}
	t.handlers[cmd] = cmdHandler

	return nil
}

func (t *Telegram) HandleSubscribers() error {
	ch, err := t.Api.GetUpdatesChan(tgbotapi.NewUpdate(0))
	if err != nil {
		return err
	}
	go func() {
		for upd := range ch {

			if cmd, exist := t.handlers[upd.Message.Command()]; exist {
				cmd.HandleFunc(t, upd)

			}
		}
	}()

	return nil
}

func (t *Telegram) AnswerText(upd tgbotapi.Update, msg string) {
	answer := tgbotapi.NewMessage(upd.Message.Chat.ID, msg)
	t.Api.Send(answer)
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

			for sub := range t.VideoSubscribers {
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
