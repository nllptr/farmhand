FROM golang:latest AS builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# GCO_ENABLED=0 is required to produce a statically linked
# binary file. Without static linking, it will probably not
# run correctly.
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/auth

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /build/auth .
CMD ["./auth"]
