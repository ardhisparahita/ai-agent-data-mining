"""
Sales Agent — menjawab pertanyaan terkait penjualan & order.
"""
import os
import requests
from langchain_core.tools import tool

BASE_URL = os.getenv("ECOMMERCE_API_BASE_URL", "http://localhost:3000/api/v1")
TOKEN = os.getenv("ECOMMERCE_API_TOKEN", "")


@tool
def sales_data_tool(query: str) -> str:
    """Gunakan untuk menjawab pertanyaan terkait penjualan, order, dan revenue."""
    try:
        resp = requests.get(
            f"{BASE_URL}/orders",
            headers={"Authorization": f"Bearer {TOKEN}"},
            timeout=10,
        )
        resp.raise_for_status()
        return f"Data order terbaru: {resp.json()}"
    except Exception as e:
        return f"Gagal mengambil data sales: {e}"