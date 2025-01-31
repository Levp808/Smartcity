FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o orchestrator_service ./cmd/main

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/orchestrator_service /orchestrator_service

COPY --from=builder /app/configs/scheduler.env /configs/scheduler.env

COPY --from=builder /app/ai_assist/prompt.json /ai_assist/prompt.json

EXPOSE 50051

CMD ["./orchestrator_service"]