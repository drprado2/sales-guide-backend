resource "aws_api_gateway_rest_api" "apiGateway" {
  name = "${var.environment}-s3-proxy-github-example"
  description = "${var.environment} s3 proxy"
  binary_media_types = "${var.supported_binary_media_types}"
}

resource "aws_api_gateway_deployment" "s3-proxy-api-deployment-example" {
  depends_on = [
    "aws_api_gateway_integration.itemPutMethod-ApiProxyIntegration",
    "aws_api_gateway_integration.itemGetMethod-ApiProxyIntegration",
    "aws_api_gateway_integration.itemOptionsMethod-ApiProxyIntegration",
  ]

  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"

  stage_name = "${var.environment}"
}

resource "aws_api_gateway_usage_plan" "s3-proxy-usage-plan" {
  name = "s3-proxy-usage-plan-github-example-${var.environment}"
  description = "usage plan for s3 proxy"

  api_stages {
    api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
    stage = "${aws_api_gateway_deployment.s3-proxy-api-deployment-example.stage_name}"
  }
}

resource "aws_api_gateway_api_key" "s3-api-key" {
  name = "s3-proxy-gitgub-example-apikey-${var.environment}"
}

resource "aws_api_gateway_usage_plan_key" "s3-proxy-usage-plan-key" {
  key_id = "${aws_api_gateway_api_key.s3-api-key.id}"
  key_type = "API_KEY"
  usage_plan_id = "${aws_api_gateway_usage_plan.s3-proxy-usage-plan.id}"
}
#iam role
resource "aws_iam_role" "s3_proxy_role" {
  name = "${var.environment}-s3-proxy-role-example"
  path = "/"
  assume_role_policy = "${data.aws_iam_policy_document.s3_proxy_policy.json}"
}

data "aws_iam_policy_document" "s3_proxy_policy" {
  statement {
    actions = [
      "sts:AssumeRole"]

    principals {
      type = "Service"
      identifiers = [
        "apigateway.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "s3_proxy_role_file_upload_attachment" {
  depends_on = [
    "aws_iam_policy.s3_file_upload_policy",
  ]

  role = "${aws_iam_role.s3_proxy_role.name}"
  policy_arn = "${aws_iam_policy.s3_file_upload_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "s3_proxy_role_api_gateway_attachment" {
  depends_on = [
    "aws_iam_policy.s3_file_upload_policy",
  ]

  role = "${aws_iam_role.s3_proxy_role.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess"
}

resource "aws_s3_bucket" "file_upload_bucket" {
  bucket = "file-upload-bucket-${var.environment}"
  acl = "private"

  tags = {
    Name = "file-upload-bucket-${var.environment}"
    Environment = "${var.environment}"
  }

  depends_on = [
    "aws_iam_policy.s3_file_upload_policy",
  ]
}

resource "aws_iam_policy" "s3_file_upload_policy" {
  name = "${var.environment}-github-s3-file-upload-policy"
  path = "/"
  description = "${var.environment} s3 file upload policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
                "s3:PutObject",
                "s3:GetObject"
            ],
      "Effect": "Allow",
      "Resource": [
                "arn:aws:s3:::file-upload-bucket-${var.environment}/*"
            ]
    }
  ]
}
EOF
}

#api resources
resource "aws_api_gateway_resource" "bucketResource" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  parent_id = "${aws_api_gateway_rest_api.apiGateway.root_resource_id}"
  path_part = "{bucket}"
}

resource "aws_api_gateway_resource" "folderResource" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  parent_id = "${aws_api_gateway_resource.bucketResource.id}"
  path_part = "{folder}"
}

resource "aws_api_gateway_resource" "itemResource" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  parent_id = "${aws_api_gateway_resource.folderResource.id}"
  path_part = "{item}"
}

#proxy integration
resource "aws_api_gateway_integration" "itemPutMethod-ApiProxyIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "${aws_api_gateway_method.itemPutMethod.http_method}"

  type = "AWS"
  integration_http_method = "PUT"
  credentials = "${aws_iam_role.s3_proxy_role.arn}"
  uri = "arn:aws:apigateway:${var.aws_region}:s3:path/{bucket}/{folder}/{item}"

  request_parameters = {
    "integration.request.header.x-amz-meta-fileinfo" = "method.request.header.x-amz-meta-fileinfo"
    "integration.request.header.Accept" = "method.request.header.Accept"
    "integration.request.header.Content-Type" = "method.request.header.Content-Type"

    "integration.request.path.item" = "method.request.path.item"
    "integration.request.path.folder" = "method.request.path.folder"
    "integration.request.path.bucket" = "method.request.path.bucket"
  }
}

