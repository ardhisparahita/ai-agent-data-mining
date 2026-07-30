"""
Multi-Agent Orchestrator — LLM memilih agent (tool) yang relevan lalu memanggilnya.
"""
import os
from dotenv import load_dotenv
load_dotenv()   # <-- HARUS di atas, sebelum import agent-agent di bawah ini

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
)
llm_with_tools = llm.bind_tools(tools)


def handle_query(user_query: str) -> str:
    messages = [("human", user_query)]
    ai_msg = llm_with_tools.invoke(messages)

    if not ai_msg.tool_calls:
        return ai_msg.content

    messages.append(ai_msg)

    for tool_call in ai_msg.tool_calls:
        selected_tool = tool_map[tool_call["name"]]
        tool_output = selected_tool.invoke(tool_call["args"])
        print(f"[orchestrator] Agent dipanggil: {tool_call['name']} -> {tool_output}")
        messages.append(ToolMessage(content=str(tool_output), tool_call_id=tool_call["id"]))

    final_response = llm_with_tools.invoke(messages)
    return final_response.content


if __name__ == "__main__":
    print(handle_query("Berapa lama saya bisa retur barang?"))     # test SupportAgent (RAG)
    print(handle_query("Ada berapa order yang masuk?"))            # test SalesAgent