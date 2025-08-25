import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Button,
  Snackbar,
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Switch,
  FormControlLabel,
  Chip,
  InputAdornment,
  Grid,
  Card,
  CardContent,
  Breadcrumbs,
  Link,
} from '@mui/material';
import { 
  Edit as EditIcon, 
  Delete as DeleteIcon,
  Search as SearchIcon,
} from '@mui/icons-material';
import { DataGrid, GridColDef, GridActionsCellItem } from '@mui/x-data-grid';
import apiService from '@/services/api';
import { Category, CreateCategoryRequest, UpdateCategoryRequest } from '@/types/api';
import NiPlus from "@/icons/nexture/ni-plus"; // Import custom icon

const CategoriesPage: React.FC = () => {
  const navigate = useNavigate();
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' } | null>(null);
  const [confirmDialogOpen, setConfirmDialogOpen] = useState(false);
  const [categoryToDelete, setCategoryToDelete] = useState<number | null>(null);
  const [searchTerm, setSearchTerm] = useState('');

  const [formData, setFormData] = useState<CreateCategoryRequest>({
    name: '',
    description: '',
    color: '#CCCCCC',
    is_active: true,
  });

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      const data = await apiService.getCategories();
      setCategories(data);
    } catch (error) {
      setSnackbar({ open: true, message: 'Failed to load categories', severity: 'error' });
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const filteredCategories = categories.filter(category =>
    category.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (category.description && category.description.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  const handleOpenModal = (category: Category | null = null) => {
    if (category) {
      setEditingCategory(category);
      setFormData({
        name: category.name,
        description: category.description,
        color: category.color,
        is_active: category.is_active,
      });
    } else {
      setEditingCategory(null);
      setFormData({ name: '', description: '', color: '#CCCCCC', is_active: true });
    }
    setIsModalOpen(true);
  };

  const handleFormChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = event.target;
    // ✅ FIX: Explicitly typed 'prev' parameter
    setFormData((prev: CreateCategoryRequest) => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleSubmit = async () => {
    try {
      if (editingCategory) {
        await apiService.updateCategory(editingCategory.id, formData as UpdateCategoryRequest);
        setSnackbar({ open: true, message: 'Category updated successfully!', severity: 'success' });
      } else {
        await apiService.createCategory(formData);
        setSnackbar({ open: true, message: 'Category created successfully!', severity: 'success' });
      }
      setIsModalOpen(false);
      loadData();
    } catch (error) {
      setSnackbar({ open: true, message: 'Failed to save category', severity: 'error' });
    }
  };

  const handleDeleteClick = (id: number) => {
    setCategoryToDelete(id);
    setConfirmDialogOpen(true);
  };

  const handleConfirmDelete = async () => {
    if (categoryToDelete !== null) {
      try {
        await apiService.deleteCategory(categoryToDelete);
        setSnackbar({ open: true, message: 'Category deleted successfully!', severity: 'success' });
        loadData();
      } catch (error) {
        setSnackbar({ open: true, message: 'Failed to delete category', severity: 'error' });
      } finally {
        setConfirmDialogOpen(false);
        setCategoryToDelete(null);
      }
    }
  };

  const columns: GridColDef[] = [
    { field: 'name', headerName: 'Name', flex: 1, minWidth: 150 },
    { field: 'description', headerName: 'Description', flex: 2, minWidth: 250 },
    {
      field: 'color', headerName: 'Color', width: 120,
      renderCell: (params) => (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Box sx={{ width: 20, height: 20, borderRadius: '4px', backgroundColor: params.value, border: '1px solid #ccc' }} />
          {params.value}
        </Box>
      ),
    },
    {
      field: 'is_active', headerName: 'Status', width: 120,
      renderCell: (params) => (
        <Chip label={params.value ? 'Active' : 'Inactive'} color={params.value ? 'success' : 'default'} size="small" />
      ),
    },
    {
      field: 'actions', type: 'actions', headerName: 'Actions', width: 100,
      getActions: (params) => [
        <GridActionsCellItem icon={<EditIcon />} label="Edit" onClick={() => handleOpenModal(params.row as Category)} />,
        <GridActionsCellItem icon={<DeleteIcon />} label="Delete" onClick={() => handleDeleteClick(params.row.id as number)} />,
      ],
    },
  ];

  return (
    <Grid container spacing={5}>
      {/* Page Header */}
      <Grid size={{ xs: 12 }}>
        <Typography variant="h1" component="h1" sx={{ mb: 0 }}>
          Category Management
        </Typography>
        <Breadcrumbs aria-label="breadcrumb">
          <Link component="button" variant="body2" onClick={() => navigate('/dashboard')}>
            Dashboard
          </Link>
          <Typography color="text.primary">Categories</Typography>
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
                placeholder="Search by name or description..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                InputProps={{
                  startAdornment: <InputAdornment position="start"><SearchIcon /></InputAdornment>,
                }}
              />
              <Button variant="contained" startIcon={<NiPlus />} onClick={() => handleOpenModal()}>
                Add Category
              </Button>
            </Box>

            {/* Data Grid */}
            <Box sx={{ height: '70vh', width: '100%' }}>
              <DataGrid rows={filteredCategories} columns={columns} loading={loading} sx={{ border: 'none' }} />
            </Box>
          </CardContent>
        </Card>
      </Grid>
      
      {/* Add/Edit Dialog */}
      <Dialog open={isModalOpen} onClose={() => setIsModalOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>{editingCategory ? 'Edit Category' : 'Add New Category'}</DialogTitle>
        <DialogContent>
          <TextField autoFocus margin="dense" name="name" label="Category Name" type="text" fullWidth variant="outlined" value={formData.name} onChange={handleFormChange} />
          <TextField margin="dense" name="description" label="Description" type="text" fullWidth multiline rows={3} variant="outlined" value={formData.description} onChange={handleFormChange} />
          <TextField margin="dense" name="color" label="Color" type="color" fullWidth variant="outlined" value={formData.color} onChange={handleFormChange} sx={{ maxWidth: 150 }} />
          <FormControlLabel control={<Switch checked={formData.is_active} onChange={handleFormChange} name="is_active" />} label="Active" />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setIsModalOpen(false)}>Cancel</Button>
          <Button onClick={handleSubmit} variant="contained">{editingCategory ? 'Update' : 'Create'}</Button>
        </DialogActions>
      </Dialog>
      
      {/* Delete Confirmation Dialog */}
      <Dialog open={confirmDialogOpen} onClose={() => setConfirmDialogOpen(false)}>
        <DialogTitle>Confirm Deletion</DialogTitle>
        <DialogContent>
          <Typography>Are you sure you want to delete this category? This cannot be undone.</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setConfirmDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleConfirmDelete} color="error" variant="contained">Delete</Button>
        </DialogActions>
      </Dialog>
      
      {/* Snackbar */}
      {snackbar && (
        <Snackbar open={snackbar.open} autoHideDuration={6000} onClose={() => setSnackbar(null)}>
          <Alert onClose={() => setSnackbar(null)} severity={snackbar.severity} sx={{ width: '100%' }} elevation={6} variant="filled">
            {snackbar.message}
          </Alert>
        </Snackbar>
      )}
    </Grid>
  );
};

export default CategoriesPage;