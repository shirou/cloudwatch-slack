# cloudwatch to slack

This package is an AWS Lambda function compatible tool which can send CloudWatch Events to slack.

## requirement

- Events should be sent via EventBridge, not from.
- Run on AWS Lambda

## Variables

### Environment Variable

- SLACK_WEBHOOK_URL

You can specify a Slack webhook URL to `SLACK_WEBHOOK_URL`, like `https://hooks.slack.com/services/A00000/B0000000/dddddddd`. And if set `ssm:foo_bar`, starts with `ssm:`, get webhook url from AWS Parameter store with key(this example get from `foo_bar`).

- SLACK_TEMPLATE_DIR (default: <binary_path>/templates/)

You can specify message template directory by this environment variable.



## Limitations

Note: This can not set channel name, because Slack Incoming-Webhook by App should set one channel per webhook.

## License

Apache 2