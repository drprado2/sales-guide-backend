resource "aws_sqs_queue" "test_queue" {
  name = "test_queue"
  delay_seconds = 0
  max_message_size = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 0
  visibility_timeout_seconds = 4
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.test_queue_deadletter.arn
    maxReceiveCount = 4
  })

  tags = {
    Environment = "develop"
  }
}

resource "aws_sqs_queue" "test_queue_deadletter" {
  name = "test_queue_dlq"
  delay_seconds = 90
  max_message_size = 2048
  message_retention_seconds = 604800
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 10

  tags = {
    Environment = "develop"
  }
}
