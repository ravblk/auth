package users

import (
	"auth/model"
	"context"
)

type Service interface {
	UserLogin(ctx context.Context, l *model.Login, ss *model.Session) error
	UserRegistration(ctx context.Context, u *model.User, ss *model.Session) error
}
