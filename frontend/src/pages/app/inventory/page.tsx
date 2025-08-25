import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Typography,
  Button,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  CircularProgress,
  InputAdornment,
  MenuItem,
  Grid,
  Breadcrumbs,
  Link,
  Card,
  CardContent,
  Snackbar,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
  Analytics as AnalyticsIcon,
} from '@mui/icons-material';
import { DataGrid, GridColDef, GridActionsCellItem, GridRenderCellParams } from '@mui/x-data-grid';
import { useNavigate } from 'react-router-dom';
import apiService from '@/services/api';
import { InventoryItem, CreateInventoryItemRequest, UpdateInventoryItemRequest, Category } from '@/types/api';
import ItemUsageAnalytics from '@/components/ItemUsageAnalytics';

// Import custom icons from your theme
import NiPlus from "@/icons/nexture/ni-plus";
import NiLayout from "@/icons/nexture/ni-layout";

const Inventory: React.FC = () => {
  const navigate = useNavigate();
  const [items, setItems] = useState<InventoryItem[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  
  // Dialog and Form State
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<InventoryItem | null>(null);
  const [formData, setFormData] = useState<CreateInventoryItemRequest>({
    name: '',
    unit: '',
    cost_per_unit: 0,
    preferred_vendor: '',
    min_stock_level: 0,
    max_stock_level: 0,
    min_weeks_stock: 2,
    max_weeks_stock: 8,
    category_id: undefined,
  });

  // Analytics Dialog State
  const [analyticsDialogOpen, setAnalyticsDialogOpen] = useState(false);
  const [selectedItemForAnalytics, setSelectedItemForAnalytics] = useState<InventoryItem | null>(null);

  // Delete Confirmation Dialog State
  const [confirmDialogOpen, setConfirmDialogOpen] = useState(false);
  const [itemToDelete, setItemToDelete] = useState<number | null>(null);

  // Snackbar State
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' | 'warning' }>({
    open: false,
    message: '',
    severity: 'success',
  });

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      const [inventoryData, categoryData] = await Promise.all([
        apiService.getInventoryItems(),
        apiService.getActiveCategories()
      ]);
      setItems(inventoryData);
      setCategories(categoryData || []);
    } catch (error) {
      setSnackbar({ open: true, message: 'Failed to load data', severity: 'error' });
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const filteredItems = Array.isArray(items) ? items.filter(item =>
    item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (item.preferred_vendor && item.preferred_vendor.toLowerCase().includes(searchTerm.toLowerCase()))
  ) : [];

  const handleInputChange = (field: keyof CreateInventoryItemRequest, value: string | number | undefined) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleAddClick = () => {
    setEditingItem(null);
    setFormData({
      name: '', unit: '', cost_per_unit: 0, preferred_vendor: '',
      min_stock_level: 0, max_stock_level: 0, min_weeks_stock: 2,
      max_weeks_stock: 8, category_id: undefined,
    });
    setDialogOpen(true);
  };

  const handleEditClick = (item: InventoryItem) => {
    setEditingItem(item);
    setFormData({
      name: item.name, unit: item.unit, cost_per_unit: item.cost_per_unit,
      preferred_vendor: item.preferred_vendor, min_stock_level: item.min_stock_level,
      max_stock_level: item.max_stock_level, min_weeks_stock: item.min_weeks_stock,
      max_weeks_stock: item.max_weeks_stock, category_id: item.category_id,
    });
    setDialogOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const payload = { ...formData, category_id: formData.category_id ? Number(formData.category_id) : undefined };
      if (editingItem) {
        await apiService.updateInventoryItem(editingItem.id, payload as UpdateInventoryItemRequest);
        setSnackbar({ open: true, message: 'Item updated successfully!', severity: 'success' });
      } else {
        await apiService.createInventoryItem(payload);
        setSnackbar({ open: true, message: 'Item created successfully!', severity: 'success' });
      }
      setDialogOpen(false);
      loadData();
    } catch (error: any) {
      setSnackbar({ open: true, message: error.response?.data?.error || 'Failed to save item', severity: 'error' });
    }
  };

  const handleDeleteClick = (id: number) => {
    setItemToDelete(id);
    setConfirmDialogOpen(true);
  };
  
  const handleConfirmDelete = async () => {
    if (itemToDelete === null) return;
    try {
      await apiService.deleteInventoryItem(itemToDelete);
      setSnackbar({ open: true, message: 'Item deleted successfully!', severity: 'success' });
      loadData();
    } catch (error: any) {
      setSnackbar({ open: true, message: error.response?.data?.error || 'Failed to delete item', severity: 'error' });
    } finally {
      setConfirmDialogOpen(false);
      setItemToDelete(null);
    }
  };

  const handleAnalyticsClick = (item: InventoryItem) => {
    setSelectedItemForAnalytics(item);
    setAnalyticsDialogOpen(true);
  };

  const columns: GridColDef[] = [
    { field: 'name', headerName: 'Name', flex: 1, minWidth: 200 },
    {
      field: 'current_stock',
      headerName: 'Current Stock',
      width: 130,
      renderCell: (params: GridRenderCellParams<any, number>) => {
        const stock = params.value ?? 0;
        const unit = params.row.unit || '';
        return `${stock} ${unit}`;
      },
      cellClassName: (params) => {
        const { min_stock_level, current_stock } = params.row;
        if (current_stock < min_stock_level) {
          return 'low-stock-cell';
        }
        return '';
      }
    },
    { field: 'min_stock_level', headerName: 'Min Stock', width: 120 },
    { field: 'cost_per_unit', headerName: 'Cost/Unit', width: 120, valueFormatter: (value) => value ? `$${Number(value).toFixed(2)}` : '$0.00' },
    { field: 'preferred_vendor', headerName: 'Vendor', flex: 1, minWidth: 150 },
    {
      field: 'category_id',
      headerName: 'Category',
      width: 150,
      valueGetter: (value) => categories.find(cat => cat.id === value)?.name || 'N/A'
    },
    {
      field: 'actions', type: 'actions', headerName: 'Actions', width: 120,
      getActions: (params) => [
        <GridActionsCellItem icon={<AnalyticsIcon />} label="Analytics" onClick={() => handleAnalyticsClick(params.row)} />,
        <GridActionsCellItem icon={<EditIcon />} label="Edit" onClick={() => handleEditClick(params.row)} />,
        <GridActionsCellItem icon={<DeleteIcon />} label="Delete" onClick={() => handleDeleteClick(params.row.id)} />,
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
    <Grid container spacing={5}>
      {/* Page Header */}
      <Grid size={{ xs: 12 }}>
        <Typography variant="h1" component="h1" sx={{ mb: 0 }}>
          Inventory Management
        </Typography>
        <Breadcrumbs aria-label="breadcrumb">
          <Link component="button" variant="body2" onClick={() => navigate('/dashboard')}>
            Dashboard
          </Link>
          <Typography color="text.primary">Inventory</Typography>
        </Breadcrumbs>
      </Grid>

      {/* Main Content */}
      <Grid size={{ xs: 12 }}>
        <Card>
          <CardContent>
            {/* Toolbar */}
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 2, mb: 2 }}>
              <TextField
                sx={{ flexGrow: 1, minWidth: '300px' }}
                variant="standard"
                size="small"
                placeholder="Search by name or vendor..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                InputProps={{ startAdornment: <InputAdornment position="start"><SearchIcon /></InputAdornment> }}
              />
              <Box sx={{ display: 'flex', gap: 1 }}>
                <Button variant="outlined" color="grey" startIcon={<NiLayout />} onClick={() => navigate('/categories')}>
                  Manage Categories
                </Button>
                <Button variant="contained" startIcon={<NiPlus />} onClick={handleAddClick}>
                  Add Item
                </Button>
              </Box>
            </Box>

            {/* Data Grid */}
            <Box sx={{ height: '70vh', width: '100%', '.low-stock-cell': { backgroundColor: 'rgba(255, 170, 0, 0.2)' } }}>
              <DataGrid
                rows={filteredItems}
                columns={columns}
                disableRowSelectionOnClick
                initialState={{ pagination: { paginationModel: { pageSize: 15 } } }}
                pageSizeOptions={[15, 30, 50]}
                sx={{ border: 'none' }}
              />
            </Box>
          </CardContent>
        </Card>
      </Grid>
      
      {/* Add/Edit Dialog */}
      <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>{editingItem ? 'Edit Inventory Item' : 'Add New Item'}</DialogTitle>
        <DialogContent>
          <TextField fullWidth label="Name" value={formData.name} onChange={(e) => handleInputChange('name', e.target.value)} margin="normal" required />
          <TextField select fullWidth label="Category" value={formData.category_id || ''} onChange={(e) => handleInputChange('category_id', e.target.value ? parseInt(e.target.value) : undefined)} margin="normal">
            <MenuItem value=""><em>None</em></MenuItem>
            {categories.map((cat) => <MenuItem key={cat.id} value={cat.id}>{cat.name}</MenuItem>)}
          </TextField>
          <TextField fullWidth label="Unit (e.g., kg, L, pcs)" value={formData.unit} onChange={(e) => handleInputChange('unit', e.target.value)} margin="normal" required />
          <TextField fullWidth label="Cost per Unit" type="number" value={formData.cost_per_unit} onChange={(e) => handleInputChange('cost_per_unit', parseFloat(e.target.value) || 0)} margin="normal" required inputProps={{ min: 0, step: 0.01 }} />
          <TextField fullWidth label="Preferred Vendor" value={formData.preferred_vendor} onChange={(e) => handleInputChange('preferred_vendor', e.target.value)} margin="normal" />
          <TextField fullWidth label="Min Stock Level" type="number" value={formData.min_stock_level} onChange={(e) => handleInputChange('min_stock_level', parseFloat(e.target.value) || 0)} margin="normal" inputProps={{ min: 0 }} />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleSubmit} variant="contained">{editingItem ? 'Update' : 'Create'}</Button>
        </DialogActions>
      </Dialog>
      
      {/* Delete Confirmation Dialog */}
      <Dialog open={confirmDialogOpen} onClose={() => setConfirmDialogOpen(false)}>
        <DialogTitle>Confirm Deletion</DialogTitle>
        <DialogContent><Typography>Are you sure you want to delete this item? This action is permanent.</Typography></DialogContent>
        <DialogActions>
          <Button onClick={() => setConfirmDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleConfirmDelete} color="error" variant="contained">Delete</Button>
        </DialogActions>
      </Dialog>
      
      {/* Analytics Dialog */}
      <Dialog open={analyticsDialogOpen} onClose={() => setAnalyticsDialogOpen(false)} maxWidth="xl" fullWidth PaperProps={{ sx: { height: '90vh' } }}>
        <DialogTitle>{selectedItemForAnalytics?.name} - Usage Analytics</DialogTitle>
        <DialogContent dividers>
          {selectedItemForAnalytics && <ItemUsageAnalytics item={selectedItemForAnalytics} onClose={() => setAnalyticsDialogOpen(false)} />}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setAnalyticsDialogOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>

      {/* Snackbar for Notifications */}
      <Snackbar open={snackbar.open} autoHideDuration={6000} onClose={() => setSnackbar({ ...snackbar, open: false })} anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
        <Alert onClose={() => setSnackbar({ ...snackbar, open: false })} severity={snackbar.severity} sx={{ width: '100%' }} elevation={6} variant="filled">
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Grid>
  );
};

export default Inventory;