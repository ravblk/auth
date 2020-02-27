package config

import (
	"errors"
	"strings"

	"auth/storage/pg"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
)

var ErrConfig = errors.New("Wrong config")

type Service struct {
	API   API
	DB    pg.DBConfig
	Debug bool
}

type API struct {
	Port               string `valid:"required"`
	MaxRequestBodySize int    `valid:"required"`
	TTL                int    `valid:"required"`
}

func Read() (*Service, error) {
	sc := &Service{}
	svc := viper.New()

	if err := svc.ReadInConfig(); err != nil {
		return nil, err
	}
	svc.SetEnvPrefix("AUTH")
	svc.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//----------------
	svc.BindEnv("Api.Port")
	svc.BindEnv("Api.Maxrequestbodysize")
	svc.BindEnv("Api.TTL")
	//----------
	svc.BindEnv("DB.Host")
	svc.BindEnv("DB.Port")
	svc.BindEnv("DB.User")
	svc.BindEnv("DB.Password")
	svc.BindEnv("DB.DatabaseName")
	svc.BindEnv("DB.SSL")
	//------------
	svc.BindEnv("Debug")
	//--------------
	err := svc.Unmarshal(sc)
	if err != nil {
		return nil, err
	}
	if _, err = govalidator.ValidateStruct(sc); err != nil {
		return nil, err
	}
	return sc, nil
}
