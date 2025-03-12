#!/bin/bash

# Exit on error
set -e

# Load environment variables
source .env

# Build and push Docker images
echo "Building and pushing Docker images..."
docker-compose build
docker-compose push

# Connect to production server
echo "Connecting to production server..."
ssh $PROD_SERVER_USER@$PROD_SERVER_HOST << 'EOF'
    # Pull latest images
    cd /opt/crypto-monitor
    docker-compose pull
    
    # Run database migrations
    docker-compose run --rm backend ./migrations up
    
    # Restart services
    docker-compose down
    docker-compose up -d
    
    # Check logs for errors
    docker-compose logs --tail=100
EOF

echo "Deployment completed successfully!"