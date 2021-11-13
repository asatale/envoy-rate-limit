FROM golang:1.17 AS Builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0  go build -o server server.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /go/src/app/bin
COPY --from=Builder /go/src/app/server .

