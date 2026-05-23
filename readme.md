# SPK Pemilihan Cloud Provider
### Sistem Pendukung Keputusan – Metode Weighted Product (WP)

Aplikasi web standalone berbasis **Python + Flask + SQLite + pywebview** untuk membantu pengambilan keputusan pemilihan penyedia layanan cloud menggunakan metode **Weighted Product**.

---

## 📦 Struktur Proyek

```
spk_cloud/
├── app.py              # Entry point: Flask routes + pywebview launcher
├── db.py               # Inisialisasi & helper SQLite
├── spk_cloud.db        # Database SQLite (dibuat otomatis saat pertama jalan)
├── requirements.txt    # Dependensi Python
├── templates/
│   └── index.html      # UI lengkap (Single Page App)
└── README.md
```

---

## 🚀 Cara Menjalankan

### 1. Pastikan Python 3.9+ sudah terpasang

```bash
python --version
# Python 3.9.x atau lebih baru
```

### 2. Buat virtual environment (opsional tapi direkomendasikan)

```bash
python -m venv venv

# Aktivasi di Windows:
venv\Scripts\activate

# Aktivasi di macOS/Linux:
source venv/bin/activate
```

### 3. Install dependensi

```bash
pip install -r requirements.txt
```

> **Catatan untuk Linux:** pywebview membutuhkan WebKit2GTK. Install dengan:
> ```bash
> # Ubuntu/Debian:
> sudo apt install python3-gi python3-gi-cairo gir1.2-gtk-3.0 gir1.2-webkit2-4.0
>
> # Fedora/RHEL:
> sudo dnf install python3-gobject webkit2gtk3
> ```

### 4. Jalankan aplikasi

```bash
cd spk_cloud
python app.py
```

Jendela desktop akan terbuka otomatis menampilkan aplikasi.

---

## 🌐 Akses via Browser (mode Flask saja, tanpa pywebview)

Jika ingin akses lewat browser biasa tanpa jendela desktop:

```bash
# Edit app.py bagian bawah:
# Ganti webview.start() menjadi:
app.run(host='127.0.0.1', port=5050, debug=True)

# Lalu buka browser ke:
http://127.0.0.1:5050
```

---

## 🗄️ Database

Database SQLite (`spk_cloud.db`) dibuat otomatis pada saat pertama kali aplikasi dijalankan.  
Data awal (seed) dari file Excel langsung dimasukkan ke database:

| Tabel            | Keterangan                                 |
|------------------|--------------------------------------------|
| `kriteria`       | 10 kriteria penilaian beserta bobot & tipe |
| `alternatif`     | 6 penyedia cloud                           |
| `penilaian`      | Nilai tiap alternatif terhadap tiap kriteria |
| `vektor_s_detail`| Detail perhitungan komponen Vektor S       |
| `hasil`          | Hasil akhir: Vektor S, V, dan Ranking      |

---

## 📊 Fitur Aplikasi

| Halaman       | Fitur                                                              |
|---------------|--------------------------------------------------------------------|
| Dashboard     | Statistik ringkasan + kartu ranking dengan progress bar            |
| Kriteria      | CRUD kriteria: kode, nama, tipe (Benefit/Cost), bobot             |
| Alternatif    | CRUD penyedia cloud                                                |
| Penilaian     | Input matriks nilai (skala 1–10) secara tabel                     |
| Perhitungan   | Detail 3 langkah WP: normalisasi bobot, Vektor S, Vektor V        |
| Hasil         | Ranking akhir + rekomendasi terbaik                               |

---

## ⚙️ Metode: Weighted Product (WP)

```
1. Normalisasi Bobot  →  W*j = Wj / ΣWj

2. Vektor S           →  Si = Π (Xij ^ W*j)
                          pangkat (+) untuk Benefit
                          pangkat (−) untuk Cost

3. Vektor V           →  Vi = Si / ΣSi
                          Alternatif terbaik = Vi tertinggi
```

---

## 🔧 Troubleshooting

| Masalah                          | Solusi                                              |
|----------------------------------|-----------------------------------------------------|
| `ModuleNotFoundError: webview`   | `pip install pywebview`                             |
| `ModuleNotFoundError: flask`     | `pip install flask`                                 |
| Jendela tidak muncul (Linux)     | Install WebKit2GTK (lihat langkah 3 di atas)       |
| Port 5050 sudah dipakai          | Ganti port di `app.py` baris `port=5050`           |
| Database corrupt                 | Hapus `spk_cloud.db`, jalankan ulang app           |
