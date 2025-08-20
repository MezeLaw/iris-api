#!/bin/bash

echo "Creating DynamoDB table for Patients..."

awslocal dynamodb create-table \
    --table-name Patients \
    --attribute-definitions \
        AttributeName=PatientID,AttributeType=S \
    --key-schema \
        AttributeName=PatientID,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --region us-east-1

echo "DynamoDB table 'Patients' created successfully!"

# List tables to verify creation
echo "Listing DynamoDB tables:"
awslocal dynamodb list-tables --region us-east-1