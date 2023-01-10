package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func HandleSub(t *Telegram, upd tgbotapi.Update) {
	if _, exist := t.VideoSubscribers[*upd.Message.Chat]; !exist {
		log.Println("appended sub", upd.Message.Chat.UserName)
		t.VideoSubscribers[*upd.Message.Chat] = struct{}{}
		t.AnswerText(upd, "Вы успешно подписаны")
	}
}

func HandleHelp(t *Telegram, upd tgbotapi.Update) {
	var answer string
	for cmd, h := range t.handlers {
		desc := fmt.Sprintf("/%s: %s", cmd, h.Description)
		answer = fmt.Sprintf("%s \n %s", answer, desc)
	}
	t.AnswerText(upd, answer)
}
