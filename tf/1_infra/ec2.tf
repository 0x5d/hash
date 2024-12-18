resource "aws_instance" "bastion" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = var.bastion_node_instance_type
  availability_zone           = module.vpc.azs[0]
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.bastion.id
  key_name                    = var.bastion_ssh_key
  vpc_security_group_ids      = [aws_security_group.bastion_sg.id]
  subnet_id                   = module.vpc.public_subnets[0]
}

resource "aws_security_group" "bastion_sg" {
  name   = "bastion_sg"
  vpc_id = module.vpc.vpc_id
}

resource "aws_vpc_security_group_egress_rule" "bastion_egress" {
  security_group_id = aws_security_group.bastion_sg.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = -1
}

resource "aws_vpc_security_group_ingress_rule" "bastion_ingress" {
  security_group_id = aws_security_group.bastion_sg.id
  cidr_ipv4         = var.my_ip
  ip_protocol       = -1
}

resource "aws_vpc_security_group_ingress_rule" "eks_node_ingress" {
  security_group_id = module.eks.node_security_group_id
  cidr_ipv4         = "${aws_instance.bastion.private_ip}/32"
  ip_protocol       = -1
}

resource "aws_vpc_security_group_ingress_rule" "eks_cluster_ingress" {
  security_group_id = module.eks.cluster_primary_security_group_id
  cidr_ipv4         = "${aws_instance.bastion.private_ip}/32"
  ip_protocol       = -1
}
