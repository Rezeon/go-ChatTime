```markdown
# Go ChatTime

Backend chat application menggunakan Go, WebSocket, Gin, Claudinary

---

##  Fitur Utama

- Menyediakan WebSocket endpoint untuk komunikasi real-time antar client  
- Broadcast pesan dari server ke semua client yang terkoneksi  
- Struktur modular: controllers, routes, models, middleware, utils, dan ws  
- Integrasi dengan database (opsional tergantung implementasi)  
- Middleware untuk proteksi, autentikasi
- Upload image ke claudinary

---

##  Struktur Direktori

```

├── controllers     # Handler logika bisnis (user, message, dll)
├── database        # Inisialisasi koneksi DB, migrasi
├── docs            # Dokumentasi (contoh API, diagram)
├── middleware      # Middleware: autentikasi, validasi, dll
├── models          # Model / struct entitas
├── routes          # Routing HTTP / WebSocket
├── utils           # Fungsi pembantu umum
├── ws              # Modul WebSocket: koneksi, broadcast
├── main.go         # Entry point aplikasi
├── go.mod
├── .env.example    # Template konfigurasi environment
└── .gitignore

````

---

## 🚀 Cara Menjalankan

1. **Clone repository**
   ```bash
   git clone https://github.com/Rezeon/go-ChatTime.git
   cd go-ChatTime
````

2. **Isi konfigurasi environment**
   Salin `.env.example` ke `.env` dan sesuaikan variabel (misalnya database URL, port, secret key)

3. **Install dependensi**

   ```bash
   go mod download
   ```

4. **Jalankan aplikasi**

   ```bash
   go run main.go
   ```

5. **Gunakan WebSocket endpoint**
   Misalnya: `ws://localhost:8080/ws` (tergantung routing di aplikasi kamu)

---

##  Cara Kerja WebSocket (Modul `ws`)

* `HandleConnections`
  Menerima koneksi WebSocket dari client, menyimpan koneksi ke map `clients`, dan mendengarkan pesan masuk (walau saat ini pesan dari client belum disebarkan).
* `HandleMessage`
  Menunggu pesan dari channel `broadcast`, lalu menuliskan (broadcast) pesan ke semua client aktif.
* `SendToClients(msg interface{})`
  Fungsi publik untuk memasukkan pesan ke channel `broadcast` sehingga nanti akan disebarkan ke client.

---

##  Penyesuaian / Pengembangan

* Forward pesan dari client ke broadcast agar client bisa saling mengobrol
* Tambahkan mekanisme autentikasi (JWT, session) agar hanya user terverifikasi yang boleh connect
* Validasi payload WebSocket (struktur pesan, format JSON, format FormData)
* Monitoring koneksi client (ping/pong)
* Logging lebih baik, error handling yang lebih robust
* Skalabilitas: misalnya clustering, distribusi broadcast antar instance

---

##  Lisensi & Kontribusi

Project ini bebas digunakan. Kalau kamu ingin kontribute:

1. Fork repository
2. Buat branch fitur baru (`feature/…`)
3. Push perubahan
4. Buka pull request
