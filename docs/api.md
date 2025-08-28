# Data Manager API

The below details how the Data Manager API functions.

## Running Locally

Export required environment variables:
```bash
export API_HOST=localhost
export FRONT_END_URL=localhost
export RESOURCES_BUCKET_NAME=<your-resources-bucket>
```
From the project root, run:
```bash
go run api/main.go
```
Access the sandbox UI on your browser via:
```
http://localhost:8080/sandbox
```
As there will not be a database running, expect that queries may fail with a TCP connection error.

## Environment Variables

- **API_HOST**: Hostname for the API (e.g., `api.example.com`)
- **FRONT_END_URL**: Hostname for the frontend (e.g., `app.example.com`)
- **RESOURCES_BUCKET_NAME**: S3 bucket name used for data preview resources

## HTTP Endpoints

- **GET /** and **GET /health_check**  
  Health check endpoint; returns a simple HTML response.

- **GET /preview**  
  Data Preview REST API. Query parameters:
  - `vendor` (required): Vendor name (URL-encoded, spaces replaced with `+`)
  - `feedname` (required): Feed name
  - `filename` (required): File name (including extension)  
  Example:  
  ```
  GET https://localhost:8080/preview?vendor=Coyne+Enterprises&feedname=Orders&filename=orders.csv
  ```
  Returns `201 Created` and a JSON payload of preview results.

- **POST /graphql**  
  GraphQL API endpoint. Accepts JSON bodies with a `query` field. CORS is enabled for the configured frontend origin.

- **GET /sandbox**  
  Embeddable GraphQL sandbox UI (Apollo Studio Explorer).

## GraphQL Schemas

_(See `api/schemas/factory.go` for type definitions.)_

## Resolvers

_(See `api/feed/resolver.go` and `api/datapreview/preview.go` for resolver implementations.)_

## PostgreSQL Layer

_(See `api/feed/postgress.go` for the Postgres service and repository implementations.)_