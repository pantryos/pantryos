import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

// Mock the auth context to avoid authentication issues in tests
jest.mock('./contexts/AuthContext', () => ({
  AuthProvider: ({ children }: { children: React.ReactNode }) => <div data-testid="auth-provider">{children}</div>,
  useAuth: () => ({
    user: null,
    loading: false,
    login: jest.fn(),
    register: jest.fn(),
    logout: jest.fn(),
    isAuthenticated: false,
  }),
}));

// Mock react-router-dom
jest.mock('react-router-dom', () => ({
  BrowserRouter: ({ children }: { children: React.ReactNode }) => <div data-testid="router">{children}</div>,
  Routes: ({ children }: { children: React.ReactNode }) => <div data-testid="routes">{children}</div>,
  Route: ({ children }: { children: React.ReactNode }) => <div data-testid="route">{children}</div>,
  Navigate: () => <div data-testid="navigate">Navigate</div>,
}));

test('renders app with providers', () => {
  render(<App />);
  
  // Check that the auth provider is rendered
  expect(screen.getByTestId('auth-provider')).toBeInTheDocument();
  
  // Check that the router is rendered
  expect(screen.getByTestId('router')).toBeInTheDocument();
  
  // Check that routes are rendered
  expect(screen.getByTestId('routes')).toBeInTheDocument();
});
