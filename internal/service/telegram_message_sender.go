package service

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type TelegramMessageSender struct{}

func NewTelegramMessageSender() *TelegramMessageSender {
	return &TelegramMessageSender{}
}

func (s *TelegramMessageSender) Send(_ context.Context, token string, userId int64, msg bots.Message) error {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	// fmt.Printf("%v\n", msg.Options())
	m := tgbotapi.NewMessage(userId, msg.Text)
	m.ParseMode = tgbotapi.ModeHTML
	if len(msg.Buttons) > 0 {
		keyboard := buildInlineKeyboardMarkup(msg.Buttons)
		m.ReplyMarkup = keyboard
	} else {
		m.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}

	_, err = api.Send(m)
	return err
}

func buildInlineKeyboardMarkup(opts []string) tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, len(opts))
	for i, opt := range opts {
		rows[i] = []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton(opt),
		}
	}
	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true
	return keyboard
}
