resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "hello_world" {
  function_name = "HelloWorld"
  handler          = "hello_world"
  memory_size      = 128
  timeout          = 2

  # The bucket name as created earlier with "aws s3api create-bucket"
  s3_bucket = var.s3_bucket_name
  s3_key    = "v1.0.2/hello_wolrd.zip"

  runtime = "go1.x"

  role          = aws_iam_role.iam_for_lambda.arn
}

resource "aws_lambda_function" "good_by" {
  function_name = "GoodBy"
  handler          = "good_by"
  memory_size      = 128
  timeout          = 2

  # The bucket name as created earlier with "aws s3api create-bucket"
  s3_bucket = var.s3_bucket_name
  s3_key    = "v1.0.0/good_by.zip"

  runtime = "go1.x"

  role          = aws_iam_role.iam_for_lambda.arn
}
