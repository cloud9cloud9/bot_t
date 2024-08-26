package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"telegram-bot/pkg/repository"
)

func (b *Bot) initAuthProcess(message *tgbotapi.Message) error {
	authLink, err := b.genAuthorizationLink(message.Chat.ID)

	if err != nil {
		b.log.Error("unable to generate auth link", zap.Error(err))
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf(b.messages.Start, authLink))

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessToken)
}

func (b *Bot) genAuthorizationLink(chatID int64) (string, error) {
	redirectURI := b.generateRedirectURI(chatID)

	token, err := b.pocketClient.GetRequestToken(context.Background(), redirectURI)
	if err != nil {
		b.log.Error("unable to get request token", zap.Error(err))
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, token, repository.RequestToken); err != nil {
		b.log.Error("unable to save token", zap.Error(err))
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(token, redirectURI)

}

func (b *Bot) generateRedirectURI(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURI, chatID)
}
