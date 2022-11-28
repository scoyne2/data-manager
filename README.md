# Data Manager

## Deploy Containers Locally
From the root of the project run: `make deploy_containers`

### Front End
In your browser navigate to http://localhost:3000/
TODO pass in graphql url as a variable, currently its hardcoded in frontend/src/pages/_app.tsx as "http://localhost:8080/graphql"

### GraphQLAPI
In your browser navigate to http://localhost/sandbox

## Database Admin
In your browser navigate to http://localhost:5050/
Use dummy credentials username: admin@admin.com password: root

## Sysops
TODO add documentation

# Deploy infrastructure
Run `make deploy_infra`

# Future Work
Implement postgres resolver for GraphQL
Connect front end to use GraphQL for Feed CRUD

Stub out api integrations for:
Execute source setup
Execute destination setup
File received
Read file status
Read logs
Sla check

Implement features for:
Data preview
Data quality