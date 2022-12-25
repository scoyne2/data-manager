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

build_lambda:
	@echo '${GREEN}Building Lambda${RESET}'
	cd sysops/terraform/lambda && GOARCH=amd64 GOOS=linux  go build -v -ldflags '-s -w' -a -tags netgo -installsuffix netgo -o ../build/bin/app .

init:
	@echo '${GREEN}Terraform init${RESET}'
	cd sysops/terraform && terraform init

plan:
	@echo '${GREEN}Planning Terraform${RESET}'
	cd sysops/terraform && terraform plan

apply:
	@echo '${GREEN}Applying Terraform${RESET}'
	cd scripts && bash tf_apply_output.sh

destroy:
	@echo '${GREEN}Destroying Terraform${RESET}'
	cd sysops/terraform &&  terraform destroy --auto-approve