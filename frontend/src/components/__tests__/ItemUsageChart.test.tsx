import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import ItemUsageChart, { UsageData } from '../ItemUsageChart';

const theme = createTheme();

// Wrapper component to provide necessary context
const TestWrapper = ({ children }: { children: React.ReactNode }) => (
  <ThemeProvider theme={theme}>
    {children}
  </ThemeProvider>
);

describe('ItemUsageChart Component', () => {
  const mockData: UsageData[] = [
    { date: 'Jan 1', quantity: 10, cost: 25.50 },
    { date: 'Jan 2', quantity: 15, cost: 38.25 },
    { date: 'Jan 3', quantity: 8, cost: 20.40 },
    { date: 'Jan 4', quantity: 12, cost: 30.60 },
    { date: 'Jan 5', quantity: 20, cost: 51.00 },
  ];

  const defaultProps = {
    title: 'Test Chart',
    data: mockData,
    itemName: 'Test Item',
    currentStock: 50,
    maxStock: 100,
    minStock: 10,
    timeRange: '30d' as const,
    onTimeRangeChange: jest.fn(),
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('rendering', () => {
    it('should render chart with title and item name', () => {
      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} />
        </TestWrapper>
      );

      expect(screen.getByText('Test Chart')).toBeInTheDocument();
      expect(screen.getByText('Test Item')).toBeInTheDocument();
    });

    it('should display time range selector', () => {
      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} />
        </TestWrapper>
      );

      // Use getAllByText to handle multiple elements with the same text
      const timeRangeElements = screen.getAllByText('Time Range');
      expect(timeRangeElements.length).toBeGreaterThan(0);
      expect(screen.getByText('Last 30 Days')).toBeInTheDocument();
    });

    it('should display summary statistics cards', () => {
      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} />
        </TestWrapper>
      );

      expect(screen.getByText('Current Stock')).toBeInTheDocument();
      expect(screen.getByText('Avg Usage')).toBeInTheDocument();
      expect(screen.getByText('Total Cost')).toBeInTheDocument();
    });

    it('should display stock level information', () => {
      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} />
        </TestWrapper>
      );

      expect(screen.getByText('50 / 100')).toBeInTheDocument();
      expect(screen.getByText('Normal (50.0%)')).toBeInTheDocument();
    });

    it('should display usage statistics', () => {
      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} />
        </TestWrapper>
      );

      // Total usage: 10 + 15 + 8 + 12 + 20 = 65
      // Average usage: 65 / 5 = 13
      expect(screen.getByText('13.0 / day')).toBeInTheDocument();
      expect(screen.getByText('Total: 65')).toBeInTheDocument();
    });

    it('should display cost statistics', () => {
      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} />
        </TestWrapper>
      );

      // Total cost: 25.50 + 38.25 + 20.40 + 30.60 + 51.00 = 165.75
      expect(screen.getByText('$165.75')).toBeInTheDocument();
      expect(screen.getByText('Avg: $33.15')).toBeInTheDocument();
    });

    it('should display chart placeholder in test environment', () => {
      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} />
        </TestWrapper>
      );

      expect(screen.getByText('Chart visualization (disabled in test environment)')).toBeInTheDocument();
    });
  });

  describe('interactions', () => {
    it('should handle time range change', () => {
      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} />
        </TestWrapper>
      );

      const timeRangeSelect = screen.getByRole('combobox');
      fireEvent.mouseDown(timeRangeSelect);

      const sevenDaysOption = screen.getByText('Last 7 Days');
      fireEvent.click(sevenDaysOption);

      expect(defaultProps.onTimeRangeChange).toHaveBeenCalledWith('7d');
    });

    it('should handle all time range options', () => {
      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} />
        </TestWrapper>
      );

      const timeRangeSelect = screen.getByRole('combobox');
      
      // Test 7 days
      fireEvent.mouseDown(timeRangeSelect);
      fireEvent.click(screen.getByText('Last 7 Days'));
      expect(defaultProps.onTimeRangeChange).toHaveBeenCalledWith('7d');

      // Test 90 days
      fireEvent.mouseDown(timeRangeSelect);
      fireEvent.click(screen.getByText('Last 90 Days'));
      expect(defaultProps.onTimeRangeChange).toHaveBeenCalledWith('90d');
    });
  });

  describe('stock status calculations', () => {
    it('should show low stock status when current stock is at or below minimum', () => {
      const lowStockProps = {
        ...defaultProps,
        currentStock: 5,
        minStock: 10,
      };

      render(
        <TestWrapper>
          <ItemUsageChart {...lowStockProps} />
        </TestWrapper>
      );

      expect(screen.getByText('5 / 100')).toBeInTheDocument();
      expect(screen.getByText('Low Stock (5.0%)')).toBeInTheDocument();
    });

    it('should show high stock status when current stock is 80% or more of max', () => {
      const highStockProps = {
        ...defaultProps,
        currentStock: 85,
        maxStock: 100,
      };

      render(
        <TestWrapper>
          <ItemUsageChart {...highStockProps} />
        </TestWrapper>
      );

      expect(screen.getByText('85 / 100')).toBeInTheDocument();
      expect(screen.getByText('High Stock (85.0%)')).toBeInTheDocument();
    });

    it('should show normal stock status when current stock is between min and 80%', () => {
      const normalStockProps = {
        ...defaultProps,
        currentStock: 50,
        maxStock: 100,
        minStock: 10,
      };

      render(
        <TestWrapper>
          <ItemUsageChart {...normalStockProps} />
        </TestWrapper>
      );

      expect(screen.getByText('50 / 100')).toBeInTheDocument();
      expect(screen.getByText('Normal (50.0%)')).toBeInTheDocument();
    });
  });

  describe('trend calculations', () => {
    it('should show stable trend for consistent data', () => {
      const stableData: UsageData[] = [
        { date: 'Jan 1', quantity: 10, cost: 25 },
        { date: 'Jan 2', quantity: 10, cost: 25 },
        { date: 'Jan 3', quantity: 10, cost: 25 },
        { date: 'Jan 4', quantity: 10, cost: 25 },
        { date: 'Jan 5', quantity: 10, cost: 25 },
        { date: 'Jan 6', quantity: 10, cost: 25 },
      ];

      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} data={stableData} />
        </TestWrapper>
      );

      expect(screen.getByText('Stable')).toBeInTheDocument();
    });

    it('should show increasing trend for rising data', () => {
      const increasingData: UsageData[] = [
        { date: 'Jan 1', quantity: 5, cost: 12.5 },
        { date: 'Jan 2', quantity: 5, cost: 12.5 },
        { date: 'Jan 3', quantity: 5, cost: 12.5 },
        { date: 'Jan 4', quantity: 15, cost: 37.5 },
        { date: 'Jan 5', quantity: 15, cost: 37.5 },
        { date: 'Jan 6', quantity: 15, cost: 37.5 },
      ];

      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} data={increasingData} />
        </TestWrapper>
      );

      expect(screen.getByText('Increasing')).toBeInTheDocument();
    });

    it('should show decreasing trend for falling data', () => {
      const decreasingData: UsageData[] = [
        { date: 'Jan 1', quantity: 15, cost: 37.5 },
        { date: 'Jan 2', quantity: 15, cost: 37.5 },
        { date: 'Jan 3', quantity: 15, cost: 37.5 },
        { date: 'Jan 4', quantity: 5, cost: 12.5 },
        { date: 'Jan 5', quantity: 5, cost: 12.5 },
        { date: 'Jan 6', quantity: 5, cost: 12.5 },
      ];

      render(
        <TestWrapper>
          <ItemUsageChart {...defaultProps} data={decreasingData} />
        </TestWrapper>
      );

      expect(screen.getByText('Decreasing')).toBeInTheDocument();
    });
  });

  describe('edge cases', () => {
    it('should handle empty data array', () => {
      const emptyDataProps = {
        ...defaultProps,
        data: [],
      };

      render(
        <TestWrapper>
          <ItemUsageChart {...emptyDataProps} />
        </TestWrapper>
      );

      expect(screen.getByText('0.0 / day')).toBeInTheDocument();
      expect(screen.getByText('Total: 0')).toBeInTheDocument();
      expect(screen.getByText('$0.00')).toBeInTheDocument();
      expect(screen.getByText('No trend')).toBeInTheDocument();
    });

    it('should handle single data point', () => {
      const singleDataProps = {
        ...defaultProps,
        data: [{ date: 'Jan 1', quantity: 10, cost: 25.50 }],
      };

      render(
        <TestWrapper>
          <ItemUsageChart {...singleDataProps} />
        </TestWrapper>
      );

      expect(screen.getByText('10.0 / day')).toBeInTheDocument();
      expect(screen.getByText('Total: 10')).toBeInTheDocument();
      expect(screen.getByText('$25.50')).toBeInTheDocument();
      expect(screen.getByText('No trend')).toBeInTheDocument();
    });

    it('should handle zero stock levels', () => {
      const zeroStockProps = {
        ...defaultProps,
        currentStock: 0,
        maxStock: 0,
        minStock: 0,
      };

      render(
        <TestWrapper>
          <ItemUsageChart {...zeroStockProps} />
        </TestWrapper>
      );

      expect(screen.getByText('0 / 0')).toBeInTheDocument();
      expect(screen.getByText('Invalid Max Stock (0.0%)')).toBeInTheDocument();
    });
  });
}); 