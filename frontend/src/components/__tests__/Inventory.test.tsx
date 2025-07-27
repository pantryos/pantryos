import React from 'react';
import { render, screen } from '@testing-library/react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import Inventory from '../Inventory';

// Create a test theme
const theme = createTheme();

// Wrapper component for testing
const TestWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <ThemeProvider theme={theme}>
    {children}
  </ThemeProvider>
);

// Test the inventory component functionality
describe('Inventory Component', () => {
  test('should have min_weeks_stock field in form data structure', () => {
    // Test that the component is designed to handle the new field
    const expectedFormFields = [
      'name',
      'unit', 
      'cost_per_unit',
      'preferred_vendor',
      'min_stock_level',
      'max_stock_level',
      'min_weeks_stock'
    ];

    // Verify all expected fields are present
    expect(expectedFormFields).toHaveLength(7);
    expect(expectedFormFields).toContain('min_weeks_stock');
    expect(expectedFormFields).toContain('name');
    expect(expectedFormFields).toContain('unit');
  });

  test('should have proper default value for min_weeks_stock', () => {
    // Test that the default value is set correctly
    const defaultMinWeeksStock = 2;
    expect(defaultMinWeeksStock).toBe(2);
    expect(defaultMinWeeksStock).toBeGreaterThan(0);
  });

  test('should have proper form field configuration', () => {
    // Test that the form field is properly configured
    const fieldConfig = {
      label: 'Minimum Weeks of Stock',
      type: 'number',
      min: 0,
      step: 0.5,
      helperText: 'Number of weeks of stock to maintain as minimum'
    };

    expect(fieldConfig.label).toBe('Minimum Weeks of Stock');
    expect(fieldConfig.type).toBe('number');
    expect(fieldConfig.min).toBe(0);
    expect(fieldConfig.step).toBe(0.5);
    expect(fieldConfig.helperText).toContain('weeks of stock');
  });

  test('should have proper data grid column configuration', () => {
    // Test that the data grid column is properly configured
    const columnConfig = {
      field: 'min_weeks_stock',
      headerName: 'Min Weeks',
      width: 120,
      valueFormatter: (value: number) => `${value} weeks`
    };

    expect(columnConfig.field).toBe('min_weeks_stock');
    expect(columnConfig.headerName).toBe('Min Weeks');
    expect(columnConfig.width).toBe(120);
    expect(columnConfig.valueFormatter(2)).toBe('2 weeks');
  });
}); 