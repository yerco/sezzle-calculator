import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useCalculator } from './useCalculator';
import * as api from '../api/calculator';

vi.mock('../api/calculator');

const mockCalculate = vi.mocked(api.calculate);
const mockFetchHistory = vi.mocked(api.fetchHistory);

beforeEach(() => {
  vi.clearAllMocks();
  mockFetchHistory.mockResolvedValue([]);
});

describe('useCalculator — initial state', () => {
  it('displays 0 on mount', () => {
    const { result } = renderHook(() => useCalculator());
    expect(result.current.display).toBe('0');
    expect(result.current.error).toBeNull();
    expect(result.current.pendingOp).toBeNull();
  });
});

describe('useCalculator — digit input', () => {
  it('replaces leading zero with digit', () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.inputDigit('5'));
    expect(result.current.display).toBe('5');
  });

  it('appends digits', () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.inputDigit('1'));
    act(() => result.current.inputDigit('2'));
    act(() => result.current.inputDigit('3'));
    expect(result.current.display).toBe('123');
  });

  it('prevents duplicate decimal points', () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.inputDigit('3'));
    act(() => result.current.inputDigit('.'));
    act(() => result.current.inputDigit('.'));
    expect(result.current.display).toBe('3.');
  });

  it('handles decimal-first input as 0.', () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.inputDigit('.'));
    expect(result.current.display).toBe('0.');
  });
});

describe('useCalculator — clear & backspace', () => {
  it('clear resets state', () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.inputDigit('9'));
    act(() => result.current.clear());
    expect(result.current.display).toBe('0');
  });

  it('backspace removes last character', () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.inputDigit('4'));
    act(() => result.current.inputDigit('2'));
    act(() => result.current.backspace());
    expect(result.current.display).toBe('4');
  });

  it('backspace bottoms out at 0', () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.inputDigit('7'));
    act(() => result.current.backspace());
    expect(result.current.display).toBe('0');
  });
});

describe('useCalculator — binary operations', () => {
  it('calls API with correct operands on equals', async () => {
    mockCalculate.mockResolvedValueOnce({ result: 7 });
    const { result } = renderHook(() => useCalculator());

    act(() => result.current.inputDigit('3'));
    await act(() => result.current.inputOperator('add'));
    act(() => result.current.inputDigit('4'));
    await act(() => result.current.pressEquals());

    expect(mockCalculate).toHaveBeenCalledWith({ a: 3, b: 4, operation: 'add' });
    expect(result.current.display).toBe('7');
  });

  it('shows a friendly message when fetch throws (network failure)', async () => {
    mockCalculate.mockRejectedValueOnce(new Error('Network Error'));
    const { result } = renderHook(() => useCalculator());

    act(() => result.current.inputDigit('5'));
    await act(() => result.current.inputOperator('add'));
    act(() => result.current.inputDigit('3'));
    await act(() => result.current.pressEquals());

    expect(result.current.error).toBe('Could not reach the server. Please try again.');
    expect(result.current.loading).toBe(false);
  });

  it('shows error message on API error', async () => {
    mockCalculate.mockResolvedValueOnce({ error: 'division by zero' });
    const { result } = renderHook(() => useCalculator());

    act(() => result.current.inputDigit('5'));
    await act(() => result.current.inputOperator('divide'));
    act(() => result.current.inputDigit('0'));
    await act(() => result.current.pressEquals());

    expect(result.current.error).toBe('division by zero');
  });

  it('chains operations — calculates intermediate result', async () => {
    mockCalculate.mockResolvedValueOnce({ result: 10 });
    const { result } = renderHook(() => useCalculator());

    act(() => result.current.inputDigit('3'));
    await act(() => result.current.inputOperator('add'));
    act(() => result.current.inputDigit('7'));
    // Pressing another operator should chain
    await act(() => result.current.inputOperator('multiply'));

    expect(mockCalculate).toHaveBeenCalledWith({ a: 3, b: 7, operation: 'add' });
    expect(result.current.display).toBe('10');
  });

  it('handles result of zero correctly', async () => {
    mockCalculate.mockResolvedValueOnce({ result: 0 });
    const { result } = renderHook(() => useCalculator());

    act(() => result.current.inputDigit('5'));
    await act(() => result.current.inputOperator('subtract'));
    act(() => result.current.inputDigit('5'));
    await act(() => result.current.pressEquals());

    expect(result.current.display).toBe('0');
    expect(result.current.error).toBeNull();
  });
});

describe('useCalculator — unary operations', () => {
  it('calls API with b=0 for sqrt', async () => {
    mockCalculate.mockResolvedValueOnce({ result: 3 });
    const { result } = renderHook(() => useCalculator());

    act(() => result.current.inputDigit('9'));
    await act(() => result.current.pressUnary('sqrt'));

    expect(mockCalculate).toHaveBeenCalledWith({ a: 9, b: 0, operation: 'sqrt' });
    expect(result.current.display).toBe('3');
  });

  it('calls API for percentage', async () => {
    mockCalculate.mockResolvedValueOnce({ result: 0.25 });
    const { result } = renderHook(() => useCalculator());

    act(() => result.current.inputDigit('2'));
    act(() => result.current.inputDigit('5'));
    await act(() => result.current.pressUnary('percentage'));

    expect(mockCalculate).toHaveBeenCalledWith({ a: 25, b: 0, operation: 'percentage' });
    expect(result.current.display).toBe('0.25');
  });
});

describe('useCalculator — expression display', () => {
  it('shows expression after operator pressed', async () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.inputDigit('8'));
    await act(() => result.current.inputOperator('divide'));
    expect(result.current.expression).toContain('8');
    expect(result.current.expression).toContain('÷');
  });
});

describe('useCalculator — history', () => {
  it('loads history on mount', async () => {
    const entries = [
      { a: 1, b: 2, operation: 'add', result: 3, timestamp: new Date().toISOString() },
    ];
    mockFetchHistory.mockResolvedValueOnce(entries);
    const { result } = renderHook(() => useCalculator());
    // Wait for the useEffect to resolve
    await act(async () => { await Promise.resolve(); });
    expect(result.current.history).toEqual(entries);
  });

  it('refreshes history after a calculation', async () => {
    const entry = { a: 2, b: 3, operation: 'multiply', result: 6, timestamp: new Date().toISOString() };
    mockCalculate.mockResolvedValueOnce({ result: 6 });
    mockFetchHistory.mockResolvedValueOnce([]).mockResolvedValueOnce([entry]);
    const { result } = renderHook(() => useCalculator());

    act(() => result.current.inputDigit('2'));
    await act(() => result.current.inputOperator('multiply'));
    act(() => result.current.inputDigit('3'));
    await act(() => result.current.pressEquals());

    expect(result.current.history).toEqual([entry]);
  });
});
