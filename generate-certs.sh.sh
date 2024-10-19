#!/bin/bash

set -e

# Create a directory for certificates
mkdir -p certs
cd certs

# Generate CA key and certificate
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 365 -key ca.key -subj "/O=Admission Controller CA/CN=Admission Controller CA" -out ca.crt

# Generate server key and certificate signing request (CSR)
openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj "/O=Admission Controller/CN=admission-controller-webhook-service.admission-system.svc" -out server.csr

# Create a config file for the server certificate
cat > server.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = admission-controller-webhook-service.admission-system.svc
EOF

# Generate server certificate
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -extfile server.ext

# Create Kubernetes secrets
kubectl create secret tls admission-webhook-tls \
    --cert=server.crt \
    --key=server.key \
    --namespace=admission-system \
    --dry-run=client -o yaml | kubectl apply -f -

# Update the ValidatingWebhookConfiguration with the CA bundle
CA_BUNDLE=$(cat ca.crt | base64 | tr -d '\n')
sed -i'' -e "s/caBundle: .*/caBundle: ${CA_BUNDLE}/" ../deploy/webhook-configuration.yaml

echo "Certificates generated and Kubernetes resources updated successfully."