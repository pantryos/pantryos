# Low Stock Banner Implementation

## Overview
A banner component has been added to the frontend dashboard that displays menu items (inventory items) that are low on stock. The banner provides a quick overview of items that need attention and allows users to navigate to the inventory management page.

## Features Implemented

### 1. Backend Changes

#### New API Endpoints
- **GET `/api/v1/inventory/items`** - Now returns inventory items with current stock levels calculated from the latest inventory snapshot
- **GET `/api/v1/inventory/items/low-stock`** - Returns only inventory items that are currently low on stock

#### New Service Methods
- `GetInventoryItemsWithCurrentStock(accountID int)` - Retrieves inventory items with current stock levels from the latest snapshot
- `InventoryItemWithStock` struct - Extends the base inventory item with current stock information

#### Database Integration
- Uses existing `InventorySnapshot` model to calculate current stock levels
- Falls back to zero stock if no snapshots exist
- Maintains data integrity and security through account scoping

### 2. Frontend Changes

#### New Components
- **`LowStockBanner`** - Main banner component that displays low stock items
  - Shows warning alert with item details
  - Displays current stock vs minimum stock levels
  - Includes vendor information
  - Expandable list for multiple items
  - Navigation button to inventory page
  - Color-coded stock status (error for zero stock, warning for low stock)

#### Updated Components
- **`Dashboard`** - Now includes the low stock banner at the top
- **`Inventory`** - Added current stock column with visual indicators for low stock items
- **`api.ts`** - Added `getLowStockItems()` method for the new API endpoint
- **`api.ts` (types)** - Updated `InventoryItem` interface to include `current_stock` field

### 3. Key Features

#### Stock Level Calculation
- Current stock is calculated from the most recent inventory snapshot
- Items are considered "low stock" when `current_stock < min_stock_level`
- Zero stock items are highlighted with red indicators
- Low stock items are highlighted with yellow indicators

#### User Experience
- Banner only appears when there are low stock items
- Shows item name, current stock, unit, and minimum stock level
- Displays vendor information for each item
- Expandable interface for viewing multiple items
- Direct navigation to inventory management
- Responsive design that works on different screen sizes

#### Visual Indicators
- Warning alert styling for the banner
- Color-coded chips showing stock levels
- Highlighted cells in the inventory table for low stock items
- Expand/collapse functionality for long lists

## Technical Implementation

### Backend Architecture
```
InventoryHandler.GetInventoryItems() 
  → Service.GetInventoryItemsWithCurrentStock()
    → InventoryItemRepository.GetByAccountID()
    → InventorySnapshotRepository.GetLatestByAccountID()
    → Calculate current stock from snapshot.Counts map
```

### Frontend Architecture
```
Dashboard 
  → LowStockBanner
    → apiService.getLowStockItems()
    → Display items with stock status
    → Navigate to inventory on button click
```

### Data Flow
1. User visits dashboard
2. LowStockBanner component loads
3. API call to `/api/v1/inventory/items/low-stock`
4. Backend calculates current stock from latest snapshot
5. Returns filtered list of low stock items
6. Frontend displays items in banner with status indicators

## Usage

### For Users
- The banner automatically appears on the dashboard when there are low stock items
- Click "View All" to navigate to the inventory management page
- Use the expand/collapse button to see more items if there are many
- Check the inventory table for highlighted low stock cells

### For Developers
- The banner component is reusable and can be added to other pages
- The `maxItems` prop controls how many items are shown initially
- The API endpoints can be used independently for other features
- Stock level calculations are handled server-side for consistency

## Future Enhancements
- Add notifications for critical stock levels
- Implement automatic reorder suggestions
- Add stock level trends and forecasting
- Include cost impact of low stock situations
- Add bulk actions for reordering multiple items

## Testing
- Backend endpoints are tested through the existing test suite
- Frontend component logic is implemented but tests need setup configuration
- Manual testing confirms functionality works as expected
- API integration is verified through browser developer tools 