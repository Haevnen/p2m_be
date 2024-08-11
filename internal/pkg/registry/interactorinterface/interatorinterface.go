package interactorinterface

import (
	"context"
	"time"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	NickName  string    `json:"nick_name"`
	IsAdmin   bool      `json:"is_admin"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UserManagementInterface interface {
	GetAllUser(context.Context) (*model.User, error)
	LoginUser(context.Context, p2m_api.UserLoginBody) (p2m_api.UserLoginResponse, error)
	LogoutUser(context.Context, p2m_api.RefreshTokenBody) error
	RefreshToken(context.Context, p2m_api.RefreshTokenBody) (string, error)
}

type Maker interface {
	// Return token, payload and error
	CreateToken(nickname string, isAdmin bool, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
