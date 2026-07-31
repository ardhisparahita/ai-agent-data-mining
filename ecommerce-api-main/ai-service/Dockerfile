# Dockerfile untuk ai-service (Python/FastAPI)
# Taruh file ini di dalam folder ai-service/ (sejajar dengan main.py, requirements.txt)

FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8000

CMD ["sh", "-c", "python rag/ingest.py && uvicorn main:app --host 0.0.0.0 --port ${PORT:-8000}"]
