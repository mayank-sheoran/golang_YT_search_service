FROM golang:latest
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go mod tidy
RUN go build -o main .
EXPOSE 8090
CMD ["./main"]