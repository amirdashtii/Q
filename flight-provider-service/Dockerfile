# Dockerfile
FROM golang:alpine3.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /flight-provider-service ./cmd/main.go

EXPOSE 8080

CMD [ "/flight-provider-service" ]
