variable "cluster_name" {
  description = "The name of the EKS cluster."
  type        = string
}

variable "eks_token" {
  description = "The token to authenticate with the EKS cluster."
  type        = string
}

variable "region" {
  description = "The AWS region where the EKS cluster will be created."
  type        = string
}

variable "vpc_id" {
  description = "The VPC ID."
  type        = string
}

variable "http_advertised_addr" {
  description = "The advertised address for the HTTP server."
  type        = string
}

variable "db_url" {
  description = "The URL for the database connection."
  type        = string
}

variable "db_id_table_id" {
  description = "The ID table ID to generate URL ids."
  type        = string
}

variable "cache_ttl" {
  description = "The time-to-live for cached items."
  type        = string
}

variable "cache_read_timeout" {
  description = "The read timeout for the cache."
  type        = string
}

variable "oidc_provider_id" {
  description = <<DESC
  Obtained by running
  aws eks describe-cluster --region <region> --name <cluster name> --query "cluster.identity.oidc.issuer" --output text | cut -d '/' -f 5
DESC
  type        = string
}

variable "account_id" {
  description = "AWS account ID"
  type        = string
}
