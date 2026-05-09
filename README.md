# Casa da Vó — Reservation System

A full-stack restaurant reservation system with role-based access for customers and managers.

## Features

**Customers (`cliente`)**
- Register and log in
- Browse available dates and time slots
- Select a table from a visual map (filtered by party size and date availability)
- Create, view, and cancel reservations

**Managers (`gerente`)**
- Dashboard with pending/confirmed reservation counts and full reservation list
- Confirm or refuse pending reservations
- Manage tables (add, edit, activate/deactivate, delete)
- Manage availability: open dates, time slots, and auto-confirm toggle
- Table layout per date: block specific tables on specific days without affecting the default layout

## Tech Stack

| Layer | Tech |
|---|---|
| Backend | Go 1.23, [chi](https://github.com/go-chi/chi), JWT (`golang-jwt/jwt`) |
| Database | SQLite (via `go-sqlite3`, WAL mode, foreign keys on) |
| Frontend | SvelteKit + `adapter-static` |
| Web server | nginx (reverse-proxies `/api/` to backend) |
| Deployment | Docker Compose |

## Project Structure

```
casadavo/
├── backend/
│   ├── cmd/server/main.go          # Entry point, router setup
│   └── internal/
│       ├── handler/                # HTTP handlers
│       ├── middleware/             # JWT auth + role enforcement
│       ├── model/                  # Shared structs
│       ├── repository/             # SQLite queries
│       └── service/                # Business logic
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api.js              # Fetch wrapper
│   │   │   └── stores/auth.js      # Auth state (localStorage)
│   │   └── routes/
│   │       ├── login/
│   │       ├── register/
│   │       ├── reservations/       # Customer reservation list
│   │       ├── reserve/            # New reservation form
│   │       └── manager/
│   │           ├── +page.svelte    # Dashboard
│   │           ├── tables/         # Table management
│   │           ├── availability/   # Date & slot management
│   │           └── layout/         # Per-date table exclusions
│   └── nginx.conf
└── docker-compose.yml
```

## Running with Docker

The easiest way to run the app. Images are published on Docker Hub.

**1. Create a `.env` file:**
```env
JWT_SECRET=your_long_random_secret
```

**2. Start:**
```bash
docker compose up -d
```

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api

The SQLite database is persisted in a named Docker volume (`sqlite_data`).

## Running Locally (Development)

**Backend:**
```bash
cd backend
JWT_SECRET=dev-secret go run ./cmd/server
# API at http://localhost:8080
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
# UI at http://localhost:5173
```

> In dev mode the frontend calls `/api` relative to itself. Configure a proxy in `vite.config.js` pointing to `http://localhost:8080` if needed.

## API Overview

All endpoints (except auth) require `Authorization: Bearer <token>`.

| Method | Path | Role | Description |
|---|---|---|---|
| POST | `/api/auth/register` | public | Create account |
| POST | `/api/auth/login` | public | Get JWT token |
| GET | `/api/auth/me` | any | Current user info |
| GET | `/api/tables` | any | List tables |
| POST/PUT/DELETE | `/api/tables/{id}` | gerente | Manage tables |
| GET | `/api/availability` | any | List open dates |
| GET | `/api/availability/{id}/slots` | any | List time slots |
| POST/PUT | `/api/availability/{id}` | gerente | Manage dates |
| POST/DELETE | `/api/availability/{id}/slots` | gerente | Manage slots |
| GET | `/api/layout?date=YYYY-MM-DD` | any | Table exclusions for a date |
| POST | `/api/layout` | gerente | Block a table on a date |
| DELETE | `/api/layout/{id}` | gerente | Unblock a table |
| GET | `/api/reservations` | any | List reservations (own; all for gerente) |
| POST | `/api/reservations` | any | Create reservation |
| PUT | `/api/reservations/{id}` | any | Update reservation |
| DELETE | `/api/reservations/{id}` | any | Cancel reservation |
| PUT | `/api/reservations/{id}/status` | gerente | Confirm or refuse |

## Docker Hub

```
pokkew/casadavo-backend:latest
pokkew/casadavo-frontend:latest
```
