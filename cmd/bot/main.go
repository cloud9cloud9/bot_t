package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"go.uber.org/zap"
	"telegram-bot/pkg/config"
	"telegram-bot/pkg/repository"
	"telegram-bot/pkg/repository/boltdb"
	"telegram-bot/pkg/server"
	"telegram-bot/pkg/telegram"
)

func main() {
	logg := zap.Must(zap.NewProduction())

	logg.Info("starting reading cfg files...")
	cfg, err := config.Init()

	if err != nil {
		logg.Fatal("unable to read cfg files", zap.Error(err))
	}
	logg.Warn("cfg loaded", zap.Any("cfg", cfg))

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		logg.Fatal("unable to create bot", zap.Error(err))
	}
	bot.Debug = true

	logg.Info("Authorized on account")

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		logg.Fatal("unable to create pocket client", zap.Error(err))
	}

	db, err := initDB(cfg)
	if err != nil {
		logg.Fatal("unable to init db", zap.Error(err))
	}

	tokenRepository := boltdb.NewTokenRepository(db, logg)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURI, cfg.Messages, logg)

	authServer := server.NewAuthServer(cfg.TelegramBotURI, tokenRepository, pocketClient, logg)

	go func() {
		if err := telegramBot.Start(); err != nil {
			logg.Fatal("unable to start telegram bot", zap.Error(err))
		}
	}()
	logg.Info("Telegram bot started, starting listening...")
	err = authServer.Start()
	if err != nil {
		logg.Fatal("unable to start auth server", zap.Error(err))
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessToken))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestToken))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
