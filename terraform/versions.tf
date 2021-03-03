terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.30.0"
    }
    archive = {
      souce   = "hashicorp/archive"
      version = "~> 2.1.0"
    }
  }
  required_version = "~> 0.13.6"
}
