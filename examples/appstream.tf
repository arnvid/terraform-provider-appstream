provider "appstream" {
  version = "v1.0.8"
  assume_role {
    role_arn = var.assume_role_arn
  }
  region = var.region_primary
}


resource "appstream_image_builder" "test-image-builder" {
  name                           = "test-image-builder"
  appstream_agent_version        = "LATEST"
  description                    = "test image builder"
  display_name                   = "test-image-builder"
  enable_default_internet_access = true
  image_name                     = "Base-Image-Builder-05-02-2018"
  instance_type                  = "stream.standard.large"
  vpc_config {
    security_group_ids = "sg-b5af81d3"
    subnet_ids         = "subnet-7a5f4b51"
  }
  state = "RUNNING"
}


resource "appstream_stack" "test-stack" {
  name         = "test-stack"
  description  = "appstream test stack"
  display_name = "test-stack"
  feedback_url = "http://example1.com"
  redirect_url = "http://example1.com"
  storage_connectors {
    connector_type = "HOMEFOLDERS"
  }
  tags {
    Env  = "lab"
    Role = "appstream-stack"
  }
}

resource "appstream_fleet" "test-fleet" {
  name       = "test-fleet"
  stack_name = appstream_stack.test-stack.name
  compute_capacity {
    desired_instances = 1
  }
  description                    = "test fleet"
  disconnect_timeout             = 300
  display_name                   = "test-fleet"
  enable_default_internet_access = true
  fleet_type                     = "ON_DEMAND"
  image_name                     = "arn:aws:appstream:eu-west-1:1231241241:image/Base-Image-Builder-05-02-2018"
  instance_type                  = "stream.standard.large"
  max_user_duration              = 600
  vpc_config {
    security_group_ids = "sg-b5af81d3"
    subnet_ids         = "subnet-7a5f4b51,subnet-7a5f1231"
  }
  tags {
    Env  = "lab"
    Role = "appstream-fleet"
  }
  state = "RUNNING"
}
