#!/bin/bash

# Update system packages
sudo yum update -y

# Install Docker
sudo yum install -y docker
sudo service docker start
sudo usermod -a -G docker ec2-user

# Install AWS CLI
sudo yum install -y aws-cli

# Create application directory
sudo mkdir -p /home/ec2-user/workdone
sudo chown ec2-user:ec2-user /home/ec2-user/workdone

# Create environment file
cat > /home/ec2-user/workdone/.env << EOL
ENV=production
JWT_ACCESS_SIGN_KEY=${JWT_ACCESS_SIGN_KEY}
JWT_REFRESH_SIGN_KEY=${JWT_REFRESH_SIGN_KEY}
JWT_ISSUER=workdone.me
APP_HOST=0.0.0.0
APP_PORT=9090
DB_HOST=${DB_HOST}
DB_PORT=3306
DB_DRIVER=mysql
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME}
GOOGLE_OAUTH_CLIENT_ID=${GOOGLE_OAUTH_CLIENT_ID}
GOOGLE_OAUTH_CLIENT_SECRET=${GOOGLE_OAUTH_CLIENT_SECRET}
GOOGLE_OAUTH_REDIRECTION_URL=${GOOGLE_OAUTH_REDIRECTION_URL}
EOL

# Set proper permissions
sudo chown ec2-user:ec2-user /home/ec2-user/workdone/.env
sudo chmod 600 /home/ec2-user/workdone/.env

# Configure AWS CLI (if using ECR)
aws configure set default.region ap-south-1

# Install Docker Compose (optional)
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create a systemd services for the application
cat > /etc/systemd/system/workdone.service << EOL
[Unit]
Description=Workdone AI API
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/ec2-user/workdone
ExecStart=/usr/local/bin/docker-compose up -d
ExecStop=/usr/local/bin/docker-compose down
User=ec2-user

[Install]
WantedBy=multi-user.target
EOL

# Reload systemd and enable services
sudo systemctl daemon-reload
sudo systemctl enable workdone.service

# Setup log rotation
cat > /etc/logrotate.d/workdone << EOL
/home/ec2-user/workdone/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 ec2-user ec2-user
}
EOL

# Create docker-compose file
cat > /home/ec2-user/workdone/docker-compose.yml << EOL
version: '3.8'

services:
  api:
    image: ${ECR_REGISTRY}/${ECR_REPOSITORY}:latest
    container_name: workdone-api
    restart: unless-stopped
    ports:
      - "9090:9090"
    env_file:
      - .env
    volumes:
      - ./logs:/app/logs
EOL

# Create logs directory
mkdir -p /home/ec2-user/workdone/logs
chown ec2-user:ec2-user /home/ec2-user/workdone/logs

echo "EC2 setup completed successfully!" 