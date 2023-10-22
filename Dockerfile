FROM golang:1.21.3

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /rangeplusclient

EXPOSE 9000

CMD ["/rangeplusclient"]