# cloudwatch to slack

This package is an AWS Lambda function compatible tool which can send CloudWatch Events to slack.

## Requirement

- Events should be sent via EventBridge, not from SNS.
- Run on AWS Lambda

## Variables

### Environment Variable

- SLACK_WEBHOOK_URL

You can specify a Slack webhook URL to `SLACK_WEBHOOK_URL`, like `https://hooks.slack.com/services/A00000/B0000000/dddddddd`. And if set `ssm:foo_bar`, starts with `ssm:`, get webhook url from AWS Parameter store with key(this example get from `foo_bar`).

- SLACK_TEMPLATE_DIR (default: <binary_path>/templates/)

You can specify message template directory by this environment variable.

## Limitations

We can not set channel name. Because Slack Incoming-Webhook by App should set one channel per webhook.



## Terraform sample to add Cloudwatch Event Rule

```
### Event Watch Trigger
resource "aws_cloudwatch_event_rule" "cloudwatch_alarm_event" {
  name          = "cloudwatch_alarm_event"
  description   = "Fires if Cloudwatch Alarm state changed"
  event_pattern = <<EOT
{
  "detail-type": [
    "CloudWatch Alarm State Change"
  ]
}
EOT
}

resource "aws_cloudwatch_event_target" "cloudwatch_alarm_event" {
  rule      = aws_cloudwatch_event_rule.cloudwatch_alarm_event.name
  target_id = aws_cloudwatch_event_rule.cloudwatch_alarm_event.id
  arn       = aws_lambda_function.cloudwatch_slack.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_slack" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.cloudwatch_slack.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.cloudwatch_alarm_event.arn
}
```

## License

Apache 2