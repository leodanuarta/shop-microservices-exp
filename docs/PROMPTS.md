# Contoh prompt untuk coding agent

```
Context:
- Go microservices (product-service, order-service), stdlib only. Keep current structure and dependencies.

Task:
- Add GET /orders/{id} to order-service.

Acceptance criteria:
- 200 + order JSON if found, 404 + {"error": "..."} if not.
- Add tests for both cases.

Constraints:
- Do not change unrelated files.
- Report changed files and actual test results.
- Do not claim tests passed unless you ran them.
```
