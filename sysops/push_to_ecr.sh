## run this from the root directory like: bash sysops/push_to_ecr.sh

set -a
source .env
set +a

docker build -t data-manager-api data-manager-api/. --platform linux/amd64
docker tag data-manager-api:latest $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com/data-manager-api:latest

docker build -t data-manager-frontend frontend/. --platform linux/amd64
docker tag data-manager-frontend:latest $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com/data-manager-frontend:latest

$(aws ecr get-login --region us-west-2 --no-include-email --profile $AWS_PROFILE)

docker push $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com/data-manager-api:latest
docker push $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com/data-manager-frontend:latest