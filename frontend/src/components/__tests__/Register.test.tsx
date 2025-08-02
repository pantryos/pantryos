import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import Register from '../Register';
import { useAuth } from '../../contexts/AuthContext';

// Mock the auth context
jest.mock('../../contexts/AuthContext');
const mockedUseAuth = useAuth as jest.MockedFunction<typeof useAuth>;

// Mock react-router-dom
jest.mock('react-router-dom', () => ({
  BrowserRouter: ({ children }: any) => <div>{children}</div>,
  useNavigate: () => jest.fn(),
}));

const theme = createTheme();

// Wrapper component to provide necessary context
const TestWrapper = ({ children }: { children: React.ReactNode }) => (
  <BrowserRouter>
    <ThemeProvider theme={theme}>
      {children}
    </ThemeProvider>
  </BrowserRouter>
);

describe('Register Component', () => {
  const mockRegister = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
    mockedUseAuth.mockReturnValue({
      user: null,
      loading: false,
      login: jest.fn(),
      register: mockRegister,
      logout: jest.fn(),
      isAuthenticated: false,
    });
  });

  describe('rendering', () => {
    it('should render registration form', () => {
      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      expect(screen.getByRole('heading', { name: 'Create Account' })).toBeInTheDocument();
      expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/^password\s*\*?$/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/^confirm password\s*\*?$/i)).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /create account/i })).toBeInTheDocument();
      expect(screen.getByText(/already have an account/i)).toBeInTheDocument();
    });

    it('should show loading state during form submission', async () => {
      // Mock a delayed register response to test loading state
      mockRegister.mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)));

      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      const emailInput = screen.getByLabelText(/email/i);
      const passwordInput = screen.getByLabelText(/^password\s*\*?$/i);
      const confirmPasswordInput = screen.getByLabelText(/^confirm password\s*\*?$/i);
      const registerButton = screen.getByRole('button', { name: /create account/i });

      fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
      fireEvent.change(passwordInput, { target: { value: 'password123' } });
      fireEvent.change(confirmPasswordInput, { target: { value: 'password123' } });
      fireEvent.click(registerButton);

      // Check that button is disabled during loading
      expect(registerButton).toBeDisabled();
      
      // Wait for the promise to resolve
      await waitFor(() => {
        expect(mockRegister).toHaveBeenCalled();
      });
    });
  });

  describe('form validation', () => {
    it('should show error for mismatched passwords', async () => {
      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      const emailInput = screen.getByLabelText(/email/i);
      const passwordInput = screen.getByLabelText(/^password\s*\*?$/i);
      const confirmPasswordInput = screen.getByLabelText(/^confirm password\s*\*?$/i);
      
      fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
      fireEvent.change(passwordInput, { target: { value: 'password123' } });
      fireEvent.change(confirmPasswordInput, { target: { value: 'differentpassword' } });

      const registerButton = screen.getByRole('button', { name: /create account/i });
      fireEvent.click(registerButton);

      await waitFor(() => {
        expect(screen.getByText(/passwords do not match/i)).toBeInTheDocument();
      });
    });
  });

  describe('form submission', () => {
    it('should call register function with valid credentials', async () => {
      mockRegister.mockResolvedValue(undefined);

      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      const emailInput = screen.getByLabelText(/email/i);
      const passwordInput = screen.getByLabelText(/^password\s*\*?$/i);
      const confirmPasswordInput = screen.getByLabelText(/^confirm password\s*\*?$/i);
      const registerButton = screen.getByRole('button', { name: /create account/i });

      fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
      fireEvent.change(passwordInput, { target: { value: 'password123' } });
      fireEvent.change(confirmPasswordInput, { target: { value: 'password123' } });
      fireEvent.click(registerButton);

      await waitFor(() => {
        expect(mockRegister).toHaveBeenCalledWith('test@example.com', 'password123');
      });
    });

    it('should handle registration error', async () => {
      const error = new Error('Email already exists');
      mockRegister.mockRejectedValue(error);

      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      const emailInput = screen.getByLabelText(/email/i);
      const passwordInput = screen.getByLabelText(/^password\s*\*?$/i);
      const confirmPasswordInput = screen.getByLabelText(/^confirm password\s*\*?$/i);
      const registerButton = screen.getByRole('button', { name: /create account/i });

      fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
      fireEvent.change(passwordInput, { target: { value: 'password123' } });
      fireEvent.change(confirmPasswordInput, { target: { value: 'password123' } });
      fireEvent.click(registerButton);

      await waitFor(() => {
        expect(mockRegister).toHaveBeenCalledWith('test@example.com', 'password123');
      });
    });

    it('should show error message when registration fails', async () => {
      const error = { response: { data: { error: 'Email already exists' } } };
      mockRegister.mockRejectedValue(error);

      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      const emailInput = screen.getByLabelText(/email/i);
      const passwordInput = screen.getByLabelText(/^password\s*\*?$/i);
      const confirmPasswordInput = screen.getByLabelText(/^confirm password\s*\*?$/i);
      const registerButton = screen.getByRole('button', { name: /create account/i });

      fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
      fireEvent.change(passwordInput, { target: { value: 'password123' } });
      fireEvent.change(confirmPasswordInput, { target: { value: 'password123' } });
      fireEvent.click(registerButton);

      await waitFor(() => {
        expect(screen.getByText('Email already exists')).toBeInTheDocument();
      });
    });

    it('should show generic error message when registration fails without specific error', async () => {
      const error = new Error('Network error');
      mockRegister.mockRejectedValue(error);

      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      const emailInput = screen.getByLabelText(/email/i);
      const passwordInput = screen.getByLabelText(/^password\s*\*?$/i);
      const confirmPasswordInput = screen.getByLabelText(/^confirm password\s*\*?$/i);
      const registerButton = screen.getByRole('button', { name: /create account/i });

      fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
      fireEvent.change(passwordInput, { target: { value: 'password123' } });
      fireEvent.change(confirmPasswordInput, { target: { value: 'password123' } });
      fireEvent.click(registerButton);

      await waitFor(() => {
        expect(screen.getByText('Registration failed. Please try again.')).toBeInTheDocument();
      });
    });
  });

  describe('navigation', () => {
    it('should have login button', () => {
      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      const loginButton = screen.getByText(/sign in/i);
      expect(loginButton).toBeInTheDocument();
      expect(loginButton.tagName).toBe('BUTTON');
    });
  });

  describe('accessibility', () => {
    it('should have proper form labels', () => {
      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/^password\s*\*?$/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/^confirm password\s*\*?$/i)).toBeInTheDocument();
    });

    it('should have proper button roles', () => {
      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      expect(screen.getByRole('button', { name: /create account/i })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument();
    });
  });

  describe('invitation-based registration', () => {
    it('should display invitation message', () => {
      render(
        <TestWrapper>
          <Register />
        </TestWrapper>
      );

      expect(screen.getByText(/join stok using your invitation/i)).toBeInTheDocument();
      expect(screen.getByText(/you must have been invited by an account administrator/i)).toBeInTheDocument();
    });
  });
}); 