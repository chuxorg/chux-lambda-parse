provider "aws" {
  region = "us-east-1" # Change this to your desired AWS region
}

locals {
  function_name = "chux-lambda-parse"
}

resource "aws_iam_role" "lambda_role" {
  name = "${local.function_name}_execution_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "s3_secretsmanager_cloudwatch" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
  role       = aws_iam_role.lambda_role.id
}

resource "aws_iam_role_policy_attachment" "secretsmanager_policy" {
  policy_arn = "arn:aws:iam::aws:policy/SecretsManagerReadWrite"
  role       = aws_iam_role.lambda_role.id
}

resource "aws_iam_role_policy_attachment" "cloudwatch_policy" {
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchFullAccess"
  role       = aws_iam_role.lambda_role.id
}

resource "aws_security_group" "lambda_sg" {
  name        = "${local.function_name}_sg"
  description = "Security group for Lambda function to access the internet"
  vpc_id      = "vpc-0d29c91c33cb0acd7"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/8"] # Adjust this to your VPC's CIDR block
  }
}

resource "aws_lambda_function" "chux_lambda_parse" {
  function_name = local.function_name
  handler       = "parseHandler" # Update this to the correct handler in your Go package
  runtime       = "go1.x"
  role          = aws_iam_role.lambda_role.arn

  filename = "chux-lambda-parser.zip" # Make sure to create the deployment package

  vpc_config {
    subnet_ids         = ["subnet-009f7d01c00791a01"]
    security_group_ids = [aws_security_group.lambda_sg.id]
  }

  environment {
    variables = {
      AWS_ACCESS_KEY_ID     = var.aws_access_key_id
      AWS_SECRET_ACCESS_KEY = var.aws_secret_access_key
      AWS_REGION            = var.aws_region
    }
  }
}
