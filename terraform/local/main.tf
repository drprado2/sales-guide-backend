# local
terraform {
  backend "s3" {
    bucket                      = "terraform-state"
    key                         = "dev/terraform.tfstate"
    region                      = "sa-east-1"
    endpoint                    = "http://localhost:4566"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    force_path_style            = true
    encrypt                     = true
  }
}

provider "aws" {
  access_key                  = var.aws_access_key
  secret_key                  = var.aws_secret_key
  region                      = var.aws_region
  s3_force_path_style = true
  skip_credentials_validation = true
  skip_metadata_api_check = true
  skip_requesting_account_id = true

  endpoints {
    s3 = "http://localhost:4566"
    dynamodb = "http://localhost:4566"
    rds = "http://localhost:4566"
    route53 = "http://localhost:4566"
    secretsmanager = "http://localhost:4566"
    sqs = "http://localhost:4566"
    stepfunctions = "http://localhost:4566"
  }
}