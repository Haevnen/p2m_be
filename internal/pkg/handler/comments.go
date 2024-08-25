package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
)

type commentHandler struct {
	commentManagementInteractor interactorinterface.CommentManagementInterface
}

func newCommentHandler(registry *registry.Registry) commentHandler {
	return commentHandler{
		commentManagementInteractor: registry.CommentManagementInterface(),
	}
}

// Update comment by id
// (PUT /comments/update/{comment_id})
func (h commentHandler) InternalUpdateComment(c *gin.Context, commentId int64, params p2m_api.InternalUpdateCommentParams) {
	var updateCommentBody p2m_api.UpdateCommentBody
	if err := bindRequestBody(c, &updateCommentBody); err != nil {
		SendError(c, err.Error(), apperror.ErrInvalidRequestInput)
		return
	}

	err := h.commentManagementInteractor.UpdateComment(c, commentId, updateCommentBody)
	if err != nil {
		SendError(c, "update comment error", err)
		return
	}

	c.JSON(http.StatusOK, "update comment successfully")
}

// Delete comment by id
// (DELETE /comments/{comment_id})
func (h commentHandler) InternalRemoveComment(c *gin.Context, commentId int64, params p2m_api.InternalRemoveCommentParams) {
	err := h.commentManagementInteractor.DeleteComment(c, commentId)
	if err != nil {
		SendError(c, "delete comment error", err)
		return
	}

	c.JSON(http.StatusOK, "delete comment successfully")
}

// Get all comments
// (GET /comments/{ticket_id})
func (h commentHandler) InternalGetComments(c *gin.Context, ticketId int64) {
	comments, err := h.commentManagementInteractor.GetAllComment(c)
	if err != nil {
		SendError(c, "get all comment error", err)
		return
	}

	c.JSON(http.StatusOK, comments)
}

// Create new comment
// (POST /comments)
func (h commentHandler) InternalCreateComment(c *gin.Context, params p2m_api.InternalCreateCommentParams) {
	var client p2m_api.CreateCommentBody
	if err := bindRequestBody(c, &client); err != nil {
		SendError(c, "bind error", err)
		return
	}

	res, err := h.commentManagementInteractor.CreateComment(c, client)
	if err != nil {
		SendError(c, "create comment error", err)
		return
	}

	c.JSON(http.StatusOK, res)
}
