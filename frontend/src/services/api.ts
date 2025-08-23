import axios, { AxiosInstance, AxiosResponse } from 'axios';
import {
  User,
  InventoryItem,
  MenuItem,
  Delivery,
  InventorySnapshot,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  CreateInventoryItemRequest,
  UpdateInventoryItemRequest,
  CreateMenuItemRequest,
  CreateDeliveryRequest,
  StandardAPIResponse,
  Category,
  CreateCategoryRequest,
  UpdateCategoryRequest
} from '../types/api';

// API service class for communicating with the Go backend
// Handles authentication, requests, and error handling
class ApiService {
  private api: AxiosInstance;
  private baseURL: string;

  constructor() {
    // Use environment variable or default to localhost
    this.baseURL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

    this.api = axios.create({
      baseURL: this.baseURL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add request interceptor to include auth token
    this.api.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Add response interceptor for error handling
    this.api.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Token expired or invalid, redirect to login
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          window.dispatchEvent(new Event('auth-error'));
        }
        return Promise.reject(error);
      }
    );
  }

  // Authentication methods
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response: AxiosResponse<AuthResponse> = await this.api.post('/auth/login', credentials);
    return response.data;
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    const response: AxiosResponse<AuthResponse> = await this.api.post('/auth/register', userData);
    return response.data;
  }

  async getCurrentUser(): Promise<User> {
    const response: AxiosResponse<User> = await this.api.get('/api/v1/me');
    return response.data;
  }

  //fetch categories
  async getActiveCategories(): Promise<Category[]> {
    const response: AxiosResponse<StandardAPIResponse<Category[]>> = await this.api.get('/api/v1/categories/active');
    return response.data.data || [];
  }

    async getCategories(): Promise<Category[]> {
    const response: AxiosResponse<StandardAPIResponse<Category[]>> = await this.api.get('/api/v1/categories');
    return response.data.data || [];
  }

  async createCategory(categoryData: CreateCategoryRequest): Promise<Category> {
    const response: AxiosResponse<StandardAPIResponse<Category>> = await this.api.post('/api/v1/categories', categoryData);
    return response.data.data;
  }

  async updateCategory(id: number, categoryData: UpdateCategoryRequest): Promise<Category> {
    const response: AxiosResponse<StandardAPIResponse<Category>> = await this.api.put(`/api/v1/categories/${id}`, categoryData);
    return response.data.data;
  }

  async deleteCategory(id: number): Promise<void> {
    await this.api.delete(`/api/v1/categories/${id}`);
  }

  async getInventoryItems(): Promise<InventoryItem[]> {
    try {
      const response: AxiosResponse<StandardAPIResponse<InventoryItem[]>> = await this.api.get('/api/v1/inventory/items');
      return response.data.data || [];
    } catch (error) {
      console.error("Error fetching inventory items:", error);
      return [];
    }
  }

  async getLowStockItems(): Promise<InventoryItem[]> {
    const response: AxiosResponse<StandardAPIResponse<InventoryItem[]>> = await this.api.get('/api/v1/inventory/items/low-stock');
    return response.data.data || [];
  }

  async getInventoryItem(id: number): Promise<InventoryItem> {
    const response: AxiosResponse<InventoryItem> = await this.api.get(`/api/v1/inventory/items/${id}`);
    return response.data;
  }

  async createInventoryItem(item: CreateInventoryItemRequest): Promise<InventoryItem> {
    const response: AxiosResponse<InventoryItem> = await this.api.post('/api/v1/inventory/items', item);
    return response.data;
  }

  async updateInventoryItem(id: number, item: UpdateInventoryItemRequest): Promise<InventoryItem> {
    const response: AxiosResponse<InventoryItem> = await this.api.put(`/api/v1/inventory/items/${id}`, item);
    return response.data;
  }

  async deleteInventoryItem(id: number): Promise<void> {
    await this.api.delete(`/api/v1/inventory/items/${id}`);
  }

  async getInventoryItemsByVendor(vendor: string): Promise<InventoryItem[]> {
    const response: AxiosResponse<InventoryItem[]> = await this.api.get(`/api/v1/inventory/vendor/${vendor}`);
    return response.data;
  }

  // Menu methods
  async getMenuItems(): Promise<MenuItem[]> {
    const response: AxiosResponse<StandardAPIResponse<MenuItem[]>> = await this.api.get('/api/v1/menu/items');
    return response.data.data || [];
  }

  async createMenuItem(item: CreateMenuItemRequest): Promise<MenuItem> {
    const response: AxiosResponse<MenuItem> = await this.api.post('/api/v1/menu/items', item);
    return response.data;
  }

  // Delivery methods
  async getDeliveries(vendor?: string): Promise<Delivery[]> {
    // Start with the base URL
    let url = '/api/v1/deliveries';

    // If a vendor search term is provided, add it as a query parameter
    if (vendor && vendor.trim() !== '') {
      url += `?vendor=${encodeURIComponent(vendor)}`;
    }

    const response: AxiosResponse<StandardAPIResponse<Delivery[]>> = await this.api.get(url);
    return response.data.data || [];
  }

  async logDelivery(delivery: CreateDeliveryRequest): Promise<Delivery> {
    const response: AxiosResponse<Delivery> = await this.api.post('/api/v1/deliveries', delivery);
    return response.data;
  }

  async getDeliveriesByVendor(vendor: string): Promise<Delivery[]> {
    const response: AxiosResponse<Delivery[]> = await this.api.get(`/api/v1/deliveries/vendor/${vendor}`);
    return response.data;
  }

  // Snapshot methods
  async getInventorySnapshots(): Promise<InventorySnapshot[]> {
    const response: AxiosResponse<InventorySnapshot[]> = await this.api.get('/api/v1/snapshots');
    return response.data;
  }

  async createInventorySnapshot(snapshot: { counts: Record<number, number> }): Promise<InventorySnapshot> {
    const response: AxiosResponse<InventorySnapshot> = await this.api.post('/api/v1/snapshots', snapshot);
    return response.data;
  }

  // async getEmailSchedules(accountId: number): Promise<{ schedules: EmailSchedule[] }> {
  //   const response: AxiosResponse<{ schedules: EmailSchedule[] }> = await this.api.get(`/accounts/${accountId}/email-schedules`);
  //   return response.data;
  // }

  // Utility methods
  setAuthToken(token: string): void {
    localStorage.setItem('token', token);
  }

  getAuthToken(): string | null {
    return localStorage.getItem('token');
  }

  clearAuthToken(): void {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }

  isAuthenticated(): boolean {
    return !!this.getAuthToken();
  }
}

// Export singleton instance
export const apiService = new ApiService();
export default apiService; 