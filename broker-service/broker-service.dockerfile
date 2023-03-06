# base go image
FROM golang:1.19-alpine as builder

# create inside docker
RUN mkdir /app

# copy everythink to app from .
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp



# build a tiny docker image
FROM alpine:latest

COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]