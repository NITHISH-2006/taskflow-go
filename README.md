# TaskFlow

TaskFlow is a minimal Go backend service for task management.

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

## Domain models

The service defines core task models and request payloads:

- `Task`
- `CreateTaskRequest`
- `UpdateTaskRequest`
- `TaskResponse`

## Validation helpers

Validation utilities cover:

- required title values
- priority bounds between 1 and 5
- allowed status values
- RFC3339 timestamp parsing
