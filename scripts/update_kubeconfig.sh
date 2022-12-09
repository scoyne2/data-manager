## run this from the root directory like: bash scripts/update_kubeconfig.sh

set -a
source .env
set +a

aws eks --region $AWS_REGION update-kubeconfig --name data-manager-eks --profile $AWS_PROFILE
