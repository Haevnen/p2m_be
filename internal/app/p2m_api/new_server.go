package p2m_api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apiModel "github.com/Haevnen/p2m_be/gen/api"
)

// optional code omitted

type Server struct{}

func NewServer() Server {
	return Server{}
}

// PostUsersRegister (POST /users/register)
func (Server) PostUsersRegister(c *gin.Context, params apiModel.PostUsersRegisterParams) {

	c.JSON(http.StatusOK, nil)
}
