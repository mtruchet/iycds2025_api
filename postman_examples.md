# Ejemplos Postman para IYCDS2025 API

Este archivo contiene ejemplos de requests para probar la API con Postman.

## Request de ejemplo

### Ping
```
GET http://localhost:8080/ping
```

**Respuesta esperada:**
```json
{
    "message": "pong"
}
```

## Importar en Postman

1. Abrir Postman
2. Click en "Import"
3. Crear una nueva colección llamada "IYCDS2025 API"
4. Agregar el request:

```json
{
    "info": {
        "name": "IYCDS2025 API",
        "description": "Colección de endpoints para IYCDS2025 API",
        "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
    },
    "item": [
        {
            "name": "Ping",
            "request": {
                "method": "GET",
                "header": [],
                "url": {
                    "raw": "http://localhost:8080/ping",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["ping"]
                }
            }
        }
    ]
}
```

## Variables de entorno en Postman

Configurar estas variables en Postman para facilitar el testing:

- `base_url`: `http://localhost:8080`

Luego usar en las URLs: `{{base_url}}/ping`
