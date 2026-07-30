"""
ai-service — AI Multi-Agent Layer untuk ecommerce-api
"""
import os
from fastapi import FastAPI
from pydantic import BaseModel
from dotenv import load_dotenv

load_dotenv()

app = FastAPI(
    title="Ecommerce AI Multi-Agent Service",
    description="Layer AI (RAG + Multi-Agent) yang berdampingan dengan ecommerce-api",
    version="0.1.0",
)


class QueryRequest(BaseModel):
    query: str


class QueryResponse(BaseModel):
    query: str
    answer: str


@app.get("/health")
def health_check():
    return {"status": "ok", "service": "ai-service"}


@app.post("/query", response_model=QueryResponse)
def query_agent(request: QueryRequest):
    from agents.orchestrator import handle_query
    answer = handle_query(request.query)
    return QueryResponse(query=request.query, answer=answer)


if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("AI_SERVICE_PORT", 8000))
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)