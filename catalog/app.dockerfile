FROM golang:1.24-alpine3.21 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/fabian-emmanuel/go-ms
COPY go.mod go.sum ./
COPY vendor vendor
COPY catalog catalog
RUN go build -mod=vendor -o /go/bin/app ./catalog/cmd/catalog

FROM alpine:3.21
WORKDIR /usr/bin
COPY --from=build /go/bin/app .
EXPOSE 8080
CMD ["app"]