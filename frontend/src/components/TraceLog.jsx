import { useEffect, useRef } from 'react'

export default function TraceLog({ entries }) {
  const scrollRef = useRef(null)

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [entries])

  return (
    <div className="trace">
      <div className="trace-header">Route Trace</div>
      <div className="trace-log" ref={scrollRef}>
        {entries.length === 0 && (
          <div className="trace-empty">— menunggu query pertama —</div>
        )}
        {entries.map((entry, i) => (
          <div key={i} className={`trace-line ${entry.type || ''}`}>
            <span className="trace-time">{entry.time}</span>
            <span className="trace-text">{entry.text}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
