FROM golang:latest

COPY go.mod ./
RUN go mod download

COPY ./ ./
RUN go build server.go

CMD ["sh", "./startup.sh"]
