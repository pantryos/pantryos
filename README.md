# Stok - Coffee Shop Inventory Management System

A comprehensive inventory management system for coffee shops and restaurants, built with Go and Gin framework. The system provides multi-tenant support, inventory tracking, menu management, delivery logging, and variance reporting with Toast POS integration.

## Features

- **Multi-tenant Architecture**: Secure account-based data isolation with organization hierarchy
- **User Management**: JWT-based authentication and role-based authorization (user, manager, admin, org_admin)
- **Inventory Management**: Track inventory items, quantities, costs, and stock levels
- **Menu Management**: Create menu items with recipe ingredients and categories
- **Delivery Tracking**: Log incoming deliveries, costs, and vendor information
- **Inventory Snapshots**: Track inventory levels over time for historical analysis
- **Variance Reports**: Compare actual vs theoretical inventory usage (Toast POS integration planned)
- **RESTful API**: Clean, documented API endpoints with Swagger documentation
- **Comprehensive Testing**: Full test coverage with SQLite in-memory database for reliable testing

## Tech Stack

- **Backend**: Go 1.21+, Gin framework
- **Database**: PostgreSQL with GORM ORM (SQLite for testing)
- **Authentication**: JWT tokens with bcrypt password hashing
- **Testing**: Testify framework with comprehensive integration tests
- **Documentation**: Swagger/OpenAPI with interactive API documentation
- **Containerization**: Docker & Docker Compose

## Project Structure

```
stok/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point with server configuration
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP request handlers with comprehensive error handling
│   │   │   ├── auth_handler.go      # Authentication endpoints (register, login)
│   │   │   ├── auth_handler_test.go # Comprehensive auth handler tests
│   │   │   └── inventory_handler.go # Inventory management endpoints
│   │   ├── middleware/          # Authentication and authorization middleware
│   │   └── router.go            # Route definitions and API setup
│   ├── auth/
│   │   └── jwt.go               # JWT utilities with secure token generation
│   ├── database/
│   │   ├── database.go          # Database connection and configuration
│   │   ├── repositories.go      # Repository interfaces and implementations
│   │   ├── service.go           # Business logic service layer with validation
│   │   ├── test_setup.go        # Test database setup with SQLite in-memory
│   │   ├── database_test.go     # Database operation tests
│   │   └── service_integration_test.go # Comprehensive integration tests
│   └── models/
│       └── models.go            # Data models with GORM tags and JSON serialization
├── pkg/
│   └── utils/
│       └── password.go          # Password utilities with bcrypt hashing
├── docker-compose.yml           # Database setup with PostgreSQL
├── go.mod                       # Go module dependencies
└── README.md                    # This file
```

## Database Schema

### Core Tables

- **organizations**: Multi-tenant organization management
- **accounts**: Business locations within organizations
- **users**: User accounts with JWT authentication and role-based access
- **inventory_items**: Trackable inventory items with stock level management
- **menu_items**: Menu items with pricing and categories
- **recipe_ingredients**: Links menu items to inventory ingredients
- **deliveries**: Incoming inventory deliveries with vendor tracking
- **inventory_snapshots**: Point-in-time inventory counts for historical analysis
- **orders**: Purchase orders with approval workflow
- **order_items**: Individual items within orders
- **order_requests**: Request workflow for inventory items
- **request_items**: Items within order requests

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for database)
- PostgreSQL (or use Docker)

### 1. Clone and Setup

```bash
git clone <repository-url>
cd stok
go mod tidy
```

### 2. Start Database

```bash
docker-compose up -d postgres
```

This will start PostgreSQL on port 5432 with:
- Database: `stok_db`
- Username: `postgres`
- Password: `password`

### 3. Environment Variables

Create a `.env` file (optional, defaults are provided):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=stok_db
DB_SSLMODE=disable
PORT=8080
JWT_SECRET=your-secret-key
USE_RAMSQL=false
```

**Note**: The application now uses PostgreSQL by default. Ramsql is only available for testing.

### 4. Run the Application

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### 5. Access API Documentation

Once the server is running, you can access the interactive Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

This provides a complete interactive API documentation where you can:
- View all available endpoints
- Test API calls directly from the browser
- See request/response schemas
- Authenticate with JWT tokens

### 6. Run Tests

```bash
# Run all tests
go test ./...

# Run database tests specifically
go test ./internal/database/...

# Run auth handler tests
go test ./internal/api/handlers/...

