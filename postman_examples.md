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

### Eliminar Servicio (DELETE definitivo)
```
DELETE http://localhost:8080/api/services/1
Authorization: Bearer {token}
```

**Respuesta esperada (200 OK):**
```json
{
    "message": "Service deleted successfully"
}
```

**Errores posibles:**
- 401 Unauthorized: Token inválido o no proporcionado
- 403 Forbidden: El servicio no pertenece al usuario
- 404 Not Found: Servicio no encontrado
- 500 Internal Server Error: Error del servidor

**Nota importante:** Este endpoint elimina el servicio de manera DEFINITIVA de la base de datos. Para cambiar el estado del servicio a inactivo, usar el endpoint PATCH `/services/{id}/status`.

### Actualizar Perfil de Usuario
```
PUT http://localhost:8080/api/user/profile
Authorization: Bearer {token}
Content-Type: application/json

{
    "name": "Juan Carlos Pérez",
    "email": "juan.carlos@example.com",
    "locality": "La Plata", 
    "province": "Buenos Aires",
    "phone": "+54 221 1234-5678"
}
```

**Respuesta esperada (200 OK):**
```json
{
    "message": "User updated successfully",
    "data": {
        "id": 1,
        "name": "Juan Carlos Pérez",
        "email": "juan.carlos@example.com",
        "locality": "La Plata",
        "province": "Buenos Aires", 
        "phone": "+54 221 1234-5678",
        "created_at": "2025-01-01T10:00:00Z",
        "updated_at": "2025-01-01T15:30:00Z"
    }
}
```

**Errores posibles:**
- 400 Bad Request: Datos inválidos o email ya en uso por otro usuario
- 401 Unauthorized: Token inválido o no proporcionado
- 404 Not Found: Usuario no encontrado
- 409 Conflict: Email ya está en uso
- 500 Internal Server Error: Error del servidor

**Nota:** Todos los campos son opcionales. Solo se actualizarán los campos enviados en la petición.

### Obtener Disponibilidad de Servicio
```
GET http://localhost:8080/api/services/1/availability?date=2025-10-15
```

**Respuesta esperada (200 OK):**
```json
{
    "message": "Service availability retrieved successfully",
    "data": {
        "date": "2025-10-15",
        "day_of_week": "tuesday",
        "time_slots": [
            {"time": "08:00-08:30", "available": true},
            {"time": "08:30-09:00", "available": false},
            {"time": "09:00-09:30", "available": true}
        ]
    }
}
```

### Obtener Calendario de Servicio (30 días)
```
GET http://localhost:8080/api/services/1/calendar
```

**Respuesta esperada (200 OK):**
```json
{
    "message": "Service calendar retrieved successfully",
    "data": {
        "service_id": 1,
        "service_title": "Clases de Programación",
        "start_date": "2025-10-02",
        "end_date": "2025-10-31",
        "days": [
            {
                "date": "2025-10-02",
                "day_of_week": "miércoles",
                "has_availability": true,
                "available_slots": 8,
                "total_slots": 8
            },
            {
                "date": "2025-10-03",
                "day_of_week": "jueves",
                "has_availability": true,
                "available_slots": 6,
                "total_slots": 8
            },
            {
                "date": "2025-10-04",
                "day_of_week": "viernes",
                "has_availability": false,
                "available_slots": 0,
                "total_slots": 0
            }
        ]
    }
}
```

**Descripción del endpoint:**
- Retorna los próximos 30 días desde la fecha actual
- `has_availability`: indica si el día tiene horarios configurados
- `available_slots`: número de slots libres (sin appointments)
- `total_slots`: número total de slots configurados para ese día
- `day_of_week`: día de la semana en español
- Útil para mostrar calendarios en el frontend con vista general de disponibilidad

### Crear Cita/Appointment
```
POST http://localhost:8080/api/appointments
Authorization: Bearer {token}
Content-Type: application/json

{
    "service_id": 1,
    "date": "2025-10-15",
    "time_slot": "09:00-09:30",
    "notes": "Necesito ayuda con programación en Go"
}
```

