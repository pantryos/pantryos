import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Container,
  Paper,
  Typography,
  Button,
  AppBar,
  Toolbar,
  IconButton,
  Snackbar,
  Alert,
  TextField,
  InputAdornment,
} from '@mui/material';
import { Add as AddIcon, ArrowBack as ArrowBackIcon, Search as SearchIcon } from '@mui/icons-material';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import apiService from '../services/api';
import { Delivery, InventoryItem, CreateDeliveryRequest } from '../types/api';
import LogDelivery from './LogDelivery'; // Re-use the modal component

const DeliveriesPage: React.FC = () => {
  const navigate = useNavigate();

  // --- State Management ---
  const [deliveries, setDeliveries] = useState<Delivery[]>([]);
  const [inventoryItems, setInventoryItems] = useState<InventoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' } | null>(null);
  const [searchTerm, setSearchTerm] = useState('');

  // --- Data Fetching ---
  const loadData = useCallback(async (vendorQuery?: string) => {
    try {
      setLoading(true);
      // Fetch both deliveries and inventory items at the same time
      const [deliveriesData, itemsData] = await Promise.all([
        apiService.getDeliveries(vendorQuery),
        apiService.getInventoryItems(),
      ]);
      setDeliveries(deliveriesData || []);
      setInventoryItems(itemsData || []);
    } catch (error) {
      setSnackbar({ open: true, message: 'Failed to load data', severity: 'error' });
    } finally {
      setLoading(false);
    }
  }, []);

  // useEffect hook to handle debounced searching
  useEffect(() => {
    // Set a timer to fetch data 500ms after the user stops typing
    const delayDebounceFn = setTimeout(() => {
      loadData(searchTerm);
    }, 500);

    // Cleanup function to cancel the timer if the user types again
    return () => clearTimeout(delayDebounceFn);
  }, [searchTerm, loadData]); // This effect re-runs whenever searchTerm changes

  // --- Event Handlers ---
  const handleSubmitDelivery = async (deliveryData: CreateDeliveryRequest) => {
    if (deliveryData.inventory_item_id === 0 || deliveryData.quantity <= 0) {
      setSnackbar({ open: true, message: 'Please select an item and enter a valid quantity.', severity: 'error' });
      return;
    }

    // Convert the date to a full ISO string for the backend
    const submissionData = {
      ...deliveryData,
      delivery_date: new Date(deliveryData.delivery_date).toISOString(),
    };

    try {
      await apiService.logDelivery(submissionData);
      setSnackbar({ open: true, message: 'Delivery logged successfully!', severity: 'success' });
      setIsModalOpen(false); // Close the modal
      loadData(); // Refresh the list of deliveries
    } catch (error) {
      setSnackbar({ open: true, message: 'Failed to log delivery', severity: 'error' });
    }
  };

  // --- Data Grid Columns ---
  const columns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 90 },
    {
      field: 'inventory_item_id',
      headerName: 'Item Name',
      flex: 1,
      minWidth: 200,
      valueGetter: (value, row) => {
        // Find the item name from the inventoryItems list
        const item = inventoryItems.find(i => i.id === row.inventory_item_id);
        return item ? item.name : 'Unknown Item';
      },
    },
    { field: 'vendor', headerName: 'Vendor', width: 200 },
    { field: 'quantity', headerName: 'Quantity', width: 120 },
    {
      field: 'cost',
      headerName: 'Total Cost',
      width: 150,
      valueFormatter: (value: number) => `$${value.toFixed(2)}`,
    },
    {
      field: 'delivery_date',
      headerName: 'Delivery Date',
      width: 180,
      valueFormatter: (value: string) => new Date(value).toLocaleDateString(),
    },
  ];

  return (
    <Box>
      <AppBar position="static">
        <Toolbar>
          <IconButton edge="start" color="inherit" onClick={() => navigate(-1)} sx={{ mr: 2 }}>
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Delivery Management
          </Typography>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Paper elevation={3} sx={{ p: 2 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2, flexWrap: 'wrap', gap: 2 }}>
            <TextField
              sx={{ flexGrow: 1, minWidth: '300px' }}
              variant="outlined"
              size="small"
              placeholder="Search by vendor..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
              }}
            />
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => setIsModalOpen(true)}
            >
              Log New Delivery
            </Button>
          </Box>
          <Box sx={{ height: '70vh', width: '100%' }}>
            <DataGrid
              rows={deliveries}
              columns={columns}
              loading={loading}
              initialState={{ pagination: { paginationModel: { pageSize: 10 } } }}
              pageSizeOptions={[10, 25, 50]}
            />
          </Box>
        </Paper>
      </Container>

      {/* The reusable modal component for adding a new delivery */}
      <LogDelivery
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        inventoryItems={inventoryItems}
        onSubmit={handleSubmitDelivery}
      />

      {snackbar && (
        <Snackbar
          open={snackbar.open}
          autoHideDuration={6000}
          onClose={() => setSnackbar(null)}
          anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        >
          <Alert onClose={() => setSnackbar(null)} severity={snackbar.severity} sx={{ width: '100%' }}>
            {snackbar.message}
          </Alert>
        </Snackbar>
      )}
    </Box>
  );
};

export default DeliveriesPage;