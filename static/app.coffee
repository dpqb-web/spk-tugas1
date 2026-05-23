API = ""
_kriteria = []
_alternatif = []
_penilaian = []
_hasil = []
_norm = []
_vsDetail = []

# ── Navigation ──────────────────────────────────────────────────────────────
showPage = (el, id) ->
  for p in document.querySelectorAll ".page"
    p.classList.remove "active"
  for n in document.querySelectorAll ".nav-item"
    n.classList.remove "active"
  document.getElementById("page-#{id}").classList.add "active"
  el.currentTarget.classList.add "active"
  loadPage id

showTab = (el, id) ->
  for t in document.querySelectorAll ".tab"
    t.classList.remove "active"
  el.currentTarget.classList.add "active"
  for t in ["tab-normalisasi", "tab-vektors", "tab-vektorv"]
    document.getElementById(t).style.display = if t is id then "" else "none"

# ── Fetch helpers ────────────────────────────────────────────────────────────
api = (path, opts = {}) ->
  mergedOpts = Object.assign(
    headers: { "Content-Type": "application/json" }
  , opts,
    body: if opts.body then JSON.stringify(opts.body) else undefined
  )
  res = await fetch API + path, mergedOpts
  res.json()

loadAll = ->
  [_kriteria, _alternatif, _penilaian, _hasil, _norm, _vsDetail] = await Promise.all [
    api "/api/kriteria"
    api "/api/alternatif"
    api "/api/penilaian"
    api "/api/hasil"
    api "/api/normalisasi"
    api "/api/vektor_s_detail"
  ]

# ── Toast ────────────────────────────────────────────────────────────────────
toast = (msg, type = "") ->
  el = document.getElementById "toast"
  el.textContent = msg
  el.className = "toast show#{if type then " " + type else ""}"
  setTimeout (-> el.className = "toast"), 2500

# ── Pages ────────────────────────────────────────────────────────────────────
loadPage = (id) ->
  await loadAll()
  switch id
    when "dashboard" then renderDashboard()
    when "kriteria" then renderKriteria()
    when "alternatif" then renderAlternatif()
    when "penilaian" then renderPenilaian()
    when "perhitungan" then renderPerhitungan()
    when "hasil" then renderHasil()
    else console.error 'Halaman yang dipilih tidak ditemukan'

# ── Dashboard ────────────────────────────────────────────────────────────────
renderDashboard = ->
  document.getElementById("stat-kriteria").textContent = _kriteria.length
  document.getElementById("stat-alternatif").textContent = _alternatif.length
  best = _hasil[0]

  if best
    document.getElementById("stat-terbaik").textContent = best.nama
    document.getElementById("stat-terbaik-v").textContent = "V = #{best.vektor_v.toFixed(5)}"

  maxV = if _hasil.length then _hasil[0].vektor_v else 1
  grid = document.getElementById "ranking-grid"
  grid.innerHTML = _hasil.map((h) ->
    pct = ((h.vektor_v / maxV) * 100).toFixed(1)
    sPct = ((h.vektor_s / (_hasil[0].vektor_s or 1)) * 100).toFixed(1)
    """
    <div class="rank-card rank-#{h.ranking}">
      <div class="rank-num">Peringkat #{h.ranking}</div>
      <div class="rank-name">#{h.nama}</div>
      <div class="rank-code">#{h.alternatif_kode}</div>
      <div class="rank-bars">
        <div class="rank-bar-row">
          <span class="rank-bar-label">Vektor V</span>
          <div class="rank-bar-track"><div class="rank-bar-fill" style="width:#{pct}%"></div></div>
          <span class="rank-bar-val">#{h.vektor_v.toFixed(5)}</span>
        </div>
        <div class="rank-bar-row">
          <span class="rank-bar-label">Vektor S</span>
          <div class="rank-bar-track"><div class="rank-bar-fill" style="width:#{sPct}%;background:#7c3aed"></div></div>
          <span class="rank-bar-val">#{h.vektor_s.toFixed(4)}</span>
        </div>
      </div>
    </div>
    """
  ).join ""

