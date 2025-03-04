FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod init service-b

RUN go mod tidy

RUN go build -o service-b .

EXPOSE 8081

CMD ["./service-b"]