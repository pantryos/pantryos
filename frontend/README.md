# Stok Frontend

A modern React TypeScript frontend for the Stok inventory management system.

## Features

- **Modern UI**: Built with Material-UI (MUI) for a beautiful, responsive interface
- **Authentication**: Secure login/register system with JWT tokens
- **Inventory Management**: Full CRUD operations for inventory items
- **Dashboard**: Overview with statistics and quick actions
- **TypeScript**: Full type safety throughout the application
- **Responsive Design**: Works on desktop, tablet, and mobile devices

## Tech Stack

- **React 18** with TypeScript
- **Material-UI (MUI)** for UI components
- **React Router** for navigation
- **Axios** for API communication
- **MUI X Data Grid** for data tables

## Getting Started

### Prerequisites

- Node.js (v18 or higher)
- npm or yarn
- The Stok backend server running on `http://localhost:8080`

### Installation

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm start
   ```

4. Open [http://localhost:3000](http://localhost:3000) in your browser.

### Environment Variables

Create a `.env` file in the frontend directory to configure the API URL:

```env
REACT_APP_API_URL=http://localhost:8080
```

## Available Scripts

- `npm start` - Start the development server
- `npm run build` - Build the app for production
- `npm test` - Run tests
- `npm run eject` - Eject from Create React App (not recommended)

## Project Structure

```
src/
├── components/          # React components
│   ├── Dashboard.tsx    # Main dashboard
│   ├── Inventory.tsx    # Inventory management
│   ├── Login.tsx        # Login form
│   └── Register.tsx     # Registration form
├── contexts/            # React contexts
│   └── AuthContext.tsx  # Authentication context
├── services/            # API services
│   └── api.ts          # API client
├── types/               # TypeScript type definitions
│   └── api.ts          # API types
└── App.tsx             # Main app component
```

## Features Overview

### Authentication
- User registration with email, password, and account ID
- Secure login with JWT token storage
- Protected routes that redirect to login if not authenticated
- Automatic token refresh and logout on expiration

### Dashboard
- Overview statistics (inventory items, menu items, recent deliveries)
- Quick action buttons for common tasks
- Navigation to different sections of the app

### Inventory Management
- View all inventory items in a searchable data grid
- Add new inventory items with detailed information
- Edit existing items
- Delete items with confirmation
- Search by name or vendor
- Pagination and sorting

## API Integration

The frontend communicates with the Go backend through RESTful APIs:

- **Authentication**: `/auth/login`, `/auth/register`
- **Inventory**: `/api/v1/inventory/items`
- **Menu Items**: `/api/v1/menu/items`
- **Deliveries**: `/api/v1/deliveries`
- **Snapshots**: `/api/v1/snapshots`

All API calls include JWT authentication headers and handle errors gracefully.

## Development

### Adding New Components

1. Create a new component in the `src/components/` directory
2. Add TypeScript interfaces in `src/types/api.ts` if needed
3. Add API methods in `src/services/api.ts`
4. Add routes in `src/App.tsx`
5. Update navigation in the Dashboard component

### Styling

The app uses Material-UI's `sx` prop for styling. Follow the existing patterns for consistency.

### State Management

- Local component state for UI interactions
- React Context for global authentication state
- API calls for data fetching and persistence

## Building for Production

```bash
npm run build
```

This creates a `build` folder with optimized production files that can be served by any static file server.

## Contributing

1. Follow the existing code style and patterns
2. Add TypeScript types for all new features
3. Include error handling for API calls
4. Test the UI on different screen sizes
5. Add comments for complex logic

## Troubleshooting

### Common Issues

1. **API Connection Errors**: Ensure the backend server is running on the correct port
2. **Authentication Issues**: Check that JWT tokens are being sent correctly
3. **Build Errors**: Make sure all dependencies are installed

### Development Tips

- Use the browser's developer tools to debug API calls
- Check the console for TypeScript errors
- Use React Developer Tools for component debugging
