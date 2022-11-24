# Data Manager

## Deploy Locally
From the root of the project run: `docker-compose up`

### Front End
In your browser navigate to http://localhost:3000/

### GraphQLAPI
In your browser navigate to http://localhost:8080/sandbox

### Database
To access the db once its running:
Get the container ID: `docker ps -f "name=postgres"`
Access the container, replaceing CONTAINER_ID: `docker exec -it CONTAINER_ID /bin/bash`
Bring up psql commands: `psql -U postgres`

## Sysops
TODO add documentation



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