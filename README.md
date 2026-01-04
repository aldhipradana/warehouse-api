# Go RESTful Products API

This is a RESTful API built with Go, Gin, and GORM for managing products. It supports CRUD operations, filtering, searching, and pagination.

## Features

- **CRUD Operations**: Create, Read, Update, and Delete products.
- **Filtering**: Advanced filtering with operators, functions, and JSON-based filters.
- **Search**: Global search across specified fields.
- **Pagination**: Paginated responses with total counts.
- **Sorting**: Sort results by any field in ascending or descending order.
- **Relation Loading**: Eager load related data.
- **Authentication & Authorization**: JWT-based authentication with role-based access control (admin, user).
- **User Management**: Complete user registration, login, and profile management.
- **Action Logging**: Automatically logs all data-modifying requests (POST, PUT, DELETE) to daily log files with payload and query capture.

## Project Structure

```text
.gitignore
go.mod
main.go
config.toml         # Application configuration
config/
  config.go         # Configuration loader and structs
database/
  seed/
    main/
      main.go         # Seeder entry point
    product_seeder.go # Initial data seeder
    user_seeder.go    # User data seeder
docs/
  bruno.json
  environments/
    local.bru
  products/
    create-product.bru
    delete-product.bru
    get-product.bru
    list-product.bru
    update-product.bru
models/
  product.go        # Product model definition
  user.go           # User model with password hashing
middleware/
  auth.go           # JWT authentication middleware
  logger.go         # Action logger middleware
log/
  YYYY-MM-DD.log    # Daily action logs
restful/
  controller.go     # Generic controller logic
  interface.go
  scopes.go
routes/
  api.go            # Main route entry point
  auth.go           # Authentication routes (register, login)
  user.go           # User management routes
  product.go        # Product-specific routes
```

### Key Files

