# Используем официальный образ Golang для сборки приложения
FROM golang:1.18 as builder

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
COPY . .

# Сборка приложения
RUN go build -o order-service cmd/apiserver/main.go

# Используем минимальный образ для выполнения
FROM alpine:latest

# Установка зависимостей, необходимых для выполнения приложения
RUN apk --no-cache add ca-certificates

# Копирование исполняемого файла из предыдущего образа
COPY --from=builder /app/order-service /app/order-service

# Запуск приложения
CMD ["/app/order-service"]
