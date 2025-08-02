import React, { useState, useEffect } from 'react';
import {
  Box,
  Container,
  Grid,
  Paper,
  Typography,
  Card,
  CardContent,
  CardActions,
  Button,
  AppBar,
  Toolbar,
  IconButton,
  Menu,
  MenuItem,
  Avatar,
  Chip,
  CircularProgress
} from '@mui/material';
import {
  Inventory as InventoryIcon,
  Restaurant as MenuIcon,
  LocalShipping as DeliveryIcon,
  Assessment as AnalyticsIcon,
  AccountCircle,
  Logout
} from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import apiService from '../services/api';
import { InventoryItem, MenuItem as MenuItemType, Delivery } from '../types/api';
import LowStockBanner from './LowStockBanner';

// Dashboard component with overview cards and navigation
// Provides quick access to all major features of the inventory system
const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    inventoryItems: 0,
    menuItems: 0,
    recentDeliveries: 0,
  });

  // Load dashboard statistics
  useEffect(() => {
    const loadStats = async () => {
      try {
        const [inventory, menu, deliveries] = await Promise.all([
          apiService.getInventoryItems(),
          apiService.getMenuItems(),
          apiService.getDeliveries(),
        ]);

        setStats({
          inventoryItems: inventory.length,
          menuItems: menu.length,
          recentDeliveries: deliveries.filter(d => 
            new Date(d.delivery_date) > new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
          ).length,
        });
      } catch (error) {
        console.error('Failed to load dashboard stats:', error);
      } finally {
        setLoading(false);
      }
    };

    loadStats();
  }, []);

  // Handle user menu
  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  // Navigation functions
  const navigateTo = (path: string) => {
    navigate(path);
  };

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
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            PantryOS - Inventory Management
          </Typography>
          <IconButton
            size="large"
            edge="end"
            color="inherit"
            onClick={handleMenuOpen}
          >
            <AccountCircle />
          </IconButton>
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
          >
            <MenuItem disabled>
              <Typography variant="body2">
                {user?.email}
              </Typography>
            </MenuItem>
            <MenuItem onClick={handleLogout}>
              <Logout sx={{ mr: 1 }} />
              Logout
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      {/* Main Content */}
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        {/* Welcome Section */}
        <Paper sx={{ p: 3, mb: 3 }}>
          <Typography variant="h4" gutterBottom>
            Welcome back!
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Manage your inventory, menu items, and track deliveries from your dashboard.
          </Typography>
        </Paper>

        {/* Low Stock Banner */}
        <LowStockBanner maxItems={5} />

        {/* Stats Cards */}
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' }, gap: 3, mb: 4 }}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" justifyContent="space-between">
                <Box>
                  <Typography color="textSecondary" gutterBottom>
                    Inventory Items
                  </Typography>
                  <Typography variant="h4">
                    {stats.inventoryItems}
                  </Typography>
                </Box>
                <InventoryIcon color="primary" sx={{ fontSize: 40 }} />
              </Box>
            </CardContent>
            <CardActions>
              <Button size="small" onClick={() => navigateTo('/inventory')}>
                View All
              </Button>
              <Button size="small" onClick={() => navigateTo('/inventory/new')}>
                Add New
              </Button>
            </CardActions>
          </Card>

          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" justifyContent="space-between">
                <Box>
                  <Typography color="textSecondary" gutterBottom>
                    Menu Items
                  </Typography>
                  <Typography variant="h4">
                    {stats.menuItems}
                  </Typography>
                </Box>
                <MenuIcon color="primary" sx={{ fontSize: 40 }} />
              </Box>
            </CardContent>
            <CardActions>
              <Button size="small" onClick={() => navigateTo('/menu')}>
                View All
              </Button>
              <Button size="small" onClick={() => navigateTo('/menu/new')}>
                Add New
              </Button>
            </CardActions>
          </Card>

          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" justifyContent="space-between">
                <Box>
                  <Typography color="textSecondary" gutterBottom>
                    Recent Deliveries
                  </Typography>
                  <Typography variant="h4">
                    {stats.recentDeliveries}
                  </Typography>
                  <Typography variant="caption" color="textSecondary">
                    Last 7 days
                  </Typography>
                </Box>
                <DeliveryIcon color="primary" sx={{ fontSize: 40 }} />
              </Box>
            </CardContent>
            <CardActions>
              <Button size="small" onClick={() => navigateTo('/deliveries')}>
                View All
              </Button>
              <Button size="small" onClick={() => navigateTo('/deliveries/new')}>
                Log Delivery
              </Button>
            </CardActions>
          </Card>
        </Box>

        {/* Quick Actions */}
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Quick Actions
          </Typography>
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(4, 1fr)' }, gap: 2 }}>
            <Button
              fullWidth
              variant="outlined"
              startIcon={<InventoryIcon />}
              onClick={() => navigateTo('/inventory/new')}
              sx={{ height: 56 }}
            >
              Add Inventory Item
            </Button>
            <Button
              fullWidth
              variant="outlined"
              startIcon={<MenuIcon />}
              onClick={() => navigateTo('/menu/new')}
              sx={{ height: 56 }}
            >
              Add Menu Item
            </Button>
            <Button
              fullWidth
              variant="outlined"
              startIcon={<DeliveryIcon />}
              onClick={() => navigateTo('/deliveries/new')}
              sx={{ height: 56 }}
            >
              Log Delivery
            </Button>
            <Button
              fullWidth
              variant="outlined"
              startIcon={<AnalyticsIcon />}
              onClick={() => navigateTo('/analytics')}
              sx={{ height: 56 }}
            >
              View Analytics
            </Button>
          </Box>
        </Paper>
      </Container>
    </Box>
  );
};

export default Dashboard; 