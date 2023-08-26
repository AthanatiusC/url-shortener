FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o binary cmd/main.go

EXPOSE 8080

ENTRYPOINT ["/app/binary"]