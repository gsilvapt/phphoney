FROM golang:1.27.1 AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download
COPY *.go ./
COPY public ./public
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/srv

FROM scratch
COPY --from=builder /app/srv /app

LABEL description="Capture all the tings" 
EXPOSE 8080
ENTRYPOINT ["/app"]
