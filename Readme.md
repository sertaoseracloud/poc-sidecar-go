# Multicloud Sidecar PoC

This proof of concept demonstrates a sidecar architecture for federated authentication across AWS and Azure using RabbitMQ for asynchronous messaging.

## Components

 - **sidecar-go** – HTTP/gRPC server implemented in Go that publishes authentication requests to RabbitMQ and returns credentials to the caller.
 - **auth-worker** – Go service that consumes requests, obtains credentials from AWS STS or Azure Entra ID via the adapters and sends the results back.
 - **identity-adapters** – Go module with basic adapters for AWS and Azure SDKs (returns fake credentials when the SDK calls fail).
 - **mock-app** – Go client that requests credentials from the sidecar.
- **infra** – Docker Compose setup running RabbitMQ and all services.

## Running locally

```bash
go build ./...
cd packages/infra
docker compose up --build
```

The mock application will perform both HTTP and gRPC requests to the sidecar and print the returned credentials.

## Running on Kubernetes (AKS/EKS)

The `packages/infra/k8s` directory contains manifest files for deploying all services to a Kubernetes cluster. Build and push the Docker images to a registry accessible from your cluster and then apply the manifests:

```bash
docker build -t sidecar-go:latest packages/sidecar-go
docker build -t auth-worker:latest packages/auth-worker
docker build -t mock-app:latest packages/mock-app
kubectl apply -f packages/infra/k8s
```

The `sidecar-go` service is exposed as a LoadBalancer and publishes ports 3000 (HTTP) and 50051 (gRPC). These manifests are provider agnostic and can be used on both AKS and EKS.

## Provisioning clusters with Terraform

Example Terraform modules are available under `packages/infra/terraform` for creating an EKS or AKS cluster. The modules assume that networking resources already exist.

```bash
cd packages/infra/terraform/eks
terraform init
terraform apply -var=aws_region=us-east-1 -var=cluster_name=poc -var='subnet_ids=["subnet-123","subnet-456"]'
```

```bash
cd packages/infra/terraform/aks
terraform init
terraform apply -var=cluster_name=poc -var=location=eastus -var=resource_group=my-rg -var=dns_prefix=poc
```
