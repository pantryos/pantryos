import React from 'react';

// Mock for react-router-dom
export const useNavigate = jest.fn();
export const useLocation = jest.fn();
export const useParams = jest.fn();
export const useSearchParams = jest.fn();
export const Link = ({ children, to, ...props }: any) => React.createElement('a', { href: to, ...props }, children);
export const NavLink = ({ children, to, ...props }: any) => React.createElement('a', { href: to, ...props }, children);
export const BrowserRouter = ({ children }: any) => React.createElement('div', {}, children);
export const Routes = ({ children }: any) => React.createElement('div', {}, children);
export const Route = ({ children }: any) => React.createElement('div', {}, children);
export const Navigate = ({ to }: any) => React.createElement('div', { 'data-testid': 'navigate', 'data-to': to });
export const Outlet = () => React.createElement('div', { 'data-testid': 'outlet' }); 