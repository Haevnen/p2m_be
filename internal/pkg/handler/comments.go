package handler

import (
	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/gin-gonic/gin"
	"net/http"
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
func (h commentHandler) InternalUpdateComment(c *gin.Context, commentId int, params p2m_api.InternalUpdateCommentParams) {
}

// Delete comment by id
// (DELETE /comments/{comment_id})
func (h commentHandler) InternalRemoveComment(c *gin.Context, commentId int, params p2m_api.InternalRemoveCommentParams) {
}

// Get all comments
// (GET /comments/{ticket_id})
func (h commentHandler) InternalGetComments(c *gin.Context, ticketId int64) {
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
