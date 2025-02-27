#!/bin/bash
set -e

# Configuration
CONTAINER_NAME="summarizer-backend-container"
IMAGE_NAME="summarizer-backend"
HOST_PORT=80
CONTAINER_PORT=8080
ENV_FILE_PATH="/home/anpurnama/project/summarizer-backend/.env"
DB_DIR_PATH="/home/anpurnama/project/summarizer-backend/db"

# Build the Docker image
echo "Building Docker image: $IMAGE_NAME"
docker build -t $IMAGE_NAME .

# Check if container with the same name exists
if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    echo "Stopping existing container: $CONTAINER_NAME"
    docker stop $CONTAINER_NAME
    
    echo "Removing existing container: $CONTAINER_NAME"
    docker rm $CONTAINER_NAME
fi

# Run the new container
echo "Starting new container: $CONTAINER_NAME"
docker run \
    --name $CONTAINER_NAME \
    -v $ENV_FILE_PATH:/app/.env \
    -v $DB_DIR_PATH:/app/db \
    -d \
    -p $HOST_PORT:$CONTAINER_PORT \
    $IMAGE_NAME

echo "Container $CONTAINER_NAME is now running"