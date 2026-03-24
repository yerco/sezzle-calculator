const API_URL = (import.meta.env.VITE_API_URL as string | undefined) ?? 'http://localhost:8080';

export interface CalculationRequest {
  a: number;
  b: number;
  operation: string;
}

export interface CalculationResponse {
  result?: number;
  error?: string;
}

export interface HistoryEntry {
  a: number;
  b: number;
  operation: string;
  result: number;
  timestamp: string;
}

export async function calculate(req: CalculationRequest): Promise<CalculationResponse> {
  const res = await fetch(`${API_URL}/calculate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  });
  return res.json() as Promise<CalculationResponse>;
}

export async function fetchHistory(): Promise<HistoryEntry[]> {
  const res = await fetch(`${API_URL}/history`);
  return res.json() as Promise<HistoryEntry[]>;
}
