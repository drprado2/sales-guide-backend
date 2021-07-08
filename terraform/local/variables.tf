variable "aws_region" {
  type    = string
  default = "sa-east-1"
}

variable "aws_access_key" {
  type    = string
  default = "test"
}

variable "aws_secret_key" {
  type    = string
  default = "test"
}

variable "s3_bucket_name" {
  type    = string
  default = "sandpit-sample"
}

variable "account" {
  default = "sales-guide"
}

variable "rest_api_domain_name" {
  default     = "example.com"
  description = "Domain name of the API Gateway REST API for self-signed TLS certificate"
  type        = string
}

variable "rest_api_name" {
  default     = "api-gateway-rest-api-openapi-example"
  description = "Name of the API Gateway REST API (can be used to trigger redeployments)"
  type        = string
}

variable "rest_api_path" {
  default     = "/path1"
  description = "Path to create in the API Gateway REST API (can be used to trigger redeployments)"
  type        = string
}