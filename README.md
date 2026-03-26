[![Go CI](https://github.com/pconk/auth-service/actions/workflows/ci.yml/badge.svg)](https://github.com/pconk/auth-service/actions/workflows/ci.yml)

### Identity & Access Management (Auth Service) 🔐

**Auth Service** adalah layanan mikro pusat (*Identity Provider*) yang mengelola seluruh siklus hidup pengguna dan keamanan akses di dalam ekosistem *Warehouse*. Layanan ini bertindak sebagai *Single Source of Truth* untuk autentikasi dan otorisasi.

---

🚀 **Fitur Utama**

*   **Centralized Identity Management**: Memusatkan data kredensial pengguna agar dapat digunakan oleh berbagai layanan (Gateway, Audit, dll) tanpa redundansi.
*   **JWT Issue & Validation**: Menggunakan *JSON Web Tokens* (JWT) untuk autentikasi *stateless*. Payload mencakup `user_id`, `username`, dan `role`.
*   **Role-Based Access Control (RBAC)**: Mendefinisikan hak akses secara granular (`admin`, `staff`).
*   **Dual Protocol Support**: 
    *   **HTTP (REST)**: Untuk interaksi publik (Login & Register).
    *   **gRPC**: Untuk komunikasi internal antar-service yang sangat cepat (Token Validation & Profiling).
*   **Bcrypt Security**: Pengamanan password menggunakan algoritma *hashing* Bcrypt yang kuat.
*   **Structured Logging**: Implementasi `slog` dengan pelacakan `request_id` (Correlation ID) di layer HTTP dan gRPC.
*   **Graceful Shutdown**: Menjamin server berhenti dengan aman tanpa memutus koneksi atau transaksi database yang sedang berjalan.

---

🛠️ **Tech Stack**

*   **Language**: Go (Golang) 1.26
*   **Communication**: gRPC (Protobuf) & REST (chi)
*   **Database**: MySQL dengan **GORM** (Object-Relational Mapping)
*   **Security**: JWT (golang-jwt) & Bcrypt
*   **Logging**: slog (Structured Logging) & Google UUID

---

🏗 **Architecture Overview**

Project menggunakan pendekatan **Clean Architecture** untuk memisahkan logika bisnis dari infrastruktur.

```text
      [ Clients / Gateway ]
              │
      ┌───────┴───────┐
      │               │
      ▼               ▼
 [ HTTP 8081 ]   [ gRPC 50051 ]
 (Public API)    (Internal RPC)
      │               │
      └───────┬───────┘
              ▼
    ┌───────────────────┐
    │     Handlers      │ ───▶ Middleware (Logger, Recovery, Interceptors)
    └────────┬──────────┘
             ▼
    ┌───────────────────┐
    │  Service Layer    │ ───▶ (Business Logic, JWT, Bcrypt)
    └────────┬──────────┘
             ▼
    ┌───────────────────┐
    │ Repository Layer  │ ───▶ (GORM / MySQL)
    └───────────────────┘
```

---

📁 **Struktur Project**

```text
.
├── cmd/api/main.go          # Entry point (HTTP & gRPC Runner)
├── internal/
│   ├── config/              # Konfigurasi environment
│   ├── entity/              # GORM Models (User)
│   ├── handler/             # HTTP & gRPC Handlers
│   ├── helper/              # Standard Web Response
│   ├── middleware/          # Logger (HTTP) & Interceptor (gRPC)
│   ├── pb/                  # Generated code dari Protobuf
│   ├── repository/          # Database Access Layer
│   └── service/             # Auth Logic (Login, Register, Validate)
├── pkg/
│   └── database/            # MySQL Connection driver
├── proto/                   # File definisi Protobuf (.proto)
└── deskripsi-auth-service.txt # Referensi dokumentasi sistem
```

---

🚥 **API Contract**

**Public Endpoints (HTTP - Port 8081)**
| Method | Endpoint         | Description                        |
| :----- | :--------------- | :--------------------------------- |
| POST   | `/auth/register` | Mendaftarkan pengguna baru         |
| POST   | `/auth/login`    | Menukar kredensial dengan JWT      |

**Internal Endpoints (gRPC - Port 50051)**
| RPC              | Request              | Description                               |
| :--------------- | :------------------- | :---------------------------------------- |
| `ValidateToken`  | `TokenRequest`       | Memverifikasi JWT dan return user info    |
| `GetUserProfile` | `GetUserRequest`     | Mengambil detail user berdasarkan ID      |

#### Example Responses

**HTTP Responses:**

*   **/auth/login**
    ```json
    {
      "code": 200,
      "status": "OK",
      "message": "Login successful",
      "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluX2d1ZGFuZyIsInJvbGUiOiJhZG1pbiIsImlzcyI6ImF1dGgtc2VydmljZSIsImV4cCI6MTc3NDU4MzMxNX0.URIbPxvQo93a7VS1cZ3Vi6MbCnDlu0YeliZ8RiklPsM"
      }
    }
    ```
*   **/auth/register**
    ```json
    {
      "code": 201,
      "status": "Created",
      "message": "User registered successfully"
    }
    ```

**gRPC Responses:**

*   **--- 2. Testing gRPC: GetUserProfile (ID=1) ---**
    ```json
    {
      "id": 1,
      "username": "admin_gudang",
      "role": "admin",
      "position": "Head of Warehouse",
      "email": "admin@warehouse.com"
    }
    ```
*   **--- 3. Testing gRPC: ValidateToken ---**
    ```json
    {
      "id": 1,
      "username": "admin_gudang",
      "role": "admin"
    }
    ```

---

⚙️ **Cara Menjalankan**

1.  **Persiapan Database**:
    Pastikan MySQL berjalan dan buat database bernama `warehouse_auth`.

2.  **Environment Setup**:
    Sesuaikan `.env` atau konfigurasi di `internal/config/config.go`. 
    Isi .env
    ```bash
    DB_HOST=localhost
    DB_PORT=3306
    DB_USER=root
    DB_PASSWORD=secret
    DB_NAME=warehouse_auth
    JWT_SECRET=supersecretkey
    APP_PORT=8081
    GRPC_PORT=50051
    ```

3.  **Generate Protobuf** (Opsional jika ada perubahan):
    ```bash
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/auth.proto
    ```

4.  **Jalankan Aplikasi**:
    ```bash
    go mod tidy
    go run cmd/api/main.go
    ```

5.  **Testing**:
    Gunakan file `api_http.test` (REST Client) atau `test-grpcurl.ps1` (PowerShell) untuk verifikasi server.