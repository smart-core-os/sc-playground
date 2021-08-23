#!/usr/bin/env bash
set -e

# Help for those not using openssl all the time

# create a CSR (cert signing request)
# openssl req\
#   # instruct openssl to output a self signed cert instead of a CSR
#   -x509 \
#   # generate a new CSR and private key (as opposed to reading one in)
#   -newkey rsa:4096 \
#   # (no des) do not encrypt the private key file
#   -nodes \
#   # specify the filename for the CAs private key file
#   -keyout ca.key \
#   # specify the filename for the CAs CSR file
#   -out ca.csr \
#   # the contents of the certificate that will be signed by the private key
#   -subj "/CN=localhost/"

# Generate the (self signed) CA private key and certificate.
# This is the root of trust.
openssl req \
  -x509 \
  -newkey rsa:4096 \
  -nodes \
  -keyout ca.key \
  -out ca.crt \
  -subj "/CN=localhost/"

# Generate server and client certs and private keys.
# Both certs are signed using the CAs private key,
# which means the CA cert is all that is needed for either client or server to verify
# the integrity of their peers certificate.

# Generate a CSR (cert sign request) for the server.
# We want the certificate to allow clients to verify the server is who it says it is
openssl req \
  -newkey rsa:4096 \
  -nodes \
  -keyout server.key \
  -out server.csr \
  -subj "/CN=localhost/"
# Use the CA to sign the servers csr
openssl x509 \
  -req \
  -in server.csr \
  -days 365 \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -extfile server-ext.conf \
  -out server.crt


# Generate a CSR (cert sign request) for the client.
# We want the certificate to allow servers to verify the client is who it says it is
openssl req \
  -newkey rsa:4096 \
  -nodes \
  -keyout client.key \
  -out client.csr \
  -subj "/CN=client1/"
# Use the CA to sign the clients csr
openssl x509 \
  -req \
  -in client.csr \
  -days 365 \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -extfile client-ext.conf \
  -out client.crt
