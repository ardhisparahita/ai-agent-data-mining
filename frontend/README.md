# Dispatch Console — Frontend

Frontend React (Vite) untuk berinteraksi dengan sistem multi-agent `ai-service` yang
mengoordinasikan Sales Agent, Inventory Agent, dan Support Agent (RAG) di atas
`ecommerce-api`.

## Fitur

- **Agent Roster** — menampilkan status real-time (idle/aktif/selesai) tiap agent
  berdasarkan agent yang benar-benar dipanggil orchestrator untuk setiap query.
- **Console Chat** — antarmuka percakapan untuk mengajukan pertanyaan.
- **Route Trace** — log timeline yang menampilkan alur routing query secara transparan.

## Menjalankan secara lokal

```bash
npm install
cp .env.example .env
```

Edit `.env` dan sesuaikan URL `ai-service` Anda:
```
VITE_AI_SERVICE_URL=http://localhost:8000
```

Jalankan:
```bash
npm run dev
```
Buka `http://localhost:5173`.

> **Penting**: pastikan `ai-service` (backend Python) sudah berjalan dan CORS
> sudah diaktifkan (lihat `ai-service/main.py` — sudah termasuk `CORSMiddleware`).

## Build untuk production

```bash
npm run build
```
Hasil build ada di folder `dist/`, siap di-deploy ke hosting statis mana pun.

## Deploy ke Vercel

1. Push project ini (folder `frontend/`) ke repository GitHub.
2. Buka [vercel.com](https://vercel.com) → **New Project** → import repo Anda.
3. Set **Root Directory** ke `frontend` (jika project ini ada di dalam repo `ecommerce-api`).
4. Tambahkan Environment Variable:
   - `VITE_AI_SERVICE_URL` = URL `ai-service` Anda yang sudah di-deploy (misal Railway/Render — lihat bagian di bawah).
5. Klik **Deploy**.

Konfigurasi `vercel.json` sudah disertakan, jadi build command dan output directory otomatis terdeteksi.

## Deploy ke Netlify

1. Push project ke GitHub.
2. Buka [netlify.com](https://netlify.com) → **Add new site** → **Import an existing project**.
3. Set **Base directory** ke `frontend`.
4. Build command: `npm run build`, Publish directory: `dist` (sudah diatur di `netlify.toml`).
5. Tambahkan Environment Variable `VITE_AI_SERVICE_URL` di **Site settings → Environment variables**.

## Men-deploy `ai-service` (backend Python) agar bisa diakses frontend yang sudah live

Frontend yang di-deploy (Vercel/Netlify) **tidak bisa** memanggil `http://localhost:8000` —
`ai-service` juga harus di-deploy ke server yang bisa diakses publik. Opsi gratis/murah:

- **Railway** (https://railway.app) — deploy langsung dari GitHub, mendukung Python/FastAPI, ada free tier.
- **Render** (https://render.com) — free tier untuk web service Python.

Langkah umum di kedua platform tersebut:
1. Hubungkan repo GitHub Anda, pilih folder `ai-service/`.
2. Set start command: `uvicorn main:app --host 0.0.0.0 --port $PORT`.
3. Tambahkan Environment Variables: `GROQ_API_KEY`, `ECOMMERCE_API_BASE_URL`, `ECOMMERCE_API_TOKEN`, `ALLOWED_ORIGINS` (isi dengan URL frontend Vercel/Netlify Anda, misal `https://dispatch-console.vercel.app`).
4. Setelah deploy selesai, copy URL publik yang diberikan (misal `https://ai-service-production.up.railway.app`), lalu gunakan itu sebagai nilai `VITE_AI_SERVICE_URL` di pengaturan environment variable frontend Anda.

## Struktur Project

```
frontend/
├── src/
│   ├── components/
│   │   ├── AgentRoster.jsx   # Panel status agent (kiri)
│   │   ├── ChatConsole.jsx   # Panel percakapan (tengah)
│   │   └── TraceLog.jsx      # Panel log routing (kanan)
│   ├── api.js                # Client untuk memanggil ai-service
│   ├── App.jsx                # Komponen utama + state management
│   ├── main.jsx                # Entry point React
│   └── styles.css             # Design tokens & styling
├── index.html
├── package.json
├── vite.config.js
├── vercel.json
├── netlify.toml
└── .env.example
```
