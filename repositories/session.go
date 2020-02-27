package repository

import (
	"auth/model"
	"auth/storage/pg"
	"context"

	"go.uber.org/zap"
)

type Session interface {
	SessionPut(ctx context.Context, l *model.Session) error
}

type sessRep struct {
	db  *pg.Database
	log *zap.Logger
}

func NewSessionRepository(log *zap.Logger, db *pg.Database) *sessRep {
	return &sessRep{db: db, log: log}
}

func (r *sessRep) SessionPut(ctx context.Context, s *model.Session) error {
	if err := r.db.Client.QueryRowxContext(ctx, `INSERT INTO sessions(
		session_hash,
		user_uuid,
		ip_address,
		user_agent,
		created_at)
	VALUES(
		$1,
		$2,
		$3,
		$4,
		NOW()) 
	returning guid;`, s.SessionHash, s.UserID, s.IP, s.UserAgent).Scan(&s.UserID); err != nil {
		return err
	}
	return nil
}
