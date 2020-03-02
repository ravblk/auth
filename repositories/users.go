package repository

import (
	"auth/model"
	"auth/storage/pg"
	"context"

	"go.uber.org/zap"
)

type User interface {
	UserPut(ctx context.Context, u *model.User) error
	UserGet(ctx context.Context, email string) (*model.User, error)
}
type usrRep struct {
	db  *pg.Database
	log *zap.Logger
}

func NewUserRepository(log *zap.Logger, db *pg.Database) *usrRep {
	return &usrRep{db: db, log: log}
}

func (r *usrRep) UserPut(ctx context.Context, u *model.User) error {
	if err := r.db.Client.QueryRowxContext(ctx, `INSERT INTO users(
		uuid,
		email,
		password,
		first_name,
		last_name,
		created_at)
	VALUES(
		gen_random_uuid(),
		$1,
		$2,
		$3,
		$4,
		NOW()) 
	returning uuid;`, u.Email, u.Password, u.FirstName, u.LastName).Scan(&u.UUID); err != nil {
		return err
	}
	return nil
}

func (r *usrRep) UserGet(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	if err := r.db.Client.GetContext(ctx, u, `select * from  users
	WHERE 
		email = $1 `, email); err != nil {
		return nil, err
	}
	return u, nil
}
