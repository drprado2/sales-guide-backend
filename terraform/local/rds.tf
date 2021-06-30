resource "aws_db_instance" "test_mysql" {
  allocated_storage    = 5
  max_allocated_storage = 15
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t3.micro"
  name                 = "mydb"
  username             = "foo"
  password             = "foobarbaz"
  parameter_group_name = "default.mysql5.7"
  skip_final_snapshot  = true
  monitoring_interval = 15
}