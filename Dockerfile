# Build stage
FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o rate-my .

# Final stage
FROM scratch AS runtime
WORKDIR /app
COPY --from=builder /app/rate-my ./rate-my
COPY static ./static
EXPOSE 8080
ENV PORT=8080
CMD ["./rate-my"]
