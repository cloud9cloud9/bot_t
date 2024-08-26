package boltdb

import (
	"errors"
	"github.com/boltdb/bolt"
	"go.uber.org/zap"
	"strconv"
	repo "telegram-bot/pkg/repository"
)

type TokenRepository struct {
	db  *bolt.DB
	log *zap.Logger
}

func NewTokenRepository(db *bolt.DB, log *zap.Logger) *TokenRepository {
	return &TokenRepository{db: db, log: log}
}

func (r *TokenRepository) Get(chatID int64, bucket repo.Bucket) (string, error) {
	var token string

	r.log.Debug("get token", zap.Int64("chat_id", chatID))
	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		data := bucket.Get(intToByte(chatID))
		token = string(data)
		return nil
	})
	if err != nil {
		r.log.Error("unable to get token", zap.Error(err))
		return "", err
	}

	if token == "" {
		r.log.Debug("token not found", zap.Int64("chat_id", chatID))
		return "", errors.New("token not found")
	}

	r.log.Debug("token found", zap.String("token", token))
	return token, nil
}

func (r *TokenRepository) Save(chatID int64, token string, bucket repo.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		r.log.Debug("saving token", zap.Int64("chat_id", chatID), zap.String("token", token))
		bucket := tx.Bucket([]byte(bucket))
		return bucket.Put(intToByte(chatID), []byte(token))
	})
}

func intToByte(value int64) []byte {
	return []byte(strconv.FormatInt(value, 10))
}
