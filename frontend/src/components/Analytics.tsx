import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Container,
  Grid,
  Paper,
  Typography,
  Card,
  CardContent,
  CardHeader,
  Button,
  AppBar,
  Toolbar,
  IconButton,
  Menu,
  MenuItem,
  Avatar,
  Chip,
  CircularProgress,
  FormControl,
  InputLabel,
  Select,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
  Alert,
  Divider,
  Stack,
  LinearProgress,
} from '@mui/material';
import {
  Analytics as AnalyticsIcon,
  TrendingUp,
  TrendingDown,
  Remove,
  Compare,
  FilterList,
  Refresh,
  AccountCircle,
  Logout,
  BarChart,
  ShowChart,
  PieChart,
  TableChart,
} from '@mui/icons-material';
import {
  LineChart,
  Line,
  BarChart as RechartsBarChart,
  Bar,
  PieChart as RechartsPieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  AreaChart,
  Area,
  ComposedChart,
} from 'recharts';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import apiService from '../services/api';
import { InventoryItem, Delivery, InventorySnapshot } from '../types/api';
import ItemComparison from './ItemComparison';

// Types for analytics data
interface AnalyticsData {
  date: string;
  stockLevel: number;
  usage: number;
  cost: number;
  deliveries: number;
}

interface ComparisonData {
  item1: { name: string; data: AnalyticsData[] };
  item2: { name: string; data: AnalyticsData[] };
}

interface SummaryStats {
  totalItems: number;
  lowStockItems: number;
  totalValue: number;
  avgUtilization: number;
  topPerformingItems: Array<{ name: string; utilization: number }>;
  recentDeliveries: number;
}

// Mock data generators for demonstration
const generateAnalyticsData = (days: number, itemName: string): AnalyticsData[] => {
  const data: AnalyticsData[] = [];
  const today = new Date();
  
  for (let i = days - 1; i >= 0; i--) {
    const date = new Date(today);
    date.setDate(date.getDate() - i);
    
    const baseStock = 50 + Math.random() * 30;
    const usage = Math.random() * 10 + 2;
    const stockLevel = Math.max(0, baseStock - usage);
    const cost = usage * (Math.random() * 2 + 1);
    const deliveries = Math.random() > 0.8 ? 1 : 0;
    
    data.push({
      date: date.toLocaleDateString('en-US', { 
        month: 'short', 
        day: 'numeric' 
      }),
      stockLevel: Math.round(stockLevel),
      usage: Math.round(usage * 10) / 10,
      cost: Math.round(cost * 100) / 100,
      deliveries,
    });
  }
  
  return data;
};

const generateComparisonData = (days: number): ComparisonData => {
  return {
    item1: {
      name: 'Tomatoes',
      data: generateAnalyticsData(days, 'Tomatoes')
    },
    item2: {
      name: 'Onions',
      data: generateAnalyticsData(days, 'Onions')
    }
  };
};

const generateSummaryStats = (): SummaryStats => {
  return {
    totalItems: 45,
    lowStockItems: 8,
    totalValue: 12500,
    avgUtilization: 78.5,
    topPerformingItems: [
      { name: 'Tomatoes', utilization: 95.2 },
      { name: 'Onions', utilization: 88.7 },
      { name: 'Potatoes', utilization: 82.1 },
      { name: 'Carrots', utilization: 79.3 },
      { name: 'Lettuce', utilization: 76.8 },
    ],
    recentDeliveries: 12,
  };
};

// Custom tooltip component
const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <Box
        sx={{
          backgroundColor: 'background.paper',
          border: '1px solid',
          borderColor: 'divider',
          borderRadius: 1,
          p: 2,
          boxShadow: 3,
        }}
      >
        <Typography variant="body2" color="text.secondary">
          Date: {label}
        </Typography>
        {payload.map((entry: any, index: number) => (
          <Typography key={index} variant="body2" color={entry.color}>
            {entry.name}: {entry.value}
          </Typography>
        ))}
      </Box>
    );
  }
  return null;
};

