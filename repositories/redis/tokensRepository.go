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
	conn *Connection
}

func NewTokensRepository(conn *Connection) *TokensRepository {
	return &TokensRepository{
		conn: conn,
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
		return errs.New(err) //TODO handle not found
	}

	token.UserID = val
	return nil
}
