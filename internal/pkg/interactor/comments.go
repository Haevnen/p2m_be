package interactor

import (
	"errors"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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

func (ci *CommentManagement) UpdateComment(ctx *gin.Context, commentID int64, body p2mapi.UpdateCommentBody) error {
	c := dal.Q.Comment

	comment, err := c.WithContext(ctx).Where(c.ID.Eq(commentID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrRecordNotFound
		}
		return err
	}

	comment.Comment = body.Comment

	_, err = c.WithContext(ctx).Where(c.ID.Eq(commentID)).UpdateColumns(comment)
	return err
}

func (ci *CommentManagement) DeleteComment(ctx *gin.Context, commentID int64) error {
	c := dal.Q.Comment
	comment, err := c.WithContext(ctx).Where(c.ID.Eq(commentID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrRecordNotFound
		}
	}

	_, err = c.WithContext(ctx).Where(c.ID.Eq(commentID)).Delete(comment)
	return err
}

func (ci *CommentManagement) GetAllComment(ctx *gin.Context) ([]*p2mapi.CommentResponse, error) {
	c := dal.Q.Comment
	var comments []*model.Comment
	comments, err := c.WithContext(ctx).Order(c.CreatedAt.Desc()).Find()

	if err != nil {
		return nil, err
	}

	res := make([]*p2mapi.CommentResponse, 0, len(comments))
	for _, client := range comments {
		res = append(res, client.FromComment())
	}

	return res, nil
}
