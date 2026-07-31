# 🏢 Studi Kasus Enterprise — Ecommerce API Multi-Agent System

## 1. Latar Belakang Masalah

`ecommerce-api` saat ini melayani seluruh proses bisnis inti sebuah perusahaan e-commerce: autentikasi pengguna, manajemen produk, keranjang belanja, checkout, dan pengelolaan pesanan/pembayaran. Namun, seluruh proses ini masih ditangani secara **manual atau melalui panggilan API langsung** oleh masing-masing tim internal.

Dalam operasional nyata, sebuah perusahaan e-commerce terbagi ke dalam beberapa **divisi** yang masing-masing menghadapi masalah spesifik yang sulit diselesaikan hanya dengan CRUD API biasa:

| Divisi | Masalah Bisnis | Kebutuhan Data | Solusi yang Dibutuhkan |
|---|---|---|---|
| **Sales & Marketing** | Tim sales kesulitan merangkum performa penjualan harian/mingguan dan mengidentifikasi tren produk terlaris tanpa menulis query manual | Data `orders`, `order_items`, `payments` | Agent yang bisa menjawab pertanyaan natural language seperti "produk apa yang paling laku minggu ini?" |
| **Inventory / Gudang** | Staf gudang perlu memantau stok produk secara cepat dan mendapat rekomendasi kapan harus restock | Data `products`, `categories` | Agent yang mengecek stok real-time dan memberi peringatan/prediksi kebutuhan restock |
| **Customer Support (CS)** | Tim CS menerima banyak pertanyaan berulang (kebijakan retur, status pengiriman, cara pembayaran) yang memakan waktu jika dijawab manual satu per satu | Dokumen FAQ, kebijakan retur, data `orders` milik pelanggan | Agent berbasis RAG yang menjawab pertanyaan pelanggan secara otomatis dan akurat berdasarkan basis pengetahuan resmi perusahaan |
| **Finance / Payment** | Tim finance perlu rekap status pembayaran (paid, failed, pending) untuk rekonsiliasi harian | Data `payments`, `orders` | Agent yang merangkum status transaksi dan mendeteksi anomali pembayaran |

## 2. Solusi yang Diusulkan

Untuk menjawab kebutuhan lintas-divisi di atas, dibangun **layer AI tambahan** (`ai-service/`) yang berjalan berdampingan dengan `ecommerce-api` (Go/Fiber) yang sudah ada. Layer ini terdiri dari **beberapa agent AI otonom**, masing-masing mewakili satu divisi, yang berkoordinasi melalui satu *orchestrator* untuk menjawab pertanyaan pengguna (baik dari internal staff maupun pelanggan).

```
┌─────────────────────────────┐
│   User / Staff Query        │
│ ("Produk apa yang stoknya   │
│   menipis minggu ini?")     │
└─────────────┬───────────────┘
              ▼
┌─────────────────────────────┐
│      Agent Orchestrator     │
│   (LangChain Multi-Agent)   │
└──────┬──────────┬───────────┘
       │          │          │
       ▼          ▼          ▼
 ┌───────────┐┌───────────┐┌────────────────┐
 │Sales Agent││Inventory  ││Support Agent    │
 │           ││Agent      ││(RAG + VectorDB) │
 └─────┬─────┘└─────┬─────┘└────────┬────────┘
       │            │               │
       ▼            ▼               ▼
┌─────────────────────────────────────────────┐
│      ecommerce-api (Go/Fiber REST API)      │
│   /orders  /products  /payments  /carts     │
└─────────────────────────────────────────────┘
```

## 3. Pemetaan Divisi → Agent

| Divisi Enterprise | Nama Agent | Tools yang Digunakan | File Implementasi |
|---|---|---|---|
| Sales & Marketing | `SalesAgent` | Query REST API `/api/v1/orders` | `ai-service/agents/sales_agent.py` |
| Inventory / Gudang | `InventoryAgent` | Query REST API `/api/v1/products` | `ai-service/agents/inventory_agent.py` |
| Customer Support | `SupportAgent` | RAG (ChromaDB + retriever) atas dokumen FAQ/kebijakan | `ai-service/agents/support_agent.py` |
| Koordinasi Lintas-Divisi | `Orchestrator` | Routing pertanyaan ke agent yang relevan | `ai-service/agents/orchestrator.py` |

## 4. Manfaat bagi Enterprise

- **Sales**: Mempercepat pengambilan keputusan berbasis data tanpa perlu menulis query SQL manual.
- **Inventory**: Mengurangi risiko stockout melalui deteksi dini dan rekomendasi restock.
- **Customer Support**: Menurunkan beban tiket berulang dengan jawaban otomatis yang konsisten dan berbasis dokumen resmi (mengurangi risiko informasi yang salah/hallucination karena dijawab berbasis RAG, bukan model murni).
- **Finance**: Rekonsiliasi pembayaran lebih cepat dengan ringkasan otomatis.

Dengan pendekatan ini, `ecommerce-api` tidak lagi sekadar backend CRUD, tetapi menjadi **fondasi data** bagi sistem multi-agent AI yang menyelesaikan masalah nyata di berbagai divisi perusahaan e-commerce.
