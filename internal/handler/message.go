package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/config"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/dto"
	"io/ioutil"
	"net/http"
)

// author zhasulan
// created on 27.10.20 14:29

type Message interface {
	GetFile() ([]byte, error)
}

type message struct {
	config          *config.Configuration
	ProxyURL        string
	cdnProxyRequest dto.CDNProxyRequest
	authentication  dto.LoginResult
}

func MessageHandler(config *config.Configuration, ProxyURL string, cdnProxyRequest dto.CDNProxyRequest, auth dto.LoginResult) Message {
	return &message{
		config:          config,
		ProxyURL:        ProxyURL,
		cdnProxyRequest: cdnProxyRequest,
		authentication:  auth,
	}
}

func (a message) GetFile() ([]byte, error) {
	cdnProxyRequestBuffer := new(bytes.Buffer)
	err := json.NewEncoder(cdnProxyRequestBuffer).Encode(a.cdnProxyRequest)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", a.ProxyURL, cdnProxyRequestBuffer)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.authentication.Access.Hash))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return bodyBytes, nil
	} else {
		return nil, errors.New(fmt.Sprintf("CDN proxy response error. Code: %d", response.StatusCode))
	}
}
