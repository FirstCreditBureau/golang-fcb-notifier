package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/config"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/controllers"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/handler"
)

// author zhasulan
// created on 27.10.20 10:31

type App interface {
	Init()
	Run()
}

type app struct {
	config *config.Configuration
	router *gin.Engine
	log    *logrus.Logger
}

func AppHandler(config *config.Configuration) App {
	return &app{
		config: config,
		router: gin.New(),
		log:    config.Log,
	}
}

func (a *app) Init() {

	authentication := handler.AuthenticationHandler(a.config)
	auth, err := authentication.Auth()
	if err != nil {
		a.log.Error(err)
	}

	endpoint := controllers.EndpointHandler(a.config, auth)
	a.router.POST("/endpoint", endpoint.Endpoint)
}

func (a *app) Run() {
	if err := a.router.Run(":" + a.config.AppPort); err != nil {
		a.log.Fatal("Error running")
	}
}
