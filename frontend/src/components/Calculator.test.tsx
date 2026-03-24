import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import Calculator from './Calculator';
import * as api from '../api/calculator';

vi.mock('../api/calculator');

const mockCalculate = vi.mocked(api.calculate);
const mockFetchHistory = vi.mocked(api.fetchHistory);

beforeEach(() => {
  vi.clearAllMocks();
  mockFetchHistory.mockResolvedValue([]);
});

describe('Calculator — rendering', () => {
  it('renders the display with initial value 0', () => {
    render(<Calculator />);
    expect(screen.getByTestId('display-value')).toHaveTextContent('0');
  });

  it('renders all digit buttons', () => {
    render(<Calculator />);
    ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9'].forEach((d) => {
      expect(screen.getByRole('button', { name: d })).toBeInTheDocument();
    });
  });

  it('renders operator buttons', () => {
    render(<Calculator />);
    expect(screen.getByRole('button', { name: 'add' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'divide' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'multiply' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'subtract' })).toBeInTheDocument();
  });

  it('renders the history panel', () => {
    render(<Calculator />);
    expect(screen.getByText('History')).toBeInTheDocument();
    expect(screen.getByText('No calculations yet')).toBeInTheDocument();
  });
});

describe('Calculator — interactions', () => {
  it('updates display when digit buttons are clicked', () => {
    render(<Calculator />);
    fireEvent.click(screen.getByRole('button', { name: '4' }));
    fireEvent.click(screen.getByRole('button', { name: '2' }));
    expect(screen.getByTestId('display-value')).toHaveTextContent('42');
  });

  it('clears the display when C is pressed', () => {
    render(<Calculator />);
    fireEvent.click(screen.getByRole('button', { name: '9' }));
    fireEvent.click(screen.getByRole('button', { name: 'C' }));
    expect(screen.getByTestId('display-value')).toHaveTextContent('0');
  });

  it('shows result after a calculation', async () => {
    mockCalculate.mockResolvedValueOnce({ result: 15 });
    render(<Calculator />);

    fireEvent.click(screen.getByRole('button', { name: '1' }));
    fireEvent.click(screen.getByRole('button', { name: '0' }));
    fireEvent.click(screen.getByRole('button', { name: 'add' }));
    fireEvent.click(screen.getByRole('button', { name: '5' }));
    fireEvent.click(screen.getByRole('button', { name: '=' }));

    await screen.findByText('15');
  });

  it('shows error message for division by zero', async () => {
    mockCalculate.mockResolvedValueOnce({ error: 'division by zero' });
    render(<Calculator />);

    fireEvent.click(screen.getByRole('button', { name: '5' }));
    fireEvent.click(screen.getByRole('button', { name: 'divide' }));
    fireEvent.click(screen.getByRole('button', { name: '0' }));
    fireEvent.click(screen.getByRole('button', { name: '=' }));

    await screen.findByText('division by zero');
  });

  it('disables buttons while loading', async () => {
    mockCalculate.mockImplementation(() => new Promise(() => {})); // never resolves
    render(<Calculator />);

    fireEvent.click(screen.getByRole('button', { name: '9' }));
    fireEvent.click(screen.getByRole('button', { name: 'sqrt' }));

    // All buttons should be disabled during load
    const buttons = screen.getAllByRole('button');
    expect(buttons.some((b) => b.hasAttribute('disabled'))).toBe(true);
  });
});
