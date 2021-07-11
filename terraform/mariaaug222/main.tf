# local
terraform {
  backend "s3" {
    bucket                      = "terraform-state-616f65f0-33bd-40ea-a46c-c2a4da7ce7f0"
    key                         = "dev/terraform.tfstate"
    region                      = "us-east-2"
    force_path_style            = true
    encrypt                     = true
  }
}

provider "aws" {
  access_key                  = var.aws_access_key
  secret_key                  = var.aws_secret_key
  region                      = var.aws_region
  s3_force_path_style = true
  skip_credentials_validation = false
  skip_metadata_api_check = false
  skip_requesting_account_id = false
}