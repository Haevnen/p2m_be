package main

import (
	"errors"
	"github.com/Haevnen/p2m_be/gen/api"

	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"

	"github.com/Haevnen/p2m_be/internal/app/p2m_api"
	"github.com/Haevnen/p2m_be/internal/configuration"
)

type App struct {
	appConfig *configuration.Config
	ginApp    *gin.Engine
}

func (a *App) StartProfile() {
	_ = http.ListenAndServe(a.appConfig.Server.Host+":"+a.appConfig.Server.Monitor, nil)
}

// StartGin start gin-gonic app
func (a *App) StartGin() error {
	server := p2m_api.NewServer()
	r := gin.Default()
	api.RegisterHandlers(r, server)

	// And we serve HTTP until the world ends.
	s := &http.Server{
		Handler: r,
		Addr:    a.appConfig.Server.Host + ":" + a.appConfig.Server.Port,
	}

	return s.ListenAndServe()
	if a.ginApp == nil {
		return errors.New("initialize ginApp failed")
	}
	return a.ginApp.Lio(a.appConfig.Server.Host + ":" + a.appConfig.Server.Port)
}
