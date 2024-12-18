resource "aws_iam_role_policy" "eks_lb_controller" {
  name = "eks_lb_controller"
  role = aws_iam_role.eks_lb_controller.name

  policy = file("${path.module}/files/lb_iam_policy.json")
}

resource "aws_iam_role" "eks_lb_controller" {
  name = "eks_lb_controller"
  assume_role_policy = templatefile("${path.module}/files/lb_assume_role.json", {
    region           = var.region
    oidc_provider_id = var.oidc_provider_id
    account_id       = var.account_id
  })
}

resource "kubernetes_service_account" "eks_lb_controller_svc_account" {
  metadata {
    name      = "aws-load-balancer-controller"
    namespace = "kube-system"
    annotations = {
      "app.kubernetes.io/component" = "controller"
      "eks.amazonaws.com/role-arn"  = aws_iam_role.eks_lb_controller.arn
    }
  }
}