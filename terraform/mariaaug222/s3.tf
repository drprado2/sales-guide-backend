resource "aws_s3_bucket" "public_assets_s3" {
  bucket = local.public_assets_s3_bucket
  acl    = "public-read-write"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST", "GET", "DELETE"]
    allowed_origins = ["http://localhost:3000"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

resource "aws_s3_bucket" "private_assets_s3" {
  bucket = local.private_assets_s3_bucket
  acl    = "private"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST", "GET", "DELETE"]
    allowed_origins = ["http://localhost:3000"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
