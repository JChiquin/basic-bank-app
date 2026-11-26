# AGENTS.md

Educational bank app: two independent apps, no root workspace or shared tooling. Run all commands from inside the relevant app directory.

- `backend/` — Go 1.23 REST API (`bank-service`), gorilla/mux + GORM + PostgreSQL 16. See `backend/AGENTS.md`.
- `frontend/` — React 18 (Create React App, JavaScript, MUI 5 + Redux Toolkit). See `frontend/AGENTS.md`.

## Cross-app facts

- Both apps are developed via their own `docker-compose.yml`; each needs a `.env` copied from its `.env.example` first. Use `docker compose` (v2) syntax.
- Ports: backend API `:3000`, frontend `:8080`. Postgres is not published to the host by default (uncomment its mapping to use a local DB client); pgAdmin is opt-in (`docker compose --profile tools up -d pgadmin`, `:5050`). The frontend must run on 8080 (backend CORS `WHITE_LIST`); a bare CRA `npm start` defaults to 3000 and collides with the backend.
- Backend migrations run automatically when its app container starts (no manual step, no external docker networks).
- Frontend talks to backend via `REACT_APP_API_URL_BACKEND` (default `http://localhost:3000`). Auth is JWT sent as `Bearer` header.
- API reference: Postman collection `Bank-service` (linked in root README).