#!/bin/bash
set -e

# Function to deploy the CloudFormation stack
deploy_stack() {
    local stack_name="$1"
    local stage="$2"
    local region="$3"
    local s3_bucket="$4"
    local postmark_token="$5"
    local dry_run="$6"

    # Package the application
    echo "Packaging application..."
    ./scripts/package.sh

    # Upload deployment package to S3
    echo "Uploading deployment package to S3..."
    aws s3 cp build/function.zip "s3://$s3_bucket/$stack_name/function.zip"

    # Base command for CloudFormation deployment
    local deploy_cmd=(aws cloudformation deploy \
        --template-file iac/base.yml \
        --stack-name "$stack_name" \
        --region "$region" \
        --capabilities CAPABILITY_IAM \
        --parameter-overrides \
            StackName="$stack_name" \
            Stage="$stage" \
            PostmarkToken="$postmark_token" \
            S3Bucket="$s3_bucket" \
            S3Key="$stack_name/function.zip")

    # Add no-execute-changeset flag for dry runs
    if [ "$dry_run" = true ]; then
        deploy_cmd+=(--no-execute-changeset)
        echo "Running in dry-run mode (--no-execute-changeset)"
    fi

    # Execute the deployment command
    echo "Deploying CloudFormation stack..."
    "${deploy_cmd[@]}"

    # Only get and display the API URL if not in dry-run mode
    if [ "$dry_run" = false ]; then
        local api_url
        api_url=$(aws cloudformation describe-stacks \
            --stack-name "$stack_name" \
            --region "$region" \
            --query 'Stacks[0].Outputs[?OutputKey==`ApiEndpoint`].OutputValue' \
            --output text)

        echo "Deployment complete!"
        echo "API Gateway URL: $api_url"
    else
        echo "Dry run complete! No changes were applied."
    fi
}