# ── Kriteria ─────────────────────────────────────────────────────────────────
renderKriteria = ->
  totalBobot = _kriteria.reduce ((s, k) -> s + k.bobot), 0
  badge = document.getElementById "total-bobot-badge"
  badge.textContent = "Total Bobot: #{totalBobot.toFixed(4)}"
  badge.className = "badge " + (if Math.abs(totalBobot - 1) < 0.001 then "badge-benefit" else "badge-cost")

  totalNorm = _norm.reduce ((s, n) -> s + n.bobot_norm), 0
  document.getElementById("kriteria-tbody").innerHTML = _kriteria.map((k) ->
    norm = _norm.find (n) -> n.kode is k.kode
    bNorm = if norm then norm.bobot_norm.toFixed(4) else "–"
    """
    <tr>
      <td><b>#{k.kode}</b></td>
      <td>#{k.nama}</td>
      <td><span class="badge badge-#{k.tipe.toLowerCase()}">#{k.tipe}</span></td>
      <td>#{k.bobot}</td>
      <td>#{bNorm}</td>
      <td>
        <button class="btn btn-sm btn-ghost" onclick="openEditKriteria('#{k.kode}')">Edit</button>
        <button class="btn btn-sm btn-danger" onclick="deleteKriteria('#{k.kode}')">Hapus</button>
      </td>
    </tr>
    """
  ).join ""

addKriteria = ->
  kode = document.getElementById("k-kode").value.trim()
  nama = document.getElementById("k-nama").value.trim()
  tipe = document.getElementById("k-tipe").value
  bobot = parseFloat document.getElementById("k-bobot").value

  return toast("Lengkapi semua field!", "error") if not kode or not nama or isNaN(bobot)

  await api "/api/kriteria",
    method: "POST"
    body: { kode, nama, tipe, bobot }

  toast "Kriteria ditambahkan"
  for id in ["k-kode", "k-nama", "k-bobot"]
    document.getElementById(id).value = ""
  loadPage "kriteria"

openEditKriteria = (kode) ->
  k = _kriteria.find (x) -> x.kode is kode
  document.getElementById("edit-k-kode").value = k.kode
  document.getElementById("edit-k-nama").value = k.nama
  document.getElementById("edit-k-tipe").value = k.tipe
  document.getElementById("edit-k-bobot").value = k.bobot
  document.getElementById("modal-kriteria").classList.add "open"

saveKriteria = ->
  kode = document.getElementById("edit-k-kode").value
  await api "/api/kriteria/#{kode}",
    method: "PUT"
    body:
      nama: document.getElementById("edit-k-nama").value
      tipe: document.getElementById("edit-k-tipe").value
      bobot: parseFloat document.getElementById("edit-k-bobot").value

  closeModal "modal-kriteria"
  toast "Kriteria diperbarui"
  loadPage "kriteria"

deleteKriteria = (kode) ->
  return unless confirm "Hapus kriteria #{kode}? Data penilaian terkait juga akan dihapus."
  await api "/api/kriteria/#{kode}", method: "DELETE"
  toast "Kriteria dihapus", "error"
  loadPage "kriteria"

# ── Alternatif ───────────────────────────────────────────────────────────────
renderAlternatif = ->
  document.getElementById("alternatif-tbody").innerHTML = _alternatif.map((a) ->
    """
    <tr>
      <td><b>#{a.kode}</b></td>
      <td>#{a.nama}</td>
      <td>
        <button class="btn btn-sm btn-ghost" onclick="openEditAlternatif('#{a.kode}')">Edit</button>
        <button class="btn btn-sm btn-danger" onclick="deleteAlternatif('#{a.kode}')">Hapus</button>
      </td>
    </tr>
    """
  ).join ""

addAlternatif = ->
  kode = document.getElementById("a-kode").value.trim()
  nama = document.getElementById("a-nama").value.trim()
  return toast("Lengkapi semua field!", "error") if not kode or not nama

  await api "/api/alternatif", method: "POST", body: { kode, nama }
  toast "Alternatif ditambahkan"
  for id in ["a-kode", "a-nama"]
    document.getElementById(id).value = ""
  loadPage "alternatif"

