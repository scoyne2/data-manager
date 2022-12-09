sed -i -e "s|.*ACM_ARN.*|ACM_ARN=FAKE_ARN|" ../.env
terraform -chdir="../sysops/terraform" apply --auto-approve
ACM_ARN=$(terraform -chdir="../sysops/terraform" output -raw aws_acm_certificate_arn)
sed -i -e "s|FAKE_ARN|$ACM_ARN|g" ../.env
rm ../.env-e