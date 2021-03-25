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

func LambdaHandler(context context.Context, event events.CloudWatchEvent) error {
	//	fmt.Printf("Detail = %s\n", event.Detail)
	j, _ := json.MarshalIndent(event, "", "  ")
	fmt.Printf("Source = %s\n", string(j))

	msgBody, err := GenerateMessage(event.Source, event.DetailType, event.Detail)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("GenerateMessage, %w", err)
	}

	adapter, err := GetSecretAdapter()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("GetSecretAdapter, %w", err)
	}

	ss, err := NewSlackClient(GetSlackName(), adapter)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("NewSlackClient, %w", err)
	}

	if err := ss.sendHttpRequest(msgBody); err != nil {
		fmt.Println(err)
		return fmt.Errorf("sendHttpRequest, %w", err)
	}

	return nil
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
