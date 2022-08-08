FROM golang:1.19.0-bullseye AS builder

RUN mkdir /app
COPY ./local_endpoint.go /app
WORKDIR /app
RUN go mod init github.com/Rashoru-Infinity/local_endpoint && \
    go mod tidy && \
    go build local_endpoint.go && \
    chmod 755 local_endpoint

FROM gcr.io/distroless/base

COPY --from=builder /app/local_endpoint /app/local_endpoint
ENTRYPOINT ["/app/local_endpoint"]