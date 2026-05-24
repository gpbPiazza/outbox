package envs

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/caarlos0/env/v11"
)

type Env struct {
	AppName string `env:"NAME" envDefault:"wes-api"`
	DB      Database
}

type Database struct {
	// DatabaseURL example: "postgres://username:password@localhost:5432/database_name"
	DatabaseURL         string        `env:"DATABASE_URL,required"`
	DBMateMigrationsDir string        `env:"DBMATE_MIGRATIONS_DIR" envDefault:"./infra/databases/sql/migrations"`
	MaxConns            int           `env:"DATABASE_POOL_MAX_CONNS" envDefault:"25"`
	MinConns            int           `env:"DATABASE_POOL_MIN_CONNS" envDefault:"25"`
	MaxConnsLifeTime    time.Duration `env:"DATABASE_POOL_MAX_CONN_LIFETIME" envDefault:"30m"`
}

func Envs() Env {
	valRef := reflect.ValueOf(globalEnv)

	if valRef.IsZero() {
		if err := initEnvs(); err != nil {
			log.Fatal(err)
		}
	}

	return globalEnv
}

var globalEnv Env

func initEnvs() error {
	var config Env

	if err := env.Parse(&config); err != nil {
		return fmt.Errorf("error on init envs variables - %s", err)
	}

	globalEnv = config
	return nil
}
