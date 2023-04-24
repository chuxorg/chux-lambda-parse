terraform {
  backend "s3" {
    bucket = "chux-terraform-state"
    key    = "chux-ecs-terraform.tfstate"
    region = "us-east-1"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

locals {
  service_name = "chux-ecs-service"
}

resource "aws_iam_role" "task_role" {
  name = "${local.service_name}_task_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role" "execution_role" {
  name = "${local.service_name}_execution_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "task_role_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
  role       = aws_iam_role.task_role.id
}

resource "aws_iam_role_policy_attachment" "execution_role_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  role       = aws_iam_role.execution_role.id
}

resource "aws_security_group" "ecs_service_sg" {
  name        = "${local.service_name}_sg"
  description = "Security group for the ECS service"
  vpc_id      = "vpc-0d29c91c33cb0acd7"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }
}

resource "aws_ecs_cluster" "chux_cluster" {
  name = "chux-cluster"
}

resource "aws_ecs_task_definition" "chux_task" {
  family                   = "chux-task-family"
  container_definitions    = jsonencode([{
    name  = "chux-container"
    image = var.image_uri
    essential = true
  }])
  task_role_arn            = aws_iam_role.task_role.arn
  execution_role_arn       = aws_iam_role.execution_role.arn
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
}

resource "aws_ecs_service" "chux_service" {
  name            = local.service_name
  cluster         = aws_ecs_cluster.chux_cluster.id
  task_definition = aws_ecs_task_definition.chux_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = ["subnet-0555ae6b617d99d25"]
    security_groups  = [aws_security_group.ecs_service.id]
    
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.chux_tg.arn
    container_name   = "chux-container"
    container_port   = 80
  }

  depends_on = [aws_iam_role_policy_attachment.task_role_policy, aws_iam_role_policy_attachment.execution_role_policy]
}

resource "aws_security_group" "alb_sg" {
  name        = "chux-alb-sg"
  description = "Security group for the Application Load Balancer"
  vpc_id      = "vpc-0d29c91c33cb0acd7"
}

resource "aws_security_group_rule" "allow_http" {
  security_group_id = aws_security_group.alb_sg.id

  type        = "ingress"
  from_port   = 80
  to_port     = 80
  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_route53_record" "chux_alb_record" {
  zone_id = "Z0984989O0LNST2M7R4F"
  name    = "chuxtone.com"
  type    = "A"

  alias {
    name                   = aws_lb.chux_alb.dns_name
    zone_id                = aws_lb.chux_alb.zone_id
    evaluate_target_health = false
  }
}

resource "aws_lb" "chux_alb" {
  name               = "chux-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = ["subnet-0555ae6b617d99d25"]

  tags = {
    Name = "chux-alb"
  }
}

resource "aws_lb_target_group" "chux_tg" {
  name     = "chux-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = var.vpc_id

  health_check {
    enabled             = true
    interval            = 30
    path                = "/"
    timeout             = 5
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
}

resource "aws_security_group" "ecs_service" {
  name        = "chux-ecs-service-sg"
  description = "Security group for the ECS service"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "chux-ecs-service-sg"
  }
}
