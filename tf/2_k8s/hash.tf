locals {
  port = 8080
}

resource "kubernetes_deployment" "hash" {
  depends_on = [kubernetes_service.redis]
  metadata {
    name = "hash"
    namespace = kubernetes_namespace.backend.metadata[0].name
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        "app" = "hash"
      }
    }
    template {
      metadata {
        labels = {
          "app" = "hash"
        }
      }
      spec {
        node_selector = {
            "0x5d.io/purpose" = "api"
        }
        container {
          image = "ox5d/hash-89ffed742e210a369c46bda720795ff5@sha256:8c3a053da24ec32c61747caeefac5b45a949f28eb11989229f3b520b3ffe710a"
          image_pull_policy = "Always"
          name = "hash"
          resources {
            requests = {
              "cpu" = "1700m"
              "memory" = "500Mi"
            }
          }
          port {
            name = "http"
            container_port = local.port
          }
          env {
            name = "HTTP_ADV_ADDR"
            value = kubernetes_service.hash.status[0].load_balancer[0].ingress[0].hostname
          }
          env {
            name = "DB_URL"
            value = var.db_url
          }
          env {
            name = "DB_ID_TABLE_ID"
            value = var.db_id_table_id
          }
          env {
            name = "DB_MIGRATIONS_TABLE"
            value = "migrations"
          }
          env {
            name =  "CACHE_ADDR"
            value = "${kubernetes_service.redis.metadata[0].name}.${kubernetes_namespace.backend.metadata[0].name}.svc.cluster.local:${kubernetes_service.redis.spec[0].port[0].port}"
          }
          env {
            name = "CACHE_TTL"
            value = "1d"
          }
          env {
            name = "CACHE_READ_TIMEOUT"
            value = "1s"
          }
          env {
            name = "LOG_LEVEL"
            value = "INFO"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "hash" {
  depends_on = [helm_release.eks_lb_controller]
  metadata {
    name      = "hash"
    namespace = kubernetes_namespace.backend.metadata[0].name
    annotations = {
      "service.beta.kubernetes.io/aws-load-balancer-type" : "external"
      "service.beta.kubernetes.io/aws-load-balancer-nlb-target-type" : "instance"
      "service.beta.kubernetes.io/aws-load-balancer-scheme" : "internet-facing"
    }
  }
  spec {
    type = "LoadBalancer"
    selector = {
      "app" = "hash"
    }
    port {
      name = "http"
      port = 80
      target_port = "http"
    }
  }
  wait_for_load_balancer = true
}
