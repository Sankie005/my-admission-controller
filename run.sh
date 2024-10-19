chmod +x *.sh
echo "Deploying Admission Controller"
./deploy-admission-controller.sh
echo "Generating TLS certifications"
./generate-certs.sh
echo "Testing Admission Controller"
./test-admission-controller.sh