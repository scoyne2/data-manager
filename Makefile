ifneq (,$(wildcard ./.env))
    include .env
    export
endif

deploy_infra:
	cd cdk/infra && TRIGGER_BUCKET_NAME=${TRIGGER_BUCKET_NAME} cdk deploy

deploy_containers:
	docker-compose up

deploy_helm:
	bash sysops/push_to_ecr.sh
	cd sysops/helm/data-manager && helm upgrade --install data-manager \
	 . --values values.yaml --set api.image.repository=${AWS_ACCOUNT}.dkr.ecr.us-west-2.amazonaws.com/data-manager-api \
	 --set frontend.image.repository=${AWS_ACCOUNT}.dkr.ecr.us-west-2.amazonaws.com/data-manager-frontend \
	 --set postgresql.auth.postgresPassword=${POSTGRES_PASSWORD} --set postgresql.auth.password=${POSTGRES_PASSWORD} \
	 --set postgresql.primary.initdb.password=${POSTGRES_PASSWORD} --set postgresql.primary.initdb.user=${POSTGRES_USER}