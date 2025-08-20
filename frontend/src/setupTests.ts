// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom';
import './__mocks__/muiDataGrid'; 

// Mock TextEncoder and TextDecoder for Node.js environment
global.TextEncoder = class {
  encode(input: string): Uint8Array {
    return new Uint8Array(Buffer.from(input, 'utf8'));
  }
} as any;

global.TextDecoder = class {
  decode(input?: Uint8Array): string {
    return Buffer.from(input || []).toString('utf8');
  }
} as any;

// Mock localStorage
const localStorageMock: any = {
  store: {} as Record<string, string>,
  getItem: jest.fn((key: string): string | null => {
    return localStorageMock.store[key] || null;
  }),
  setItem: jest.fn((key: string, value: string): void => {
    localStorageMock.store[key] = value;
  }),
  removeItem: jest.fn((key: string): void => {
    delete localStorageMock.store[key];
  }),
  clear: jest.fn((): void => {
    localStorageMock.store = {};
  }),
  length: 0,
  key: jest.fn(),
};

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
  writable: true,
});

// Make localStorage methods available globally for tests
(global as any).localStorage = localStorageMock;

// Mock ResizeObserver for Recharts
global.ResizeObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));
