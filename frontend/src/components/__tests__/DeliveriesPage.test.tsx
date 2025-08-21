import React from 'react';
import userEvent from '@testing-library/user-event';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { BrowserRouter } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';

import DeliveriesPage from '../DeliveriesPage';
import apiService from '../../services/api';
import { Delivery, InventoryItem } from '../../types/api';

// Mock api service
jest.mock('../../services/api');
const mockApiService = apiService as jest.Mocked<typeof apiService>;

const mockInventoryItems: InventoryItem[] = [
  { id: 1, account_id: 1, name: 'Coffee Beans', unit: 'kg', current_stock: 10, min_stock_level: 5, max_stock_level: 20, cost_per_unit: 10, preferred_vendor: 'A' },
  { id: 2, account_id: 1, name: 'Milk', unit: 'liter', current_stock: 10, min_stock_level: 5, max_stock_level: 20, cost_per_unit: 1, preferred_vendor: 'B' },
] as InventoryItem[];

const mockDeliveries: Delivery[] = [
  { id: 1, inventory_item_id: 1, vendor: 'Mountain Roasters', quantity: 50, cost: 200, delivery_date: '2025-08-20T10:00:00Z', account_id: 1 },
  { id: 2, inventory_item_id: 2, vendor: 'Farm Fresh Dairy', quantity: 100, cost: 150, delivery_date: '2025-08-19T11:00:00Z', account_id: 1 },
];

const theme = createTheme({
  components: {
    MuiButtonBase: {
      defaultProps: {
        disableRipple: true, // Disable the ripple globally
      },
    },
  },
});

const renderWithProviders = (component: React.ReactElement) => {
  return render(
    <ThemeProvider theme={theme}>
      <BrowserRouter>{component}</BrowserRouter>
    </ThemeProvider>
  );
};

describe('DeliveriesPage Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockApiService.getDeliveries.mockResolvedValue([...mockDeliveries]);
    mockApiService.getInventoryItems.mockResolvedValue(mockInventoryItems);
    mockApiService.logDelivery.mockResolvedValue({ ...mockDeliveries[0], id: 3 });
  });

  test('renders loading state and then displays the list of deliveries', async () => {
    renderWithProviders(<DeliveriesPage />);

    expect(await screen.findByText('Mountain Roasters')).toBeInTheDocument();
    expect(screen.getByText('Farm Fresh Dairy')).toBeInTheDocument();
  });

  // test('filters deliveries when a user types in the search bar', async () => {
  //   jest.useFakeTimers();

  //   renderWithProviders(<DeliveriesPage />);

  //   expect(await screen.findByText('Mountain Roasters')).toBeInTheDocument();
  //   expect(await screen.findByText('Farm Fresh Dairy')).toBeInTheDocument();

  //   const searchInput = screen.getByPlaceholderText(/search by vendor/i);
  //   await userEvent.type(searchInput, 'Mountain');

  //   // Manually run timers inside act because it causes a state update
  //   act(() => {
  //     jest.runAllTimers(); // flush debounce
  //   });

  //   expect(await screen.findByText('Mountain Roasters')).toBeInTheDocument();

  //   // Now that the UI has settled, we can check what should NOT be there.
  //   // expect(screen.queryByText('Farm Fresh Dairy')).not.toBeInTheDocument();

  //   jest.useRealTimers()
  // });

  test('allows a user to log a new delivery via the modal', async () => {
    renderWithProviders(<DeliveriesPage />);

    await screen.findByText('Mountain Roasters'); // Ensure page is loaded

    await userEvent.click(screen.getByRole('button', { name: /log new delivery/i }));

    expect(await screen.findByRole('heading', { name: 'Log New Delivery' })).toBeVisible();

    await userEvent.type(screen.getByLabelText(/quantity/i), '25');
    await userEvent.type(screen.getByLabelText(/total cost/i), '99.99');
    await userEvent.click(screen.getByLabelText(/inventory item/i));

    const option = await screen.findByRole('option', { name: 'Coffee Beans' });
    await userEvent.click(option);

    await userEvent.click(screen.getByRole('button', { name: 'Log Delivery' }));

    await waitFor(() => {
      expect(mockApiService.logDelivery).toHaveBeenCalledWith(
        expect.objectContaining({
          inventory_item_id: 1,
          quantity: 25,
          cost: 99.99,
        })
      );
    });

    expect(await screen.findByText('Delivery logged successfully!')).toBeInTheDocument();
  });
});