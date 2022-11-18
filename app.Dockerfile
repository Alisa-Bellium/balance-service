FROM golang:1.19

WORKDIR /app/src
COPY . .
RUN go build -o /app/build ./cmd/main.go

EXPOSE 8080

CMD ["/app/build"]