FROM golang:1.25

WORKDIR /app

# Copy source
COPY . .

# Init module if not exists
RUN [ -f go.mod ] || go mod init queue-api
RUN go mod tidy

# Install air (correct repo)
RUN go install github.com/air-verse/air@latest

# Add Go bin to PATH
ENV PATH="/go/bin:${PATH}"

EXPOSE 8080

CMD ["air"]
