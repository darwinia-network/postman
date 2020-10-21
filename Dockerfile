# First layer as builder
FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN apk add --no-cache make \
    && make

# Fetch the binary
FROM alpine:latest
COPY --from=builder /src/target/postman /usr/local/bin
CMD ["postman"]
