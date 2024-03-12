#!/bin/bash
# Define variables
git_repo_user="${GIT_REPO_USER:?Environment variable GIT_REPO_USER is not set}"
repo_path="shorturls"

# Rest of the script...

# Check if the repository exists
if [ -d "$repo_path" ]; then
    # If it exists, navigate to it and pull the latest changes
    cd $repo_path
    git pull origin master
    cd ..
else
    # If it doesn't exist, clone the repository
    git clone https://github.com/$git_repo_user/$repo_path.git
fi

# Check if configuration files exist
if [ ! -f config.json ] || [ ! -f wordlist.json ]; then
    echo "Configuration files not found. Exiting."
    exit 1
fi

# Copy configuration files
cp config.json $repo_path/api-server/go/config.json
cp wordlist.json $repo_path/api-server/go/wordlist.json

# Build the server application
cd $repo_path/api-server/go
docker build -t shorturls-api .
cd ../../webui

# Deploy the server application
docker stop shorturls-api || true
docker rm shorturls-api || true
docker run -d --name shorturls-api -p 8080:8080 shorturls-api


# Install dependencies and build the frontend
yarn
yarn prod

# Move the built frontend to the web server's root directory
sudo rm -rf /var/www/html/*
sudo mv dist/* /var/www/html/