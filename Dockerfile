FROM golang:1.15-alpine AS builder

WORKDIR $GOPATH/src/github.com/jfreeland/ec2-network-monitor

# Download module dependencies to take advantage of Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir output && \
      CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o output/ ./...

FROM golang:1.15-alpine
COPY --from=builder /go/src/github.com/jfreeland/ec2-network-monitor/output/ec2nm /usr/local/bin
CMD ["ec2nm"]
