import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Typography,
  Alert,
  CircularProgress,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
} from '@mui/material';
import { Analytics, Close, Refresh } from '@mui/icons-material';
import ItemUsageChart, { UsageData } from './ItemUsageChart';
import { InventoryItem } from '../types/api';

interface ItemUsageAnalyticsProps {
  item: InventoryItem;
  onClose?: () => void;
}

// Mock data generator for demonstration
const generateMockUsageData = (days: number): UsageData[] => {
  const data: UsageData[] = [];
  const today = new Date();
  
  for (let i = days - 1; i >= 0; i--) {
    const date = new Date(today);
    date.setDate(date.getDate() - i);
    
    // Generate realistic usage data with some variation
    const baseQuantity = Math.random() * 10 + 5; // 5-15 units
    const quantity = Math.round(baseQuantity + (Math.random() - 0.5) * 3);
    const cost = quantity * (Math.random() * 2 + 1); // $1-3 per unit
    
    data.push({
      date: date.toLocaleDateString('en-US', { 
        month: 'short', 
        day: 'numeric' 
      }),
      quantity: Math.max(0, quantity),
      cost: Math.round(cost * 100) / 100,
    });
  }
  
  return data;
};

const ItemUsageAnalytics: React.FC<ItemUsageAnalyticsProps> = ({ item, onClose }) => {
  const [timeRange, setTimeRange] = useState<'7d' | '30d' | '90d'>('30d');
  const [usageData, setUsageData] = useState<UsageData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  // Fetch usage data
  const fetchUsageData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      // TODO: Replace with actual API call when backend supports it
      // const data = await apiService.getItemUsage(item.id, timeRange);
      
      // For now, use mock data
      const days = timeRange === '7d' ? 7 : timeRange === '30d' ? 30 : 90;
      const mockData = generateMockUsageData(days);
      
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 500));
      
      setUsageData(mockData);
    } catch (err) {
      setError('Failed to load usage data. Please try again.');
      console.error('Error fetching usage data:', err);
    } finally {
      setLoading(false);
    }
  }, [timeRange, item.id]);

  useEffect(() => {
    fetchUsageData();
  }, [timeRange, item.id, fetchUsageData]);

  const handleTimeRangeChange = (range: '7d' | '30d' | '90d') => {
    setTimeRange(range);
  };

  const handleRefresh = () => {
    fetchUsageData();
  };

  const handleOpenDialog = () => {
    setIsDialogOpen(true);
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false);
    onClose?.();
  };

  const chartContent = (
    <Box sx={{ height: '100%', minHeight: 500 }}>
      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="100%">
          <CircularProgress />
        </Box>
      ) : error ? (
        <Alert severity="error" action={
          <Button color="inherit" size="small" onClick={handleRefresh}>
            Retry
          </Button>
        }>
          {error}
        </Alert>
      ) : (
        <ItemUsageChart
          title="Item Usage Analytics"
          data={usageData}
          itemName={item.name}
          currentStock={item.current_stock || 0}
          maxStock={item.max_stock_level}
          minStock={item.min_stock_level}
          timeRange={timeRange}
          onTimeRangeChange={handleTimeRangeChange}
        />
      )}
    </Box>
  );

  return (
    <>
      {/* Inline view */}
      <Box sx={{ mb: 2 }}>
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
          <Typography variant="h6" component="h2">
            Usage Analytics
          </Typography>
          <Box>
            <IconButton onClick={handleRefresh} disabled={loading} aria-label="refresh">
              <Refresh />
            </IconButton>
            <Button
              variant="outlined"
              startIcon={<Analytics />}
              onClick={handleOpenDialog}
              size="small"
            >
              Full View
            </Button>
          </Box>
        </Box>
        
        <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 2 }}>
          <Box sx={{ flex: { xs: '1', md: '2' } }}>
            {chartContent}
          </Box>
          
          <Box sx={{ flex: { xs: '1', md: '1' } }}>
            <Box sx={{ p: 2, bgcolor: 'background.paper', borderRadius: 1, height: '100%' }}>
              <Typography variant="h6" gutterBottom>
                Quick Stats
              </Typography>
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" color="text.secondary">
                  Stock Level
                </Typography>
                <Typography variant="h4" color="primary">
                  {item.current_stock || 0}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  of {item.max_stock_level} max
                </Typography>
              </Box>
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" color="text.secondary">
                  Reorder Point
                </Typography>
                <Typography variant="h6" color="warning.main">
                  {item.min_stock_level}
                </Typography>
              </Box>
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" color="text.secondary">
                  Cost per Unit
                </Typography>
                <Typography variant="h6" color="secondary">
                  ${item.cost_per_unit}
                </Typography>
              </Box>
              
              <Box>
                <Typography variant="body2" color="text.secondary">
                  Preferred Vendor
                </Typography>
                <Typography variant="body1">
                  {item.preferred_vendor}
                </Typography>
              </Box>
            </Box>
          </Box>
        </Box>
      </Box>

      {/* Full screen dialog */}
      <Dialog
        open={isDialogOpen}
        onClose={handleCloseDialog}
        maxWidth="xl"
        fullWidth
        PaperProps={{
          sx: { height: '90vh' }
        }}
      >
        <DialogTitle>
          <Box display="flex" alignItems="center" justifyContent="space-between">
            <Typography variant="h5">
              {item.name} - Usage Analytics
            </Typography>
            <IconButton onClick={handleCloseDialog}>
              <Close />
            </IconButton>
          </Box>
        </DialogTitle>
        
        <DialogContent dividers>
          {chartContent}
        </DialogContent>
        
        <DialogActions>
          <Button onClick={handleRefresh} disabled={loading}>
            Refresh Data
          </Button>
          <Button onClick={handleCloseDialog}>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default ItemUsageAnalytics; 