// Analytics Dashboard Component
const Analytics: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // State for filters and data
  const [timeRange, setTimeRange] = useState<'7d' | '30d' | '90d'>('30d');
  const [chartType, setChartType] = useState<'line' | 'bar' | 'area' | 'composed'>('line');
  const [selectedItems, setSelectedItems] = useState<string[]>([]);
  const [inventoryItems, setInventoryItems] = useState<InventoryItem[]>([]);
  const [analyticsData, setAnalyticsData] = useState<AnalyticsData[]>([]);
  const [comparisonData, setComparisonData] = useState<ComparisonData | null>(null);
  const [summaryStats, setSummaryStats] = useState<SummaryStats | null>(null);
  const [comparisonDialogOpen, setComparisonDialogOpen] = useState(false);

  // Load data
  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Load inventory items
      const items = await apiService.getInventoryItems();
      setInventoryItems(items);
      
      // Generate mock analytics data
      const days = timeRange === '7d' ? 7 : timeRange === '30d' ? 30 : 90;
      const data = generateAnalyticsData(days, 'Overall');
      setAnalyticsData(data);
      
      // Generate comparison data
      const comparison = generateComparisonData(days);
      setComparisonData(comparison);
      
      // Generate summary stats
      const stats = generateSummaryStats();
      setSummaryStats(stats);
      
    } catch (err) {
      setError('Failed to load analytics data. Please try again.');
      console.error('Error loading analytics data:', err);
    } finally {
      setLoading(false);
    }
  }, [timeRange]);

  useEffect(() => {
    loadData();
  }, [loadData]);

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

  const handleRefresh = () => {
    loadData();
  };

  const handleItemSelection = (itemName: string) => {
    setSelectedItems(prev => 
      prev.includes(itemName) 
        ? prev.filter(item => item !== itemName)
        : [...prev, itemName]
    );
  };

  const handleOpenComparison = () => {
    setComparisonDialogOpen(true);
  };

  const handleCloseComparison = () => {
    setComparisonDialogOpen(false);
  };

  const getChartComponent = () => {
    const commonProps = {
      data: analyticsData,
      margin: { top: 20, right: 30, left: 20, bottom: 5 },
    };

    switch (chartType) {
      case 'line':
        return (
          <LineChart {...commonProps}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis yAxisId="left" />
            <YAxis yAxisId="right" orientation="right" />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            <Line yAxisId="left" type="monotone" dataKey="stockLevel" stroke="#1976d2" name="Stock Level" />
            <Line yAxisId="left" type="monotone" dataKey="usage" stroke="#dc004e" name="Usage" />
            <Line yAxisId="right" type="monotone" dataKey="cost" stroke="#4caf50" name="Cost ($)" />
          </LineChart>
        );
      
      case 'bar':
        return (
          <RechartsBarChart {...commonProps}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            <Bar dataKey="stockLevel" fill="#1976d2" name="Stock Level" />
            <Bar dataKey="usage" fill="#dc004e" name="Usage" />
          </RechartsBarChart>
        );
      
      case 'area':
        return (
          <AreaChart {...commonProps}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            <Area type="monotone" dataKey="stockLevel" stackId="1" stroke="#1976d2" fill="#1976d2" name="Stock Level" />
            <Area type="monotone" dataKey="usage" stackId="2" stroke="#dc004e" fill="#dc004e" name="Usage" />
          </AreaChart>
        );
      
      case 'composed':
        return (
          <ComposedChart {...commonProps}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis yAxisId="left" />
            <YAxis yAxisId="right" orientation="right" />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            <Bar yAxisId="left" dataKey="usage" fill="#dc004e" name="Usage" />
            <Line yAxisId="right" type="monotone" dataKey="cost" stroke="#4caf50" name="Cost ($)" />
          </ComposedChart>
        );
      
      default:
        return null;
    }
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
          <AnalyticsIcon sx={{ mr: 2 }} />
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Analytics Dashboard
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
      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        {/* Error Alert */}
        {error && (
          <Alert severity="error" sx={{ mb: 3 }} action={
            <Button color="inherit" size="small" onClick={handleRefresh}>
              Retry
            </Button>
          }>
            {error}
          </Alert>
        )}

        {/* Controls */}
        <Paper sx={{ p: 3, mb: 3 }}>
          <Grid container spacing={3} alignItems="center">
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Time Range</InputLabel>
                <Select
                  value={timeRange}
                  label="Time Range"
                  onChange={(e) => setTimeRange(e.target.value as '7d' | '30d' | '90d')}
                >
                  <MenuItem value="7d">Last 7 Days</MenuItem>
                  <MenuItem value="30d">Last 30 Days</MenuItem>
                  <MenuItem value="90d">Last 90 Days</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Chart Type</InputLabel>
                <Select
                  value={chartType}
                  label="Chart Type"
                  onChange={(e) => setChartType(e.target.value as any)}
                >
                  <MenuItem value="line">Line Chart</MenuItem>
                  <MenuItem value="bar">Bar Chart</MenuItem>
                  <MenuItem value="area">Area Chart</MenuItem>
                  <MenuItem value="composed">Composed Chart</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            
            <Grid item xs={12} sm={6} md={3}>
              <Button
                fullWidth
                variant="outlined"
                startIcon={<Refresh />}
                onClick={handleRefresh}
                disabled={loading}
              >
                Refresh
              </Button>
            </Grid>
            
            <Grid item xs={12} sm={6} md={3}>
                              <Button
                  fullWidth
                  variant="contained"
                  startIcon={<Compare />}
                  onClick={handleOpenComparison}
                >
                  Compare Items
                </Button>
            </Grid>
          </Grid>
        </Paper>

        {/* Summary Statistics */}
        {summaryStats && (
          <Grid container spacing={3} sx={{ mb: 3 }}>
            <Grid item xs={12} sm={6} md={3}>
              <Card>
                <CardContent>
                  <Typography color="textSecondary" gutterBottom>
                    Total Items
                  </Typography>
                  <Typography variant="h4">
                    {summaryStats.totalItems}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    In inventory
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} sm={6} md={3}>
              <Card>
                <CardContent>
                  <Typography color="textSecondary" gutterBottom>
                    Low Stock Items
                  </Typography>
                  <Typography variant="h4" color="warning.main">
                    {summaryStats.lowStockItems}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Need reorder
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} sm={6} md={3}>
              <Card>
                <CardContent>
                  <Typography color="textSecondary" gutterBottom>
                    Total Value
                  </Typography>
                  <Typography variant="h4" color="success.main">
                    ${summaryStats.totalValue.toLocaleString()}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Inventory value
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} sm={6} md={3}>
              <Card>
                <CardContent>
                  <Typography color="textSecondary" gutterBottom>
                    Avg Utilization
                  </Typography>
                  <Typography variant="h4" color="primary">
                    {summaryStats.avgUtilization}%
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Across all items
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        )}

        {/* Main Chart */}
        <Grid container spacing={3}>
          <Grid item xs={12} lg={8}>
            <Card>
              <CardHeader
                title="Inventory Utilization Over Time"
                subheader={`Showing data for the last ${timeRange === '7d' ? '7 days' : timeRange === '30d' ? '30 days' : '90 days'}`}
                action={
                  <Box display="flex" gap={1}>
                    <ToggleButtonGroup
                      value={chartType}
                      exclusive
                      onChange={(e, value) => value && setChartType(value)}
                      size="small"
                    >
                      <ToggleButton value="line">
                        <ShowChart />
                      </ToggleButton>
                      <ToggleButton value="bar">
                        <BarChart />
                      </ToggleButton>
                      <ToggleButton value="area">
                        <PieChart />
                      </ToggleButton>
                      <ToggleButton value="composed">
                        <TableChart />
                      </ToggleButton>
                    </ToggleButtonGroup>
                  </Box>
                }
              />
              <CardContent>
                <Box sx={{ height: 400 }}>
                  <ResponsiveContainer width="100%" height="100%">
                    {getChartComponent()}
                  </ResponsiveContainer>
                </Box>
              </CardContent>
            </Card>
          </Grid>
          
          <Grid item xs={12} lg={4}>
            <Stack spacing={3}>
              {/* Top Performing Items */}
              <Card>
                <CardHeader title="Top Performing Items" />
                <CardContent>
                  {summaryStats?.topPerformingItems.map((item, index) => (
                    <Box key={item.name} sx={{ mb: 2 }}>
                      <Box display="flex" justifyContent="space-between" alignItems="center">
                        <Typography variant="body2">
                          {item.name}
                        </Typography>
                        <Typography variant="body2" color="primary">
                          {item.utilization}%
                        </Typography>
                      </Box>
                      <LinearProgress 
                        variant="determinate" 
                        value={item.utilization} 
                        sx={{ mt: 1 }}
                      />
                    </Box>
                  ))}
                </CardContent>
              </Card>
              
              {/* Quick Actions */}
              <Card>
                <CardHeader title="Quick Actions" />
                <CardContent>
                  <Stack spacing={2}>
                    <Button
                      fullWidth
                      variant="outlined"
                      onClick={() => navigate('/inventory')}
                    >
                      View Inventory
                    </Button>
                    <Button
                      fullWidth
                      variant="outlined"
                      onClick={() => navigate('/dashboard')}
                    >
                      Back to Dashboard
                    </Button>
                  </Stack>
                </CardContent>
              </Card>
            </Stack>
          </Grid>
        </Grid>

        {/* Comparison Chart */}
        {comparisonData && (
          <Card sx={{ mt: 3 }}>
            <CardHeader title="Item Comparison" />
            <CardContent>
              <Box sx={{ height: 400 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <ComposedChart data={comparisonData.item1.data}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="date" />
                    <YAxis yAxisId="left" />
                    <YAxis yAxisId="right" orientation="right" />
                    <Tooltip content={<CustomTooltip />} />
                    <Legend />
                    <Bar yAxisId="left" dataKey="stockLevel" fill="#1976d2" name={`${comparisonData.item1.name} Stock`} />
                    <Line yAxisId="right" type="monotone" dataKey="usage" stroke="#dc004e" name={`${comparisonData.item2.name} Usage`} />
                  </ComposedChart>
                </ResponsiveContainer>
              </Box>
            </CardContent>
          </Card>
                 )}

        {/* Item Comparison Dialog */}
        <ItemComparison
          open={comparisonDialogOpen}
          onClose={handleCloseComparison}
        />
      </Container>
    </Box>
  );
};

export default Analytics; 