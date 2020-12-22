FROM golang:1.15-alpine AS builder

WORKDIR /src

# Download module dependencies to take advantage of Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags="-w -s" -o /go/bin/ec2nm cmd/ec2nm/main.go

FROM scratch
# SSL Certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=builder /go/bin/ec2nm /usr/local/bin/ec2nm

# Run the binary
ENTRYPOINT ["/usr/local/bin/ec2nm"]
