variable "cluster_name" {
  description = "The name of the EKS cluster."
  type        = string
}

variable "region" {
  description = "The AWS region where the EKS cluster will be created."
  type        = string
}

variable "kubernetes_version" {
  description = "The Kubernetes version for the cluster."
  type        = string
  default     = "1.27"
}

variable "api_node_instance_type" {
  description = "The instance type for the API worker nodes."
  type        = string
  default     = "m5.large"
}

variable "cache_node_instance_type" {
  description = "The instance type for the cache worker nodes."
  type        = string
  default     = "m5.large"
}

variable "bastion_node_instance_type" {
  description = "The instance type for the worker nodes."
  type        = string
  default     = "m5.large"
}

variable "api_node_count" {
  description = "The number of worker nodes in the default node group."
  type        = number
  default     = 1
}

variable "cache_node_count" {
  description = "The number of worker nodes in the default node group."
  type        = number
  default     = 1
}

variable "vpc_cidr" {
  description = "The CIDR block for the VPC."
  type        = string
  default     = "10.0.0.0/16"
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for the subnets."
  type        = list(string)
  default     = ["10.0.1.0/24"]
}

variable "subnet_availability_zones" {
  description = "List of availability zones for the subnets."
  type        = list(string)
}

variable "bastion_ssh_key" {
  description = "The bastion instance SSH key."
  type        = string
}

variable "my_ip" {
  description = "Your IP address to allow SSH access to the bastion."
  type        = string
}

variable "db_instance_class" {
  description = "The instance class for the RDS instance."
  type        = string
  default     = "db.t3.medium"
}

variable "db_capacity" {
  description = "The allocated storage in GB for the RDS instance."
  type        = number
  default     = 10
}

variable "db_name" {
  description = "The name of the database to create."
  type        = string
  default     = "backend"
}

variable "db_username" {
  description = "The username for the database."
  type        = string
}

variable "db_password" {
  description = "The password for the database."
  type        = string
  sensitive   = true
}