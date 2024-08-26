package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"go.uber.org/zap"
	"telegram-bot/pkg/config"
	"telegram-bot/pkg/repository"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURI     string

	messages config.Messages
	log      *zap.Logger
}

func NewBot(
	bot *tgbotapi.BotAPI,
	pocketClient *pocket.Client,
	tokenRepository repository.TokenRepository,
	redirectURI string,
	messages config.Messages,
	log *zap.Logger,
) *Bot {
	return &Bot{
		bot:             bot,
		pocketClient:    pocketClient,
		tokenRepository: tokenRepository,
		redirectURI:     redirectURI,
		messages:        messages,
		log:             log,
	}
}

func (b *Bot) Start() error {
	b.log.Info("Starting bot...")
	updates, err := b.initUpdatesChan()

	if err != nil {
		b.log.Error("unable to init updates channel", zap.Error(err))
		return err
	}
	b.handleUpdates(updates)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {

		if update.Message.IsCommand() {
			b.log.Info("Received command", zap.String("command", update.Message.Text))
			if err := b.handleCommand(update.Message); err != nil {
				b.log.Error("unable to handle command", zap.Error(err))
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.log.Error("unable to handle message", zap.Error(err))
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}

func (b *Bot) initUpdatesChan() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
