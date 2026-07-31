"""
Embedding + Vector DB ingestion pipeline.
Menggunakan embedding lokal gratis (HuggingFace sentence-transformers) — tidak butuh API key.
"""
import os
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_huggingface import HuggingFaceEmbeddings
from langchain_community.vectorstores import Chroma
from dotenv import load_dotenv

load_dotenv()

PERSIST_DIR = os.getenv("CHROMA_PERSIST_DIR", "./chroma_db")
EMBEDDING_MODEL = os.getenv("EMBEDDING_MODEL", "sentence-transformers/all-MiniLM-L6-v2")


def get_embeddings():
    return HuggingFaceEmbeddings(model_name=EMBEDDING_MODEL)


def ingest_product_catalog(products: list[dict]):
    docs = [
        f"{p['name']} - {p.get('description', '')} - Kategori: {p.get('category', '')} "
        f"- Stok: {p.get('stock', 'tidak diketahui')}"
        for p in products
    ]

    splitter = RecursiveCharacterTextSplitter(chunk_size=500, chunk_overlap=50)
    chunks = splitter.create_documents(docs)

    vectorstore = Chroma.from_documents(
        chunks, get_embeddings(), persist_directory=PERSIST_DIR, collection_name="products"
    )
    vectorstore.persist()
    print(f"[ingest] {len(chunks)} chunk produk berhasil di-embed ke {PERSIST_DIR}")
    return vectorstore


def ingest_faq_documents(faq_texts: list[str]):
    splitter = RecursiveCharacterTextSplitter(chunk_size=500, chunk_overlap=50)
    chunks = splitter.create_documents(faq_texts)

    vectorstore = Chroma.from_documents(
        chunks, get_embeddings(), persist_directory=PERSIST_DIR, collection_name="faq"
    )
    vectorstore.persist()
    print(f"[ingest] {len(chunks)} chunk FAQ berhasil di-embed ke {PERSIST_DIR}")
    return vectorstore


if __name__ == "__main__":
    sample_products = [
        {"name": "Kaos Polos Hitam", "description": "Kaos cotton combed 30s", "category": "Fashion", "stock": 12},
        {"name": "Sepatu Lari X1", "description": "Sepatu running ringan", "category": "Olahraga", "stock": 3},
    ]
    sample_faq = [
        "Kebijakan retur: produk dapat dikembalikan dalam 7 hari sejak diterima, dengan kondisi masih tersegel.",
        "Estimasi pengiriman: 2-4 hari kerja untuk area Jawa, 5-7 hari untuk luar Jawa.",
    ]

    ingest_product_catalog(sample_products)
    ingest_faq_documents(sample_faq)