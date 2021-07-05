resource "aws_sns_topic" "sns-test" {
  name = "sns-test"
}

resource "aws_sns_topic_subscription" "sns_test_sqs_test" {
  topic_arn = aws_sns_topic.sns-test.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.test_queue.arn
}

resource "aws_sns_topic_subscription" "sns_test_email_maria" {
  topic_arn = aws_sns_topic.sns-test.arn
  protocol  = "email"
  endpoint  = "mariaaug222@gmail.com"
}
