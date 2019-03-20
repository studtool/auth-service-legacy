FROM golang:1.12 as base
WORKDIR /tmp/auth-service
COPY . .
RUN go build -o /tmp/service .

FROM alpine:3.9
WORKDIR /tmp
COPY --from=base /tmp/service ./service
ENTRYPOINT ./service
EXPOSE 80