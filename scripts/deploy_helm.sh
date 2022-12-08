## run this from the root directory like: bash scripts/deploy_helm.sh

set -a
source .env
set +a

cd sysops/helm/data-manager && helm upgrade --install data-manager .  --set domainName=$DOMAIN_NAME --set hostedZoneId=$HOSTED_ZONE_ID \
    --values values.yaml --set api.image.repository=$AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com/data-manager-api \
    --set frontend.image.repository=$AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com/data-manager-frontend \
    --set postgresql.auth.postgresPassword=$POSTGRES_PASSWORD --set postgresql.auth.password=$POSTGRES_PASSWORD \
    --set postgresql.primary.initdb.password=$POSTGRES_PASSWORD --set postgresql.primary.initdb.user=$POSTGRES_USER  \
    --set api.env.postgres_password=$POSTGRES_PASSWORD \
    --set pgadmin4.env.email=$PGADMIN_DEFAULT_EMAIL --set pgadmin4.env.password=$PGADMIN_DEFAULT_PASSWORD \
    --set certificateARN=$ACM_ARN