ifneq (,$(wildcard ./.env))
    include .env
    export
endif

deploy:
	bash scripts/push_to_ecr.sh
	bash scripts/update_kubeconfig.sh
	bash scripts/deploy_helm.sh

build_lambda:
	cd sysops/terraform/lambda && GOARCH=amd64 GOOS=linux  go build -v -ldflags '-s -w' -a -tags netgo -installsuffix netgo -o ../build/bin/app .

init:
	cd sysops/terraform && terraform init

plan:
	cd sysops/terraform && terraform plan

apply:
	cd scripts && bash tf_apply_output.sh

destroy:
	cd sysops/terraform &&  terraform destroy --auto-approve