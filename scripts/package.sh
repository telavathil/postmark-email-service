#!/bin/bash
set -e

# Script configuration
BINARY_NAME="bootstrap"
OUTPUT_DIR="build"
ZIP_NAME="function.zip"

# Create build directory if it doesn't exist
mkdir -p $OUTPUT_DIR

echo "Building Go binary for AWS Lambda..."
GOOS=linux GOARCH=amd64 go build -o $OUTPUT_DIR/$BINARY_NAME main.go

echo "Creating deployment package..."
cd $OUTPUT_DIR
zip $ZIP_NAME $BINARY_NAME

echo "Cleaning up..."
rm $BINARY_NAME

echo "Package created at $OUTPUT_DIR/$ZIP_NAME"
