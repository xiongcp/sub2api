# Repository Guidelines

## Project Structure & Module Organization
This repository is split into `backend/` and `frontend/`. The Go service entrypoint is `backend/cmd/server`, with business code under `backend/internal/{handler,service,repository,...}`, schema generation in `backend/ent`, SQL migrations in `backend/migrations`, and embedded web assets in `backend/internal/web`. The Vue app lives in `frontend/src`, organized by `api/`, `components/`, `composables/`, `stores/`, `views/`, and colocated `__tests__/`. Deployment assets are in `deploy/`, supporting scripts in `tools/`, and reference docs in `docs/`.

## Build, Test, and Development Commands
Run from the repository root unless noted:

- `make build`: build backend and frontend together.
- `make test`: run backend tests plus frontend lint and type checks.
- `make secret-scan`: scan staged content for leaked secrets.
- `cd backend && make build`: compile `bin/server` with version ldflags.
- `cd backend && make test-unit` / `make test-integration`: run tagged Go suites.
- `cd frontend && pnpm install`: install frontend dependencies. Use `pnpm`, not `npm`.
- `cd frontend && pnpm dev`: start the Vite dev server.
- `cd frontend && pnpm test:run` or `pnpm test:coverage`: run Vitest once, optionally with coverage.

## Coding Style & Naming Conventions
Follow existing style instead of reformatting unrelated code. Go code should remain `gofmt`-clean and package-oriented; keep tests in `*_test.go`. Vue and TypeScript use 2-space indentation, no semicolon style, PascalCase component files, and camelCase for composables/utilities, e.g. `UseKeyModal.vue`, `useRoutePrefetch.ts`. Frontend linting is enforced by `frontend/.eslintrc.cjs`; run `pnpm lint` before submitting UI-heavy changes.

## Testing Guidelines
Backend tests use Go's `testing` package, with `unit`, `integration`, and `e2e` tags defined in `backend/Makefile`. Frontend tests use Vitest and usually live in `src/**/__tests__/*.spec.ts`. Add or update tests for touched logic; there is no published coverage gate, but regressions in critical gateway, billing, auth, and admin UI paths should be covered explicitly.

## Commit & Pull Request Guidelines
Recent history favors Conventional Commit prefixes such as `fix(scope): ...`, `fix: ...`, and `chore: ...`. Keep commit subjects imperative and scoped when useful; use `[skip ci]` only for release/version housekeeping. PRs should include a short problem statement, the chosen fix, validation commands, and screenshots for frontend changes. Link related issues or migrations when applicable.

## Security & Configuration Tips
Do not commit real secrets, `.env` files, or local database dumps. Prefer `deploy/docker-compose.local.yml` for local stacks, and note that CI validates Go from `backend/go.mod` and frontend installs with `pnpm --frozen-lockfile`.
