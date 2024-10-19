#!/bin/bash

set -e

# Function to wait for a resource to be ready
wait_for_resource() {
    echo "Waiting for $1 to be ready..."
    kubectl wait --for=condition=ready $1 --timeout=60s
}

# Test the admission controller
test_admission_controller() {
    echo "Testing admission controller..."

    # Create a test pod
    cat << EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: nginx
    image: nginx
EOF

    # Wait for the pod to be created
    wait_for_resource "pod/test-pod"

    # Check if the pod was admitted
    if kubectl get pod test-pod &> /dev/null; then
        echo "Test pod was admitted successfully."
    else
        echo "Test pod was not admitted. Check the admission controller logs for more information."
        exit 1
    fi

    # Clean up
    kubectl delete pod test-pod

    echo "Admission controller test completed successfully."
}

# Main execution
main() {
    test_admission_controller
}

main