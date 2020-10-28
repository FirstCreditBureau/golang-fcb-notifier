package main

import (
	"gitlab.com/golang-libs/mosk.git"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/config"
)

// author zhasulan
// created on 27.10.20 10:19

func main() {
	mosk.LoadLocalConfig(config.Config)
	log := config.GetLog()
	err := config.CheckValidity()
	if err != nil {
		log.Fatal(err)
	}
	conf := config.Config
	conf.Log = log

	app := internal.AppHandler(conf)
	app.Init()
	app.Run()
}
