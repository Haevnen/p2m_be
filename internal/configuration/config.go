package configuration

import (
	"fmt"

	"github.com/Haevnen/p2m_be/pkg/gormdb"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DBPort             int    `env:"DB_PORT"`
	DBHost             string `env:"DB_HOST"`
	MySQLRootPassword  string `env:"MYSQL_ROOT_PASSWORD"`
	MySQLUser          string `env:"MYSQL_USER"`
	MySQLPassword      string `env:"MYSQL_PASSWORD"`
	MySQLDatabase      string `env:"MYSQL_DATABASE"`
	MySQLContainerName string `env:"MYSQL_CONTAINER_NAME"`

	ProjectName           string `env:"PROJECT_NAME"`
	AppAPIDir             string `env:"APP_API_DIR"`
	SpecDir               string `env:"SPEC_DIR"`
	DockerBuildTagVersion string `env:"DOCKER_BUILD_TAG_VERSION"`

	APIPort int    `env:"API_PORT"`
	APIHost string `env:"API_HOST"`
	Mode    string `env:"RUN_MODE"`
}

func LoadConfig() (config Config, err error) {
	// Load the .env file
	err = godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}

	err = env.Parse(&config)
	if err != nil {
		panic(err)
	}

	return
}

func (e *Config) GetGORMConfig() *gormdb.Config {
	return &gormdb.Config{
		DBHost:            e.DBHost,
		DBPort:            e.DBPort,
		DBUser:            e.MySQLUser,
		DBPass:            e.MySQLPassword,
		DBName:            e.MySQLDatabase,
		LogSQL:            true,
		MaxOpenConn:       10,
		MaxLifetimeSecond: 300,
	}
}

// GetURLBase build server config from env
func (e *Config) GetURLBase() string {
	return fmt.Sprintf("%s:%d", e.APIHost, e.APIPort)
}
