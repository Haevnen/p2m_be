package constants

import "time"

const (
	BaseURL              = "/api/v1"
	AccessTokenDuration  = 30 * time.Minute
	RefreshTokenDuration = 48 * time.Hour
	DateTimeFormat       = "2006-01-02T15:04:05Z07:00"
	DateFormat           = "2006-01-02"
)
