package cmd

import (
	"auth/config"
	repository "auth/repositories"
	"auth/server/transport/http"
	"auth/server/transport/http/handlers"
	"auth/services"
	"auth/services/users"
	"auth/storage/pg"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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
	atom := zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	log := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	undo := zap.ReplaceGlobals(log)
	defer undo()

	cfg, err := config.Read()
	if err != nil {
		log.Fatal("", zap.Error(err))
	}
	if cfg.Debug {
		atom.SetLevel(zap.DebugLevel)
	}
	db := pg.New(log)
	if err := db.Connect(cfg.DB); err != nil {
		log.Fatal("", zap.Error(err))
	}
	sessRep := repository.NewSessionRepository(log, db)
	usrRep := repository.NewUserRepository(log, db)

	usrSvc := users.NewService(log, &cfg.API, usrRep, sessRep)

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
