package pg

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	maxRetry = 10
	ttlRetry = 1 * time.Second
)

type Config struct {
	Host         string `valid:"required"`
	Port         string `valid:"required"`
	User         string `valid:"required"`
	Password     string `valid:"required"`
	DatabaseName string
	Schema       string
	SSL          string
	MaxIdleConns int
	MaxOpenConns int
}

type Database struct {
	Client *sqlx.DB
	log    *zap.Logger
}

func New(log *zap.Logger) *Database {
	return &Database{log: log}
}
func (d *Database) Connect(dbc *Config) error {
	args := fmt.Sprintf(
		"sslmode=%s host=%s port=%s user=%s password='%s' dbname=%s",
		dbc.SSL,
		dbc.Host,
		dbc.Port,
		dbc.User,
		dbc.Password,
		dbc.DatabaseName,
	)

	var (
		conn *sqlx.DB
		err  error
	)
	retry := 1
	for retry < maxRetry {
		conn, err = sqlx.Connect("postgres", args)
		if err != nil {
			d.log.Error("", zap.Int("#retrying:", retry), zap.String("second:", ttlRetry.String()))
			retry++
			time.Sleep(ttlRetry)
			continue
		}
		break
	}
	if conn != nil {
		if dbc.MaxIdleConns > 0 {
			conn.SetMaxIdleConns(dbc.MaxIdleConns)
		}
		if dbc.MaxOpenConns > 0 {
			conn.SetMaxOpenConns(dbc.MaxOpenConns)
		}
		d.Client = conn
	}

	return err
}
