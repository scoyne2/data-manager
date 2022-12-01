ifneq (,$(wildcard ./.env))
    include .env
    export
endif

deploy_helm:
	bash scripts/push_to_ecr.sh
	bash scripts/update_kubeconfig.sh
	bash scripts/deploy_helm.sh
