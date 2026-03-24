import { useState, useEffect } from 'react';
import { calculate, fetchHistory } from '../api/calculator';
import type { HistoryEntry } from '../api/calculator';

export const OP_SYMBOLS: Record<string, string> = {
  add: '+',
  subtract: '−',
  multiply: '×',
  divide: '÷',
  power: '^',
  sqrt: '√',
  percentage: '%',
};

function formatNumber(n: number): string {
  if (!isFinite(n)) return 'Error';
  if (Number.isInteger(n) && Math.abs(n) < 1e15) return String(n);
  return String(parseFloat(n.toPrecision(10)));
}

export function useCalculator() {
  const [display, setDisplay] = useState('0');
  const [a, setA] = useState(0);
  const [pendingOp, setPendingOp] = useState<string | null>(null);
  const [waitingForB, setWaitingForB] = useState(false);
  const [justCalculated, setJustCalculated] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [history, setHistory] = useState<HistoryEntry[]>([]);

  useEffect(() => {
    void loadHistory();
  }, []);

  async function loadHistory() {
    try {
      setHistory(await fetchHistory());
    } catch {
      // backend may not be running yet
    }
  }

  async function runCalculation(
    opA: number,
    opB: number,
    op: string,
  ): Promise<number | null> {
    setLoading(true);
    setError(null);
    try {
      const res = await calculate({ a: opA, b: opB, operation: op });
      if (res.error) {
        setError(res.error);
        return null;
      }
      await loadHistory();
      return res.result ?? 0;
    } catch {
      setError('Could not reach the server. Please try again.');
      return null;
    } finally {
      setLoading(false);
    }
  }

  function inputDigit(digit: string) {
    setError(null);
    if (waitingForB || justCalculated) {
      setDisplay(digit === '.' ? '0.' : digit);
      setWaitingForB(false);
      setJustCalculated(false);
      return;
    }
    setDisplay((prev) => {
      if (digit === '.' && prev.includes('.')) return prev;
      if (prev === '0' && digit !== '.') return digit;
      return prev + digit;
    });
  }

  async function inputOperator(op: string) {
    setError(null);
    setJustCalculated(false);

    if (pendingOp && !waitingForB) {
      // Chain: calculate intermediate result first
      const result = await runCalculation(a, parseFloat(display), pendingOp);
      if (result === null) return;
      setDisplay(formatNumber(result));
      setA(result);
      setPendingOp(op);
      setWaitingForB(true);
    } else if (pendingOp && waitingForB) {
      // User changed their mind about the operator
      setPendingOp(op);
    } else {
      setA(parseFloat(display));
      setPendingOp(op);
      setWaitingForB(true);
    }
  }

  async function pressEquals() {
    if (!pendingOp) return;
    const b = waitingForB ? a : parseFloat(display);
    const result = await runCalculation(a, b, pendingOp);
    if (result === null) return;
    setDisplay(formatNumber(result));
    setA(result);
    setPendingOp(null);
    setWaitingForB(false);
    setJustCalculated(true);
  }

  async function pressUnary(op: string) {
    const result = await runCalculation(parseFloat(display), 0, op);
    if (result === null) return;
    setDisplay(formatNumber(result));
    setA(result);
    setPendingOp(null);
    setWaitingForB(false);
    setJustCalculated(true);
  }

  function clear() {
    setDisplay('0');
    setA(0);
    setPendingOp(null);
    setWaitingForB(false);
    setJustCalculated(false);
    setError(null);
  }

  function backspace() {
    if (waitingForB || justCalculated) return;
    setError(null);
    setDisplay((prev) => (prev.length > 1 ? prev.slice(0, -1) : '0'));
  }

  const expression = pendingOp
    ? `${formatNumber(a)} ${OP_SYMBOLS[pendingOp] ?? pendingOp}${!waitingForB ? ` ${display}` : ''}`
    : '';

  return {
    display,
    expression,
    pendingOp,
    error,
    loading,
    history,
    loadHistory,
    inputDigit,
    inputOperator,
    pressEquals,
    pressUnary,
    clear,
    backspace,
  };
}
