package main

import (
	"context"
	"errors"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "go.uber.org/automaxprocs"

	"github.com/Haevnen/p2m_be/internal/app/p2m_api"
	apiModel "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/configuration"
	"github.com/Haevnen/p2m_be/internal/di"
	"github.com/Haevnen/p2m_be/pkg/gormdb"
	"github.com/Haevnen/p2m_be/pkg/logger"
)

var config configuration.Config

func init() {
	viper.SetConfigFile(`configs/config.yml`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func main() {
	go func() {
		_ = http.ListenAndServe(config.GetURLProfile(), nil)
	}()

	ctx := context.Background()

	db, closeFunc, err := gormdb.New(config.GetGORMConfig())
	if err != nil {
		logger.Fatal(err, "init DB")
	}
	defer func() {
		if err := closeFunc(); err != nil {
			logger.Error(ctx, err, "close DB")
		}
	}()

	di := di.New(db)
	serverHandler := p2m_api.NewServer(di)
	r := gin.Default()
	gin.SetMode(config.App.Mode)
	apiModel.RegisterHandlers(r, serverHandler)

	// And we serve HTTP until the world ends.

	s := &http.Server{
		Handler: r,
		Addr:    config.GetURLBase(),
	}

	go func() {
		// service connections
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logger.Info("timeout of 5 seconds.")
	}
	logger.Info("Server exiting")
}
