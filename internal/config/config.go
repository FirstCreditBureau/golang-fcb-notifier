package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/dto"
	"gopkg.in/validator.v2"
)

// author zhasulan
// created on 11.08.20 22:20

type Configuration struct {
	AppPort         string         `json:"app_port" validate:"nonzero"`
	AppName         string         `json:"app_name" validate:"nonzero"`
	Auth            dto.Auth       `json:"auth" validate:"nonzero"`
	Log             *logrus.Logger `json:"log"`
	FCBNotifierHost string         `json:"fcb_notifier_host"`
}

var Config = &Configuration{}

func GetLog() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
	})
	return log
}

func CheckValidity() error {
	err := validator.Validate(Config)
	if err != nil {
		return errors.Errorf("Config inconsistent: %v", err)
	}

	return nil
}
