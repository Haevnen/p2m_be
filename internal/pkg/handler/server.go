package handler

import (
	"net/http"

	apiModel "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userHandler
	clientHandler
	linkHandler
	commentHandler
	historyHandler
	ticketHandler
}

func New(reg *registry.Registry) Handler {
	h := Handler{}
	h.userHandler = newUserHandler(reg)
	h.clientHandler = newClientHandler(reg)
	h.linkHandler = newLinkHandler(reg)
	h.commentHandler = newCommentHandler(reg)
	h.historyHandler = newHistoryHandler(reg)
	h.ticketHandler = newTicketHandler(reg)
	return h
}

// GetPing (GET /ping)
func (Handler) GetPing(ctx *gin.Context) {
	resp := apiModel.Pong{
		Ping: "pong",
	}

	ctx.JSON(http.StatusOK, resp)
}

func bindRequestBody(ctx *gin.Context, body interface{}) error {
	if err := ctx.ShouldBindJSON(body); err != nil {
		return err
	}
	return nil
}

func SendError(ctx *gin.Context, title string, err error) {

	appErr := apperror.New(ctx, err)

	ctx.JSON(appErr.HTTPStatus(), apiModel.Error{
		Type:   apiModel.ErrorType(appErr.ResType()),
		Title:  title,
		Code:   appErr.ErrorCode(),
		Detail: appErr.Detail(),
	})
}
