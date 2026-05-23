package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS kriteria (
    kode        TEXT PRIMARY KEY,
    nama        TEXT NOT NULL,
    tipe        TEXT NOT NULL CHECK(tipe IN ('Benefit','Cost')),
    bobot       REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS alternatif (
    kode        TEXT PRIMARY KEY,
    nama        TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS penilaian (
    alternatif_kode TEXT NOT NULL REFERENCES alternatif(kode),
    kriteria_kode   TEXT NOT NULL REFERENCES kriteria(kode),
    nilai           REAL NOT NULL,
    PRIMARY KEY (alternatif_kode, kriteria_kode)
);

CREATE TABLE IF NOT EXISTS vektor_s_detail (
    alternatif_kode TEXT NOT NULL,
    kriteria_kode   TEXT NOT NULL,
    pangkat         REAL,
    komponen        REAL,
    PRIMARY KEY (alternatif_kode, kriteria_kode)
);

CREATE TABLE IF NOT EXISTS hasil (
    alternatif_kode TEXT PRIMARY KEY,
    vektor_s        REAL,
    vektor_v        REAL,
    ranking         INTEGER
);
`

type Kriteria struct {
	Kode  string  `json:"kode"`
	Nama  string  `json:"nama"`
	Tipe  string  `json:"tipe"`
	Bobot float64 `json:"bobot"`
}

type Alternatif struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type Penilaian struct {
	AlternatifKode string  `json:"alternatif_kode"`
	KriteriaKode   string  `json:"kriteria_kode"`
	Nilai          float64 `json:"nilai"`
}

type Hasil struct {
	AlternatifKode string  `json:"alternatif_kode"`
	Nama           string  `json:"nama"`
	VektorS        float64 `json:"vektor_s"`
	VektorV        float64 `json:"vektor_v"`
	Ranking        int     `json:"ranking"`
}

type Normalisasi struct {
	Kode      string  `json:"kode"`
	Nama      string  `json:"nama"`
	Tipe      string  `json:"tipe"`
	Bobot     float64 `json:"bobot"`
	BobotNorm float64 `json:"bobot_norm"`
}

type VektorSDetail struct {
	AlternatifKode string  `json:"alternatif_kode"`
	KriteriaKode   string  `json:"kriteria_kode"`
	Pangkat        float64 `json:"pangkat"`
	Komponen       float64 `json:"komponen"`
}

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "db.sqlite?_foreign_keys=on")
	if err != nil {
		return err
	}

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM kriteria").Scan(&count)
	if count > 0 {
		return nil
	}

	return seed()
}

func closeDB() {
	if db != nil {
		db.Close()
	}
}

func seed() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	kriteria := []struct {
		kode, nama, tipe string
		bobot            float64
	}{
		{"C1", "Harga / Biaya Langganan", "Cost", 0.15},
		{"C2", "Keamanan & Enkripsi Data", "Benefit", 0.15},
		{"C3", "Ketersediaan / Uptime (SLA)", "Benefit", 0.12},
		{"C4", "Skalabilitas Layanan", "Benefit", 0.12},
		{"C5", "Kecepatan Transfer Data", "Benefit", 0.10},
		{"C6", "Kemudahan Integrasi API", "Benefit", 0.10},
		{"C7", "Dukungan Teknis (Support)", "Benefit", 0.08},
		{"C8", "Kepatuhan Regulasi (Compliance)", "Benefit", 0.08},
		{"C9", "Lokasi Data Center", "Benefit", 0.05},
		{"C10", "Reputasi & Ulasan Pengguna", "Benefit", 0.05},
	}

	for _, k := range kriteria {
		_, err = tx.Exec("INSERT OR IGNORE INTO kriteria VALUES (?,?,?,?)", k.kode, k.nama, k.tipe, k.bobot)
		if err != nil {
			return err
		}
	}

	alternatifs := []struct{ kode, nama string }{
		{"A1", "Amazon Web Services (AWS)"},
		{"A2", "Microsoft Azure"},
		{"A3", "Google Cloud Platform (GCP)"},
		{"A4", "Alibaba Cloud"},
		{"A5", "IBM Cloud"},
		{"A6", "DigitalOcean"},
	}

	for _, a := range alternatifs {
		_, err = tx.Exec("INSERT OR IGNORE INTO alternatif VALUES (?,?)", a.kode, a.nama)
		if err != nil {
			return err
		}
	}

	type seedPenilaian struct {
		altKode string
		values  [10]float64
	}

	penilaians := []seedPenilaian{
		{"A1", [10]float64{6, 9, 9, 9, 8, 9, 8, 9, 8, 9}},
		{"A2", [10]float64{7, 9, 9, 9, 8, 9, 9, 9, 7, 9}},
		{"A3", [10]float64{7, 8, 9, 9, 9, 9, 8, 8, 7, 8}},
		{"A4", [10]float64{9, 7, 8, 8, 7, 7, 7, 7, 6, 7}},
		{"A5", [10]float64{5, 8, 8, 7, 7, 7, 9, 9, 6, 7}},
		{"A6", [10]float64{9, 7, 8, 7, 7, 8, 7, 6, 5, 7}},
	}

	kKodes := []string{"C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9", "C10"}
	for _, p := range penilaians {
		for i, kk := range kKodes {
			_, err = tx.Exec("INSERT OR IGNORE INTO penilaian VALUES (?,?,?)", p.altKode, kk, p.values[i])
			if err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	recalculate()
	log.Println("Database seeded successfully")
	return nil
}
