#!/bin/bash

# Generate self-signed certificates for JWT authentication and HTTPS
# This script creates RSA private and public keys plus a self-signed certificate

set -e

CERT_DIR="certs"
PRIVATE_KEY="$CERT_DIR/private.pem"
PUBLIC_KEY="$CERT_DIR/public.pem"
CERTIFICATE="$CERT_DIR/certificate.pem"
CSR="$CERT_DIR/certificate.csr"

echo "Creating certificates directory..."
mkdir -p "$CERT_DIR"

echo "Generating RSA private key (PKCS1 format)..."
openssl genrsa -traditional -out "$PRIVATE_KEY" 2048

echo "Generating RSA public key..."
openssl rsa -in "$PRIVATE_KEY" -pubout -out "$PUBLIC_KEY"

echo "Generating Certificate Signing Request (CSR)..."
openssl req -new -key "$PRIVATE_KEY" -out "$CSR" -subj "/C=US/ST=State/L=City/O=Organization/OU=Unit/CN=localhost"

echo "Generating self-signed certificate..."
openssl x509 -req -in "$CSR" -signkey "$PRIVATE_KEY" -out "$CERTIFICATE" -days 365

echo "Setting proper permissions..."
chmod 600 "$PRIVATE_KEY"
chmod 644 "$PUBLIC_KEY"
chmod 644 "$CERTIFICATE"
chmod 644 "$CSR"

echo "Certificate generation complete!"
echo "Private key: $PRIVATE_KEY"
echo "Public key: $PUBLIC_KEY"
echo "Certificate: $CERTIFICATE"
echo "CSR: $CSR"

# Verify the key formats
echo ""
echo "Verifying key formats..."
echo "Private key header:"
head -1 "$PRIVATE_KEY"
echo "Public key header:"
head -1 "$PUBLIC_KEY"
echo "Certificate header:"
head -1 "$CERTIFICATE"

# Display certificate information
echo ""
echo "Certificate information:"
openssl x509 -in "$CERTIFICATE" -text -noout | grep -E "(Subject:|Issuer:|Not Before|Not After)"
