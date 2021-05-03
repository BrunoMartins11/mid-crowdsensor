FROM golang:1.16.3-alpine3.13

RUN mkdir /app

ADD . /app

WORKDIR /app/cmd

RUN go mod download
COPY .env .env

RUN go build -o main .

CMD ["/app/cmd/main"]
