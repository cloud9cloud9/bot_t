package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidURI   = errors.New("invalid URI")
	errUnauthorized = errors.New("user is not authorized")
	errUnableToSave = errors.New("unable to save")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.messages.Default)
	switch {
	case errors.Is(err, errInvalidURI):
		msg.Text = b.messages.InvalidURI
		b.bot.Send(msg)
	case errors.Is(err, errUnauthorized):
		msg.Text = b.messages.Unauthorized
		b.bot.Send(msg)
	case errors.Is(err, errUnableToSave):
		msg.Text = b.messages.UnableToSave
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}
}
