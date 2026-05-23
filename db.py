import os
import sqlite3

DB_PATH = os.path.join(os.path.dirname(__file__), "db.sqlite")

SCHEMA = """
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
"""

SEED_KRITERIA = [
    ("C1", "Harga / Biaya Langganan", "Cost", 0.15),
    ("C2", "Keamanan & Enkripsi Data", "Benefit", 0.15),
    ("C3", "Ketersediaan / Uptime (SLA)", "Benefit", 0.12),
    ("C4", "Skalabilitas Layanan", "Benefit", 0.12),
    ("C5", "Kecepatan Transfer Data", "Benefit", 0.10),
    ("C6", "Kemudahan Integrasi API", "Benefit", 0.10),
    ("C7", "Dukungan Teknis (Support)", "Benefit", 0.08),
    ("C8", "Kepatuhan Regulasi (Compliance)", "Benefit", 0.08),
    ("C9", "Lokasi Data Center", "Benefit", 0.05),
    ("C10", "Reputasi & Ulasan Pengguna", "Benefit", 0.05),
]

SEED_ALTERNATIF = [
    ("A1", "Amazon Web Services (AWS)"),
    ("A2", "Microsoft Azure"),
    ("A3", "Google Cloud Platform (GCP)"),
    ("A4", "Alibaba Cloud"),
    ("A5", "IBM Cloud"),
    ("A6", "DigitalOcean"),
]

# rows: alt_kode, [C1..C10]
SEED_PENILAIAN = [
    ("A1", [6, 9, 9, 9, 8, 9, 8, 9, 8, 9]),
    ("A2", [7, 9, 9, 9, 8, 9, 9, 9, 7, 9]),
    ("A3", [7, 8, 9, 9, 9, 9, 8, 8, 7, 8]),
    ("A4", [9, 7, 8, 8, 7, 7, 7, 7, 6, 7]),
    ("A5", [5, 8, 8, 7, 7, 7, 9, 9, 6, 7]),
    ("A6", [9, 7, 8, 7, 7, 8, 7, 6, 5, 7]),
]


def _conn():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    conn.execute("PRAGMA foreign_keys = ON")
    return conn


def execute(sql, params=()):
    conn = _conn()
    with conn:
        conn.execute(sql, params)
    conn.close()


def query(sql, params=()):
    conn = _conn()
    rows = conn.execute(sql, params).fetchall()
    conn.close()
    return rows


def query_one(sql, params=()):
    conn = _conn()
    row = conn.execute(sql, params).fetchone()
    conn.close()
    return row


def init():
    conn = _conn()
    with conn:
        conn.executescript(SCHEMA)
    conn.close()

    # Only seed if empty
    if not query("SELECT 1 FROM kriteria LIMIT 1"):
        for k in SEED_KRITERIA:
            execute("INSERT OR IGNORE INTO kriteria VALUES (?,?,?,?)", k)
        for a in SEED_ALTERNATIF:
            execute("INSERT OR IGNORE INTO alternatif VALUES (?,?)", a)
        kriteria_kodes = [k[0] for k in SEED_KRITERIA]
        for alt_kode, vals in SEED_PENILAIAN:
            for k_kode, nilai in zip(kriteria_kodes, vals):
                execute(
                    "INSERT OR IGNORE INTO penilaian VALUES (?,?,?)",
                    (alt_kode, k_kode, nilai),
                )
        # Trigger initial calculation
        import app as _app

        _app._recalculate()