openEditAlternatif = (kode) ->
  a = _alternatif.find (x) -> x.kode is kode
  document.getElementById("edit-a-kode").value = a.kode
  document.getElementById("edit-a-nama").value = a.nama
  document.getElementById("modal-alternatif").classList.add "open"

saveAlternatif = ->
  kode = document.getElementById("edit-a-kode").value
  await api "/api/alternatif/#{kode}",
    method: "PUT"
    body:
      nama: document.getElementById("edit-a-nama").value

  closeModal "modal-alternatif"
  toast "Alternatif diperbarui"
  loadPage "alternatif"

deleteAlternatif = (kode) ->
  return unless confirm "Hapus alternatif #{kode}?"
  await api "/api/alternatif/#{kode}", method: "DELETE"
  toast "Alternatif dihapus", "error"
  loadPage "alternatif"

# ── Penilaian ─────────────────────────────────────────────────────────────────
renderPenilaian = ->
  head = document.getElementById "penilaian-head"
  body = document.getElementById "penilaian-body"

  nilaiMap = {}
  for p in _penilaian
    nilaiMap["#{p.alternatif_kode}|#{p.kriteria_kode}"] = p.nilai

  head.innerHTML = "<tr><th>Kode</th><th>Nama Penyedia</th>" +
    _kriteria.map((k) -> """
      <th title="#{k.nama}">#{k.kode}<br><small style="color:var(--muted);font-size:10px">#{k.tipe}</small></th>
    """).join("") + "</tr>"

  body.innerHTML = _alternatif.map((a) ->
    cells = _kriteria.map((k) ->
      val = nilaiMap["#{a.kode}|#{k.kode}"] or ""
      """
      <td><input class="nilai-input" type="number" min="1" max="10" step="1"
        data-alt="#{a.kode}" data-krit="#{k.kode}" value="#{val}" placeholder="–"/></td>
      """
    ).join ""
    "<tr><td><b>#{a.kode}</b></td><td>#{a.nama}</td>#{cells}</tr>"
  ).join ""

savePenilaian = ->
  inputs = document.querySelectorAll "#penilaian-body .nilai-input"
  items = []

  for inp in inputs
    val = parseFloat inp.value
    unless isNaN val
      items.push
        alternatif_kode: inp.dataset.alt
        kriteria_kode: inp.dataset.krit
        nilai: val

  await api "/api/penilaian", method: "POST", body: items
  toast "Penilaian tersimpan & dihitung ulang"
  loadPage "penilaian"

# ── Perhitungan ───────────────────────────────────────────────────────────────
renderPerhitungan = ->
  # Normalisasi
  tbody = document.getElementById "norm-tbody"
  totalB = _norm.reduce ((s, n) -> s + n.bobot), 0
  totalN = _norm.reduce ((s, n) -> s + n.bobot_norm), 0

  tbody.innerHTML = _norm.map((n) ->
    """
    <tr>
      <td><b>#{n.kode}</b></td><td>#{n.nama}</td>
      <td><span class="badge badge-#{n.tipe.toLowerCase()}">#{n.tipe}</span></td>
      <td>#{n.bobot}</td>
      <td><b>#{n.bobot_norm.toFixed(4)}</b></td>
    </tr>
    """
  ).join("") +
  """
    <tr style="font-weight:700;background:var(--surface)">
      <td colspan="3">TOTAL</td><td>#{totalB.toFixed(4)}</td><td>#{totalN.toFixed(4)}</td>
    </tr>
  """

  # Vektor S header + rows
  vsHead = document.getElementById "vektors-head"
  vsTbody = document.getElementById "vektors-tbody"

  detailMap = {}
  for d in _vsDetail
    detailMap["#{d.alternatif_kode}|#{d.kriteria_kode}"] = d

  # Pangkat row
  vsHead.innerHTML = "<tr><th>Kode</th><th>Nama</th>" +
    _kriteria.map((k) ->
      norm = _norm.find (n) -> n.kode is k.kode
      pangkat = if norm then (if k.tipe is "Cost" then -norm.bobot_norm else norm.bobot_norm).toFixed(4) else "–"
      "<th>#{k.kode}<br><small style=\"color:var(--accent);font-size:10px\">(#{pangkat})</small></th>"
    ).join("") + "<th>Nilai S</th></tr>"

  vsTbody.innerHTML = _hasil.map((h) ->
    cells = _kriteria.map((k) ->
      d = detailMap["#{h.alternatif_kode}|#{k.kode}"]
      "<td style=\"font-size:12px\">#{if d then d.komponen.toFixed(6) else "–"}</td>"
    ).join ""
    "<tr><td><b>#{h.alternatif_kode}</b></td><td>#{h.nama}</td>#{cells}<td><b>#{h.vektor_s.toFixed(6)}</b></td></tr>"
  ).join ""

  # Vektor V
  maxV = if _hasil.length then _hasil[0].vektor_v else 1
  document.getElementById("vektorv-tbody").innerHTML = _hasil.map((h) ->
    pct = ((h.vektor_v / maxV) * 100).toFixed(1)
    """
    <tr>
      <td><b>##{h.ranking}</b></td>
      <td>#{h.alternatif_kode}</td>
      <td>#{h.nama}</td>
      <td>#{h.vektor_s.toFixed(6)}</td>
      <td><b>#{h.vektor_v.toFixed(6)}</b></td>
      <td><div class="progress-bar"><div class="progress-fill" style="width:#{pct}%"></div></div>#{pct}%</td>
    </tr>
    """
  ).join ""

