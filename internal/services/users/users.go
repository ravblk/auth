package users

import (
	"auth/internal/config"
	"auth/internal/model"
	repository "auth/internal/repositories"
	"auth/internal/services/users/passwords"
	"auth/tools/generator"
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/asaskevich/govalidator"
)

var (
	ErrEmail     = errors.New("incorrect email")
	ErrAut       = errors.New("authentication forbidden")
	ErrEmailBusy = errors.New("email busy")
)

type service struct {
	usrRep  repository.User
	sessRep repository.Session
	cfg     *config.API
	log     *zap.Logger
}

func NewService(log *zap.Logger, cfg *config.API, usrRep repository.User, sessRep repository.Session) *service {
	return &service{usrRep: usrRep, cfg: cfg, sessRep: sessRep, log: log}
}

func (s *service) UserLogin(ctx context.Context, l *model.Login, ss *model.Session) error {
	if !govalidator.IsEmail(l.Email) {
		return ErrEmail
	}
	u, err := s.usrRep.UserGet(ctx, l.Email)
	if err != nil {
		return err
	}
	valid, err := passwords.ValidateMD5(u.Password, l.Password)
	if err != nil {
		return err
	}
	if !valid {
		return ErrAut
	}
	ss.UserID = u.UUID
	if ss.SessionHash, err = generator.Token(); err != nil {
		return err
	}
	if err := s.sessRep.SessionPut(ctx, ss); err != nil {
		return err
	}
	return nil
}
func (s *service) UserRegistration(ctx context.Context, u *model.User, ss *model.Session) error {
	if !govalidator.IsEmail(u.Email) {
		return ErrEmail
	}
	ur, err := s.usrRep.UserGet(ctx, u.Email)
	if err != nil {
		return err
	}
	if ur.UUID != "" {
		return ErrEmailBusy
	}
	if err := passwords.Verify(u.Password); err != nil {
		return err
	}
	p, err := passwords.ToMD5(u.Password)
	if err != nil {
		return err
	}
	u.Password = p
	if err := s.usrRep.UserPut(ctx, u); err != nil {
		return err
	}
	hash, err := generator.Token()
	if err != nil {
		return err
	}
	ss.UserID = u.UUID
	ss.SessionHash = hash
	if err := s.sessRep.SessionPut(ctx, ss); err != nil {
		return err
	}
	return nil
}
