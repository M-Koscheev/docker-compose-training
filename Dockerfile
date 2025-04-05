FROM golang:1.23-alpine
RUN go version
ENV GOPATH /
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

# build go app
RUN go build -o docker-compose-training ./cmd/main.go

RUN chmod +x docker-compose-training

EXPOSE 8080

CMD ["./docker-compose-training"]
