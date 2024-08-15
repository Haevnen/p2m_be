package interactorinterface

import (
	"context"
	"time"

	"github.com/google/uuid"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	IsAdmin   bool      `json:"is_admin"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UserManagementInterface interface {
	GetAllUser(ctx context.Context, includeDeActive *bool) ([]*p2mapi.User, error)
	CreateUser(ctx context.Context, user p2mapi.User) (*p2mapi.User, error)
	RemoveUser(ctx context.Context, nickName string) error
	LoginUser(ctx context.Context, body p2mapi.UserLoginBody) (p2mapi.UserLoginResponse, error)
	LogoutUser(context.Context, p2mapi.RefreshTokenBody) error
	RefreshToken(context.Context, p2mapi.RefreshTokenBody) (p2mapi.RefreshTokenResponse, error)
}

type ClientManagementInterface interface {
	CreateClient(ctx context.Context, client p2mapi.ClientBody) (*p2mapi.ClientResponse, error)
	RemoveClient(ctx context.Context, id string) error
	GetAllClient(ctx context.Context, includeDeActive *bool) ([]*p2mapi.ClientResponse, error)
	UpdateClient(ctx context.Context, clientID string, body p2mapi.UpdateClientBody) error
}

type Maker interface {
	// Return token, payload and error
	CreateToken(userID string, isAdmin bool, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
