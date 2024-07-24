# Build stage
FROM golang:1.22.4-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN go mod download
RUN go build -o main .

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY .env ./
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/assests ./assests
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/controllers ./controllers
COPY --from=builder /app/models ./models
COPY --from=builder /app/middleware ./middleware
COPY --from=builder /app/utils ./utils
COPY --from=builder /app/initializers ./initializers
COPY --from=builder /app/wait-for-it.sh ./
# Install necessary tools
RUN apk add --no-cache netcat-openbsd

CMD ["./wait-for-it.sh", "database", "5432", "--", "./main"]