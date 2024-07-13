#!/bin/bash

# Set the output directory
OUTPUT_DIR="./builds"

# Set the app name
APP_NAME="git-branch-updater"

# Create the output directory if it doesn't exist
mkdir -p $OUTPUT_DIR

# Build GO
echo "Building for GO..."
go build -o $OUTPUT_DIR/$APP_NAME

# Build for Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-windows-amd64.exe

# Build for macOS
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-darwin-amd64

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-linux-amd64

echo "Build completed. Binaries are in the $OUTPUT_DIR directory."
