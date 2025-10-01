FROM golang:1.25 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first to enable caching dependencies
COPY go.mod go.sum ./
COPY .env .

# Download dependencies.
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application.
# -o /usr/local/bin/server: Names the output binary 'server' and places it in a common directory.
# CGO_ENABLED=0: Disables CGO for a truly static binary, which is ideal for minimal base images like 'alpine'.
# -ldflags -s -w: Strips debug symbols and dwarf table for a smaller final binary size.
RUN CGO_ENABLED=0 go build -a -tags netgo \
  -ldflags '-w -s' \
  -o /bin/cli \
  ./cmd

# RUN chmod +x /bin/cli
# RUN chmod +x /app
# Use the minimal Alpine image for the final container, making it much smaller and more secure.
FROM alpine:latest

# Expose the port that the net/http server will listen on (8080 is a common default)
EXPOSE 8080

# Set a non-root user for security best practices
RUN adduser -D -g '' appnonrootuser

# Copy the compiled binary from the build stage
COPY --from=build --chown=appnonrootuser:appnonrootuser /app/bin /app

# # Explicitly set ownership of the binary to the non-root user (appnonrootuser).
# RUN chown appnonrootuser:appnonrootuser /usr/local/bin/run

# Switch to the non-root user for execution.
USER appnonrootuser

# Set the entry point to run the compiled application.
# The Cobra-based app should be designed to handle the 'serve' command itself.
ENTRYPOINT ["/app/bin/"]

CMD ["run"]
