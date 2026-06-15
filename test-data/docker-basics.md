# Docker Basics

Docker is a platform for developing, shipping, and running applications in containers.

## Key Concepts

- **Image**: A read-only template with instructions for creating a container
- **Container**: A runnable instance of an image
- **Dockerfile**: A text file with instructions to build an image
- **Volume**: Persistent data storage for containers

## Common Commands

```bash
# Build an image
docker build -t myapp .

# Run a container
docker run -d -p 8080:8080 myapp

# List running containers
docker ps

# Stop a container
docker stop <container-id>
```

## Docker Compose

Docker Compose lets you define multi-container applications. See [[kubernetes-intro]] for orchestration at scale.

## Best Practices

1. Use multi-stage builds to keep images small
2. Don't run containers as root
3. Use `.dockerignore` to exclude unnecessary files
4. Pin image versions (don't use `latest` in production)