**Respuesta esperada (201 Created):**
```json
{
    "message": "Appointment created successfully",
    "data": {
        "id": 1,
        "service": {
            "id": 1,
            "title": "Clases de Programación",
            "description": "Clases personalizadas...",
            "category": "educacion",
            "price": 2500.0
        },
        "client_id": 2,
        "provider_id": 1,
        "date": "2025-10-15",
        "time_slot": "09:00-09:30",
        "status": "pending",
        "notes": "Necesito ayuda con programación en Go"
    }
}
```

### Listar Mis Citas (Cliente)
```
GET http://localhost:8080/api/my-appointments
Authorization: Bearer {token}
```

### Ver Citas de Mi Servicio (Proveedor)
```
GET http://localhost:8080/api/services/1/appointments
Authorization: Bearer {token}
```

### Aceptar/Rechazar Cita (Proveedor)
```
PUT http://localhost:8080/api/appointments/1/status
Authorization: Bearer {token}
Content-Type: application/json

{
    "status": "accepted"
}
```

**Estados válidos:**
- `accepted` - Solo proveedor puede aceptar citas pendientes
- `rejected` - Solo proveedor puede rechazar citas pendientes
- `cancelled` - Cliente o proveedor pueden cancelar
- `completed` - Solo proveedor puede marcar como completada

**Errores comunes:**
- 400 Bad Request: Datos inválidos o estado no permitido
- 401 Unauthorized: Token inválido
- 403 Forbidden: Sin permisos para la acción
- 404 Not Found: Cita o servicio no encontrado
- 409 Conflict: Horario ya ocupado
- 500 Internal Server Error: Error del servidor

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
        },
        {
            "name": "Delete Service",
            "request": {
                "method": "DELETE",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {token}"
                    }
                ],
                "url": {
                    "raw": "http://localhost:8080/api/services/1",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "services", "1"]
                }
            }
        },
        {
            "name": "Update User Profile",
            "request": {
                "method": "PUT",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {token}"
                    },
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": "{\n    \"name\": \"Juan Carlos Pérez\",\n    \"email\": \"juan.carlos@example.com\",\n    \"locality\": \"La Plata\",\n    \"province\": \"Buenos Aires\",\n    \"phone\": \"+54 221 1234-5678\"\n}"
                },
                "url": {
                    "raw": "http://localhost:8080/api/user/profile",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "user", "profile"]
                }
            }
        },
        {
            "name": "Get Service Availability",
            "request": {
                "method": "GET",
                "header": [],
                "url": {
                    "raw": "http://localhost:8080/api/services/1/availability?date=2025-10-15",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "services", "1", "availability"],
                    "query": [
                        {
                            "key": "date",
                            "value": "2025-10-15"
                        }
                    ]
                }
            }
        },
        {
            "name": "Get Service Calendar",
            "request": {
                "method": "GET",
                "header": [],
                "url": {
                    "raw": "http://localhost:8080/api/services/1/calendar",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "services", "1", "calendar"]
                }
            }
        },
        {
            "name": "Create Appointment",
            "request": {
                "method": "POST",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {token}"
                    },
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": "{\n    \"service_id\": 1,\n    \"date\": \"2025-10-15\",\n    \"time_slot\": \"09:00-09:30\",\n    \"notes\": \"Necesito ayuda con programación en Go\"\n}"
                },
                "url": {
                    "raw": "http://localhost:8080/api/appointments",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "appointments"]
                }
            }
        },
        {
            "name": "List My Appointments",
            "request": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {token}"
                    }
                ],
                "url": {
                    "raw": "http://localhost:8080/api/my-appointments",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "my-appointments"]
                }
            }
        },
        {
            "name": "View Service Appointments",
            "request": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {token}"
                    }
                ],
                "url": {
                    "raw": "http://localhost:8080/api/services/1/appointments",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "services", "1", "appointments"]
                }
            }
        },
        {
            "name": "Update Appointment Status",
            "request": {
                "method": "PUT",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "Bearer {token}"
                    },
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": "{\n    \"status\": \"accepted\"\n}"
                },
                "url": {
                    "raw": "http://localhost:8080/api/appointments/1/status",
                    "protocol": "http",
                    "host": ["localhost"],
                    "port": "8080",
                    "path": ["api", "appointments", "1", "status"]
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
