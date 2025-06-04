# FROM golang:1.24

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN go build -o app

# CMD ["./app"]

# === STAGE 1: Build ===
FROM golang:1.24 AS builder

WORKDIR /app

# Nur Abhängigkeiten zuerst kopieren (für Layer-Caching)
COPY go.mod go.sum ./
RUN go mod download

# Quellcode kopieren und kompilieren
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# === STAGE 2: Runtime ===
FROM alpine:latest

WORKDIR /app

# Nur das Binary vom Builder übernehmen
COPY --from=builder /app/app .

# Port freigeben (optional)
# EXPOSE 8080

# Binary ausführen
CMD ["./app"]
