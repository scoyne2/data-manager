#!/bin/bash

# Setup AWS credentials and config
mkdir -p ~/.aws/
touch ~/.aws/credentials
touch ~/.aws/config

cat << EOF > ~/.aws/credentials
[default]
aws_access_key_id = ${AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${AWS_SECRET_ACCESS_KEY}
EOF

cat << EOF > ~/.aws/config
[default]
region = us-west-2
EOF

# Install next js
cd frontend && npm install next && yarn install

# Set default port for local testing of Go API
export POSTGRES_PORT=5432
export API_HOST=localhost