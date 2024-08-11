package interactorinterface

import (
	"context"
	"time"

	"github.com/google/uuid"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	NickName  string    `json:"nick_name"`
	IsAdmin   bool      `json:"is_admin"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UserManagementInterface interface {
	GetAllUser(ctx context.Context, includeDeActive *bool) ([]*p2mapi.User, error)
	CreateUser(ctx context.Context, user p2mapi.User) (*p2mapi.User, error)
	RemoveUser(ctx context.Context, nickName string) error
}

type Maker interface {
	// Return token, payload and error
	CreateToken(nickname string, isAdmin bool, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
