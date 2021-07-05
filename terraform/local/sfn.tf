resource "aws_iam_role" "sfn_role" {
  name = local.role_sfn
  assume_role_policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AssumeRole",
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
                "Service": "states.amazonaws.com"
            }
        }
    ]
}
POLICY
}

resource "aws_sfn_state_machine" "state_machine" {
  name       = "SendMessageToSqsInFuture"
  role_arn = "arn:aws:iam::012345678901:role/DummyRole"
  definition = <<EOF
{
  "Comment": "Send message to SQS in the future",
  "StartAt": "Verify",
  "States": {
    "Verify": {
      "Comment": "The next state is going to verify de inputed data.",
      "Type": "Pass",
      "Next": "Check"
    },
    "Check": {
      "Comment": "Checks if payload, dueDate and queueUrl fields are present in the input",
      "Type": "Choice",
      "Default": "Success",
      "Choices": [
        {
          "Variable": "$.payload",
          "BooleanEquals": false,
          "Next": "Failed"
        },
        {
          "Variable": "$.dueDate",
          "BooleanEquals": false,
          "Next": "Failed"
        },
        {
          "Variable": "$.queueUrl",
          "BooleanEquals": false,
          "Next": "Failed"
        }
      ]
    },
    "Success": {
      "Type": "Pass",
      "Next": "WaitForDueDate"
    },
    "Failed": {
      "Type": "Fail",
      "Cause": "Invalid input, mandatory: payload, dueDate and queueUrl"
    },
    "WaitForDueDate": {
      "Comment": "Wait until the scheduled time",
      "Type": "Wait",
      "TimestampPath": "$.dueDate",
      "Next": "PushMessage"
    },
    "PushMessage": {
      "Type": "Task",
      "Resource": "arn:aws:states:::sqs:sendMessage",
      "Parameters": {
        "QueueUrl.$": "$.queueUrl",
        "MessageBody.$": "$.payload"
      },
      "Retry": [{
          "ErrorEquals": ["States.TaskFailed"],
          "IntervalSeconds": 30,
          "MaxAttempts": 10,
          "BackoffRate": 2.0
      }],
      "End": true
    }
  }
}
EOF

  tags = {
    env     = var.account
  }
}

//resource "aws_sfn_state_machine" "state_machine_lambda_example" {
//  name       = "SfnLambdaExample"
//  role_arn = "arn:aws:iam::012345678901:role/DummyRole"
//  definition = <<EOF
//  {
//  "Comment": "State Machine to calculate net run rates for given teams.",
//  "StartAt": "CreateMatches",
//  "States": {
//    "CreateMatches": {
//      "Type": "Task",
//      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:CreateMatchesFunction",
//      "Next": "CalculateRunRate"
//    },
//    "CalculateRunRate": {
//      "Type": "Task",
//      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:CalculateRunRateFunction",
//      "End": true
//    }
//  }
//}
//EOF
//
//  tags = {
//    env     = var.account
//  }
//}
