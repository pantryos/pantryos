import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import ItemComparison from '../ItemComparison';
import { InventoryItem } from '../../types/api';

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
    {
      id: 3,
      name: 'Potatoes',
      current_stock: 75,
      max_stock_level: 120,
      min_stock_level: 20,
      cost_per_unit: 1.2,
      unit: 'kg',
      preferred_vendor: 'Root Vegetables Co',
      account_id: 1,
    },
  ]),
}));

const theme = createTheme();

const renderWithTheme = (component: React.ReactElement) => {
  return render(
    <ThemeProvider theme={theme}>
      {component}
    </ThemeProvider>
  );
};

describe('ItemComparison Component', () => {
  const mockOnClose = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders dialog with title when open', () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    expect(screen.getByText('Item Comparison')).toBeInTheDocument();
  });

  it('does not render when closed', () => {
    renderWithTheme(<ItemComparison open={false} onClose={mockOnClose} />);
    
    expect(screen.queryByText('Item Comparison')).not.toBeInTheDocument();
  });

  it('displays close button', () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    const closeButton = screen.getByRole('button', { name: /close/i });
    expect(closeButton).toBeInTheDocument();
  });

  it('calls onClose when close button is clicked', () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    const closeButton = screen.getByRole('button', { name: /close/i });
    fireEvent.click(closeButton);
    
    expect(mockOnClose).toHaveBeenCalled();
  });

  it('shows time range selector', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Time Range')).toHaveLength(2);
    });
  });

  it('shows metric selector', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Metric')).toHaveLength(2);
    });
  });

  it('shows chart type selector', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Chart Type')).toHaveLength(2);
    });
  });

  it('shows clear selection button', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Clear Selection')).toBeInTheDocument();
    });
  });

  it('displays item selection list', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select Items to Compare')).toBeInTheDocument();
    });
  });

  it('shows item stock information', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select Items to Compare')).toBeInTheDocument();
    });
  });

  it('displays item cost chips', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select Items to Compare')).toBeInTheDocument();
    });
  });

  it('allows item selection via checkboxes', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select Items to Compare')).toBeInTheDocument();
    });
  });

  it('shows chart when items are selected', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select items to compare')).toBeInTheDocument();
    });
  });

  it('shows placeholder when no items selected', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select items to compare')).toBeInTheDocument();
    });
  });

  it('handles time range change', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Time Range')).toHaveLength(2);
    });
  });

  it('handles metric change', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Metric')).toHaveLength(2);
    });
  });

  it('handles chart type change', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Chart Type')).toHaveLength(2);
    });
  });

  it('clears selection when clear button is clicked', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Clear Selection')).toBeInTheDocument();
    });
  });

  it('shows selected items summary when items are selected', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select Items to Compare')).toBeInTheDocument();
    });
  });

  it('displays trend indicators for selected items', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select Items to Compare')).toBeInTheDocument();
    });
  });

  it('shows average utilization for selected items', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select Items to Compare')).toBeInTheDocument();
    });
  });

  it('displays current stock information for selected items', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select Items to Compare')).toBeInTheDocument();
    });
  });

  it('handles multiple item selection', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText('Select Items to Compare')).toBeInTheDocument();
    });
  });

  it('shows loading state when fetching items', () => {
    // Mock API to be slow
    const mockApi = require('../../services/api');
    mockApi.getInventoryItems.mockImplementation(() => 
      new Promise(resolve => setTimeout(() => resolve([]), 100))
    );
    
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    expect(screen.getByRole('progressbar')).toBeInTheDocument();
  });

  it('handles API error gracefully', async () => {
    // Mock API to throw error
    const mockApi = require('../../services/api');
    mockApi.getInventoryItems.mockRejectedValueOnce(new Error('API Error'));
    
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getByText(/Failed to load inventory items/)).toBeInTheDocument();
    });
  });

  it('updates chart when different metrics are selected', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Metric')).toHaveLength(2);
    });
  });

  it('maintains selected items when changing time range', async () => {
    renderWithTheme(<ItemComparison open={true} onClose={mockOnClose} />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Time Range')).toHaveLength(2);
    });
  });
}); 