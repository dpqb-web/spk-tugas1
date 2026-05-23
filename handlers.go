package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"math"
	"net/http"
	"sort"
)

//go:embed public
var publicPath embed.FS

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// ─── Kriteria ──────────────────────────────────────────────────────────────

func handleGetKriteria(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT kode, nama, tipe, bobot FROM kriteria ORDER BY kode")
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	list := []Kriteria{}
	for rows.Next() {
		var k Kriteria
		if err := rows.Scan(&k.Kode, &k.Nama, &k.Tipe, &k.Bobot); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
		list = append(list, k)
	}
	writeJSON(w, list)
}

func handleAddKriteria(w http.ResponseWriter, r *http.Request) {
	var k Kriteria
	if err := json.NewDecoder(r.Body).Decode(&k); err != nil {
		writeError(w, "invalid JSON", 400)
		return
	}
	_, err := db.Exec("INSERT INTO kriteria (kode, nama, tipe, bobot) VALUES (?,?,?,?)",
		k.Kode, k.Nama, k.Tipe, k.Bobot)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	recalculate()
	writeJSON(w, map[string]bool{"ok": true})
}

func handleUpdateKriteria(w http.ResponseWriter, r *http.Request) {
	kode := r.PathValue("kode")
	var k Kriteria
	if err := json.NewDecoder(r.Body).Decode(&k); err != nil {
		writeError(w, "invalid JSON", 400)
		return
	}
	_, err := db.Exec("UPDATE kriteria SET nama=?, tipe=?, bobot=? WHERE kode=?",
		k.Nama, k.Tipe, k.Bobot, kode)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	recalculate()
	writeJSON(w, map[string]bool{"ok": true})
}

func handleDeleteKriteria(w http.ResponseWriter, r *http.Request) {
	kode := r.PathValue("kode")
	tx, err := db.Begin()
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM kriteria WHERE kode=?", kode)
	tx.Exec("DELETE FROM penilaian WHERE kriteria_kode=?", kode)

	if err := tx.Commit(); err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	recalculate()
	writeJSON(w, map[string]bool{"ok": true})
}

// ─── Alternatif ────────────────────────────────────────────────────────────

func handleGetAlternatif(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT kode, nama FROM alternatif ORDER BY kode")
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	list := []Alternatif{}
	for rows.Next() {
		var a Alternatif
		if err := rows.Scan(&a.Kode, &a.Nama); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
		list = append(list, a)
	}
	writeJSON(w, list)
}

func handleAddAlternatif(w http.ResponseWriter, r *http.Request) {
	var a Alternatif
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeError(w, "invalid JSON", 400)
		return
	}
	_, err := db.Exec("INSERT INTO alternatif (kode, nama) VALUES (?,?)", a.Kode, a.Nama)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	recalculate()
	writeJSON(w, map[string]bool{"ok": true})
}

func handleUpdateAlternatif(w http.ResponseWriter, r *http.Request) {
	kode := r.PathValue("kode")
	var a Alternatif
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeError(w, "invalid JSON", 400)
		return
	}
	_, err := db.Exec("UPDATE alternatif SET nama=? WHERE kode=?", a.Nama, kode)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	recalculate()
	writeJSON(w, map[string]bool{"ok": true})
}

func handleDeleteAlternatif(w http.ResponseWriter, r *http.Request) {
	kode := r.PathValue("kode")
	tx, err := db.Begin()
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM alternatif WHERE kode=?", kode)
	tx.Exec("DELETE FROM penilaian WHERE alternatif_kode=?", kode)

	if err := tx.Commit(); err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	recalculate()
	writeJSON(w, map[string]bool{"ok": true})
}

// ─── Penilaian ─────────────────────────────────────────────────────────────

func handleGetPenilaian(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(
		"SELECT alternatif_kode, kriteria_kode, nilai FROM penilaian ORDER BY alternatif_kode, kriteria_kode")
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	list := []Penilaian{}
	for rows.Next() {
		var p Penilaian
		if err := rows.Scan(&p.AlternatifKode, &p.KriteriaKode, &p.Nilai); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
		list = append(list, p)
	}
	writeJSON(w, list)
}

func handleSavePenilaian(w http.ResponseWriter, r *http.Request) {
	var items []Penilaian
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		writeError(w, "invalid JSON", 400)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO penilaian (alternatif_kode, kriteria_kode, nilai)
		VALUES (?,?,?) ON CONFLICT(alternatif_kode, kriteria_kode) DO UPDATE SET nilai=excluded.nilai`)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	defer stmt.Close()

	for _, item := range items {
		if _, err := stmt.Exec(item.AlternatifKode, item.KriteriaKode, item.Nilai); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	recalculate()
	writeJSON(w, map[string]bool{"ok": true})
}

// ─── Hasil ─────────────────────────────────────────────────────────────────

func handleGetHasil(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT h.alternatif_kode, a.nama, h.vektor_s, h.vektor_v, h.ranking
		FROM hasil h JOIN alternatif a ON h.alternatif_kode=a.kode ORDER BY h.ranking`)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	list := []Hasil{}
	for rows.Next() {
		var h Hasil
		if err := rows.Scan(&h.AlternatifKode, &h.Nama, &h.VektorS, &h.VektorV, &h.Ranking); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
		list = append(list, h)
	}
	writeJSON(w, list)
}

