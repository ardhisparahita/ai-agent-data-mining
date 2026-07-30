"""
Support Agent — menjawab pertanyaan pelanggan menggunakan RAG.
"""
from langchain_core.tools import tool
from rag.retriever import answer_with_rag


@tool
def support_rag_tool(question: str) -> str:
    """Gunakan untuk menjawab pertanyaan pelanggan seputar FAQ, kebijakan retur, dan pengiriman menggunakan RAG."""
    result = answer_with_rag(question, collection_name="faq")
    return result["result"]