# First layer as builder
FROM golang:alpine as builder
COPY . postman
RUN apk add --no-cache make\
    && cd postman \
    && make

# Fetch the binary
FROM alpine:latest
COPY --from=builder /go/postman/target/postman /usr/local/bin
ENTRYPOINT ["postman"]
