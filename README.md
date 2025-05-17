# Metro Management System (Backend)

Sistem manajemen perusahaan berbasis multi-company. Fitur:
- Auth (Login & Register dengan JWT)
- Management Karyawan (CRUD)
- Absensi (Foto & Lokasi)
- Pengajuan & Persetujuan Cuti
- Penerbitan Surat Otomatis (HTML + PDF ready)
- Slip Gaji per Bulan

## 🚀 Tech Stack
- Golang + Gin Framework
- MySQL 8
- Docker + Docker Compose
- JWT for auth

## ⚙️ Setup (Dev)

```bash
git clone https://github.com/username/metro-backend.git
cd metro-backend
cp .env.example .env
docker-compose up --build
