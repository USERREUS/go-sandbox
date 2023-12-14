FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -v ./cmd/apiserver

CMD ["./apiserver"]
