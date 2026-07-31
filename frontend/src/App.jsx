import { useEffect, useState, useCallback } from 'react'
import AgentRoster from './components/AgentRoster'
import ChatConsole from './components/ChatConsole'
import TraceLog from './components/TraceLog'
import { checkHealth, sendQuery } from './api'

function timeNow() {
  return new Date().toTimeString().slice(0, 8)
}

export default function App() {
  const [isConnected, setIsConnected] = useState(null) // null = checking, true/false = result
  const [messages, setMessages] = useState([])
  const [traceEntries, setTraceEntries] = useState([])
  const [agentStates, setAgentStates] = useState({})
  const [isLoading, setIsLoading] = useState(false)

  const addTrace = useCallback((text, type = '') => {
    setTraceEntries((prev) => [...prev, { time: timeNow(), text, type }])
  }, [])

  const pollHealth = useCallback(async () => {
    const ok = await checkHealth()
    setIsConnected(ok)
  }, [])

  useEffect(() => {
    pollHealth()
    const interval = setInterval(pollHealth, 15000)
    return () => clearInterval(interval)
  }, [pollHealth])

  const resetAgentLights = () => setAgentStates({})

  const handleSend = async (question) => {
    setMessages((prev) => [...prev, { role: 'user', text: question }])
    addTrace(`Query diterima: "${question}"`)
    setIsLoading(true)
    resetAgentLights()

    try {
      addTrace('Orchestrator menganalisis intent...', 'step-route')
      const data = await sendQuery(question)

      // Tandai agent yang benar-benar dipanggil berdasarkan response backend
      const nextStates = {}
      ;(data.agents_called || []).forEach((call) => {
        nextStates[call.tool_name] = 'active'
        addTrace(`Agent dipanggil: ${call.tool_name}`, 'step-route')
      })
      setAgentStates(nextStates)

      // Setelah sedikit jeda, tandai selesai (hijau)
      setTimeout(() => {
        const doneStates = {}
        Object.keys(nextStates).forEach((k) => (doneStates[k] = 'done'))
        setAgentStates(doneStates)
      }, 400)

     setMessages((prev) => [...prev, { role: 'agent', text: data.answer || '(tidak ada jawaban)' }])
      addTrace('Jawaban terkirim ke console.', 'step-done')
    } catch (err) {
      setMessages((prev) => [
        ...prev,
        {
          role: 'agent',
          isError: true,
          text: `Gagal terhubung ke ai-service: ${err.message}. Pastikan backend berjalan dan URL VITE_AI_SERVICE_URL sudah benar.`,
        },
      ])
      addTrace(`Error: ${err.message}`, 'step-error')
    } finally {
      setIsLoading(false)
      setTimeout(resetAgentLights, 2500)
    }
  }

  return (
    <div className="app">
      <div className="topbar">
        <div className="brand">
          <span className="brand-mark">Dispatch</span>
          <span className="brand-sub">ECOMMERCE AI · MULTI-AGENT CONSOLE</span>
        </div>
        <div className="conn-status">
          <span
            className={`conn-dot ${isConnected === null ? '' : isConnected ? 'live' : 'down'}`}
          ></span>
          <span>
            {isConnected === null
              ? 'memeriksa koneksi...'
              : isConnected
              ? 'ai-service tersambung'
              : 'ai-service tidak terjangkau'}
          </span>
        </div>
      </div>

      <div className="layout">
        <AgentRoster agentStates={agentStates} />
        <ChatConsole messages={messages} onSend={handleSend} isLoading={isLoading} />
        <TraceLog entries={traceEntries} />
      </div>
    </div>
  )
}
