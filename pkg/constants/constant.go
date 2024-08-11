package constants

import "time"

const (
	BaseURL              = "/api/v1"
	AccessTokenDuration  = 30 * time.Minute
	RefreshTokenDuration = 48 * time.Hour
)
