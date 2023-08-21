## run this from the root directory like: bash scripts/deploy_helm.sh

set -a
source .env
set +a

# Setup datadog
kubectl delete secret datadog
kubectl create secret generic datadog --from-literal api-key=$DATADOG_API_KEY --from-literal app-key=$DATADOG_APP_KEY
cd sysops/helm/data-manager && helm dependency build && helm upgrade --install datadog-agent -f datadog-values.yaml --set targetSystem=linux datadog/datadog

# Setup datamanager deployment
helm dependency build && helm upgrade --install data-manager . --set domainName=$DOMAIN_NAME --set hostedZoneId=$HOSTED_ZONE_ID \
    --values values.yaml --set api.image.repository=$AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/data-manager-api \
    --set frontend.image.repository=$AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/data-manager-frontend \
    --set sla.image.repository=$AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/data-manager-sla \
    --set postgresql.auth.postgresPassword=$POSTGRES_PASSWORD --set postgresql.auth.password=$POSTGRES_PASSWORD \
    --set postgresql.primary.initdb.password=$POSTGRES_PASSWORD --set postgresql.primary.initdb.user=$POSTGRES_USER  \
    --set api.env.postgres_password=$POSTGRES_PASSWORD --set api.env.postgres_user=$POSTGRES_USER \
    --set pgadmin4.env.email=$PGADMIN_DEFAULT_EMAIL --set pgadmin4.env.password=$PGADMIN_DEFAULT_PASSWORD \
    --set certificateARN=$ACM_ARN --set wafARN=$WAF_ARN --set awsRegion=$AWS_REGION  \
    --set pgadmin4.serverDefinitions.servers.dataManagerServer.Username=$POSTGRES_USER \
    --set pgadmin4.ingress.annotations."external-dns\.alpha\.kubernetes\.io\/hostname"=pgadmin4.$DOMAIN_NAME \
    --set resourceBucketName=$RESOURCES_BUCKET_NAME 