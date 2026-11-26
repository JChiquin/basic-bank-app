# AGENTS.md — frontend (bank-client)

Create React App (react-scripts 5), React 18, plain JavaScript (no TypeScript), MUI 5 + Redux Toolkit + axios. Docker image uses Node 16 — keep dependencies compatible. Never run `npm run eject`.

## Dev workflow

Preferred (Docker, from this directory):

```bash
cp .env.example .env   # REACT_APP_API_URL_BACKEND=http://localhost:3000, PORT=8080
docker-compose build && docker-compose up -d   # app at http://localhost:8080
```

- `node_modules` lives in the named volume `reserved`, not the bind mount. Installing/updating a dependency requires `docker-compose build` (npm i runs at image build) — `npm install` on the host won't reach the container.
- Env vars must be prefixed `REACT_APP_` and are read at dev-server start; restart the container/server after changing `.env`.
- Without Docker: `npm start`, but set `PORT=8080` — the default 3000 collides with the backend, and the backend CORS whitelist expects 8080.

## Commands

- No lint script: ESLint runs via the `eslintConfig` in `package.json` (`react-app` preset) during start/build.
- Tests: `npm test` is interactive watch mode. Single run: `CI=true npm test`. Single file: `CI=true npm test -- src/App.test.js`. Only a CRA smoke test exists; there is no meaningful coverage.
- Production build: `npm run build` → `build/`.

## Architecture / conventions

- ALL HTTP goes through `apiHttp()` in `src/api/axiosApi.js`: it attaches the JWT from localStorage (`src/utils/localStorage`) as `Bearer` header and **never throws** — on HTTP errors it returns the error response body. Callers must inspect the returned body for errors, not use try/catch.
- Endpoint wrappers per domain in `src/api/modules/` (user, movement); add new endpoints there, not ad-hoc axios calls.
- State: Redux Toolkit slices in `src/redux/`, store in `src/redux/store.js`.
- UI: page components in `src/components/`, reusable pieces in `src/components/prebuilt/`; MUI + notistack for notifications.