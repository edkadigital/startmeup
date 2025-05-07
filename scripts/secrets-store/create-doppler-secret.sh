#!/bin/bash

# Usage: ./create-doppler-secret.sh <secret-name> <doppler-token>

SECRET_NAME="$1"
DOPPLER_TOKEN="$2"

if [ -z "$SECRET_NAME" ]; then
  read -p "Enter the Kubernetes secret name: " SECRET_NAME
fi

if [ -z "$DOPPLER_TOKEN" ]; then
  read -sp "Enter the Doppler token: " DOPPLER_TOKEN
  echo
fi

HISTIGNORE='*kubectl*' kubectl --namespace=external-secrets create secret generic \
    "$SECRET_NAME" \
    --from-literal dopplerToken="$DOPPLER_TOKEN"

if [ $? -eq 0 ]; then
  echo "Secret '$SECRET_NAME' created successfully."
else
  echo "Failed to create secret '$SECRET_NAME'." >&2
  exit 1
fi 