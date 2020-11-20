package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/config"
	"gitlabce.1cb.kz/notifier/golang-fcb-notifier/internal/dto"
	"io/ioutil"
	"net/http"
	"time"
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
		return loginResult, fmt.Errorf("authentication server response error. Code: %d", response.StatusCode)
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
		return loginResult, fmt.Errorf("authentication server response error. Code: %d", response.StatusCode)
	}
}

func GetBearerToken(config *config.Configuration, auth dto.LoginResult) (string, error) {
	if time.Now().After(auth.Access.ExpiresAt) || auth == (dto.LoginResult{}) {
		if time.Now().After(auth.Refresh.ExpiresAt) || auth == (dto.LoginResult{}) {
			authentication := AuthenticationHandler(config)
			a, err := authentication.Auth()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Bearer %s", a.Access.Hash), nil
		} else {
			authentication := AuthenticationHandler(config)
			a, err := authentication.Refresh(dto.TokenRefresh{TokenHash: auth.Refresh.Hash})
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Bearer %s", a.Access.Hash), nil
		}
	}

	return fmt.Sprintf("Bearer %s", auth.Access.Hash), nil
}
