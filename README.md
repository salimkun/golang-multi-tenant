# Multi-Tenant Messaging Application

Aplikasi ini adalah implementasi sistem pesan multi-tenant menggunakan Go, RabbitMQ, dan PostgreSQL. Aplikasi ini mendukung pengelolaan konsumen dinamis, penyimpanan data terpartisi, dan pengaturan tingkat konkurensi.

## Fitur

- **Multi-Tenant Messaging**: Mendukung banyak tenant dengan isolasi pesan.
- **Dynamic Consumer Management**: Secara otomatis membuat dan mengelola konsumen RabbitMQ untuk setiap tenant.
- **Partitioned Data Storage**: Menggunakan PostgreSQL untuk menyimpan pesan dengan partisi berdasarkan `tenant_id`.
- **Configurable Concurrency**: Mendukung pengaturan jumlah pekerja untuk setiap tenant.
- **Graceful Shutdown**: Memastikan transaksi yang sedang berjalan selesai sebelum aplikasi berhenti.
- **Cursor Pagination**: Mendukung pengambilan pesan dengan pagination berbasis cursor.

## Struktur Proyek

```
multi-tenant-messaging-app
├── cmd
│   └── main.go                 # Entry point aplikasi
├── internal
│   │── handlers
│   │   ├── tenant.go           # Handler API terkait tenant
│   │   └── message.go          # Handler API terkait pesan
│   ├── config
│   │   ├── config.go           # Manajemen konfigurasi aplikasi
│   │   ├── rabbitmq.go         # Koneksi RabbitMQ
│   │   └── db.go               # Koneksi DB
│   ├── repository
│   │   └── message_repo.go     # Manajemen dan operasi database pesan
│   ├── model
│   │   └── message.go          # Model Struct data
│   ├── payload
│   │   └── tenant.go           # Payload Struct data
│   ├── server
│   │   └── router.go           # Management route API
│   ├── migration
│   │   ├── xxxxx_xxxx.down.go  # Migration down table
│   │   └── xxxxx_xxxx.up.sql   # Migration up table
├── go.mod                      # Definisi modul dan dependensi
├── go.sum                      # Checksum untuk dependensi modul
├── makefile                    # Makefile untuk mempermudah
├── env.example                 # Template ENV and remove .example
└── README.md                   # Dokumentasi proyek
```

## Instruksi Setup

1. **Clone repository**:
   ```
   git clone <repository-url>
   cd multi-tenant-messaging-app
   ```

2. **Install dependensi**:
   ```
   go mod tidy
   ```

3. **Konfigurasi aplikasi**:
   Rename .env.example menjadi .env dan sesuaikan configurasinya dengan detail koneksi RabbitMQ dan PostgreSQL Anda.

4. **Jalankan aplikasi**:
   ```
   go run cmd/main.go
   ```

## Penggunaan

- Aplikasi ini menyediakan RESTful API untuk pengelolaan tenant dan operasi pesan.
- Gunakan alat seperti Postman atau `curl` untuk berinteraksi dengan endpoint API.

### **Endpoint Utama**

1. **Tenant Management**:
   - **Buat Tenant dan Kirim Pesan**: `POST /api/tenants`
   - **Hapus Tenant**: `DELETE /api/tenants/{tenant_id}`
   - **Atur Concurrency**: `PUT /api/tenants/{tenant_id}/config/concurrency`

2. **Messaging**:
   - **Ambil Pesan**: `GET /api/tenants/{tenant_id}/messages?cursor={cursor}`

## Pengujian

Jalankan pengujian otomatis menggunakan:
```
go test ./internal/tests/...
```

## Kontribusi

Kontribusi sangat diterima! Silakan buka issue atau kirim pull request untuk perbaikan atau penambahan fitur.

## Lisensi

Proyek ini dilisensikan di bawah MIT License. Lihat file LICENSE untuk detail lebih lanjut.