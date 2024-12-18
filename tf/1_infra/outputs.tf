output "kubeconfig_cmd" {
  description = "Command to add credentials to .kubeconfig file"
  value       = "aws eks update-kubeconfig --region ${var.region} --name ${var.cluster_name}"
}

output "bastion_ssh_cmd" {
  value = "ssh -i ${var.bastion_ssh_key}.pem ubuntu@${aws_instance.bastion.public_dns}"
}

output "psql_cmd" {
  value = "psql -d postgres -h ${aws_db_instance.backend_db.address} -U ${var.db_username}"
}

output "vpc_id" {
  value = module.vpc.vpc_id
}