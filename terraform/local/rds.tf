resource "aws_db_instance" "test_postgres" {
  allocated_storage    = 5
  max_allocated_storage = 0
  engine               = "postgres"
  engine_version       = "13.3"
  instance_class       = "db.t3.micro"
  name                 = "mydb-ps"
  username             = "foo"
  password             = "foobarbaz"
  skip_final_snapshot  = true
  monitoring_interval = 0
}
