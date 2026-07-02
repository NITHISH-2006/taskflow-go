# TaskFlow

TaskFlow is a minimal Go backend service that exposes a health check endpoint.

## Run the server

```bash
go run ./cmd/api
```

The API listens on `:8080`.

## Endpoints

### GET /health

Returns a status payload:

```json
{ "status": "ok" }
```
