# Go_CRM_Backend

## Description
This project is a Customer Relationship Management (CRM) Backend Server built with Go, Gorilla Mux, and PostgreSQL. The server provides endpoints for managing customer data through CRUD operations.

## Features
### API Endpoints
- **GET /customers**: Retrieve all customers.
- **GET /customers/{id}**: Retrieve a customer by ID.
- **POST /customers**: Add a new customer.
- **PUT /customers/{id}**: Update an existing customer.
- **DELETE /customers/{id}**: Delete a customer.

## Installation

1. **Clone the repository**:
    ```bash
    git clone https://github.com/yourusername/project-name.git
    cd project-name
    ```

2. **Install dependencies**:
   - **Initialize a Go module** (if no `go.mod` file exists):
     ```bash
     go mod init <module-name>
     ```
     > **Note**: Replace `<module-name>` with the name of your project, e.g., `go mod init main`. Ensure `go.mod` is in the root directory of the project.

   - **Update dependencies**:
     ```bash
     go mod tidy
     ```

   - **Import the Gorilla Mux package**:
     Make sure to include the following import statement in your code:
     ```go
     import "github.com/gorilla/mux"
     ```

## Usage
Configure the server at `http://localhost:3000` by including the following in func main():
```go
http.ListenAndServe(":3000", router)
```
Start the server by running:
```bash
go run main.go
```
