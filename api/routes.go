package api

import (
	"LocalDex/api/auth"
	"net/http"
)

var ApiRoutes = map[string]http.HandlerFunc{
	"POST /auth/login":      auth.SendOTPHandler,
	"POST /auth/verify_otp": auth.VerifyOTPHandler,
	"GET /auth/status":      auth.VerifyAuthStatus,
}
