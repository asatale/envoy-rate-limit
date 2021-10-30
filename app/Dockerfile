FROM golang:1.17 AS Builder
WORKDIR /go/src/app
COPY . .
RUN make all

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /go/src/app/bin
COPY --from=Builder /go/src/app/bin .

