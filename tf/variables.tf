variable "aws_access_key_id" {
  description = "AWS access key ID"
  type        = string
}

variable "aws_secret_access_key" {
  description = "AWS secret access key"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "image_uri" {
  type        = string
  description = "The ECR image URI for the Lambda function"
}

# variable "vpc_id" {
#   type        = string
#   description = "The ID of the VPC where the ECS Fargate cluster and other resources will be deployed"
#   default     = "vpc-0d29c91c33cb0acd7"
# }

