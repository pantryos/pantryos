import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import Menu from '../Menu'; // Adjust the import path as needed
import apiService from '../../services/api';
import { Category, MenuItem } from '../../types/api';

// Mock the apiService to avoid actual network calls
jest.mock('../../services/api');
const mockApiService = apiService as jest.Mocked<typeof apiService>;

// Mock react-router-dom's useNavigate hook
const mockedNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
    useNavigate: () => mockedNavigate,
}));

// Mock data for our tests
const mockCategories: Category[] = [
    { id: 1, account_id: 1, name: 'Hot Drinks', description: '', color: '', is_active: true, created_at: '', updated_at: '' },
    { id: 2, account_id: 1, name: 'Pastries', description: '', color: '', is_active: true, created_at: '', updated_at: '' },
];

const mockMenuItems: MenuItem[] = [
    { id: 1, account_id: 1, name: 'Espresso', price: 3.50, category: 'Hot Drinks', category_id: 1 },
    { id: 2, account_id: 1, name: 'Croissant', price: 2.75, category: 'Pastries', category_id: 2 },
];


describe('Menu Component', () => {

    // Before each test, reset mocks and set default successful responses
    beforeEach(() => {
        jest.clearAllMocks();
        mockApiService.getMenuItems.mockResolvedValue(mockMenuItems);
        mockApiService.getActiveCategories.mockResolvedValue(mockCategories);
        mockApiService.createMenuItem.mockResolvedValue({ ...mockMenuItems[0], id: 3, name: 'New Item' });
        // mockApiService.deleteMenuItem.mockResolvedValue();
    });

    // Test 1: Initial Render and Data Display
    test('renders loading state initially and then displays menu items', async () => {
        render(<Menu />);
        // Initially, nothing is shown because it's loading (or a spinner, which is harder to test)
        // We use `findBy` to wait for the async data loading to complete
        expect(await screen.findByText('Espresso')).toBeInTheDocument();
        expect(screen.getByText('Croissant')).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /add menu item/i })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /Manage Categories/i })).toBeInTheDocument();

    });
    // Test 2: Adding a New Menu Item
    test('should open a dialog, allow user input, and call the create API on submit', async () => {
        // Arrange: Mock the necessary API calls
        mockApiService.getMenuItems.mockResolvedValue(mockMenuItems);
        mockApiService.getActiveCategories.mockResolvedValue(mockCategories);
        // Mock the create function to simulate a successful creation
        mockApiService.createMenuItem.mockResolvedValue({
            id: 103,
            name: 'New Latte',
            price: 4.25,
            category: 'Hot Drinks',
            account_id: 1,
            category_id: 1,
        });

        // Act
        render(<Menu />);

        // 1. Wait for the page to load fully
        await screen.findByText('Espresso');

        // 2. Click the "Add Menu Item" button to open the dialog
        fireEvent.click(screen.getByRole('button', { name: /add menu item/i }));

        // 3. Wait for the dialog to appear and fill out the form
        await screen.findByRole('heading', { name: /add new menu item/i });
        fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'New Latte' } });
        fireEvent.change(screen.getByLabelText(/price/i), { target: { value: 4.25 } });

        // 4. Select a category from the MUI Select component
        const categorySelect = screen.getByLabelText(/category/i);
        fireEvent.mouseDown(categorySelect); // Opens the dropdown
        fireEvent.click(await screen.findByRole('option', { name: 'Hot Drinks' }));

        // 5. Click the "Create" button to submit the form
        fireEvent.click(screen.getByRole('button', { name: /create/i }));

        // Assert: Verify the API was called with the correct data from the form
        await waitFor(() => {
            expect(mockApiService.createMenuItem).toHaveBeenCalledTimes(1);
            expect(mockApiService.createMenuItem).toHaveBeenCalledWith({
                name: 'New Latte',
                price: '4.25',
                category: 'Hot Drinks',
                category_id: 1,
            });
        });
    });


    test('displays an error message if fetching data fails', async () => {
        // Arrange:
        // 1. Temporarily silence console.error for this test
        const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation(() => { });

        // 2. Override the mock to simulate a network error
        mockApiService.getMenuItems.mockRejectedValue(new Error('Network Error'));

        // Act:
        render(<Menu />);

        // Assert:
        expect(await screen.findByText('Failed to load menu data')).toBeInTheDocument();

        // Cleanup: Restore the original console.error function
        consoleErrorSpy.mockRestore();

    });


    test('should filter menu items based on the search term', async () => {
        // Arrange: Set up the API to return our mock data
        mockApiService.getMenuItems.mockResolvedValue(mockMenuItems);
        mockApiService.getActiveCategories.mockResolvedValue(mockCategories);

        // Act: Render the component
        render(<Menu />);

        // 1. Wait for the initial data to load
        await screen.findByText('Espresso');

        // Assert (Initial State): Both items should be visible before searching
        expect(screen.getByText('Espresso')).toBeInTheDocument();
        expect(screen.getByText('Croissant')).toBeInTheDocument();

        // 2. Find the search input and simulate typing
        const searchInput = screen.getByPlaceholderText(/search by name or category/i);
        fireEvent.change(searchInput, { target: { value: 'Croissant' } });

        // Assert (After Search): Only the searched item should be visible
        expect(screen.getByText('Croissant')).toBeInTheDocument();
        // Use `queryByText` to assert an element is NOT present
        expect(screen.queryByText('Espresso')).not.toBeInTheDocument();

        // 3. Simulate clearing the search to see if all items return
        fireEvent.change(searchInput, { target: { value: '' } });

        // Assert (After Clearing): Both items should be visible again
        expect(screen.getByText('Espresso')).toBeInTheDocument();
        expect(screen.getByText('Croissant')).toBeInTheDocument();
    });

    test('should navigate to the categories page when "Manage Categories" is clicked', async () => {
        // Arrange
        mockApiService.getMenuItems.mockResolvedValue(mockMenuItems);
        mockApiService.getActiveCategories.mockResolvedValue(mockCategories);
        render(<Menu />);
        await screen.findByText('Espresso');

        // Act: Click the new button
        fireEvent.click(screen.getByRole('button', { name: /manage categories/i }));

        // Assert: Check if the navigate function was called with the correct path
        expect(mockedNavigate).toHaveBeenCalledWith('/categories');
    });

});