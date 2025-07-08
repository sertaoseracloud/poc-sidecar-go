# Terraform Configuration

This directory contains example Terraform modules for provisioning Kubernetes clusters on AWS (EKS) and Azure (AKS).

Each module requires pre-existing networking resources (such as subnets or resource groups) and credentials configured for the respective provider.

## EKS

```
cd eks
terraform init
terraform apply -var=aws_region=us-east-1 -var=cluster_name=poc -var='subnet_ids=["subnet-123","subnet-456"]'
```

## AKS

```
cd aks
terraform init
terraform apply -var=cluster_name=poc -var=location=eastus -var=resource_group=my-rg -var=dns_prefix=poc
```
