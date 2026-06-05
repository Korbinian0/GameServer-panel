# GameServer Panel

Eine professionelle Server-Management-Plattform bestehend aus:

- Gateway (Go) — HTTP REST, WebSocket, gRPC-Registry, PostgreSQL-Persistenz
- Frontend (Vue 3 + TypeScript) — Dashboard, Live-Konsole, Server-Manager
- Datenbank (Postgres)

Architektur (vereinfacht):

- Browser → Frontend → Gateway (REST / WebSocket / gRPC) → Backend-Nodes

Wichtigste Pfade im Repository

- `backend/` — Go-Gateway (Quellcode, proto, Migrationen)
	- `backend/cmd/gateway` — Gateway Entrypoint
	- `backend/internal/adapters/http` — REST & WebSocket Handler
	- `backend/internal/adapters/websocket` — Hub / Client
	- `backend/internal/adapters/grpc` — gRPC-Registry & Clients
	- `backend/internal/adapters/migration/schema.sql` — DB-Schema
	- `backend/proto/gateway.proto` — gRPC-Service-Definitionen
- `frontend/` — Vue 3 App (Vite)
- `docker-compose.yml` — Entwickler-Docker-Stack (Postgres, Gateway, Frontend)

Kurzüberblick über vorhandene APIs

- REST
	- `POST /api/login` — Benutzeranmeldung (gibt JWT zurück)
	- `GET /api/health` — Healthcheck
	- `GET /api/nodes` — Liste registrierter Nodes
	- `GET /api/nodes/{id}` — Node-Details
	- `POST /api/nodes` — Node-Registrierung
	- `POST /api/nodes/{id}/heartbeat` — Heartbeat
	- `POST /api/users` — (geschützt) Benutzer anlegen
	- `GET /api/roles` — (geschützt) Rollen
- WebSocket
	- `GET /ws/events?token=...` — Live-Events / Broadcasts (JWT im Query erlaubt)
- gRPC
	- Gateway gRPC-Port: `:50051` (siehe `backend/proto/gateway.proto`)

Lokal ausführen (ohne Docker)

Voraussetzungen:

- Go >= 1.25
- Node.js & npm (für Frontend)
- PostgreSQL (z. B. lokal auf 5432)
- `protoc` + Go-Plugins falls du die gRPC-Clients/Server stützen willst

1) Datenbank anlegen & Migration ausführen

```bash
# Beispiel mit lokalem psql (ersetze URL/Dateipfade nach Bedarf)
psql "postgres://gateway:gatewaypass@localhost:5432/gateway?sslmode=disable" -f backend/internal/adapters/migration/schema.sql
```

2) TLS-Zertifikate für gRPC (lokale Entwicklung)

Das Gateway erwartet Pfade `GRPC_TLS_CERT` und `GRPC_TLS_KEY`. Erzeuge selbstsignierte Zertifikate für lokale Nutzung:

```bash
mkdir -p backend/certs
openssl req -x509 -newkey rsa:4096 -nodes -keyout backend/certs/server.key -out backend/certs/server.crt -days 365 -subj "/CN=localhost"
```

3) Gateway (Backend) starten

Setze Umgebungsvariablen (Beispiel) und starte das Gateway:

```bash
export DATABASE_URL="postgres://gateway:gatewaypass@localhost:5432/gateway?sslmode=disable"
export JWT_SECRET="changeme"
export GRPC_TLS_CERT="$(pwd)/backend/certs/server.crt"
export GRPC_TLS_KEY="$(pwd)/backend/certs/server.key"
cd backend
go run ./cmd/gateway
```

Alternativ builden und starten:

```bash
cd backend
go build -o gateway ./cmd/gateway
./gateway
```

4) Frontend lokal starten

```bash
cd frontend
npm install
npm run dev    # Entwicklung
# oder
npm run build   # Produktions-Build
```

Web UI (Dev): üblicherweise `http://localhost:4173` (Vite), im Docker-Setup ist es auf `:3000` gemappt.

gRPC-Code generieren (optional)

Wenn du gRPC-Clients/Server generieren willst, benutze protoc mit den Go-Plugins. Beispiel (ersetze Pfade):

```bash
protoc --go_out=. --go-grpc_out=. backend/proto/gateway.proto
```

Mit Docker (empfohlener schneller Start)

```bash
docker compose up --build
```

Wichtigste Dienste (Ports)

- Frontend: `3000` → mapped zu Vite `4173`
- Gateway REST / WebSocket: `8080`
- Gateway gRPC: `50051`
- Postgres: `5432`

WebSocket-Client-Verbindung (Beispiel)

```js
// Browser / frontend
const ws = new WebSocket(`ws://${window.location.host}/ws/events?token=${yourJwt}`)
```

Weitere Hinweise / Troubleshooting

- Wenn das Gateway beim Start TLS-Dateien verlangt, prüfe `GRPC_TLS_CERT` / `GRPC_TLS_KEY`.
- DB-Verbindungen: URL-Format `postgres://user:pass@host:5432/dbname?sslmode=disable`.
- Bei Problemen mit CORS / WebSocket-Proxys: prüfe, ob ein Reverse-Proxy die Header verändert.

Mitwirken / Entwicklung

- Tests und Lint: derzeit keine automatischen Tests im Repo, nutze `go test ./...` im `backend`-Ordner.
- Kontributionen: Fork → Branch → PR.

Lizenz

Dieses Projekt steht unter der Lizenz in der `LICENSE`-Datei.

---

Wenn du willst, ergänze ich noch eine `CONTRIBUTING.md` und ein kurzes Dev-Skript, das lokale Starts (DB, certs, env) automatisiert.
