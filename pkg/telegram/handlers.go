package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"go.uber.org/zap"
	"net/url"
)

const (
	commandStart = "start"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}

}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		b.log.Error("invalid URI", zap.Error(err))
		return errInvalidURI
	}
	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		b.log.Error("unauthorized", zap.Error(err))
		return errUnauthorized
	}

	if err = b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		b.log.Error("unable to save", zap.Error(err))
		return errUnableToSave
	}

	b.log.Info("successfully saved", zap.String("url", message.Text))
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.SavedSuccessfully)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		b.log.Info("user is not authorized", zap.Error(err))
		return b.initAuthProcess(message)
	}

	b.log.Info("user is already authorized")
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.AlreadyAuthorized)
	b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}
