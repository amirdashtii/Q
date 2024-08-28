# Dockerfile
FROM golang:alpine3.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /flight-ticket-service ./cmd/main.go

EXPOSE 8081

CMD [ "/flight-ticket-service" ]