# Run with verbose output
go test ./... -v
```

## API Documentation

### Interactive Swagger Documentation

The API includes comprehensive interactive documentation powered by Swagger/OpenAPI. Once the server is running, you can access it at:

```
http://localhost:8080/swagger/index.html
```

### Key Features of the Swagger Documentation:

- **Interactive Testing**: Test API endpoints directly from the browser
- **Authentication**: Built-in JWT token authentication for protected endpoints
- **Request/Response Examples**: See exact request and response formats
- **Error Codes**: Complete list of possible error responses
- **Model Schemas**: Detailed view of all data models

### API Endpoints Overview

#### Authentication
- `POST /auth/register` - Register new user with account validation
- `POST /auth/login` - User login with JWT token generation
- `GET /api/v1/me` - Get current user information (protected)

#### Inventory Management (Protected)
- `GET /api/v1/inventory/items` - List inventory items for account
- `POST /api/v1/inventory/items` - Create inventory item with validation
- `GET /api/v1/inventory/items/{id}` - Get specific inventory item
- `PUT /api/v1/inventory/items/{id}` - Update inventory item
- `DELETE /api/v1/inventory/items/{id}` - Delete inventory item
- `GET /api/v1/inventory/vendor/{vendor}` - Get items by vendor
- `GET /api/v1/inventory/low-stock` - Get items below minimum stock level

#### Menu Management (Protected)
- `GET /api/v1/menu/items` - List menu items for account
- `POST /api/v1/menu/items` - Create menu item with category
- `GET /api/v1/menu/items/category/{category}` - Get menu items by category

#### Deliveries (Protected)
- `GET /api/v1/deliveries` - List all deliveries for account
- `POST /api/v1/deliveries` - Log delivery with vendor and cost tracking
- `GET /api/v1/deliveries/vendor/{vendor}` - Get deliveries by vendor
- `GET /api/v1/deliveries/date-range` - Get deliveries within date range

#### Snapshots (Protected)
- `GET /api/v1/snapshots` - List all inventory snapshots for account
- `POST /api/v1/snapshots` - Create inventory snapshot with quantity mapping
- `GET /api/v1/snapshots/latest` - Get most recent inventory snapshot
- `GET /api/v1/snapshots/date-range` - Get snapshots within date range

## Development

### Database Migrations

The application uses GORM's auto-migration feature. Tables are automatically created when the application starts. The migration includes:

- All core tables with proper relationships
- Indexes for performance optimization
- JSON serialization for complex data types (CountsMap)
- Timestamps for audit trails

### Testing Strategy

The project includes comprehensive testing with multiple layers:

1. **Unit Tests**: Individual function testing
2. **Integration Tests**: Database operations and business logic
3. **Handler Tests**: HTTP endpoint testing with authentication
4. **Test Database**: SQLite in-memory database for fast, reliable tests

#### Test Features:
- **Isolated Test Environment**: Each test uses a fresh database instance
- **Helper Functions**: Reusable test data creation functions
- **Comprehensive Coverage**: Tests for success cases, error cases, and edge cases
- **Authentication Testing**: JWT token generation and validation testing

### Adding New Features

1. **Models**: Add new structs to `internal/models/models.go` with proper GORM tags
2. **Repository**: Add repository interface and implementation in `internal/database/repositories.go`
3. **Service**: Add business logic with validation to `internal/database/service.go`
4. **Handler**: Add HTTP handlers with proper error handling in `internal/api/handlers/`
5. **Routes**: Add routes in `internal/api/router.go`
6. **Tests**: Add comprehensive tests for new functionality
7. **Documentation**: Add Swagger annotations to handlers

### Code Quality

The codebase follows Go best practices:

- **Comprehensive Comments**: All functions and complex logic are documented
- **Error Handling**: Proper error handling with meaningful error messages
- **Validation**: Input validation at multiple layers (handler, service, repository)
- **Security**: Password hashing, JWT tokens, and access control
- **Testing**: High test coverage with both unit and integration tests

## Recent Improvements

### Fixed Issues:
1. **Database Test Setup**: Fixed function name inconsistencies (`setupTestDB` → `SetupTestDB`)
2. **JWT Token Generation**: Added default JWT secret for testing environment
3. **Password Hashing**: Fixed login tests to use properly hashed passwords
4. **Error Handling**: Improved error messages and HTTP status codes
5. **Access Control**: Enhanced organization access validation with proper return values

### Added Features:
1. **Comprehensive Comments**: Added detailed documentation throughout the codebase
2. **Enhanced Testing**: Improved test coverage and reliability
3. **Better Error Messages**: More descriptive error responses
4. **Security Improvements**: Proper password hashing and JWT configuration

## Docker

### Build and Run with Docker

```bash
# Build the application
docker build -t stok .

# Run with database
docker-compose up
```

### Database Management

Access pgAdmin at `http://localhost:8081`:
- Email: `admin@stok.com`
- Password: `admin`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite: `go test ./... -v`
6. Ensure all tests pass
7. Submit a pull request

## License

[Add your license here]

## Support

For questions or issues, please open a GitHub issue or contact the development team.