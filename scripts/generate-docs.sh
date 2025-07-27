#!/bin/bash

# Generate Swagger documentation
echo "Generating Swagger documentation..."

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate docs
swag init -g cmd/server/main.go

echo "Swagger documentation generated successfully!"
echo "You can view it at: http://localhost:8080/swagger/index.html" 