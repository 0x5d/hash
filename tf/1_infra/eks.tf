module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.31"

  cluster_name    = var.cluster_name
  cluster_version = var.kubernetes_version

  enable_cluster_creator_admin_permissions = true

  eks_managed_node_groups = {
    api = {
      ami_type       = "AL2023_x86_64_STANDARD"
      instance_types = [var.api_node_instance_type]

      desired_size = var.api_node_count
      min_size     = var.api_node_count
      max_size     = var.api_node_count + 1
      labels = {
        "0x5d.io/purpose" = "api"
      }
    }

    cache = {
      ami_type       = "AL2023_x86_64_STANDARD"
      instance_types = [var.cache_node_instance_type]

      desired_size = var.cache_node_count
      min_size     = var.cache_node_count
      max_size     = var.cache_node_count + 1
      labels = {
        "0x5d.io/purpose" = "cache"
      }
    }
  }

  vpc_id      = module.vpc.vpc_id
  subnet_ids  = module.vpc.private_subnets
  enable_irsa = true

  tags = {
    Terraform = "true"
  }
}
