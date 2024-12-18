resource "aws_iam_instance_profile" "bastion" {
  name = "bastion"
  role = aws_iam_role.bastion.name
}

resource "aws_iam_role" "bastion" {
  name               = "bastion_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy" "eks_policy" {
  name = "bastion_eks_admin"
  role = aws_iam_role.bastion.id

  policy = data.aws_iam_policy_document.eks_admin.json
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }

    actions = [
      "sts:AssumeRole",
    ]
  }
}

data "aws_iam_policy_document" "eks_admin" {
  statement {
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "eks:*",
    ]
  }
}