- **main.go**: Entry point of the application.
- **config.toml**: Application configuration file (server, database, JWT settings).
- **config/**: Contains configuration loader and structs for TOML parsing.
- **models/**: Contains data models (Product, User) with validation and hooks.
- **restful/**: Contains reusable components for controllers, filters, and database operations.
- **middleware/**: Contains authentication (JWT) and action logging middleware.
- **routes/**: Organizes API routes by feature (auth, users, products).
- **log/**: Stores daily log files capturing request details, queries, and payloads.
- **docs/**: Contains API documentation and test cases in .bru format.

## Getting Started

### Prerequisites

- Go 1.23 or higher
- SQLite (or any other database supported by GORM)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/aldhipradana/warehouse-api.git
   cd warehouse-api
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

4. (Optional) Seed the database:
   ```bash
   go run database/seed/main/main.go
   ```

5. The API will be available at http://localhost:8080.

### Building for Production

To build the application for deployment:

**For Linux (AMD64):**
```powershell
$env:GOOS="linux"
$env:GOARCH="amd64"
$env:CGO_ENABLED="0"
go build -trimpath -ldflags="-s -w" -o warehouse-api-linux .
```

**For Windows:**
```powershell
go build -trimpath -ldflags="-s -w" -o warehouse-api.exe .
```

**Build flags explained:**
- `CGO_ENABLED=0`: Disables CGO for static binary (no external dependencies)
- `-trimpath`: Removes file system paths from the binary
- `-ldflags="-s -w"`: Strips debug information to reduce binary size
- `-o warehouse-api-linux`: Output filename

The resulting binary is self-contained and can be deployed directly to your server.

### Deployment

After uploading the binary to your server, make it executable:

```bash
chmod u+x $HOME/www/warehouse-api/warehouse-api-linux
```

Then run the application:
```bash
cd $HOME/www/warehouse-api
./warehouse-api-linux
```

Make sure your `config.toml` file is in the same directory as the binary.

### Configuration

The application uses a `config.toml` file for configuration. Create or modify it in the root directory:

```toml
[server]
port = 8080
debug = true

[database]
driver = "sqlite"
path = "database/test.db"
# For PostgreSQL:
# driver = "postgres"
# host = "localhost"
# port = 5432
# name = "warehouse_db"
# user = "postgres"
# password = "password"

[jwt]
secret = "your-secret-key-change-this-in-production"
token_expiry_hours = 24
```

**Configuration Options:**

- **Server**:
  - `port`: Server port (default: 8080)
  - `debug`: Enable debug mode for detailed logs (default: true)

- **Database**:
  - `driver`: Database driver (`sqlite`, `postgres`, `mysql`)
  - `path`: Database file path (for SQLite)
  - `host`, `port`, `name`, `user`, `password`: Database connection details (for PostgreSQL/MySQL)

- **JWT**:
  - `secret`: Secret key for signing JWT tokens (change in production!)
  - `token_expiry_hours`: Token expiration time in hours (default: 24)

### API Endpoints

#### Authentication

| Method | Endpoint         | Description                | Auth Required |
|--------|------------------|----------------------------|---------------|
| POST   | /api/auth/register | Register a new user      | No            |
| POST   | /api/auth/login    | Login and get JWT token  | No            |
| GET    | /api/auth/me       | Get current user info    | Yes           |

#### Users

| Method | Endpoint         | Description                | Auth Required | Role Required |
|--------|------------------|----------------------------|---------------|---------------|
| GET    | /api/users       | List all users             | Yes           | Admin         |
| GET    | /api/users/:id   | Get a user by ID           | Yes           | Admin         |
| PUT    | /api/users/:id   | Update a user by ID        | Yes           | Any           |
| DELETE | /api/users/:id   | Delete a user by ID        | Yes           | Admin         |

#### Products

| Method | Endpoint         | Description                |
|--------|------------------|----------------------------|
| GET    | /api/products  | List products with filters |
| GET    | /api/products/:id | Get a product by ID       |
| POST   | /api/products  | Create a new product       |
| PUT    | /api/products/:id | Update a product by ID    |
| DELETE | /api/products/:id | Delete a product by ID    |

### Query Parameters for Listing

- **Pagination**:
  - page: Page number (default: 1)
  - limit: Items per page (default: 20)
- **Sorting**:
  - sort: Field to sort by (default: created_at)
  - order: Sort order (asc or desc, default: desc)
- **Search**:
  - q: Search term for global search across fields.
- **Filters**:
  - filter: JSON object for advanced filtering.

#### Filter Examples

- Simple Equality:
  ```
  ?filter={"status": "active"}
  ```
- Operators:
  ```
  ?filter={"price": {"operator": ">", "value": 100}}
  ```
- LIKE Search:
  ```
  ?filter={"name": {"function": "like", "value": "Product"}}
  ```
- IN Clause:
  ```
  ?filter={"status": {"function": "in", "value": "active,pending"}}
  ```
- BETWEEN:
  ```
  ?filter={"price": {"function": "between", "value": "100,500"}}
  ```

### Example Requests

#### Register a New User
```json
POST /api/auth/register
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Login
```json
POST /api/auth/login
{
  "email": "john@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Get Current User (Protected Route)
```bash
GET /api/auth/me
Headers:
  Authorization: Bearer <your-jwt-token>
```

#### Create Product
```json
POST /api/products
{
  "name": "New Product",
  "price": 99.99,
  "status": "active"
}
```

#### List Products
```json
GET /api/products?page=1&limit=10&filter={"status":"active"}&q=Sample
```

#### Update Product
```json
PUT /api/products/1
{
  "name": "Updated Product",
  "price": 149.99,
  "status": "inactive"
}
```

#### Delete Product
```json
DELETE /api/products/1
```

## Testing

The docs/ folder contains .bru files for testing the API using [Bruno](https://www.usebruno.com/), a lightweight API client.

1. Install Bruno.
2. Open the docs/ folder in Bruno.
3. Run the requests to test the API.

## Authentication

The API uses **JWT (JSON Web Tokens)** for authentication. After logging in or registering, you'll receive a token that must be included in the `Authorization` header for protected routes.

### Using the Token

For protected endpoints, include the token in your request headers:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Default Users (After Seeding)

| Email                | Password     | Role    |
|----------------------|--------------|---------|
| admin@example.com    | admin123     | admin   |
| john@example.com     | password123  | user    |
| jane@example.com     | password123  | user    |
| bob@example.com      | password123  | manager |

### Security Notes

- The JWT secret is configured in `config.toml` under the `[jwt]` section.
- **Important**: Change the default JWT secret in production!
- Token expiry time is configurable via `token_expiry_hours` in `config.toml` (default: 24 hours).
- Passwords are hashed using bcrypt before storage.
- Never commit your `config.toml` file with production secrets to version control.
- Consider using environment variables or secure vaults for sensitive production configurations.

## Logging

The application includes a custom middleware that logs every `POST`, `PUT`, and `DELETE` request. 

- **Location**: Logs are saved in the `log/` directory.
- **Rotation**: A new log file is created for each day (e.g., `2026-01-04.log`).
- **Content**: Each log entry includes:
  - Timestamp
  - HTTP Method and Path
  - Status Code
  - Latency (processing time)
  - Client IP
  - URL Query Parameters
  - Request Payload (JSON body)

Example log entry:
`15:04:05 [ACTION] PUT /api/products/3 | Status: 200 | Latency: 32.1934ms | IP: 127.0.0.1 | Query: none | Payload: {"name": "Updated Product", "price": 149.99}`