func handleGetNormalisasi(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT k.kode, k.nama, k.tipe, k.bobot, (k.bobot / total.s) AS bobot_norm
		FROM kriteria k, (SELECT SUM(bobot) AS s FROM kriteria) total ORDER BY k.kode`)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	list := []Normalisasi{}
	for rows.Next() {
		var n Normalisasi
		if err := rows.Scan(&n.Kode, &n.Nama, &n.Tipe, &n.Bobot, &n.BobotNorm); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
		list = append(list, n)
	}
	writeJSON(w, list)
}

func handleGetVektorSDetail(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT alternatif_kode, kriteria_kode, pangkat, komponen FROM vektor_s_detail ORDER BY alternatif_kode, kriteria_kode")
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	list := []VektorSDetail{}
	for rows.Next() {
		var d VektorSDetail
		if err := rows.Scan(&d.AlternatifKode, &d.KriteriaKode, &d.Pangkat, &d.Komponen); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
		list = append(list, d)
	}
	writeJSON(w, list)
}

// ─── WP Calculation ────────────────────────────────────────────────────────

func recalculate() {
	kRows, err := db.Query("SELECT kode, nama, tipe, bobot FROM kriteria ORDER BY kode")
	if err != nil {
		log.Printf("recalculate: query kriteria: %v", err)
		return
	}
	defer kRows.Close()

	type kNorm struct {
		Kriteria
		bobotNorm float64
	}

	var kriteria []kNorm
	totalBobot := 0.0
	for kRows.Next() {
		var k Kriteria
		if err := kRows.Scan(&k.Kode, &k.Nama, &k.Tipe, &k.Bobot); err != nil {
			log.Printf("recalculate: scan kriteria: %v", err)
			return
		}
		totalBobot += k.Bobot
		kriteria = append(kriteria, kNorm{Kriteria: k})
	}
	if totalBobot == 0 {
		return
	}

	for i := range kriteria {
		kriteria[i].bobotNorm = kriteria[i].Bobot / totalBobot
	}

	aRows, err := db.Query("SELECT kode, nama FROM alternatif ORDER BY kode")
	if err != nil {
		log.Printf("recalculate: query alternatif: %v", err)
		return
	}
	defer aRows.Close()

	var alternatif []Alternatif
	for aRows.Next() {
		var a Alternatif
		if err := aRows.Scan(&a.Kode, &a.Nama); err != nil {
			log.Printf("recalculate: scan alternatif: %v", err)
			return
		}
		alternatif = append(alternatif, a)
	}

	db.Exec("DELETE FROM vektor_s_detail")

	vektorS := make(map[string]float64)
	for _, alt := range alternatif {
		s := 1.0
		for _, k := range kriteria {
			var nilai float64
			err := db.QueryRow("SELECT nilai FROM penilaian WHERE alternatif_kode=? AND kriteria_kode=?", alt.Kode, k.Kode).Scan(&nilai)
			if err == sql.ErrNoRows {
				nilai = 1
			} else if err != nil {
				log.Printf("recalculate: query penilaian: %v", err)
				return
			}

			pangkat := k.bobotNorm
			if k.Tipe == "Cost" {
				pangkat = -k.bobotNorm
			}
			komponen := math.Pow(nilai, pangkat)

			db.Exec("INSERT INTO vektor_s_detail (alternatif_kode, kriteria_kode, pangkat, komponen) VALUES (?,?,?,?)",
				alt.Kode, k.Kode, pangkat, komponen)
			s *= komponen
		}
		vektorS[alt.Kode] = s
	}

	totalS := 0.0
	for _, s := range vektorS {
		totalS += s
	}

	db.Exec("DELETE FROM hasil")

	type scored struct {
		kode string
		v    float64
	}
	var scoredList []scored
	for kode, s := range vektorS {
		scoredList = append(scoredList, scored{kode, s / totalS})
	}
	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].v > scoredList[j].v
	})

	for rank, sc := range scoredList {
		db.Exec("INSERT INTO hasil (alternatif_kode, vektor_s, vektor_v, ranking) VALUES (?,?,?,?)",
			sc.kode, vektorS[sc.kode], sc.v, rank+1)
	}
}

// ─── Routes ────────────────────────────────────────────────────────────────

func setupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/kriteria", handleGetKriteria)
	mux.HandleFunc("POST /api/kriteria", handleAddKriteria)
	mux.HandleFunc("PUT /api/kriteria/{kode}", handleUpdateKriteria)
	mux.HandleFunc("DELETE /api/kriteria/{kode}", handleDeleteKriteria)

	mux.HandleFunc("GET /api/alternatif", handleGetAlternatif)
	mux.HandleFunc("POST /api/alternatif", handleAddAlternatif)
	mux.HandleFunc("PUT /api/alternatif/{kode}", handleUpdateAlternatif)
	mux.HandleFunc("DELETE /api/alternatif/{kode}", handleDeleteAlternatif)

	mux.HandleFunc("GET /api/penilaian", handleGetPenilaian)
	mux.HandleFunc("POST /api/penilaian", handleSavePenilaian)

	mux.HandleFunc("GET /api/hasil", handleGetHasil)
	mux.HandleFunc("GET /api/normalisasi", handleGetNormalisasi)
	mux.HandleFunc("GET /api/vektor_s_detail", handleGetVektorSDetail)

	webFS, err := fs.Sub(publicPath, "public")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("GET /", http.FileServer(http.FS(webFS)))

	return mux
}
