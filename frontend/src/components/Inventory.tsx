import React, { useState, useEffect } from 'react';
import {
  Box,
  Container,
  Paper,
  Typography,
  Button,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Alert,
  CircularProgress,
  AppBar,
  Toolbar,
  InputAdornment
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
  ArrowBack as ArrowBackIcon,
  Analytics as AnalyticsIcon
} from '@mui/icons-material';
import { DataGrid, GridColDef, GridActionsCellItem } from '@mui/x-data-grid';
import { useNavigate } from 'react-router-dom';
import apiService from '../services/api';
import { InventoryItem, CreateInventoryItemRequest, UpdateInventoryItemRequest } from '../types/api';
import ItemUsageAnalytics from './ItemUsageAnalytics';

// Inventory management component with full CRUD operations
// Features: list view, search, add, edit, delete inventory items
const Inventory: React.FC = () => {
  const navigate = useNavigate();
  const [items, setItems] = useState<InventoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<InventoryItem | null>(null);
  const [error, setError] = useState('');
  const [analyticsDialogOpen, setAnalyticsDialogOpen] = useState(false);
  const [selectedItemForAnalytics, setSelectedItemForAnalytics] = useState<InventoryItem | null>(null);
  const [formData, setFormData] = useState<CreateInventoryItemRequest>({
    name: '',
    unit: '',
    cost_per_unit: 0,
    preferred_vendor: '',
    min_stock_level: 0,
    max_stock_level: 0,
  });

  // Load inventory items on component mount
  useEffect(() => {
    loadItems();
  }, []);

  // Load inventory items from API
  const loadItems = async () => {
    try {
      setLoading(true);
      const data = await apiService.getInventoryItems();
      setItems(data);
    } catch (error) {
      setError('Failed to load inventory items');
      console.error('Error loading inventory items:', error);
    } finally {
      setLoading(false);
    }
  };

  // Filter items based on search term
  const filteredItems = items.filter(item =>
    item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    item.preferred_vendor.toLowerCase().includes(searchTerm.toLowerCase())
  );

  // Handle form input changes
  const handleInputChange = (field: keyof CreateInventoryItemRequest, value: string | number) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }));
  };

  // Open dialog for adding new item
  const handleAddClick = () => {
    setEditingItem(null);
    setFormData({
      name: '',
      unit: '',
      cost_per_unit: 0,
      preferred_vendor: '',
      min_stock_level: 0,
      max_stock_level: 0,
    });
    setDialogOpen(true);
  };

  // Open dialog for editing item
  const handleEditClick = (item: InventoryItem) => {
    setEditingItem(item);
    setFormData({
      name: item.name,
      unit: item.unit,
      cost_per_unit: item.cost_per_unit,
      preferred_vendor: item.preferred_vendor,
      min_stock_level: item.min_stock_level,
      max_stock_level: item.max_stock_level,
    });
    setDialogOpen(true);
  };

  // Handle form submission (create or update)
  const handleSubmit = async () => {
    try {
      if (editingItem) {
        await apiService.updateInventoryItem(editingItem.id, formData as UpdateInventoryItemRequest);
      } else {
        await apiService.createInventoryItem(formData);
      }
      setDialogOpen(false);
      loadItems(); // Reload items
    } catch (error: any) {
      setError(error.response?.data?.error || 'Failed to save item');
    }
  };

  // Handle item deletion
  const handleDeleteClick = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this item?')) {
      try {
        await apiService.deleteInventoryItem(id);
        loadItems(); // Reload items
      } catch (error: any) {
        setError(error.response?.data?.error || 'Failed to delete item');
      }
    }
  };

  // Handle analytics view
  const handleAnalyticsClick = (item: InventoryItem) => {
    setSelectedItemForAnalytics(item);
    setAnalyticsDialogOpen(true);
  };

  // Handle analytics dialog close
  const handleAnalyticsClose = () => {
    setAnalyticsDialogOpen(false);
    setSelectedItemForAnalytics(null);
  };

  // Data grid columns configuration
  const columns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 70 },
    { field: 'name', headerName: 'Name', width: 200 },
    { field: 'unit', headerName: 'Unit', width: 100 },
    { 
      field: 'cost_per_unit', 
      headerName: 'Cost/Unit', 
      width: 120,
      valueFormatter: (params: any) => `$${params.value.toFixed(2)}`
    },
    { field: 'preferred_vendor', headerName: 'Vendor', width: 200 },
    { 
      field: 'min_stock_level', 
      headerName: 'Min Stock', 
      width: 120,
      valueFormatter: (params: any) => params.value.toString()
    },
    { 
      field: 'max_stock_level', 
      headerName: 'Max Stock', 
      width: 120,
      valueFormatter: (params: any) => params.value.toString()
    },
    {
      field: 'actions',
      type: 'actions',
      headerName: 'Actions',
      width: 180,
      getActions: (params) => [
        <GridActionsCellItem
          icon={<AnalyticsIcon />}
          label="Analytics"
          onClick={() => handleAnalyticsClick(params.row)}
        />,
        <GridActionsCellItem
          icon={<EditIcon />}
          label="Edit"
          onClick={() => handleEditClick(params.row)}
        />,
        <GridActionsCellItem
          icon={<DeleteIcon />}
          label="Delete"
          onClick={() => handleDeleteClick(params.row.id)}
        />,
      ],
    },
  ];

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ flexGrow: 1 }}>
      {/* App Bar */}
      <AppBar position="static">
        <Toolbar>
          <IconButton
            edge="start"
            color="inherit"
            onClick={() => navigate('/dashboard')}
            sx={{ mr: 2 }}
          >
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Inventory Management
          </Typography>
        </Toolbar>
      </AppBar>

      {/* Main Content */}
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        {/* Header */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" gutterBottom>
            Inventory Items
          </Typography>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={handleAddClick}
          >
            Add Item
          </Button>
        </Box>

        {/* Error Alert */}
        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
            {error}
          </Alert>
        )}

        {/* Search Bar */}
        <Paper sx={{ p: 2, mb: 3 }}>
          <TextField
            fullWidth
            placeholder="Search by name or vendor..."
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
        </Paper>

        {/* Data Grid */}
        <Paper sx={{ height: 600, width: '100%' }}>
          <DataGrid
            rows={filteredItems}
            columns={columns}
            initialState={{
              pagination: {
                paginationModel: { page: 0, pageSize: 10 },
              },
            }}
            pageSizeOptions={[10, 25, 50]}
            disableRowSelectionOnClick
            sx={{
              '& .MuiDataGrid-cell:focus': {
                outline: 'none',
              },
            }}
          />
        </Paper>

        {/* Add/Edit Dialog */}
        <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>
            {editingItem ? 'Edit Inventory Item' : 'Add New Inventory Item'}
          </DialogTitle>
          <DialogContent>
            <Box sx={{ pt: 1 }}>
              <TextField
                fullWidth
                label="Name"
                value={formData.name}
                onChange={(e) => handleInputChange('name', e.target.value)}
                margin="normal"
                required
              />
              <TextField
                fullWidth
                label="Unit"
                value={formData.unit}
                onChange={(e) => handleInputChange('unit', e.target.value)}
                margin="normal"
                required
                placeholder="e.g., kg, liters, pieces"
              />
              <TextField
                fullWidth
                label="Cost per Unit"
                type="number"
                value={formData.cost_per_unit}
                onChange={(e) => handleInputChange('cost_per_unit', parseFloat(e.target.value) || 0)}
                margin="normal"
                required
                inputProps={{ min: 0, step: 0.01 }}
              />
              <TextField
                fullWidth
                label="Preferred Vendor"
                value={formData.preferred_vendor}
                onChange={(e) => handleInputChange('preferred_vendor', e.target.value)}
                margin="normal"
              />
              <TextField
                fullWidth
                label="Minimum Stock Level"
                type="number"
                value={formData.min_stock_level}
                onChange={(e) => handleInputChange('min_stock_level', parseFloat(e.target.value) || 0)}
                margin="normal"
                inputProps={{ min: 0 }}
              />
              <TextField
                fullWidth
                label="Maximum Stock Level"
                type="number"
                value={formData.max_stock_level}
                onChange={(e) => handleInputChange('max_stock_level', parseFloat(e.target.value) || 0)}
                margin="normal"
                inputProps={{ min: 0 }}
              />
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setDialogOpen(false)}>Cancel</Button>
            <Button onClick={handleSubmit} variant="contained">
              {editingItem ? 'Update' : 'Create'}
            </Button>
          </DialogActions>
        </Dialog>

        {/* Analytics Dialog */}
        <Dialog 
          open={analyticsDialogOpen} 
          onClose={handleAnalyticsClose} 
          maxWidth="xl" 
          fullWidth
          PaperProps={{
            sx: { height: '90vh' }
          }}
        >
          <DialogTitle>
            <Typography variant="h5">
              {selectedItemForAnalytics?.name} - Usage Analytics
            </Typography>
          </DialogTitle>
          <DialogContent dividers>
            {selectedItemForAnalytics && (
              <ItemUsageAnalytics 
                item={selectedItemForAnalytics} 
                onClose={handleAnalyticsClose}
              />
            )}
          </DialogContent>
          <DialogActions>
            <Button onClick={handleAnalyticsClose}>
              Close
            </Button>
          </DialogActions>
        </Dialog>
      </Container>
    </Box>
  );
};

export default Inventory; 