#!/bin/bash

set -e

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    echo "Checking prerequisites..."
    for cmd in docker minikube kubectl; do
        if ! command_exists $cmd; then
            echo "$cmd is not installed. Please install it and try again."
            exit 1
        fi
    done
    echo "All prerequisites are installed."
}

# Start Minikube if it's not running
start_minikube() {
    echo "Checking Minikube status..."
    if ! minikube status | grep -q "Running"; then
        echo "Starting Minikube..."
        minikube start
    else
        echo "Minikube is already running."
    fi
}

# Build and push Docker image to Minikube
# Build and push Docker image to Minikube
build_and_push_image() {
    echo "Building Docker image..."
    eval $(minikube docker-env)
    docker build -t admission-controller:latest .
    echo "Image built successfully."
}
# Deploy the admission controller
deploy_admission_controller() {
    echo "Deploying admission controller..."
    
    # Create namespace
    kubectl create namespace admission-system --dry-run=client -o yaml | kubectl apply -f -

    # Apply deployment
    cat << EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: admission-controller
  namespace: admission-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: admission-controller
  template:
    metadata:
      labels:
        app: admission-controller
    spec:
      containers:
      - name: admission-controller
        image: admission-controller:latest
        imagePullPolicy: Never
        ports:
        - containerPort: 8443
EOF

    # Apply service
    kubectl apply -f deploy/service.yaml

    # Generate certificates
    ./generate-certs.sh

    # Apply webhook configuration
    kubectl apply -f deploy/webhook-configuration.yaml

    echo "Admission controller deployed successfully."
}

# Main execution
main() {
    check_prerequisites
    start_minikube
    build_and_push_image
    deploy_admission_controller

    echo "Admission controller has been built, pushed, and deployed to Minikube."
    echo "You can now test it by creating resources that trigger the admission webhook."
}

main