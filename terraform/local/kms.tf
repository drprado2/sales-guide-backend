resource "aws_kms_key" "test_key" {
  description             = "Key for test"
  key_usage = "ENCRYPT_DECRYPT"
  deletion_window_in_days = 7
}