build_lambda:
	cd lambda && GOARCH=amd64 GOOS=linux  go build -v -ldflags '-s -w' -a -tags netgo -installsuffix netgo -o ../build/bin/app .

init:
	terraform init

plan:
	terraform plan

apply:
	terraform apply --auto-approve

destroy:
	terraform destroy --auto-approve