# Используем официальный образ Golang для сборки приложения
FROM golang:1.18 AS builder

# Установка переменных среды
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Создание директории приложения внутри образа
WORKDIR /app

# Копирование файла go.mod и go.sum и загрузка зависимостей
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копирование исходного кода приложения в образ
COPY ./ ./

# Сборка приложения
RUN go build -o product-service cmd/apiserver/main.go

# Используем минимальный образ для выполнения
FROM alpine:latest

# Установка зависимостей, необходимых для выполнения приложения
RUN apk --no-cache add ca-certificates

# Копирование исполняемого файла из предыдущего образа
COPY --from=builder /app/product-service /app/product-service

# Установка переменных среды для подключения к PostgreSQL
ENV DB_HOST=$DB_HOST \
    DB_PORT=$DB_PORT \
    DB_USER=$DB_USER \
    DB_PASSWORD=$DB_PASSWORD \
    DB_NAME=$DB_NAME \
    LOG_LEVEL=$LOG_LEVEL

# Указание порта, на котором будет работать приложение
EXPOSE 8080

# Запуск приложения
CMD ["/app/product-service"]
