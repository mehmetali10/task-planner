#!/bin/bash

VERSION="0.0.1"

PARENT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Building Docker images with version $VERSION..."

IMAGES_BUILT=()

find "$PARENT_DIR" -type f -name 'Dockerfile*' | while read dockerfile; do
    IMAGE_NAME=$(basename "$dockerfile" | sed 's/^Dockerfile\.//')

    echo "Building Docker image: $IMAGE_NAME with version $VERSION..."

    cd "$PARENT_DIR"
    docker build -t "$IMAGE_NAME:$VERSION" -f "$dockerfile" .

    if [[ $? -eq 0 ]]; then
        echo "Image $IMAGE_NAME:$VERSION built successfully!"
        IMAGES_BUILT+=("$IMAGE_NAME:$VERSION")  
    else
        echo "Failed to build image $IMAGE_NAME:$VERSION."
    fi
done


echo "Build process complete."
