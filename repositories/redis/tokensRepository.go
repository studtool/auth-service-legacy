package redis

import (
	"time"

	"github.com/studtool/common/errs"
	"github.com/studtool/common/utils"

	"github.com/studtool/auth-service/models"
)

const (
	TokenLen = 100
	TokenExp = 10 * time.Minute
)

type TokensRepository struct {
	conn        *Connection
	notFoundErr *errs.Error
}

func NewTokensRepository(conn *Connection) *TokensRepository {
	return &TokensRepository{
		conn:        conn,
		notFoundErr: errs.NewNotFoundError("token not found"),
	}
}

func (r *TokensRepository) SetToken(token *models.Token) *errs.Error {
	token.Token = utils.RandString(TokenLen)

	err := r.conn.client.Set(token.Token, token.UserID, TokenExp).Err()
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *TokensRepository) GetToken(token *models.Token) *errs.Error {
	val, err := r.conn.client.Get(token.Token).Result()
	if err != nil {
		if r.isErrNotFound(err) {
			return r.notFoundErr
		}
		return errs.New(err)
	}

	token.UserID = val
	return nil
}

func (r *TokensRepository) DeleteToken(token *models.Token) *errs.Error {
	if err := r.conn.client.Del(token.Token).Err(); err != nil {
		return errs.New(err)
	}
	return nil
}

func (r *TokensRepository) isErrNotFound(err error) bool {
	return err.Error() == "redis: nil"
}
