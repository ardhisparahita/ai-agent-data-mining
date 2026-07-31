"""
Inventory Agent — mengecek stok & detail produk, termasuk deteksi stok menipis.
Selalu mengambil SEMUA produk dari API, lalu filter dilakukan di sisi Python
(tidak bergantung pada parameter search bawaan backend) agar hasil lebih konsisten.
"""
import os
import requests
from langchain_core.tools import tool

BASE_URL = os.getenv("ECOMMERCE_API_BASE_URL", "http://localhost:3000/api/v1")
TOKEN = os.getenv("ECOMMERCE_API_TOKEN", "")
LOW_STOCK_THRESHOLD = 5


def _fetch_all_products() -> list[dict]:
    resp = requests.get(
        f"{BASE_URL}/products",
        headers={"Authorization": f"Bearer {TOKEN}"},
        timeout=10,
    )
    resp.raise_for_status()
    data = resp.json()
    return data.get("data", {}).get("items", [])


@tool
def inventory_tool(intent: str = "all") -> str:
    """Gunakan untuk pertanyaan seputar stok atau daftar produk.
    Isi 'intent' dengan salah satu dari: 'low_stock' (untuk produk yang stoknya menipis/hampir habis),
    'all' (untuk menampilkan semua produk), atau ketikkan nama/kata kunci produk tertentu
    untuk mencari produk spesifik."""
    try:
        items = _fetch_all_products()
        if not items:
            return "Tidak ada data produk sama sekali di database."

        intent_lower = intent.lower().strip()

        if intent_lower in ("low_stock", "stok menipis", "menipis", ""):
            filtered = [p for p in items if p.get("stock", 999) < LOW_STOCK_THRESHOLD]
            if not filtered:
                return f"Tidak ada produk dengan stok di bawah {LOW_STOCK_THRESHOLD}."
            return f"Produk dengan stok menipis (<{LOW_STOCK_THRESHOLD}): {filtered}"

        if intent_lower == "all":
            return f"Semua produk: {items}"

        # Cari berdasarkan nama produk
        matched = [p for p in items if intent_lower in p.get("name", "").lower()]
        if not matched:
            return f"Produk dengan nama mengandung '{intent}' tidak ditemukan."
        return f"Produk ditemukan: {matched}"

    except Exception as e:
        return f"Gagal mengambil data inventory: {e}"