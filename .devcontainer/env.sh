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

# Install helm
brew install helm
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add runix https://helm.runix.net
helm repo add eks https://aws.github.io/eks-charts

# Install kubectl
brew install kubectl
alias k='kubectl'

# Terraform init
make init

# Build lambda
make build_lambda

# Install next js
cd frontend && npm install next && yarn install

# Set default for local testing of Go API
export API_HOST=localhost

# Install lambda and spark testing requirements
pip install boto3==1.26.109
pip install great-expectations==0.15.5
pip install pytest-cov==4.0.0
pip install pyspark==3.3.2
pip install PyYAML==6.0
pip install requests==2.26.0 
export APPLICATION_ID="999999999999"
export JOB_ROLE_ARN="arn:aws:iam::999999999999:role/fake"
export SCRIPT_LOCATION="test_script_location"
export OUTPUT_BUCKET="test_outputbucket"
export RESOURCE_BUCKET="test_resource_bucket"

# Add codecov
curl -Os https://uploader.codecov.io/latest/linux/codecov
chmod +x codecov
export CODECOV_TOKEN=${CODECOV_TOKEN}