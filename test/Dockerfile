FROM golang:latest

WORKDIR /app

COPY . .

RUN go get -u github.com/lib/pq

CMD ["go", "run", "main.go"]
