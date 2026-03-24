import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { calculate, fetchHistory } from './calculator';

const mockFetch = vi.fn();
vi.stubGlobal('fetch', mockFetch);

beforeEach(() => vi.clearAllMocks());
afterEach(() => vi.restoreAllMocks());

function mockJson(body: unknown, status = 200) {
  mockFetch.mockResolvedValueOnce({
    status,
    json: () => Promise.resolve(body),
  });
}

describe('calculate()', () => {
  it('sends a POST request to /calculate with the correct body', async () => {
    mockJson({ result: 7 });
    await calculate({ a: 3, b: 4, operation: 'add' });

    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/calculate',
      expect.objectContaining({
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ a: 3, b: 4, operation: 'add' }),
      }),
    );
  });

  it('returns the result on success', async () => {
    mockJson({ result: 15 });
    const res = await calculate({ a: 10, b: 5, operation: 'add' });
    expect(res.result).toBe(15);
    expect(res.error).toBeUndefined();
  });

  it('returns zero result without dropping it', async () => {
    mockJson({ result: 0 });
    const res = await calculate({ a: 0, b: 0, operation: 'add' });
    expect(res.result).toBe(0);
  });

  it('returns the error field on API error', async () => {
    mockJson({ error: 'division by zero' }, 422);
    const res = await calculate({ a: 5, b: 0, operation: 'divide' });
    expect(res.error).toBe('division by zero');
    expect(res.result).toBeUndefined();
  });

  it('propagates network errors for the caller to handle', async () => {
    mockFetch.mockRejectedValueOnce(new Error('Network Error'));
    await expect(calculate({ a: 1, b: 1, operation: 'add' })).rejects.toThrow('Network Error');
  });
});

describe('fetchHistory()', () => {
  it('sends a GET request to /history', async () => {
    mockJson([]);
    await fetchHistory();

    expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/history');
  });

  it('returns an array of history entries', async () => {
    const entries = [
      { a: 3, b: 4, operation: 'add', result: 7, timestamp: '2026-03-24T10:00:00Z' },
    ];
    mockJson(entries);
    const res = await fetchHistory();
    expect(res).toEqual(entries);
  });

  it('returns an empty array when history is empty', async () => {
    mockJson([]);
    const res = await fetchHistory();
    expect(res).toEqual([]);
  });

  it('propagates network errors for the caller to handle', async () => {
    mockFetch.mockRejectedValueOnce(new Error('Network Error'));
    await expect(fetchHistory()).rejects.toThrow('Network Error');
  });
});
