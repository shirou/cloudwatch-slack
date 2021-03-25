package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	EnvKeyWebhookURL = "SLACK_WEBHOOK_URL"
)

type SecretAdapter interface {
	get(string) (string, error)
}

func GetSecretAdapter() (SecretAdapter, error) {
	p := os.Getenv(EnvKeyWebhookURL)
	// if environmet variable has slack url, use it
	if strings.HasPrefix(p, "https://hooks.slack.com") {
		return NewEnvStore()
	}

	// if starts from "ssm:", use Parameter Store
	if strings.HasPrefix(p, "ssm:") {
		key := strings.Replace(p, "ssm:", "", 1)
		return NewParameterStore(key)
	}

	return nil, fmt.Errorf("can not find Secret Adapter, %s. please set %s", p, EnvKeyWebhookURL)
}

// ParameterStore uses AWS System Manager Parameter Store to get secrets
type ParameterStore struct {
	svc *ssm.SSM

	webhook_key string
}

// NewParameterStore returns ParameterStore SecretAdapter.
// This returns secret from AWS Parameter Store.
func NewParameterStore(webhook_key string) (ParameterStore, error) {
	sess, err := session.NewSessionWithOptions(session.Options{})
	if err != nil {
		return ParameterStore{}, err
	}

	svc := ssm.New(sess)
	ret := ParameterStore{
		svc:         svc,
		webhook_key: webhook_key,
	}
	return ret, nil
}

func (p ParameterStore) get(key string) (string, error) {
	if key == EnvKeyWebhookURL {
		key = p.webhook_key
	}

	res, err := p.svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", nil
	}

	return *res.Parameter.Value, nil

}

// EnvStore uses Environment Variables to get secrets, not recommended
type EnvStore struct {
}

func NewEnvStore() (EnvStore, error) {
	return EnvStore{}, nil
}

func (p EnvStore) get(key string) (string, error) {
	return os.Getenv(key), nil
}
