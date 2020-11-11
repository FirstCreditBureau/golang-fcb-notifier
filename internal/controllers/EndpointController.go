package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/config"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/dto"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/handler"
	"net/http"
	"time"
)

// author zhasulan
// created on 27.10.20 13:04

type EndpointController interface {
	Endpoint(context *gin.Context)
}

type endpointController struct {
	config         *config.Configuration
	authentication dto.LoginResult
}

func EndpointHandler(config *config.Configuration, auth dto.LoginResult) EndpointController {
	return &endpointController{
		config:         config,
		authentication: auth,
	}
}

func (controller *endpointController) Endpoint(context *gin.Context) {
	var notifierRequest dto.NotifierRequest
	if err := context.ShouldBindJSON(&notifierRequest); err != nil {
		context.JSON(http.StatusBadRequest, dto.Response{Code: http.StatusBadRequest, Message: err.Error()})
		_ = context.Error(err)
		return
	}

	// check token expire
	if time.Now().After(controller.authentication.Access.ExpiresAt) {
		authentication := handler.AuthenticationHandler(controller.config)
		auth, err := authentication.Auth()
		if err != nil {
			context.JSON(http.StatusInternalServerError, dto.Response{Code: http.StatusInternalServerError, Message: err.Error()})
			_ = context.Error(err)
			return
		}
		controller.authentication = auth
	}

	message := handler.MessageHandler(controller.config, notifierRequest.ProxyURL, dto.CDNProxyRequest{Code: notifierRequest.Code, Filename: notifierRequest.Filename}, controller.authentication)
	content, err := message.GetFile()
	if err != nil {
		context.JSON(http.StatusInternalServerError, dto.Response{Code: http.StatusInternalServerError, Message: err.Error()})
		_ = context.Error(err)
		return
	}

	//content_reader := bytes.NewReader(content)
	//sha := sha256.New()
	//r := io.TeeReader(content_reader, sha)
	//if _, err = io.Copy(content_reader, r); err != nil {
	//	return "", err
	//}

	enc := sha256.Sum256(content)
	checksum := hex.EncodeToString(enc[:])
	if notifierRequest.Checksum != string(checksum) {
		err := errors.New("checksum not equal to calculated hash sum from requested file")
		context.JSON(http.StatusInternalServerError, dto.Response{Code: http.StatusInternalServerError, Message: err.Error()})
		_ = context.Error(err)
		return
	}

	// todo check RSA

	mainHandler := handler.CreateMainHandler(notifierRequest.Code, content, controller.config.Log)
	go mainHandler.Handler()

}
