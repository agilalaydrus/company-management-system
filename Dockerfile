FROM golang:1.24.3

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o metro-app ./cmd/main.go

EXPOSE 8080
CMD ["./metro-app"]
