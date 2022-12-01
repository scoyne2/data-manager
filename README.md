# Data Manager
Data manager is a platform that allows ingesting flat files into a datalake. It requires use of AWS and uses the following technology: terraform, helm, go, nextjs, react, typescript.
<img width="1728" alt="Data Manager Demo" src="docs/images/demo.png">


## Deploy Containers Locally
From the root of the project run: `make deploy_containers`

### Front End
In your browser navigate to http://localhost:3000/
TODO pass in graphql url as a variable, currently its hardcoded in frontend/src/pages/_app.tsx as "http://localhost/graphql"

### GraphQLAPI
In your browser navigate to http://localhost/sandbox

## Database Admin
In your browser navigate to http://localhost:5050/
Use dummy credentials username: admin@admin.com password: root


# Deploy infrastructure
Run `make deploy_infra`