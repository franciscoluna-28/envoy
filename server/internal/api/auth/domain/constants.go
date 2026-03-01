package auth

import "time"

const (
	JWTExpirationHours    = 168
	JWTExpirationDuration = 24 * 7 * time.Hour
	AuthCookieName        = "session_token"
)
