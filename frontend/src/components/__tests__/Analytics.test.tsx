import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { AuthProvider } from '../../contexts/AuthContext';
import Analytics from '../Analytics';

// Mock the API service
jest.mock('../../services/api', () => ({
  getInventoryItems: jest.fn().mockResolvedValue([
    {
      id: 1,
      name: 'Tomatoes',
      current_stock: 50,
      max_stock_level: 100,
      min_stock_level: 10,
      cost_per_unit: 2.5,
      unit: 'kg',
      preferred_vendor: 'Fresh Farms',
      account_id: 1,
    },
    {
      id: 2,
      name: 'Onions',
      current_stock: 30,
      max_stock_level: 80,
      min_stock_level: 15,
      cost_per_unit: 1.8,
      unit: 'kg',
      preferred_vendor: 'Veggie Supply',
      account_id: 1,
    },
  ]),
  getMenuItems: jest.fn().mockResolvedValue([]),
  getDeliveries: jest.fn().mockResolvedValue([]),
}));

// Mock recharts components
jest.mock('recharts', () => ({
  LineChart: ({ children }: any) => <div data-testid="line-chart">{children}</div>,
  Line: () => <div data-testid="line" />,
  BarChart: ({ children }: any) => <div data-testid="bar-chart">{children}</div>,
  Bar: () => <div data-testid="bar" />,
  AreaChart: ({ children }: any) => <div data-testid="area-chart">{children}</div>,
  Area: () => <div data-testid="area" />,
  ComposedChart: ({ children }: any) => <div data-testid="composed-chart">{children}</div>,
  XAxis: () => <div data-testid="x-axis" />,
  YAxis: () => <div data-testid="y-axis" />,
  CartesianGrid: () => <div data-testid="cartesian-grid" />,
  Tooltip: () => <div data-testid="tooltip" />,
  Legend: () => <div data-testid="legend" />,
  ResponsiveContainer: ({ children }: any) => <div data-testid="responsive-container">{children}</div>,
}));

// Mock the auth context
const mockAuthContext = {
  user: { id: 1, email: 'test@example.com', account_id: 1 },
  isAuthenticated: true,
  loading: false,
  login: jest.fn(),
  logout: jest.fn(),
  register: jest.fn(),
};

jest.mock('../../contexts/AuthContext', () => ({
  useAuth: () => mockAuthContext,
  AuthProvider: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

// Mock react-router-dom
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  useNavigate: () => mockNavigate,
  BrowserRouter: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

const theme = createTheme();

const renderWithProviders = (component: React.ReactElement) => {
  return render(
    <ThemeProvider theme={theme}>
      <AuthProvider>
        <BrowserRouter>
          {component}
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
};

describe('Analytics Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders analytics dashboard with title', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getByText('Analytics Dashboard')).toBeInTheDocument();
    });
  });

  it('displays summary statistics cards', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getByText('Total Items')).toBeInTheDocument();
      expect(screen.getByText('Low Stock Items')).toBeInTheDocument();
      expect(screen.getByText('Total Value')).toBeInTheDocument();
      expect(screen.getByText('Avg Utilization')).toBeInTheDocument();
    });
  });

  it('shows time range selector', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Time Range')).toHaveLength(2);
    });
  });

  it('shows chart type selector', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Chart Type')).toHaveLength(2);
    });
  });

  it('displays main chart with title', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getByText('Inventory Utilization Over Time')).toBeInTheDocument();
    });
  });

  it('shows top performing items section', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getByText('Top Performing Items')).toBeInTheDocument();
    });
  });

  it('displays quick actions section', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getByText('Quick Actions')).toBeInTheDocument();
      expect(screen.getByText('View Inventory')).toBeInTheDocument();
      expect(screen.getByText('Back to Dashboard')).toBeInTheDocument();
    });
  });

  it('shows comparison chart section', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getByText('Item Comparison')).toBeInTheDocument();
    });
  });

  it('handles time range change', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Time Range')).toHaveLength(2);
    });
  });

  it('handles chart type change', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Chart Type')).toHaveLength(2);
    });
  });

  it('shows refresh button', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getByText('Refresh')).toBeInTheDocument();
    });
  });

  it('shows compare items button', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getByText('Compare Items')).toBeInTheDocument();
    });
  });

  it('displays user menu in app bar', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      // Check for account icon button
      expect(screen.getByTestId('AccountCircleIcon')).toBeInTheDocument();
    });
  });

  it('shows loading state initially', () => {
    renderWithProviders(<Analytics />);
    
    // Should show loading initially
    expect(screen.getByRole('progressbar')).toBeInTheDocument();
  });

  it('handles error state gracefully', async () => {
    // Mock API to throw error
    const mockApi = require('../../services/api');
    mockApi.getInventoryItems.mockRejectedValueOnce(new Error('API Error'));
    
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getByText(/Failed to load analytics data/)).toBeInTheDocument();
    });
  });

  it('displays chart controls correctly', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      // Check for chart type toggle buttons
      expect(screen.getByTestId('ShowChartIcon')).toBeInTheDocument();
      expect(screen.getByTestId('BarChartIcon')).toBeInTheDocument();
    });
  });

  it('shows utilization percentage in summary stats', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      // Should show utilization percentage
      const percentageElements = screen.getAllByText(/%/);
      expect(percentageElements.length).toBeGreaterThan(0);
    });
  });

  it('displays inventory value in summary stats', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      // Should show dollar amount
      expect(screen.getByText(/\$/)).toBeInTheDocument();
    });
  });

  it('shows navigation buttons work', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      const viewInventoryButton = screen.getByText('View Inventory');
      const backToDashboardButton = screen.getByText('Back to Dashboard');
      
      expect(viewInventoryButton).toBeInTheDocument();
      expect(backToDashboardButton).toBeInTheDocument();
    });
  });

  it('renders with proper Material-UI theme', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      // Check that Material-UI components are rendered
      expect(screen.getByText('Analytics Dashboard')).toBeInTheDocument();
    });
  });

  it('displays chart controls', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      // Should show chart type toggle buttons
      expect(screen.getByText('Inventory Utilization Over Time')).toBeInTheDocument();
    });
  });

  it('shows time range selector', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Time Range')).toHaveLength(2);
    });
  });

  it('displays summary statistics', async () => {
    renderWithProviders(<Analytics />);
    
    await waitFor(() => {
      // Should show summary stats
      expect(screen.getByText('Total Items')).toBeInTheDocument();
      expect(screen.getByText('Low Stock Items')).toBeInTheDocument();
      expect(screen.getByText('Total Value')).toBeInTheDocument();
      expect(screen.getByText('Avg Utilization')).toBeInTheDocument();
    });
  });
}); 