const AGENTS = [
  { toolName: 'sales_data_tool', name: 'Sales Agent', division: 'Divisi Sales & Marketing' },
  { toolName: 'inventory_tool', name: 'Inventory Agent', division: 'Divisi Gudang' },
  { toolName: 'support_rag_tool', name: 'Support Agent', division: 'Customer Support · RAG' },
]

export default function AgentRoster({ agentStates }) {
  return (
    <div className="roster">
      <div className="roster-label">Agent Roster</div>
      {AGENTS.map((agent) => (
        <div
          key={agent.toolName}
          className={`agent-card ${agentStates[agent.toolName] || ''}`}
        >
          <div className="agent-top">
            <span className="agent-light"></span>
            <span className="agent-name">{agent.name}</span>
          </div>
          <span className="agent-division">{agent.division}</span>
        </div>
      ))}
    </div>
  )
}

export { AGENTS }