resource "aws_api_gateway_integration" "itemGetMethod-ApiProxyIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "${aws_api_gateway_method.itemGetMethod.http_method}"

  type = "AWS"
  integration_http_method = "GET"
  credentials = "${aws_iam_role.s3_proxy_role.arn}"
  uri = "arn:aws:apigateway:${var.aws_region}:s3:path/{bucket}/{folder}/{item}"

  request_parameters = {
    "integration.request.path.item" = "method.request.path.item"
    "integration.request.path.folder" = "method.request.path.folder"
    "integration.request.path.bucket" = "method.request.path.bucket"
  }
}

resource "aws_api_gateway_integration" "itemOptionsMethod-ApiProxyIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "${aws_api_gateway_method.itemOptionsMethod.http_method}"
  type = "MOCK"
  depends_on = [
    "aws_api_gateway_method.itemOptionsMethod"]

  request_templates = {
    "application/json" = <<EOF
        {
        "statusCode" : 200
        }
    EOF
  }
}

#api methods
resource "aws_api_gateway_method" "itemPutMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "PUT"
  authorization = "NONE"
  api_key_required = true

  request_parameters = {
    "method.request.header.Accept" = false
    "method.request.header.Content-Type" = false
    "method.request.header.x-amz-meta-fileinfo" = false

    "method.request.path.bucket" = true
    "method.request.path.folder" = true
    "method.request.path.item" = true
  }
}

resource "aws_api_gateway_method" "itemOptionsMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "OPTIONS"
  authorization = "NONE"

  request_parameters = {
    "method.request.header.x-amz-meta-fileinfo" = false
  }
}

resource "aws_api_gateway_method" "itemGetMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "GET"
  authorization = "NONE"
  api_key_required = true

  request_parameters = {
    "method.request.path.bucket" = true
    "method.request.path.folder" = true
    "method.request.path.item" = true
  }
}

#method responses
resource "aws_api_gateway_method_response" "itemPutMethod200Response" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "${aws_api_gateway_method.itemPutMethod.http_method}"
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = true
  }

  depends_on = [
    "aws_api_gateway_method.itemPutMethod"]
}

resource "aws_api_gateway_method_response" "itemOptionsMethod200Response" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "${aws_api_gateway_method.itemOptionsMethod.http_method}"
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Origin" = true
  }

  depends_on = [
    "aws_api_gateway_method.itemOptionsMethod"]
}

resource "aws_api_gateway_method_response" "itemGetMethod200Response" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "${aws_api_gateway_method.itemGetMethod.http_method}"
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = true
  }

  depends_on = [
    "aws_api_gateway_method.itemGetMethod"]
}

# integration responses
resource "aws_api_gateway_integration_response" "itemPutMethod-IntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "${aws_api_gateway_method.itemPutMethod.http_method}"

  status_code = "${aws_api_gateway_method_response.itemPutMethod200Response.status_code}"

  response_templates = {
    "application/json" = ""
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = "'*'"
  }
}

resource "aws_api_gateway_integration_response" "itemGetMethod-IntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "${aws_api_gateway_method.itemGetMethod.http_method}"

  status_code = "${aws_api_gateway_method_response.itemGetMethod200Response.status_code}"

  response_templates = {
    "application/json" = ""
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = "'*'"
  }
}

resource "aws_api_gateway_integration_response" "itemOptionsMethod-IntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.apiGateway.id}"
  resource_id = "${aws_api_gateway_resource.itemResource.id}"
  http_method = "${aws_api_gateway_method.itemOptionsMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.itemOptionsMethod200Response.status_code}"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,x-amz-meta-fileinfo'"
    "method.response.header.Access-Control-Allow-Methods" = "'GET,PUT,OPTIONS'"
    "method.response.header.Access-Control-Allow-Origin" = "'*'"
  }

  response_templates = {
    "application/json" = ""
  }

  depends_on = [
    "aws_api_gateway_method_response.itemOptionsMethod200Response",
    "aws_api_gateway_integration.itemOptionsMethod-ApiProxyIntegration"]
}