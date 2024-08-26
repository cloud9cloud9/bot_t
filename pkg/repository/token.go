package repository

type Bucket string

const (
	AccessToken  Bucket = "access_tokens"
	RequestToken Bucket = "request_tokens"
)

type TokenRepository interface {
	Get(chatID int64, bucket Bucket) (string, error)
	Save(chatID int64, token string, bucket Bucket) error
}
