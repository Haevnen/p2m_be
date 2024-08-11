package interactor

import (
	"time"

	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto        *paseto.V2
	symmetric_key []byte
}

func NewPasetoMaker(symmetric_key string) (interactorinterface.Maker, error) {
	if len(symmetric_key) != chacha20poly1305.KeySize {
		return nil, apperror.ErrInvalidKeySize
	}

	pasetoMaker := &PasetoMaker{
		paseto:        paseto.NewV2(),
		symmetric_key: []byte(symmetric_key),
	}

	return pasetoMaker, nil
}

func (p *PasetoMaker) CreateToken(nickname string, isAdmin bool, duration time.Duration) (string, *interactorinterface.Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", nil, err
	}

	payload := &interactorinterface.Payload{
		ID:        tokenID,
		NickName:  nickname,
		IsAdmin:   isAdmin,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	token, err := p.paseto.Encrypt(p.symmetric_key, payload, nil)
	return token, payload, err
}

func (p *PasetoMaker) VerifyToken(token string) (*interactorinterface.Payload, error) {
	payload := &interactorinterface.Payload{}

	err := p.paseto.Decrypt(token, p.symmetric_key, payload, nil)
	if err != nil {
		return nil, apperror.ErrInvalidToken
	}

	if time.Now().After(payload.ExpiredAt) {
		return nil, apperror.ErrInvalidToken
	}

	return payload, nil

}
