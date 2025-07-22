# Dockerfile para IYCDS2025 API

# Etapa de construcción
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iycds2025_api src/api/main.go

# Etapa de ejecución
FROM alpine:latest

# Instalar certificados CA para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/iycds2025_api .

# Exponer puerto
EXPOSE 8080

# Variables de entorno por defecto
ENV APP_ENV=production
ENV PORT=8080
ENV GIN_MODE=release

# Ejecutar la aplicación
CMD ["./iycds2025_api"]
