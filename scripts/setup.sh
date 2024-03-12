#!/bin/bash
# Exit script if any command fails and treat unset variables as an error
set -euo pipefail

# Define required environment variables
git_repo_user="${GIT_REPO_USER:?Environment variable GIT_REPO_USER is not set}"
domain_name="${DOMAIN_NAME:?Environment variable DOMAIN_NAME is not set}"

# Update package lists for upgrades and new package installations
sudo apt-get update

# Install necessary packages
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# Check if Docker is installed, if not, install it
if ! command -v docker &> /dev/null
then
    # Add Docker’s official GPG key
    curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    # Set up the Docker stable repository
    echo \
      "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    # Update the apt package index, and install Docker
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io
fi


# Check if Certbot is installed, if not, install it
if ! command -v certbot &> /dev/null
then
    sudo apt-get install -y certbot python3-certbot-nginx
fi


# Run Certbot with the Namecheap plugin to obtain an SSL certificate
sudo certbot --manual --preferred-challenges dns --nginx -d $domain_name


# Check if Nginx is installed, if not, install it
if ! command -v nginx &> /dev/null
then
    sudo apt-get install -y nginx
fi

# Configure Nginx for the domain
echo "server {
    listen 80;
    server_name $domain_name;
    location / {
        return 301 https://\$host\$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name $domain_name;

    ssl_certificate /etc/letsencrypt/live/$domain_name/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$domain_name/privkey.pem;

    location /u/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location / {
        root /var/www/html;
        try_files \$uri /index.html;
    }
}" | sudo tee /etc/nginx/sites-available/$domain_name

# Enable the Nginx configuration by creating a symbolic link
sudo ln -s /etc/nginx/sites-available/$domain_name /etc/nginx/sites-enabled/
# Test Nginx configuration
sudo nginx -t
# Reload Nginx to apply the changes
sudo systemctl reload nginx


# Add a cron job to renew the SSL certificate automatically
echo "0 12 * * * root certbot renew --quiet --deploy-hook 'systemctl reload nginx'" | sudo tee -a /etc/crontab > /dev/null

# Check if Git is installed, if not, install it
if ! command -v git &> /dev/null
then
    sudo apt-get install -y git
fi

# Check if Node.js is installed, if not, install it
if ! command -v node &> /dev/null
then
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
    sudo apt-get install -y nodejs
fi

# Check if Yarn is installed, if not, install it
if ! command -v yarn &> /dev/null
then
    sudo npm install -g yarn
fi

# Run the deployment script
bash deploy.sh