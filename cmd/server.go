package cmd

import (
	"auth/internal/config"
	repository "auth/internal/repositories"
	"auth/internal/server/transport/http"
	"auth/internal/server/transport/http/handlers"
	"auth/internal/services"
	"auth/internal/services/users"
	"auth/internal/storage/pg"

	"go.uber.org/zap"

	_ "github.com/lib/pq"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server auth",
	Long:  `auth server with services for registration, authorization, authentication`,
	Run: func(cmd *cobra.Command, args []string) {
		RunServer()
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func RunServer() {
	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("", zap.Error(err))
	}
	undo := zap.ReplaceGlobals(log)
	defer undo()

	cfg, err := config.Read()
	if err != nil {
		log.Fatal("", zap.Error(err))
	}
	if cfg.Debug {
		zap.L().Sync()
		undo()
		log = zap.NewNop()
		undo = zap.ReplaceGlobals(log)
	}
	db := pg.New(log)
	if err := db.Connect(cfg.DB); err != nil {
		log.Fatal("", zap.Error(err))
	}
	sessionRepo := repository.NewSessionRepository(log, db)
	userRepo := repository.NewUserRepository(log, db)

	usrSvc := users.NewService(log, &cfg.API, userRepo, sessionRepo)

	authsvc := services.Auth{
		UsrSvc: usrSvc,
		Log:    log,
	}

	hs := handlers.New(authsvc)
	s, err := http.NewServer(hs, cfg.API.MaxRequestBodySize)
	if err != nil {
		log.Fatal("", zap.Error(err))
	}

	s.RoutesInit()
	log.Info("server started")
	if err := s.Run(cfg.API.Port); err != nil {
		log.Fatal("", zap.Error(err))
	}
}
