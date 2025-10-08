FROM golang:1.24.6-alpine
WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../ .
RUN go build -o sensor ./cmd/sensor/main.go
CMD ["./sensor"]