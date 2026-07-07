FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o profile-worker ./cmd/profile-worker

FROM alpine:3.19

RUN apk --no-cache add ca-certificates && \
    addgroup -S worker && adduser -S -G worker worker

WORKDIR /app

COPY --from=builder --chown=worker:worker /app/profile-worker .

USER worker

EXPOSE 8080

CMD ["./profile-worker"]
