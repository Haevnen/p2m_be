package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/configuration"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/handler"
	"github.com/Haevnen/p2m_be/internal/pkg/handler/middleware"
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/pkg/constants"
	"github.com/Haevnen/p2m_be/pkg/gormdb"
	"github.com/Haevnen/p2m_be/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Start() int {
	config, err := configuration.LoadConfig()
	if err != nil {
		fmt.Println("cannot load config: ", err)
		return 1
	}
	ctx := context.Background()
	gormConfig := config.GetGORMConfig()
	db, closeFunc, err := gormdb.New(gormConfig)
	if err != nil {
		logger.Fatal("cannot connect to db: ", err)
	}
	defer func() {
		if err := closeFunc(); err != nil {
			logger.Error(ctx, err, "close DB")
		}
	}()

	dal.SetDefault(db)

	reg := registry.New(config.TokenSymmetricKey)

	serverHandler := handler.New(reg)
	r := gin.Default()
	gin.SetMode(config.Mode)

	p2m_api.RegisterHandlersWithOptions(r, serverHandler, p2m_api.GinServerOptions{
		BaseURL: constants.BaseURL,
		Middlewares: []p2m_api.MiddlewareFunc{
			middleware.Authentication(reg.PasetoMaker()),
			middleware.Authorization(),
		},
	})

	s := &http.Server{
		Handler: r,
		Addr:    config.GetURLBase(),
	}

	go func() {
		// service connections
		fmt.Println("Starting server at " + config.GetURLBase())
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
	return 0
}

func main() {
	Start()
}
