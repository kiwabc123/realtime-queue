
FROM golang:1.25

WORKDIR /app

# Copy source
COPY . .

# Init module if not exists
RUN [ -f go.mod ] || go mod init queue-api
RUN go mod tidy


# Build Go binary
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
