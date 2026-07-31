const AI_SERVICE_URL = import.meta.env.VITE_AI_SERVICE_URL || 'http://localhost:8000'

export async function checkHealth() {
  try {
    const res = await fetch(`${AI_SERVICE_URL}/health`, { signal: AbortSignal.timeout(4000) })
    return res.ok
  } catch {
    return false
  }
}

export async function sendQuery(query) {
  const res = await fetch(`${AI_SERVICE_URL}/query`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query }),
  })
  if (!res.ok) {
    throw new Error(`HTTP ${res.status}`)
  }
  return res.json()
}
