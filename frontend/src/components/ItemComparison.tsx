import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  CardHeader,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Grid,
  Alert,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Checkbox,
  Divider,
  Stack,
  LinearProgress,
} from '@mui/material';
import {
  Compare,
  Close,
  Refresh,
  TrendingUp,
  TrendingDown,
  Remove,
  BarChart,
  ShowChart,
  PieChart,
} from '@mui/icons-material';
import {
  LineChart,
  Line,
  BarChart as RechartsBarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  ComposedChart,
  AreaChart,
  Area,
} from 'recharts';
import { InventoryItem } from '../types/api';
import apiService from '../services/api';

// Types for comparison data
interface ComparisonItem {
  item: InventoryItem;
  data: Array<{
    date: string;
    stockLevel: number;
    usage: number;
    cost: number;
    utilization: number;
  }>;
  color: string;
}

interface ComparisonProps {
  open: boolean;
  onClose: () => void;
}

// Color palette for different items
const COLORS = [
  '#1976d2', '#dc004e', '#4caf50', '#ff9800', '#9c27b0',
  '#f44336', '#2196f3', '#ff5722', '#795548', '#607d8b'
];

// Mock data generator for item comparison
const generateItemData = (item: InventoryItem, days: number): ComparisonItem['data'] => {
  const data: ComparisonItem['data'] = [];
  const today = new Date();
  
  for (let i = days - 1; i >= 0; i--) {
    const date = new Date(today);
    date.setDate(date.getDate() - i);
    
    const baseStock = item.max_stock_level;
    const usage = Math.random() * (baseStock * 0.3) + (baseStock * 0.1);
    const stockLevel = Math.max(0, baseStock - usage);
    const cost = usage * item.cost_per_unit;
    const utilization = (usage / baseStock) * 100;
    
    data.push({
      date: date.toLocaleDateString('en-US', { 
        month: 'short', 
        day: 'numeric' 
      }),
      stockLevel: Math.round(stockLevel),
      usage: Math.round(usage * 10) / 10,
      cost: Math.round(cost * 100) / 100,
      utilization: Math.round(utilization * 10) / 10,
    });
  }
  
  return data;
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

// Calculate trend for an item
const calculateTrend = (data: ComparisonItem['data']) => {
  if (data.length < 3) return { icon: <Remove />, color: 'default', text: 'Insufficient data' };
  
  const recent = data.slice(-3).reduce((sum, item) => sum + item.utilization, 0) / 3;
  const older = data.slice(-6, -3).reduce((sum, item) => sum + item.utilization, 0) / 3;
  
  if (recent > older * 1.1) {
    return { icon: <TrendingUp />, color: 'success', text: 'Increasing' };
  } else if (recent < older * 0.9) {
    return { icon: <TrendingDown />, color: 'error', text: 'Decreasing' };
  } else {
    return { icon: <Remove />, color: 'default', text: 'Stable' };
  }
};

const ItemComparison: React.FC<ComparisonProps> = ({ open, onClose }) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [inventoryItems, setInventoryItems] = useState<InventoryItem[]>([]);
  const [selectedItems, setSelectedItems] = useState<InventoryItem[]>([]);
  const [comparisonData, setComparisonData] = useState<ComparisonItem[]>([]);
  const [timeRange, setTimeRange] = useState<'7d' | '30d' | '90d'>('30d');
  const [chartType, setChartType] = useState<'line' | 'bar' | 'area' | 'composed'>('line');
  const [metric, setMetric] = useState<'utilization' | 'stockLevel' | 'usage' | 'cost'>('utilization');

  // Load inventory items
  const loadInventoryItems = useCallback(async () => {
    try {
      setLoading(true);
      const items = await apiService.getInventoryItems();
      setInventoryItems(items);
    } catch (err) {
      setError('Failed to load inventory items');
      console.error('Error loading inventory items:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (open) {
      loadInventoryItems();
    }
  }, [open, loadInventoryItems]);

  // Generate comparison data when selected items change
  useEffect(() => {
    if (selectedItems.length === 0) {
      setComparisonData([]);
      return;
    }

    const days = timeRange === '7d' ? 7 : timeRange === '30d' ? 30 : 90;
    const data: ComparisonItem[] = selectedItems.map((item, index) => ({
      item,
      data: generateItemData(item, days),
      color: COLORS[index % COLORS.length],
    }));
    
    setComparisonData(data);
  }, [selectedItems, timeRange]);

  const handleItemToggle = (item: InventoryItem) => {
    setSelectedItems(prev => {
      const isSelected = prev.some(selected => selected.id === item.id);
      if (isSelected) {
        return prev.filter(selected => selected.id !== item.id);
      } else {
        return [...prev, item];
      }
    });
  };

  const handleClearSelection = () => {
    setSelectedItems([]);
  };

  const getChartComponent = () => {
    if (comparisonData.length === 0) return <></>;;

    const commonProps = {
      data: comparisonData[0].data,
      margin: { top: 20, right: 30, left: 20, bottom: 5 },
    };

    switch (chartType) {
      case 'line':
        return (
          <LineChart {...commonProps}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            {comparisonData.map((item, index) => (
              <Line
                key={item.item.id}
                type="monotone"
                dataKey={metric}
                stroke={item.color}
                name={item.item.name}
                strokeWidth={2}
              />
            ))}
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
            {comparisonData.map((item, index) => (
              <Bar
                key={item.item.id}
                dataKey={metric}
                fill={item.color}
                name={item.item.name}
              />
            ))}
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
            {comparisonData.map((item, index) => (
              <Area
                key={item.item.id}
                type="monotone"
                dataKey={metric}
                stroke={item.color}
                fill={item.color}
                name={item.item.name}
                fillOpacity={0.3}
              />
            ))}
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
            {comparisonData.map((item, index) => (
              <React.Fragment key={item.item.id}>
                <Bar
                  yAxisId="left"
                  dataKey="stockLevel"
                  fill={item.color}
                  name={`${item.item.name} Stock`}
                />
                <Line
                  yAxisId="right"
                  type="monotone"
                  dataKey="usage"
                  stroke={item.color}
                  name={`${item.item.name} Usage`}
                />
              </React.Fragment>
            ))}
          </ComposedChart>
        );
      
      default:
        return <></>;
    }
  };

  const getMetricLabel = () => {
    switch (metric) {
      case 'utilization': return 'Utilization (%)';
      case 'stockLevel': return 'Stock Level';
      case 'usage': return 'Usage';
      case 'cost': return 'Cost ($)';
      default: return '';
    }
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="xl"
      fullWidth
      PaperProps={{
        sx: { height: '90vh' }
      }}
    >
      <DialogTitle>
        <Box display="flex" alignItems="center" justifyContent="space-between">
          <Typography variant="h5">
            Item Comparison
          </Typography>
          <IconButton onClick={onClose}>
            <Close />
          </IconButton>
        </Box>
      </DialogTitle>
      
      <DialogContent dividers>
        <Grid container spacing={3}>
          {/* Controls */}
          <Grid size={{xs:12}}>
            <Card>
              <CardContent>
                <Grid container spacing={2} alignItems="center">
                  <Grid size={{ xs: 12, sm: 6, md: 3 }}>
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
                  
                  <Grid size={{ xs: 12, sm: 6, md: 3 }}>
                    <FormControl fullWidth size="small">
                      <InputLabel>Metric</InputLabel>
                      <Select
                        value={metric}
                        label="Metric"
                        onChange={(e) => setMetric(e.target.value as any)}
                      >
                        <MenuItem value="utilization">Utilization</MenuItem>
                        <MenuItem value="stockLevel">Stock Level</MenuItem>
                        <MenuItem value="usage">Usage</MenuItem>
                        <MenuItem value="cost">Cost</MenuItem>
                      </Select>
                    </FormControl>
                  </Grid>
                  
                  <Grid size={{ xs: 12, sm: 6, md: 3 }}>
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
                  
                  <Grid size={{ xs: 12, sm: 6, md: 3 }}>
                    <Button
                      fullWidth
                      variant="outlined"
                      onClick={handleClearSelection}
                      disabled={selectedItems.length === 0}
                    >
                      Clear Selection
                    </Button>
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>

          {/* Error Alert */}
          {error && (
            <Grid size={{xs:12}}>
              <Alert severity="error" action={
                <Button color="inherit" size="small" onClick={loadInventoryItems}>
                  Retry
                </Button>
              }>
                {error}
              </Alert>
            </Grid>
          )}

          {/* Item Selection */}
          <Grid size={{xs:12, md:4}}>
            <Card>
              <CardHeader title="Select Items to Compare" />
              <CardContent>
                {loading ? (
                  <Box display="flex" justifyContent="center" p={2}>
                    <CircularProgress />
                  </Box>
                ) : (
                  <List>
                    {inventoryItems?.map((item) => {
                      const isSelected = selectedItems.some(selected => selected.id === item.id);
                      return (
                        <ListItem key={item.id} dense>
                          <Checkbox
                            checked={isSelected}
                            onChange={() => handleItemToggle(item)}
                          />
                          <ListItemText
                            primary={item.name}
                            secondary={`Stock: ${item.current_stock || 0}/${item.max_stock_level}`}
                          />
                          <ListItemSecondaryAction>
                            <Chip
                              size="small"
                              label={`$${item.cost_per_unit}`}
                              color="primary"
                              variant="outlined"
                            />
                          </ListItemSecondaryAction>
                        </ListItem>
                      );
                    })}
                  </List>
                )}
              </CardContent>
            </Card>
          </Grid>

          {/* Chart */}
          <Grid size={{xs:12, md:4}}>
            <Card>
              <CardHeader
                title={`${getMetricLabel()} Comparison`}
                subheader={`Comparing ${selectedItems.length} items over ${timeRange === '7d' ? '7 days' : timeRange === '30d' ? '30 days' : '90 days'}`}
              />
              <CardContent>
                {selectedItems.length === 0 ? (
                  <Box display="flex" justifyContent="center" alignItems="center" height={400}>
                    <Typography variant="body1" color="text.secondary">
                      Select items to compare
                    </Typography>
                  </Box>
                ) : (
                  <Box sx={{ height: 400 }}>
                    <ResponsiveContainer width="100%" height="100%">
                      {getChartComponent()}
                    </ResponsiveContainer>
                  </Box>
                )}
              </CardContent>
            </Card>
          </Grid>

          {/* Selected Items Summary */}
          {selectedItems.length > 0 && (
            <Grid size={{xs:12}}>
              <Card>
                <CardHeader title="Selected Items Summary" />
                <CardContent>
                  <Grid container spacing={2}>
                    {comparisonData.map((item, index) => {
                      const trend = calculateTrend(item.data);
                      const avgUtilization = item.data.reduce((sum, d) => sum + d.utilization, 0) / item.data.length;
                      
                      return (
                        <Grid size={{ xs: 12, sm: 6, md: 3 }} key={item.item.id}>
                          <Card variant="outlined">
                            <CardContent>
                              <Box display="flex" alignItems="center" justifyContent="space-between" mb={1}>
                                <Typography variant="h6" color="primary">
                                  {item.item.name}
                                </Typography>
                                <Chip
                                  icon={trend.icon}
                                  label={trend.text}
                                  color={trend.color as any}
                                  size="small"
                                />
                              </Box>
                              
                              <Box mb={2}>
                                <Typography variant="body2" color="text.secondary">
                                  Avg Utilization
                                </Typography>
                                <Typography variant="h6" color="primary">
                                  {avgUtilization.toFixed(1)}%
                                </Typography>
                              </Box>
                              
                              <LinearProgress
                                variant="determinate"
                                value={Math.min(avgUtilization, 100)}
                                sx={{ mb: 1 }}
                              />
                              
                              <Typography variant="body2" color="text.secondary">
                                Current Stock: {item.item.current_stock || 0} / {item.item.max_stock_level}
                              </Typography>
                            </CardContent>
                          </Card>
                        </Grid>
                      );
                    })}
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
          )}
        </Grid>
      </DialogContent>
      
      <DialogActions>
        <Button onClick={onClose}>
          Close
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ItemComparison; 