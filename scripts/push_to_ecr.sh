## run this from the root directory like: bash scripts/push_to_ecr.sh
unset NEXT_PUBLIC_DOMAIN_NAME
unset NEXT_PUBLIC_S3_RESOURCE_BUCKET
unset NEXT_PUBLIC_AWS_REGION
unset DOMAIN_NAME

set -a
source .env
set +a

# AWS CLI v2
aws ecr get-login-password --region $AWS_REGION --profile $AWS_PROFILE | docker login --username AWS --password-stdin $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com

# AWS CLI v1
# $(aws ecr get-login-password --region $AWS_REGION --no-include-email --profile $AWS_PROFILE)

docker build -t data-manager-api api/. --platform linux/amd64
docker tag data-manager-api:latest $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/data-manager-api:latest

docker build -t data-manager-sla sla/. --platform linux/amd64
docker tag data-manager-sla:latest $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/data-manager-sla:latest

docker build -t data-manager-frontend frontend/. --platform linux/amd64 --build-arg NEXT_PUBLIC_DOMAIN_NAME=${DOMAIN_NAME} --build-arg NEXT_PUBLIC_S3_RESOURCE_BUCKET=${RESOURCES_BUCKET_NAME}  --build-arg NEXT_PUBLIC_AWS_REGION=${AWS_REGION}   
docker tag data-manager-frontend:latest $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/data-manager-frontend:latest

docker push $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/data-manager-api:latest
docker push $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/data-manager-sla:latest
docker push $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/data-manager-frontend:latest