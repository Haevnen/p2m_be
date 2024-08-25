package interactor

import (
	"github.com/gin-gonic/gin"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
)

type CommentManagement struct {
}

func NewCommentManagement() *CommentManagement {
	return &CommentManagement{}
}

func (ci *CommentManagement) CreateComment(ctx *gin.Context, client p2mapi.CreateCommentBody) (*p2mapi.CommentResponse, error) {
	c := dal.Q.Comment

	var commentDb model.Comment
	commentDb.ToComment(client)
	// add userId
	payload := ctx.MustGet(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return nil, apperror.ErrUserNotExists
	}
	commentDb.UserID = payload.UserID

	err := c.WithContext(ctx).Save(&commentDb)
	if err != nil {
		return nil, err
	}

	return commentDb.FromComment(), nil
}
