FROM golang:1.22.7-alpine3.20 AS builder
WORKDIR /build
COPY ./go.mod . 
RUN go mod download
COPY . .
RUN go build -o main cmd/app/main.go

FROM alpine:3.20.3
RUN apk add --no-cache tzdata
COPY ./configs /configs
COPY --from=builder /build/main /bin/main
ENTRYPOINT ["/bin/main"]