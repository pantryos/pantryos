import React, { useState, useEffect } from 'react';
import {
  Alert,
  AlertTitle,
  Box,
  Chip,
  Collapse,
  IconButton,
  List,
  ListItem,
  ListItemText,
  Typography,
  Button
} from '@mui/material';
import {
  Warning as WarningIcon,
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
  Inventory as InventoryIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import apiService from '../services/api';
import { InventoryItem } from '../types/api';

// Utility function to get stock status color
const getStockStatusColor = (item: InventoryItem): 'error' | 'warning' | 'info' => {
  if (item.current_stock === 0) {
    return 'error';
  }
  
  if (item.current_stock < item.min_stock_level) {
    return 'warning';
  }
  
  return 'info';
};

interface LowStockBannerProps {
  maxItems?: number; // Maximum number of items to show in the banner
}

const LowStockBanner: React.FC<LowStockBannerProps> = ({ maxItems = 5 }) => {
  const navigate = useNavigate();
  const [lowStockItems, setLowStockItems] = useState<InventoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState(false);
  const [error, setError] = useState('');

  // Load low stock items from API
  useEffect(() => {
    const loadLowStockItems = async () => {
      try {
        setLoading(true);
        const items = await apiService.getLowStockItems();
        
        // Sort by urgency (lowest stock first)
        const sortedItems = items.sort((a, b) => {
          return a.current_stock - b.current_stock;
        }).slice(0, maxItems);
        
        setLowStockItems(sortedItems);
      } catch (error) {
        console.error('Failed to load low stock items:', error);
        setError('Failed to load inventory data');
      } finally {
        setLoading(false);
      }
    };

    loadLowStockItems();
  }, [maxItems]);

  // Don't render if no low stock items
  if (loading || lowStockItems.length === 0) {
    return null;
  }

  const handleExpandClick = () => {
    setExpanded(!expanded);
  };

  const handleViewInventory = () => {
    navigate('/inventory');
  };

  const displayedItems = expanded ? lowStockItems : lowStockItems.slice(0, 3);
  const hasMoreItems = lowStockItems.length > 3;

  return (
    <Alert 
      severity="warning" 
      icon={<WarningIcon />}
      sx={{ 
        mb: 3,
        '& .MuiAlert-message': {
          width: '100%'
        }
      }}
    >
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', width: '100%' }}>
        <Box sx={{ flex: 1 }}>
          <AlertTitle>
            Low Stock Alert - {lowStockItems.length} item{lowStockItems.length !== 1 ? 's' : ''} need{lowStockItems.length !== 1 ? '' : 's'} attention
          </AlertTitle>
          
          <List dense sx={{ py: 0 }}>
            {displayedItems.map((item) => (
              <ListItem key={item.id} sx={{ py: 0.5, px: 0 }}>
                <ListItemText
                  primary={
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="body2" component="span">
                        {item.name}
                      </Typography>
                      <Chip
                        label={`${item.current_stock || 0} ${item.unit}`}
                        size="small"
                        color={getStockStatusColor(item)}
                        variant="outlined"
                      />
                      <Typography variant="caption" color="text.secondary">
                        (min: {item.min_stock_level} {item.unit})
                      </Typography>
                    </Box>
                  }
                  secondary={
                    <Typography variant="caption" color="text.secondary">
                      Vendor: {item.preferred_vendor || 'Not specified'}
                    </Typography>
                  }
                />
              </ListItem>
            ))}
          </List>
          
          {hasMoreItems && !expanded && (
            <Typography variant="caption" color="text.secondary">
              ... and {lowStockItems.length - 3} more items
            </Typography>
          )}
        </Box>
        
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, ml: 2 }}>
          {hasMoreItems && (
            <IconButton
              size="small"
              onClick={handleExpandClick}
              sx={{ alignSelf: 'flex-end' }}
            >
              {expanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
            </IconButton>
          )}
          
          <Button
            size="small"
            variant="outlined"
            startIcon={<InventoryIcon />}
            onClick={handleViewInventory}
            sx={{ 
              minWidth: 'auto',
              px: 2,
              py: 0.5,
              fontSize: '0.75rem'
            }}
          >
            View All
          </Button>
        </Box>
      </Box>
    </Alert>
  );
};

export default LowStockBanner; 