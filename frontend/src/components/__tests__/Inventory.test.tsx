import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import Inventory from '../Inventory';
import apiService from '../../services/api';
import { InventoryItem } from '../../types/api';

// Mock the API service
jest.mock('../../services/api');
const mockedApiService = apiService as jest.Mocked<typeof apiService>;



// Mock ItemUsageAnalytics component
jest.mock('../ItemUsageAnalytics', () => {
  return function MockItemUsageAnalytics({ item, onClose }: { item: InventoryItem; onClose?: () => void }) {
    return (
      <div data-testid="item-usage-analytics">
        <h3>Analytics for {item.name}</h3>
        <button onClick={onClose}>Close Analytics</button>
      </div>
    );
  };
});

const theme = createTheme();

// Wrapper component to provide necessary context
const TestWrapper = ({ children }: { children: React.ReactNode }) => (
  <BrowserRouter>
    <ThemeProvider theme={theme}>
      {children}
    </ThemeProvider>
  </BrowserRouter>
);

describe('Inventory Component', () => {
  const mockItems: InventoryItem[] = [
    {
      id: 1,
      account_id: 1,
      name: 'Test Item 1',
      unit: 'kg',
      cost_per_unit: 10.50,
      preferred_vendor: 'Vendor A',
      min_stock_level: 10,
      max_stock_level: 100,
    },
    {
      id: 2,
      account_id: 1,
      name: 'Test Item 2',
      unit: 'pieces',
      cost_per_unit: 5.25,
      preferred_vendor: 'Vendor B',
      min_stock_level: 5,
      max_stock_level: 50,
    },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    mockedApiService.getInventoryItems.mockResolvedValue(mockItems);
    mockedApiService.createInventoryItem.mockResolvedValue(mockItems[0]);
    mockedApiService.updateInventoryItem.mockResolvedValue(mockItems[0]);
    mockedApiService.deleteInventoryItem.mockResolvedValue(undefined);
  });

  describe('rendering', () => {
    it('should render inventory management page', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByText('Inventory Management')).toBeInTheDocument();
        expect(screen.getByText('Inventory Items')).toBeInTheDocument();
        expect(screen.getByText('Add Item')).toBeInTheDocument();
      });
    });

    it('should display inventory items in data grid', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByText('Test Item 1')).toBeInTheDocument();
        expect(screen.getByText('Test Item 2')).toBeInTheDocument();
        expect(screen.getByText('Vendor A')).toBeInTheDocument();
        expect(screen.getByText('Vendor B')).toBeInTheDocument();
      });
    });

    it('should show loading state initially', () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      // Loading state should be shown initially
      expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });
  });

  describe('analytics functionality', () => {
    it('should display analytics button for each item', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        // Check for analytics buttons (should be 2 items)
        const analyticsButtons = screen.getAllByLabelText('Analytics');
        expect(analyticsButtons).toHaveLength(2);
      });
    });

    it('should open analytics dialog when analytics button is clicked', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        const analyticsButtons = screen.getAllByLabelText('Analytics');
        fireEvent.click(analyticsButtons[0]);
      });

      // Check that analytics dialog opens
      expect(screen.getByText('Test Item 1 - Usage Analytics')).toBeInTheDocument();
      expect(screen.getByTestId('item-usage-analytics')).toBeInTheDocument();
    });

    it('should close analytics dialog when close button is clicked', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        const analyticsButtons = screen.getAllByLabelText('Analytics');
        fireEvent.click(analyticsButtons[0]);
      });

      // Verify dialog is open
      expect(screen.getByText('Test Item 1 - Usage Analytics')).toBeInTheDocument();

      // Click close button
      const closeButton = screen.getByRole('button', { name: /close/i });
      fireEvent.click(closeButton);

      // Verify dialog is closed
      await waitFor(() => {
        expect(screen.queryByText('Test Item 1 - Usage Analytics')).not.toBeInTheDocument();
      });
    });

    it('should pass correct item data to analytics component', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        const analyticsButtons = screen.getAllByLabelText('Analytics');
        fireEvent.click(analyticsButtons[0]);
      });

      // Check that the correct item name is displayed in analytics
      expect(screen.getByText('Analytics for Test Item 1')).toBeInTheDocument();
    });

    it('should close analytics dialog when analytics component calls onClose', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        const analyticsButtons = screen.getAllByLabelText('Analytics');
        fireEvent.click(analyticsButtons[0]);
      });

      // Verify dialog is open
      expect(screen.getByText('Test Item 1 - Usage Analytics')).toBeInTheDocument();

      // Click the close button inside the analytics component
      const analyticsCloseButton = screen.getByText('Close Analytics');
      fireEvent.click(analyticsCloseButton);

      // Verify dialog is closed
      await waitFor(() => {
        expect(screen.queryByText('Test Item 1 - Usage Analytics')).not.toBeInTheDocument();
      });
    });
  });

  describe('search functionality', () => {
    it('should filter items by name', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByText('Test Item 1')).toBeInTheDocument();
        expect(screen.getByText('Test Item 2')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('Search by name or vendor...');
      fireEvent.change(searchInput, { target: { value: 'Test Item 1' } });

      await waitFor(() => {
        expect(screen.getByText('Test Item 1')).toBeInTheDocument();
        expect(screen.queryByText('Test Item 2')).not.toBeInTheDocument();
      });
    });

    it('should filter items by vendor', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByText('Vendor A')).toBeInTheDocument();
        expect(screen.getByText('Vendor B')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('Search by name or vendor...');
      fireEvent.change(searchInput, { target: { value: 'Vendor A' } });

      await waitFor(() => {
        expect(screen.getByText('Vendor A')).toBeInTheDocument();
        expect(screen.queryByText('Vendor B')).not.toBeInTheDocument();
      });
    });
  });

  describe('error handling', () => {
    it('should display error message when API call fails', async () => {
      mockedApiService.getInventoryItems.mockRejectedValue(new Error('API Error'));

      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByText('Failed to load inventory items')).toBeInTheDocument();
      });
    });

    it('should clear error when error alert is closed', async () => {
      mockedApiService.getInventoryItems.mockRejectedValue(new Error('API Error'));

      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByText('Failed to load inventory items')).toBeInTheDocument();
      });

      const closeButton = screen.getByRole('button', { name: /close/i });
      fireEvent.click(closeButton);

      await waitFor(() => {
        expect(screen.queryByText('Failed to load inventory items')).not.toBeInTheDocument();
      });
    });
  });

  describe('CRUD operations', () => {
    it('should open add item dialog when Add Item button is clicked', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        const addButton = screen.getByText('Add Item');
        fireEvent.click(addButton);
      });

      expect(screen.getByText('Add New Inventory Item')).toBeInTheDocument();
    });

    it('should open edit dialog when edit button is clicked', async () => {
      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        const editButtons = screen.getAllByLabelText('Edit');
        fireEvent.click(editButtons[0]);
      });

      expect(screen.getByText('Edit Inventory Item')).toBeInTheDocument();
    });

    it('should show delete confirmation when delete button is clicked', async () => {
      // Mock window.confirm
      const mockConfirm = jest.spyOn(window, 'confirm').mockReturnValue(true);

      render(
        <TestWrapper>
          <Inventory />
        </TestWrapper>
      );

      await waitFor(() => {
        const deleteButtons = screen.getAllByLabelText('Delete');
        fireEvent.click(deleteButtons[0]);
      });

      expect(mockConfirm).toHaveBeenCalledWith('Are you sure you want to delete this item?');
      mockConfirm.mockRestore();
    });
  });
}); 