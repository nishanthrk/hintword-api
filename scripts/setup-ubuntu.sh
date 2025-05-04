#!/bin/bash

# Update system packages
sudo apt-get update
sudo apt-get upgrade -y

# Install Docker
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y docker-ce
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker ubuntu

# Install AWS CLI
sudo apt-get install -y unzip
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
rm -rf aws awscliv2.zip

# Create application directory
mkdir -p ~/workdone
chmod 755 ~/workdone

# Create environment file template
cat > ~/workdone/.env << EOL
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
chmod 600 ~/workdone/.env

# Configure AWS CLI
aws configure set default.region ap-south-1

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create logs directory
mkdir -p ~/workdone/logs
chmod 755 ~/workdone/logs

# Setup log rotation
sudo tee /etc/logrotate.d/workdone << EOL
/home/ubuntu/workdone/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 ubuntu ubuntu
}
EOL

# Create docker-compose file
cat > ~/workdone/docker-compose.yml << EOL
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

echo "EC2 setup completed successfully!"

# Print next steps
echo "
Next steps:
1. Update ~/workdone/.env with your actual values
2. Configure AWS credentials using: aws configure
3. Test docker with: docker run hello-world
4. Logout and login again for docker permissions to take effect
" 