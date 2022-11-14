package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Mailer struct {
	Url string
}

type Message struct {
	From    string `json:"from"`
	To      string `json:"to" example:"me@gmail.com,you@gmail.com"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Authentication struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type User struct {
	Username string
	Password string
}

type Info struct {
	Version string
	Secure  bool
}

type Request struct {
	Message        Message         `json:"message"`
	Authentication *Authentication `json:"authentication,omitempty"`
}

func (mailer *Mailer) Info() (*Info, error) {

	httpClient := http.Client{Timeout: 5 * time.Second}

	// info url
	infoUrl, err := url.JoinPath(mailer.Url, "info")
	if err != nil {
		return nil, err
	}

	// request info data
	infoRequest, err := http.NewRequest(http.MethodGet, infoUrl, nil)
	infoRequest.Header.Set("Content-Type", "application/json")
	infoResponse, err := httpClient.Do(infoRequest)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(infoResponse.Body)
	if infoResponse.StatusCode != http.StatusOK {
		return nil, errors.New("retrieving mailer info")
	}

	// Info struct
	i := &Info{}
	err = json.NewDecoder(infoResponse.Body).Decode(&i)
	if err != nil {
		return nil, err
	}

	return i, nil

}

func (mailer *Mailer) IsSecure() (bool, error) {

	info, err := mailer.Info()

	if err != nil {
		return false, err
	}

	return info.Secure, nil
}

func (mailer *Mailer) SendMailRequest(request Request) error {
	return mailer.SendMail(request.Message, request.Authentication)
}

func (mailer *Mailer) SendMail(message Message, authentication *Authentication) error {

	httpClient := http.Client{Timeout: 5 * time.Second}

	// send url
	sendUrl, err := url.JoinPath(mailer.Url, "send")
	if err != nil {
		return err
	}

	// convert message request to json
	jsonRequest, err := json.Marshal(&Request{Message: message, Authentication: authentication})
	if err != nil {
		return err
	}

	// send mail
	mailerRequest, _ := http.NewRequest(http.MethodPost, sendUrl, bytes.NewBuffer(jsonRequest))
	mailerRequest.Header.Set("Content-Type", "application/json")
	_, err = httpClient.Do(mailerRequest)
	if err != nil {
		return err
	}

	return nil
}

func (mailer *Mailer) SendSecureMailRequest(user User, request Request) error {
	return mailer.SendSecureMail(user, request.Message, request.Authentication)
}

func (mailer *Mailer) SendSecureMail(user User, message Message, authentication *Authentication) error {

	httpClient := http.Client{Timeout: 5 * time.Second}

	// token url
	tokenUrl, err := url.JoinPath(mailer.Url, "token")
	if err != nil {
		return err
	}

	// user
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// request user token
	tokenRequest, err := http.NewRequest(http.MethodPost, tokenUrl, bytes.NewBuffer(jsonUser))
	tokenRequest.Header.Set("Content-Type", "application/json")
	tokenResponse, err := httpClient.Do(tokenRequest)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(tokenResponse.Body)
	if tokenResponse.StatusCode != http.StatusCreated {
		return errors.New("retrieving access token")
	}

	// get token
	var token map[string]interface{}
	err = json.NewDecoder(tokenResponse.Body).Decode(&token)
	if err != nil {
		return err
	}

	// sen url
	sendUrl, err := url.JoinPath(mailer.Url, "send")
	if err != nil {
		return err
	}

	// convert message request to json
	jsonRequest, err := json.Marshal(&Request{Message: message, Authentication: authentication})
	if err != nil {
		return err
	}

	// send secure mail
	mailerRequest, err := http.NewRequest(http.MethodPost, sendUrl, bytes.NewBuffer(jsonRequest))
	mailerRequest.Header.Set("Content-Type", "application/json")
	mailerRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token["access_token"]))
	_, err = httpClient.Do(mailerRequest)
	if err != nil {
		return err
	}

	return nil
}
