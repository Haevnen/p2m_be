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
	GetAllUser(ctx context.Context, IncludingUnassigned *bool) ([]*p2mapi.User, error)
	CreateUser(ctx context.Context, user p2mapi.CreateUserBody) (*p2mapi.User, error)
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
	GetSingleClient(ctx context.Context, clientID string) (*p2mapi.ClientResponse, error)
}

type LinkManagementInterface interface {
	CreateLink(ctx context.Context, link p2mapi.LinkBody) (*p2mapi.LinkResponse, error)
	GetAllLink(ctx context.Context, ticketID int64) ([]*p2mapi.LinkResponse, error)
	RemoveLink(ctx context.Context, linkID int64) error
	UpdateLink(ctx context.Context, linkID int64, body p2mapi.LinkBody) error
}

type CommentManagementInterface interface {
	CreateComment(ctx context.Context, client p2mapi.CreateCommentBody) (*p2mapi.CommentResponse, error)
	UpdateComment(ctx context.Context, commentID int64, body p2mapi.UpdateCommentBody) error
	DeleteComment(ctx context.Context, commentID int64) error
	GetAllComment(ctx context.Context, ticketID int64) ([]*p2mapi.CommentResponse, error)
}

type HistoryManagementInterface interface {
	GetAllHistoriesByTicket(ctx context.Context, ticketID int64) ([]*p2mapi.HistoryResponse, error)
}

type TicketManagementInterface interface {
	AddTicketManual(ctx context.Context, body p2mapi.CreateTicketBody) error
	UpdateTicket(ctx context.Context, ticketID int64, body p2mapi.UpdateTicketBody) error
}
type Maker interface {
	// Return token, payload and error
	CreateToken(userID string, isAdmin bool, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

type TxManager interface {
	TransactionExec(ctx context.Context, fn func(context.Context) error) error
}
