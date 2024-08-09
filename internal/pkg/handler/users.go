package handler

import (
	"net/http"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
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
func (h userHandler) InternalUserLogin(c *gin.Context, params p2m_api.InternalUserLoginParams) {}

// User logout
// (POST /logout)
func (h userHandler) InternalUserLogout(c *gin.Context, params p2m_api.InternalUserLogoutParams) {}

// Renew refresh_token and access token
// (POST /refresh-token)
func (h userHandler) InternalRefreshToken(c *gin.Context, params p2m_api.InternalRefreshTokenParams) {
}

// Get all users
// (GET /users)
func (h userHandler) InternalGetAllUsers(c *gin.Context, params p2m_api.InternalGetAllUsersParams) {
	user, err := h.userManagementInteractor.GetAllUser(c)
	if err != nil {
		sendError(c, "get all user error", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Register new user
// (POST /users/register)
func (h userHandler) InternalRegisterUser(c *gin.Context, params p2m_api.InternalRegisterUserParams) {
	var user p2m_api.User
	if err := bindRequestBody(c, &user); err != nil {
		sendError(c, "bind error", err)
		return
	}

	c.JSON(http.StatusOK, "Register very successful")
}

// Delete user by name
// (DELETE /users/{name})
func (h userHandler) InternalRemoveUser(c *gin.Context, name string, params p2m_api.InternalRemoveUserParams) {
}
