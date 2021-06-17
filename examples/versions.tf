terraform {
  required_providers {
    appstream = {
      source  = "arnvid/appstream"
      version = "2.0.0"
    }
    aws = {
      source  = "hashicorp/aws"
    }
  }
  required_version = ">= 1.0"
}
