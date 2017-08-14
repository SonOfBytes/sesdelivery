variable "aws_region" {
  type = "string"
  default = "eu-west-1"
}

provider "aws" {
  region = "${var.aws_region}"
}

data "aws_caller_identity" "current" {}