import type { HistoryEntry } from '../api/calculator';
import { OP_SYMBOLS } from '../hooks/useCalculator';

function formatEntry(entry: HistoryEntry): string {
  const { a, b, operation, result } = entry;
  const sym = OP_SYMBOLS[operation] ?? operation;
  if (operation === 'sqrt')       return `√ ${a} = ${result}`;
  if (operation === 'percentage') return `${a} % = ${result}`;
  return `${a} ${sym} ${b} = ${result}`;
}

function timeAgo(timestamp: string): string {
  const diff = Date.now() - new Date(timestamp).getTime();
  const s = Math.floor(diff / 1000);
  if (s < 60)  return `${s}s ago`;
  const m = Math.floor(s / 60);
  if (m < 60)  return `${m}m ago`;
  const h = Math.floor(m / 60);
  return `${h}h ago`;
}

interface Props {
  entries: HistoryEntry[];
  onRefresh: () => void;
}

export default function History({ entries, onRefresh }: Props) {
  return (
    <aside className="history-panel">
      <div className="history-header">
        <h2 className="history-title">History</h2>
        <button
          className="history-refresh"
          onClick={onRefresh}
          aria-label="Refresh history"
          title="Refresh"
        >
          ↻
        </button>
      </div>

      {entries.length === 0 ? (
        <p className="history-empty">No calculations yet</p>
      ) : (
        <ul className="history-list">
          {[...entries].reverse().map((entry, i) => (
            <li key={i} className="history-item">
              <span className="history-expr">{formatEntry(entry)}</span>
              <span className="history-time">{timeAgo(entry.timestamp)}</span>
            </li>
          ))}
        </ul>
      )}
    </aside>
  );
}
