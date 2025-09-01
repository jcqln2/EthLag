
FROM golang:1.24.1
WORKDIR /app
COPY go.mod go.sum main.go ./
RUN go mod download
RUN go build -o main main.go
CMD ["./main"]
