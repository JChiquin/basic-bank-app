# AGENTS.md

Educational bank app: two independent apps, no root workspace or shared tooling. Run all commands from inside the relevant app directory.

- `backend/` — Go 1.16 REST API (`bank-service`), gorilla/mux + GORM + PostgreSQL 12. See `backend/AGENTS.md`.
- `frontend/` — React 18 (Create React App, JavaScript, MUI 5 + Redux Toolkit). See `frontend/AGENTS.md`.

## Cross-app facts

- Both apps are developed via their own `docker-compose.yml`; each needs a `.env` copied from its `.env.example` first.
- Backend compose requires a pre-created external network: `docker network create bank-net` (once per machine).
- Ports: backend API `:3000`, Postgres `:5432`, pgAdmin `:80`, frontend `:8080`. The frontend must run on 8080 (backend CORS `WHITE_LIST`); a bare CRA `npm start` defaults to 3000 and collides with the backend.
- Frontend talks to backend via `REACT_APP_API_URL_BACKEND` (default `http://localhost:3000`). Auth is JWT sent as `Bearer` header.
- API reference: Postman collection `Bank-service` (linked in root README).