# ── Build stage ──────────────────────────────────────────────────────────────
FROM --platform=$BUILDPLATFORM golang:1.27.0-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /build/rate-my .

# ── Runtime stage (distroless) ────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

COPY --from=builder /build/rate-my /app/rate-my
COPY --chown=nonroot:nonroot ./static ./static

USER nonroot:nonroot

EXPOSE 8080
ENV PORT=8080

ENTRYPOINT ["/app/rate-my"]
