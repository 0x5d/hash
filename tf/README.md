# Infraestructura

La infraestructura para la prueba de concepto está sobre AWS y es manejada con Terraform.

- /1_infra contiene la infraestructura base: VPC, subredes, un "bastion", un cluster de K8s, y una instancia de RDS.
- /2_k8s contiene las aplicaciones que corren sobre el cluster de k8s y los roles de IAM necesarios.

Vale aclarar que estos módulos no son considerados "de producción".
