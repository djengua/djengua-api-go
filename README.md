# Portal API (Golang)

API REST en Go para un portal de e‑commerce. Expone productos, categorías, colecciones, órdenes y ventas, con autenticación JWT y métricas Prometheus en `/metrics`. Este README está pensado para servir como **fuente de verdad** para construir un front (por ejemplo el mock de `index.html`) usando **React + Next.js**, **Zustand** y **Zod**.

## Funcionalidad principal

- **Catálogo**: CRUD de productos, categorías y colecciones.
- **Autenticación**: registro, login y perfil (`/auth/*`).
- **Órdenes**: crear órdenes y listar las órdenes del usuario autenticado.
- **Ventas**: registrar ventas asociadas a una orden.
- **Observabilidad**: healthcheck y métricas para Prometheus.

## Base URL y rutas globales

- Base URL: `/api/v1`
- Healthcheck: `GET /healthz`
- Métricas: `GET /metrics`

## Requisitos

- Go 1.22+
- MongoDB (local o Docker)

### Quick start (local)

1) Levanta MongoDB (ejemplo con Docker):

```bash
docker run --rm -p 27017:27017 --name portal-mongo mongo:7
```

2) Configura variables de entorno y ejecuta la API:

```bash
export PORT_API=8080
export MONGO_URI='mongodb://localhost:27017'
export MONGO_DB='portal'
export CORS_ORIGINS='http://localhost:3000'
export JWT_SECRET='super-secret'
export JWT_ISSUER='djengua-api'
export JWT_TTL_MINUTES=60

go run ./cmd/api
```

3) Prueba salud y métricas:

```bash
curl -i http://localhost:8080/healthz
curl -s http://localhost:8080/metrics | head
```

## Configuración por variables de entorno

| Variable | Default | Descripción |
| --- | --- | --- |
| `PORT_API` | `8080` | Puerto HTTP de la API. |
| `MONGO_URI` | `mongodb://localhost:27017` | Cadena de conexión MongoDB. |
| `MONGO_DB` | `portal` | Base de datos MongoDB. |
| `CORS_ORIGINS` | `http://localhost:3000` | Orígenes permitidos en CORS. |
| `JWT_SECRET` | _requerido_ | Secreto para firmar JWT. |
| `JWT_ISSUER` | `djengua-api` | Emisor del JWT. |
| `JWT_TTL_MINUTES` | `60` | Tiempo de vida del token en minutos. |

## Autenticación

- **Header**: `Authorization: Bearer <token>`
- Los endpoints de **lectura** son públicos; los de **escritura** requieren JWT.

### Endpoints de auth

- `POST /auth/register`
  - Body: `{ "email", "password", "name" }`
  - Respuesta: `{ "token", "user" }`
- `POST /auth/login`
  - Body: `{ "email", "password" }`
  - Respuesta: `{ "token", "user" }`
- `GET /auth/me` *(JWT requerido)*
  - Respuesta: `User`

## Modelos y enums

> Los ejemplos representan el shape JSON que devuelve la API.

### Producto (`Product`)

- `status`: `active | inactive | draft`
- `type`: `physical | digital | service`
- `size`: `XS | S | M | L | XL`

```json
{
  "id": "uuid",
  "sku": "SKU-001",
  "name": "Producto",
  "description": "...",
  "price": 120.0,
  "cost": 60.0,
  "currency": "USD",
  "type": "physical",
  "tags": ["new", "gift"],
  "status": "active",
  "size": "M",
  "category_id": "uuid",
  "collection_ids": ["uuid"],
  "images": ["https://..."],
  "stock": 10,
  "attributes": { "color": "red" },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "deleted_at": null
}
```

### Categoría (`Category`)

```json
{
  "id": "uuid",
  "name": "Ropa",
  "slug": "ropa",
  "description": "...",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "deleted_at": null
}
```

### Colección (`Collection`)

```json
{
  "id": "uuid",
  "name": "Temporada",
  "slug": "temporada",
  "description": "...",
  "status": "active",
  "start_date": "2024-10-01T00:00:00Z",
  "end_date": "2024-12-31T00:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "deleted_at": null
}
```

### Orden (`Order`) y venta (`Sale`)

- `order.status`: `draft | pending | paid | cancelled`
- `sale.method`: `cash | card | transfer | other`

```json
{
  "id": "uuid",
  "number": "ORD-10001",
  "user_id": "uuid",
  "status": "pending",
  "currency": "USD",
  "items": [
    {
      "product_id": "uuid",
      "sku": "SKU-001",
      "name": "Producto",
      "qty": 1,
      "unit_price": 120.0,
      "line_total": 120.0
    }
  ],
  "subtotal": 120.0,
  "tax": 0.0,
  "total": 120.0,
  "notes": "...",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "deleted_at": null
}
```

```json
{
  "id": "uuid",
  "order_id": "uuid",
  "user_id": "uuid",
  "method": "card",
  "amount": 120.0,
  "currency": "USD",
  "reference": "STRIPE-123",
  "sold_at": "2024-01-01T00:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "deleted_at": null
}
```

### Usuario (`User`)

- `role`: `admin | user`

```json
{
  "id": "uuid",
  "email": "user@mail.com",
  "name": "User",
  "role": "user",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "deleted_at": null
}
```

## Endpoints

> Prefijo: `/api/v1`

### Productos

- `GET /products` (paginado + filtros)
  - `page`, `page_size`
  - `status`, `category_id`, `collection_id`, `min_price`, `max_price`, `q`
