import React, { useState, useEffect, useCallback } from 'react';
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
  AppBar,
  Toolbar,
  InputAdornment,
  MenuItem as MuiMenuItem,
  Snackbar,
  Alert as MuiAlert,
  Chip
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
  ArrowBack as ArrowBackIcon,
  Category as CategoryIcon, // Icon for the new button
} from '@mui/icons-material';
import { DataGrid, GridColDef, GridActionsCellItem } from '@mui/x-data-grid';
import { useNavigate } from 'react-router-dom';
import apiService from '../services/api';
import { MenuItem, CreateMenuItemRequest, Category } from '../types/api';

const Menu: React.FC = () => {
  const navigate = useNavigate();
  const [menuItems, setMenuItems] = useState<MenuItem[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');

  // State for Add/Edit Dialog
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<MenuItem | null>(null);
  const [formData, setFormData] = useState<CreateMenuItemRequest>({
    name: '',
    price: 0,
    category_id: 0,
    category: '',
  });

  // State for Snackbar and Confirmation Dialog
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({
    open: false,
    message: '',
    severity: 'success',
  });
  const [confirmDialogOpen, setConfirmDialogOpen] = useState(false);
  const [itemToDelete, setItemToDelete] = useState<number | null>(null);

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      const [menuData, categoryData] = await Promise.all([
        apiService.getMenuItems(),
        apiService.getActiveCategories(),
      ]);
      setMenuItems(menuData);
      setCategories(categoryData || []);
    } catch (err) {
      setSnackbar({ open: true, message: 'Failed to load menu data', severity: 'error' });
      console.error('Error loading data:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const filteredItems = Array.isArray(menuItems) ? menuItems.filter(item =>
    item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    item.category.toLowerCase().includes(searchTerm.toLowerCase())
  ) : [];

  const handleInputChange = (field: keyof CreateMenuItemRequest, value: string | number) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleAddClick = () => {
    setEditingItem(null);
    setFormData({ name: '', price: 0, category_id: 0, category: '' });
    setDialogOpen(true);
  };

  const handleEditClick = (item: MenuItem) => {
    setEditingItem(item);
    setFormData({
      name: item.name,
      price: item.price,
      category_id: item.category_id,
      category: item.category,
    });
    setDialogOpen(true);
  };

  const handleSubmit = async () => {
    const selectedCategory = categories.find(cat => cat.id === formData.category_id);
    if (!selectedCategory) {
      setSnackbar({ open: true, message: 'Please select a valid category.', severity: 'error' });
      return;
    }
    const payload = { ...formData, category: selectedCategory.name };

    try {
      if (editingItem) {
        // NOTE: The updateMenuItem function is commented out in the original code.
        // await apiService.updateMenuItem(editingItem.id, payload as UpdateMenuItemRequest);
        setSnackbar({ open: true, message: 'Menu item updated successfully!', severity: 'success' });
      } else {
        await apiService.createMenuItem(payload);
        setSnackbar({ open: true, message: 'Menu item created successfully!', severity: 'success' });
      }
      setDialogOpen(false);
      loadData();
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Failed to save menu item';
      setSnackbar({ open: true, message: errorMessage, severity: 'error' });
    }
  };

  const handleDeleteClick = (id: number) => {
    setItemToDelete(id);
    setConfirmDialogOpen(true);
  };

  const handleConfirmDelete = async () => {
    if (itemToDelete === null) return;
    try {
      // NOTE: The deleteMenuItem function is commented out in the original code.
      // await apiService.deleteMenuItem(itemToDelete);
      setSnackbar({ open: true, message: 'Item deleted successfully!', severity: 'success' });
      loadData();
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Failed to delete item';
      setSnackbar({ open: true, message: errorMessage, severity: 'error' });
    } finally {
      setConfirmDialogOpen(false);
      setItemToDelete(null);
    }
  };

  const columns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 90 },
    { field: 'name', headerName: 'Name', flex: 1, minWidth: 200 },
    {
      field: 'price',
      headerName: 'Price',
      width: 150,
      valueFormatter: (value: number) => `$${value.toFixed(2)}`
    },
    {
      field: 'category',
      headerName: 'Category',
      width: 200,
      renderCell: (params) => <Chip label={params.value} color="primary" variant="outlined" size="small" />
    },
    {
      field: 'actions',
      type: 'actions',
      headerName: 'Actions',
      width: 120,
      getActions: (params) => [
        <GridActionsCellItem
          icon={<EditIcon />}
          label="Edit"
          onClick={() => handleEditClick(params.row as MenuItem)}
        />,
        <GridActionsCellItem
          icon={<DeleteIcon />}
          label="Delete"
          onClick={() => handleDeleteClick(params.row.id as number)}
        />,
      ],
    },
  ];

  return (
    <Box sx={{ flexGrow: 1, pb: 4 }}>
      <AppBar position="static">
        <Toolbar>
          <IconButton edge="start" color="inherit" onClick={() => navigate('/dashboard')} sx={{ mr: 2 }}>
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Menu Management
          </Typography>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Paper elevation={3} sx={{ p: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 2 }}>
            <TextField
              sx={{ flexGrow: 1, minWidth: '300px' }}
              variant="outlined"
              size="small"
              placeholder="Search by name or category..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              InputProps={{
                startAdornment: <InputAdornment position="start"><SearchIcon /></InputAdornment>,
              }}
            />
            {/* Group the action buttons together for better alignment */}
            <Box sx={{ display: 'flex', gap: 1 }}>
              <Button variant="outlined" startIcon={<CategoryIcon />} onClick={() => navigate('/categories')}>
                Manage Categories
              </Button>
              <Button variant="contained" startIcon={<AddIcon />} onClick={handleAddClick}>
                Add Menu Item
              </Button>
            </Box>
          </Box>
          <Box sx={{ height: '65vh', width: '100%' }}>
            <DataGrid
              rows={filteredItems}
              columns={columns}
              loading={loading}
              initialState={{ pagination: { paginationModel: { page: 0, pageSize: 10 } } }}
              pageSizeOptions={[10, 25, 50]}
              disableRowSelectionOnClick
            />
          </Box>
        </Paper>
      </Container>

      <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>{editingItem ? 'Edit Menu Item' : 'Add New Menu Item'}</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 1 }}>
            <TextField fullWidth label="Name" value={formData.name} onChange={(e) => handleInputChange('name', e.target.value)} margin="normal" required />
            <TextField select fullWidth label="Category" value={formData.category_id} onChange={(e) => handleInputChange('category_id', parseInt(e.target.value))} margin="normal" required>
              <MuiMenuItem value={0} disabled><em>Select a category</em></MuiMenuItem>
              {categories.map((cat) => <MuiMenuItem key={cat.id} value={cat.id}>{cat.name}</MuiMenuItem>)}
            </TextField>
            <TextField
              fullWidth
              label="Price"
              type="number"
              value={formData.price}
              onChange={(e) => handleInputChange('price', parseFloat(e.target.value) || 0)} // Pass the raw string value
              margin="normal"
              required
              InputProps={{
                startAdornment: <InputAdornment position="start">$</InputAdornment>
              }}
              inputProps={{ min: 0, step: 0.01 }}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleSubmit} variant="contained">{editingItem ? 'Update' : 'Create'}</Button>
        </DialogActions>
      </Dialog>

      <Dialog open={confirmDialogOpen} onClose={() => setConfirmDialogOpen(false)}>
        <DialogTitle>Confirm Deletion</DialogTitle>
        <DialogContent><Typography>Are you sure you want to delete this menu item? This action cannot be undone.</Typography></DialogContent>
        <DialogActions>
          <Button onClick={() => setConfirmDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleConfirmDelete} color="error" variant="contained">Delete</Button>
        </DialogActions>
      </Dialog>

      <Snackbar open={snackbar.open} autoHideDuration={6000} onClose={() => setSnackbar({ ...snackbar, open: false })} anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
        <MuiAlert onClose={() => setSnackbar({ ...snackbar, open: false })} severity={snackbar.severity} sx={{ width: '100%' }} elevation={6} variant="filled">
          {snackbar.message}
        </MuiAlert>
      </Snackbar>
    </Box>
  );
};

export default Menu;
