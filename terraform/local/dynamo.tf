resource "aws_dynamodb_table" "table_1" {
  name           = "table_1"
  read_capacity  = "20"
  write_capacity = "20"
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}