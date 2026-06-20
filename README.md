# JumpaPay — Aplikasi Order Perpanjangan STNK

Aplikasi full-stack untuk layanan perpanjangan STNK online. Customer login dengan Google, submit order beserta dokumen (KTP & STNK), dan memantau status. Admin mengelola order dan mengubah statusnya lewat dashboard.

**Live Demo:** https://miniproject.onelabs.my.id

---

## Daftar Isi

- [Fitur](#fitur)
- [Tech Stack](#tech-stack)
- [Arsitektur](#arsitektur)
- [Struktur Folder](#struktur-folder)
- [Setup Lokal](#setup-lokal)
- [Environment Variables](#environment-variables)
- [Cara Setup Google OAuth](#cara-setup-google-oauth)
- [Setup Midtrans (Bonus)](#setup-midtrans-bonus)
- [API Endpoints](#api-endpoints)
- [Status Flow Order](#status-flow-order)
- [Deployment (Docker)](#deployment-docker)

---

## Fitur

- ✅ **Login Google OAuth 2.0** — session via JWT (httpOnly cookie)
- ✅ **Form order** — validasi nomor WA, plat kendaraan, 5 digit nomor rangka
- ✅ **Upload dokumen** — KTP & STNK (JPG/PNG, maks 2MB) ke object storage
- ✅ **Tracking order** — status real-time + timeline perubahan status
- ✅ **Admin dashboard** — daftar order, filter status, detail + foto, ubah status
- ✅ **Proteksi route** — admin & customer punya akses berbeda (admin via whitelist email)
- ⭐ **Bonus: Midtrans Sandbox** — pembayaran QRIS/VA, auto-update status saat lunas

---

## Tech Stack

| Layer | Teknologi |
|---|---|
| Frontend | React 19 + Vite + TypeScript + Tailwind CSS |
| Backend | Go 1.25 + Gin (REST API) |
| Database | PostgreSQL 16 |
| Auth | Google OAuth 2.0 + JWT |
| Storage | MinIO (S3-compatible) |
| Payment | Midtrans Snap (Sandbox) |
| Deploy | Docker Compose + Nginx |

---

## Arsitektur

```
                    ┌─────────────────────────────┐
   Internet ───────▶│  Nginx Proxy Manager (HTTPS) │
                    └──────────────┬──────────────┘
                                   │ :8877
                    ┌──────────────▼──────────────┐
                    │  frontend (nginx)            │
                    │  - serve React build (/)     │
                    │  - proxy /api,/auth,/payment │
                    └──────────────┬──────────────┘
                                   │ http://backend:8080 (internal network)
                    ┌──────────────▼──────────────┐
                    │  backend (Go + Gin)          │
                    └──────┬───────────────┬───────┘
                           │               │
                  ┌────────▼──────┐  ┌─────▼──────────┐
                  │ PostgreSQL    │  │ MinIO (extern) │
                  │ (container)   │  │ Midtrans API   │
                  └───────────────┘  └────────────────┘
```

Backend tidak meng-expose port ke host — frontend (nginx) yang menjadi satu-satunya pintu masuk dan mem-proxy request API ke backend lewat Docker network internal.

---

## Struktur Folder

```
miniproject/
├── backend/
│   ├── cmd/main.go              # entrypoint, wiring & routes
│   ├── internal/
│   │   ├── handler/             # HTTP layer (thin)
│   │   ├── service/             # business logic
│   │   ├── repository/          # query DB
│   │   ├── model/               # domain structs
│   │   └── middleware/          # JWT auth, admin guard, CORS
│   ├── pkg/
│   │   ├── database/            # koneksi PostgreSQL (pgx)
│   │   ├── storage/             # MinIO client
│   │   ├── oauth/               # Google OAuth helper
│   │   └── payment/             # Midtrans Snap client
│   ├── migrations/              # SQL migration (auto-run di Docker)
│   ├── Dockerfile
│   └── .env.example
├── frontend/
│   ├── src/
│   │   ├── pages/               # Login, OrderForm, Tracking, Payment, admin/
│   │   ├── components/          # Layout, ProtectedRoute, StatusBadge, FileUpload
│   │   ├── hooks/               # useAuth + AuthProvider
│   │   └── lib/                 # axios, types, validation, format, snap
│   ├── nginx.conf              # serve static + reverse proxy
│   ├── Dockerfile
│   └── .env.example
└── docker-compose.yml
```

---

## Setup Lokal

### Prasyarat
- Go 1.25+
- Node.js 20+
- PostgreSQL 16 (lokal) dengan database bernama `miniproject`
- Instance MinIO (atau S3-compatible) — opsional untuk fitur upload
- Akun Google Cloud (untuk OAuth)

### 1. Clone & siapkan env

```bash
git clone https://github.com/kurniawa9157/miniproject-interview.git
cd miniproject-interview

# Backend
cp backend/.env.example backend/.env
# Frontend
cp frontend/.env.example frontend/.env
```

Isi `backend/.env` dan `frontend/.env` (lihat [Environment Variables](#environment-variables)).

### 2. Database

Buat database lalu jalankan migration:

```bash
createdb miniproject   # atau via psql/pgAdmin

psql -U postgres -d miniproject -f backend/migrations/001_init.sql
psql -U postgres -d miniproject -f backend/migrations/002_payment.sql
```

### 3. Backend

```bash
cd backend
go mod download
go run cmd/main.go
# Server berjalan di http://localhost:8080
```

### 4. Frontend

```bash
cd frontend
npm install
npm run dev
# App berjalan di http://localhost:5173
```

Buka http://localhost:5173 dan login dengan Google.

---

## Environment Variables

### `backend/.env`

| Variable | Deskripsi |
|---|---|
| `PORT` | Port backend (default 8080) |
| `APP_ENV` | `development` / `production` (produksi → cookie Secure/HTTPS) |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | Koneksi PostgreSQL |
| `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` | Untuk service `db` di Docker (samakan dengan `DB_*`) |
| `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` | Kredensial Google OAuth |
| `GOOGLE_REDIRECT_URL` | Callback OAuth (harus cocok dengan Google Console) |
| `JWT_SECRET` | Secret untuk sign JWT (min 32 karakter) |
| `ADMIN_EMAILS` | Whitelist email admin (pisah koma) |
| `MINIO_ENDPOINT`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `MINIO_BUCKET`, `MINIO_USE_SSL` | Konfigurasi storage |
| `FRONTEND_URL` | URL frontend (untuk CORS & redirect setelah login) |
| `MIDTRANS_SERVER_KEY`, `MIDTRANS_CLIENT_KEY` | Kredensial Midtrans (bonus) |

### `frontend/.env`

| Variable | Deskripsi |
|---|---|
| `VITE_API_URL` | Base URL backend. **Dev:** `http://localhost:8080`. **Produksi:** kosongkan (`""`) agar request relatif & diproxy nginx |
| `VITE_GOOGLE_CLIENT_ID` | Google Client ID (opsional di frontend) |

> File `.env` tidak di-commit. Gunakan `.env.example` sebagai template.

---

## Cara Setup Google OAuth

1. Buka [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
2. Buat **OAuth Client ID** → tipe **Web application**
3. **Authorized JavaScript origins:**
   - Lokal: `http://localhost:5173`
   - Produksi: `https://miniproject.onelabs.my.id`
4. **Authorized redirect URIs:**
   - Lokal: `http://localhost:8080/auth/google/callback`
   - Produksi: `https://miniproject.onelabs.my.id/auth/google/callback`
5. Salin **Client ID** & **Client Secret** ke `backend/.env`
6. Pastikan `GOOGLE_REDIRECT_URL` di `.env` **sama persis** dengan redirect URI di Console
7. Jika OAuth consent screen masih "Testing", tambahkan email kamu sebagai **Test user**

> **Admin** ditentukan oleh `ADMIN_EMAILS` di `.env`. Email yang login dan ada di daftar ini otomatis jadi admin.

---

## Setup Midtrans (Bonus)

1. Daftar di [dashboard.midtrans.com](https://dashboard.midtrans.com) → pilih environment **Sandbox**
2. **Settings → Access Keys** → salin **Server Key** & **Client Key** ke `backend/.env`
3. **Settings → Configuration → Payment Notification URL:**
   ```
   https://<domain>/payment/notification
   ```
4. Flow: submit order → halaman konfirmasi (Rp 150.000) → Snap popup → bayar → webhook → status order otomatis `PENDING → IN_PROCESS`, payment `PAID`

Kartu test sandbox: `4811 1111 1111 1114`, CVV bebas, OTP `112233`. Atau QRIS sandbox.

---

## API Endpoints

### Auth
| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | `/auth/google` | Redirect ke Google |
| GET | `/auth/google/callback` | Callback, set JWT cookie |
| POST | `/auth/logout` | Hapus session |
| GET | `/auth/me` | Info user login |

### Customer (perlu auth)
| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | `/api/orders` | Submit order (multipart/form-data) |
| GET | `/api/orders` | List order milik sendiri |
| GET | `/api/orders/:id` | Detail + tracking (hanya milik sendiri) |
| GET | `/api/payment/config` | Client key Midtrans |
| POST | `/api/orders/:id/pay` | Buat transaksi pembayaran (Snap token) |

### Admin (perlu auth + whitelist)
| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | `/api/admin/orders?status=` | List semua order + filter |
| GET | `/api/admin/orders/:id` | Detail order |
| PATCH | `/api/admin/orders/:id/status` | Ubah status (validasi transisi) |

### Webhook
| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | `/payment/notification` | Webhook Midtrans (diverifikasi signature) |

**Validasi order:**
- Nomor WA: angka, diawali `08`/`+62`, min 10 digit
- Nomor plat: format Indonesia (`D 1234 ABC`)
- Nomor rangka: tepat 5 karakter alphanumeric
- KTP/STNK: JPG/PNG, maks 2MB

---

## Status Flow Order

```
PENDING ──▶ IN_PROCESS ──▶ DONE
   │             │
   └──▶ CANCELLED ◀──┘
```

- Transisi valid: `PENDING→IN_PROCESS`, `PENDING→CANCELLED`, `IN_PROCESS→DONE`, `IN_PROCESS→CANCELLED`
- `DONE` & `CANCELLED` bersifat **final** (tidak bisa diubah)
- Status tidak bisa mundur atau melompat
- Hanya admin yang bisa mengubah status
- Order ID format: `JP-YYYYMMDD-XXXX`

---

## Deployment (Docker)

Deployment menggunakan Docker Compose: PostgreSQL + backend + frontend (nginx). DB & MinIO bisa internal (container) atau eksternal.

```bash
# Di server
git clone https://github.com/kurniawa9157/miniproject-interview.git
cd miniproject-interview

# Siapkan env (set APP_ENV=production, domain produksi, dll)
cp backend/.env.example backend/.env
nano backend/.env

# Build & jalankan
docker compose up -d --build
```

Hal penting saat produksi:
- `APP_ENV=production`, `FRONTEND_URL` & `GOOGLE_REDIRECT_URL` pakai domain HTTPS
- `POSTGRES_*` di `backend/.env` harus sama dengan `DB_*` (service `db` membuat database darinya)
- Migration jalan **otomatis** saat volume PostgreSQL pertama kali dibuat
- Frontend di-publish di port `8877` — arahkan reverse proxy (mis. Nginx Proxy Manager) ke `<ip-server>:8877`
- Daftarkan redirect URI produksi di Google Console & Notification URL di Midtrans

**Update aplikasi:**
```bash
git pull
docker compose up -d --build
```
