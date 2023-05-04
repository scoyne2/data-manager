GREEN  := $(shell tput -Txterm setaf 2)
RESET  := $(shell tput -Txterm sgr0)

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

deploy:
	@echo '${GREEN}Building Images and pushing to ECR${RESET}'
	@bash scripts/push_to_ecr.sh
	@echo '${GREEN}Updating Kube Config${RESET}'
	@bash scripts/update_kubeconfig.sh
	@echo '${GREEN}Deploying Helm Chart${RESET}'
	@bash scripts/deploy_helm.sh

init:
	@echo '${GREEN}Terraform init${RESET}'
	cd sysops/terraform && terraform init

plan:
	 @echo '${GREEN}Building EMR Python Deps${RESET}'
	cd pyspark/requirements && DOCKER_BUILDKIT=1 docker build --file Dockerfile --output out .
	@echo '${GREEN}Planning Terraform${RESET}'
	cd sysops/terraform && terraform plan

apply:
	@echo '${GREEN}Building EMR Python Deps${RESET}'
	cd pyspark/requirements && DOCKER_BUILDKIT=1 docker build --file Dockerfile --output out .
	@echo '${GREEN}Applying Terraform${RESET}'
	cd scripts && bash tf_apply_output.sh

destroy:
	@echo '${GREEN}Destroying Terraform${RESET}'
	cd sysops/terraform &&  terraform destroy --auto-approve

test:
	@echo '${GREEN}Running Tests${RESET}'
	cd sysops/terraform/lambda/ && pytest --cov --cov-report lcov
	cd pyspark/ && pytest --cov --cov-report lcov
	cd api && go test ./... -coverprofile=coverage.txt
	./codecov
