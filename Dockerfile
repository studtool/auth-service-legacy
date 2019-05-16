FROM golang:1.12.5-alpine3.9 as base
WORKDIR /tmp/auth-service
COPY . .
RUN go build -mod vendor -o /tmp/service .

FROM alpine:3.9.4
WORKDIR /tmp
COPY --from=base /tmp/service ./service
CMD ./service
EXPOSE 80
