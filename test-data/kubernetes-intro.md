# Kubernetes Introduction

Kubernetes (K8s) is an open-source container orchestration platform that automates deployment, scaling, and management of containerized applications.

## Core Components

| Component | Description |
|-----------|-------------|
| Pod | Smallest deployable unit, wraps one or more containers |
| Service | Stable network endpoint for a set of Pods |
| Deployment | Manages ReplicaSets and rolling updates |
| ConfigMap | External configuration for Pods |
| Secret | Sensitive data (passwords, tokens) |

## Architecture

```
┌─────────────────────────────────┐
│         Control Plane           │
│  ┌─────────┐  ┌──────────────┐ │
│  │ API     │  │ Scheduler    │ │
│  │ Server  │  └──────────────┘ │
│  └─────────┘  ┌──────────────┐ │
│  ┌─────────┐  │ Controller   │ │
│  │ etcd    │  │ Manager      │ │
│  └─────────┘  └──────────────┘ │
└─────────────────────────────────┘
```

## Workload Types

- **Deployment**: Stateless apps (web servers, APIs)
- **StatefulSet**: Stateful apps (databases)
- **DaemonSet**: One per node (log collectors)
- **CronJob**: Scheduled tasks

See [[docker-basics]] for container fundamentals and [[kafka-events]] for event-driven patterns in K8s.

## Getting Started with Kind

```bash
# Install kind
go install sigs.k8s.io/kind@latest

# Create a local cluster
kind create cluster --name dev

# Check it works
kubectl get nodes
```
