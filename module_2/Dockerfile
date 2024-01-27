FROM golang:1.17 as builder

WORKDIR /app

COPY * ./

RUN CGO_ENABLED=0 GOOS=linux go build -v -mod=readonly -o http_server

FROM alpine:3
COPY --from=builder /app/http_server /http_server
EXPOSE 8080
CMD ["/http_server"]