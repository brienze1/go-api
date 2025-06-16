#!/bin/bash

echo "########### Creating profile ###########"

aws configure set aws_access_key_id default_access_key
aws configure set aws_secret_access_key default_secret_key
aws configure set region us-east-2

echo "########### Creating resources ###########"

aws --endpoint-url=http://localhost:4566 secretsmanager create-secret \
  --name notes-database-secret \
  --secret-string '{"username":"postgres","password":"postgres","engine":"postgres","host":"localhost","port":5432,"dbname":"postgres","dbInstanceIdentifier":"notes"}' \
  --region us-east-2

aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name notes-queue --region us-east-2