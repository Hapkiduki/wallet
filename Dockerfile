FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

# Build a statically linked binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/api

FROM scratch

COPY --from=builder /app/main /app/main

EXPOSE 8080

ENTRYPOINT ["/app/main"]