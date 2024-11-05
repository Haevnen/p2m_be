package interactor

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/Haevnen/p2m_be/pkg/logger"
)

type CommentManagement struct {
	txManager interactorinterface.TxManager
}

func NewCommentManagement(
	txManager interactorinterface.TxManager,
) *CommentManagement {
	return &CommentManagement{txManager: txManager}
}

func (ci *CommentManagement) CreateComment(ctx context.Context, client p2mapi.CreateCommentBody) (*p2mapi.CommentResponse, error) {
	u := dal.Q.User

	var commentDb model.Comment
	commentDb.ToComment(client)
	// add userId
	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return nil, apperror.ErrUserNotExists
	}
	commentDb.UserID = payload.UserID

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

	err = ci.txManager.TransactionExec(ctx, func(childCtx context.Context) error {
		tx := childCtx.Value(txTransactionKey).(*dal.QueryTx)

		// Create new comment
		err := tx.Comment.WithContext(childCtx).Create(&commentDb)
		if err != nil {
			return err
		}

		// Create new history
		return tx.History.WithContext(childCtx).Create(&model.History{
			TicketID:    commentDb.TicketID,
			Action:      fmt.Sprintf("User %s create comment for ticket: %v", nickName, commentDb.TicketID),
			PerformedBy: payload.UserID,
		})
	})

	if err != nil {
		return nil, err
	}

	return commentDb.FromComment(nickName).FromCommentWithName(), nil
}

func (ci *CommentManagement) UpdateComment(ctx context.Context, commentID int64, body p2mapi.UpdateCommentBody) error {
	u := dal.Q.User
	c := dal.Q.Comment

	comment, err := c.WithContext(ctx).Where(c.ID.Eq(commentID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrRecordNotFound
		}
		return err
	}

	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return apperror.ErrUserNotExists
	}

	if payload.UserID != comment.UserID {
		return apperror.ErrPermissionDenied
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

	comment.Comment = body.Comment

	return ci.txManager.TransactionExec(ctx, func(childCtx context.Context) error {
		tx := childCtx.Value(txTransactionKey).(*dal.QueryTx)

		_, err = tx.Comment.WithContext(childCtx).Where(c.ID.Eq(commentID)).UpdateColumns(comment)
		if err != nil {
			return err
		}

		// Create new history
		return tx.History.WithContext(childCtx).Create(&model.History{
			TicketID:    comment.TicketID,
			Action:      fmt.Sprintf("User %s update comment for ticket: %v", nickName, comment.TicketID),
			PerformedBy: payload.UserID,
		})
	})
}

func (ci *CommentManagement) DeleteComment(ctx context.Context, commentID int64) error {
	u := dal.Q.User
	c := dal.Q.Comment
	comment, err := c.WithContext(ctx).Where(c.ID.Eq(commentID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrRecordNotFound
		}
	}

	payload := ctx.Value(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
	if payload == nil {
		return apperror.ErrUserNotExists
	}

	if payload.UserID != comment.UserID {
		return apperror.ErrPermissionDenied
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

	return ci.txManager.TransactionExec(ctx, func(childCtx context.Context) error {
		tx := childCtx.Value(txTransactionKey).(*dal.QueryTx)

		_, err = tx.Comment.WithContext(childCtx).Where(c.ID.Eq(commentID)).Delete(comment)
		if err != nil {
			return err
		}

		// Create new history
		return tx.History.WithContext(childCtx).Create(&model.History{
			TicketID:    comment.TicketID,
			Action:      fmt.Sprintf("User %s delete comment for ticket: %v", nickName, comment.TicketID),
			PerformedBy: payload.UserID,
		})
	})
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
