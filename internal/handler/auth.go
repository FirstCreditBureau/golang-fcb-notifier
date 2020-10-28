package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/config"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/dto"
	"io/ioutil"
	"net/http"
)

// author zhasulan
// created on 27.10.20 14:00

type Authentication interface {
	getUserB64Auth() string
	Auth() (dto.LoginResult, error)
	Refresh(tokenRefresh dto.TokenRefresh) (dto.LoginResult, error)
}

type authentication struct {
	config *config.Configuration
}

func (a authentication) getUserB64Auth() string {
	str := a.config.Auth.Username + ":" + a.config.Auth.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(str))
}

func AuthenticationHandler(config *config.Configuration) Authentication {
	return &authentication{
		config: config,
	}
}

func (a authentication) Auth() (dto.LoginResult, error) {
	var loginResult dto.LoginResult
	request, err := http.NewRequest("POST", a.config.FCBNotifierHost+a.config.Auth.Login, nil)
	if err != nil {
		return loginResult, err
	}
	request.Header.Add("Authorization", a.getUserB64Auth())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return loginResult, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return loginResult, err
		}

		err = json.Unmarshal(bodyBytes, &loginResult)
		if err != nil {
			return loginResult, err
		}

		return loginResult, nil
	} else {
		return loginResult, errors.New("authentication server response error")
	}
}

func (a authentication) Refresh(tokenRefresh dto.TokenRefresh) (dto.LoginResult, error) {
	var loginResult dto.LoginResult

	tokenRefreshBuffer := new(bytes.Buffer)
	err := json.NewEncoder(tokenRefreshBuffer).Encode(tokenRefresh)
	if err != nil {
		return loginResult, err
	}

	request, err := http.NewRequest("POST", a.config.FCBNotifierHost+a.config.Auth.Refresh, tokenRefreshBuffer)
	if err != nil {
		return loginResult, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return loginResult, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return loginResult, err
		}

		err = json.Unmarshal(bodyBytes, &loginResult)
		if err != nil {
			return loginResult, err
		}

		return loginResult, nil
	} else {
		return loginResult, errors.New("authentication server response error")
	}
}
