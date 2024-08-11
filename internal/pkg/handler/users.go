package handler

import (
	"net/http"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userManagementInteractor interactorinterface.UserManagementInterface
}

func newUserHandler(registry *registry.Registry) userHandler {
	return userHandler{
		userManagementInteractor: registry.UserManagementInteractor(),
	}
}

// User login
// (POST /login)
func (h userHandler) InternalUserLogin(c *gin.Context, params p2m_api.InternalUserLoginParams) {
	var body p2m_api.UserLoginBody
	if err := bindRequestBody(c, &body); err != nil {
		SendError(c, err.Error(), apperror.ErrInvalidRequestInput)
		return
	}

	userLoginResponse, err := h.userManagementInteractor.LoginUser(c, body)
	if err != nil {
		SendError(c, "can not login user", err)
		return
	}

	c.JSON(http.StatusOK, userLoginResponse)
}

// User logout
// (POST /logout)
func (h userHandler) InternalUserLogout(c *gin.Context, params p2m_api.InternalUserLogoutParams) {
	var body p2m_api.RefreshTokenBody
	if err := bindRequestBody(c, &body); err != nil {
		SendError(c, err.Error(), apperror.ErrInvalidRequestInput)
		return
	}

	err := h.userManagementInteractor.LogoutUser(c, body)
	if err != nil {
		SendError(c, "can not logout user", err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// Renew refresh_token and access token
// (POST /refresh-token)
func (h userHandler) InternalRefreshToken(c *gin.Context, params p2m_api.InternalRefreshTokenParams) {
	var oldRefreshToken p2m_api.RefreshTokenBody
	if err := bindRequestBody(c, &oldRefreshToken); err != nil {
		SendError(c, err.Error(), apperror.ErrInvalidRequestInput)
		return
	}

	newAccessToken, err := h.userManagementInteractor.RefreshToken(c, oldRefreshToken)
	if err != nil {
		SendError(c, "can not refresh token", err)
		return
	}

	c.JSON(http.StatusOK, p2m_api.RefreshTokenResponse{
		AccessToken: newAccessToken,
	})

}

// Get all users
// (GET /users)
func (h userHandler) InternalGetAllUsers(c *gin.Context, params p2m_api.InternalGetAllUsersParams) {
	user, err := h.userManagementInteractor.GetAllUser(c)
	if err != nil {
		SendError(c, "get all user error", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Register new user
// (POST /users/register)
func (h userHandler) InternalRegisterUser(c *gin.Context, params p2m_api.InternalRegisterUserParams) {
	var user p2m_api.User
	if err := bindRequestBody(c, &user); err != nil {
		SendError(c, "bind error", err)
		return
	}

	c.JSON(http.StatusOK, "Register very successful")
}

// Delete user by name
// (DELETE /users/{name})
func (h userHandler) InternalRemoveUser(c *gin.Context, name string, params p2m_api.InternalRemoveUserParams) {
}
