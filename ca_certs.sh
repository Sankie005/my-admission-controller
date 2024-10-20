#!/bin/bash

# Set variables
CERT_DIR="/tmp/k8s-webhook-server/serving-certs"
CA_CERT="${CERT_DIR}/ca.crt"
CA_KEY="${CERT_DIR}/ca.key"

# Create directory for certs
mkdir -p ${CERT_DIR}

# Generate CA certificate
openssl genpkey -algorithm RSA -out ${CA_KEY}
openssl req -x509 -new -nodes -key ${CA_KEY} -subj "/CN=ca" -days 365 -out ${CA_CERT}

# Base64 encode the CA certificate
BASE64_CA_CERT=$(cat ${CA_CERT} | base64 | tr -d '\n')

# Print the base64 encoded CA certificate
echo "Base64 Encoded CA Certificate:"
echo ${BASE64_CA_CERT}
