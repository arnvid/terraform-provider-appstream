variable "region" {
  type = string
  default = "eu-west-1"
}

variable "assume_role_arn" {
  type = string
}

provider "appstream" {
  region = var.region
  assume_role {
    role_arn = var.assume_role_arn
  }
}

provider "aws" {
  region = var.region
  assume_role {
    role_arn = var.assume_role_arn
  }
}

variable "tags" {
  description = "A map of tags to add to all resources"
  default     = {}
}

locals {
  appstream_tags = {
    ProviderVersion = "v2.0.0"
    Env             = "test"
  }
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "appstream-test-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["eu-central-1a", "eu-central-1b", "eu-central-1c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway = true
  enable_vpn_gateway = false

  tags = {
    Terraform = "true"
    Environment = "dev"
  }
}

resource "aws_security_group" "appstream_security_group" {
  name       = "appstream-test-security-group"
  vpc_id     = module.vpc.vpc_id
}

resource "aws_security_group_rule" "allow_all_outbound" {
  type        = "egress"
  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
  security_group_id = aws_security_group.appstream_security_group.id
}

resource "aws_security_group_rule" "allow_all_inbbound" {
  type        = "ingress"
  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
  security_group_id = aws_security_group.appstream_security_group.id
}

/*
resource "appstream_image_builder" "test-image-builder" {
  count = 0
  name                           = "test-image-builder"
  appstream_agent_version        = "05-17-2021"
  description                    = "test image builder"
  display_name                   = "test-image-builder"
  enable_default_internet_access = false
  image_name                     = "AppStream-WinServer2019-06-01-2021"
  instance_type                  = "stream.standard.large"
  vpc_config {
    security_group_ids = aws_security_group.appstream_security_group.id
    subnet_ids         = module.vpc.private_subnets[0]
  }
  state = "RUNNING"
}
*/
resource "appstream_stack" "test_stack" {
  name         = "appstream-test-stack"
  description  = "appstream test stack"
  display_name = "test-stack"
  feedback_url = "http://example1.com"
  redirect_url = "http://example1.com"
  storage_connectors {
    connector_type = "HOMEFOLDERS"
  }
  tags = {
    Env  = "lab"
    Role = "appstream-test-stack"
    Terraform = "true"
    Environment = "dev"
    Project = "ISD"
  }
}

resource "appstream_fleet" "test_fleet" {
  name       = "test-fleet"
  stack_name = appstream_stack.test_stack.name
  compute_capacity {
    desired_instances = 1
  }
  description                    = "test fleet"
  disconnect_timeout             = 300
  display_name                   = "test-fleet"
  enable_default_internet_access = false
  fleet_type                     = "ON_DEMAND"
# Use image ARN when you share images betwen multiple accounts in an organization
# Build your image in a shared account.
#  image_arn                     = "arn:aws:appstream:eu-west-1:1231241241:image/Base-Image-Builder-05-02-2018"
  image_name                     = "AppStream-WinServer2019-06-01-2021"
  instance_type                  = "stream.standard.large"
  max_user_duration              = 600
  vpc_config {
    security_group_ids = [aws_security_group.appstream_security_group.id]
#    subnet_ids         = join(", ", module.vpc.private_subnets)
    subnet_ids         = module.vpc.private_subnets
  }
  tags = {
    Env  = "lab"
    Role = "appstream-fleet"
    Terraform = "true"
    Environment = "dev"
  }
  state = "RUNNING"
}

resource "appstream_stack_attachment" "test_attachment" {
  appstream_stack_id = appstream_stack.test_stack.name
  appstream_fleet_id = appstream_fleet.test_fleet.name
}

