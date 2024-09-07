package interactor

import (
	"context"
	"errors"

	"gorm.io/gorm"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/Haevnen/p2m_be/pkg/logger"
)

type CommentManagement struct {
}

func NewCommentManagement() *CommentManagement {
	return &CommentManagement{}
}

func (ci *CommentManagement) CreateComment(ctx context.Context, client p2mapi.CreateCommentBody) (*p2mapi.CommentResponse, error) {
	c := dal.Q.Comment
	u := dal.Q.User

	var commentDb model.Comment
	commentDb.ToComment(client)
	// add userId
	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return nil, apperror.ErrUserNotExists
	}
	commentDb.UserID = payload.UserID

	err := c.WithContext(ctx).Save(&commentDb)
	if err != nil {
		return nil, err
	}

	// get nick-name from userID
	nickName := ""
	user, err := u.WithContext(ctx).Where(u.UserID.Eq(payload.UserID)).First()
	if err != nil {
		logger.Errorf("find user when create comment has error: %v", err)
	}
	if user != nil {
		nickName = user.NickName
	} else {
		logger.Error("find user when create comment has error: user is nil")
	}

	return commentDb.FromComment(nickName).FromCommentWithName(), nil
}

func (ci *CommentManagement) UpdateComment(ctx context.Context, commentID int64, body p2mapi.UpdateCommentBody) error {
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

func (ci *CommentManagement) DeleteComment(ctx context.Context, commentID int64) error {
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

func (ci *CommentManagement) GetAllComment(ctx context.Context, ticketID int64) ([]*p2mapi.CommentResponse, error) {
	c := dal.Q.Comment
	u := dal.Q.User

	var comments []*model.CommentWithName
	err := c.WithContext(ctx).Select(c.ID, c.TicketID, c.UserID, c.Comment, c.CreatedAt, u.NickName).
		Join(u, c.UserID.EqCol(u.UserID)).
		Where(c.TicketID.Eq(ticketID)).Order(c.CreatedAt.Desc()).Scan(&comments)

	if err != nil {
		return nil, err
	}

	res := make([]*p2mapi.CommentResponse, 0, len(comments))
	for _, client := range comments {
		res = append(res, client.FromCommentWithName())
	}

	return res, nil
}
