import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Box,
  MenuItem as MuiMenuItem,
  InputAdornment,
} from '@mui/material';
import { InventoryItem, CreateDeliveryRequest } from '../types/api';

// Define the props the component will accept
interface LogDeliveryProps {
  open: boolean;
  onClose: () => void;
  inventoryItems: InventoryItem[];
  onSubmit: (deliveryData: CreateDeliveryRequest) => Promise<void>;
}

const LogDelivery: React.FC<LogDeliveryProps> = ({ open, onClose, inventoryItems, onSubmit }) => {
  const [formData, setFormData] = React.useState<CreateDeliveryRequest>({
    inventory_item_id: 0,
    vendor: '',
    quantity: 0,
    delivery_date: new Date().toISOString().split('T')[0],
    cost: 0,
  });

  const handleFormChange = (field: keyof CreateDeliveryRequest, value: string | number) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleSubmit = () => {
    const submissionData = {
      ...formData,
      delivery_date: new Date(formData.delivery_date).toISOString(),
    };
    onSubmit(submissionData);
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Log New Delivery</DialogTitle>
      <DialogContent>
        <Box sx={{ pt: 1 }}>
          <TextField
            select
            fullWidth
            label="Inventory Item"
            value={formData.inventory_item_id}
            onChange={(e) => handleFormChange('inventory_item_id', parseInt(e.target.value))}
            margin="normal"
            required
          >
            <MuiMenuItem value={0} disabled>
              <em>Select an item</em>
            </MuiMenuItem>
            {inventoryItems.map((item) => (
              <MuiMenuItem key={item.id} value={item.id}>
                {item.name}
              </MuiMenuItem>
            ))}
          </TextField>
          <TextField
            fullWidth
            label="Vendor"
            value={formData.vendor}
            onChange={(e) => handleFormChange('vendor', e.target.value)}
            margin="normal"
          />
          <TextField
            fullWidth
            label="Quantity"
            type="number"
            value={formData.quantity}
            onChange={(e) => handleFormChange('quantity', parseFloat(e.target.value) || 0)}
            margin="normal"
            required
            inputProps={{ min: 0 }}
          />
          <TextField
            fullWidth
            label="Total Cost"
            type="number"
            value={formData.cost}
            onChange={(e) => handleFormChange('cost', parseFloat(e.target.value) || 0)}
            margin="normal"
            required
            InputProps={{ startAdornment: <InputAdornment position="start">$</InputAdornment> }}
            inputProps={{ min: 0, step: 0.01 }}
          />
          <TextField
            fullWidth
            label="Delivery Date"
            type="date"
            value={formData.delivery_date}
            onChange={(e) => handleFormChange('delivery_date', e.target.value)}
            margin="normal"
            InputLabelProps={{ shrink: true }}
          />
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button onClick={handleSubmit} variant="contained">
          Log Delivery
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default LogDelivery;