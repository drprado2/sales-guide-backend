resource "aws_iam_user" "user" {
  name = local.user_name
}

resource "aws_iam_policy_attachment" "test_queue_policy_attach" {
  name = "test_queue_policy_attachment"
  users = [aws_iam_user.user.name]
  policy_arn = aws_iam_policy.test_queue_policy.arn
}

resource "aws_iam_policy" "test_queue_policy" {
name = "test_queue_policy"

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
        "${aws_sqs_queue.test_queue.arn}",
        "${aws_sqs_queue.test_queue_deadletter.arn}"
      ]
    }
  ]
}
POLICY
}

resource "aws_iam_policy" "sfn_policy" {
  name = "api-sales_guide-stf-policy"

  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "stateMachine",
            "Effect": "Allow",
            "Action": [
                "states:StartExecution"
            ],
            "Resource": [
                "${aws_sfn_state_machine.state_machine.arn}"
            ]
        }
    ]
}
POLICY
}

resource "aws_iam_policy_attachment" "sales_guide_sfn_role_policy_attach" {
  name = "sales_guide_sfn_role_policy-attachment"
  users = [aws_iam_user.user.name]
  roles = [local.role_sfn]
  policy_arn = aws_iam_policy.sfn_policy.arn
}
