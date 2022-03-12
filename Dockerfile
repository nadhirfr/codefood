#build stage
FROM golang:alpine AS builder
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./app ./main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /api/app /app
ENTRYPOINT /app
LABEL Name=codefood Version=0.0.2
EXPOSE 3030
