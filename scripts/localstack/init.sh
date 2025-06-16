#!/bin/bash

aws --endpoint-url=http://localhost:4566 secretsmanager create-secret \
  --name notes-database-secret \
  --secret-string '{"username":"postgres","password":"postgres","engine":"postgres","host":"localhost","port":5432,"dbname":"postgres","dbInstanceIdentifier":"notes"}' \
  --region us-east-2

aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name notes-queue --region us-east-2