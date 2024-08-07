package p2m_api

import (
	apiModel "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/di"
	"github.com/gin-gonic/gin"
	"net/http"
)

// optional code omitted

type Server struct {
	userHandler
}

func NewServer(getter di.Getter) Server {
	s := Server{}
	s.userHandler = newUserHandler(getter)

	return s
}

// GetPing (GET /ping)
func (Server) GetPing(ctx *gin.Context) {
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

func sendError(ctx *gin.Context, title string, err error) {

	appErr := apperror.New(ctx, err)

	ctx.JSON(appErr.HTTPStatus(), apiModel.Error{
		Type:   apiModel.ErrorType(appErr.ResType()),
		Title:  title,
		Code:   appErr.ErrorCode(),
		Detail: appErr.Detail(),
	})
}
