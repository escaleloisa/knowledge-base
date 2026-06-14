# Deployment Guide — AWS EC2

## Overview

The knowledge base is deployed as a Docker Compose stack on a single EC2 instance. All services (Go backend, Kafka, Elasticsearch, Redis, PostgreSQL, Next.js frontend) run as containers on the same machine.

## Prerequisites

- AWS account
- SSH client (built into Windows PowerShell)
- Git repository pushed to GitHub

## Step 1: Launch EC2 Instance

1. Open https://console.aws.amazon.com/ec2
2. Select region: **ap-southeast-1** (Singapore)
3. Click **Launch Instance**

### Configuration

| Setting | Value |
|---------|-------|
| Name | `knowledge-base-demo` |
| AMI | Ubuntu Server 26.04 LTS (HVM), EBS General Purpose (SSD) |
| Architecture | 64-bit (x86) |
| Instance type | `t2.medium` (2 vCPU, 4GB RAM) |
| Key pair | RSA, `.pem` format — name: `kb-demo` |
| Storage | 30 GB gp3 |

### Security Group (Inbound Rules)

| Type | Port | Source | Purpose |
|------|------|--------|---------|
| SSH | 22 | My IP | SSH access |
| Custom TCP | 3000 | 0.0.0.0/0 | Next.js frontend |
| Custom TCP | 8080 | 0.0.0.0/0 | Note Service API |
| Custom TCP | 8081 | 0.0.0.0/0 | Search API |
| Custom TCP | 8082 | 0.0.0.0/0 | Tag Aggregator API |

4. Click **Launch Instance**
5. Go to Instances → copy the **Public IPv4 address**

## Step 2: Fix PEM File Permissions (Windows)

The `.pem` key file needs restricted permissions or SSH will reject it:

```powershell
icacls "D:\eloisa-demo\kb-demo.pem" /inheritance:r /remove:g "Everyone" /remove:g "BUILTIN\Users" /remove:g "NT AUTHORITY\Authenticated Users"
icacls "D:\eloisa-demo\kb-demo.pem" /grant:r "${env:USERNAME}:(R)"
```

## Step 3: SSH into the Instance

```powershell
ssh -i "D:\eloisa-demo\kb-demo.pem" ubuntu@<EC2-PUBLIC-IP>
```

Type `yes` when prompted about host key authenticity.

## Step 4: Install Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Install Docker Compose plugin
sudo apt-get update && sudo apt-get install docker-compose-plugin -y

# Log out and back in for docker group to take effect
exit
```

SSH back in:

```powershell
ssh -i "D:\eloisa-demo\kb-demo.pem" ubuntu@<EC2-PUBLIC-IP>
```

Verify Docker works:

```bash
docker --version
docker compose version
```

## Step 5: Clone and Configure

```bash
# Clone the repo
git clone https://github.com/escaleloisa/knowledge-base.git
cd knowledge-base

# Switch to the feature branch
git checkout feature/web-ui

# Create .env (replace YOUR_EC2_IP with actual public IP)
cat > .env << 'EOF'
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=knowledgebase
PORT=8080
NEXT_PUBLIC_NOTE_SERVICE_URL=http://YOUR_EC2_IP:8080
NEXT_PUBLIC_SEARCH_API_URL=http://YOUR_EC2_IP:8081
NEXT_PUBLIC_TAG_API_URL=http://YOUR_EC2_IP:8082
EOF
```

## Step 6: Build and Run

```bash
docker compose up -d --build
```

This takes 5-10 minutes on first run (downloading images + building Go binaries + building Next.js).

Verify all services are running:

```bash
docker compose ps
```

All services should show status `Up`:
- postgres, kafka, zookeeper, elasticsearch, redis (infrastructure)
- note-service, search-api, tag-aggregator, search-indexer, backlink-extractor (backend)
- web-ui (frontend)

## Step 7: Access the App

Open in browser: `http://<EC2-PUBLIC-IP>:3000`

## Updating After Code Changes

SSH into the instance and run:

```bash
cd ~/knowledge-base
git pull
docker compose up -d --build
```

To rebuild only specific services (faster):

```bash
# Frontend only
docker compose up -d --build web-ui

# Backend services only
docker compose up -d --build note-service search-api tag-aggregator backlink-extractor
```

## Useful Commands

```bash
# View logs for a service
docker compose logs note-service --tail=30
docker compose logs web-ui --tail=30

# Restart a service
docker compose restart note-service

# Stop everything
docker compose down

# Stop and remove all data (fresh start)
docker compose down -v

# Flush Redis (clear orphaned tags)
docker compose exec redis redis-cli FLUSHALL

# Reset Postgres data
docker compose exec postgres psql -U postgres -d knowledgebase -c "DELETE FROM notes;"
docker compose exec postgres psql -U postgres -d knowledgebase -c "DELETE FROM collections;"
```

## Managing the EC2 Instance

### Stop (save money, keep data)

AWS Console → EC2 → Instances → Select instance → Instance state → **Stop instance**

- Storage cost while stopped: ~$3/month
- No compute charges
- Data persists

### Start again

AWS Console → Instance state → **Start instance**

**Important:** The public IP changes on restart. Update the `.env` with the new IP and rebuild the frontend:

```bash
cd ~/knowledge-base
# Edit .env with new IP
nano .env
# Rebuild frontend with new API URLs
docker compose up -d --build web-ui
```

### Terminate (delete everything)

AWS Console → Instance state → **Terminate instance**

Permanently deletes the instance and storage. Zero cost after.

## Cost

| Usage | Approximate cost |
|-------|-----------------|
| Running 24/7 for a month | ~$35 |
| Running for 3 days (demo) | ~$3-4 |
| Stopped (storage only) | ~$3/month |
| Terminated | $0 |
