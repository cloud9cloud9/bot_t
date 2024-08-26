package server

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zhashkevych/go-pocket-sdk"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"telegram-bot/pkg/repository"
)

type AuthServer struct {
	server *http.Server

	storage repository.TokenRepository
	client  *pocket.Client

	redirectUrl string
	log         *zap.Logger
}

func NewAuthServer(
	redirectUrl string,
	storage repository.TokenRepository,
	client *pocket.Client,
	log *zap.Logger,
) *AuthServer {
	return &AuthServer{
		redirectUrl: redirectUrl,
		storage:     storage,
		client:      client,

		log: log,
	}
}

func (s *AuthServer) Start() error {
	s.server = &http.Server{
		Handler: s,
		Addr:    ":80",
	}

	return s.server.ListenAndServe()
}

func (s *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.log.Info("method not allowed", zap.String("method", r.Method))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	chatIDQuery := r.URL.Query().Get("chat_id")
	if chatIDQuery == "" {
		s.log.Info("chat_id not found", zap.String("query", chatIDQuery))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDQuery, 10, 64)
	if err != nil {
		s.log.Info("failed to parse chat_id", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.createAccessToken(r.Context(), chatID); err != nil {
		s.log.Info("failed to create access token", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.log.Info("access token created", zap.Int64("chat_id", chatID))
	w.Header().Set("Location", s.redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}

func (s *AuthServer) createAccessToken(ctx context.Context, chatID int64) error {
	requestToken, err := s.storage.Get(chatID, repository.RequestToken)
	if err != nil {
		s.log.Error("failed to get request token", zap.Error(err))
		return errors.WithMessage(err, "failed to get request token")
	}

	authResp, err := s.client.Authorize(ctx, requestToken)
	if err != nil {
		s.log.Error("failed to authorize at Pocket", zap.Error(err))
		return errors.WithMessage(err, "failed to authorize at Pocket")
	}

	if err := s.storage.Save(chatID, authResp.AccessToken, repository.AccessToken); err != nil {
		s.log.Error("failed to save access token", zap.Error(err))
		return err
	}

	s.log.Info("access token saved", zap.Int64("chat_id", chatID))
	return nil
}
