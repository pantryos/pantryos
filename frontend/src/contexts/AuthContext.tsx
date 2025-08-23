import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User } from '../types/api';
import apiService from '../services/api';

// Interface for the authentication context
interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

// Create the context with a default value
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Props for the AuthProvider component
interface AuthProviderProps {
  children: ReactNode;
}

// Provider component that wraps the app and makes auth object available to any child component
export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // Check if user is logged in on app start
  useEffect(() => {
    const initializeAuth = () => {
      // Check for a stored token and user in localStorage first.
      const token = apiService.getAuthToken();
      const storedUser = localStorage.getItem('user');

      if (token && storedUser) {
        try {
          // If a token exists, and we have a stored user object,
          // assume they are authenticated for now.
          const parsedUser: User = JSON.parse(storedUser);
          setUser(parsedUser);
        } catch (e) {
          console.error('Failed to parse user from localStorage:', e);
          apiService.clearAuthToken();
        }
      }

      // Crucially, set loading to false after this initial synchronous check.
      // The API interceptor will handle token validation on the next API call.
      setLoading(false);
    };

    initializeAuth();
  }, []);

  // Login function
  const login = async (email: string, password: string) => {
    try {
      // setLoading(true);
      const response = await apiService.login({ email, password });
      apiService.setAuthToken(response.token);
      setUser(response.user);
      localStorage.setItem('user', JSON.stringify(response.user)); // Store the user object
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    } finally {
      // setLoading(false);
    }
  };

  // Register function
  const register = async (email: string, password: string) => {
    try {
      setLoading(true);
      const response = await apiService.register({ email, password });
      apiService.setAuthToken(response.token);
      setUser(response.user);
      localStorage.setItem('user', JSON.stringify(response.user)); // Store the user object
    } catch (error) {
      console.error('Registration failed:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  // Logout function
  const logout = () => {
    apiService.clearAuthToken();
    setUser(null);
  };

    useEffect(() => {
    const handleAuthError = () => {
      logout();
    };

    window.addEventListener('auth-error', handleAuthError);
    return () => {
      window.removeEventListener('auth-error', handleAuthError);
    };
  }, []); 

  // Value object that will be given to the context
  const value: AuthContextType = {
    user,
    loading,
    login,
    register,
    logout,
    isAuthenticated: !!user,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

// Hook for child components to get the auth object and re-render when it changes
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};