# ── Hasil ─────────────────────────────────────────────────────────────────────
renderHasil = ->
  best = _hasil[0]
  if best
    document.getElementById("rekomendasi-content").innerHTML = """
    <div style="display:flex;align-items:center;gap:20px;flex-wrap:wrap;">
      <div>
        <div style="font-size:12px;color:var(--muted);text-transform:uppercase;letter-spacing:.5px;">Rekomendasi Terbaik</div>
        <div style="font-size:22px;font-weight:700;color:var(--rank1)">#{best.nama}</div>
        <div style="font-size:13px;color:var(--muted);margin-top:4px;">
          Kode: <b>#{best.alternatif_kode}</b> &nbsp;|&nbsp;
          Nilai V: <b style="color:var(--green)">#{best.vektor_v.toFixed(6)}</b> &nbsp;|&nbsp;
          Nilai S: <b>#{best.vektor_s.toFixed(6)}</b>
        </div>
        <div style="font-size:12px;color:var(--muted);margin-top:8px;">
          Penyedia ini memiliki nilai Vektor V tertinggi berdasarkan metode Weighted Product,
          dengan mempertimbangkan #{_kriteria.length} kriteria dan #{_alternatif.length} alternatif.
        </div>
      </div>
    </div>
    """

  maxV = if _hasil.length then _hasil[0].vektor_v else 1
  document.getElementById("hasil-tbody").innerHTML = _hasil.map((h) ->
    pct = ((h.vektor_v / maxV) * 100).toFixed(1)
    rowStyle = if h.ranking is 1 then "background:rgba(251,191,36,.05)" else ""
    """
    <tr style="#{rowStyle}">
      <td><b>##{h.ranking}</b></td>
      <td><b>#{h.alternatif_kode}</b></td>
      <td>#{h.nama}</td>
      <td>#{h.vektor_s.toFixed(6)}</td>
      <td><b style="color:#{if h.ranking is 1 then 'var(--rank1)' else 'var(--text)'}">#{h.vektor_v.toFixed(6)}</b></td>
      <td>
        <div class="progress-bar"><div class="progress-fill" style="width:#{pct}%"></div></div>
        #{pct}%
      </td>
    </tr>
    """
  ).join ""

# ── Modal ─────────────────────────────────────────────────────────────────────
closeModal = (id) ->
  document.getElementById(id).classList.remove "open"

for m in document.querySelectorAll ".modal-overlay"
  m.addEventListener "click", (e) ->
    m.classList.remove "open" if e.target is m

# ── Init ──────────────────────────────────────────────────────────────────────
for x in document.querySelectorAll '.nav-item[data-page]'
  x.addEventListener 'click', (e) ->
    showPage e, @dataset.page
    
for x in document.querySelectorAll '.tab[data-tab]'
  x.addEventListener 'click', (e) ->
    showTab e, @dataset.tab

loadPage "dashboard"
