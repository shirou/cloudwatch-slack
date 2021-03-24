package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const DefaultSlackTimeout = 5 * time.Second

type SlackClient struct {
	WebHookUrl string
	UserName   string
	Channel    string
	TimeOut    time.Duration
}

func NewSlackClient(username string) (SlackClient, error) {
	/*
		ss, err := NewParameterStore()
		if err != nil {
			return SlackClient{}, fmt.Errorf("NewParameterStore, %w", err)
		}
		webHookUrl, err := ss.get("webhook")
		if err != nil {
			return SlackClient{}, fmt.Errorf("get webhook url, %w", err)
		}
	*/
	ss, _ := NewEnvStore()
	webHookUrl, _ := ss.get("webhook")

	return SlackClient{
		WebHookUrl: webHookUrl,
		UserName:   username,
		TimeOut:    DefaultSlackTimeout,
	}, nil
}

func (sc SlackClient) sendHttpRequest(msgBody string) error {
	req, err := http.NewRequest(http.MethodPost, sc.WebHookUrl, bytes.NewBuffer([]byte(msgBody)))
	if err != nil {
		return fmt.Errorf("sendHttpRequest, %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	if sc.TimeOut == 0 {
		sc.TimeOut = DefaultSlackTimeout
	}
	client := &http.Client{Timeout: sc.TimeOut}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if buf.String() != "ok" {
		return fmt.Errorf("Non-ok response returned from Slack, %s", buf.String())
	}
	return nil
}
