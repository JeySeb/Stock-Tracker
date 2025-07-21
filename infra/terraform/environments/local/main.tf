terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

provider "docker" {}

module "database" {
  source = "../../modules/database"
  
  environment = "local"
  is_local    = true
  cockroachdb_connection_string = "postgresql://jeyseb:<ENTER-SQL-USER-PASSWORD>@hiring-test-stock-cluster-13493.j77.aws-us-east-1.cockroachlabs.cloud:26257/stockdb?sslmode=verify-full&sslrootcert=certs/cc-ca.crt"
}

module "redis" {
  source = "../../modules/redis"
  
  environment = "local"
  is_local    = true
}

output "database_url" {
  value = module.database.connection_string