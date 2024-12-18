provider "kubernetes" {
  token       = var.eks_token
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    token       = var.eks_token
    config_path = "~/.kube/config"
  }
}