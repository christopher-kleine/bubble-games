# Specifies a parent image
FROM golang:1.20-alpine AS builder
 
# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copies everything from your root directory into /app
COPY . .

# Installs Git
RUN apk add --no-cache git
 
# Installs Go dependencies
RUN go mod download
 
# Builds your app with optional configuration
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" .

# -----------------------------

FROM scratch

# Copies the binary from the build stage into the root directory
COPY --from=builder /app/bubble-games /bubble-games
 
# Tells Docker which network port your container listens on
EXPOSE 2222
 
# Specifies the executable command that runs when the container starts
ENTRYPOINT ["/bubble-games"]
