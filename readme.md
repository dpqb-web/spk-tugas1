# SPK Pemilihan Cloud Provider
### Sistem Pendukung Keputusan – Metode Weighted Product (WP)

Aplikasi desktop standalone berbasis **Go + SQLite + WebView** untuk membantu pengambilan keputusan pemilihan penyedia layanan cloud menggunakan metode **Weighted Product**.

---

## Struktur Proyek

```
spk/
├── main.go              # Entry point: HTTP server + WebView launcher
├── db.go                # Inisialisasi SQLite, model struct, seed data
├── handlers.go          # HTTP handlers & logika WP
├── public/              # Frontend (embedded via embed.FS)
│   ├── index.html       # SPA utama (dikompilasi dari Pug)
│   ├── app.js           # Logic frontend (dikompilasi dari CoffeeScript)
│   ├── style.css        # Styling (dikompilasi dari Sass)
├── templates/           # Source Pug (pre-build)
├── go.mod / go.sum      # Dependensi Go
├── db.sqlite            # Database SQLite (dibuat otomatis saat pertama jalan)
└── readme.md            # Catatan ini
```

---

## Cara Menjalankan

### 1. Pastikan Go 1.21+ sudah terpasang

```bash
go version
```

### 2. Clone & build

```bash
git clone https://github.com/dpqb-web/spk-tugas1.git
cd spk
go build -o spk
```

**Untuk Windows** (sembunyikan jendela terminal):

```bash
go build -ldflags="-H windowsgui" -o spk.exe
```

### 3. Jalankan

```bash
./spk
# atau di Windows:
spk.exe
```

Jendela desktop akan terbuka otomatis menampilkan aplikasi.
Database `db.sqlite` akan dibuat dan diisi data contoh saat pertama kali dijalankan.

---

## Database

Database SQLite (`db.sqlite`) dibuat otomatis saat pertama kali aplikasi dijalankan.
Data awal langsung dimasukkan melalui fungsi `seed()` di `db.go`:

| Tabel              | Keterangan                                 |
|--------------------|--------------------------------------------|
| `kriteria`         | 10 kriteria penilaian beserta bobot & tipe |
| `alternatif`       | 6 penyedia cloud                           |
| `penilaian`        | Nilai tiap alternatif terhadap tiap kriteria |
| `vektor_s_detail`  | Detail perhitungan komponen Vektor S       |
| `hasil`            | Hasil akhir: Vektor S, V, dan Ranking      |

---

## Fitur Aplikasi

| Halaman       | Fitur                                                              |
|---------------|--------------------------------------------------------------------|
| Dashboard     | Statistik ringkasan + kartu ranking dengan progress bar            |
| Kriteria      | CRUD kriteria: kode, nama, tipe (Benefit/Cost), bobot             |
| Alternatif    | CRUD penyedia cloud                                                |
| Penilaian     | Input matriks nilai (skala 1–10) secara tabel                     |
| Perhitungan   | Detail 3 langkah WP: normalisasi bobot, Vektor S, Vektor V        |
| Hasil         | Ranking akhir + rekomendasi terbaik                               |

---

## Metode: Weighted Product (WP)

```
1. Normalisasi Bobot  →  W*j = Wj / ΣWj

2. Vektor S           →  Si = Π (Xij ^ W*j)
                          pangkat (+) untuk Benefit
                          pangkat (−) untuk Cost

3. Vektor V           →  Vi = Si / ΣSi
                          Alternatif terbaik = Vi tertinggi
```

---

## Build Catatan

| Platform   | Perintah                                |
|------------|-----------------------------------------|
| Linux      | `go build -o spk`                      |
| macOS      | `go build -o spk`                      |
| Windows    | `go build -ldflags="-H windowsgui" -o spk.exe` |

---

## Troubleshooting

| Masalah                          | Solusi                                              |
|----------------------------------|-----------------------------------------------------|
| `undefined: webview.New`         | Pastikan `go mod tidy` sudah dijalankan             |
| Database corrupt                 | Hapus `db.sqlite`, jalankan ulang aplikasi          |
| Port sudah dipakai               | Program otomatis memilih port kosong (`:0`)         |
| WebView tidak muncul (Linux)     | Install WebKit2GTK: `sudo apt install libwebkit2gtk-4.1-dev` |
| WebView tidak muncul (macOS)     | Pastikan Xcode Command Line Tools terinstall        |
