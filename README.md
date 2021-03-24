# cloudwatch to slack

This package is an AWS Lambda function compatible tool which can send CloudWatch Events to slack.

## requirement

- Events should be sent via SNS
- Run on AWS Lambda

## Variables

### Environment Variables

- SLACK_NAME (default: "alert")

These variables should be stored in System Manager Parameter store

- SlackWebhook URL

Note: This can not set channel name, because Slack Incoming-Webhook by App should set one channel per webhook.