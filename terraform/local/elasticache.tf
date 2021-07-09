resource "aws_elasticache_cluster" "redis-cluster-test" {
  cluster_id = "redis-cluster-test"
  engine = "redis"
  node_type = "cache.m4.large"
  num_cache_nodes = 1
  parameter_group_name = "default.redis3.2"
  engine_version = "3.2.10"
  port = 4530
  apply_immediately = true
}