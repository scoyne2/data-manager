## run this from the root directory like: bash scripts/push_to_ecr.sh

set -a
source .env
set +a

# AWS CLI v2
aws ecr get-login-password --region us-west-2 --profile $AWS_PROFILE | docker login --username AWS --password-stdin 852056369035.dkr.ecr.us-west-2.amazonaws.com

# AWS CLI v1
# $(aws ecr get-login-password --region us-west-2 --no-include-email --profile $AWS_PROFILE)

docker build -t data-manager-api api/. --platform linux/amd64
docker tag data-manager-api:latest $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com/data-manager-api:latest

docker build -t data-manager-frontend frontend/. --platform linux/amd64 --build-arg NEXT_PUBLIC_DOMAIN_NAME=${DOMAIN_NAME}
docker tag data-manager-frontend:latest $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com/data-manager-frontend:latest

docker push $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com/data-manager-api:latest
docker push $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com/data-manager-frontend:latest