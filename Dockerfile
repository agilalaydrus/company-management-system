FROM golang:1.24.1

WORKDIR /app

# 1. Copy dan download dependency dulu, biar cache efektif
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy semua source setelah dependency selesai di-download
COPY . .

# 3. Gabungkan install apt dan hapus cache supaya layer lebih ringan dan cepat
RUN apt-get update && apt-get install -y \
    xfonts-75dpi \
    xfonts-base \
    wkhtmltopdf \
    && rm -rf /var/lib/apt/lists/*

# 4. Build binary Go
RUN go build -o metro-app ./cmd/main.go

EXPOSE 8080

CMD ["./metro-app"]
