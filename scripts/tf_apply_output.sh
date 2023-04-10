sed -i -e "s|.*ACM_ARN.*|ACM_ARN=FAKEACM|" ../.env
sed -i -e "s|.*WAF_ARN.*|WAF_ARN=FAKEWAF|" ../.env

terraform -chdir="../sysops/terraform" apply --auto-approve
ACM_ARN=$(terraform -chdir="../sysops/terraform" output -raw aws_acm_certificate_arn)
sed -i -e "s|FAKEACM|$ACM_ARN|g" ../.env

WAF_ARN=$(terraform -chdir="../sysops/terraform" output -raw aws_wafv2_web_acl_arn)
sed -i -e "s|FAKEWAF|$WAF_ARN|g" ../.env