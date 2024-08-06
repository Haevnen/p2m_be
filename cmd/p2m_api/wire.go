package main

import (
	"context"

	"github.com/google/wire"
	"gitlab.citigo.com.vn/kship/kship-add-on/internal/app/voucher"
	"gitlab.citigo.com.vn/kship/kship-add-on/internal/pkg/configuration"
	"gitlab.citigo.com.vn/kship/kship-add-on/internal/pkg/connection/mysql"
	"gitlab.citigo.com.vn/kship/kship-add-on/internal/pkg/connection/redis"
	"gitlab.citigo.com.vn/kship/kship-add-on/internal/pkg/logging"
)

func NewAppConfig(configPath string) (*configuration.Config, error) {
	appConfig, err := configuration.NewConfig(configPath)
	if err != nil {
		return nil, err
	}
	// Init logging
	logging.InitializeLogging(appConfig.App.LogLevel, appConfig.App.LogMode, appConfig.App.LogFile, appConfig.App.LogFileAudit)
	return appConfig, nil
}

func initApp(configPath string, ctx context.Context) (*App, error) {
	wire.Build(
		// config
		NewAppConfig,
		// Redis
		redis.InitRedisClient,
		// mysql
		mysql.InitDb,

		// http server
		voucher.NewServer,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
