package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

type Request struct {
	Message        Message         `json:"message"`
	Authentication *Authentication `json:"authentication,omitempty"`
}

func (mailer *Mailer) SendMailRequest(request Request) {
	mailer.SendMail(request.Message, request.Authentication)
}

func (mailer *Mailer) SendMail(message Message, authentication *Authentication) error {

	httpClient := http.Client{Timeout: 5 * time.Second}

	// send url
	sendUrl, err := url.JoinPath(mailer.Url, "send")
	if err != nil {
		return err
	}

	// Request to send
	jsonRequest, err := json.Marshal(&Request{Message: message, Authentication: authentication})
	if err != nil {
		return err
	}

	// Sending
	mailerRequest, _ := http.NewRequest(http.MethodPost, sendUrl, bytes.NewBuffer(jsonRequest))
	mailerRequest.Header.Set("Content-Type", "application/json")
	_, err = httpClient.Do(mailerRequest)
	if err != nil {
		return err
	}

	return nil
}

func (mailer *Mailer) SendSecureMailRequest(user User, request Request) {
	mailer.SendSecureMail(user, request.Message, request.Authentication)
}

func (mailer *Mailer) SendSecureMail(user User, message Message, authentication *Authentication) error {

	httpClient := http.Client{Timeout: 5 * time.Second}

	// Request token
	tokenUrl, err := url.JoinPath(mailer.Url, "token")
	if err != nil {
		return err
	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return err
	}
	tokenRequest, err := http.NewRequest(http.MethodPost, tokenUrl, bytes.NewBuffer(jsonUser))
	tokenRequest.Header.Set("Content-Type", "application/json")
	tokenResponse, err := httpClient.Do(tokenRequest)
	if err != nil {
		return err
	}
	if tokenResponse.StatusCode != http.StatusCreated {
		return errors.New("retrieving access token")
	}
	var token map[string]interface{}
	json.NewDecoder(tokenResponse.Body).Decode(&token)

	// Send secure mail
	sendUrl, err := url.JoinPath(mailer.Url, "send")
	if err != nil {
		return err
	}
	jsonRequest, err := json.Marshal(&Request{Message: message, Authentication: authentication})
	if err != nil {
		return err
	}
	mailerRequest, err := http.NewRequest(http.MethodPost, sendUrl, bytes.NewBuffer(jsonRequest))
	mailerRequest.Header.Set("Content-Type", "application/json")
	mailerRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token["access_token"]))
	_, err = httpClient.Do(mailerRequest)
	if err != nil {
		return err
	}

	return nil
}
