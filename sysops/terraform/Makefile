build_lambda:
	cd data-manager-file-notifications && GOARCH=amd64 GOOS=linux  go build -v -ldflags '-s -w' -a -tags netgo -installsuffix netgo -o ../data-manager-file-notifications/build/bin/app .

init:
	cd data-manager-file-notifications && terraform init

plan:
	cd data-manager-file-notifications && terraform plan

apply:
	cd data-manager-file-notifications && terraform apply --auto-approve

destroy:
	cd data-manager-file-notifications && terraform destroy --auto-approve