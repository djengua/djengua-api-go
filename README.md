# Portal API (Golang) + Prometheus

API REST en Go basada en la definición de entidades y endpoints del documento `API_ENTIDADES.md` (Productos, Categorías y Colecciones), con **métricas Prometheus en `/metrics`**.

## Características

- Base URL: `/api/v1`
- Recursos:
  - `/products`
  - `/categories`
  - `/collections`
- CRUD completo (GET list/get, POST, PUT, PATCH, DELETE con soft-delete)
- PostgreSQL (con migraciones SQL)
- Observabilidad:
  - `GET /healthz`
  - `GET /metrics` (Prometheus)
    - `http_requests_total{method,route,status}`
    - `http_request_duration_seconds{method,route}`

## Requisitos

- Go 1.22+
- Docker + Docker Compose (recomendado para correr todo local)
- PostgreSQL 16 (si corres sin Docker)

## Estructura del proyecto

```
.
├── cmd/api/main.go
├── internal/
│   ├── config/
│   ├── db/
│   ├── domain/
│   └── httpapi/
│       ├── handlers/
│       ├── middleware/
│       └── router/
├── migrations/
├── deploy/
│   ├── docker-compose.yml
│   └── prometheus.yml
├── Dockerfile
└── Makefile
```

## Quick start (Docker)

1) Levanta Postgres + API + Prometheus:

```bash
make docker-up
# o
docker compose -f deploy/docker-compose.yml up --build
```

2) Prueba salud y métricas:

```bash
curl -i http://localhost:8080/healthz
curl -s http://localhost:8080/metrics | head
```

3) Abre Prometheus:

- Prometheus: http://localhost:9090
- Target esperado: `portal-api` scrapeando `api:8080/metrics`

> Las migraciones se ejecutan automáticamente al iniciar `postgres` porque `migrations/` se monta en `/docker-entrypoint-initdb.d`.

## Configuración por variables de entorno

- `PORT` (default: `8080`)
- `DATABASE_URL` (default: `postgres://postgres:postgres@localhost:5432/portal?sslmode=disable`)

Ejemplo:

```bash
export PORT=8080
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/portal?sslmode=disable'
```

## Endpoints

> Prefijo: `/api/v1`

### Productos

- `GET /products` (paginado + filtros)
  - `page`, `page_size`
  - `status`, `category_id`, `collection_id`, `min_price`, `max_price`, `q`
- `GET /products/{id}`
- `POST /products`
- `PUT /products/{id}`
- `PATCH /products/{id}`
- `DELETE /products/{id}` (soft-delete)

Ejemplo **crear**:

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H 'Content-Type: application/json' \
  -d '{
    "sku":"ROS-001",
    "name":"Ramo de rosas rojas",
    "description":"Ramo premium de 24 rosas rojas.",
    "price":120.0,
    "currency":"USD",
    "status":"active",
    "category_id":"<UUID_CATEGORIA>",
    "collection_ids":["<UUID_COLECCION>"]
  }'
```

### Categorías

- `GET /categories`
- `GET /categories/{id}`
- `POST /categories`
- `PUT /categories/{id}`
- `PATCH /categories/{id}`
- `DELETE /categories/{id}` (soft-delete)

### Colecciones

- `GET /collections`
- `GET /collections/{id}`
- `POST /collections`
- `PUT /collections/{id}`
- `PATCH /collections/{id}`
- `DELETE /collections/{id}` (soft-delete)

## Métricas Prometheus

La API expone métricas en `GET /metrics`. De forma predeterminada:

- `http_requests_total{method,route,status}`: contador por método/ruta/código.
- `http_request_duration_seconds{method,route}`: histograma de latencias.

Ejemplos de consultas en Prometheus:

- QPS por ruta:

```
sum(rate(http_requests_total[5m])) by (route)
```

- Latencia p95 por ruta:

```
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, route))
```

## Guía para crear el proyecto desde cero (paso a paso)

Si quieres recrearlo en tu repo desde una carpeta vacía:

```bash
mkdir portal-api && cd portal-api

go mod init <tu-dominio>/<tu-proyecto>

# dependencias
# (en tu máquina, con internet)
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5
go get github.com/google/uuid
go get github.com/prometheus/client_golang

# pega la estructura de carpetas y el código
# luego:
go fmt ./...
go test ./...
```

Después puedes incorporar `deploy/docker-compose.yml` para levantar todo.

## Notas

- El proyecto usa **soft-delete** (`deleted_at`) y no elimina físicamente los registros.
- Los campos de fecha se manejan en UTC.
- En producción, agrega autenticación/autorización para endpoints de escritura.
