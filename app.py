import webview
from flask import Flask, jsonify, render_template, request

import db

app = Flask(__name__)
app.secret_key = "spk_cloud_secret"


@app.route("/")
def index():
    return render_template("index.html")


# ─── Kriteria ────────────────────────────────────────────────────────────────


@app.route("/api/kriteria", methods=["GET"])
def get_kriteria():
    rows = db.query("SELECT * FROM kriteria ORDER BY kode")
    return jsonify([dict(r) for r in rows])


@app.route("/api/kriteria", methods=["POST"])
def add_kriteria():
    d = request.json
    db.execute(
        "INSERT INTO kriteria (kode, nama, tipe, bobot) VALUES (?,?,?,?)",
        (d["kode"], d["nama"], d["tipe"], float(d["bobot"])),
    )
    _recalculate()
    return jsonify({"ok": True})


@app.route("/api/kriteria/<kode>", methods=["PUT"])
def update_kriteria(kode):
    d = request.json
    db.execute(
        "UPDATE kriteria SET nama=?, tipe=?, bobot=? WHERE kode=?",
        (d["nama"], d["tipe"], float(d["bobot"]), kode),
    )
    _recalculate()
    return jsonify({"ok": True})


@app.route("/api/kriteria/<kode>", methods=["DELETE"])
def delete_kriteria(kode):
    db.execute("DELETE FROM kriteria WHERE kode=?", (kode,))
    db.execute("DELETE FROM penilaian WHERE kriteria_kode=?", (kode,))
    _recalculate()
    return jsonify({"ok": True})


# ─── Alternatif ──────────────────────────────────────────────────────────────


@app.route("/api/alternatif", methods=["GET"])
def get_alternatif():
    rows = db.query("SELECT * FROM alternatif ORDER BY kode")
    return jsonify([dict(r) for r in rows])


@app.route("/api/alternatif", methods=["POST"])
def add_alternatif():
    d = request.json
    db.execute(
        "INSERT INTO alternatif (kode, nama) VALUES (?,?)", (d["kode"], d["nama"])
    )
    _recalculate()
    return jsonify({"ok": True})


@app.route("/api/alternatif/<kode>", methods=["PUT"])
def update_alternatif(kode):
    d = request.json
    db.execute("UPDATE alternatif SET nama=? WHERE kode=?", (d["nama"], kode))
    _recalculate()
    return jsonify({"ok": True})


@app.route("/api/alternatif/<kode>", methods=["DELETE"])
def delete_alternatif(kode):
    db.execute("DELETE FROM alternatif WHERE kode=?", (kode,))
    db.execute("DELETE FROM penilaian WHERE alternatif_kode=?", (kode,))
    _recalculate()
    return jsonify({"ok": True})


# ─── Penilaian ───────────────────────────────────────────────────────────────


@app.route("/api/penilaian", methods=["GET"])
def get_penilaian():
    rows = db.query(
        "SELECT p.alternatif_kode, p.kriteria_kode, p.nilai "
        "FROM penilaian p "
        "ORDER BY p.alternatif_kode, p.kriteria_kode"
    )
    return jsonify([dict(r) for r in rows])


@app.route("/api/penilaian", methods=["POST"])
def save_penilaian():
    """Bulk upsert: body = [{alternatif_kode, kriteria_kode, nilai}, ...]"""
    items = request.json
    for item in items:
        db.execute(
            "INSERT INTO penilaian (alternatif_kode, kriteria_kode, nilai) "
            "VALUES (?,?,?) "
            "ON CONFLICT(alternatif_kode, kriteria_kode) DO UPDATE SET nilai=excluded.nilai",
            (item["alternatif_kode"], item["kriteria_kode"], float(item["nilai"])),
        )
    _recalculate()
    return jsonify({"ok": True})


# ─── Hasil ───────────────────────────────────────────────────────────────────


@app.route("/api/hasil", methods=["GET"])
def get_hasil():
    rows = db.query(
        "SELECT h.alternatif_kode, a.nama, h.vektor_s, h.vektor_v, h.ranking "
        "FROM hasil h JOIN alternatif a ON h.alternatif_kode=a.kode "
        "ORDER BY h.ranking"
    )
    return jsonify([dict(r) for r in rows])


@app.route("/api/normalisasi", methods=["GET"])
def get_normalisasi():
    rows = db.query(
        "SELECT k.kode, k.nama, k.tipe, k.bobot, "
        "(k.bobot / total.s) AS bobot_norm "
        "FROM kriteria k, (SELECT SUM(bobot) AS s FROM kriteria) total "
        "ORDER BY k.kode"
    )
    return jsonify([dict(r) for r in rows])


@app.route("/api/vektor_s_detail", methods=["GET"])
def get_vektor_s_detail():
    rows = db.query(
        "SELECT * FROM vektor_s_detail ORDER BY alternatif_kode, kriteria_kode"
    )
    return jsonify([dict(r) for r in rows])


# ─── Kalkulasi WP ────────────────────────────────────────────────────────────


def _recalculate():
    kriteria = [dict(r) for r in db.query("SELECT * FROM kriteria ORDER BY kode")]
    alternatif = [dict(r) for r in db.query("SELECT * FROM alternatif ORDER BY kode")]

    total_bobot = sum(k["bobot"] for k in kriteria)
    if total_bobot == 0:
        return

    # Normalisasi bobot
    for k in kriteria:
        k["bobot_norm"] = k["bobot"] / total_bobot

    # Hitung Vektor S
    db.execute("DELETE FROM vektor_s_detail")
    vektor_s = {}
    for alt in alternatif:
        s = 1.0
        for k in kriteria:
            row = db.query_one(
                "SELECT nilai FROM penilaian WHERE alternatif_kode=? AND kriteria_kode=?",
                (alt["kode"], k["kode"]),
            )
            nilai = row["nilai"] if row else 1
            pangkat = (
                -k["bobot_norm"] if k["tipe"].upper() == "COST" else k["bobot_norm"]
            )
            komponen = nilai**pangkat
            db.execute(
                "INSERT INTO vektor_s_detail (alternatif_kode, kriteria_kode, pangkat, komponen) "
                "VALUES (?,?,?,?)",
                (alt["kode"], k["kode"], pangkat, komponen),
            )
            s *= komponen
        vektor_s[alt["kode"]] = s

    total_s = sum(vektor_s.values())

    # Simpan hasil
    db.execute("DELETE FROM hasil")
    vektor_v = {kode: s / total_s for kode, s in vektor_s.items()}
    ranking = sorted(vektor_v, key=vektor_v.get, reverse=True)
    for rank, kode in enumerate(ranking, 1):
        db.execute(
            "INSERT INTO hasil (alternatif_kode, vektor_s, vektor_v, ranking) VALUES (?,?,?,?)",
            (kode, vektor_s[kode], vektor_v[kode], rank),
        )


# ─── Entry point ─────────────────────────────────────────────────────────────


if __name__ == "__main__":
    db.init()
    webview.create_window(
        "Sistem Pendukung Keputusan – Pemilihan Penyedia Layanan Cloud",
        app,
        min_size=(900, 600)
    )
    webview.start(ssl=True)
