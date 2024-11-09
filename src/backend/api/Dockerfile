# Usa una imagen base de Golang
FROM golang:1.23

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos del proyecto al directorio de trabajo
COPY . .

# Descarga las dependencias del módulo Go
# RUN go mod download
#Update dependencies
RUN go get -u ./...

# Compila la aplicación
RUN go mod tidy
RUN go build -o api .

# Comando para ejecutar la aplicación
CMD ["./api"]