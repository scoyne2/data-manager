# Data Manager
Data manager is a platform that allows ingesting flat files into a datalake. It requires use of AWS and uses the following technology: terraform, helm, go, nextjs, react, typescript.
<img width="1728" alt="Data Manager Demo" src="docs/images/demo.png">

## Getting Started
1. To start you must had a domain name with AWS Route 53, and have an active Hosted Zone.
Take the Hosted Zone Id and update ``sysops/helm/data-manager/values.yaml``, update ``hostedZoneId`` with your ID. TODO: move this to the .env file
2. Rename ``.env_template`` to ``.env`` and update all necessary fields. In particular DOMAIN_NAME must align with the domain name you own via Route 53, AWS_ACCOUNT must be the account ID youwill be setting up resources in. AWS_PROFILE must be the name of the profile in your local ~/.aws/credentials that you will use to create resources (It will need access to all of the systems that Terraform will build)
3. From ``sysops/terraform`` run command ``make init`` then build the lambda with ``make build_lambda`` then plan the creation of infra  ``make plan`` and finally ``make apply`` to build the infra.
4. From project root build and push ECR images by running `` ``
5. Update kubeconfig by running ``aws eks --region us-west-2 update-kubeconfig --name data-manager-eks --profile personal`` change your AWS profile name as needed.
6. From project root deploy helm chart by running ``bash scripts/deploy_helm.sh``



### Front End
In your browser navigate to http://localhost:3000/
TODO pass in graphql url as a variable, currently its hardcoded in frontend/src/pages/_app.tsx as "http://localhost/graphql"

### GraphQLAPI
In your browser navigate to http://localhost/sandbox

## Database Admin
In your browser navigate to http://localhost:5050/
Use dummy credentials username: ``admin@admin.com`` password: ``admin``
