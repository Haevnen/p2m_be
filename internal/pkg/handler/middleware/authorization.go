package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/handler"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/Haevnen/p2m_be/pkg/constants"
)

var authorizationMap map[string]bool = map[string]bool{
	"/users":                     true,
	"/users/register":            true,
	"/users/:name":               true,
	"/clients":                   true,
	"/clients/:client_id":        true,
	"/clients/register":          true,
	"/clients/update/:client_id": true,
	"/links":                     true,
	"/links/update/:link_id":     true,
	"/links/:link_id":            true,
	"/tickets/add":               true,
	"/tickets/remove/:ticket_id": true,
}

func needsAuthorization(ctx *gin.Context) bool {
	relativePath := strings.TrimPrefix(ctx.FullPath(), constants.BaseURL)
	return authorizationMap[relativePath]
}

func Authorization() p2m_api.MiddlewareFunc {
	return func(ctx *gin.Context) {
		if needsAuthorization(ctx) {
			payload := ctx.MustGet(model.AuthorizationPayloadKey).(*interactorinterface.Payload)
			if !payload.IsAdmin {
				ctx.Abort()
				handler.SendError(ctx, "forbidden", apperror.ErrForbidden)
				return
			}
		}
		ctx.Next()
	}
}
