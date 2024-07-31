FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN god mod download

COPY . .

RUN go build -o index .

CMD ["./index"]