locals {
  redis_port = 6379
}

resource "kubernetes_deployment" "redis" {
  metadata {
    name = "redis"
    namespace = kubernetes_namespace.backend.metadata[0].name
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        "app" = "redis"
      }
    }
    template {
      metadata {
        labels = {
          "app" = "redis"
        }
      }
      spec {
        node_selector = {
            "0x5d.io/purpose" = "cache"
        }
        container {
          image = "redis:latest"
          image_pull_policy = "Always"
          name = "redis"
          resources {
            requests = {
              "cpu" = "1"
              "memory" = "2Gi"
            }
          }
          port {
            name = "redis"
            container_port = local.redis_port
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "redis" {
  depends_on = [kubernetes_deployment.redis]
  metadata {
    name      = "redis"
    namespace = kubernetes_namespace.backend.metadata[0].name
  }
  spec {
    type = "ClusterIP"
    selector = {
      "app" = "redis"
    }
    port {
      name = "redis"
      port = local.redis_port
      target_port = "redis"
    }
  }
}
