#!/bin/bash

# Set variables
CERT_DIR="/tmp/k8s-webhook-server/serving-certs"
SERVICE="admission-controller-service"
NAMESPACE="default"
SECRET_NAME="webhook-server-tls"

# Create directory for certs
mkdir -p ${CERT_DIR}

# Generate CA certificate
openssl genpkey -algorithm RSA -out ${CERT_DIR}/ca.key
openssl req -x509 -new -nodes -key ${CERT_DIR}/ca.key -subj "/CN=ca" -days 365 -out ${CERT_DIR}/ca.crt

# Generate server certificate
openssl genpkey -algorithm RSA -out ${CERT_DIR}/tls.key
openssl req -new -key ${CERT_DIR}/tls.key -subj "/CN=${SERVICE}.${NAMESPACE}.svc" -out ${CERT_DIR}/tls.csr
openssl x509 -req -in ${CERT_DIR}/tls.csr -CA ${CERT_DIR}/ca.crt -CAkey ${CERT_DIR}/ca.key -CAcreateserial -out ${CERT_DIR}/tls.crt -days 365

# Create Kubernetes secret
kubectl create secret tls ${SECRET_NAME} --cert=${CERT_DIR}/tls.crt --key=${CERT_DIR}/tls.key --namespace ${NAMESPACE}
