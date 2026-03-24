import { useCalculator } from '../hooks/useCalculator';
import History from './History';

type ButtonType = 'digit' | 'operator' | 'unary' | 'clear' | 'backspace' | 'equals';

interface ButtonConfig {
  label: string;
  type: ButtonType;
  op?: string;
  span?: number;
}

const BUTTONS: ButtonConfig[] = [
  { label: 'C',  type: 'clear' },
  { label: '⌫',  type: 'backspace' },
  { label: '^',  type: 'operator', op: 'power' },
  { label: '÷',  type: 'operator', op: 'divide' },
  { label: '7',  type: 'digit' },
  { label: '8',  type: 'digit' },
  { label: '9',  type: 'digit' },
  { label: '×',  type: 'operator', op: 'multiply' },
  { label: '4',  type: 'digit' },
  { label: '5',  type: 'digit' },
  { label: '6',  type: 'digit' },
  { label: '−',  type: 'operator', op: 'subtract' },
  { label: '1',  type: 'digit' },
  { label: '2',  type: 'digit' },
  { label: '3',  type: 'digit' },
  { label: '+',  type: 'operator', op: 'add' },
  { label: '0',  type: 'digit' },
  { label: '.',  type: 'digit' },
  { label: '√',  type: 'unary',  op: 'sqrt' },
  { label: '%',  type: 'unary',  op: 'percentage' },
  { label: '=',  type: 'equals', span: 4 },
];

export default function Calculator() {
  const calc = useCalculator();

  function handleButton(btn: ButtonConfig) {
    switch (btn.type) {
      case 'digit':     calc.inputDigit(btn.label); break;
      case 'operator':  void calc.inputOperator(btn.op!); break;
      case 'unary':     void calc.pressUnary(btn.op!); break;
      case 'equals':    void calc.pressEquals(); break;
      case 'clear':     calc.clear(); break;
      case 'backspace': calc.backspace(); break;
    }
  }

  return (
    <div className="app-layout">
      <div className="calculator-card">
        <header className="calc-header">
          <span className="calc-brand">sezzle</span>
          <span className="calc-subtitle">calculator</span>
        </header>

        <div
          className={[
            'display',
            calc.error   ? 'display--error'   : '',
            calc.loading ? 'display--loading' : '',
          ].join(' ').trim()}
          aria-live="polite"
          aria-atomic="true"
        >
          <div className="display-expression">
            {calc.expression || '\u00A0'}
          </div>
          <div className="display-value" data-testid="display-value">
            {calc.error ?? calc.display}
          </div>
        </div>

        <div className="button-grid" role="group" aria-label="Calculator buttons">
          {BUTTONS.map((btn) => (
            <button
              key={btn.label}
              className={`btn btn--${btn.type}`}
              style={btn.span ? { gridColumn: `span ${btn.span}` } : undefined}
              onClick={() => handleButton(btn)}
              disabled={calc.loading}
              aria-label={btn.op ?? btn.label}
            >
              {btn.label}
            </button>
          ))}
        </div>
      </div>

      <History entries={calc.history} onRefresh={() => void calc.loadHistory()} />
    </div>
  );
}
