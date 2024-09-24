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

func needsAuthentication(ctx *gin.Context) bool {
	// There are some API that don't need authentication
	// 1. Login
	// 2. Refresh token

	relativePath := strings.TrimPrefix(ctx.FullPath(), constants.BaseURL)
	return !(relativePath == "/login" || relativePath == "/refresh-token" || relativePath == "/tickets/add_auto")
}

func Authentication(tokenMaker interactorinterface.Maker) p2m_api.MiddlewareFunc {
	return func(ctx *gin.Context) {
		if needsAuthentication(ctx) {
			authorizationHeader := ctx.GetHeader(model.AuthorizationHeaderKey)
			if len(authorizationHeader) == 0 {
				ctx.Abort()
				handler.SendError(ctx, "authorization header is not provided", apperror.ErrNotProvidedAuthenticationHeader)
				return
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				ctx.Abort()
				handler.SendError(ctx, "invalid authorization header format", apperror.ErrInvalidAuthorizationHeaderFormat)
				return
			}

			authorizationType := strings.ToLower((fields[0]))
			if authorizationType != model.AuthorizationTypeBearer {
				ctx.Abort()
				handler.SendError(ctx, "unsupported authorization type", apperror.ErrUnsupportedAuthorizationType)
				return
			}

			accessToken := fields[1]
			payload, err := tokenMaker.VerifyToken(accessToken)
			if err != nil {
				ctx.Abort()
				handler.SendError(ctx, err.Error(), apperror.ErrInvalidToken)
				return
			}

			ctx.Set(model.AuthorizationPayloadKey, payload)
		}
		ctx.Next()
	}
}