- `GET /products/{id}`
- `POST /products` *(JWT)*
- `PUT /products/{id}` *(JWT)*
- `PATCH /products/{id}` *(JWT)*
- `DELETE /products/{id}` *(JWT, soft-delete)*

**Body para crear/actualizar:**

```json
{
  "sku": "ROS-001",
  "name": "Ramo de rosas rojas",
  "description": "Ramo premium de 24 rosas rojas.",
  "price": 120.0,
  "cost": 80.0,
  "currency": "USD",
  "type": "physical",
  "tags": ["gift"],
  "status": "active",
  "size": "M",
  "category_id": "<UUID_CATEGORIA>",
  "collection_ids": ["<UUID_COLECCION>"],
  "images": ["https://..."],
  "stock": 10,
  "attributes": { "color": "red" }
}
```

### Categorías

- `GET /categories` (paginado)
  - `page`, `page_size`
- `GET /categories/{id}`
- `POST /categories` *(JWT)*
- `PUT /categories/{id}` *(JWT)*
- `PATCH /categories/{id}` *(JWT)*
- `DELETE /categories/{id}` *(JWT, soft-delete)*

**Body:**

```json
{
  "name": "Ropa",
  "slug": "ropa",
  "description": "...",
  "status": "active"
}
```

### Colecciones

- `GET /collections` (paginado)
  - `page`, `page_size`
- `GET /collections/{id}`
- `POST /collections` *(JWT)*
- `PUT /collections/{id}` *(JWT)*
- `PATCH /collections/{id}` *(JWT)*
- `DELETE /collections/{id}` *(JWT, soft-delete)*

**Body:**

```json
{
  "name": "Temporada",
  "slug": "temporada",
  "description": "...",
  "status": "active",
  "start_date": "2024-10-01T00:00:00Z",
  "end_date": "2024-12-31T00:00:00Z"
}
```

### Órdenes

- `POST /orders` *(JWT)*
- `GET /orders/me` *(JWT)*

**Body:**

```json
{
  "currency": "USD",
  "items": [
    { "product_id": "uuid", "sku": "SKU-001", "name": "Producto", "qty": 1, "unit_price": 120.0, "line_total": 120.0 }
  ],
  "notes": "..."
}
```

### Ventas

- `POST /sales` *(JWT)*

**Body:**

```json
{
  "order_id": "uuid",
  "method": "card",
  "amount": 120.0,
  "currency": "USD",
  "reference": "STRIPE-123",
  "sold_at": "2024-01-01T00:00:00Z"
}
```

## Reglas de validación relevantes

- `currency` debe tener **3 letras** (ej. `USD`).
- `status` de productos: `active | inactive | draft`.
- `status` de categorías/colecciones: `active | inactive`.
- `type` de producto: `physical | digital | service` (opcional).
- `size` de producto: `XS | S | M | L | XL` (opcional).
- `slug` debe ser URL-friendly: minúsculas, números y guiones.
- `start_date` debe ser **<=** `end_date` si ambos existen.

## Observabilidad (Prometheus)

La API expone métricas en `GET /metrics`:

- `http_requests_total{method,route,status}`
- `http_request_duration_seconds{method,route}`

Ejemplos de consultas en Prometheus:

```
sum(rate(http_requests_total[5m])) by (route)
```

```
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, route))
```

## Guía para construir el front (React + Next.js + Zustand + Zod)

> Objetivo: recrear la UX de `index.html` con datos reales de esta API.

### Stack recomendado

- **Next.js (App Router)**: rutas `/`, `/product/[id]`, `/checkout`.
- **Zustand**: estado global (carrito, filtros, auth token, UI).
- **Zod**: validación y parseo de respuestas del API en el client.
- **Fetch + Server Actions o Route Handlers** para llamadas a la API.

### Mapeo UI (index.html → API)

- **Hero + categorías**: consumir `GET /categories`.
- **Chips/filtros**: usar `GET /products` con `category_id`, `min_price`, `max_price`, `q`.
- **Productos**: listados con `GET /products`; detalle con `GET /products/{id}`.
- **Carrito**: estado local en Zustand; al confirmar, crear `POST /orders` y luego `POST /sales`.
- **Login/registro**: `POST /auth/register` y `POST /auth/login`.

### Estructura sugerida

```
app/
  page.tsx               // home (hero, filtros, listado)
  product/[id]/page.tsx  // detalle
  checkout/page.tsx
lib/
  api.ts                 // client fetch + token
  schemas.ts             // zod schemas
store/
  cart.ts                // zustand store
  filters.ts
  auth.ts
components/
  ProductCard.tsx
  Filters.tsx
  CartDrawer.tsx
```

### Zod schemas (ejemplo conceptual)

- `ProductSchema`, `CategorySchema`, `CollectionSchema`, `OrderSchema`, `SaleSchema`.
- `ApiList<T>` si quieres normalizar paginación.

### Estado en Zustand (ejemplo conceptual)

- `cart.items`, `cart.total`, `cart.addItem/removeItem`.
- `filters.query`, `filters.categoryId`, `filters.maxPrice`.
- `auth.token`, `auth.user` (persistencia en `localStorage`).

### Consideraciones

- Guardar el token en memoria + `localStorage` y mandarlo como `Authorization: Bearer`.
- Validar respuestas con Zod antes de renderizar.
- Debounce para búsquedas y filtros.

## Notas

- El proyecto usa **soft-delete** (`deleted_at`) y no elimina físicamente los registros.
- Los campos de fecha se manejan en UTC.
- En producción, agrega autenticación/autorización más robusta y rate limiting.
