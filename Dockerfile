FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
COPY client/client.go ./client/
RUN go build -o app main.go
RUN go build -o client ./client/client.go

EXPOSE 8080

CMD ["./app"]
