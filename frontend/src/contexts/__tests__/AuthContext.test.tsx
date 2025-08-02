import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { AuthProvider, useAuth } from '../AuthContext';

// Mock the API service as a default export
jest.mock('../../services/api', () => ({
  __esModule: true,
  default: {
    login: jest.fn(),
    register: jest.fn(),
    getCurrentUser: jest.fn(),
    setAuthToken: jest.fn(),
    getAuthToken: jest.fn(),
    clearAuthToken: jest.fn(),
    isAuthenticated: jest.fn(),
  },
}));

// Import the mocked service
import apiService from '../../services/api';
const mockedApiService = apiService as jest.Mocked<typeof apiService>;

// Test component to access auth context
const TestComponent = () => {
  const { user, isAuthenticated, login, register, logout } = useAuth();
  
  return (
    <div>
      <div data-testid="user">{user ? user.email : 'No user'}</div>
      <div data-testid="authenticated">{isAuthenticated ? 'true' : 'false'}</div>
      <button onClick={async () => {
        try {
          await login('test@example.com', 'password');
        } catch (error) {
          // Error is expected in error tests
        }
      }}>
        Login
      </button>
      <button onClick={async () => {
        try {
          await register('test@example.com', 'password');
        } catch (error) {
          // Error is expected in error tests
        }
      }}>
        Register
      </button>
      <button onClick={logout}>Logout</button>
    </div>
  );
};

describe('AuthContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    localStorage.clear();
    
    // Mock the API service methods
    mockedApiService.setAuthToken.mockImplementation((token: string) => {
      localStorage.setItem('token', token);
    });
    
    mockedApiService.clearAuthToken.mockImplementation(() => {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    });
    
    mockedApiService.getAuthToken.mockImplementation(() => {
      return localStorage.getItem('token');
    });
    
    mockedApiService.isAuthenticated.mockImplementation(() => {
      return !!localStorage.getItem('token');
    });
  });

  describe('AuthProvider', () => {
    it('should render children', () => {
      render(
        <AuthProvider>
          <div data-testid="child">Child Component</div>
        </AuthProvider>
      );
      
      expect(screen.getByTestId('child')).toBeInTheDocument();
    });

    it('should initialize with no user when no token exists', () => {
      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );
      
      expect(screen.getByTestId('user')).toHaveTextContent('No user');
      expect(screen.getByTestId('authenticated')).toHaveTextContent('false');
    });
  });

  describe('login functionality', () => {
    it('should login user successfully', async () => {
      const mockUser = { id: 1, email: 'test@example.com', account_id: 1, role: 'user', created_at: '2024-01-01' };
      const mockResponse = { token: 'test-token', user: mockUser };
      
      mockedApiService.login.mockResolvedValue(mockResponse);

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      fireEvent.click(screen.getByText('Login'));

      await waitFor(() => {
        expect(mockedApiService.login).toHaveBeenCalledWith({
          email: 'test@example.com',
          password: 'password'
        });
        expect(screen.getByTestId('user')).toHaveTextContent('test@example.com');
        expect(screen.getByTestId('authenticated')).toHaveTextContent('true');
      });
    });

    it('should handle login error', async () => {
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      const error = new Error('Invalid credentials');
      mockedApiService.login.mockRejectedValue(error);

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      fireEvent.click(screen.getByText('Login'));

      await waitFor(() => {
        expect(mockedApiService.login).toHaveBeenCalled();
        expect(screen.getByTestId('user')).toHaveTextContent('No user');
        expect(screen.getByTestId('authenticated')).toHaveTextContent('false');
      });

      consoleSpy.mockRestore();
    });
  });

  describe('register functionality', () => {
    it('should register user successfully', async () => {
      const mockUser = { id: 1, email: 'test@example.com', account_id: 1, role: 'user', created_at: '2024-01-01' };
      const mockResponse = { token: 'test-token', user: mockUser };
      
      mockedApiService.register.mockResolvedValue(mockResponse);

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      fireEvent.click(screen.getByText('Register'));

      await waitFor(() => {
        expect(mockedApiService.register).toHaveBeenCalledWith({
          email: 'test@example.com',
          password: 'password'
        });
        expect(screen.getByTestId('user')).toHaveTextContent('test@example.com');
        expect(screen.getByTestId('authenticated')).toHaveTextContent('true');
      });
    });

    it('should handle register error', async () => {
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      const error = new Error('Email already exists');
      mockedApiService.register.mockRejectedValue(error);

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      fireEvent.click(screen.getByText('Register'));

      await waitFor(() => {
        expect(mockedApiService.register).toHaveBeenCalled();
        expect(screen.getByTestId('user')).toHaveTextContent('No user');
        expect(screen.getByTestId('authenticated')).toHaveTextContent('false');
      });

      consoleSpy.mockRestore();
    });
  });

  describe('logout functionality', () => {
    it('should logout user successfully', async () => {
      const mockUser = { id: 1, email: 'test@example.com', account_id: 1, role: 'user', created_at: '2024-01-01' };
      const mockResponse = { token: 'test-token', user: mockUser };
      
      mockedApiService.login.mockResolvedValue(mockResponse);

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      // First login
      fireEvent.click(screen.getByText('Login'));
      await waitFor(() => {
        expect(screen.getByTestId('authenticated')).toHaveTextContent('true');
      });

      // Then logout
      fireEvent.click(screen.getByText('Logout'));

      await waitFor(() => {
        expect(screen.getByTestId('user')).toHaveTextContent('No user');
        expect(screen.getByTestId('authenticated')).toHaveTextContent('false');
      });
    });
  });

  describe('error handling', () => {
    it('should handle API errors gracefully', async () => {
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      const error = new Error('Network error');
      mockedApiService.login.mockRejectedValue(error);

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      fireEvent.click(screen.getByText('Login'));

      await waitFor(() => {
        expect(consoleSpy).toHaveBeenCalledWith('Login failed:', error);
      });

      consoleSpy.mockRestore();
    });
  });
}); 