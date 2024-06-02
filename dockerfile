# Usar una imagen base de Golang
FROM golang:1.16

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el módulo Go y el archivo sum
COPY go.mod go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar el resto de los archivos
COPY . .

# Construir la aplicación
RUN go build -o main .

# Definir el comando por defecto para ejecutar la aplicación
CMD ["./main"]
