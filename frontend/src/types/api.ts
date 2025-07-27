// TypeScript interfaces matching the Go backend models
// These ensure type safety when communicating with the API

export interface User {
  id: number;
  email: string;
  account_id: number;
  role: string;
  created_at: string;
}

export interface Account {
  id: number;
  organization_id: number;
  name: string;
  location: string;
  phone: string;
  email: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface InventoryItem {
  id: number;
  account_id: number;
  name: string;
  unit: string;
  cost_per_unit: number;
  preferred_vendor: string;
  min_stock_level: number;
  max_stock_level: number;
}

export interface MenuItem {
  id: number;
  account_id: number;
  name: string;
  price: number;
  category: string;
}

export interface Delivery {
  id: number;
  account_id: number;
  inventory_item_id: number;
  vendor: string;
  quantity: number;
  delivery_date: string;
  cost: number;
}

export interface InventorySnapshot {
  id: number;
  account_id: number;
  timestamp: string;
  counts: Record<number, number>; // Map of inventory item ID to quantity
}

export interface Order {
  id: number;
  organization_id: number;
  status: string;
  order_date: string;
  expected_date: string;
  total_cost: number;
  notes: string;
  created_by: number;
  approved_by?: number;
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  id: number;
  order_id: number;
  inventory_item_id: number;
  quantity: number;
  unit_cost: number;
  total_cost: number;
  vendor: string;
  notes: string;
}

export interface OrderRequest {
  id: number;
  organization_id: number;
  account_id: number;
  status: string;
  priority: string;
  request_date: string;
  needed_by: string;
  notes: string;
  created_by: number;
  approved_by?: number;
  created_at: string;
  updated_at: string;
}

export interface RequestItem {
  id: number;
  order_request_id: number;
  inventory_item_id: number;
  quantity: number;
  reason: string;
  priority: string;
}

// API Response types
export interface ApiResponse<T> {
  data?: T;
  message?: string;
  error?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  account_id: number;
}

export interface AuthResponse {
  token: string;
  user: User;
}

// Form types for creating/updating items
export interface CreateInventoryItemRequest {
  name: string;
  unit: string;
  cost_per_unit: number;
  preferred_vendor: string;
  min_stock_level: number;
  max_stock_level: number;
}

export interface UpdateInventoryItemRequest extends CreateInventoryItemRequest {
  id: number;
}

export interface CreateMenuItemRequest {
  name: string;
  price: number;
  category: string;
}

export interface CreateDeliveryRequest {
  inventory_item_id: number;
  vendor: string;
  quantity: number;
  delivery_date: string;
  cost: number;
} 