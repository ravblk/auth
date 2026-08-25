package services

import (
	"auth/internal/services/users"

	"go.uber.org/zap"
)

type Auth struct {
	UsrSvc users.Service
	Log    *zap.Logger
}
