package config_test

import (
	"gitlab.com/golang-libs/mosk.git"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/config"
	"testing"
)

func TestCheckValidity(t *testing.T) {
	log := config.GetLog()
	mosk.LoadConfig("../../config/conf.json", config.Config)
	err := config.CheckValidity()
	if err != nil {
		log.Error(err)
		t.FailNow()
	}
}
