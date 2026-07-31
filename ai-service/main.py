"""
ai-service — AI Multi-Agent Layer untuk ecommerce-api
Menyediakan endpoint HTTP agar frontend (React) dan orchestrator multi-agent
dapat berkomunikasi, dengan CORS aktif untuk akses dari domain frontend yang di-deploy.
"""

import os
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from dotenv import load_dotenv

load_dotenv()

app = FastAPI(
    title="Ecommerce AI Multi-Agent Service",
    description="Layer AI (RAG + Multi-Agent) yang berdampingan dengan ecommerce-api",
    version="0.2.0",
)

# CORS — izinkan frontend (dev maupun hasil deploy) memanggil service ini.
# Untuk production, ganti allow_origins dengan domain frontend spesifik Anda,
# misal: ["https://dispatch-console.vercel.app"]
ALLOWED_ORIGINS = os.getenv("ALLOWED_ORIGINS", "*").split(",")

app.add_middleware(
    CORSMiddleware,
    allow_origins=ALLOWED_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class QueryRequest(BaseModel):
    query: str


class AgentCall(BaseModel):
    tool_name: str
    args: dict
    output: str


class QueryResponse(BaseModel):
    query: str
    answer: str
    agents_called: list[AgentCall]


@app.get("/health")
def health_check():
    return {"status": "ok", "service": "ai-service"}


@app.post("/query", response_model=QueryResponse)
def query_agent(request: QueryRequest):
    from agents.orchestrator import handle_query

    result = handle_query(request.query)
    return QueryResponse(
        query=request.query,
        answer=result["answer"],
        agents_called=result["agents_called"],
    )


if __name__ == "__main__":
    import uvicorn

    port = int(os.getenv("AI_SERVICE_PORT", 8000))
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)
