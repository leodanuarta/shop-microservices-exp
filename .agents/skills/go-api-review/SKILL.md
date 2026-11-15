---
name: go-api-review
description: Review or modify Go HTTP handlers in this repo. Checks HTTP method, status codes, input validation, error handling, tests, and file-change scope.
---

# go-api-review

1. Read `docs/SPEC.md` first; the acceptance criteria there are the contract.
2. For each handler check: HTTP method/route, status codes, `Content-Type: application/json`, input validation, error responses via `httpx.Error`.
3. Keep stdlib only; do not add dependencies.
4. Every behavior change needs a table-driven test (success, validation error, empty result).
5. Only touch files relevant to the task; list them at the end.
6. Run `go vet ./... && go test ./...` and report the real output. Never claim tests passed without running them.
