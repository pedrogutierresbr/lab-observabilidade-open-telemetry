FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod init service-a

RUN go mod tidy

RUN go build -o service-a .

EXPOSE 8080

CMD ["./service-a"]