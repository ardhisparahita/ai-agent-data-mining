"""
RAG retrieval chain — implementasi manual (retrieve -> augment -> generate)
menggunakan Groq (gratis) sebagai LLM dan HuggingFace sebagai embedding.
Tidak menggunakan RetrievalQA (deprecated) untuk menghindari konflik versi langchain.
"""
import os
from langchain_groq import ChatGroq
from langchain_huggingface import HuggingFaceEmbeddings
from langchain_community.vectorstores import Chroma
from dotenv import load_dotenv

load_dotenv()

PERSIST_DIR = os.getenv("CHROMA_PERSIST_DIR", "./chroma_db")
EMBEDDING_MODEL = os.getenv("EMBEDDING_MODEL", "sentence-transformers/all-MiniLM-L6-v2")


def get_retriever(collection_name: str = "faq", k: int = 4):
    embeddings = HuggingFaceEmbeddings(model_name=EMBEDDING_MODEL)
    vectorstore = Chroma(
        persist_directory=PERSIST_DIR,
        embedding_function=embeddings,
        collection_name=collection_name,
    )
    return vectorstore.as_retriever(search_kwargs={"k": k})


def get_llm():
    return ChatGroq(
        model=os.getenv("GROQ_MODEL", "llama-3.3-70b-versatile"),
        temperature=0,
        api_key=os.getenv("GROQ_API_KEY"),
    )


def answer_with_rag(question: str, collection_name: str = "faq") -> dict:
    """Alur RAG eksplisit: Retrieve -> Augment -> Generate."""
    # 1. RETRIEVE — ambil dokumen relevan dari vector DB
    retriever = get_retriever(collection_name)
    docs = retriever.invoke(question)
    context = "\n\n".join([doc.page_content for doc in docs])

    # 2. AUGMENT — susun prompt dengan konteks yang di-retrieve
    prompt = f"""Jawab pertanyaan berikut HANYA berdasarkan konteks di bawah ini.
Jika jawabannya tidak ada di konteks, katakan tidak tahu.

Konteks:
{context}

Pertanyaan: {question}

Jawaban:"""

    # 3. GENERATE — panggil LLM untuk menghasilkan jawaban
    llm = get_llm()
    response = llm.invoke(prompt)

    return {
        "result": response.content,
        "source_documents": docs,
    }


def build_rag_chain(collection_name: str = "faq"):
    """Wrapper agar kompatibel dengan pemanggilan lama: chain({"query": ...})"""
    def chain(inputs: dict):
        return answer_with_rag(inputs["query"], collection_name)
    return chain


if __name__ == "__main__":
    result = answer_with_rag("Berapa lama saya bisa retur barang?")
    print("Jawaban:", result["result"])
    print("Sumber:", [doc.page_content[:80] for doc in result["source_documents"]])