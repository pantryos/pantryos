import React  from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Button,
  Snackbar,
  Alert,
  TextField,
  InputAdornment,
  Grid,
  Breadcrumbs,
  Link,
  Card,
  CardContent,
} from '@mui/material';
import { Add as AddIcon, Search as SearchIcon } from '@mui/icons-material';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import apiService from '@/services/api';
import { Delivery, InventoryItem, CreateDeliveryRequest } from '@/types/api';
import LogDelivery from '@/components/LogDelivery';

// Import custom icons from your theme
import NiPlus from "@/icons/nexture/ni-plus";
import { useCallback, useEffect, useState } from 'react';

const DeliveriesPage: React.FC = () => {
  const navigate = useNavigate();

  const [deliveries, setDeliveries] = useState<Delivery[]>([]);
  const [inventoryItems, setInventoryItems] = useState<InventoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' } | null>(null);
  const [searchTerm, setSearchTerm] = useState('');

  const loadData = useCallback(async (vendorQuery?: string) => {
    try {
      setLoading(true);
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

  useEffect(() => {
    const delayDebounceFn = setTimeout(() => {
      loadData(searchTerm);
    }, 500);
    return () => clearTimeout(delayDebounceFn);
  }, [searchTerm, loadData]);

  const handleSubmitDelivery = async (deliveryData: CreateDeliveryRequest) => {
    if (deliveryData.inventory_item_id === 0 || deliveryData.quantity <= 0) {
      setSnackbar({ open: true, message: 'Please select an item and enter a valid quantity.', severity: 'error' });
      return;
    }
    const submissionData = { ...deliveryData, delivery_date: new Date(deliveryData.delivery_date).toISOString() };

    try {
      await apiService.logDelivery(submissionData);
      setSnackbar({ open: true, message: 'Delivery logged successfully!', severity: 'success' });
      setIsModalOpen(false);
      loadData();
    } catch (error) {
      setSnackbar({ open: true, message: 'Failed to log delivery', severity: 'error' });
    }
  };

  const columns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 90 },
    {
      field: 'inventory_item_id',
      headerName: 'Item Name',
      flex: 1,
      minWidth: 200,
      valueGetter: (value, row) => inventoryItems.find((i: InventoryItem) => i.id === row.inventory_item_id)?.name || 'Unknown Item',
    },
    { field: 'vendor', headerName: 'Vendor', width: 200 },
    { field: 'quantity', headerName: 'Quantity', width: 120 },
    { field: 'cost', headerName: 'Total Cost', width: 150, valueFormatter: (value: number) => `$${value.toFixed(2)}` },
    { field: 'delivery_date', headerName: 'Delivery Date', width: 180, valueFormatter: (value: string) => new Date(value).toLocaleDateString() },
  ];

  return (
    <Grid container spacing={5}>
      {/* Page Header */}
      <Grid size={{ xs: 12 }}>
        <Typography variant="h1" component="h1" sx={{ mb: 0 }}>
          Delivery Management
        </Typography>
        <Breadcrumbs aria-label="breadcrumb">
          <Link component="button" variant="body2" onClick={() => navigate('/dashboard')}>
            Dashboard
          </Link>
          <Typography color="text.primary">Deliveries</Typography>
        </Breadcrumbs>
      </Grid>

      {/* Main Content */}
      <Grid size={{ xs: 12 }}>
        <Card>
          <CardContent>
            {/* Toolbar */}
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2, flexWrap: 'wrap', gap: 2 }}>
              <TextField
                sx={{ flexGrow: 1, minWidth: '300px' }}
                variant="standard"
                size="small"
                placeholder="Search by vendor..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                InputProps={{
                  startAdornment: <InputAdornment position="start"><SearchIcon /></InputAdornment>,
                }}
              />
              <Button
                variant="contained"
                startIcon={<NiPlus />}
                onClick={() => setIsModalOpen(true)}
              >
                Log Delivery
              </Button>
            </Box>

            {/* Data Grid */}
            <Box sx={{ height: '70vh', width: '100%' }}>
              <DataGrid
                rows={deliveries}
                columns={columns}
                loading={loading}
                initialState={{ pagination: { paginationModel: { pageSize: 15 } } }}
                pageSizeOptions={[15, 30, 50]}
                disableRowSelectionOnClick
                sx={{ border: 'none' }}
              />
            </Box>
          </CardContent>
        </Card>
      </Grid>
      
      {/* Reusable modal component */}
      <LogDelivery
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        inventoryItems={inventoryItems}
        onSubmit={handleSubmitDelivery}
      />
      
      {/* Snackbar for notifications */}
      {snackbar && (
        <Snackbar
          open={snackbar.open}
          autoHideDuration={6000}
          onClose={() => setSnackbar(null)}
          anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        >
          <Alert onClose={() => setSnackbar(null)} severity={snackbar.severity} sx={{ width: '100%' }} elevation={6} variant="filled">
            {snackbar.message}
          </Alert>
        </Snackbar>
      )}
    </Grid>
  );
};

export default DeliveriesPage;