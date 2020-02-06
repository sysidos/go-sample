# Builder
FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o invoice main.go


# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app

WORKDIR /app
EXPOSE 8080
COPY --from=builder /app/invoice /app

ENTRYPOINT ["./invoice"]
