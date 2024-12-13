FROM golang:1.22.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN god mod download

COPY . .

RUN go build -o main ./cmd/server/main.go

RUN chmod +x main

EXPOSE 3200

CMD ["./main"]
