package p2m_api

import (
	apiModel "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/di"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/user_service"
	"github.com/gin-gonic/gin"
	"net/http"

	api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
)

type userHandler struct {
	userUseCase user_service.UseCase
}

func newUserHandler(getter di.Getter) userHandler {
	return userHandler{userUseCase: getter.UserUseCase()}
}

// InternalUserLogin // (POST /login)
func (h userHandler) InternalUserLogin(c *gin.Context, params api.InternalUserLoginParams) {

	c.JSON(http.StatusOK, nil)
}

// InternalUserLogout // (POST /logout)
func (h userHandler) InternalUserLogout(c *gin.Context, params api.InternalUserLogoutParams) {

	c.JSON(http.StatusOK, nil)
}

// InternalRefreshToken // (POST /refresh-token)
func (h userHandler) InternalRefreshToken(c *gin.Context, params api.InternalRefreshTokenParams) {

	c.JSON(http.StatusOK, nil)
}

// InternalGetAllUsers // (GET /users)
func (h userHandler) InternalGetAllUsers(c *gin.Context, params api.InternalGetAllUsersParams) {

	c.JSON(http.StatusOK, nil)
}

// InternalRegisterUser (POST /users/register)
func (h userHandler) InternalRegisterUser(c *gin.Context, params api.InternalRegisterUserParams) {
	var user apiModel.User
	if err := bindRequestBody(c, &user); err != nil {
		sendError(c, "bind error", err)
		return
	}

	// validate api model implementation

	// convert api model to use case model
	userDb := model.User{}

	// use case call
	err := h.userUseCase.Add(c, userDb.ToDB(user))
	if err != nil {
		sendError(c, "create error", err)
		return
	}

	c.JSON(http.StatusOK, userDb.ToAPI())
}

// InternalRemoveUser // (DELETE /users/{name})
func (h userHandler) InternalRemoveUser(c *gin.Context, name string, params api.InternalRemoveUserParams) {

	c.JSON(http.StatusOK, nil)
}
