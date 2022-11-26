ifneq (,$(wildcard ./.env))
    include .env
    export
endif

deploy_infra:
	cd cdk/infra && TRIGGER_BUCKET_NAME=${TRIGGER_BUCKET_NAME} cdk deploy

deploy_containers:
	docker-compose up