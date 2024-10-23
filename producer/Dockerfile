 # Go builder
FROM golang:1.23 AS builder

WORKDIR /usr/src/webapp

# pre-downloade dependencies (only redownloading if they change)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/webapp

CMD ["webapp"]

# Production image
# FROM alpine:latest

# WORKDIR /root
# COPY --from=builder /usr/local/bin/webapp .

# RUN pwd && whoami && ls -l /root

# CMD ["/root/webapp"]
