# Usa una imagen base de Golang
FROM golang:1.23

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia solo el main
COPY . .

# Descarga las dependencias del módulo Go
RUN go mod download
RUN go mod tidy
# Compila la aplicación
RUN go build -o master .

# Expone el puerto en el que la aplicación escuchará
EXPOSE 8081

# Comando para ejecutar la aplicación
CMD ["./master"]