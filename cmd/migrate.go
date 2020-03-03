package cmd

import (
	"auth/config"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	migrateCmd = &cobra.Command{
		Use:   "migrate [sub]",
		Short: "migration db auth",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				zap.L().Fatal("Wrong arrgument enter up or down")
			}
			switch args[0] {
			case "up":
				upCmd.Run(cmd, args)
			case "down":
				downCmd.Run(cmd, args)
			}
		},
	}
	upCmd = &cobra.Command{
		Use:   "up [no options!]",
		Short: "migrate up",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			defer func() {
				if err != nil {
					zap.L().Fatal("", zap.Error(err))
				}
			}()

			db, err := migrationInit()
			if err != nil {
				return
			}
			defer db.Close()
			up(db)
		},
	}
	downCmd = &cobra.Command{
		Use:   "down [no options!]",
		Short: "migrate down",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			defer func() {
				if err != nil {
					zap.L().Fatal("", zap.Error(err))
				}
			}()

			db, err := migrationInit()
			if err != nil {
				return
			}
			defer db.Close()
			down(db)
		},
	}
	migrations = &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			&migrate.Migration{
				Id: "1",
				Up: []string{
					`CREATE TABLE users(
					uid           uuid   PRIMARY KEY  NOT NULL,
					email          text    NOT NULL,
					password        text   NOT NULL,
					first_name      text,
					last_name 		text,
					created_at  TIMESTAMP
				 );`,
					`CREATE TABLE sessions(
					session_hash  text   PRIMARY KEY,
					user_id           uuid references users(uid),
					ip_address      text not null,
					user_agent text default '',
					created_at  TIMESTAMP
				 );`},
				Down: []string{
					"DROP TABLE sessions;",
					"DROP TABLE users;"},
			},
		},
	}
)

func init() {
	RootCmd.AddCommand(migrateCmd)
}

func migrationInit() (*sql.DB, error) {
	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("", zap.Error(err))
	}

	cfg, err := config.Read()
	if err != nil {
		log.Warn("", zap.Error(errors.New("wrong cfg")))
		return nil, err
	}
	args := fmt.Sprintf(
		"sslmode=%s host=%s port=%s user=%s password='%s' dbname=%s",
		cfg.DB.SSL,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.DatabaseName,
	)
	postDB, err := sql.Open("postgres", args)
	if err != nil {
		zap.L().Warn("", zap.Error(err))
		return nil, err
	}
	return postDB, nil
}

func up(postDB *sql.DB) {
	affected, err := migrate.Exec(postDB, "postgres", migrations, migrate.Up)
	if err != nil {
		zap.L().Warn("", zap.Error(err))
		return
	}
	zap.L().Info("", zap.Int("up migration applied:", affected))
}
func down(postDB *sql.DB) {
	affected, err := migrate.ExecMax(postDB, "postgres", migrations, migrate.Down, 1)
	if err != nil {
		zap.L().Warn("", zap.Error(err))
		return
	}
	zap.L().Info("", zap.Int("down migration applied:", affected))
}
