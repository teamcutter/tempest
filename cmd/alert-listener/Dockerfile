FROM golang:1.24.6-alpine
WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../ .
RUN go build -o alert-listener ./cmd/alert-listener/main.go
CMD ["./alert-listener"]