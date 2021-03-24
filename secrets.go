package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SecretAdapter interface {
	get(string) (string, error)
}

// ParameterStore uses AWS System Manager Parameter Store to get secrets
type ParameterStore struct {
	svc *ssm.SSM
}

func NewParameterStore() (ParameterStore, error) {
	ret := ParameterStore{}

	sess, err := session.NewSessionWithOptions(session.Options{})
	if err != nil {
		return ret, err
	}

	svc := ssm.New(sess)
	ret.svc = svc
	return ret, nil
}

func (p ParameterStore) get(key string) (string, error) {
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
