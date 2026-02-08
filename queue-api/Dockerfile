FROM golang:1.23

WORKDIR /app

COPY go.mod ./
COPY . ./

RUN go mod tidy
RUN go mod download

RUN apt-get update && apt-get install -y sqlite3
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]