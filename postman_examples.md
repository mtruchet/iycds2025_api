# Ejemplos Postman para IYCDS2025 API

Este archivo contiene ejemplos de requests para probar la API con Postman.

## Requests de ejemplo

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

### Login de Usuario
```
POST http://localhost:8080/api/user/login
Content-Type: application/json

{
    "email": "test@example.com",
    "password": "testpassword123"
}
```

**Respuesta esperada:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Registro de Usuario
```
POST http://localhost:8080/api/user/register
Content-Type: application/json

{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "password": "password123",
    "password_confirm": "password123",
    "locality": "Buenos Aires",
    "province": "Buenos Aires",
    "phone": "+54 11 1234-5678"
}
```

**Respuesta esperada:**
```json
{
    "user": {
        "id": 1,
        "name": "Juan Pérez",
        "email": "juan@example.com",
        "locality": "Buenos Aires",
        "province": "Buenos Aires",
        "phone": "+54 11 1234-5678",
        "first_login": true,
        "created_at": "2025-01-13T10:30:00Z",
        "updated_at": "2025-01-13T10:30:00Z"
    }
}
```

### Recuperación de Contraseña
```
POST http://localhost:8080/api/user/forgot-password
Content-Type: application/json

{
    "email": "test@example.com"
}
```

**Respuesta esperada:**
```json
{
    "message": "If your email is registered in our system, you will receive instructions to reset your password"
}
```

### Restablecer Contraseña
```
POST http://localhost:8080/api/user/reset-password
Content-Type: application/json

{
    "token": "abc123def456...",
    "new_password": "newpassword123"
}
```

**Respuesta esperada:**
```json
{
    "message": "Password reset successfully"
}
```

## Importar en Postman

1. Abrir Postman
2. Click en "Import"
3. Crear una nueva colección llamada "IYCDS2025 API"
4. Agregar los requests usando la siguiente colección:

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
        },
        {
            "name": "Login",
            "request": {
                "method": "POST",
                "header": [
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": "{\n    \"email\": \"test@example.com\",\n    \"password\": \"testpassword123\"\n}"
                },
                "url": {
                    "raw": "http://localhost:8080/api/auth/login",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "auth", "login"]
                }
            }
        },
        {
            "name": "Forgot Password",
            "request": {
                "method": "POST",
                "header": [
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": "{\n    \"email\": \"test@example.com\"\n}"
                },
                "url": {
                    "raw": "http://localhost:8080/api/auth/forgot-password",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "auth", "forgot-password"]
                }
            }
        },
        {
            "name": "Reset Password",
            "request": {
                "method": "POST",
                "header": [
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": "{\n    \"token\": \"your_reset_token_here\",\n    \"new_password\": \"newpassword123\"\n}"
                },
                "url": {
                    "raw": "http://localhost:8080/api/auth/reset-password",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "auth", "reset-password"]
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
