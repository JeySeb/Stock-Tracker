# Database module for CockroachDB Cloud integration
terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "cockroachdb_connection_string" {
  description = "CockroachDB Cloud connection string"
  type        = string
  sensitive   = true
}

variable "is_local" {
  description = "Whether this is a local deployment"
  type        = bool
  default     = false
}

# Local deployment using Docker (for development only)
resource "docker_container" "cockroachdb" {
  count = var.is_local ? 1 : 0
  
  image = "cockroachdb/cockroach:v23.1.0"
  name  = "stock-cockroachdb-${var.environment}"
  
  command = [
    "start-single-node",
    "--insecure",
    "--listen-addr=0.0.0.0:26257",
    "--http-addr=0.0.0.0:8080"
  ]
  
  ports {
    internal = 26257
    external = 26257
  }
  
  ports {
    internal = 8080
    external = 8081
  }
  
  volumes {
    container_path = "/cockroach/cockroach-data"
    volume_name    = "cockroach_data_${var.environment}"
  }
}

# Output connection string (uses CockroachDB Cloud for production)
output "connection_string" {
  value = var.is_local ? 
    "postgres://root@localhost:26257/stockdb?sslmode=disable" :
    var.cockroachdb_connection_string
}

output "database_host" {
  value = var.is_local ? 
    "localhost" :
    "hiring-test-stock-cluster-13493.j77.aws-us-east-1.cockroachlabs.cloud"
}