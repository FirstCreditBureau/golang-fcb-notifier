package handler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// author zhasulan
// created on 28.10.20 16:40

type MainHandler interface {
	Handler()
}

type mainHandler struct {
	code    string
	content []byte
	log     *logrus.Logger
}

func CreateMainHandler(code string, content []byte, log *logrus.Logger) MainHandler {
	return &mainHandler{
		code:    code,
		content: content,
		log:     log,
	}
}

func (handler mainHandler) Handler() {
	time.Sleep(5 * time.Second)
	/*
	   Здесь вы должны писать собственный обработчик

	   Каждый code - это отдельный сервис ПКБ и каждый из них должен обрабатываться отдельно или высылаться в другую
	   программу

	   if code == 'MON':  # Банковский мониторинг
	       mon_handler(content)
	   if code == 'BCO':  # Бэкофис ПКБ
	       bco_handler(content)

	   и т.д.
	*/
	handler.log.Info(fmt.Sprintf("Сообщение обработано. Код: %s; Контент: %s", handler.code, string(handler.content)))
	time.Sleep(5 * time.Second)
}
