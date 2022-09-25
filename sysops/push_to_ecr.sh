## run this from the root director like: bash sysops/push_to_ecr.sh

docker build -t data-manager . --platform linux/amd64

docker tag data-manager:latest 852056369035.dkr.ecr.us-west-2.amazonaws.com/data-manager:latest

$(aws ecr get-login --region us-west-2 --no-include-email --profile personal)

docker push 852056369035.dkr.ecr.us-west-2.amazonaws.com/data-manager:latest