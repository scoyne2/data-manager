build_lambda:
	cd data-manager-file-notifications && GOARCH=amd64 GOOS=linux  go build -v -ldflags '-s -w' -a -tags netgo -installsuffix netgo -o ../data-manager-file-notifications/build/bin/app .

init:
	terraform init

plan:
	terraform plan

apply:
	terraform apply --auto-approve

destroy:
	terraform destroy --auto-approve