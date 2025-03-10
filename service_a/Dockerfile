FROM golang:1.24.1

WORKDIR /app

COPY . .

RUN if [ ! -f go.mod ]; then go mod init service_a; fi

RUN go mod tidy

RUN go build -o service_a .

EXPOSE 8080

CMD ["./service_a"]