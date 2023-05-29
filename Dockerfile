FROM golang:1.20-alpine as builder
WORKDIR /home
ADD . .
RUN go build -o producer cmd/main.go

FROM ubuntu:20.04

WORKDIR /home
COPY --from=builder /home/producer /home/producer

ENTRYPOINT ["/home/producer"]
