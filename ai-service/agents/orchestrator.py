"""
Multi-Agent Orchestrator — LLM memilih agent (tool) yang relevan lalu memanggilnya.
Menggunakan loop agar bisa menangani kasus LLM memanggil lebih dari satu tool
secara bertahap sebelum memberi jawaban akhir.
"""
import os
from dotenv import load_dotenv
load_dotenv()

from langchain_groq import ChatGroq
from langchain_core.messages import ToolMessage

from agents.sales_agent import sales_data_tool
from agents.inventory_agent import inventory_tool
from agents.support_agent import support_rag_tool

tools = [sales_data_tool, inventory_tool, support_rag_tool]
tool_map = {t.name: t for t in tools}

llm = ChatGroq(
    model=os.getenv("GROQ_MODEL", "llama-3.3-70b-versatile"),
    temperature=0,
    api_key=os.getenv("GROQ_API_KEY"),
    timeout=30,       # cegah request menggantung lama
    max_retries=1,    # kurangi retry berlebihan yang bikin lama
)
llm_with_tools = llm.bind_tools(tools)

MAX_ITERATIONS = 4


def handle_query(user_query: str) -> dict:
    messages = [("human", user_query)]
    agents_called = []

    for _ in range(MAX_ITERATIONS):
        ai_msg = llm_with_tools.invoke(messages)

        if not ai_msg.tool_calls:
            # LLM sudah selesai — ini jawaban final
            answer = ai_msg.content.strip() if ai_msg.content else ""
            if not answer and agents_called:
                # fallback: kalau LLM tidak memberi teks tapi tool sudah menghasilkan data
                answer = agents_called[-1]["output"]
            return {"answer": answer or "Maaf, tidak ada jawaban yang dihasilkan.", "agents_called": agents_called}

        messages.append(ai_msg)

        for tool_call in ai_msg.tool_calls:
            selected_tool = tool_map[tool_call["name"]]
            tool_output = selected_tool.invoke(tool_call["args"])
            print(f"[orchestrator] Agent dipanggil: {tool_call['name']} -> {tool_output}")

            agents_called.append({
                "tool_name": tool_call["name"],
                "args": tool_call["args"],
                "output": str(tool_output),
            })
            messages.append(ToolMessage(content=str(tool_output), tool_call_id=tool_call["id"]))

    return {
        "answer": "Maaf, proses memakan terlalu banyak langkah. Coba pertanyaan yang lebih spesifik.",
        "agents_called": agents_called,
    }


if __name__ == "__main__":
    result = handle_query("Produk apa yang stoknya menipis minggu ini?")
    print(result)