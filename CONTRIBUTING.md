# Contributing

Danke für dein Interesse, zum GameServer Panel beizutragen! Diese Datei beschreibt kurz, wie du lokal entwickelst, Tests ausführst und Pull Requests vorbereitest.

Kurz: Fork → Feature-Branch → Commit → Pull Request

Lokale Entwicklung

- Lies die Haupt-README (`README.md`) für Startanweisungen.
- Es gibt ein Hilfs-Skript `scripts/dev.sh`, das lokale Dev-Setup (Zertifikate, optionale DB-Migration) und die Starts für Backend/Frontend automatisiert.

Code-Stil und Tests

- Backend (Go)
  - Formatiere Code mit `gofmt` bzw. `go fmt`.
  - Führe `go test ./...` im `backend/`-Verzeichnis aus.
- Frontend (Vue)
  - Nutze `npm run build` zum Überprüfen des Typescript- und Bundling-Status.

Pull Request Workflow

1. Forke das Repo und erstelle einen Feature-Branch vom `main`-Branch.
2. Schreibe kleine, fokussierte Commits mit erklärenden Nachrichten.
3. Stelle sicher, dass `go test ./...` und `npm run build` lokal laufen.
4. Öffne einen Pull Request mit einer kurzen Beschreibung, was geändert wurde und warum.

Code of Conduct

Bitte verhalte dich respektvoll. Bei Unklarheiten diskutiere Änderungen in der PR-Beschreibung.

Kontakt

Wenn du größere Änderungen planst (z. B. API-Änderungen), eröffne vorher ein Issue, damit wir das Design abstimmen können.
