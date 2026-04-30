# Repository Guidelines

## Project Structure & Module Organization

Clear Bill is a monorepo with a Go backend and a React admin frontend.

- `mgr/server`: Gin + GORM service. Entry point: `cmd/clearbill/main.go`.
- `mgr/webui`: React 18 + Vite + TypeScript UI. Source is in `src/`, static assets in `public/`.
- `scripts/`: bootstrap and packaging helpers.
- `tools/`: healthcheck and release-side utilities.
- Root `Makefile`: shared entry for build, test, lint, and release.

Keep backend layers aligned with the current structure: `action/`, `bll/`, `dal/`, and `dbmodel/`. Place tests beside the code they verify, for example `mgr/server/api/router/router_test.go`.

## Build, Test, and Development Commands

- `make server-start`: run the backend from the repository root.
- `make webui-dev`: start the frontend dev server.
- `make test`: run tests across submodules.
- `make lint`: run static checks across submodules.
- `make dist`: build backend and frontend artifacts.
- `make release`: package outputs into `release/clear-bill-<version>/`.
- `cd mgr/server && make swagger`: regenerate Swagger docs after API changes.
- `cd mgr/server && make wire`: regenerate dependency injection code.
- `cd mgr/webui && pnpm install && pnpm dev`: install dependencies and run the UI locally.
- `cd mgr/webui && pnpm api`: regenerate typed frontend API clients.

## Coding Style & Naming Conventions

Use `gofmt` defaults and tabs in Go files. Use 2-space indentation in TypeScript, PascalCase for React components, camelCase for functions, and `index.module.css` for page-local styles. Keep frontend service files grouped by domain under `mgr/webui/src/services/clear-bill/`. Prefer clear domain names over abbreviations and avoid editing generated files by hand.

## Testing Guidelines

Backend tests use Go standard tooling with `_test.go` filenames. Add focused tests for router, business, and data-access changes. Frontend validation currently relies on `pnpm lint` (`tsc --noEmit`), so UI changes should include manual verification notes in the pull request.

## Commit & Pull Request Guidelines

Recent history uses Conventional Commit prefixes such as `feat:` and `style:`, typically with concise Chinese subjects. Follow the same pattern for new commits. Pull requests should summarize scope, affected modules, commands run, and screenshots for UI changes. Note config changes and whether `swagger`, `wire`, or generated API files were updated.

## Security & Configuration Tips

Default backend config lives in `mgr/server/configs/config.toml`. Prefer environment overrides such as `CLEAR_BILL_CONFIG`, `CLEAR_BILL_HTTP_ADDR`, and `CLEAR_BILL_WEB_ROOT` over hardcoded local values. Never commit secrets, database passwords, or environment-specific endpoints.
