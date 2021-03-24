package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func LambdaHandler(context context.Context, event events.CloudWatchEvent) (int, error) {
	//	fmt.Printf("Detail = %s\n", event.Detail)
	j, _ := json.MarshalIndent(event, "", "  ")
	fmt.Printf("Source = %s\n", string(j))

	msgBody, err := GenerateMessage(event.Source, event.DetailType, event.Detail)
	if err != nil {
		return 0, fmt.Errorf("GenerateMessage, %w", err)
	}

	ss, err := NewSlackClient("Alert")
	if err != nil {
		return 0, fmt.Errorf("NewSlackClient, %w", err)
	}

	if err := ss.sendHttpRequest(msgBody); err != nil {
		return 0, fmt.Errorf("sendHttpRequest, %w", err)
	}

	return 0, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
