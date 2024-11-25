#!/bin/bash
set -e

# Source the common deployment functions
source "$(dirname "$0")/deploy-common.sh"

# Default values
STACK_NAME="postmark-email-service"
STAGE="test"
REGION="us-east-1"
S3_BUCKET=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --stack-name)
      STACK_NAME="$2"
      shift 2
      ;;
    --stage)
      STAGE="$2"
      shift 2
      ;;
    --region)
      REGION="$2"
      shift 2
      ;;
    --s3-bucket)
      S3_BUCKET="$2"
      shift 2
      ;;
    --postmark-token)
      POSTMARK_TOKEN="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Validate required parameters
if [ -z "$S3_BUCKET" ]; then
  echo "Error: --s3-bucket is required"
  exit 1
fi

if [ -z "$POSTMARK_TOKEN" ]; then
  echo "Error: --postmark-token is required"
  exit 1
fi

# Call the deploy function with dry_run=true
deploy_stack "$STACK_NAME" "$STAGE" "$REGION" "$S3_BUCKET" "$POSTMARK_TOKEN" true
