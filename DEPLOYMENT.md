# Deployment Guide for Workdone AI

## Prerequisites

1. AWS Account with:
   - EC2 instance
   - ECR repository
   - IAM user with appropriate permissions

2. GitHub repository with:
   - Source code
   - GitHub Actions enabled

## Setup Steps

### 1. AWS Setup

#### Create ECR Repository
```bash
aws ecr create-repository --repository-name workdone-api --region ap-south-1
```

#### EC2 Instance Requirements
- Amazon Linux 2
- t2.micro or larger
- Security Group with ports:
  - 22 (SSH)
  - 9090 (API)
  - 3306 (MySQL)

### 2. GitHub Secrets

Add the following secrets to your GitHub repository:

```
AWS_ACCESS_KEY_ID=<your-aws-access-key>
AWS_SECRET_ACCESS_KEY=<your-aws-secret-key>
EC2_INSTANCE_IP=<your-ec2-public-ip>
EC2_USERNAME=ec2-user
SSH_PRIVATE_KEY=<your-ec2-private-key>
```

### 3. EC2 Instance Setup

1. SSH into your EC2 instance:
   ```bash
   ssh -i your-key.pem ec2-user@your-ec2-ip
   ```

2. Copy the setup script:
   ```bash
   scp -i your-key.pem scripts/setup-ec2.sh ec2-user@your-ec2-ip:~/
   ```

3. Run the setup script:
   ```bash
   chmod +x setup-ec2.sh
   ./setup-ec2.sh
   ```

### 4. Environment Variables

Create a `.env` file in `/home/ec2-user/workdone/` with:

```env
ENV=production
JWT_ACCESS_SIGN_KEY=<your-jwt-key>
JWT_REFRESH_SIGN_KEY=<your-jwt-refresh-key>
JWT_ISSUER=workdone.me
APP_HOST=0.0.0.0
APP_PORT=9090
DB_HOST=<your-db-host>
DB_PORT=3306
DB_DRIVER=mysql
DB_USER=<your-db-user>
DB_PASSWORD=<your-db-password>
DB_NAME=<your-db-name>
GOOGLE_OAUTH_CLIENT_ID=<your-google-oauth-client-id>
GOOGLE_OAUTH_CLIENT_SECRET=<your-google-oauth-client-secret>
GOOGLE_OAUTH_REDIRECTION_URL=<your-oauth-redirect-url>
```

## Deployment Process

1. Push to main branch triggers automatic deployment
2. GitHub Actions will:
   - Build Docker image
   - Push to ECR
   - Deploy to EC2

## Manual Deployment

If needed, you can manually deploy:

```bash
# SSH into EC2
ssh -i your-key.pem ec2-user@your-ec2-ip

# Pull latest image
docker pull <ecr-repo-url>/workdone-api:latest

# Restart services
sudo systemctl restart workdone
```

## Monitoring

1. Check container logs:
   ```bash
   docker logs workdone-api
   ```

2. Check service status:
   ```bash
   sudo systemctl status workdone
   ```

3. View application logs:
   ```bash
   tail -f /home/ec2-user/workdone/logs/app.log
   ```

## Troubleshooting

1. If container fails to start:
   ```bash
   docker ps -a
   docker logs workdone-api
   ```

2. Check service logs:
   ```bash
   journalctl -u workdone
   ```

3. Check Docker status:
   ```bash
   sudo systemctl status docker
   ```

## Backup and Maintenance

1. Database backup:
   ```bash
   # Add your backup script here
   ```

2. Log rotation is automatic (configured in /etc/logrotate.d/workdone)

3. Container cleanup:
   ```bash
   docker system prune -af
   ``` 