# Lab-VerteilteSysteme

## Docker Container starten 
```bash
docker-compose build
docker-compose up
```

## Projektstruktur
- **frontendFlask**: Flask-Frontend mit einem HTML-Template zur Interaktion mit dem Shop.
- **shop**: Go Backend basierend auf Fiber mit einer PostgreSQL als Datenbank.
- **docker-compose.yml**: Koordiniert den Start aller Services (Backend, Frontend, Datenbank).

## Verteilte Systeme
Dieses Projekt besteht aus mehreren eigenständigen Komponenten, die jeweils in einem eigenen Docker-Container laufen und über ein gemeinsames Docker-Netzwerk miteinander kommunizieren.
Das Projekt besteht aus folgenden Komponenten:
- Frontend-Service (Flask)
- Backend-Service (Go Fiber API)
- Datenbank-Service (PostgreSQL)
Diese Trennung erlaubt unabhängiges Deployen, Skalieren und Entwickeln einzelner Komponenten.

## 12-Factor App Methodology
- **1. Codebase**: Eine einzige Codebasis pro Anwendung in Github, versioniert mit Git.
- **2. Dependencies**: Klare Deklaration der Abhängigkeiten über go.mod für dass Go Backend und requirements.txt (bei Python).
- **3. Config**: Konfigurationswerte (z.B. API-Keys/Passwörter, Ports) werden über Umgebungsvariablen bereitgestellt – gesteuert durch Docker.
- **4. Backing Services**: Die Datenbank (PostgreSQL) wird als externer, entkoppelter Service behandelt und kann problemlos ausgetauscht werden.
- **5. Build, Release, Run**: Build-Prozess, Release-Erstellung und Ausführung sind klar voneinander getrennt – umgesetzt durch Docker-Container.
- **6. Processes**: Die Anwendung läuft in zustandslosen (stateless) Prozessen, wodurch horizontale Skalierung problemlos möglich ist.
- **7. Port Binding**: Services binden an Ports und werden via HTTP verfügbar gemacht.
- **8. Concurrency**: Services können parallel über mehrere Container oder Prozesse betrieben werden.
- **9. Disposability**: Container starten schnell, können jederzeit beendet werden und hinterlassen keinen Zustand – ideal für automatische Deployments.
- **10. Dev/Prod Parity**: Entwicklung, Test und Produktion laufen in weitgehend identischen Umgebungen dank Containerisierung.
- **11. Logs**: Logs werden auf STDOUT/STDERR ausgegeben und zentral über Docker verarbeitet – keine Log-Dateien im Container selbst.
- **12. Admin Processes**: Administrative Tasks (wie DB-Migrationen) könnten als einmalige Jobs gestartet werden.
