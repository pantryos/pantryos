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
  ApiResponse
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
          window.location.href = '/login';
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

  // Inventory methods
  async getInventoryItems(): Promise<InventoryItem[]> {
    const response: AxiosResponse<InventoryItem[]> = await this.api.get('/api/v1/inventory/items');
    return response.data;
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
    const response: AxiosResponse<MenuItem[]> = await this.api.get('/api/v1/menu/items');
    return response.data;
  }

  async createMenuItem(item: CreateMenuItemRequest): Promise<MenuItem> {
    const response: AxiosResponse<MenuItem> = await this.api.post('/api/v1/menu/items', item);
    return response.data;
  }

  // Delivery methods
  async getDeliveries(): Promise<Delivery[]> {
    const response: AxiosResponse<Delivery[]> = await this.api.get('/api/v1/deliveries');
    return response.data;
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