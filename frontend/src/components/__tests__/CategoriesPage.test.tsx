import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom';

import CategoriesPage from '../CategoriesPage';
import apiService from '../../services/api';
import { Category } from '../../types/api';

// Mock ResizeObserver for Material-UI components
global.ResizeObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

// Also mock it on window object as a fallback
Object.defineProperty(window, 'ResizeObserver', {
  writable: true,
  value: jest.fn().mockImplementation(() => ({
    observe: jest.fn(),
    unobserve: jest.fn(),
    disconnect: jest.fn(),
  }))
});

// Mock the specific TextareaAutosize component that's causing issues
jest.mock('@mui/material/TextareaAutosize', () => {
  return function MockTextareaAutosize(props: any) {
    // Filter out MUI-specific props that aren't valid for textarea
    const { minRows, maxRows, cacheMeasurements, ...textareaProps } = props;
    return <textarea {...textareaProps} />;
  };
});

// Mock MUI DataGrid to avoid complex rendering issues in tests
jest.mock('@mui/x-data-grid', () => ({
  DataGrid: ({ rows, columns, loading }: any) => {
    if (loading) return <div data-testid="loading">Loading...</div>;
    
    return (
      <div data-testid="data-grid">
        <div data-testid="data-grid-header">
          {columns.map((col: any) => (
            <span key={col.field}>{col.headerName}</span>
          ))}
        </div>
        <div data-testid="data-grid-rows">
          {rows.map((row: any) => (
            <div key={row.id} data-testid={`data-grid-row-${row.id}`}>
              {columns.map((col: any) => (
                <span key={col.field} data-testid={`cell-${col.field}`}>
                  {col.field === 'actions' ? (
                    <div>
                      <button 
                        aria-label="Edit" 
                        onClick={() => col.getActions({ row }).find((action: any) => action.props.label === 'Edit')?.props.onClick()}
                      >
                        Edit
                      </button>
                      <button 
                        aria-label="Delete" 
                        onClick={() => col.getActions({ row }).find((action: any) => action.props.label === 'Delete')?.props.onClick()}
                      >
                        Delete
                      </button>
                    </div>
                  ) : col.renderCell ? (
                    col.renderCell({ value: row[col.field] })
                  ) : (
                    row[col.field]
                  )}
                </span>
              ))}
            </div>
          ))}
        </div>
      </div>
    );
  },
  GridColDef: {},
  GridActionsCellItem: ({ icon, label, onClick }: any) => ({ props: { icon, label, onClick } }),
}));

// Mocks and mock data
jest.mock('../../services/api');
const mockApiService = apiService as jest.Mocked<typeof apiService>;

const mockCategories: Category[] = [
  { id: 1, account_id: 1, name: 'Supplies', description: 'Non-food items.', color: '#90A4AE', is_active: true, created_at: '', updated_at: '' },
  { id: 2, account_id: 1, name: 'Drinks', description: 'Coffee, tea, juices.', color: '#03A9F4', is_active: false, created_at: '', updated_at: '' },
];

describe('CategoriesPage Component', () => {

  beforeEach(() => {
    jest.clearAllMocks();
    mockApiService.getCategories.mockResolvedValue([...mockCategories]);
    mockApiService.createCategory.mockResolvedValue({ ...mockCategories[0], id: 3, name: 'New Category' });
    mockApiService.deleteCategory.mockResolvedValue();
  });

  test('renders the list of categories after fetching data', async () => {
    render(<CategoriesPage />);
    expect(await screen.findByText('Supplies')).toBeInTheDocument();
    expect(screen.getByText('Drinks')).toBeInTheDocument();
  });

  test('filters categories when a user types in the search bar', async () => {
    // Arrange
    render(<CategoriesPage />);
    expect(await screen.findByText('Supplies')).toBeInTheDocument();

    // Act: Type using userEvent v13 syntax (no setup needed)
    const searchInput = screen.getByPlaceholderText(/search by name or description/i);
    userEvent.type(searchInput, 'Supplies');

    // Assert: Check the filtered results
    await waitFor(() => {
        expect(screen.getByText('Supplies')).toBeInTheDocument();
    });
    expect(screen.queryByText('Drinks')).not.toBeInTheDocument();
  });

  test('allows a user to add a new category via the modal', async () => {
    // Arrange
    render(<CategoriesPage />);
    await screen.findByText('Supplies');

    // Act
    fireEvent.click(screen.getByRole('button', { name: /add category/i }));
    const nameInput = await screen.findByLabelText(/category name/i);
    userEvent.type(nameInput, 'Food'); // userEvent v13 syntax
    
    fireEvent.click(screen.getByRole('button', { name: /create/i }));

    // Assert
    await waitFor(() => {
      expect(mockApiService.createCategory).toHaveBeenCalledWith(
        expect.objectContaining({ name: 'Food' })
      );
    });
    expect(await screen.findByText('Category created successfully!')).toBeInTheDocument();
  });

  test('allows a user to delete a category after confirmation', async () => {
    // Arrange
    render(<CategoriesPage />);
    await screen.findByText('Supplies');

    // Act - Find delete button using aria-label
    const deleteButton = screen.getAllByLabelText(/delete/i)[0];
    fireEvent.click(deleteButton);
    
    // Find and click the confirm delete button in the dialog
    const confirmButton = await screen.findByRole('button', { name: /delete/i });
    fireEvent.click(confirmButton);

    // Assert
    await waitFor(() => {
      expect(mockApiService.deleteCategory).toHaveBeenCalledWith(mockCategories[0].id);
    });
    expect(await screen.findByText('Category deleted successfully!')).toBeInTheDocument();
  });
});