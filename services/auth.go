package services

import (
	"auth/config"
	"auth/services/users"

	"go.uber.org/zap"
)

type Auth struct {
	UsrSvc users.Service
	Log    *zap.Logger
	Cfg    *config.Service
}
