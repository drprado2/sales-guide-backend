resource "aws_sqs_queue" "app_updated_queue" {
  name = local.app_updated_queue
  delay_seconds = 0
  max_message_size = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 0
  visibility_timeout_seconds = 4
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.app_updated_deadletter.arn
    maxReceiveCount = 2
  })

  tags = {
    Environment = "sales-guide"
  }
}

resource "aws_sqs_queue" "app_updated_deadletter" {
  name = local.app_updated_dlq
  delay_seconds = 0
  max_message_size = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 0
  visibility_timeout_seconds = 4

  tags = {
    Environment = "sales-guide"
  }
}
