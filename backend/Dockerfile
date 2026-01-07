# build project
FROM golang:latest AS builder

WORKDIR /app

# Only copy go.mod & go.sum first (improves caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o backend_server .


# copy binary to runtime container
# FROM alpine:3.20
# this is even smaller than aline
FROM gcr.io/distroless/static-debian12:latest
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/backend_server .

# Ensure binary can run (not always needed, but safe)
# doesn't work on distroless, so comment out for now.  uncomment if go back to alpine
# RUN chmod +x ./backend_server

EXPOSE 8080
CMD ["./backend_server"]
