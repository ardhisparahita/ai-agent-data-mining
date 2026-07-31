import { useEffect, useRef, useState } from 'react'

export default function ChatConsole({ messages, onSend, isLoading }) {
  const [input, setInput] = useState('')
  const scrollRef = useRef(null)

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [messages, isLoading])

  const handleSend = () => {
    const trimmed = input.trim()
    if (!trimmed || isLoading) return
    onSend(trimmed)
    setInput('')
  }

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') handleSend()
  }

  return (
    <div className="console">
      <div className="chat-scroll" ref={scrollRef}>
        {messages.length === 0 && (
          <div className="empty-state">
            <h2>Ajukan pertanyaan ke sistem</h2>
            <p>
              Orchestrator akan menentukan agent mana yang menangani — coba tanyakan
              soal stok, penjualan, atau kebijakan retur.
            </p>
          </div>
        )}

        {messages.map((msg, i) => (
          <div key={i} className={`msg ${msg.role}${msg.isError ? ' error' : ''}`}>
            <div className="msg-label">{msg.role === 'user' ? 'Anda' : 'Agent System'}</div>
            <div>{msg.text}</div>
          </div>
        ))}

        {isLoading && (
          <div className="typing">
            <span></span>
            <span></span>
            <span></span>
          </div>
        )}
      </div>

      <div className="input-row">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Contoh: produk apa yang stoknya menipis minggu ini?"
        />
        <button onClick={handleSend} disabled={isLoading}>
          Kirim
        </button>
      </div>
    </div>
  )
}
