# Usa la imagen oficial de Go para compilar el binario
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Copia los archivos de dependencias y descarga los paquetes
COPY go.mod go.sum ./
RUN go mod download

# Copia el código fuente y compila el binario
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Usa una imagen mínima de alpine para producción
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

# Expone el puerto y ejecuta el binario
EXPOSE 8080
CMD ["./main"]