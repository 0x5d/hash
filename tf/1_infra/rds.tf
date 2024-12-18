resource "aws_db_instance" "backend_db" {
  identifier             = var.db_name
  engine                 = "postgres"
  engine_version         = "16.3"
  instance_class         = var.db_instance_class
  username               = var.db_username
  password               = var.db_password
  allocated_storage      = var.db_capacity
  db_subnet_group_name   = aws_db_subnet_group.backend_db.name
  vpc_security_group_ids = [aws_security_group.backend_db_sg.id]

  storage_encrypted   = true
  skip_final_snapshot = true
  publicly_accessible = false

  tags = {
    Name = "postgres-db"
  }
}

resource "aws_db_subnet_group" "backend_db" {
  name       = "main"
  subnet_ids = module.vpc.private_subnets
}

resource "aws_security_group" "backend_db_sg" {
  name   = "backend_db_sg"
  vpc_id = module.vpc.vpc_id
}

resource "aws_vpc_security_group_ingress_rule" "backend_db_ingress" {
  security_group_id = aws_security_group.backend_db_sg.id
  cidr_ipv4         = module.vpc.vpc_cidr_block
  ip_protocol       = "tcp"
  from_port         = 5432
  to_port           = 5432
}