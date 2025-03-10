FROM golang:1.24.1

WORKDIR /app

COPY . .

RUN if [ ! -f go.mod ]; then go mod init service_b; fi

RUN go mod tidy

RUN go build -o service_b ./cmd

EXPOSE 8081

CMD ["./service_b"]