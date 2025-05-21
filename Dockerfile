FROM golang:1.24.1

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Install wkhtmltopdf dependencies dan wkhtmltopdf
RUN apt-get update && apt-get install -y \
    xfonts-75dpi \
    xfonts-base \
    wkhtmltopdf

RUN go build -o metro-app ./cmd/main.go

EXPOSE 8080

CMD ["./metro-app"]
