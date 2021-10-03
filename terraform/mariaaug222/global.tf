resource "aws_iam_user" "user" {
  name = local.user_name
}

resource "aws_iam_policy_attachment" "app_updated_queue_policy_attach" {
  name = "app_updated_queue_policy_attach"
  users = [aws_iam_user.user.name]
  policy_arn = aws_iam_policy.test_queue_policy.arn
}

resource "aws_iam_policy" "test_queue_policy" {
  name = "sales-guide-sqs-policy"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "sqs",
      "Effect": "Allow",
      "Action": [
                "sqs:ReceiveMessage",
                "sqs:DeleteMessage",
                "sqs:GetQueueAttributes",
                "sqs:SendMessage",
                "sqs:SendMessageBatch"
            ],
      "Resource": [
        "${aws_sqs_queue.app_updated_queue.arn}",
        "${aws_sqs_queue.app_updated_deadletter.arn}"
      ]
    }
  ]
}
POLICY
}
