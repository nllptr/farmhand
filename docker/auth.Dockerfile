FROM golang:latest AS builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/auth

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /build/auth .
CMD ["./auth"]  