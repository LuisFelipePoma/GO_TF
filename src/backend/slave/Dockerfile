# Usa una imagen base de Golang
FROM golang:1.23-alpine

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos del proyecto al directorio de trabajo
COPY . .

# Descarga las dependencias del módulo Go
RUN go mod download

# Upgrade dependencies
RUN go get -u ./...
RUN go mod tidy


# Compila la aplicación
RUN go build -o slave .

# Comando para ejecutar la aplicación
CMD ["./slave"]