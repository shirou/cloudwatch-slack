package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	EnvKeySlackName = "SLACK_NAME"
)

func LambdaHandler(context context.Context, event events.CloudWatchEvent) (int, error) {
	//	fmt.Printf("Detail = %s\n", event.Detail)
	j, _ := json.MarshalIndent(event, "", "  ")
	fmt.Printf("Source = %s\n", string(j))

	msgBody, err := GenerateMessage(event.Source, event.DetailType, event.Detail)
	if err != nil {
		return 0, fmt.Errorf("GenerateMessage, %w", err)
	}

	adapter, err := GetSecretAdapter()
	if err != nil {
		return 0, fmt.Errorf("GetSecretAdapter, %w", err)
	}

	ss, err := NewSlackClient(GetSlackName(), adapter)
	if err != nil {
		return 0, fmt.Errorf("NewSlackClient, %w", err)
	}

	if err := ss.sendHttpRequest(msgBody); err != nil {
		return 0, fmt.Errorf("sendHttpRequest, %w", err)
	}

	return 0, nil
}

func GetSlackName() string {
	p := os.Getenv(EnvKeySlackName)
	if p == "" {
		return "Alert"
	}
	return p
}

func main() {
	lambda.Start(LambdaHandler)
}
