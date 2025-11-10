FROM golang:1.25.3-alpine AS builder
# Set the working directory inside the container
WORKDIR /app

# Copy the module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application's source code
COPY . .

# Build the application
#
# ---- THIS IS THE CORRECTED LINE ----
# It now builds the package located in ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build \
  -a -installsuffix cgo \
  -o /server \
  ./cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /server /server
USER nonroot
EXPOSE 8080
ENTRYPOINT ["/server"]