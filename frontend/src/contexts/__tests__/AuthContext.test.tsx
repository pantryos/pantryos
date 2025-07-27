import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { AuthProvider, useAuth } from '../AuthContext';
import { apiService } from '../../services/api';

// Mock the API service
jest.mock('../../services/api');
const mockedApiService = apiService as jest.Mocked<typeof apiService>;

// Test component to access auth context
const TestComponent = () => {
  const { user, isAuthenticated, login, register, logout } = useAuth();
  
  return (
    <div>
      <div data-testid="user">{user ? user.email : 'No user'}</div>
      <div data-testid="authenticated">{isAuthenticated ? 'true' : 'false'}</div>
      <button onClick={() => login('test@example.com', 'password')}>
        Login
      </button>
      <button onClick={() => register('test@example.com', 'password')}>
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

    it('should initialize with user when token exists', async () => {
      const mockUser = { id: 1, email: 'test@example.com', account_id: 1, role: 'user', created_at: '2024-01-01' };
      localStorage.setItem('token', 'test-token');
      localStorage.setItem('user', JSON.stringify(mockUser));
      
      mockedApiService.getCurrentUser.mockResolvedValue(mockUser);

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );
      
      await waitFor(() => {
        expect(screen.getByTestId('user')).toHaveTextContent('test@example.com');
        expect(screen.getByTestId('authenticated')).toHaveTextContent('true');
      });
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
          password: 'password',
          name: 'Test User'
        });
        expect(screen.getByTestId('user')).toHaveTextContent('test@example.com');
        expect(screen.getByTestId('authenticated')).toHaveTextContent('true');
      });
    });

    it('should handle register error', async () => {
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
        expect(localStorage.getItem('token')).toBeNull();
        expect(localStorage.getItem('user')).toBeNull();
      });
    });
  });

  describe('token management', () => {
    it('should store token and user in localStorage on successful login', async () => {
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
        expect(localStorage.getItem('token')).toBe('test-token');
        expect(localStorage.getItem('user')).toBe(JSON.stringify(mockUser));
      });
    });

    it('should clear token and user from localStorage on logout', async () => {
      const mockUser = { id: 1, email: 'test@example.com', account_id: 1, role: 'user', created_at: '2024-01-01' };
      const mockResponse = { token: 'test-token', user: mockUser };
      
      mockedApiService.login.mockResolvedValue(mockResponse);

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      // Login first
      fireEvent.click(screen.getByText('Login'));
      await waitFor(() => {
        expect(localStorage.getItem('token')).toBe('test-token');
      });

      // Then logout
      fireEvent.click(screen.getByText('Logout'));

      await waitFor(() => {
        expect(localStorage.getItem('token')).toBeNull();
        expect(localStorage.getItem('user')).toBeNull();
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
        expect(consoleSpy).toHaveBeenCalledWith('Login error:', error);
      });

      consoleSpy.mockRestore();
    });

    it('should handle invalid JSON in localStorage', () => {
      localStorage.setItem('token', 'test-token');
      localStorage.setItem('user', 'invalid-json');

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );
      
      expect(screen.getByTestId('user')).toHaveTextContent('No user');
      expect(screen.getByTestId('authenticated')).toHaveTextContent('false');
    });
  });
}); 