# Go RESTful Products API

This is a RESTful API built with Go, Gin, and GORM for managing products. It supports CRUD operations, filtering, searching, and pagination.

## Features

- **CRUD Operations**: Create, Read, Update, and Delete products.
- **Filtering**: Advanced filtering with operators, functions, and JSON-based filters.
- **Search**: Global search across specified fields.
- **Pagination**: Paginated responses with total counts.
- **Sorting**: Sort results by any field in ascending or descending order.
- **Relation Loading**: Eager load related data.

## Project Structure

```text
.gitignore
go.mod
main.go
database/
  seed/
    main/
      main.go         # Seeder entry point
    product_seeder.go # Initial data seeder
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
restful/
  controller.go     # Generic controller logic
  interface.go
  scopes.go
routes/
  api.go            # Main route entry point
  product.go        # Product-specific routes
```

### Key Files

- **main.go**: Entry point of the application.
- **models/product.go**: Defines the Product model.
- **restful/**: Contains reusable components for controllers, filters, and database operations.
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

### API Endpoints

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

