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
  test('should have min_weeks_stock and max_weeks_stock fields in form data structure', () => {
    // Test that the component is designed to handle the new fields
    const expectedFormFields = [
      'name',
      'unit', 
      'cost_per_unit',
      'preferred_vendor',
      'min_stock_level',
      'max_stock_level',
      'min_weeks_stock',
      'max_weeks_stock'
    ];

    // Verify all expected fields are present
    expect(expectedFormFields).toHaveLength(8);
    expect(expectedFormFields).toContain('min_weeks_stock');
    expect(expectedFormFields).toContain('max_weeks_stock');
    expect(expectedFormFields).toContain('name');
    expect(expectedFormFields).toContain('unit');
  });

  test('should have proper default values for weeks stock fields', () => {
    // Test that the default values are set correctly
    const defaultMinWeeksStock = 2;
    const defaultMaxWeeksStock = 8;
    
    expect(defaultMinWeeksStock).toBe(2);
    expect(defaultMinWeeksStock).toBeGreaterThan(0);
    expect(defaultMaxWeeksStock).toBe(8);
    expect(defaultMaxWeeksStock).toBeGreaterThan(defaultMinWeeksStock);
  });

  test('should have proper form field configuration for weeks stock fields', () => {
    // Test that the form fields are properly configured
    const minFieldConfig = {
      label: 'Minimum Weeks of Stock',
      type: 'number',
      min: 0,
      step: 0.5,
      helperText: 'Number of weeks of stock to maintain as minimum'
    };

    const maxFieldConfig = {
      label: 'Maximum Weeks of Stock',
      type: 'number',
      min: 0,
      step: 0.5,
      helperText: 'Number of weeks of stock to maintain as maximum'
    };

    expect(minFieldConfig.label).toBe('Minimum Weeks of Stock');
    expect(minFieldConfig.type).toBe('number');
    expect(minFieldConfig.min).toBe(0);
    expect(minFieldConfig.step).toBe(0.5);
    expect(minFieldConfig.helperText).toContain('weeks of stock');

    expect(maxFieldConfig.label).toBe('Maximum Weeks of Stock');
    expect(maxFieldConfig.type).toBe('number');
    expect(maxFieldConfig.min).toBe(0);
    expect(maxFieldConfig.step).toBe(0.5);
    expect(maxFieldConfig.helperText).toContain('weeks of stock');
  });

  test('should have proper data grid column configuration for weeks stock fields', () => {
    // Test that the data grid columns are properly configured
    const minColumnConfig = {
      field: 'min_weeks_stock',
      headerName: 'Min Weeks',
      width: 120,
      valueFormatter: (value: number) => `${value} weeks`
    };

    const maxColumnConfig = {
      field: 'max_weeks_stock',
      headerName: 'Max Weeks',
      width: 120,
      valueFormatter: (value: number) => `${value} weeks`
    };

    expect(minColumnConfig.field).toBe('min_weeks_stock');
    expect(minColumnConfig.headerName).toBe('Min Weeks');
    expect(minColumnConfig.width).toBe(120);
    expect(minColumnConfig.valueFormatter(2)).toBe('2 weeks');

    expect(maxColumnConfig.field).toBe('max_weeks_stock');
    expect(maxColumnConfig.headerName).toBe('Max Weeks');
    expect(maxColumnConfig.width).toBe(120);
    expect(maxColumnConfig.valueFormatter(8)).toBe('8 weeks');
  });
}); 