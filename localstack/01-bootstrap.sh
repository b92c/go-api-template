#!/usr/bin/env bash
set -euo pipefail

# This script runs inside LocalStack container on readiness
# You can create default resources here if needed

awslocal s3 mb s3://go-api-template-bucket || true

# Create API Gateway HTTP API and route to Lambda if Lambda exists
# Placeholder for future automation
