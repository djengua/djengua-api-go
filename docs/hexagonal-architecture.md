# Análisis de rediseño hacia arquitectura hexagonal

## Estado actual (post-migración)

- Los handlers HTTP dependen de servicios de aplicación en lugar de acceder a MongoDB directamente, reduciendo el acoplamiento con la infraestructura.
- La infraestructura de persistencia vive en `internal/adapters/mongo` y expone implementaciones de puertos definidos en `internal/ports`.
- El dominio (`internal/domain`) mantiene las entidades, mientras que los puertos y servicios permiten separar casos de uso y adaptadores.
- El ensamblaje de dependencias en `cmd/api/main.go` construye repositorios, servicios y los inyecta en el router HTTP.

## Objetivo de la arquitectura hexagonal

Separar la lógica del negocio (dominio y casos de uso) de los detalles de infraestructura (HTTP, MongoDB), definiendo puertos que permitan intercambiar adaptadores sin afectar el núcleo.

## Cambios propuestos

### 1) Crear una capa de aplicación (casos de uso)

- **Nuevo paquete** `internal/application` (o `internal/core/usecase`) con servicios orientados a casos de uso:
  - `ProductService`, `CategoryService`, `CollectionService`.
- Estos servicios deben orquestar:
  - Validaciones de negocio (antes dispersas en handlers).
  - Normalización de datos (por ejemplo, currency upper-case, trims).
  - Coordinación entre puertos (repositorios).

### 2) Definir puertos (interfaces) en el núcleo

- **Nuevo paquete** `internal/ports` (o `internal/core/ports`) con interfaces de repositorio:
  - `ProductRepository`, `CategoryRepository`, `CollectionRepository`.
  - Definir filtros (`ProductFilters`) y errores (`ErrNotFound`) en este paquete para evitar dependencia inversa con infraestructura.
- El dominio puede permanecer en `internal/domain` o moverse a `internal/core/domain` si se quiere agrupar todo el núcleo.

### 3) Convertir `internal/db` en adaptador de salida

- Mover el código de MongoDB a `internal/adapters/mongo`.
- Implementar las interfaces de repositorio definidas en `internal/ports`.
- Mantener los documentos BSON y conversiones en este adaptador.

### 4) Convertir `internal/httpapi` en adaptador de entrada

- Los handlers deben depender de servicios de aplicación, no de Mongo.
- Agregar DTOs (request/response) en el adaptador HTTP si se quiere separar el modelo de dominio de los contratos externos.
- Validaciones HTTP (schema/basic) pueden quedarse en handlers, pero las reglas de negocio deben estar en aplicación o dominio.

### 5) Reorganizar el ensamblado de dependencias

- En `cmd/api/main.go`:
  - Construir adaptador Mongo (repositorios).
  - Construir servicios de aplicación y pasarlos a handlers.
  - Router igual, pero apuntando a la capa de aplicación.

### 6) Tests

- **Unit tests** en aplicación usando mocks de puertos (sin DB ni HTTP).
- **Integration tests** en adaptadores (Mongo y HTTP) según corresponda.

## Estructura aplicada

```
internal/
  application/
  domain/
  ports/
  adapters/
    mongo/
  httpapi/
    handlers/
    middleware/
    router/
  config/
cmd/api/main.go
```

> Se mantuvo `internal/domain` y `internal/httpapi` para una migración incremental, y se añadieron `internal/ports`, `internal/application` y `internal/adapters/mongo`.

## Migración incremental (aplicada)

1) Se crearon interfaces en `internal/ports`.
2) Se introdujeron servicios en `internal/application` y se refactorizaron handlers para depender de ellos.
3) Se movió `internal/db` a `internal/adapters/mongo` y se ajustaron imports.
4) Queda pendiente revisar tests y añadir mocks.

## Impacto esperado

- **Menor acoplamiento** entre API HTTP y MongoDB.
- **Mayor testabilidad** de lógica de negocio.
- **Facilidad para agregar adaptadores** (p. ej. gRPC, colas, otra BD).

## Riesgos y consideraciones

- Cambios de paquetes implican ajustes de imports y posibles ciclos; es clave mantener la dirección de dependencias hacia el núcleo.
- Requiere disciplina en dónde se colocan validaciones y DTOs para evitar duplicidad.
