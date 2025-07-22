# IYCDS2025 API

API simple desarrollada en Go para el proyecto IYCDS2025.

## Descripción

API minimalista que proporciona un endpoint básico `/ping` que responde `pong` con conexión a MySQL local.

## Arquitectura

```
iycds2025_api/
├── configs/
│   └── database.go                    # Configuración de base de datos
├── src/
│   └── api/
│       ├── app/
│       │   ├── application.go         # Aplicación principal
│       │   └── url_mappings.go        # Configuración de rutas
│       ├── core/
│       │   └── usecases/
│       │       └── ping/
│       │           └── ping.go        # Caso de uso ping
│       ├── infrastructure/
│       │   ├── dependencies/
│       │   │   └── dependencies.go    # Inyección de dependencias
│       │   └── entrypoints/
│       │       └── api/
│       │           ├── handler.go     # Interface Handler
│       │           └── handlers/
│       │               └── ping.go    # Handler ping
│       ├── middleware/
│       │   └── cors.go                # Middleware CORS
│       ├── repositories/
│       │   └── database/
│       │       └── repository.go      # Repository base
│       └── main.go                    # Punto de entrada
├── go.mod                             # Dependencias del proyecto
├── go.sum                             # Checksums de dependencias
├── docker-compose.yml                 # Docker Compose para desarrollo
├── Dockerfile                         # Imagen Docker
├── postman_examples.md                # Ejemplos para Postman
└── README.md                          # Este archivo
```

## Tecnologías

- **Go 1.23+**: Lenguaje de programación
- **Gin**: Framework web para Go
- **MySQL**: Base de datos relacional
- **Clean Architecture**: Separación de capas y responsabilidades

## Requisitos Previos

### Software Necesario

1. **Go 1.23 o superior**
   ```bash
   go version
   ```

2. **MySQL 8.0 o superior**
   ```bash
   mysql --version
   ```

3. **Git** (para clonar el repositorio)

### Base de Datos

1. **Crear base de datos MySQL local:**
   ```sql
   CREATE DATABASE iycds2025 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

2. **Crear usuario (opcional):**
   ```sql
   CREATE USER 'iycds_user'@'localhost' IDENTIFIED BY 'iycds_password';
   GRANT ALL PRIVILEGES ON iycds2025.* TO 'iycds_user'@'localhost';
   FLUSH PRIVILEGES;
   ```

## Configuración del Ambiente de Trabajo

### 1. Clonar el repositorio

```bash
git clone [URL_DEL_REPOSITORIO]
cd iycds2025_api
```

### 2. Configurar variables de entorno

Configurar las siguientes variables de entorno:

```bash
# Entorno de ejecución
export APP_ENV=development

# Configuración de Base de Datos Local
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=root
export DB_NAME=iycds2025

# Puerto del servidor
export PORT=8080
```

### 3. Instalar dependencias

```bash
go mod download
go mod tidy
```

### 4. Ejecutar la aplicación

```bash
# Desde la raíz del proyecto
go run src/api/main.go
```

O compilar y ejecutar:

```bash
go build -o iycds2025_api src/api/main.go
./iycds2025_api
```

## Endpoint Disponible

### Ping
- **URL**: `GET /ping`
- **Descripción**: Endpoint básico que responde pong
- **Respuesta**:
  ```json
  {
    "message": "pong"
  }
  ```

## Testing

### Probar con curl

```bash
curl http://localhost:8080/ping
```

### Probar con Postman

Ver el archivo `postman_examples.md` para ejemplos detallados.

## Desarrollo con Docker

### Docker Compose para desarrollo

```bash
# Levantar todos los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down
```

## Estructura de Código

### Clean Architecture

El proyecto sigue los principios de Clean Architecture:

- **`configs/`**: Configuraciones de la aplicación
- **`src/api/core/`**: Lógica de negocio (casos de uso)
- **`src/api/infrastructure/`**: Capas de infraestructura (handlers, dependencies)
- **`src/api/repositories/`**: Acceso a datos
- **`src/api/middleware/`**: Middleware de la aplicación

### Flujo de Ejecución

1. **main.go** → Punto de entrada
2. **application.go** → Configuración del servidor
3. **dependencies.go** → Inyección de dependencias
4. **url_mappings.go** → Configuración de rutas
5. **handlers/ping.go** → Manejo de la petición
6. **usecases/ping.go** → Lógica de negocio

## Variables de Entorno

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `APP_ENV` | Entorno de ejecución | `development` |
| `PORT` | Puerto del servidor | `8080` |
| `DB_HOST` | Host de MySQL | `localhost` |
| `DB_PORT` | Puerto de MySQL | `3306` |
| `DB_USER` | Usuario de MySQL | `root` |
| `DB_PASSWORD` | Contraseña de MySQL | `root` |
| `DB_NAME` | Nombre de la base de datos | `iycds2025` |

## Troubleshooting

### Error de conexión a la base de datos
1. Verificar que MySQL esté ejecutándose
2. Comprobar credenciales en variables de entorno
3. Verificar que la base de datos `iycds2025` exista

### Error de dependencias
```bash
go mod download
go mod tidy
```

### Puerto ocupado
Cambiar el puerto con variable de entorno:
```bash
PORT=8081 go run src/api/main.go
```

## Próximos Pasos

- [ ] Agregar más endpoints según necesidades del proyecto
- [ ] Implementar autenticación (si es necesario)
- [ ] Agregar logging estructurado
- [ ] Implementar tests unitarios
- [ ] Configurar CI/CD

## Contribución

1. Fork del proyecto
2. Crear rama feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit de cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear Pull Request
